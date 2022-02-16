package rpc

import (
	"context"
	"reflect"
	"sync"

	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/cluster/common"
	"github.com/erDong01/micro-kit/cluster/etcdv3"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
	"github.com/erDong01/micro-kit/tools/vector"
)

type (
	Service    etcdv3.Service
	Master     etcdv3.Master
	Snowflake  etcdv3.Snowflake
	PlayerRaft etcdv3.PlayerRaft
	//集群包管理
	IClusterPacket interface {
		actor.IActor
		SetClusterId(uint32)
	}

	ClusterNode struct {
		*network.ClientSocket
		*common.ClusterInfo
	}

	//集群客户端
	Cluster struct {
		actor.Actor
		m_ClusterMap     map[uint32]*ClusterNode
		m_ClusterLocker  *sync.RWMutex
		m_Packet         IClusterPacket
		m_Master         *Master
		m_HashRing       *tools.HashRing //hash一致性
		m_ClusterInfoMap map[uint32]*common.ClusterInfo
		m_PacketFuncList *vector.Vector //call back
	}

	ICluster interface {
		actor.IActor
		InitCluster(info *common.ClusterInfo, Endpoints []string)
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc3.RpcHead) *ClusterNode

		BindPacket(IClusterPacket)
		BindPacketFunc(network.PacketFunc)

		RandomCluster(head rpc3.RpcHead) rpc3.RpcHead ///随机分配

		sendPoint(rpc3.RpcHead, rpc3.Packet)     //发送给集群特定服务器
		balanceSend(rpc3.RpcHead, rpc3.Packet)   //负载给集群特定服务器
		boardCastSend(rpc3.RpcHead, rpc3.Packet) //给集群广播
	}
)

//注册服务器
func NewService(info *common.ClusterInfo, Endpoints []string) *Service {
	service := &etcdv3.Service{}
	service.Init(info, Endpoints)
	return (*Service)(service)
}

//监控服务器
func NewMaster(info *common.ClusterInfo, Endpoints []string, pActor actor.IActor) *Master {
	master := &etcdv3.Master{}
	master.Init(info, Endpoints, pActor)
	return (*Master)(master)
}

//uuid生成器
func NewSnowflake(Endpoints []string) *Snowflake {
	uuid := &etcdv3.Snowflake{}
	uuid.Init(Endpoints)
	return (*Snowflake)(uuid)
}

func (this *Cluster) InitCluster(info *common.ClusterInfo, Endpoints []string) {
	this.Actor.Init()
	this.m_ClusterLocker = &sync.RWMutex{}
	this.m_ClusterMap = make(map[uint32]*ClusterNode)
	this.m_Master = NewMaster(info, Endpoints, &this.Actor)
	this.m_HashRing = tools.NewHashRing()
	this.m_ClusterInfoMap = make(map[uint32]*common.ClusterInfo)
	this.m_PacketFuncList = vector.NewVector()
	actor.MGR.RegisterActor(this)
	this.Actor.Start()
}

//集群新加member
func (this *Cluster) Cluster_Add(ctx context.Context, info *common.ClusterInfo) {
	_, bEx := this.m_ClusterInfoMap[info.Id()]
	if !bEx {
		this.AddCluster(info)
		this.m_ClusterInfoMap[info.Id()] = info
	}
}

//集群删除member
func (this *Cluster) Cluster_Del(ctx context.Context, info *common.ClusterInfo) {
	delete(this.m_ClusterInfoMap, info.Id())
	this.DelCluster(info)
}

//链接断开
func (this *Cluster) DISCONNECT(ctx context.Context, ClusterId uint32) {
	pInfo, bEx := this.m_ClusterInfoMap[ClusterId]
	if bEx {
		this.DelCluster(pInfo)
	}
	delete(this.m_ClusterInfoMap, ClusterId)
}

func (this *Cluster) AddCluster(info *common.ClusterInfo) {
	pClient := new(network.ClientSocket)
	pClient.Init(info.Ip, int(info.Port))
	packet := reflect.New(reflect.ValueOf(this.m_Packet).Elem().Type()).Interface().(IClusterPacket)
	packet.Init()
	packet.SetClusterId(info.Id())
	pClient.BindPacketFunc(actor.MGR.PacketFunc)
	for _, v := range this.m_PacketFuncList.Values() {
		pClient.BindPacketFunc(v.(network.PacketFunc))
	}
	this.m_ClusterLocker.Lock()
	this.m_ClusterMap[info.Id()] = &ClusterNode{ClientSocket: pClient, ClusterInfo: info}
	this.m_ClusterLocker.Unlock()
	this.m_HashRing.Add(info.IpString())
	pClient.Start()
}

func (this *Cluster) DelCluster(info *common.ClusterInfo) {
	this.m_ClusterLocker.RLock()
	pCluster, bEx := this.m_ClusterMap[info.Id()]
	this.m_ClusterLocker.RUnlock()
	if bEx {
		pCluster.CallMsg(rpc3.RpcHead{}, "STOP_ACTOR")
		pCluster.Stop()
	}

	this.m_ClusterLocker.Lock()
	delete(this.m_ClusterMap, info.Id())
	this.m_ClusterLocker.Unlock()
	this.m_HashRing.Remove(info.IpString())
}

func (this *Cluster) GetCluster(head rpc3.RpcHead) *ClusterNode {
	this.m_ClusterLocker.RLock()
	pCluster, bEx := this.m_ClusterMap[head.ClusterId]
	this.m_ClusterLocker.RUnlock()
	if bEx {
		return pCluster
	}
	return nil
}

func (this *Cluster) BindPacketFunc(callfunc network.PacketFunc) {
	this.m_PacketFuncList.PushBack(callfunc)
}

func (this *Cluster) BindPacket(packet IClusterPacket) {
	this.m_Packet = packet
}

func (this *Cluster) sendPoint(head rpc3.RpcHead, packet rpc3.Packet) {
	pCluster := this.GetCluster(head)
	if pCluster != nil {
		pCluster.Send(head, packet)
	}
}

func (this *Cluster) balanceSend(head rpc3.RpcHead, packet rpc3.Packet) {
	_, head.ClusterId = this.m_HashRing.Get64(head.Id)
	pClient := this.GetCluster(head)
	if pClient != nil {
		pClient.Send(head, packet)
	}
}

func (this *Cluster) boardCastSend(head rpc3.RpcHead, packet rpc3.Packet) {
	clusterList := []*ClusterNode{}
	this.m_ClusterLocker.RLock()
	for _, v := range this.m_ClusterMap {
		clusterList = append(clusterList, v)
	}
	this.m_ClusterLocker.RUnlock()
	for _, v := range clusterList {
		v.Send(head, packet)
	}
}

func (this *Cluster) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
	this.Send(head, rpc.Marshal(head, funcName, params...))
}

func (this *Cluster) Send(head rpc3.RpcHead, packet rpc3.Packet) {
	switch head.SendType {
	case rpc3.SEND_BALANCE:
		this.balanceSend(head, packet)
	case rpc3.SEND_POINT:
		this.sendPoint(head, packet)
	default:
		this.boardCastSend(head, packet)
	}
}

func (this *Cluster) RandomCluster(head rpc3.RpcHead) rpc3.RpcHead {
	if head.Id == 0 {
		head.Id = int64(uint32(tools.RAND.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = this.m_HashRing.Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
	}
	return head
}