package consumer

import (
	"atlas-drg/retry"
	"atlas-drg/topic"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

type Consumer struct {
	l                 *log.Logger
	ctx               context.Context
	groupId           string
	topicToken        string
	emptyEventCreator EmptyEventCreator
	h                 EventProcessor
}

func NewConsumer(l *log.Logger, ctx context.Context, h EventProcessor, options ...ConsumerOption) Consumer {
	c := &Consumer{}
	c.l = l
	c.ctx = ctx
	c.h = h
	for _, option := range options {
		option(c)
	}
	return *c
}

type EmptyEventCreator func() interface{}

type EventProcessor func(*log.Logger, interface{})

type ConsumerOption func(c *Consumer)

func SetGroupId(groupId string) func(c *Consumer) {
	return func(c *Consumer) {
		c.groupId = groupId
	}
}

func SetTopicToken(topicToken string) func(c *Consumer) {
	return func(c *Consumer) {
		c.topicToken = topicToken
	}
}

func SetEmptyEventCreator(f EmptyEventCreator) func(c *Consumer) {
	return func(c *Consumer) {
		c.emptyEventCreator = f
	}
}

func (c Consumer) Init() {
	td, err := topic.TopicRequests(c.l).GetTopic(c.topicToken)
	if err != nil {
		c.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
	}

	c.l.Printf("[INFO] creating topic consumer for %s", td.Attributes.Name)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
		Topic:   td.Attributes.Name,
		GroupID: c.groupId,
		MaxWait: 500 * time.Millisecond,
	})

	readMessage := func(attempt int) (bool, interface{}, error) {
		msg, err := r.ReadMessage(c.ctx)
		if err != nil {
			c.l.Printf("[WARN] could not successfully read message on topic %s, will retry", td.Attributes.Name)
			return true, nil, err
		}
		return false, &msg, err
	}

	for {
		msg, err := retry.RetryResponse(readMessage, 10)
		if err != nil {
			c.l.Fatalf("[ERROR] could not successfully read message " + err.Error())
		}
		if val, ok := msg.(*kafka.Message); ok {
			event := c.emptyEventCreator()
			err = json.Unmarshal(val.Value, &event)
			if err != nil {
				c.l.Println("[ERROR] could not unmarshal event into event class ", val.Value)
			} else {
				c.h(c.l, event)
			}
		}
	}
}
