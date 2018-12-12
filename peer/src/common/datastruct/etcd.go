package datastruct

import (
	"time"
)

type Event struct {
	Action    string      `json:"action"`
	Node      *NodeExtern `json:"node,omitempty"`
	PrevNode  *NodeExtern `json:"prevNode,omitempty"`
	EtcdIndex uint64      `json:"-"`
	Refresh   bool        `json:"refresh,omitempty"`
}

type NodeExtern struct {
	Key           string      `json:"key,omitempty"`
	Value         *string     `json:"value,omitempty"`
	Dir           bool        `json:"dir,omitempty"`
	Expiration    *time.Time  `json:"expiration,omitempty"`
	TTL           int64       `json:"ttl,omitempty"`
	Nodes         NodeExterns `json:"nodes,omitempty"`
	ModifiedIndex uint64      `json:"modifiedIndex,omitempty"`
	CreatedIndex  uint64      `json:"createdIndex,omitempty"`
}

type NodeExterns []*NodeExtern
