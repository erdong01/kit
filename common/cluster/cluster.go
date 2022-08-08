package cluster

import (
	"context"
	"errors"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/erDong01/micro-kit/base"

	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/common"
	etv3 "github.com/erDong01/micro-kit/common/cluster/etcdv3"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
	"github.com/erDong01/micro-kit/tools/vector"
	"github.com/nats-io/nats.go"
)

const (
	MAX_CLUSTER_NUM = 128
	CALL_TIME_OUT   = 50 * time.Millisecond
)

type (
	HashClusterMap       map[uint32]*common.ClusterInfo
	HashClusterSocketMap map[uint32]*common.ClusterInfo
	Op                   struct {
		mailBoxEndpoints     []string
		stubMailBoxEndpoints []string
		stub                 common.Stub
	}

	OpOption func(*Op)

	//集群服务器
	Cluster struct {
		actor.Actor
		*Service       //集群注册
		clusterMap     [MAX_CLUSTER_NUM]HashClusterMap
		clusterLocker  [MAX_CLUSTER_NUM]*sync.RWMutex
		hashRing       [MAX_CLUSTER_NUM]*base.HashRing //hash一致性
		conn           *nats.Conn
		dieChan        chan bool
		master         *Master
		clusterInfoMap map[uint32]*common.ClusterInfo
		packetFuncList *vector.Vector //call back
		MailBox        etv3.MailBox
		StubMailBox    etv3.StubMailBox
		Stub           common.Stub
	}

	ICluster interface {
		actor.IActor
		InitCluster(info *common.ClusterInfo, Endpoints []string, natsUrl string)
		RegisterClusterCall() //注册集群通用回调
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc.RpcHead) *common.ClusterInfo
		GetClusters(head rpc.RpcHead) map[uint32]*common.ClusterInfo

		BindPacketFunc(packetFunc network.PacketFunc)
		SendMsg(rpc.RpcHead, string, ...interface{})                    //发送给集群特定服务器
		CallMsg(interface{}, rpc.RpcHead, string, ...interface{}) error //同步给集群特定服务器

		RandomCluster(head rpc.RpcHead) rpc.RpcHead //随机分配
	}

	EmptyClusterInfo struct {
		common.ClusterInfo
	}

	CallFunc struct {
		Func       interface{}
		FuncType   reflect.Type
		FuncVal    reflect.Value
		FuncParams string
	}
)

func (this *EmptyClusterInfo) String() string {
	return ""
}

func (this *Cluster) InitCluster(info *common.ClusterInfo, Endpoints []string, natsUrl string) {
	this.Actor.Init()
	for i := 0; i < MAX_CLUSTER_NUM; i++ {
		this.clusterLocker[i] = &sync.RWMutex{}
		this.clusterMap[i] = make(HashClusterMap)
		this.hashRing[i] = base.NewHashRing()
	}
	//注册服务器
	this.Service = NewService(info, Endpoints)
	this.master = NewMaster(&EmptyClusterInfo{}, Endpoints, &this.Actor)
	this.clusterInfoMap = make(map[uint32]*common.ClusterInfo)
	this.packetFuncList = vector.NewVector()

	conn, err := SetupNatsConn(
		natsUrl,
		this.dieChan,
	)
	if err != nil {
		log.Fatal("nats connect error!!!!", err)
	}
	this.conn = conn

	this.conn.Subscribe(GetChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff: msg.Data})
	})

	this.conn.Subscribe(GetTopicChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff: msg.Data})
	})

	this.conn.Subscribe(GetCallChannel(*info), func(msg *nats.Msg) {
		this.HandlePacket(rpc.Packet{Buff: msg.Data, Reply: msg.Reply})
	})

	rpc.GCall = reflect.ValueOf(this.call)
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

// params[0]:rpc.RpcHead
// params[1]:error
func (this *Cluster) call(parmas ...interface{}) {
	head := *parmas[0].(*rpc.RpcHead)
	reply := head.Reply
	head.Reply = ""
	head.ClusterId = head.SrcClusterId
	if parmas[1] == nil {
		parmas[1] = ""
	} else {
		parmas[1] = parmas[1].(error).Error()
	}
	packet := rpc.Marshal(head, "", parmas[1:]...)
	this.conn.Publish(reply, packet.Buff)
}

func (this *Cluster) AddCluster(info *common.ClusterInfo) {
	this.clusterLocker[info.Type].Lock()
	this.clusterMap[info.Type][info.Id()] = info
	this.clusterLocker[info.Type].Unlock()
	this.hashRing[info.Type].Add(info.IpString())
	log.Printf("服务器[%s:%s:%d]建立连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) DelCluster(info *common.ClusterInfo) {
	this.clusterLocker[info.Type].RLock()
	_, bEx := this.clusterMap[info.Type][info.Id()]
	this.clusterLocker[info.Type].RUnlock()
	if bEx {
		this.clusterLocker[info.Type].Lock()
		delete(this.clusterMap[info.Type], info.Id())
		this.clusterLocker[info.Type].Unlock()
	}
	this.hashRing[info.Type].Remove(info.IpString())
	log.Printf("服务器[%s:%s:%d]断开连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) GetCluster(head rpc.RpcHead) *common.ClusterInfo {
	this.clusterLocker[head.DestServerType].RLock()
	defer this.clusterLocker[head.DestServerType].RUnlock()
	client, bEx := this.clusterMap[head.DestServerType][head.ClusterId]
	if bEx {
		return client
	}
	return nil
}

func (this *Cluster) GetClusters(head rpc.RpcHead) map[uint32]*common.ClusterInfo {
	this.clusterLocker[head.DestServerType].RLock()
	defer this.clusterLocker[head.DestServerType].RUnlock()
	clusters := this.clusterMap[head.DestServerType]
	return clusters
}

func (this *Cluster) BindPacketFunc(callFunc network.PacketFunc) {
	this.packetFuncList.PushBack(callFunc)
}

func (this *Cluster) HandlePacket(packet rpc.Packet) {
	for _, v := range this.packetFuncList.Values() {
		if v.(network.PacketFunc)(packet) {
			break
		}
	}
}
func (this *Cluster) GetBalanceServer(head rpc.RpcHead) *common.ClusterInfo {
	_, head.ClusterId = this.hashRing[head.DestServerType].Get64(head.Id)
	client, bEx := this.clusterMap[head.DestServerType][head.ClusterId]
	if bEx {
		return client
	}
	return nil
}

func (this *Cluster) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	head.SrcClusterId = this.Id()
	this.Send(head, rpc.Marshal(head, funcName, params...))
}

func (this *Cluster) Send(head rpc.RpcHead, packet rpc.Packet) {
	switch head.SendType {
	case rpc.SEND_BALANCE:
		_, head.ClusterId = this.hashRing[head.DestServerType].Get64(head.Id)
		this.conn.Publish(GetRpcChannel(head), packet.Buff)
	case rpc.SEND_POINT:
		this.conn.Publish(GetRpcChannel(head), packet.Buff)
	default:
		this.conn.Publish(GetRpcTopicChannel(head), packet.Buff)
	}
}

func (this *Cluster) CallMsg(cb interface{}, head rpc.RpcHead, funcName string, params ...interface{}) error {
	head.SrcClusterId = this.Id()
	packet := rpc.Marshal(head, funcName, params...)

	switch head.SendType {
	case rpc.SEND_POINT:
	default:
		_, head.ClusterId = this.hashRing[head.DestServerType].Get64(head.Id)
	}

	reply, err := this.conn.Request(GetRpcCallChannel(head), packet.Buff, CALL_TIME_OUT)
	if err == nil {
		rpcPacket, _ := rpc.Unmarshal(reply.Data)
		cf := &CallFunc{Func: cb, FuncVal: reflect.ValueOf(cb), FuncType: reflect.TypeOf(cb), FuncParams: reflect.TypeOf(cb).String()}
		f := cf.FuncVal
		k := cf.FuncType
		err, params := rpc.UnmarshalBodyCall(rpcPacket, k)
		if err != nil {
			return err
		}
		iLen := len(params)
		if iLen >= 1 {
			in := make([]reflect.Value, iLen)
			for i, param := range params {
				in[i] = reflect.ValueOf(param)
			}

			f.Call(in)
		} else {
			log.Printf("CallMsg [%s] params at least one context", funcName)
			return errors.New("callmsg params at least one context")
		}
	}
	return err
}

func (this *Cluster) RandomCluster(head rpc.RpcHead) rpc.RpcHead {
	if head.Id == 0 {
		head.Id = int64(uint32(tools.RAND.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = this.hashRing[head.DestServerType].Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
	}
	return head
}

// 集群新加member
func (this *Cluster) Cluster_Add(ctx context.Context, info *common.ClusterInfo) {
	_, bEx := this.clusterInfoMap[info.Id()]
	if !bEx {
		this.AddCluster(info)
		this.clusterInfoMap[info.Id()] = info
	}
}

// 集群删除member
func (this *Cluster) Cluster_Del(ctx context.Context, info *common.ClusterInfo) {
	delete(this.clusterInfoMap, info.Id())
	this.DelCluster(info)
}
