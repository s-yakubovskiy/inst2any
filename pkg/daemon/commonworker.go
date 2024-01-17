package daemon

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/s-yakubovskiy/inst2any/pkg/db"
// 	"github.com/s-yakubovskiy/inst2any/pkg/downloader"
// )

// func (m *VKPostWorker) processMedia(ctx context.Context) {
// 	// Fetch media ids
// 	ids, err := m.metaClient.FetchMediaIds("media")
// 	if err != nil {
// 		log.Printf("[%s] Failed to fetch media ids: %v", m.FullName(), err)
// 		return
// 	}

// 	// For each id, fetch the media details
// 	for _, id := range ids {
// 		m.syncMedia(ctx, id)
// 	}
// }

// func (m *VKPostWorker) syncMedia(ctx context.Context, id string) {
// 	synced, err := db.CheckAndInsert(id, "media", m.Name(), m.database)
// 	if err != nil {
// 		log.Printf("[%s] Failed to check and insert media id: %v", m.FullName(), err)
// 		return
// 	}

// 	if synced {
// 		log.Printf("[%s] (db) %+v %s\n", m.FullName(), id, "is already synced.")
// 		return
// 	}

// 	media, err := m.metaClient.FetchMediaDetail(id)
// 	// fmt.Printf("[insta] %+v | id %+v\n", media.MediaURL, id)
// 	if err != nil {
// 		log.Printf("[%s] Failed to fetch media details: %v", m.FullName(), err)
// 		return
// 	}

// 	// download current media to mediaReader with retry 3
// 	mediaReader, err := downloader.DownloadFile(media.MediaURL)
// 	if err != nil {
// 		fmt.Println("[", m.FullName(), "]", "Error downloading file:", err)
// 		return
// 	}

// 	// Upload the media to GCS
// 	err = m.gcsClient.Upload(ctx, "posts", id, mediaReader)

// 	if err != nil {
// 		log.Printf("[%s] Failed to upload to GCS: %v", m.FullName(), err)
// 		return
// 	}

// 	// Upload to VK
// 	urlDL := m.gcsClient.ReturnPublicURL(ctx, "posts", id)

// 	resp, err := http.Get(urlDL)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		log.Println(resp.Status)
// 		return
// 	}

// 	if media.MediaType == "IMAGE" {
// 		err = m.vkClient.UploadWallPhoto(media.Caption, media.Caption, resp.Body)
// 	} else {
// 		err = m.vkClient.UploadVideo(media.Caption, media.Caption, resp.Body)
// 	}
// 	if err != nil {
// 		log.Printf("[%s] Failed to vk upload: %+v\n", m.FullName(), err)
// 		return
// 	}

// 	// If media is successfully uploaded, update the media record as synced in the database
// 	err = db.MarkAsSynced(id, "media", m.Name(), m.database)
// 	if err != nil {
// 		log.Printf("[%s] (db) Failed to update media as synced: %v", m.FullName(), err)
// 		return
// 	}

// 	log.Printf("[%s] Successfully transferred & synced media id: %s\n", m.FullName(), id)
// }
