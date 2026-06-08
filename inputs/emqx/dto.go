package emqx

type Result struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
}

// 获取节点状态
// curl -i --basic -u admin:public -X GET "http://10.0.1.7:8081/api/v4/nodes/emqx@10.0.1.5"
// {
//   "data": {
//     "version": "4.4.11",
//     "uptime": "28 days, 21 hours, 55 minutes, 29 seconds",
//     "process_used": 670,
//     "process_available": 2097152,
//     "otp_release": "24.3.4.2/12.3.2.2",
//     "node_status": "Running",
//     "node": "emqx@10.0.1.5",
//     "memory_used": "22.50G",
//     "memory_total": "23.20G",
//     "max_fds": 1048576,
//     "load5": "0.66",
//     "load15": "0.62",
//     "load1": "0.93",
//     "connections": 25
//   },
//   "code": 0
// }
// {"data":{"node":"emqx@10.0.1.5","error":"nodedown"},"code":0}

// 节点状态
// https://www.emqx.io/docs/zh/v4.4/advanced/http-api.html#%E8%8A%82%E7%82%B9
type NodeStatusDto struct {
	//EMQX 版本
	Version string `mapstructure:"version" json:"version"`
	//EMQX 运行时间
	Uptime string `mapstructure:"uptime" json:"uptime"`
	//已占用的进程数量
	ProcessUsed      int    `mapstructure:"process_used" json:"process_used"`
	ProcessAvailable int    `mapstructure:"process_available" json:"process_available"`
	OtpRelease       string `mapstructure:"otp_release" json:"otp_release"`
	//Running
	NodeStatus string `mapstructure:"node_status" json:"node_status"`
	//节点名称
	Node string `mapstructure:"node" json:"node"`
	//VM 已占用的内存大小
	MemoryUsed string `mapstructure:"memory_used" json:"memory_used"`
	//总内存
	MemoryTotal string  `mapstructure:"memory_total" json:"memory_total"`
	MaxFds      int     `mapstructure:"max_fds" json:"max_fds"`
	Load5       float64 `mapstructure:"load5" json:"load5"`
	Load15      float64 `mapstructure:"load15" json:"load15"`
	Load1       float64 `mapstructure:"load1" json:"load1"`
	//当前接入此节点的客户端数量
	Connections int64  `mapstructure:"connections" json:"connections"`
	Error       string `mapstructure:"error" json:"error"`
}

// 指标
// https://www.emqx.io/docs/zh/v4.4/advanced/http-api.html#%E7%BB%9F%E8%AE%A1%E6%8C%87%E6%A0%87
type MetricsDto struct {
	Node    string     `json:"node"`
	Metrics MetricItem `json:"metrics"`
}

type MetricItem struct {
	PacketsPubrecMissed   int `json:"packets.pubrec.missed"`
	PacketsPubrelSent     int `json:"packets.pubrel.sent"`
	PacketsPubackReceived int `json:"packets.puback.received"`
	MessagesQos0Received  int `json:"messages.qos0.received"`
	PacketsUnsubackSent   int `json:"packets.unsuback.sent"`
	//发送时丢弃的消息总数
	DeliveryDropped  int `json:"delivery.dropped"`
	SessionCreated   int `json:"session.created"`
	MessagesQos2Sent int `json:"messages.qos2.sent"`
	//发送时由于消息队列满而被丢弃的 QoS 为 0 的消息数量
	DeliveryDroppedQos0Msg int `json:"delivery.dropped.qos0_msg"`
	ClientACLCacheHit      int `json:"client.acl.cache_hit"`
	ClientDisconnected     int `json:"client.disconnected"`
	PacketsPubrecInuse     int `json:"packets.pubrec.inuse"`
	ClientACLAllow         int `json:"client.acl.allow"`
	PacketsAuthSent        int `json:"packets.auth.sent"`
	BytesReceived          int `json:"bytes.received"`
	//发送时由于长度超过限制而被丢弃的消息数量
	DeliveryDroppedTooLarge     int `json:"delivery.dropped.too_large"`
	PacketsPublishReceived      int `json:"packets.publish.received"`
	MessagesDropped             int `json:"messages.dropped"`
	MessagesQos2Received        int `json:"messages.qos2.received"`
	PacketsPublishInuse         int `json:"packets.publish.inuse"`
	ClientConnect               int `json:"client.connect"`
	PacketsSubscribeError       int `json:"packets.subscribe.error"`
	ClientAuthSuccess           int `json:"client.auth.success"`
	PacketsSubscribeAuthError   int `json:"packets.subscribe.auth_error"`
	BytesSent                   int `json:"bytes.sent"`
	SessionTakeovered           int `json:"session.takeovered"`
	PacketsUnsubscribeError     int `json:"packets.unsubscribe.error"`
	MessagesDelivered           int `json:"messages.delivered"`
	MessagesReceived            int `json:"messages.received"`
	ClientACLDeny               int `json:"client.acl.deny"`
	PacketsUnsubscribeReceived  int `json:"packets.unsubscribe.received"`
	ExhookDefaultMessagePublish int `json:"exhook.default.message.publish"`
	SessionResumed              int `json:"session.resumed"`
	PacketsSubscribeReceived    int `json:"packets.subscribe.received"`
	ClientConnected             int `json:"client.connected"`
	PacketsPubcompMissed        int `json:"packets.pubcomp.missed"`
	PacketsPublishError         int `json:"packets.publish.error"`
	PacketsPubackSent           int `json:"packets.puback.sent"`
	//超出接收限制而被丢弃的消息数量
	PacketsPublishDropped        int `json:"packets.publish.dropped"`
	PacketsDisconnectReceived    int `json:"packets.disconnect.received"`
	ClientUnsubscribe            int `json:"client.unsubscribe"`
	PacketsPubcompSent           int `json:"packets.pubcomp.sent"`
	ClientCheckACL               int `json:"client.check_acl"`
	MessagesDroppedNoSubscribers int `json:"messages.dropped.no_subscribers"`
	//客户端断开连接次数
	ExhookDefaultClientDisconnected   int `json:"exhook.default.client.disconnected"`
	PacketsAuthReceived               int `json:"packets.auth.received"`
	MessagesDroppedAwaitPubrelTimeout int `json:"messages.dropped.await_pubrel_timeout"`
	//接收的认证失败的 CONNECT 报文数量
	PacketsConnackAuthError int `json:"packets.connack.auth_error"`
	PacketsDisconnectSent   int `json:"packets.disconnect.sent"`
	PacketsPubrelMissed     int `json:"packets.pubrel.missed"`
	PacketsPubcompInuse     int `json:"packets.pubcomp.inuse"`
	PacketsPubrecReceived   int `json:"packets.pubrec.received"`
	MessagesPublish         int `json:"messages.publish"`
	PacketsPingreqReceived  int `json:"packets.pingreq.received"`
	MessagesForward         int `json:"messages.forward"`
	//发送时由于 No Local 订阅选项而被丢弃的消息数量
	DeliveryDroppedNoLocal int `json:"delivery.dropped.no_local"`
	ClientAuthenticate     int `json:"client.authenticate"`
	MessagesQos1Sent       int `json:"messages.qos1.sent"`
	PacketsPubackMissed    int `json:"packets.puback.missed"`
	PacketsPingrespSent    int `json:"packets.pingresp.sent"`
	MessagesSent           int `json:"messages.sent"`
	PacketsReceived        int `json:"packets.received"`
	//发送时由于消息过期而被丢弃的消息数量
	DeliveryDroppedExpired int `json:"delivery.dropped.expired"`
	ClientSubscribe        int `json:"client.subscribe"`
	SessionDiscarded       int `json:"session.discarded"`
	PacketsSubackSent      int `json:"packets.suback.sent"`
	PacketsPubcompReceived int `json:"packets.pubcomp.received"`
	//发送时由于消息队列满而被丢弃的 QoS 不为 0 的消息数量
	DeliveryDroppedQueueFull   int `json:"delivery.dropped.queue_full"`
	MessagesQos1Received       int `json:"messages.qos1.received"`
	ClientAuthSuccessAnonymous int `json:"client.auth.success.anonymous"`
	//EMQX 存储的延迟发布的消息数量
	MessagesDelayed              int `json:"messages.delayed"`
	PacketsPubrelReceived        int `json:"packets.pubrel.received"`
	PacketsConnectReceived       int `json:"packets.connect.received"`
	MessagesAcked                int `json:"messages.acked"`
	ExhookDefaultClientConnected int `json:"exhook.default.client.connected"`
	PacketsPublishAuthError      int `json:"packets.publish.auth_error"`
	ClientConnack                int `json:"client.connack"`
	PacketsConnackSent           int `json:"packets.connack.sent"`
	PacketsPubackInuse           int `json:"packets.puback.inuse"`
	SessionTerminated            int `json:"session.terminated"`
	ClientAuthFailure            int `json:"client.auth.failure"`
	PacketsSent                  int `json:"packets.sent"`
	PacketsPublishSent           int `json:"packets.publish.sent"`
	PacketsConnackError          int `json:"packets.connack.error"`
	MessagesQos0Sent             int `json:"messages.qos0.sent"`
	MessagesRetained             int `json:"messages.retained"`
	PacketsPubrecSent            int `json:"packets.pubrec.sent"`
}

// 获取节点状态
// # curl -i --basic -u admin:public -X GET "http://10.0.1.7:8081/api/v4/metrics"
// {
//   "data": [
//     {
//       "node": "emqx@10.0.1.5",
//       "metrics": {
//         "packets.pubrec.missed": 0,
//         "packets.pubrel.sent": 0,
//         "packets.puback.received": 1327,
//         "messages.qos0.received": 114548,
//         "packets.unsuback.sent": 193,
//         "delivery.dropped": 0,
//         "session.created": 636,
//         "messages.qos2.sent": 0,
//         "delivery.dropped.qos0_msg": 0,
//         "client.acl.cache_hit": 117322,
//         "client.disconnected": 1224,
//         "packets.pubrec.inuse": 0,
//         "client.acl.allow": 249459,
//         "packets.auth.sent": 0,
//         "bytes.received": 50435792,
//         "delivery.dropped.too_large": 0,
//         "packets.publish.received": 204327,
//         "messages.dropped": 199551,
//         "messages.qos2.received": 0,
//         "packets.publish.inuse": 0,
//         "client.connect": 2138,
//         "packets.subscribe.error": 0,
//         "client.auth.success": 2130,
//         "packets.subscribe.auth_error": 0,
//         "bytes.sent": 11825505,
//         "session.takeovered": 1487,
//         "packets.unsubscribe.error": 0,
//         "messages.delivered": 7913,
//         "messages.received": 204327,
//         "client.acl.deny": 30117,
//         "packets.unsubscribe.received": 193,
//         "exhook.default.message.publish": 204684,
//         "session.resumed": 1494,
//         "packets.subscribe.received": 79076,
//         "client.connected": 2130,
//         "packets.pubcomp.missed": 0,
//         "packets.publish.error": 0,
//         "packets.puback.sent": 89779,
//         "packets.publish.dropped": 0,
//         "packets.disconnect.received": 522,
//         "client.unsubscribe": 193,
//         "packets.pubcomp.sent": 0,
//         "client.check_acl": 162254,
//         "messages.dropped.no_subscribers": 199551,
//         "exhook.default.client.disconnected": 1224,
//         "packets.auth.received": 0,
//         "messages.dropped.await_pubrel_timeout": 0,
//         "packets.connack.auth_error": 0,
//         "packets.disconnect.sent": 1260,
//         "packets.pubrel.missed": 0,
//         "packets.pubcomp.inuse": 0,
//         "packets.pubrec.received": 0,
//         "messages.publish": 204684,
//         "packets.pingreq.received": 1157304,
//         "messages.forward": 1967,
//         "delivery.dropped.no_local": 0,
//         "client.authenticate": 2138,
//         "messages.qos1.sent": 1331,
//         "packets.puback.missed": 0,
//         "packets.pingresp.sent": 1157304,
//         "messages.sent": 7913,
//         "packets.received": 1444890,
//         "delivery.dropped.expired": 0,
//         "client.subscribe": 79076,
//         "session.discarded": 19,
//         "packets.suback.sent": 79076,
//         "packets.pubcomp.received": 0,
//         "delivery.dropped.queue_full": 0,
//         "messages.qos1.received": 89779,
//         "client.auth.success.anonymous": 0,
//         "messages.delayed": 0,
//         "packets.pubrel.received": 0,
//         "packets.connect.received": 2138,
//         "messages.acked": 1327,
//         "exhook.default.client.connected": 2130,
//         "packets.publish.auth_error": 0,
//         "client.connack": 2138,
//         "packets.connack.sent": 644,
//         "packets.puback.inuse": 0,
//         "session.terminated": 600,
//         "client.auth.failure": 8,
//         "packets.sent": 1337663,
//         "packets.publish.sent": 7913,
//         "packets.connack.error": 8,
//         "messages.qos0.sent": 6582,
//         "messages.retained": 130427,
//         "packets.pubrec.sent": 0
//       }
//     },
//     {
//       "node": "emqx@10.0.1.7",
//       "metrics": {
//         "packets.pubrec.missed": 0,
//         "packets.pubrel.sent": 1,
//         "packets.puback.received": 1759,
//         "messages.qos0.received": 127794,
//         "packets.unsuback.sent": 78,
//         "delivery.dropped": 0,
//         "session.created": 2550,
//         "messages.qos2.sent": 1,
//         "delivery.dropped.qos0_msg": 0,
//         "client.acl.cache_hit": 940835,
//         "client.disconnected": 3691,
//         "packets.pubrec.inuse": 0,
//         "client.acl.allow": 926566,
//         "packets.auth.sent": 0,
//         "bytes.received": 143468005,
//         "delivery.dropped.too_large": 0,
//         "packets.publish.received": 441289,
//         "messages.dropped": 439157,
//         "messages.qos2.received": 55,
//         "packets.publish.inuse": 0,
//         "client.connect": 4416,
//         "packets.subscribe.error": 0,
//         "client.auth.success": 4411,
//         "packets.subscribe.auth_error": 0,
//         "bytes.sent": 25276428,
//         "session.takeovered": 1874,
//         "packets.unsubscribe.error": 0,
//         "messages.delivered": 7061,
//         "messages.received": 448504,
//         "client.acl.deny": 372239,
//         "packets.unsubscribe.received": 78,
//         "exhook.default.message.publish": 448846,
//         "session.resumed": 1867,
//         "packets.subscribe.received": 857440,
//         "client.connected": 4409,
//         "packets.pubcomp.missed": 0,
//         "packets.publish.error": 8,
//         "packets.puback.sent": 320655,
//         "packets.publish.dropped": 0,
//         "packets.disconnect.received": 2339,
//         "client.unsubscribe": 78,
//         "packets.pubcomp.sent": 55,
//         "client.check_acl": 357970,
//         "messages.dropped.no_subscribers": 439157,
//         "exhook.default.client.disconnected": 3691,
//         "packets.auth.received": 0,
//         "messages.dropped.await_pubrel_timeout": 0,
//         "packets.connack.auth_error": 5,
//         "packets.disconnect.sent": 936,
//         "packets.pubrel.missed": 0,
//         "packets.pubcomp.inuse": 0,
//         "packets.pubrec.received": 1,
//         "messages.publish": 448846,
//         "packets.pingreq.received": 1722839,
//         "messages.forward": 2525,
//         "delivery.dropped.no_local": 0,
//         "client.authenticate": 4416,
//         "messages.qos1.sent": 1786,
//         "packets.puback.missed": 0,
//         "packets.pingresp.sent": 1722839,
//         "messages.sent": 7061,
//         "packets.received": 3030221,
//         "delivery.dropped.expired": 0,
//         "client.subscribe": 857436,
//         "session.discarded": 358,
//         "packets.suback.sent": 857436,
//         "packets.pubcomp.received": 1,
//         "delivery.dropped.queue_full": 0,
//         "messages.qos1.received": 320655,
//         "client.auth.success.anonymous": 0,
//         "messages.delayed": 0,
//         "packets.pubrel.received": 55,
//         "packets.connect.received": 4416,
//         "messages.acked": 1760,
//         "exhook.default.client.connected": 4409,
//         "packets.publish.auth_error": 8,
//         "client.connack": 4416,
//         "packets.connack.sent": 2549,
//         "packets.puback.inuse": 0,
//         "session.terminated": 2154,
//         "client.auth.failure": 5,
//         "packets.sent": 2913532,
//         "packets.publish.sent": 7061,
//         "packets.connack.error": 7,
//         "messages.qos0.sent": 5274,
//         "messages.retained": 126397,
//         "packets.pubrec.sent": 55
//       }
//     }
//   ],
//   "code": 0
// }
