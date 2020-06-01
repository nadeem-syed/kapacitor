package sns

import (
	"fmt"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/influxdata/kapacitor/alert"
	"github.com/influxdata/kapacitor/keyvalue"
)

type Diagnostic interface {
	WithContext(ctx ...keyvalue.T) Diagnostic
	Error(msg string, err error)
}

type Service struct {
	configValue atomic.Value
	diag        Diagnostic
	client      Client
}

func NewService(c Config, d Diagnostic) *Service {
	s := &Service{
		diag: d,
	}
	s.configValue.Store(c)
	clientConfig, _ := c.ClientConfig()
	newClient, err := new(clientConfig)
	s.client = newClient
	if err != nil {
		s = nil
	}
	return s
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) config() Config {
	return s.configValue.Load().(Config)
}

func (s *Service) Update(newConfig []interface{}) error {
	if l := len(newConfig); l != 1 {
		return fmt.Errorf("expected only one new config object, got %d", l)
	}
	if c, ok := newConfig[0].(Config); !ok {
		return fmt.Errorf("expected config object to be of type %T, got %T", c, newConfig[0])
	} else {
		s.configValue.Store(c)
	}
	return nil
}

func (s *Service) Global() bool {
	c := s.config()
	return c.Global
}

func (s *Service) StateChangesOnly() bool {
	c := s.config()
	return c.StateChangesOnly
}

type testOptions struct {
	TopicARN string      `json:"topicarn"`
	Message  string      `json:"message"`
	Level    alert.Level `json:"level"`
}

func (s *Service) TestOptions() interface{} {
	c := s.config()
	return &testOptions{
		TopicARN: c.TopicARN,
		Message:  "test sns message",
		Level:    alert.Critical,
	}
}

func (s *Service) Test(options interface{}) error {
	o, ok := options.(*testOptions)
	if !ok {
		return fmt.Errorf("unexpected options type %T", options)
	}
	return s.Alert(o.TopicARN, o.Message, o.Level)
}

func (s *Service) Alert(topicArn string, message string, level alert.Level) error {

	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicArn),
	}

	err := s.client.Publish(input)
	if err != nil {
		fmt.Println("Publish error:", err)
		return err
	}

	return nil
}

type HandlerConfig struct {
	// SNS topic in which to post messages.
	// If empty uses the topic from the configuration.
	TopicArn string `mapstructure:"topicarn"`
}

type handler struct {
	s    *Service
	c    HandlerConfig
	diag Diagnostic
}

func (s *Service) Handler(c HandlerConfig, ctx ...keyvalue.T) alert.Handler {
	return &handler{
		s:    s,
		c:    c,
		diag: s.diag.WithContext(ctx...),
	}
}

func (h *handler) Handle(event alert.Event) {
	if err := h.s.Alert(
		h.c.TopicArn,
		event.State.Message,
		event.State.Level,
	); err != nil {
		h.diag.Error("failed to send event to SNS", err)
	}
}
