package rocketmq

import "errors"

type Option func(*Config) error

func WithAddress(addrs []string) Option {
	return func(c *Config) error {
		if len(addrs) == 0 {
			return errors.New("at least one address is required")
		}
		c.Address = addrs
		return nil
	}
}

func WithCredentials(secretKey, accessKey string) Option {
	return func(c *Config) error {
		c.SecretKey = secretKey
		c.AccessKey = accessKey
		return nil
	}
}

func WithNameSpace(ns string) Option {
	return func(c *Config) error {
		c.NameSpace = ns
		return nil
	}
}

func WithGroupName(name string) Option {
	return func(c *Config) error {
		c.GroupName = name
		return nil
	}
}

func WithTopic(topic string) Option {
	return func(c *Config) error {
		c.Topic = topic
		return nil
	}
}

func WithCache(num int) Option {
	return func(c *Config) error {
		c.MaxCache = num
		return nil
	}
}

func WithStorageFile(path string) Option {
	return func(c *Config) error {
		c.StorageFile = path
		return nil
	}
}

func WithRetry(retries int) Option {
	return func(c *Config) error {
		c.Retry = retries
		return nil
	}
}
