package main

import (
	"io"
	"log"
	"net"
	"os"
)

// 上传
const OP_UPLOAD = 1

// 下载
const OP_DOWNLOAD = 0

const UPLOAD_DIR = "../upload/"
const SCOURCE_DIR = "../source-file/"

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Println(err)
		return
	}
	defer lis.Close()
	log.Println("listen on:", lis.Addr().String())
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("new conn:", conn.RemoteAddr().String())
		go handle(conn)
	}
}

// 处理请求
func handle(conn net.Conn) {
	defer conn.Close()
	bytes := make([]byte, 512)
	_, err := conn.Read(bytes)
	if err != nil {
		log.Println(err)
		return
	}
	op := int(bytes[0])
	fileNameLen := int(bytes[1])
	fileName := string(bytes[2 : 2+fileNameLen])
	if op == OP_UPLOAD {
		// 上传文件
		err := recvFile(conn, UPLOAD_DIR+fileName)
		if err != nil {
			log.Println(err)
			return
		}
	} else if op == OP_DOWNLOAD {
		// 下载文件
		err := sendFile(conn, SCOURCE_DIR+fileName)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func recvFile(conn net.Conn, filePath string) (err error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return err
		}
		if n == 0 {
			log.Println("收到结束字节")
			break
		}
		log.Println("收到字节数：", n)
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
	buf := make([]byte, 2048)
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			log.Println(err)
			return err
		}
		if n == 0 {
			break
		}
		conn.Write(buf[:n])
	}
	return nil
}
