package mastodon

import (
	"context"
	"fmt"

	"coop_case/config"

	"github.com/go-resty/resty/v2"
)

// Interface to interact with Mastodon api
type Mastodon interface {
	GetBlogPosts(ctx context.Context, sinceId string) ([]byte, error)
}

type client struct {
	cfg         config.MastodonConfig
	restyClient *resty.Client
}

func NewClient(cfg config.MastodonConfig) (Mastodon, error) {
	c := &client{
		cfg: cfg,
	}
	rClient := resty.New().
		SetBaseURL(cfg.Url).
		SetRetryCount(cfg.RetryCount).
		SetHeader("Content-Type", "application/json")

	c.restyClient = rClient

	return c, nil
}

func (c *client) GetBlogPosts(ctx context.Context, sinceId string) ([]byte, error) {
	baseUrl := c.cfg.Url
	limit := c.cfg.Limit

	return getBlogPosts(ctx, c, baseUrl, limit, sinceId)
}

// Get Mastodon micro blogposts, timeline endpoint always returns newest blog posts. TODO: use min ID to avoid getting same id's multiple times.
func getBlogPosts(ctx context.Context, c *client, baseUrl, limit string, sinceId string) ([]byte, error) {
	restyResp, err := c.restyClient.R().
		SetQueryParams(map[string]string{
			"limit":    limit,
			"sort":     "id",
			"order":    "asc",
			"since_id": sinceId,
		}).
		SetContext(ctx).
		Get(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("get request failed: %w", err)
	}

	if restyResp.StatusCode() != 200 {
		return nil, fmt.Errorf("error when sending get request to mostodon. Got status code %d", restyResp.StatusCode())
	}

	responseBody := restyResp.Body()

	return responseBody, nil
}
