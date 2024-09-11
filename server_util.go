package iec61850

import "C"
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var sensorMap = make(map[int64]int64)

func StartServer() {
	IedModel := GetModel()
	IedServer := NewServer(IedModel)
	GetMmsServer(IedServer)
	IedServer.Start(102)
}

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
