package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	serverAddr := "127.0.0.1" // 服务器地址

	// 运行ethr命令
	cmd := exec.Command("./ethr", "-c", serverAddr)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("无法获取命令输出：%s\n", err.Error())
		return
	}

	// 启动ethr进程
	err = cmd.Start()
	if err != nil {
		fmt.Printf("无法启动ethr：%s\n", err.Error())
		return
	}

	// 读取并解析ethr输出
	bindwidthContent := ""
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.HasPrefix(line, "[") {

			bindwidthContent += line + "\n"
		}
	}

	// 检查是否有读取输出时的错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取ethr输出时发生错误：%s\n", err.Error())
	}

	// 等待ethr进程退出
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("ethr进程退出时发生错误：%s\n", err.Error())
	}
	fmt.Println(bindwidthContent)
}
