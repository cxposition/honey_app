package suricata_service

import (
	"encoding/json"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"honey_node/internal/global"
	"io"
)

type AlertType struct {
	Timestamp string `json:"timestamp"`
	FlowId    int64  `json:"flow_id"`
	InIface   string `json:"in_iface"`
	EventType string `json:"event_type"`
	SrcIp     string `json:"src_ip"`
	SrcPort   int    `json:"src_port"`
	DestIp    string `json:"dest_ip"`
	DestPort  int    `json:"dest_port"`
	Proto     string `json:"proto"`
	PktSrc    string `json:"pkt_src"`
	TxId      int    `json:"tx_id"`
	Alert     struct {
		Action      string `json:"action"`
		Gid         int    `json:"gid"`
		SignatureId int    `json:"signature_id"`
		Rev         int    `json:"rev"`
		Signature   string `json:"signature"`
		Category    string `json:"category"`
		Severity    int    `json:"severity"`
		Metadata    struct {
			Level []string `json:"level"`
		} `json:"metadata"`
	} `json:"alert"`
	Http struct {
		Hostname         string `json:"hostname"`
		Url              string `json:"url"`
		HttpUserAgent    string `json:"http_user_agent"`
		HttpContentType  string `json:"http_content_type"`
		HttpMethod       string `json:"http_method"`
		Protocol         string `json:"protocol"`
		Status           int    `json:"status"`
		Length           int    `json:"length"`
		HttpResponseBody string `json:"http_response_body"`
	} `json:"http"`
	AppProto  string `json:"app_proto"`
	Direction string `json:"direction"`
	Flow      struct {
		PktsToserver  int    `json:"pkts_toserver"`
		PktsToclient  int    `json:"pkts_toclient"`
		BytesToserver int    `json:"bytes_toserver"`
		BytesToclient int    `json:"bytes_toclient"`
		Start         string `json:"start"`
		SrcIp         string `json:"src_ip"`
		DestIp        string `json:"dest_ip"`
		SrcPort       int    `json:"src_port"`
		DestPort      int    `json:"dest_port"`
	} `json:"flow"`
	Payload string `json:"payload"`
	Stream  int    `json:"stream"`
}

func Run() {
	t, err := tail.TailFile(global.Config.System.EvePath,
		tail.Config{
			Follow: true,
			Location: &tail.SeekInfo{
				Offset: 0,
				Whence: io.SeekEnd, // 从最新的地方开始监听
			},
		},
	)
	if err != nil {
		logrus.Fatalf("suricata路径错误: %s", err)
	}

	logrus.Infof("开始监听suricata告警日志")
	for line := range t.Lines {
		var t AlertType
		err = json.Unmarshal([]byte(line.Text), &t)
		if err != nil {
			logrus.Errorf("解析suricata告警记录失败 %s %s", err, line.Text)
			continue
		}
		if t.EventType != "alert" {
			continue
		}
		logrus.Infof("%s %s => %s:%d", t.Alert.Signature, t.SrcIp, t.DestIp, t.DestPort)
		// 发送到mq
	}
}
