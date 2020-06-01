package sns

import (
	"github.com/pkg/errors"
)

type Config struct {
	// Whether sns integration is enabled.
	Enabled bool `toml:"enabled" override:"enabled"`
	// The SNS Topic ARN you want to publush to.
	TopicARN string `toml:"topicarn" override:"topicarn"`
	// Secret ID for AWS access
	SecretID string `toml:"secretid" override:"secretid"`
	// Secret key for AWS access
	SecretKey string `toml:"secretket" override:"secretkey,redact"`
	// AWS region where topic is located
	Region string `toml:"region" override:"region"`
	// Whether all alerts should automatically post to SNS
	Global bool `toml:"global" override:"global"`
	// Whether all alerts should automatically use stateChangesOnly mode.
	// Only applies if global is also set.
	StateChangesOnly bool `toml:"state-changes-only" override:"state-changes-only"`
}

func NewConfig() Config {
	return Config{}
}

func (c Config) Validate() error {
	if c.Enabled && c.TopicARN == "" {
		return errors.New("must specify url")
	}
	return nil
}

func (c Config) ClientConfig() (SNSConfig, error) {
	return SNSConfig{
		AccessKey: c.SecretKey,
		SecretKey: c.SecretID,
		Region:    c.Region,
	}, nil
}
