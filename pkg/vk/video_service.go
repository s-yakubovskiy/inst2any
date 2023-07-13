package vk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type VideoService struct {
	Client *Client
}

func NewVideoService(client *Client) *VideoService {
	return &VideoService{Client: client}
}

func (s *VideoService) UploadHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params VideoParams
	err := decoder.Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", params)

	resp, err := http.Get(params.FileURL)
	if err != nil {
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to download file, status: %s", resp.Status), http.StatusInternalServerError)
		return
	}

	err = s.Client.UploadVideo(params.Name, params.Description, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Video uploaded successfully")
}
