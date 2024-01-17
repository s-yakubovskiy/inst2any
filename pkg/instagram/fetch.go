package instagram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) getLimits(field string) int {
	if field == "media" {
		return c.limitPosts
	}
	if field == "stories" {
		return c.limitStories
	}
	return 3
}

func (c *Client) FetchMediaIds(field string) ([]string, error) {
	limiter := c.getLimits(field)
	url := fmt.Sprintf("%s/%s/%s?access_token=%s", c.api, c.id, field, c.token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var media MediaResponse
	json.Unmarshal(body, &media)

	var ids []string
	for i, m := range media.Data {
		if i >= limiter {
			break
		}
		ids = append(ids, m.ID)
	}

	return reverseSlice(ids), nil
}

func (c *Client) FetchMediaDetail(id string) (*MediaDetail, error) {
	url := fmt.Sprintf("%s/%s?fields=media_url,caption,id,media_type,permalink&access_token=%s", c.api, id, c.token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var mediaDetail MediaDetail
	json.Unmarshal(body, &mediaDetail)
	if mediaDetail.MediaURL == "" {
		return nil, fmt.Errorf("No media url for id: %+v\n", id)
	}

	return &mediaDetail, nil
}
