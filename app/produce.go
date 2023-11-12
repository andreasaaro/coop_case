package app

import (
	"context"
	"fmt"
	"strings"

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

			// Blog messages arrive in ascending order, start from end to send oldest messages first
			for m := len(msg) - 1; m >= 0; m-- {

				blogText, err := extractParagraphs(msg[m].Content)
				if err != nil {
					fmt.Println("Error:", err)
				}

				if blogText == "" {
					continue
				}

				// Concatenate a string with blog info that will be sent to kafka
				finalMessage := fmt.Sprintf("Micro-blog id: %s published by user: %s created at: %s: %s", msg[m].ID, msg[m].Account.Username, msg[m].CreatedAt, blogText)

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
			}

			logrus.Info("Micro-blogs sent to kafka")
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
