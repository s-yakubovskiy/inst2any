package vk

import "io"

// Uploader interface defines the contract for uploading a video.
type Uploader interface {
	UploadVideo(name, description string, file io.Reader) error
}
