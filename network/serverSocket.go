package network

import (
	"fmt"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

var SocketServer *ServerSocket

type IServerSocket interface {
	ISocket
	AssignClientId() uint32
	GetClientById(uint32) *ServerSocketClient
	LoadClient() *ServerSocketClient
	AddClient(*net.TCPConn, string, int) *ServerSocketClient
	DelClient(*ServerSocketClient) bool
	StopClient(uint32)
	HeartbeatCheck()
}

type ServerSocket struct {
	Socket
	TCPListener *net.TCPListener
	IdSeed      uint32
	ClientList  map[uint32]*ServerSocketClient
	ClientLock  *sync.RWMutex
	clientCount int
}

func (this *ServerSocket) Init(ip string, port int) bool {
	this.Socket.Init(ip, port)
	this.ClientLock = &sync.RWMutex{}
	this.ClientList = make(map[uint32]*ServerSocketClient)
	this.IP = ip
	this.Port = port
	SocketServer = this
	return true
}

func (this *ServerSocket) Start() bool {
	var zone string
	if this.Zone != "" {
		zone = this.Zone
	}
	var IP = this.IP
	if IP == "" {
		IP = "127.0.0.1"
	}
	var port = this.Port
	if port == 0 {
		port = 8001
	}
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(IP), port, zone})
	if err != nil {
		log.Fatalf("创建链接失败:%v", err)
		return false
	}
	this.TCPListener = listen
	fmt.Println("已初始化连接，等待客户端连接...")
	go this.Run()
	return true
}

func (this *ServerSocket) Run() bool {
	for {
		conn, err := this.TCPListener.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常：", err.Error())
			continue
		}
		fmt.Println("客户端连接:", conn.RemoteAddr().String())
		this.AddClient(conn, conn.RemoteAddr().String(), 1)
	}
}

func (this *ServerSocket) Stop() bool {
	if this.shuttingDown {
		return true
	}
	this.shuttingDown = true
	this.state = SSF_SHUT_DOWN
	return true
}

func (this *ServerSocket) Send(head rpc3.RpcHead, buff []byte) int {
	client := this.GetClientById(head.SocketId)
	if client != nil {
		client.Send(head, buff)
	}
	return 0
}

func (this *ServerSocket) SendMsg(head rpc3.RpcHead, funcName string, params ...interface{}) int {
	client := this.GetClientById(head.SocketId)
	if client != nil {
		return client.Send(head, rpc.Marshal(head, funcName, params...))
	}
	return 0
}

func (this *ServerSocket) Close() {
	defer this.TCPListener.Close()
	this.Clear()
}
func (this *ServerSocket) DelClient(client *ServerSocketClient) bool {
	this.ClientLock.Lock()
	delete(this.ClientList, client.clientId)
	this.ClientLock.Unlock()
	return true
}

func (this *ServerSocket) AssignClientId() uint32 {
	return atomic.AddUint32(&this.IdSeed, 1)
}

func (this *ServerSocket) AddClient(tcpConn *net.TCPConn, addr string, connectType int) bool {
	var socketClient ServerSocketClient
	socketClient.Init("", 0)
	socketClient.ServerSocket = this
	socketClient.ReceiveBufferSize = this.ReceiveBufferSize
	socketClient.clientId = this.AssignClientId()
	socketClient.SetTcpConn(tcpConn)
	socketClient.SetMaxPacketLen(this.GetMaxPacketLen())
	socketClient.IP = addr
	this.ClientLock.Lock()
	this.ClientList[socketClient.clientId] = &socketClient
	this.ClientLock.Unlock()
	socketClient.Start()
	this.clientCount++
	return true
}

func (this *ServerSocket) GetClientById(id uint32) *ServerSocketClient {
	this.ClientLock.RLock()
	client, exist := this.ClientList[id]
	this.ClientLock.RUnlock()
	if exist == true {
		return client
	}
	return nil
}

func (this *ServerSocket) StopClient(id uint32) {
	client := this.GetClientById(id)
	if client != nil {
		client.Stop()
	}
}

func (this *ServerSocket) LoadClient() *ServerSocketClient {
	return &ServerSocketClient{}
}

func (this *ServerSocket) Restart() bool {
	return true
}

func (this *ServerSocket) Connect() bool {
	return true
}
func (this *ServerSocket) Disconnect(bool) bool {
	return true
}

func (this *ServerSocket) OnNetFail(int) {
}
