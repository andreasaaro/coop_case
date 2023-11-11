package config

import (
	"fmt"
	"github.com/rickb777/date/period"
	"os"
	"strconv"
	"time"
)

// Configs for mastodon client and kafka producer, added some functions that enables default values and handling of int, bool.
type MastodonConfig struct {
	Url          string
	Limit        string
	RetryCount   int
	PollInterval time.Duration
}

type KafkaConfig struct {
	Brokers           string
	Topic             string
	PrintKafkaMessage bool
}

type Config struct {
	MastodonConfig MastodonConfig
	KafkaConfig    KafkaConfig
}

func NewConfig() *Config {
	return &Config{
		MastodonConfig{
			Url:          getEnvIfSet("BASE_URL", "https://mastodon.social/api/v1/timelines/public"),
			Limit:        getEnvIfSet("REQUEST_LIMIT", "40"),
			RetryCount:   getEnvIntIfSet("RETRY_COUNT", 10),
			PollInterval: getEnvDurationFromIsoPeriodIfSet("POLL_INTERVAL", 10*time.Second),
		},
		KafkaConfig{
			Brokers:           getEnvIfSet("KAFKA_BROKERS", "localhost:9094"),
			Topic:             getEnvIfSet("KAFKA_TOPIC", "mastodon_topic"),
			PrintKafkaMessage: getEnvBoolIfSet("PRINT_KAFKA_MESSAGE", false),
		},
	}
}

func getEnvIfSet(envVar, defaultValue string) string {
	if os.Getenv(envVar) != "" {
		return os.Getenv(envVar)
	}
	return defaultValue
}

func getEnvIntIfSet(envVar string, defaultValue int) int {
	str := getEnvIfSet(envVar, strconv.Itoa(defaultValue))
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse %s: %v", envVar, err))
	}
	return val
}

// GetEnvBoolIfSet returns a bool based on envVar was parseable as bool,
// defaultValue if not.
func getEnvBoolIfSet(envVar string, defaultValue bool) bool {
	env := os.Getenv(envVar)
	b, err := strconv.ParseBool(env)
	if err != nil {
		return defaultValue
	}
	return b
}

// GetEnvDurationFromIsoPeriodIfSet returns a time.Duration from a parsed iso period. e.g. "PT1M" = 1 * time.Minute
// defultValue if not, long periods might be inaccurate
func getEnvDurationFromIsoPeriodIfSet(key string, defaultValue time.Duration) time.Duration {
	p, err := period.Parse(os.Getenv(key))
	if err == nil {
		d, _ := p.Duration()
		return d
	}
	return defaultValue
}
