package cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/cluster/common"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
	"github.com/erDong01/micro-kit/tools/vector"
	"github.com/nats-io/nats.go"
	"log"
	"reflect"
	"sync"
	"time"
)

const (
	MAX_CLUSTER_NUM = int(rpc3.SERVICE_ZONESERVER) + 1
	CALL_TIME_OUT   = 50 * time.Millisecond
)

type (
	HashClusterMap       map[uint32]*common.ClusterInfo
	HashClusterSocketMap map[uint32]*common.ClusterInfo

	//集群服务器
	Cluster struct {
		actor.Actor
		*Service       //集群注册
		master         *Master
		clusterMap     [MAX_CLUSTER_NUM]HashClusterMap
		clusterLocker  [MAX_CLUSTER_NUM]*sync.RWMutex
		hashRing       [MAX_CLUSTER_NUM]*tools.HashRing //hash一致性
		conn           *nats.Conn
		dieChan        chan bool
		clusterInfoMap map[uint32]*common.ClusterInfo
		packetFuncList *vector.Vector //call back
		callBackMap    sync.Map
	}

	ICluster interface {
		Init(num int, info *common.ClusterInfo, Endpoints []string, natsUrl string)
		RegisterClusterCall() //注册集群通用回调
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc3.RpcHead) *common.ClusterInfo

		BindPacketFunc(packetFunc network.PacketFunc)
		SendMsg(rpc3.RpcHead, string, ...interface{})                    //发送给集群特定服务器
		Send(rpc3.RpcHead, []byte)                                       //发送给集群特定服务器
		CallMsg(interface{}, rpc3.RpcHead, string, ...interface{}) error //同步给集群特定服务器

		RandomCluster(head rpc3.RpcHead) rpc3.RpcHead //随机分配
	}

	EmptyClusterInfo struct {
		common.ClusterInfo
	}
)

func (this *EmptyClusterInfo) String() string {
	return ""
}

func (this *Cluster) Init(num int, info *common.ClusterInfo, Endpoints []string, natsUrl string) {
	this.Actor.Init(num)
	this.RegisterClusterCall()
	for i := 0; i < MAX_CLUSTER_NUM; i++ {
		this.clusterLocker[i] = &sync.RWMutex{}
		this.clusterMap[i] = make(HashClusterMap)
		this.hashRing[i] = tools.NewHashRing()
	}
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
	fmt.Println(GetChannel(*info))
	this.conn.Subscribe(GetChannel(*info), func(msg *nats.Msg) {
		fmt.Println("11111")
		this.HandlePacket(rpc3.Packet{Buff: msg.Data})
	})
	fmt.Println(GetTopicChannel(*info))
	this.conn.Subscribe(GetTopicChannel(*info), func(msg *nats.Msg) {
		fmt.Println("2222")
		this.HandlePacket(rpc3.Packet{Buff: msg.Data})
	})
	fmt.Println(GetCallChannel(*info))
	this.conn.Subscribe(GetCallChannel(*info), func(msg *nats.Msg) {
		fmt.Println("3333")
		this.HandlePacket(rpc3.Packet{Buff: msg.Data, Reply: msg.Reply})
	})
	rpc.GCall = reflect.ValueOf(this.call)
	this.Actor.Start()
}

//params[0]:rpc.RpcHead
//params[1]:error
func (this *Cluster) call(parmas ...interface{}) {
	head := *parmas[0].(*rpc3.RpcHead)
	reply := head.Reply
	head.Reply = ""
	head.ClusterId = head.SrcClusterId
	if parmas[1] == nil {
		parmas[1] = ""
	} else {
		parmas[1] = parmas[1].(error).Error()
	}
	buff := rpc.Marshal(head, "", parmas[1:]...)
	this.conn.Publish(reply, buff)
}

func (this *Cluster) RegisterClusterCall() {
	this.RegisterCall("Cluster_Add", func(ctx context.Context, info *common.ClusterInfo) {
		_, bEx := this.clusterInfoMap[info.Id()]
		if !bEx {
			this.AddCluster(info)
			this.clusterInfoMap[info.Id()] = info
		}
	})

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
	}
	this.hashRing[info.Type].Remove(info.IpString())
	log.Printf("服务器[%s:%s:%d]断开连接", info.String(), info.Ip, info.Port)
}

func (this *Cluster) GetCluster(head rpc3.RpcHead) *common.ClusterInfo {
	this.clusterLocker[head.DestServerType].RLock()
	defer this.clusterLocker[head.DestServerType].RUnlock()
	client, bEx := this.clusterMap[head.DestServerType][head.ClusterId]
	if bEx {
		return client
	}
	return nil
}

func (this *Cluster) BindPacketFunc(callFunc network.PacketFunc) {
	this.packetFuncList.PushBack(callFunc)
}

func (this *Cluster) HandlePacket(packet rpc3.Packet) {
	for _, v := range this.packetFuncList.Values() {
		if v.(network.PacketFunc)(packet) {
			break
		}
	}
}

func (this *Cluster) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
	head.SrcClusterId = this.Id()
	buff := rpc.Marshal(head, funcName, params...)
	this.Send(head, buff)
}

func (this *Cluster) Send(head rpc3.RpcHead, buff []byte) {

	switch head.SendType {
	case rpc3.SEND_BALANCE:
		_, head.ClusterId = this.hashRing[head.DestServerType].Get64(head.Id)
		err := this.conn.Publish(GetRpcChannel(head), buff)
		fmt.Println(err)
	case rpc3.SEND_POINT:
		this.conn.Publish(GetRpcChannel(head), buff)
	default:
		this.conn.Publish(GetRpcTopicChannel(head), buff)
	}
}

func (this *Cluster) CallMsg(cb interface{}, head rpc3.RpcHead, funcName string, params ...interface{}) error {
	head.SrcClusterId = this.Id()
	buff := rpc.Marshal(head, funcName, params...)
	switch head.SendType {
	case rpc3.SEND_POINT:
	default:
		_, head.ClusterId = this.hashRing[head.DestServerType].Get64(head.Id)
	}
	reply, err := this.conn.Request(GetRpcCallChannel(head), buff, CALL_TIME_OUT)
	if err == nil {
		rpcPacket, _ := rpc.Unmarshal(reply.Data)
		var cf *actor.CallFunc
		val, bOk := this.callBackMap.Load(funcName)
		if !bOk {
			cf = &actor.CallFunc{Func: cb, FuncVal: reflect.ValueOf(cb), FuncType: reflect.TypeOf(cb), FuncParams: reflect.TypeOf(cb).String()}
			this.callBackMap.Store(funcName, cf)
		} else {
			cf = val.(*actor.CallFunc)
		}
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
			this.Trace(funcName)
			f.Call(in)
			this.Trace("")
		} else {
			log.Printf("CallMsg [%s] params at least one context", funcName)
			return errors.New("callmsg params at least one context")
		}
	}
	return err
}

func (this *Cluster) RandomCluster(head rpc3.RpcHead) rpc3.RpcHead {
	if head.Id == 0 {
		head.Id = int64((uint32(tools.RAND.RandI(1, 0xFFFFFFFF))))
	}
	_, head.ClusterId = this.hashRing[head.DestServerType].Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
	}
	return head
}
