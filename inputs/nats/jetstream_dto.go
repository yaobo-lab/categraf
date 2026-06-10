package nats

import "time"

type JetStreamDto struct {
	ServerID        string           `json:"server_id"`
	Now             time.Time        `json:"now"`
	Config          Config           `json:"config"`
	Memory          int              `json:"memory"`
	Storage         int              `json:"storage"`
	ReservedMemory  int              `json:"reserved_memory"`
	ReservedStorage int              `json:"reserved_storage"`
	Accounts        int              `json:"accounts"`
	HaAssets        int              `json:"ha_assets"`
	API             API              `json:"api"`
	Streams         int              `json:"streams"`
	Consumers       int              `json:"consumers"`
	Messages        uint64           `json:"messages"`
	Bytes           uint64           `json:"bytes"`
	MetaCluster     MetaCluster      `json:"meta_cluster"`
	AccountDetails  []AccountDetails `json:"account_details"`
}

type Config struct {
	MaxMemory  int    `json:"max_memory"`
	MaxStorage uint64 `json:"max_storage"`
	StoreDir   string `json:"store_dir"`
	CompressOk bool   `json:"compress_ok"`
}

type API struct {
	Total  int `json:"total"`
	Errors int `json:"errors"`
}
type MetaCluster struct {
	Name        string `json:"name"`
	Leader      string `json:"leader"`
	Peer        string `json:"peer"`
	ClusterSize int    `json:"cluster_size"`
}
type Cluster struct {
	Name   string `json:"name"`
	Leader string `json:"leader"`
}
type State struct {
	Messages      uint64    `json:"messages"`
	Bytes         uint64    `json:"bytes"`
	FirstSeq      uint64    `json:"first_seq"`
	FirstTs       time.Time `json:"first_ts"`
	LastSeq       uint64    `json:"last_seq"`
	LastTs        time.Time `json:"last_ts"`
	ConsumerCount int       `json:"consumer_count"`
}
type Delivered struct {
	ConsumerSeq int `json:"consumer_seq"`
	StreamSeq   int `json:"stream_seq"`
}
type AckFloor struct {
	ConsumerSeq int `json:"consumer_seq"`
	StreamSeq   int `json:"stream_seq"`
}
type ConsumerDetail struct {
	StreamName     string    `json:"stream_name"`
	Name           string    `json:"name"`
	Created        time.Time `json:"created"`
	Delivered      Delivered `json:"delivered"`
	AckFloor       AckFloor  `json:"ack_floor"`
	NumAckPending  int       `json:"num_ack_pending"`
	NumRedelivered int       `json:"num_redelivered"`
	NumWaiting     int       `json:"num_waiting"`
	NumPending     int       `json:"num_pending"`
	Cluster        Cluster   `json:"cluster"`
}
type StreamDetail struct {
	Name           string           `json:"name"`
	Created        time.Time        `json:"created"`
	Cluster        Cluster          `json:"cluster"`
	State          State            `json:"state"`
	ConsumerDetail []ConsumerDetail `json:"consumer_detail,omitempty"`
}
type AccountDetails struct {
	Name            string         `json:"name"`
	ID              string         `json:"id"`
	Memory          int            `json:"memory"`
	Storage         int            `json:"storage"`
	ReservedMemory  uint64         `json:"reserved_memory"`
	ReservedStorage uint64         `json:"reserved_storage"`
	Accounts        int            `json:"accounts"`
	HaAssets        int            `json:"ha_assets"`
	API             API            `json:"api"`
	StreamDetail    []StreamDetail `json:"stream_detail"`
}
