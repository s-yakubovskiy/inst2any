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

type StoryWorker struct {
	cfg        *config.Config
	database   *sql.DB
	gcsClient  *storage.GCS
	metaClient *instagram.Client
	vkClient   *vk.Client
}

func NewStoryWorker(cfg *config.Config, database *sql.DB, gcsClient *storage.GCS, metaClient *instagram.Client, vkClient *vk.Client) *StoryWorker {
	return &StoryWorker{
		cfg:        cfg,
		database:   database,
		gcsClient:  gcsClient,
		metaClient: metaClient,
		vkClient:   vkClient,
	}
}

func (m *StoryWorker) Name() string {
	return "[worker:story]"
}

func (m *StoryWorker) Enabled() bool {
	return m.cfg.Workers.Instagram.Story.Enabled
}

func (m *StoryWorker) Work(ctx context.Context) {
	log.Println(m.Name(), "StoryWorker run")
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

func (m *StoryWorker) processMedia(ctx context.Context) {
	// Fetch story ids
	ids, err := m.metaClient.FetchMediaIds("stories")
	if err != nil {
		log.Printf("%s Failed to fetch story ids: %v", m.Name(), err)
		return
	}

	// For each id, fetch the story details
	for _, id := range ids {
		m.syncStory(ctx, id)
	}
}

func (m *StoryWorker) syncStory(ctx context.Context, id string) {
	synced, err := db.CheckAndInsert(id, "stories", m.database)
	if err != nil {
		log.Printf("%s Failed to check and insert story id: %v", m.Name(), err)
		return
	}

	if synced {
		log.Printf("%s (db) %+v %s\n", m.Name(), id, "is already synced.")
		return
	}

	media, err := m.metaClient.FetchMediaDetail(id)
	// fmt.Printf("[insta] %+v | id %+v\n", story.MediaURL, id)
	if err != nil {
		log.Printf("[worker:story]: Failed to fetch story details: %v", err)
		return
	}

	// download current story to storyReader with retry 3
	storyReader, err := downloader.DownloadFile(media.MediaURL)
	if err != nil {
		fmt.Println(m.Name(), "Error downloading file:", err)
		return
	}

	// Upload the story to GCS
	err = m.gcsClient.Upload(ctx, "stories", id, storyReader)

	if err != nil {
		log.Printf("%s Failed to upload to GCS: %v", m.Name(), err)
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
		log.Printf("%s Failed to vk upload: %+v\n", m.Name(), err)
		return
	}

	// If story is successfully uploaded, update the story record as synced in the database
	err = db.MarkAsSynced(id, "stories", m.database)
	if err != nil {
		log.Printf("%s (db) Failed to update story as synced: %v", m.Name(), err)
		return
	}

	log.Printf("%s Successfully transferred & synced story id: %s\n", m.Name(), id)
}
