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

type MediaWorker struct {
	cfg        *config.Config
	database   *sql.DB
	gcsClient  *storage.GCS
	metaClient *instagram.Client
	vkClient   *vk.Client
}

func NewMediaWorker(cfg *config.Config, database *sql.DB, gcsClient *storage.GCS, metaClient *instagram.Client, vkClient *vk.Client) *MediaWorker {
	return &MediaWorker{
		cfg:        cfg,
		database:   database,
		gcsClient:  gcsClient,
		metaClient: metaClient,
		vkClient:   vkClient,
	}
}

func (m *MediaWorker) Name() string {
	return "[worker:media]"
}

func (m *MediaWorker) Enabled() bool {
	return m.cfg.Workers.Instagram.Post.Enabled
}

func (m *MediaWorker) Work(ctx context.Context) {
	log.Println(m.Name(), "MediaWorker run")
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

func (m *MediaWorker) processMedia(ctx context.Context) {
	// Fetch media ids
	ids, err := m.metaClient.FetchMediaIds("media")
	if err != nil {
		log.Printf("%s Failed to fetch media ids: %v", m.Name(), err)
		return
	}

	// For each id, fetch the media details
	for _, id := range ids {
		m.syncMedia(ctx, id)
	}
}

func (m *MediaWorker) syncMedia(ctx context.Context, id string) {
	synced, err := db.CheckAndInsert(id, "media", m.database)
	if err != nil {
		log.Printf("%s Failed to check and insert media id: %v", m.Name(), err)
		return
	}

	if synced {
		log.Printf("%s (db) %+v %s\n", m.Name(), id, "is already synced.")
		return
	}

	media, err := m.metaClient.FetchMediaDetail(id)
	// fmt.Printf("[insta] %+v | id %+v\n", media.MediaURL, id)
	if err != nil {
		log.Printf("%s Failed to fetch media details: %v", m.Name(), err)
		return
	}

	// download current media to mediaReader with retry 3
	mediaReader, err := downloader.DownloadFile(media.MediaURL)
	if err != nil {
		fmt.Println(m.Name(), "Error downloading file:", err)
		return
	}

	// Upload the media to GCS
	err = m.gcsClient.Upload(ctx, "posts", id, mediaReader)

	if err != nil {
		log.Printf("%s Failed to upload to GCS: %v", m.Name(), err)
		return
	}

	// Upload to VK
	urlDL := m.gcsClient.ReturnPublicURL(ctx, "posts", id)

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
		err = m.vkClient.UploadWallPhoto(media.Caption, media.Caption, resp.Body)
	} else {
		err = m.vkClient.UploadVideo(media.Caption, media.Caption, resp.Body)
	}
	if err != nil {
		log.Printf("%s Failed to vk upload: %+v\n", m.Name(), err)
		return
	}

	// If media is successfully uploaded, update the media record as synced in the database
	err = db.MarkAsSynced(id, "media", m.database)
	if err != nil {
		log.Printf("%s (db) Failed to update media as synced: %v", m.Name(), err)
		return
	}

	log.Printf("%s Successfully transferred & synced media id: %s\n", m.Name(), id)
}
