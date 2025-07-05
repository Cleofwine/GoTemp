package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
)

// 上传
const OP_UPLOAD = 1

// 下载
const OP_DOWNLOAD = 0

var (
	up       int
	fileName string
)

const DOWNLOAD_DIR = "../download/"
const SCOURCE_DIR = "../source-file/"

func getCmdArgs() {
	flag.IntVar(&up, "up", 1, "0 下载 ，1 上传")
	flag.StringVar(&fileName, "file", "", "上传或下载的文件名称")
	flag.Parse()
}

// 使用方法：
// 在server已经启动的情况下，上传：
// go run ./client.go --up=1 --file=1.jpg
// 下载：
// go run ./client.go --up=0 --file=2.jpg
func main() {
	getCmdArgs()
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	_, fileName = path.Split(fileName)
	isUp := byte(1)
	if up != OP_UPLOAD {
		isUp = byte(0)
	}
	bytes := make([]byte, 0, 512) // 长度为0，512容量
	bytes = append(bytes, isUp)
	bytes = append(bytes, byte(len(fileName)))
	bytes = append(bytes, []byte(fileName)...)
	fmt.Println(bytes)
	// ！！！必须切到512大小，不然会被tcp底层从后面的包里面取数据重新组包到这里面，导致处理这个文件头的过程中漏读一部分数据
	if _, err := conn.Write(bytes[:512]); err != nil {
		log.Println(err)
		return
	}
	if up == OP_UPLOAD {
		// 发送文件
		err := sendFile(conn, SCOURCE_DIR+fileName)
		if err != nil {
			log.Println(err)
			return
		}
	} else if up == OP_DOWNLOAD {
		// 接收文件
		err := recvFile(conn, DOWNLOAD_DIR+fileName)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		return
	}
}

func recvFile(conn net.Conn, filePath string) (err error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return err
		}
		if n == 0 {
			break
		}
		file.Write(buf[:n])
	}
	return nil
}

func sendFile(conn net.Conn, filePath string) (err error) {
	src, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer src.Close()
	buf := make([]byte, 1024)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return err
		}
		if n == 0 {
			break
		}
		fmt.Println("读取字节数：", n)
		nc, err := conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("发送字节数：", nc)
	}
	return nil
}
