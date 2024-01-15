package instagram

import (
	"net/http"
	"os"

	"github.com/s-yakubovskiy/inst2any/pkg/config"
)

type MediaFetcher interface {
	FetchMediaIds() ([]string, error)
	FetchMediaDetail(id string) (MediaDetail, error)
}

type MediaResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

type MediaDetail struct {
	MediaURL  string `json:"media_url"`
	Caption   string `json:"caption"`
	ID        string `json:"id"`
	MediaType string `json:"media_type"`
	Permalink string `json:"permalink"`
}

type Client struct {
	httpClient   *http.Client
	limitStories int
	limitPosts   int
	token        string
	api          string
	id           string
}

func NewClient(config config.InstagramConfig) *Client {
	// token should be passed through env and fallback to config.yaml
	token := os.Getenv("META_INSTAGRAM_TOKEN")
	if token == "" {
		token = config.AccessToken
	}

	// set default fallback values
	if config.LastStoriesCount == 0 {
		config.LastStoriesCount = 5
	}
	if config.LastPostsCount == 0 {
		config.LastPostsCount = 3
	}

	return &Client{
		httpClient:   &http.Client{},
		token:        token,
		api:          config.API,
		id:           config.AccountID,
		limitStories: config.LastStoriesCount,
		limitPosts:   config.LastPostsCount,
	}
}
