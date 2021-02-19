package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var tag string

const HAND_SHAKE_MSG = "我是打洞消息"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入一个客户端标识")
		os.Exit(0)
	}

	tag = os.Args[1]
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port:9902}
	dstAddr := &net.UDPAddr{IP: net.ParseIP("192.168.56.1"), Port:9527}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	if _, err = conn.Write([]byte("hello, i am new peer:" + tag)); err != nil {
		log.Panic(err)
	}
	data := make([]byte, 1024)
	n, remoteAddr, err := conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	conn.Close()
	anotherPeer := parseAddr(string(data[:n]))
	fmt.Printf("local:%s server:%s another:%s\n", srcAddr, remoteAddr, dstAddr)
	bidirectionHole(srcAddr, &anotherPeer)


}

func parseAddr(addr string) net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP: net.ParseIP(t[0]),
		Port: int(port),
	}
}

func bidirectionHole(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()

	if _,err = conn.Write([]byte(HAND_SHAKE_MSG)); err != nil {
		log.Println("send handshake:", err)
	}

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if _, err = conn.Write([]byte("from ["+ tag +"]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()
	for {
		data := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s \n", err)
		} else {
			log.Printf("收到数据:%s\n", data[:n])
		}
	}
}
