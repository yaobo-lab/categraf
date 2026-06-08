package emqx

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"flashcat.cloud/categraf/config"
	"flashcat.cloud/categraf/inputs"
	"flashcat.cloud/categraf/types"
	"github.com/guonaihong/gout"
	"github.com/mitchellh/mapstructure"
)

const inputPrefixName = "emqx"

type EMQX struct {
	config.PluginConfig
	Instances []*Instance `toml:"instances"`
}

func (a *EMQX) Clone() inputs.Input {
	return &EMQX{}
}

func (a *EMQX) Name() string {
	return inputPrefixName
}

func init() {
	inputs.Add("emqx", func() inputs.Input {
		return &EMQX{}
	})
}

type Instance struct {
	config.InstanceConfig
	Api      string `toml:"api"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	NodeName string `toml:"nodename"`
}

// 多实例
func (m *EMQX) GetInstances() []inputs.Instance {
	ret := make([]inputs.Instance, len(m.Instances))
	for i := 0; i < len(m.Instances); i++ {
		ret[i] = m.Instances[i]
	}
	return ret
}

func (ins *Instance) Gather(slist *types.SampleList) {

	err := ins.getNodeStatus(slist)
	if err != nil {
		return
	}

	ins.getMeterics(slist)
}

func (ins *Instance) getHttpResult(url string) (*Result, error) {
	rsp := &Result{}

	err := gout.
		GET(url).
		SetBasicAuth(ins.User, ins.Password).
		SetTimeout(10 * time.Second).
		BindJSON(rsp).
		Do()

	if err == nil && rsp.Code != 0 {
		log.Println("emqx err_code 接口错误码:", rsp.Code)
		return nil, nil
	}

	return rsp, err
}

// 获取节点状态
func (ins *Instance) getNodeStatus(slist *types.SampleList) error {

	url := fmt.Sprintf("%s/nodes/%s", ins.Api, ins.NodeName)

	rsp, err := ins.getHttpResult(url)

	if err != nil {
		slist.PushSample(inputPrefixName, "healthz", 0)
		log.Println("emqx err 接口出错:", err)
		return err
	}

	if rsp == nil {
		return fmt.Errorf("rsp is null")
	}

	var node_status NodeStatusDto
	err = mapstructure.WeakDecode(rsp.Data, &node_status)
	if err != nil {
		log.Println("E! failed cover data to NodeStatusDto:", err)
		return err
	}

	//节点挂了
	if node_status.Error == "nodedown" {
		slist.PushSample(inputPrefixName, "healthz", 0)
		return fmt.Errorf("nodedown")
	}

	slist.PushSample(inputPrefixName, "healthz", 1)

	//387 days, 3 hours, 30 minutes, 42 seconds
	uptime := stringToUptime(node_status.Uptime)
	slist.PushSample(inputPrefixName, "uptime", uptime)

	//内存使用
	memory_used := stringToMemoryUse(node_status.MemoryUsed)
	slist.PushSample(inputPrefixName, "memory_used", memory_used)

	//当前接入客户端数量
	slist.PushSample(inputPrefixName, "client_total", node_status.Connections)
	//节点进程使用量
	slist.PushSample(inputPrefixName, "process_used", node_status.ProcessUsed)

	isRuning := 0
	if node_status.NodeStatus == "Running" {
		isRuning = 1
	}
	slist.PushSample(inputPrefixName, "runing", isRuning)
	return nil
}

// 获取集群下所有统计指标数据
func (ins *Instance) getMeterics(slist *types.SampleList) {

	url := fmt.Sprintf("%s/metrics", ins.Api)

	rsp, err := ins.getHttpResult(url)

	if err != nil {
		log.Println("E! failed cover data to NodeStatusDto:", err)
		return
	}

	if rsp == nil {
		return
	}

	var nodes []MetricsDto
	err = mapstructure.Decode(rsp.Data, &nodes)
	if err != nil {
		log.Println("E! failed cover data to NodeStatusDto:", err)
		return
	}

	for _, node := range nodes {
		if node.Node != ins.NodeName {
			continue
		}

		item := node.Metrics
		//发送时丢弃的消息总数
		slist.PushSample(inputPrefixName, "meteric_delivery_dropped", item.DeliveryDropped)
		//发送时由于消息队列满而被丢弃的 QoS 为 0 的消息数量
		slist.PushSample(inputPrefixName, "meteric_delivery_dropped_qos0_msg", item.DeliveryDroppedQos0Msg)
		//发送时由于长度超过限制而被丢弃的消息数量
		slist.PushSample(inputPrefixName, "meteric_delivery_dropped_too_large", item.DeliveryDroppedTooLarge)
		//客户端断开连接次数
		slist.PushSample(inputPrefixName, "meteric_client_disconnected", item.ExhookDefaultClientDisconnected)
		//发送时由于消息过期而被丢弃的消息数量
		slist.PushSample(inputPrefixName, "meteric_delivery_dropped_expired", item.DeliveryDroppedExpired)
	}
}

// 将字符串转数字类型
func stringToMemoryUse(str string) float64 {

	str = strings.ReplaceAll(str, " ", "")

	var total float64 = 0
	if strings.Contains(str, "G") {
		str = strings.ReplaceAll(str, "G", "")
		size, err := strconv.ParseFloat(str, 64)
		if err == nil {
			total = size
		}
	}
	return total
}

func stringToUptime(f string) int {
	f = strings.ReplaceAll(f, " ", "")
	values := strings.Split(f, ",")

	total := 0
	for _, value := range values {

		if strings.Contains(value, "days") {
			vstr := strings.ReplaceAll(value, "days", "")
			v, err := strconv.Atoi(vstr)
			if err == nil {
				total += v * 86400
			}
		}

		if strings.Contains(value, "hours") {
			vstr := strings.ReplaceAll(value, "hours", "")
			v, err := strconv.Atoi(vstr)
			if err == nil {
				total += v * 3600
			}
		}

		if strings.Contains(value, "minutes") {
			vstr := strings.ReplaceAll(value, "minutes", "")
			v, err := strconv.Atoi(vstr)
			if err == nil {
				total += v * 60
			}
		}

		if strings.Contains(value, "seconds") {
			vstr := strings.ReplaceAll(value, "seconds", "")
			v, err := strconv.Atoi(vstr)
			if err == nil {
				total += v
			}
		}
	}
	return total
}
