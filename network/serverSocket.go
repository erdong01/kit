package network

import (
	"fmt"
	"github.com/erDong01/micro-kit/pb/rpc3"
	"github.com/erDong01/micro-kit/rpc"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type ServerSocket struct {
	Socket
	TCPListener *net.TCPListener
	IdSeed      uint32
	ClientList  map[uint32]*ServerSocketClient
	ClientLock  *sync.RWMutex
	clientCount int
}

var (
	Req_REGISTER byte = 1 // 1 --- c register cid
	Res_REGISTER byte = 2 // 2 --- s response

	Req_HEARTBEAT byte = 3 // 3 --- s send heartbeat req
	Res_HEARTBEAT byte = 4 // 4 --- c send heartbeat res

	Req byte = 5 // 5 --- cs send data
	Res byte = 6 // 6 --- cs send ack
)

type CS struct {
	Rch chan []byte
	Wch chan []byte
	Dch chan bool
	u   string
}

var CMap map[string]*CS

func Udp() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{net.ParseIP("0.0.0.0"), 8001, ""})
	if err != nil {
		fmt.Println("Udp家庭")
		return
	}
	fmt.Println(listen)
}

func PushGRT() {
	for {
		time.Sleep(15 * time.Second)
		for k, v := range CMap {
			fmt.Println("push msg to user:" + k)
			v.Wch <- []byte{Req, '#', 'p', 'u', 's', 'h'}
		}
	}
}
func (this *ServerSocket) Init(ip string, port int) bool {
	this.Socket.Init(ip, port)
	this.ClientLock = &sync.RWMutex{}
	this.ClientList = make(map[uint32]*ServerSocketClient)

	this.IP = ip
	this.Port = port
	return true
}

func (this *ServerSocket) StartTcpServer() error {

	CMap = make(map[string]*CS)
	var zone string
	if this.Zone != "" {
		zone = this.Zone
	}
	var IP string = this.IP
	if IP == "" {
		IP = "127.0.0.1"
	}
	var port int = this.Port
	if port == 0 {
		port = 8001
	}
	listen, err := net.ListenTCP("tcp4", &net.TCPAddr{net.ParseIP(IP), port, zone})
	if err != nil {
		log.Fatalf("创建链接失败:%v", err)
		return err
	}
	this.TCPListener = listen
	fmt.Println("已初始化连接，等待客户端连接...")
	go this.Run()
	//go PushGRT()

	return err
}

func (this *ServerSocket) Run() {

	for {
		conn, err := this.TCPListener.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常：", err.Error())
			continue
		}
		fmt.Println("客户端连接来自", conn.RemoteAddr().String())

		this.AddClient(conn, conn.RemoteAddr().String())
		//go Handler(conn)
	}
}

func (this *ServerSocket) Stop() {

}

func (this *ServerSocket) CloseClient(tcpConn *net.TCPConn) error {
	return tcpConn.Close()
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

func (this *ServerSocket) AddClient(tcpConn *net.TCPConn, addr string) bool {
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

func Handler(conn net.Conn) {
	head := rpc3.RpcHead{Code: 200, Msg: "ok"}
	byteD := rpc.Marshal(head, "test")
	defer conn.Close()
	data2 := make([]byte, 128)
	var C *CS
	for {
		conn.Read(data2)
		C = NewCs(conn.LocalAddr().String())
		CMap[conn.LocalAddr().String()] = C
		p, h := rpc.Unmarshal(data2)
		fmt.Println("222222", p, h.ActorName)
		conn.Write(byteD)
		break
	}
	//go WHandler(conn, C)
	//go RHandler(conn, C)
	//go Work(C)
	select {
	case <-C.Dch:
		fmt.Println("close handler goroutine")
	}
}

func WHandler(conn net.Conn, C *CS) {
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case d := <-C.Wch:
			fmt.Println("d:", string(d))
			i, er := conn.Write(d)
			fmt.Println(i, er, "cccc")
		case <-ticker.C:
			if _, ok := CMap[C.u]; !ok {
				fmt.Println("conn die,close WHandler")
				return
			}
		}
	}
}

func RHandler(conn net.Conn, C *CS) {
	for {
		data := make([]byte, 128)
		err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println(err)
		}
		if _, derr := conn.Read(data); derr != nil {
			if data[0] == Res {
				fmt.Println("recv client data ack")
			} else if data[0] == Req {
				fmt.Println("recv client data")
				fmt.Println(data)
				conn.Write([]byte{Res, '#'})
			}
			continue
		}
		conn.Write([]byte{Req_HEARTBEAT, '#'})
		fmt.Println("send ht packet")
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		if _, herr := conn.Read(data); herr == nil {
			fmt.Println("resv ht packet ack")
		} else {
			delete(CMap, C.u)
			fmt.Println("delete user!")
			return
		}
	}
}

func NewCs(uid string) *CS {
	return &CS{Rch: make(chan []byte), Wch: make(chan []byte), u: uid}
}

func Work(C *CS) {
	time.Sleep(5 * time.Second)
	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}
	time.Sleep(15 * time.Second)

	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}
}
