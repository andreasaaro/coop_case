package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"coop_case/mastodon"

	"github.com/sirupsen/logrus"
)

// Poll Mastodon api, parse data and send to channel for further processing
func generateSourceStream(ctx context.Context, m mastodon.Mastodon, outChan chan<- []*mastodon.MastodonData) error {
	logrus.Info("Started source")

	ticker := time.NewTicker(time.Second * 5).C

	sinceId := ""

	for {
		select {
		case <-ticker:

			responseBody, err := m.GetBlogPosts(ctx, sinceId)
			if err != nil {
				return fmt.Errorf("unable to get blog posts: %v", err)
			}

			var responseData []*mastodon.MastodonData
			if err := json.Unmarshal(responseBody, &responseData); err != nil {
				return fmt.Errorf("unable to unmarshal mastodon response data %w", err)
			}

			sinceId = responseData[0].ID

			logrus.Infof("Sending %d micro-blogs for processing", len(responseData))

			select {
			case <-ctx.Done():
				return ctx.Err()
			case outChan <- responseData:
			}
		}
	}
}
