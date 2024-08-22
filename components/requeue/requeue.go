package requeue

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

const RetriesKey = "requeue_retries"

type Requeue struct {
	config Config
	router *message.Router
}

type Config struct {
	// Subscriber is the subscriber to consume messages from. Required.
	Subscriber message.Subscriber

	// SubscribeTopic is the topic related to the Subscriber to consume messages from. Required.
	SubscribeTopic string

	// Publisher is the publisher to publish requeued messages to. Required.
	Publisher message.Publisher

	// GeneratePublishTopic is the topic related to the Publisher to publish the requeued message to.
	// For example, it could be a constant, or taken from the message's metadata.
	// Required.
	GeneratePublishTopic func(msg *message.Message) (string, error)

	// Delay is the duration to wait before requeueing the message. Optional.
	// The default is no delay.
	//
	// This can be useful to avoid requeueing messages too quickly, for example, to avoid
	// requeueing a message that failed to process due to a temporary issue.
	//
	// Avoid setting this to a very high value, as it will block the message processing.
	Delay time.Duration

	// Router is the custom router to run the requeue handler on. Optional.
	Router *message.Router
}

func (c *Config) setDefaults() {
}

func (c *Config) validate() error {
	if c.Subscriber == nil {
		return errors.New("subscriber is required")
	}

	if c.SubscribeTopic == "" {
		return errors.New("subscribe topic is required")
	}

	if c.Publisher == nil {
		return errors.New("publisher is required")
	}

	if c.GeneratePublishTopic == nil {
		return errors.New("generate publish topic is required")
	}

	return nil
}

func NewRequeue(
	config Config,
	logger watermill.LoggerAdapter,
) (*Requeue, error) {
	config.setDefaults()
	err := config.validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

  router := config.Router
	if router == nil {
		router, err = message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			return nil, fmt.Errorf("could not create router: %w", err)
		}
	}

	r := &Requeue{
		config: config,
		router: router,
	}

	router.AddNoPublisherHandler(
		"requeue",
		config.SubscribeTopic,
		config.Subscriber,
		r.handler,
	)

	return r, nil
}

func (r *Requeue) handler(msg *message.Message) error {
	if r.config.Delay > 0 {
		select {
		case <-msg.Context().Done():
			return msg.Context().Err()
		case <-time.After(r.config.Delay):
		}
	}

	topic, err := r.config.GeneratePublishTopic(msg)
	if err != nil {
		return err
	}

	retriesStr := msg.Metadata.Get(RetriesKey)
	retries, err := strconv.Atoi(retriesStr)
	if err != nil {
		retries = 0
	}

	retries++

	msg.Metadata.Set(RetriesKey, strconv.Itoa(retries))

	err = r.config.Publisher.Publish(topic, msg)
	if err != nil {
		return err
	}

	return nil
}

func (r *Requeue) Run(ctx context.Context) error {
	return r.router.Run(ctx)
}
