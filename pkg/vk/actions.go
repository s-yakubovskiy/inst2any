package vk

import (
	"io"
	"os"

	"github.com/SevereCloud/vksdk/v2/api/params"
)

type VideoParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FileURL     string `json:"file_url"`
}

func GetFileReader(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (c *Client) UploadVideo(name, description string, file io.Reader) error {
	p := params.NewVideoSaveBuilder()
	p.Repeat(true)
	p.Name(name)
	p.Description(description)
	p.Confirm(true)
	p.Wallpost(true)
	p.Compression(true)

	_, err := c.vk.UploadVideo(p.Params, file)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UploadWallPhoto(name, description string, file io.Reader) error {
	_, err := c.vk.UploadWallPhoto(file)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UploadStoryVideo(file io.Reader) error {
	p := params.NewStoriesGetVideoUploadServerBuilder()
	p.AddToNews(true)

	_, err := c.vk.UploadStoriesVideo(p.Params, file)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UploadStoryPhoto(file io.Reader) error {
	p := params.NewStoriesGetPhotoUploadServerBuilder()
	p.AddToNews(true)

	_, err := c.vk.UploadStoriesPhoto(p.Params, file)
	if err != nil {
		return err
	}

	return nil
}
