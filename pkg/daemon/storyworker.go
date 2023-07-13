package daemon

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/s-yakubovskiy/inst2vk/pkg/config"
	"github.com/s-yakubovskiy/inst2vk/pkg/db"
	"github.com/s-yakubovskiy/inst2vk/pkg/downloader"
	"github.com/s-yakubovskiy/inst2vk/pkg/instagram"
	"github.com/s-yakubovskiy/inst2vk/pkg/storage"
	"github.com/s-yakubovskiy/inst2vk/pkg/vk"
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

func (m *StoryWorker) Work(ctx context.Context) {
	log.Printf("[worker:story] StoryWorker run")
	for {
		select {
		case <-ctx.Done():
			// Context was cancelled, stop the worker
			return
		default:
			m.processMedia(ctx)
			// Sleep for the configured duration before checking for new media
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

func (d *StoryWorker) processMedia(ctx context.Context) {
	// Fetch media ids
	ids, err := d.metaClient.FetchMediaIds("stories")
	if err != nil {
		log.Printf("[worker:story]: Failed to fetch media ids: %v", err)
		return
	}

	// For each id, fetch the media details
	for _, id := range ids {
		d.syncStory(ctx, id)
	}
}

func (d *StoryWorker) syncStory(ctx context.Context, id string) {
	synced, err := db.CheckAndInsert(id, "stories", d.database)
	if err != nil {
		log.Printf("[worker:story]: Failed to check and insert story id: %v", err)
		return
	}

	if synced {
		log.Printf("[worker:story:db] %+v %s\n", id, "is already synced.")
		return
	}

	media, err := d.metaClient.FetchMediaDetail(id)
	// fmt.Printf("[insta] %+v | id %+v\n", media.MediaURL, id)
	if err != nil {
		log.Printf("[worker:story]: Failed to fetch media details: %v", err)
		return
	}

	// download current media to mediaReader with retry 3
	mediaReader, err := downloader.DownloadFile(media.MediaURL)

	if err != nil {
		fmt.Println("[worker:story]: Error downloading file:", err)
		return
	}

	// Upload the media to GCS
	err = d.gcsClient.Upload(ctx, "stories", id, mediaReader)

	if err != nil {
		log.Printf("[worker:story]: Failed to upload to GCS: %v", err)
		return
	}

	// Upload to VK
	urlDL := d.gcsClient.ReturnPublicURL(ctx, "stories", id)

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
		err = d.vkClient.UploadStoryPhoto(resp.Body)
	} else {
		err = d.vkClient.UploadStoryVideo(resp.Body)
	}
	if err != nil {
		log.Printf("[worker:story] Failed to vk upload: %+v\n", err)
		return
	}

	// If media is successfully uploaded, update the media record as synced in the database
	err = db.MarkAsSynced(id, "stories", d.database)
	if err != nil {
		log.Printf("[worker:story:db]Failed to update story as synced: %v", err)
		return
	}

	log.Printf("[worker:story:inst2vk] Successfully transferred story id: %s\n", id)
}
