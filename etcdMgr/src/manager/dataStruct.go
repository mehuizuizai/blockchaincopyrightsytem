package manager

import (
	//"common/etcd"
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

type Member struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	PeerURLs   []string `json:"peerURLs"`
	ClientURLs []string `json:"clientURLs"`
}

type request struct {
	PeerURLs []string
}

//type certificate struct {
//	Identifier4Peer      string
//	IP4Peer              string
//	IP4Admin             string
//	Identifier4Admin     string
//	PublicKey4Admin      string
//	Port4Admin           string
//	ETCDPeerPort4Admin   string
//	ETCDClientPort4Admin string
//	ClusterMembers       []etcd.MemberInfo
//}
