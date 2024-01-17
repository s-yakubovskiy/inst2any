package daemon

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/s-yakubovskiy/inst2any/pkg/config"
	"github.com/s-yakubovskiy/inst2any/pkg/db"
	"github.com/s-yakubovskiy/inst2any/pkg/downloader"
	"github.com/s-yakubovskiy/inst2any/pkg/instagram"
	"github.com/s-yakubovskiy/inst2any/pkg/storage"
	"github.com/s-yakubovskiy/inst2any/pkg/vk"
)

type VKStoryWorker struct {
	cfg        *config.Config
	database   *sql.DB
	gcsClient  *storage.GCS
	metaClient *instagram.Client
	vkClient   *vk.Client
}

func NewVKStoryWorker(cfg *config.Config, database *sql.DB, gcsClient *storage.GCS, metaClient *instagram.Client, vkClient *vk.Client) *VKStoryWorker {
	return &VKStoryWorker{
		cfg:        cfg,
		database:   database,
		gcsClient:  gcsClient,
		metaClient: metaClient,
		vkClient:   vkClient,
	}
}

func (m *VKStoryWorker) Name() string {
	return "VK"
}

func (m *VKStoryWorker) FullName() string {
	return "vk:worker:story"
}

func (m *VKStoryWorker) Enabled() bool {
	return m.cfg.Workers.Vkontakte.Story.Enabled
}

func (m *VKStoryWorker) Work(ctx context.Context) {
	log.Println(m.FullName(), "run")
	for {
		select {
		case <-ctx.Done():
			// Context was cancelled, stop the worker
			return
		default:
			m.processMedia(ctx)
			// Sleep for the configured duration before checking for new story
			select {
			case <-time.After(time.Duration(m.cfg.SleepInterval) * time.Second):
			case <-ctx.Done():
				// If context is cancelled, stop sleeping and return
				return
			}
			break
		}
	}
}

func (m *VKStoryWorker) processMedia(ctx context.Context) {
	// Fetch story ids
	ids, err := m.metaClient.FetchMediaIds("stories")
	if err != nil {
		log.Printf("[%s] Failed to fetch story ids: %v", m.FullName(), err)
		return
	}

	// For each id, fetch the story details
	for _, id := range ids {
		m.syncStory(ctx, id)
	}
}

func (m *VKStoryWorker) syncStory(ctx context.Context, id string) {
	synced, err := db.CheckAndInsert(id, "stories", m.Name(), m.database)
	if err != nil {
		log.Printf("[%s] Failed to check and insert story id: %v", m.FullName(), err)
		return
	}

	if synced {
		log.Printf("[%s] (db) %+v %s\n", m.FullName(), id, "is already synced.")
		return
	}

	media, err := m.metaClient.FetchMediaDetail(id)
	// fmt.Printf("[insta] %+v | id %+v\n", story.MediaURL, id)
	if err != nil {
		log.Printf("[%s] Failed to fetch story details: %v", m.FullName(), err)
		return
	}

	// download current story to storyReader with retry 3
	storyReader, err := downloader.DownloadFile(media.MediaURL)
	if err != nil {
		fmt.Println("[", m.FullName(), "]", "Error downloading file:", err)
		return
	}

	// Upload the story to GCS
	err = m.gcsClient.Upload(ctx, "stories", id, storyReader)

	if err != nil {
		log.Printf("[%s] Failed to upload to GCS: %v", m.FullName(), err)
		return
	}

	// Upload to VK
	urlDL := m.gcsClient.ReturnPublicURL(ctx, "stories", id)

	resp, err := http.Get(urlDL)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		return
	}

	if media.MediaType == "IMAGE" {
		err = m.vkClient.UploadStoryPhoto(resp.Body)
	} else {
		err = m.vkClient.UploadStoryVideo(resp.Body)
	}
	if err != nil {
		log.Printf("[%s] Failed to vk upload: %+v\n", m.FullName(), err)
		return
	}

	// If story is successfully uploaded, update the story record as synced in the database
	err = db.MarkAsSynced(id, "stories", m.Name(), m.database)
	if err != nil {
		log.Printf("[%s] (db) Failed to update story as synced: %v", m.FullName(), err)
		return
	}

	log.Printf("[%s] Successfully transferred & synced story id: %s\n", m.FullName(), id)
}
