package nats

import "time"

type VarzDto struct {
	ServerID         string       `json:"server_id"`
	ServerName       string       `json:"server_name"`
	Version          string       `json:"version"`
	Proto            int          `json:"proto"`
	GitCommit        string       `json:"git_commit"`
	Go               string       `json:"go"`
	Host             string       `json:"host"`
	Port             int          `json:"port"`
	ConnectUrls      []string     `json:"connect_urls"`
	MaxConnections   int          `json:"max_connections"`
	PingInterval     int64        `json:"ping_interval"`
	PingMax          int          `json:"ping_max"`
	HTTPHost         string       `json:"http_host"`
	HTTPPort         int          `json:"http_port"`
	HTTPBasePath     string       `json:"http_base_path"`
	HTTPSPort        int          `json:"https_port"`
	AuthTimeout      int          `json:"auth_timeout"`
	MaxControlLine   int          `json:"max_control_line"`
	MaxPayload       int          `json:"max_payload"`
	MaxPending       int          `json:"max_pending"`
	Cluster          VarzCluster  `json:"cluster"`
	Gateway          Gateway      `json:"gateway"`
	Leaf             Leaf         `json:"leaf"`
	Mqtt             Mqtt         `json:"mqtt"`
	Websocket        Websocket    `json:"websocket"`
	Jetstream        Jetstream    `json:"jetstream"`
	TLSTimeout       int          `json:"tls_timeout"`
	WriteDeadline    int64        `json:"write_deadline"`
	Start            time.Time    `json:"start"`
	Now              time.Time    `json:"now"`
	Uptime           string       `json:"uptime"`
	Mem              int          `json:"mem"`
	Cores            int          `json:"cores"`
	Gomaxprocs       int          `json:"gomaxprocs"`
	CPU              float64      `json:"cpu"`
	Connections      int          `json:"connections"`
	TotalConnections int          `json:"total_connections"`
	Routes           int          `json:"routes"`
	Remotes          int          `json:"remotes"`
	Leafnodes        int          `json:"leafnodes"`
	InMsgs           int          `json:"in_msgs"`
	OutMsgs          int          `json:"out_msgs"`
	InBytes          int          `json:"in_bytes"`
	OutBytes         int          `json:"out_bytes"`
	SlowConsumers    int          `json:"slow_consumers"`
	Subscriptions    int          `json:"subscriptions"`
	HTTPReqStats     HTTPReqStats `json:"http_req_stats"`
	ConfigLoadTime   time.Time    `json:"config_load_time"`
	SystemAccount    string       `json:"system_account"`
}

type VarzCluster struct {
	Name        string   `json:"name"`
	Addr        string   `json:"addr"`
	ClusterPort int      `json:"cluster_port"`
	AuthTimeout int      `json:"auth_timeout"`
	Urls        []string `json:"urls"`
	TLSTimeout  int      `json:"tls_timeout"`
}
type Gateway struct {
}
type Leaf struct {
}
type Mqtt struct {
}
type Websocket struct {
}
type VarzConfig struct {
	MaxMemory  int    `json:"max_memory"`
	MaxStorage int64  `json:"max_storage"`
	StoreDir   string `json:"store_dir"`
	CompressOk bool   `json:"compress_ok"`
}
type VarzAPI struct {
	Total  int `json:"total"`
	Errors int `json:"errors"`
}
type VarzStats struct {
	Memory          int     `json:"memory"`
	Storage         int     `json:"storage"`
	ReservedMemory  int     `json:"reserved_memory"`
	ReservedStorage int     `json:"reserved_storage"`
	Accounts        int     `json:"accounts"`
	HaAssets        int     `json:"ha_assets"`
	API             VarzAPI `json:"api"`
}
type VarzReplicas struct {
	Name    string `json:"name"`
	Current bool   `json:"current"`
	Offline bool   `json:"offline,omitempty"`
	Active  int    `json:"active"`
	Peer    string `json:"peer"`
}
type VarzMeta struct {
	Name        string         `json:"name"`
	Leader      string         `json:"leader"`
	Peer        string         `json:"peer"`
	Replicas    []VarzReplicas `json:"replicas"`
	ClusterSize int            `json:"cluster_size"`
}
type Jetstream struct {
	Config VarzConfig `json:"config"`
	Stats  VarzStats  `json:"stats"`
	Meta   VarzMeta   `json:"meta"`
}
type HTTPReqStats struct {
	NAMING_FAILED int `json:"/"`
	Jsz           int `json:"/jsz"`
	Varz          int `json:"/varz"`
}
