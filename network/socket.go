package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

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

type ServerSocket struct {
	IPv4        string
	Port        int
	Zone        string
	TCPListener *net.TCPListener
}

func connen() {

}

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

func (this *ServerSocket) StartTcpServer(listen *net.TCPListener) error {

	CMap = make(map[string]*CS)
	var zone string
	if this.Zone != "" {
		zone = this.Zone
	}
	var ipv4 string
	if this.IPv4 != "" {
		ipv4 = "127.0.0.1"
	}
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ipv4), this.Port, zone})
	if err != nil {
		log.Fatalf("%v", err)
		return err
	}
	this.TCPListener = listen
	fmt.Println("已初始化连接，等待客户端连接...")
	go PushGRT()

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常：", err.Error())
			continue
		}
		fmt.Println("客户端连接来自")
		go Handler(conn)
	}
	return err
}
func Handler(conn net.Conn) {
	defer conn.Close()
	data := make([]byte, 128)
	var uid string
	var C *CS
	for {
		conn.Read(data)
		fmt.Println(data, "111111")
		fmt.Println("客户端发来数据：", string(data))
		if data[0] == Req_REGISTER {
			conn.Write([]byte{Res_REGISTER, '#', 'o', 'k'})
			uid = string(data[2:])
			C = NewCs(uid)
			CMap[uid] = C
			break
		} else {
			conn.Write([]byte{Res_REGISTER, '#', 'e', 'r'})
		}
	}
	go WHandler(conn, C)

	go RHandler(conn, C)

	go Work(C)
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
			conn.Write(d)
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
