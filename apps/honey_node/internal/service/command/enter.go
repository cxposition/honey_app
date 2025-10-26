package command

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/info"
	"honey_node/internal/utils/ip"
	"io"
	"os"
	"time"
)

// ============================
// 节点注册逻辑
// ============================

func Register() error {
	if global.Config.System.Uid == "" {
		global.Config.System.Uid = uuid.NewString()
		core.SetConfig(global.Config)
	}

	_ip, mac, err := ip.GetNetworkInfo(global.Config.System.Network)
	if err != nil {
		return err
	}

	hostName, _ := os.Hostname()
	sysInfo, err := info.GetSystemInfo()
	if err != nil {
		return err
	}

	logrus.Infof("🖥️ 系统信息: %+v", sysInfo)

	var networkList []*node_rpc.NetworkInfoMessage
	list, err := info.GetNetworkList("hy-")
	if err != nil {
		return err
	}
	for _, n := range list {
		networkList = append(networkList, &node_rpc.NetworkInfoMessage{
			Network: n.Network,
			Ip:      n.Ip,
			Net:     n.Net,
			Mask:    int32(n.Mask),
		})
	}

	req := &node_rpc.RegisterRequest{
		Ip:      _ip,
		Mac:     mac,
		NodeUid: global.Config.System.Uid,
		Version: global.Version,
		Commit:  global.Commit,
		SystemInfo: &node_rpc.SystemInfoMessage{
			HostName:            hostName,
			DistributionVersion: sysInfo.Distribution,
			CoreVersion:         sysInfo.KernelVersion,
			SystemType:          sysInfo.Arch,
			StartTime:           sysInfo.BootTime.Format(time.DateTime),
		},
		NetworkList: networkList,
	}

	if _, err := global.GrpcClient.Register(context.Background(), req); err != nil {
		return err
	}

	logrus.Info("✅ 节点注册成功")
	return nil
}

// ============================
// Command Stream 主循环（带重连机制）
// ============================

func RunCommandLoop() {
	var backoff = time.Second * 2

	for {
		// 每次重连前都重新初始化 gRPC 客户端
		global.GrpcClient = core.GetGrpcClient()

		// 节点重新注册，防止管理端重启后丢失注册信息
		if err := Register(); err != nil {
			logrus.Errorf("节点注册失败：%v", err)
			time.Sleep(backoff)
			continue
		}

		// 启动命令会话
		if err := StartCommandSession(); err != nil {
			logrus.Errorf("command 会话出错：%v", err)
		}

		logrus.Warnf("command 会话断开，%v 秒后重连...", backoff.Seconds())
		time.Sleep(backoff)

		// 指数回退，最大 1 分钟
		if backoff < time.Minute {
			backoff *= 2
		}
	}
}

// ============================
// 单次 Command 会话逻辑
// ============================

func StartCommandSession() error {
	ctx := metadata.NewOutgoingContext(context.Background(),
		metadata.Pairs("NodeID", global.Config.System.Uid))

	stream, err := global.GrpcClient.Command(ctx)
	if err != nil {
		return err
	}
	logrus.Infof("✅ command 通道建立成功")

	respChan := make(chan *node_rpc.CmdResponse, 16)
	defer close(respChan)

	// 异步发送响应协程
	done := make(chan struct{})
	go func() {
		defer close(done)
		for resp := range respChan {
			if err := stream.Send(resp); err != nil {
				logrus.Errorf("❌ 发送到服务端失败: %v", err)
				return
			}
		}
	}()

	// 心跳协程，维持长连接
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				resp := &node_rpc.CmdResponse{
					CmdType: node_rpc.CmdType_cmdPingType,
					NodeID:  global.Config.System.Uid,
				}
				select {
				case respChan <- resp:
				default:
					logrus.Warn("⚠️ 心跳通道阻塞，丢弃一次心跳包")
				}
			}
		}
	}()

	// 接收命令主循环
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			logrus.Info("🔌 服务器主动关闭 command 流")
			break
		}
		if err != nil {
			logrus.Errorf("command 流接收错误: %v", err)
			break
		}

		logrus.Infof("📩 收到命令: %+v", req)

		switch req.CmdType {
		case node_rpc.CmdType_cmdNetworkFlushType:
			HandleNetworkFlush(req, respChan)
		default:
			logrus.Warnf("⚠️ 未知命令类型: %v", req.CmdType)
		}
	}

	// 等待发送协程退出
	<-done
	logrus.Info("🔁 command 会话退出")
	return nil
}

// ============================
// 命令处理逻辑
// ============================

func HandleNetworkFlush(req *node_rpc.CmdRequest, respChan chan *node_rpc.CmdResponse) {
	if len(req.NetworkFlushInMessage.FilterNetworkName) == 0 {
		logrus.Warn("⚠️ 命令参数为空")
		return
	}
	name := req.NetworkFlushInMessage.FilterNetworkName[0]

	list, err := info.GetNetworkList(name)
	if err != nil {
		logrus.Errorf("获取网络信息失败：%v", err)
		return
	}

	var outList []*node_rpc.NetworkInfoMessage
	for _, n := range list {
		outList = append(outList, &node_rpc.NetworkInfoMessage{
			Network: n.Network,
			Ip:      n.Ip,
			Net:     n.Net,
			Mask:    int32(n.Mask),
		})
	}

	resp := &node_rpc.CmdResponse{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  req.TaskID,
		NodeID:  global.Config.System.Uid,
		NetworkFlushOutMessage: &node_rpc.NetworkFlushOutMessage{
			NetworkList: outList,
		},
	}

	select {
	case respChan <- resp:
	default:
		logrus.Warn("⚠️ 响应通道已满，丢弃响应")
	}
}
