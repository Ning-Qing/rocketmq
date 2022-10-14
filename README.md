# RocketMQ Client

## Get

```
export GOPROXY=http://10.1.30.103:8080
export GONOSUMDB=git.vonechain.com
go get github.com/Ning-Qing/rocketmq
```

## Use
```go
type Message interface {
	Serialization() ([]byte, error)
	Deserialize([]byte) error
}
```
### 实现Message接口
```go

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
```
### 发送消息
```go
    ctx := context.Background()
	addrs := []string{"10.1.40.10:9876", "10.1.40.11:9876"}
	clietn, err := rocketmq.NewRocketMQClient(ctx,
		rocketmq.WithAddress(addrs),
		rocketmq.WithTopic("vchain-cross-invoke-dev"),
		rocketmq.WithGroupName("vbaas-s-dev"),
	)
	if err != nil {
		log.Fatalf("new rocketmq client faild: %s", err.Error())
	}
	msge := &Message{}
	err = clietn.Send(msge)
	if err != nil {
		log.Fatalf("rocketmq client send faild: %s", err.Error())
	}
```