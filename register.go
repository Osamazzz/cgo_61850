package iec61850

// #include "iec61850_server.h"
import "C"
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type Tuple struct {
	Str      string
	OldParam int
	NewParam int
}

// 全局变量
var ln_12184_cfg map[string][]Tuple
var parameter_table [0x1000]map[int64]*DataAttribute
var parameter_table_new [0x1000]map[int64]*DataAttribute
var time_table [0x1000]map[int64]*DataAttribute
var time_table_new [0x1000]map[int64]*DataAttribute
var sensorMap = make(map[int64]int64)

func InitTable() {
	for i, _ := range parameter_table {
		parameter_table[i] = make(map[int64]*DataAttribute)
	}
	for i, _ := range parameter_table_new {
		parameter_table_new[i] = make(map[int64]*DataAttribute)
	}
	for i, _ := range time_table {
		time_table[i] = make(map[int64]*DataAttribute)
	}
	for i, _ := range time_table_new {
		time_table_new[i] = make(map[int64]*DataAttribute)
	}
}

func Init_12184_cfg() {
	ln_12184_cfg = map[string][]Tuple{
		"STMP": {
			{"Tmp/mag/f", 5, 0b00000000000101},
			{"Tmp/t", 5, 0b00000000000101},
		},
		"SPDC": {
			{"Tmp/mag/f", 4213, 0b01000001110101},
			{"Tmp/t", 4213, 0b01000001110101},
			{"ZTDDB/mag/f", 4214, 0b01000001110110},
			{"ZTDDB/t", 4214, 0b01000001110110},
			{"Ultrasonic/mag/f", 2054, 0b100000000110},
			{"Ultrasonic/t", 2054, 0b100000000110},
		},
	}
}

// Register_61850_12184 注册配置信息
func Register_61850_12184(iedModel *IedModel) int {
	Init_12184_cfg()
	InitTable()
	var ld *LogicalDevice
	ld = iedModel.GetFirstChild()
	if ld.device == nil {
		fmt.Println("ld is null")
		return 0
	}
	// 有问题.必须直接赋值
	ld_firstChild := &LogicalNode{
		(*C.LogicalNode)(unsafe.Pointer(ld.device.firstChild)),
	}

	if ld_firstChild.node.sibling == nil {
		fmt.Println("sb is null")
		return 0
	}
	getsibling := ld_firstChild.Getsibling()
	if getsibling == nil {
		fmt.Println("sibling is null")
		return 0
	}

	for ln := ld_firstChild; ln.node != nil; ln = ln.Getsibling() {
		// 遍历一个个逻辑节点
		if ln.node.name == nil {
			fmt.Println("name is null")
			return 0
		}
		lnname := C.GoString(ln.node.name)
		lnclass := lnname[:4] // 类别名
		inst := lnname[4:]    // 剩余部分

		if _, ok := ln_12184_cfg[lnclass]; !ok {
			continue
		}

		// 提取sensor的ID
		sensorID, err := strconv.ParseInt(inst, 10, 64)
		if err != nil {
			fmt.Println("Error converting instance to sensor ID:", err)
			continue
		}

		for _, cfg := range ln_12184_cfg[lnclass] {
			// 遍历一个类下的一个个元组
			path := cfg.Str
			index1 := cfg.OldParam
			index2 := cfg.NewParam

			partOfPath := strings.Split(path, "/")
			//var doo *DataAttribute
			//doo.attribute =
			doo := &DataAttribute{
				attribute: (*C.DataAttribute)(unsafe.Pointer(ln.node.firstChild)),
			}
			i := 0

			for doo.attribute != nil && i < len(partOfPath) {
				if C2GoStr((*C.char)(unsafe.Pointer(doo.attribute.name))) == partOfPath[i] {
					i++
					if i < len(partOfPath) {
						doo.attribute = (*C.DataAttribute)(unsafe.Pointer(doo.attribute.firstChild))
					}
				} else {
					doo.attribute = (*C.DataAttribute)(unsafe.Pointer(doo.attribute.sibling))
				}
			}

			if C2GoStr((*C.char)(unsafe.Pointer(doo.attribute.name))) != "t" {
				// 如果不是时间戳的话
				if index1 != 0 {
					index1 = (index1 & 0x1ff) ^ ((index1 >> 2) & 0xe00)
					//if _, ok := parameter_table[index1]; !ok {
					//	parameter_table[index1] = make(map[int64]*DataAttribute)
					//}
					parameter_table[index1][sensorID] = doo
				}
				if index2 != 0 {
					index2 = (index2 & 0x1ff) ^ ((index2 >> 2) & 0xe00)
					//if _, ok := parameter_table_new[index2]; !ok {
					//	parameter_table_new[index2] = make(map[int64]*DataAttribute)
					//}
					parameter_table_new[index2][sensorID] = doo
				}
			} else {
				// 处理时间戳
				if index1 != 0 {
					index1 = (index1 & 0x1ff) ^ ((index1 >> 2) & 0xe00)
					//if _, ok := time_table[index1]; !ok {
					//	time_table[index1] = make(map[int64]*DataAttribute)
					//}
					time_table[index1][sensorID] = doo
				}
				if index2 != 0 {
					index2 = (index2 & 0x1ff) ^ ((index2 >> 2) & 0xe00)
					//if _, ok := time_table_new[index2]; !ok {
					//	time_table_new[index2] = make(map[int64]*DataAttribute)
					//}
					time_table_new[index2][sensorID] = doo
				}
			}
		}
	}
	return 0
}

// RegisterSensorMap 根据config生成映射表
func RegisterSensorMap() error {
	// 打开文件
	file, err := os.Open("./config/sensormap.txt")
	if err != nil {
		return fmt.Errorf("failed to open sensormap.txt: %w", err)
	}
	defer file.Close()

	// 创建一个 scanner 用于逐行读取文件
	scanner := bufio.NewScanner(file)
	lineNum := 0

	// 遍历文件中的每一行
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 通过逗号分割每行的内容
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return fmt.Errorf("failed to parse line %d: %s", lineNum, line)
		}

		// 将第一个部分解析为十六进制的 long
		sensorID, err := strconv.ParseInt(parts[0], 16, 64)
		if err != nil {
			return fmt.Errorf("failed to parse sensor ID on line %d: %w", lineNum, err)
		}

		// 将第二个部分解析为整数
		inst, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse inst on line %d: %w", lineNum, err)
		}

		// 将结果保存到 map 中
		sensorMap[sensorID] = inst
	}

	// 检查是否在扫描过程中出现错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	return nil
}
