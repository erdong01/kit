package network

import (
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"github.com/xtaci/kcp-go/v5"
)

type IServerSocket interface {
	ISocket

	AssignClientId() uint32
	GetClientById(uint32) *ServerSocketClient
	LoadClient() *ServerSocketClient
	AddClient(*net.TCPConn, string, int) *ServerSocketClient
	DelClient(*ServerSocketClient) bool
	StopClient(uint32)
}

type ServerSocket struct {
	Socket
	clientCount int
	maxClients  int
	minClients  int
	idSeed      uint32
	clientList  map[uint32]*ServerSocketClient
	clientLock  *sync.RWMutex
	listen      *net.TCPListener
	lock        sync.Mutex
	kcpListern  net.Listener
}

type ClientChan struct {
	pClient *ServerSocketClient
	state   int
	id      int
}

type WriteChan struct {
	buff []byte
	id   int
}

func (this *ServerSocket) Init(ip string, port int, params ...OpOption) bool {
	this.Socket.Init(ip, port, params...)
	this.clientList = make(map[uint32]*ServerSocketClient)
	this.clientLock = &sync.RWMutex{}
	this.IP = ip
	this.Port = port
	return true
}

func (this *ServerSocket) Start() bool {
	if this.IP == "" {
		this.IP = "127.0.0.1"
	}

	var strRemote = fmt.Sprintf("%s:%d", this.IP, this.Port)
	//初始tcp
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
	if err != nil {
		log.Fatalf("%v", err)
	}
	this.listen, err = net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Fatalf("%v", err)
		return false
	}

	//初始kcp
	this.kcpListern, err = kcp.Listen(strRemote)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("启动监听，等待链接！\n")
	//延迟，监听关闭
	//defer ln.Close()
	go this.Run()
	go this.RunKcp()
	return true
}

func (this *ServerSocket) AssignClientId() uint32 {
	return atomic.AddUint32(&this.idSeed, 1)
}

func (this *ServerSocket) GetClientById(id uint32) *ServerSocketClient {
	this.clientLock.RLock()
	client, exist := this.clientList[id]
	this.clientLock.RUnlock()
	if exist == true {
		return client
	}
	return nil
}

func (this *ServerSocket) AddClient(tcpConn net.Conn, addr string, connectType int) *ServerSocketClient {
	socketClient := this.LoadClient()
	if socketClient != nil {
		socketClient.Init("", 0)
		socketClient.ServerSocket = this
		socketClient.ReceiveBufferSize = this.ReceiveBufferSize
		socketClient.SetMaxPacketLen(this.GetMaxPacketLen())
		socketClient.clientId = this.AssignClientId()
		socketClient.IP = addr
		socketClient.SetConnectType(connectType)
		socketClient.SetConn(tcpConn)
		socketClient.SetClientClose(this.GetClientClose()) //自己加的
		this.clientLock.Lock()
		this.clientList[socketClient.clientId] = socketClient
		this.clientLock.Unlock()
		socketClient.Start()
		this.clientCount++
		return socketClient
	} else {
		log.Printf("%s", "无法创建客户端连接对象")
	}
	return nil
}
func (this *ServerSocket) DelClient(client *ServerSocketClient) bool {
	this.clientLock.Lock()
	delete(this.clientList, client.clientId)
	this.clientLock.Unlock()
	return true
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

func (this *ServerSocket) Close() {
	defer this.listen.Close()
	this.Clear()
}

func (this *ServerSocket) Run() bool {
	for {
		conn, err := this.listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常：", err.Error())
			continue
		}
		fmt.Println("客户端连接:", conn.RemoteAddr().String())
		this.handleConn(conn, conn.RemoteAddr().String())
	}
}

func (this *ServerSocket) RunKcp() bool {
	for {
		kcpConn, err := this.kcpListern.Accept()
		handleError(err)
		if err != nil {
			return false
		}

		fmt.Printf("kcp客户端：%s已连接！\n", kcpConn.RemoteAddr().String())
		//延迟，关闭链接
		//defer kcpConn.Close()
		this.handleConn(kcpConn, kcpConn.RemoteAddr().String())
	}
}

func (this *ServerSocket) handleConn(tcpConn net.Conn, addr string) bool {
	if tcpConn == nil {
		return false
	}
	client := this.AddClient(tcpConn, addr, this.connectType)
	if client == nil {
		return false
	}
	return true
}
