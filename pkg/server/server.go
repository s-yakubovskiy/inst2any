package server

import (
	"net/http"

	"github.com/s-yakubovskiy/inst2vk/pkg/vk"
)

type UploadVideoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FileURL     string `json:"file_url"` // updated
}

func NewServer(videoService *vk.VideoService) *http.ServeMux {
	mux := http.NewServeMux()

	uploadHandler := http.HandlerFunc(videoService.UploadHandler)
	mux.Handle("/upload", AuthMiddleware(LoggingMiddleware(uploadHandler)))

	return mux
}
