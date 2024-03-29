package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/erdong01/kit/base"
	"github.com/erdong01/kit/rpc"
	"github.com/xtaci/kcp-go/v5"
	"google.golang.org/protobuf/proto"
)

type IServerSocket interface {
	ISocket

	AssignClientId() uint32
	GetClientById(uint32) *ServerSocketClient
	LoadClient() *ServerSocketClient
	AddClinet(*net.TCPConn, string, int) *ServerSocketClient
	DelClinet(*ServerSocketClient) bool
	StopClient(uint32)
}

type ServerSocket struct {
	Socket
	clientCount  int
	maxClients   int
	minClients   int
	idSeed       uint32
	clientMap    map[uint32]*ServerSocketClient
	clientLocker *sync.RWMutex
	listen       *net.TCPListener
	lock         sync.Mutex
	kcpListern   net.Listener
}

func (s *ServerSocket) Init(ip string, port int, params ...OpOption) bool {
	s.Socket.Init(ip, port, params...)
	s.clientMap = make(map[uint32]*ServerSocketClient)
	s.clientLocker = &sync.RWMutex{}
	s.ip = ip
	s.port = port
	return true
}

func (s *ServerSocket) Start() bool {
	if s.ip == "" {
		s.ip = "127.0.0.1"
	}

	var strRemote = fmt.Sprintf("%s:%d", s.ip, s.port)
	//初始tcp
	tcpAddr, err := net.ResolveTCPAddr("tcp4", strRemote)
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.listen, err = net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		log.Fatalf("%v", err)
		return false
	}

	//初始kcp
	s.kcpListern, err = kcp.Listen(strRemote)
	if err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("启动监听，等待链接！\n")
	//延迟，监听关闭
	//defer ln.Close()
	go s.Run()
	go s.RunKcp()
	return true
}

func (s *ServerSocket) AssignClientId() uint32 {
	return atomic.AddUint32(&s.idSeed, 1)
}

func (s *ServerSocket) GetClientById(id uint32) *ServerSocketClient {
	s.clientLocker.RLock()
	client, exist := s.clientMap[id]
	s.clientLocker.RUnlock()
	if exist == true {
		return client
	}

	return nil
}

func (s *ServerSocket) AddClinet(conn net.Conn, addr string, connectType int) *ServerSocketClient {
	client := s.LoadClient()
	if client != nil {
		client.Init("", 0)
		client.server = s
		client.receiveBufferSize = s.receiveBufferSize
		client.SetMaxPacketLen(s.GetMaxPacketLen())
		client.clientId = s.AssignClientId()
		client.ip = addr
		client.SetConnectType(connectType)
		client.SetConn(conn)
		s.clientLocker.Lock()
		s.clientMap[client.clientId] = client
		s.clientLocker.Unlock()
		client.Start()
		s.clientCount++
		return client
	} else {
		base.LOG.Printf("%s", "无法创建客户端连接对象")
	}
	return nil
}

func (s *ServerSocket) DelClinet(client *ServerSocketClient) bool {
	s.clientLocker.Lock()
	delete(s.clientMap, client.clientId)
	s.clientLocker.Unlock()
	return true
}

func (s *ServerSocket) StopClient(id uint32) {
	client := s.GetClientById(id)
	if client != nil {
		client.Stop()
	}
}

func (s *ServerSocket) LoadClient() *ServerSocketClient {
	se := &ServerSocketClient{}
	return se
}

func (s *ServerSocket) Send(head rpc.RpcHead, packet rpc.Packet) int {
	client := s.GetClientById(head.SocketId)
	if client != nil {
		client.Send(head, packet)
	}
	return 0
}

func (s *ServerSocket) SendMsg(head rpc.RpcHead, funcName string, params ...interface{}) {
	client := s.GetClientById(head.SocketId)
	if client != nil {
		client.Send(head, rpc.Marshal(&head, &funcName, params...))
	}
}

func (s *ServerSocket) Restart() bool {
	return true
}

func (s *ServerSocket) Connect() bool {
	return true
}

func (s *ServerSocket) Disconnect(bool) bool {
	return true
}

func (s *ServerSocket) OnNetFail(int) {
}

func (s *ServerSocket) Close() {
	defer s.listen.Close()
	defer s.kcpListern.Close()
	s.Clear()
}

func (s *ServerSocket) Run() bool {
	for {
		tcpConn, err := s.listen.AcceptTCP()
		handleError(err)
		if err != nil {
			return false
		}

		fmt.Printf("客户端：%s已连接！\n", tcpConn.RemoteAddr().String())
		//延迟，关闭链接
		//defer tcpConn.Close()
		s.handleConn(tcpConn, tcpConn.RemoteAddr().String())
	}
}

func (s *ServerSocket) RunKcp() bool {
	for {
		kcpConn, err := s.kcpListern.Accept()
		handleError(err)
		if err != nil {
			return false
		}

		fmt.Printf("kcp客户端：%s已连接！\n", kcpConn.RemoteAddr().String())
		//延迟，关闭链接
		//defer kcpConn.Close()
		s.handleConn(kcpConn, kcpConn.RemoteAddr().String())
	}
}

func (s *ServerSocket) handleConn(tcpConn net.Conn, addr string) bool {
	if tcpConn == nil {
		return false
	}

	client := s.AddClinet(tcpConn, addr, s.connectType)
	if client == nil {
		return false
	}

	return true
}

func (s *ServerSocket) SendPacket(head rpc.RpcHead, funcName string, packet proto.Message) int {
	client := s.GetClientById(head.SocketId)
	if client == nil {
		return 0
	}
	return client.SendPacket(head, funcName, packet)
}

// ClientSocket 给客户发送消息
func (s *ServerSocket) ClientSocket(ctx context.Context) *ServerSocketClient {
	rpcHead := ctx.Value("rpcHead").(rpc.RpcHead)
	return s.GetClientById(rpcHead.SocketId)
}
