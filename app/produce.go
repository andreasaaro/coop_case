package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"coop_case/config"
	"coop_case/mastodon"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

// Consume requests, loop over blog messages, parse content html string and send single blog posts as kafka messages
func produce(ctx context.Context, sourceCh <-chan []*mastodon.MastodonData, sinkCh chan<- *sarama.ProducerMessage, cfg config.KafkaConfig) error {
	logrus.Infoln("Starting produce")
	defer logrus.Infoln("Exiting produce")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-sourceCh:
			if !ok {
				return nil
			}

			for _, value := range msg {
				blogText, err := extractParagraphs(value.Content)
				if err != nil {
					fmt.Println("Error:", err)
				}

				if blogText == "" {
					continue
				}
				finalMessage := fmt.Sprintf("Micro blog id: %s published by user: %s created at: %s: %s", value.ID, value.Account.Username, value.CreatedAt, blogText)

				if cfg.PrintKafkaMessage {
					fmt.Println(finalMessage)
				}

				saramaMsg := &sarama.ProducerMessage{
					Topic: cfg.Topic,
					Value: sarama.ByteEncoder(finalMessage),
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				case sinkCh <- saramaMsg:
				}

				time.Sleep(500000 * time.Microsecond)

			}
		}
	}
}

// Parse blog content html, extract only paragraphs and write to string.
func extractParagraphs(input string) (string, error) {
	reader := strings.NewReader(input)

	doc, err := html.Parse(reader)
	if err != nil {
		return "", fmt.Errorf("unable to parse html blog")
	}

	var paragraphs string
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "p" {
			// Extract the text content of the paragraph
			blogTextContent := extractTextContent(n)
			paragraphs += blogTextContent
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}

	extractText(doc)

	finalBlogText := strings.TrimSpace(paragraphs)

	return finalBlogText, nil
}

func extractTextContent(n *html.Node) string {
	var result string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			result += c.Data
		}
	}
	return result
}
