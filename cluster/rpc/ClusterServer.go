package rpc

import (
	"log"
	"sync"

	"github.com/erDong01/micro-kit/actor"
	"github.com/erDong01/micro-kit/cluster/common"
	"github.com/erDong01/micro-kit/network"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/erDong01/micro-kit/tools"
	"golang.org/x/net/context"
)

type (
	HashClusterMap       map[uint32]*common.ClusterInfo
	HashClusterSocketMap map[uint32]*common.ClusterInfo

	//集群服务器
	ClusterServer struct {
		actor.Actor
		*Service           //集群注册
		m_ClusterMap       HashClusterMap
		m_ClusterSocketMap HashClusterSocketMap
		m_ClusterLocker    *sync.RWMutex
		m_pService         *network.ServerSocket //socket管理
		m_HashRing         *tools.HashRing       //hash一致性
	}

	IClusterServer interface {
		InitService(info *common.ClusterInfo, Endpoints []string)
		RegisterClusterCall() //注册集群通用回调
		AddCluster(info *common.ClusterInfo)
		DelCluster(info *common.ClusterInfo)
		GetCluster(rpc3.RpcHead) *common.ClusterInfo
		GetClusterBySocket(uint32) *common.ClusterInfo

		BindServer(*network.ServerSocket)
		SendMsg(rpc3.RpcHead, string, ...interface{}) //发送给集群特定服务器
		Send(rpc3.RpcHead, []byte)                    //发送给集群特定服务器

		RandomCluster(head rpc3.RpcHead) rpc3.RpcHead //随机分配

		sendPoint(rpc3.RpcHead, []byte)               //发送给集群特定服务器
		balanceSend(rpc3.RpcHead, []byte)             //负载给集群特定服务器
		boardCastSend(head rpc3.RpcHead, buff []byte) //给集群广播
	}
)

func (this *ClusterServer) InitService(info *common.ClusterInfo, Endpoints []string) {
	this.Actor.Init()
	this.m_ClusterLocker = &sync.RWMutex{}
	//注册服务器
	this.Service = NewService(info, Endpoints)
	this.m_ClusterMap = make(HashClusterMap)
	this.m_ClusterSocketMap = make(HashClusterSocketMap)
	this.m_HashRing = tools.NewHashRing()
	actor.MGR.RegisterActor(this)
}

func (this *ClusterServer) RegisterClusterCall() {
}

func (this *ClusterServer) AddCluster(info *common.ClusterInfo) {
	this.m_ClusterLocker.Lock()
	this.m_ClusterMap[info.Id()] = info
	this.m_ClusterSocketMap[info.SocketId] = info
	this.m_ClusterLocker.Unlock()
	this.m_HashRing.Add(info.IpString())
	this.m_pService.SendMsg(rpc3.RpcHead{SocketId: info.SocketId}, "COMMON_RegisterResponse")
	switch info.Type {
	case rpc3.SERVICE_GATESERVER:
		log.Printf("ADD Gate SERVER: [%d]-[%s:%d]", info.SocketId, info.Ip, info.Port)
	}
}

func (this *ClusterServer) DelCluster(info *common.ClusterInfo) {
	this.m_ClusterLocker.RLock()
	_, bEx := this.m_ClusterMap[info.Id()]
	this.m_ClusterLocker.RUnlock()
	if bEx {
		this.m_ClusterLocker.Lock()
		delete(this.m_ClusterMap, info.Id())
		delete(this.m_ClusterSocketMap, info.SocketId)
		this.m_ClusterLocker.Unlock()
	}

	this.m_HashRing.Remove(info.IpString())
	log.Printf("服务器断开连接socketid[%d]", info.SocketId)
	switch info.Type {
	case rpc3.SERVICE_GATESERVER:
		log.Printf("与Gate服务器断开连接,id[%d]", info.SocketId)
	}
}

func (this *ClusterServer) GetCluster(head rpc3.RpcHead) *common.ClusterInfo {
	this.m_ClusterLocker.RLock()
	defer this.m_ClusterLocker.RUnlock()
	pClient, bEx := this.m_ClusterMap[head.ClusterId]
	if bEx {
		return pClient
	}
	return nil
}

func (this *ClusterServer) GetClusterBySocket(socketId uint32) *common.ClusterInfo {
	this.m_ClusterLocker.RLock()
	defer this.m_ClusterLocker.RUnlock()
	pClient, bEx := this.m_ClusterSocketMap[socketId]
	if bEx {
		return pClient
	}
	return nil
}

func (this *ClusterServer) BindServer(pService *network.ServerSocket) {
	this.m_pService = pService
}

func (this *ClusterServer) sendPoint(head rpc3.RpcHead, packet rpc3.Packet) {
	pCluster := this.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
		this.m_pService.Send(head, packet)
	}
}

func (this *ClusterServer) balanceSend(head rpc3.RpcHead, packet rpc3.Packet) {
	_, head.ClusterId = this.m_HashRing.Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
		this.m_pService.Send(head, packet)
	}
}

func (this *ClusterServer) boardCastSend(head rpc3.RpcHead, packet rpc3.Packet) {
	clusterList := []*common.ClusterInfo{}
	this.m_ClusterLocker.RLock()
	for _, v := range this.m_ClusterMap {
		clusterList = append(clusterList, v)
	}
	this.m_ClusterLocker.RUnlock()
	for _, v := range clusterList {
		head.SocketId = v.SocketId
		this.m_pService.Send(head, packet)
	}
}

func (this *ClusterServer) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) {
	this.Send(head, rpc.Marshal(head, funcName, params...))
}

func (this *ClusterServer) Send(head rpc3.RpcHead, packet rpc3.Packet) {
	if head.DestServerType != rpc3.SERVICE_GATESERVER {
		this.balanceSend(head, packet)
	} else {
		switch head.SendType {
		case rpc3.SEND_BALANCE:
			this.balanceSend(head, packet)
		case rpc3.SEND_POINT:
			this.sendPoint(head, packet)
		default:
			this.boardCastSend(head, packet)
		}
	}
}

func (this *ClusterServer) RandomCluster(head rpc3.RpcHead) rpc3.RpcHead {
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

func (this *ClusterServer) COMMON_RegisterRequest(ctx context.Context, info *common.ClusterInfo) {
	pServerInfo := info
	pServerInfo.SocketId = this.GetRpcHead(ctx).SocketId

	this.AddCluster(pServerInfo)
}

//链接断开
func (this *ClusterServer) DISCONNECT(ctx context.Context, socketId uint32) {
	pCluster := this.GetClusterBySocket(socketId)
	if pCluster != nil {
		this.DelCluster(pCluster)
	}
}