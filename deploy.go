package deploy

import (
	"encoding/json"
	"time"
)

type Action string

const (
	Create   Action = "Create"
	Start    Action = "Start"
	Required Action = "Required"
)

type ResourceType string

const (
	DeliveryPipeline ResourceType = "DeliveryPipeline"
	Release          ResourceType = "Release"
)

// PubSubMessages ...
// See details https://cloud.google.com/deploy/docs/subscribe-deploy-notifications
type PubSubMessage struct {
	AckID   string   `json:"ackId"`
	Message *Message `json:"message"`
}

type Message struct {
	Attributes  *Attributes `json:"attributes"`
	MessageID   string      `json:"messageId"`
	PublishTime time.Time   `json:"publishTime"`
}

type Attributes struct {
	Action       Action       `json:"Action"`
	Resource     string       `json:"Resource"`
	ResourceType ResourceType `json:"ResourceType"`
	Rollout      string       `json:"Rollout"`
}

func ParseEvent(b []byte) (*PubSubMessage, error) {
	var msg PubSubMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
