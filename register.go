package iec61850

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

// 假设这些是一些全局变量
var ln_12184_cfg = make(map[string][]Tuple)
var parameter_table = make(map[int]map[int64]*DataAttribute)
var parameter_table_new = make(map[int]map[int64]*DataAttribute)
var time_table = make(map[int]map[int64]*DataAttribute)
var time_table_new = make(map[int]map[int64]*DataAttribute)
var sensorMap = make(map[int64]int64)

// register_61850_12184
func register_61850_12184(iedModel *IedModel) int {
	var ld *LogicalDevice

	ld = iedModel.GetFirstChild()
	if ld == nil {
		return 0
	}

	var ld_firstChild *LogicalNode
	ld_firstChild.node = (*C.LogicalNode)(unsafe.Pointer(ld.device.firstChild))

	for ln := ld_firstChild; ln.node != C.NULL; ln = ln.Getsibling() {
		// 遍历一个个逻辑节点
		lnname := C2GoStr((*C.char)(unsafe.Pointer(ln.node.name)))
		lnclass := lnname[:4] // 类别名
		inst := lnname[4:]    // 剩余部分

		if _, ok := ln_12184_cfg[lnclass]; !ok {
			continue
		}

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

			for _, part := range partOfPath {
				for C2GoStr((*C.char)(unsafe.Pointer(doo.attribute.name))) != part {
					// 比较字符串
					doo.attribute = (*C.DataAttribute)(unsafe.Pointer(doo.attribute.sibling))
				}
				if doo.attribute == C.NULL {
					break
				}
				// 更新
				doo.attribute = (*C.DataAttribute)(unsafe.Pointer(doo.attribute.firstChild))
			}

			if doo.attribute == C.NULL {
				continue
			}

			if C2GoStr((*C.char)(unsafe.Pointer(doo.attribute.name))) != "t" {
				if index1 != 0 {
					index1 = (index1 & 0x1ff) ^ ((index1 >> 2) & 0xe00)
					if _, ok := parameter_table[index1]; !ok {
						parameter_table[index1] = make(map[int64]*DataAttribute)
					}
					parameter_table[index1][sensorID] = doo
				}
				if index2 != 0 {
					index2 = (index2 & 0x1ff) ^ ((index2 >> 2) & 0xe00)
					if _, ok := parameter_table_new[index2]; !ok {
						parameter_table_new[index2] = make(map[int64]*DataAttribute)
					}
					parameter_table_new[index2][sensorID] = doo
				}
			} else {
				// 处理时间戳
				if index1 != 0 {
					index1 = (index1 & 0x1ff) ^ ((index1 >> 2) & 0xe00)
					if _, ok := time_table[index1]; !ok {
						time_table[index1] = make(map[int64]*DataAttribute)
					}
					time_table[index1][sensorID] = doo
				}
				if index2 != 0 {
					index2 = (index2 & 0x1ff) ^ ((index2 >> 2) & 0xe00)
					if _, ok := time_table_new[index2]; !ok {
						time_table_new[index2] = make(map[int64]*DataAttribute)
					}
					time_table_new[index2][sensorID] = doo
				}
			}
		}
	}
	return 0
}

// registerSensorMap 根据config生成映射表
func registerSensorMap() error {
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
