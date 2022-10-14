package rocketmq_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Ning-Qing/rocketmq"
)

type ReqParam struct {
	Method string `json:"function"`
	Args   []byte `json:"args"`
}

type Message struct {
	Status              int8     `json:"status"`
	ReqParam            ReqParam `json:"reqParam"`
	FromChainID         uint64   `json:"sourceChainId"`
	FromContractAddress string   `json:"sourceContractAddress"`
	FromTxHash          string   `json:"sourceTransactionHash"`
	MidTxHash           string   `json:"relayerTransactionHash"`
	ToChainID           uint64   `json:"targetChainId"`
	ToContractAddress   string   `json:"targetContractAddress"`
	ToTxHash            string   `json:"targetTransactionHash"`
}

func (m *Message) Serialization() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) Deserialize(data []byte) error {
	return json.Unmarshal(data, m)
}

func TestSend(t *testing.T) {
	ctx := context.Background()
	addrs := []string{"10.1.40.10:9876", "10.1.40.11:9876"}
	clietn, err := rocketmq.NewRocketMQClient(ctx,
		rocketmq.WithAddress(addrs),
		rocketmq.WithTopic("vchain-cross-invoke-dev"),
		rocketmq.WithGroupName("vbaas-s-dev"),
	)
	if err != nil {
		t.Fatalf("new rocketmq client faild: %s", err.Error())
	}
	msge := &Message{}
	err = clietn.Send(msge)
	if err != nil {
		t.Fatalf("rocketmq client send faild: %s", err.Error())
	}
}
