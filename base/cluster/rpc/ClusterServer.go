package rpc

import (
	"context"
	"github.com/erdong01/kit/actor"
	"github.com/erdong01/kit/base"
	"github.com/erdong01/kit/network"
	"github.com/erdong01/kit/rpc"
	"sync"
)

type (
	HashClusterMap       map[uint32]*rpc.ClusterInfo
	HashClusterSocketMap map[uint32]*rpc.ClusterInfo

	//集群服务器
	ClusterServer struct {
		actor.Actor
		*Service           //集群注册
		m_ClusterMap       HashClusterMap
		m_ClusterSocketMap HashClusterSocketMap
		m_ClusterLocker    *sync.RWMutex
		m_pService         *network.ServerSocket //socket管理
		m_HashRing         *base.HashRing        //hash一致性
	}

	IClusterServer interface {
		InitService(info *rpc.ClusterInfo, Endpoints []string)
		RegisterClusterCall() //注册集群通用回调
		AddCluster(info *rpc.ClusterInfo)
		DelCluster(info *rpc.ClusterInfo)
		GetCluster(rpc.RpcHead) *rpc.ClusterInfo
		GetClusterBySocket(uint32) *rpc.ClusterInfo

		BindServer(*network.ServerSocket)
		SendMsg(rpc.RpcHead, string, ...interface{}) //发送给集群特定服务器
		Send(rpc.RpcHead, []byte)                    //发送给集群特定服务器

		RandomCluster(head rpc.RpcHead) rpc.RpcHead //随机分配

		sendPoint(rpc.RpcHead, []byte)               //发送给集群特定服务器
		balanceSend(rpc.RpcHead, []byte)             //负载给集群特定服务器
		boardCastSend(head rpc.RpcHead, buff []byte) //给集群广播
	}
)

func (this *ClusterServer) InitService(info *rpc.ClusterInfo, Endpoints []string) {
	this.m_ClusterLocker = &sync.RWMutex{}
	//注册服务器
	this.Service = NewService(info, Endpoints)
	this.m_ClusterMap = make(HashClusterMap)
	this.m_ClusterSocketMap = make(HashClusterSocketMap)
	this.m_HashRing = base.NewHashRing()
}

func (this *ClusterServer) RegisterClusterCall() {
	this.RegisterCall("COMMON_RegisterRequest", func(ctx context.Context, info *rpc.ClusterInfo) {
		pServerInfo := info
		pServerInfo.SocketId = this.GetRpcHead(ctx).SocketId

		this.AddCluster(pServerInfo)
	})

	//链接断开
	this.RegisterCall("DISCONNECT", func(ctx context.Context, socketId uint32) {
		pCluster := this.GetClusterBySocket(socketId)
		if pCluster != nil {
			this.DelCluster(pCluster)
		}
	})
}

func (this *ClusterServer) AddCluster(info *rpc.ClusterInfo) {
	this.m_ClusterLocker.Lock()
	this.m_ClusterMap[info.Id()] = info
	this.m_ClusterSocketMap[info.SocketId] = info
	this.m_ClusterLocker.Unlock()
	this.m_HashRing.Add(info.IpString())
	this.m_pService.SendMsg(rpc.RpcHead{SocketId: info.SocketId}, "COMMON_RegisterResponse")
	switch info.Type {
	case rpc.SERVICE_GATESERVER:
		base.LOG.Printf("ADD Gate SERVER: [%d]-[%s:%d]", info.SocketId, info.Ip, info.Port)
	}
}

func (this *ClusterServer) DelCluster(info *rpc.ClusterInfo) {
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
	base.LOG.Printf("服务器断开连接socketid[%d]", info.SocketId)
	switch info.Type {
	case rpc.SERVICE_GATESERVER:
		base.LOG.Printf("与Gate服务器断开连接,id[%d]", info.SocketId)
	}
}

func (this *ClusterServer) GetCluster(head rpc.RpcHead) *rpc.ClusterInfo {
	this.m_ClusterLocker.RLock()
	defer this.m_ClusterLocker.RUnlock()
	pClient, bEx := this.m_ClusterMap[head.ClusterId]
	if bEx {
		return pClient
	}
	return nil
}

func (this *ClusterServer) GetClusterBySocket(socketId uint32) *rpc.ClusterInfo {
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

func (this *ClusterServer) sendPoint(head rpc.RpcHead, buff []byte) {
	pCluster := this.GetCluster(head)
	var packet rpc.Packet
	packet.Buff = buff
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
		this.m_pService.Send(head, packet)
	}
}

func (this *ClusterServer) balanceSend(head rpc.RpcHead, buff []byte) {
	_, head.ClusterId = this.m_HashRing.Get64(head.Id)
	pCluster := this.GetCluster(head)
	var packet rpc.Packet
	packet.Buff = buff
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
		this.m_pService.Send(head, packet)
	}
}

func (this *ClusterServer) boardCastSend(head rpc.RpcHead, buff []byte) {
	clusterList := []*rpc.ClusterInfo{}
	this.m_ClusterLocker.RLock()
	var packet rpc.Packet
	packet.Buff = buff
	for _, v := range this.m_ClusterMap {
		clusterList = append(clusterList, v)
	}
	this.m_ClusterLocker.RUnlock()
	for _, v := range clusterList {
		head.SocketId = v.SocketId
		this.m_pService.Send(head, packet)
	}
}

func (this *ClusterServer) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	packet := rpc.Marshal(&head, &funcName, params...)

	this.Send(head, packet.Buff)
}

func (this *ClusterServer) Send(head rpc.RpcHead, buff []byte) {
	if head.DestServerType != rpc.SERVICE_GATESERVER {
		this.balanceSend(head, buff)
	} else {
		switch head.SendType {
		case rpc.SEND_BALANCE:
			this.balanceSend(head, buff)
		case rpc.SEND_POINT:
			this.sendPoint(head, buff)
		default:
			this.boardCastSend(head, buff)
		}
	}
}

func (this *ClusterServer) RandomCluster(head rpc.RpcHead) rpc.RpcHead {
	if head.Id == 0 {
		head.Id = int64(uint32(base.RandI(1, 0xFFFFFFFF)))
	}
	_, head.ClusterId = this.m_HashRing.Get64(head.Id)
	pCluster := this.GetCluster(head)
	if pCluster != nil {
		head.SocketId = pCluster.SocketId
	}
	return head
}
