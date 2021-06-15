package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"testing"
	"time"
)

const (
	url  = "nats://192.168.2.129:4222"
	subj = "weather"
)

var (
	nc  *nats.Conn
	err error
)

func init() {
	if nc, err = nats.Connect(url); checkErr(err) {
		//
	}
}

func TestNats(t *testing.T) {

	startServer(subj, "s1")
	startServer(subj, "s2")
	startServer(subj, "s3")
	//wait for subscribe complete
	time.Sleep(1 * time.Second)

	startClient(subj)

	select {}

}

//send message to server
func startClient(subj string) {
	nc.Publish(subj, []byte("Sun"))
	time.Sleep(time.Second * 2)

	nc.Publish(subj, []byte("Rain"))
	time.Sleep(time.Second * 2)

	nc.Publish(subj, []byte("Fog"))
	time.Sleep(time.Second * 2)

	nc.Publish(subj, []byte("Cloudy"))
	time.Sleep(time.Second * 2)

}

//receive message
func startServer(subj, name string) {
	go async(nc, subj, name)
}

func async(nc *nats.Conn, subj, name string) {
	nc.Subscribe(subj, func(msg *nats.Msg) {
		fmt.Println(name, "Received a message From Async : ", string(msg.Data))
	})
}

func checkErr(err error) bool {
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
