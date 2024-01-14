package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Instagram     InstagramConfig `yaml:"instagram"`
	VK            VKConfig        `yaml:"vk"`
	Database      DatabaseConfig  `yaml:"database"`
	GCS           GCSConfig       `yaml:"gcs"`
	SleepInterval int64           `yaml:"sleep_interval"`
	Workers       Workers         `yaml:"workers"`
}

type Workers struct {
	Instagram struct {
		Story struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"story"`
		Post struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"post"`
	} `yaml:"instagram"`
}

type InstagramConfig struct {
	AccessToken      string `yaml:"access_token"`
	AccountID        string `yaml:"account_id"`
	API              string `yaml:"api"`
	LastPostsCount   int    `yaml:"last_posts_count"`
	LastStoriesCount int    `yaml:"last_stories_count"`
}

type VKConfig struct {
	AccessToken string `yaml:"access_token"`
	OwnerID     int    `yaml:"owner_id"`
}

type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

type GCSConfig struct {
	BucketName          string `yaml:"bucket_name"`
	CredentialsFilePath string `yaml:"credentials_file_path"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
