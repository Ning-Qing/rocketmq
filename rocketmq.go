package rocketmq

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"

	mq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type Client struct {
	producer  mq.Producer
	topic     string
	cacheChan chan Message

	fileLock    sync.Mutex
	storageFile *os.File

	cancel context.CancelFunc
	closed int32
}

func NewRocketMQClient(ctx context.Context, opts ...Option) (*Client, error) {
	var err error

	config := newDefaultConfig()
	for _, apply := range opts {
		err = apply(config)
	}
	if err != nil {
		return nil, fmt.Errorf("please check the parameters: %w", err)
	}

	producer, err := mq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver(config.Address)),
		producer.WithCredentials(primitive.Credentials{
			SecretKey: config.SecretKey,
			AccessKey: config.AccessKey,
		}),
		producer.WithGroupName(config.GroupName),
		producer.WithNamespace(config.NameSpace),
		producer.WithRetry(config.Retry),
	)
	if err != nil {
		return nil, fmt.Errorf("new producer faild!, %w", err)
	}

	storageFile, err := os.OpenFile(config.StorageFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open %s faild!, %w", config.StorageFile, err)
	}

	ctx, cancel := context.WithCancel(ctx)
	client := &Client{
		producer:    producer,
		topic:       config.Topic,
		cacheChan:   make(chan Message, config.MaxCache),
		fileLock:    sync.Mutex{},
		storageFile: storageFile,
		cancel:      cancel,
	}

	atomic.StoreInt32(&client.closed, int32(0))
	if err := producer.Start(); err != nil {
		return nil, fmt.Errorf("producer start faild: %w", err)
	}
	go client.run(ctx)
	return client, nil
}

func (c *Client) store(data []byte) (int, error) {
	c.fileLock.Lock()
	defer c.fileLock.Unlock()
	return c.storageFile.Write(data)
}

func (c *Client) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-c.cacheChan:
			var e error
			data, e := msg.Serialization()
			if e != nil {
				continue
			}
			res, e := c.producer.SendSync(context.TODO(), &primitive.Message{
				Topic: c.topic,
				Body:  data,
			})
			if e != nil || res.Status != primitive.SendOK {
				c.store(data)
			}
		}
	}
}

func (c *Client) isClose() bool {
	return atomic.LoadInt32(&c.closed) == int32(1)
}

func (c *Client) Close() {
	if !atomic.CompareAndSwapInt32(&c.closed, 0, int32(1)) {
		return
	}
	c.cancel()
}

func (c *Client) Send(msg Message) error {
	if c.isClose() {
		return errors.New("client is closed")
	}
	c.cacheChan <- msg
	return nil
}
