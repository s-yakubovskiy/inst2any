package vk

import (
	"os"

	"github.com/s-yakubovskiy/inst2any/pkg/config"

	"github.com/SevereCloud/vksdk/v2/api"
)

// Ensure that Client implements the Uploader interface.
var _ Uploader = (*Client)(nil)

type Client struct {
	vk      *api.VK
	token   string
	ownerID int
}

func NewClient(config config.VKConfig) *Client {
	token := os.Getenv("VK_TOKEN")
	if token == "" {
		token = config.AccessToken
	}

	return &Client{
		vk:      api.NewVK(token),
		token:   token,
		ownerID: config.OwnerID,
	}
}
