package rocketmq

type Message interface {
	Serialization() ([]byte, error)
	Deserialize([]byte) error
}
