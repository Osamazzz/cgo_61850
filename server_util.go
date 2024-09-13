package iec61850

import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer() {
	iedModel := GetModel()
	if iedModel.model != nil {
		fmt.Println("good.")
	}
	iedServer := NewServer(iedModel)
	GetMmsServer(iedServer)
	setFileStoreBasePath(iedServer, "./main/")
	iedServer.Start(102)
	defer iedServer.Destroy()
	defer iedServer.Stop()
	if !iedServer.IsRunning() {
		fmt.Println("iedServer is not running")
		iedServer.Destroy()
		return
	}
	iedModel.GetFirstChild()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	sigUdp := make(chan struct{})

	go StartListen(sigUdp)

	for {
		select {
		case <-sig:
			close(sigUdp)
			fmt.Println("exit.")
			time.Sleep(1 * time.Second)
			return
		default:
			fmt.Println("等待数据到来...")
			time.Sleep(10 * time.Second)
		}
	}
}
