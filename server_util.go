package iec61850

// #include <iec61850_server.h>
import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"
)

var iedServer *IedServer

func StartServer() {
	// 初始化ied模型
	iedModel := GetModel()
	if iedModel.model == nil {
		fmt.Println("bad.")
		return
	}
	err := RegisterSensorMap()
	if err != nil {
		fmt.Println("Error registering sensor map")
		return
	}
	iedServer = NewServer(iedModel)
	Register_61850_12184(iedModel)
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

func reverse32(x uint32) uint32 {
	x = ((x & 0xff00ff00) >> 8) | ((x & 0x00ff00ff) << 8)
	return (x >> 16) | (x << 16)
}

func castToIntAnylenR(first []byte, n int) uint32 {
	var res uint32
	copy((*[4]byte)(unsafe.Pointer(&res))[:], first[:n])
	return res
}

func castToIntAnylenUR(first []byte, n int) uint32 {
	var res uint32
	// 复制到结构体的末尾，保持高位字节在前
	copy((*[4]byte)(unsafe.Pointer(&res))[4-n:], first[:n])
	return reverse32(res)
}

func castToFloatAnylenR(first []byte, n int) float32 {
	var res float32
	copy((*[4]byte)(unsafe.Pointer(&res))[:], first[:n])
	return res
}

func reverse64(x uint64) uint64 {
	x = ((x & 0xff00ff00ff00ff00) >> 8) | ((x & 0x00ff00ff00ff00ff) << 8)
	x = ((x & 0xffff0000ffff0000) >> 16) | ((x & 0x0000ffff0000ffff) << 16)
	return (x >> 32) | (x << 32)
}

func castToLongAnylenUR(first []byte, n int) uint64 {
	var res uint64
	copy((*[8]byte)(unsafe.Pointer(&res))[8-n:], first[:n])
	return reverse64(res)
}

func msgToMmsData(server *IedServer, msg []byte, n int, timeVal int64) int {
	if n < 13 {
		return -1
	}
	for i := 0; i < n; i++ {
		fmt.Printf("%02X", msg[i])
	}
	fmt.Println()

	sensorId := castToLongAnylenUR(msg, 6)
	// 获取实例编号
	id := sensorMap[int64(sensorId)]
	dataLen := (castToIntAnylenUR(msg[6:7], 1) >> 4) & 0x0f
	versionNumber := (castToIntAnylenUR(msg[2:4], 2) >> 5) & 0x3f
	var newVersion bool
	if versionNumber > 5 && versionNumber < 11 {
		newVersion = true
	} else {
		newVersion = false
	}
	//newVersion := versionNumber > 5 && versionNumber < 11
	if dataLen == 0 {
		return 0 // special msg
	}

	p := msg[7:]
	fmt.Printf("id:%d, DataLen:%d, newver:%v\n", id, dataLen, newVersion)

	for i := 0; i < (int(dataLen)); i++ {
		index := castToIntAnylenR(p, 2)
		lengthFlag := index & 0x03
		index = ((index >> 2) & 0x1ff) ^ ((index >> 4) & 0xe00)
		length := lengthFlag
		if lengthFlag == 0 {
			length = 4
		} else {
			length = uint32(int(castToIntAnylenR(p[2:], int(lengthFlag))))
		}

		var datadst *DataAttribute
		if newVersion {
			datadst = parameter_table_new[index][id]
		} else {
			datadst = parameter_table[index][id]
		}

		if datadst == nil {
			continue
		}

		t := time_table[index][id]

		server.LockDataModel()
		data_type := C.DataAttribute_getType(datadst.attribute)
		switch DataAttributeType(data_type) {
		case IEC61850_INT32:
			value := castToIntAnylenR(p[2+lengthFlag:], int(length))
			fmt.Printf("value: (int32) %v\n", value)
			server.UpdateInt32AttributeValue(datadst, int32(value))
			break
		case IEC61850_INT64:
			value := castToIntAnylenR(p[2+lengthFlag:], int(length))
			fmt.Printf("value: (int64) %v\n", value)
			server.UpdateInt64AttributeValue(datadst, int64(value))
			break
		case IEC61850_FLOAT32:
			value := castToFloatAnylenR(p[2+lengthFlag:], int(length))
			fmt.Printf("value: (float32) %v\n", value)
			server.UpdateFloatAttributeValue(datadst, value)
			break
		case IEC61850_INT32U:
			value := castToIntAnylenR(p[2+lengthFlag:], int(length))
			fmt.Printf("value: (uint32) %v\n", value)
			break
		default:
			panic("unhandled default case")
		}

		server.UpdateUTCTimeAttributeValue(t, timeVal)
		server.UnlockDataModel()

		p = p[2+lengthFlag+length:]
	}

	return 0
}
