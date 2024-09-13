package iec61850

import (
	"fmt"
	"net"
	"os"
)

func StartListen(sig chan struct{}) {
	// 创建UDP地址对象，指定要监听的IP和端口
	addr, err := net.ResolveUDPAddr("udp", ":9998")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		os.Exit(1)
	}

	// 创建UDP监听器
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP server is listening on port 9998...")

	// 缓冲区大小
	buffer := make([]byte, 1024)

	for {
		select {
		case <-sig:
			fmt.Println("Exiting...")
			return
		default:
			// 读取客户端消息
			n, clientAddr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Error reading from UDP:", err)
				continue
			}
			// 输出收到的消息和客户端地址
			fmt.Printf("Received %s from %s\n", string(buffer[:n]), clientAddr)
			// 回复客户端
			response := []byte("Hello from UDP server!")
			_, err = conn.WriteToUDP(response, clientAddr)
			if err != nil {
				fmt.Println("Error writing to UDP:", err)
				continue
			}
			//fmt.Println("Reading...")
		}
	}
}
