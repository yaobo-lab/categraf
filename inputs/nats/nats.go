package nats

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"flashcat.cloud/categraf/config"
	"flashcat.cloud/categraf/inputs"
	"flashcat.cloud/categraf/types"
	gnatsd "github.com/nats-io/nats-server/v2/server"
)

//https://docs.nats.io/running-a-nats-service/nats_admin/monitoring#general-information-varz

const inputName = "nats"

type Nats struct {
	config.PluginConfig
	Instances []*Instance `toml:"instances"`
}

func init() {
	inputs.Add(inputName, func() inputs.Input {
		return &Nats{}
	})
}

func (n *Nats) Clone() inputs.Input {
	return &Nats{}
}

func (n *Nats) Name() string {
	return inputName
}

func (n *Nats) GetInstances() []inputs.Instance {
	ret := make([]inputs.Instance, len(n.Instances))
	for i := 0; i < len(n.Instances); i++ {
		ret[i] = n.Instances[i]
	}
	return ret
}

type Instance struct {
	Server          string          `toml:"server"`
	ResponseTimeout config.Duration `toml:"response_timeout"`
	client          *http.Client
	config.HTTPCommonConfig
	config.InstanceConfig
}

func (ins *Instance) Init() error {
	if ins.Server == "" {
		return types.ErrInstancesEmpty
	}
	if ins.ResponseTimeout <= 0 {
		ins.ResponseTimeout = config.Duration(time.Second * 5)
	}

	ins.InitHTTPClientConfig()

	var err error
	ins.client, err = ins.createHTTPClient()
	return err
}

func (ins *Instance) createHTTPClient() (*http.Client, error) {
	tr := &http.Transport{
		ResponseHeaderTimeout: time.Duration(ins.ResponseTimeout),
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(ins.ResponseTimeout),
	}
	return client, nil
}

func (ins *Instance) Gather(slist *types.SampleList) {
	err := ins.getVarz(slist)
	if err != nil {
		return
	}
	ins.getJetStream(slist)
}

// https://demo.nats.io:8222/varz
func (ins *Instance) getVarz(slist *types.SampleList) error {
	if ins.DebugMod {
		log.Println("D! nats... server:", ins.Server)
	}
	address, err := url.Parse(ins.Server)
	if err != nil {
		log.Println("E! error parseURL", err)
		return err
	}
	address.Path = path.Join(address.Path, "varz")

	resp, err := ins.client.Get(address.String())
	if err != nil {
		log.Println("E! error while polling", address.String(), err)
		return err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("E! error reading body", err)
		return err
	}

	stats := new(gnatsd.Varz)
	err = json.Unmarshal(bytes, &stats)
	if err != nil {
		log.Println("E! error parsing response", err)
		return err
	}

	fields := map[string]interface{}{
		"in_msgs":           stats.InMsgs,
		"out_msgs":          stats.OutMsgs,
		"in_bytes":          stats.InBytes,
		"out_bytes":         stats.OutBytes,
		"uptime":            stats.Now.Sub(stats.Start).Nanoseconds(),
		"cores":             stats.Cores,
		"cpu":               stats.CPU,
		"mem":               stats.Mem,
		"connections":       stats.Connections,
		"total_connections": stats.TotalConnections,
		"subscriptions":     stats.Subscriptions,
		"slow_consumers":    stats.SlowConsumers,
		"routes":            stats.Routes,
		"remotes":           stats.Remotes,
		"healthz":           1,
	}
	tags := map[string]string{
		"server": ins.Server,
	}
	slist.PushSamples(inputName, fields, tags)
	return nil
}

// https://docs.nats.io/running-a-nats-service/nats_admin/monitoring#jetstream-information-jsz
// https://demo.nats.io:8222/jsz?consumers=true
func (ins *Instance) getJetStream(slist *types.SampleList) {

	address, err := url.Parse(ins.Server)
	if err != nil {
		log.Println("E! error parseURL", err)
		return
	}
	address.Path = path.Join(address.Path, "jsz?consumers=true")

	resp, err := ins.client.Get(address.String())
	if err != nil {
		log.Println("E! error while polling", address.String(), err)
		return
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("E! error reading body", err)
		return
	}

	stats := new(gnatsd.JSInfo)
	err = json.Unmarshal(bytes, &stats)
	if err != nil {
		log.Println("E! error parsing response", err)
		return
	}

	fields := map[string]interface{}{
		"streams_total":   stats.Streams,   //streams 数量
		"consumers_total": stats.Consumers, //consumers 数量
		"msg_total":       stats.Messages,  //消息数量
	}
	if len(stats.AccountDetails) == 0 {
		return
	}

	slist.PushSamples(inputName, fields)
	for _, item := range stats.AccountDetails[0].Streams {
		tags := map[string]string{"stream_name": item.Name}
		//每个流 ：消息数量
		slist.PushSample(inputName, "stream_msg_count", item.State.Msgs, tags)
		//每个流 ：消费者数量
		slist.PushSample(inputName, "stream_consumer_count", item.State.Consumers, tags)
	}
}
