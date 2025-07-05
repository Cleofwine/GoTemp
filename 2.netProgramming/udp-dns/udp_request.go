package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strings"
)

type DNSHeader struct {
	ID            uint16
	Flags         uint16
	Questions     uint16
	AnswerRRS     uint16
	AuthorityRRS  uint16
	AdditionalRRS uint16
}

func (header *DNSHeader) SetFlag(qr, opcode, aa, tc, rd, ra, rcode uint16) {
	header.Flags = qr<<15 + opcode<<11 + aa<<10 + tc<<9 + rd<<8 + ra<<7 + rcode
}

type Query struct {
	Type  uint16
	Class uint16
}

// 域名转字节切片
func domainToBytes(domain string) []byte {
	// 将域名解析成相应的字节数组
	// 通过.切分域名，将每一段的字节数及内容保存到切片
	// 具体结构：长度 + 内容，最后以0x00结尾
	var (
		buf      bytes.Buffer
		segments = strings.Split(domain, ".")
	)
	for _, s := range segments {
		binary.Write(&buf, binary.BigEndian, byte(len(s)))
		binary.Write(&buf, binary.BigEndian, []byte(s))
	}
	binary.Write(&buf, binary.BigEndian, byte(0x00))
	return buf.Bytes()
}

func DigDomain(dnsServerAddr, domain string) (queries, answers []string) {
	header := DNSHeader{}
	header.ID = 0xffff
	header.SetFlag(0, 0, 0, 0, 1, 0, 0)
	header.Questions = 1
	header.AnswerRRS = 0
	header.AuthorityRRS = 0
	header.AdditionalRRS = 0

	query := Query{}
	query.Type = 1
	query.Class = 1

	var (
		conn net.Conn
		err  error
		buf  bytes.Buffer
	)
	binary.Write(&buf, binary.BigEndian, header)
	binary.Write(&buf, binary.BigEndian, domainToBytes(domain))
	binary.Write(&buf, binary.BigEndian, query)

	if conn, err = net.Dial("udp", dnsServerAddr); err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	if _, err = conn.Write(buf.Bytes()); err != nil {
		log.Println(err)
		return
	}
	bytes := make([]byte, 1024)
	n, err := conn.Read(bytes)
	if err != nil {
		log.Println(err)
		return
	}
	// fmt.Printf("% x \n", bytes[:n])
	return dNSResponseDecode(bytes[:n])
}

func dNSResponseDecode(res []byte) (queries, answers []string) {
	// res := []byte{255, 255, 129, 128, 0, 1, 0, 3, 0, 0, 0, 0, 3, 119, 119, 119, 5, 98, 97, 105, 100, 117, 3, 99, 111, 109, 0, 0, 1, 0, 1, 192, 12, 0, 5, 0, 1, 0, 0, 3, 31, 0, 18, 3, 119, 119, 119, 1, 97, 6, 115, 104, 105, 102, 101, 110, 3, 99, 111, 109, 0, 192, 43, 0, 1, 0, 1, 0, 0, 1, 47, 0, 4, 180, 101, 50, 242, 192, 43, 0, 1, 0, 1, 0, 0, 1, 47, 0, 4, 180, 101, 50, 188}
	header := res[:12]
	queryNum := uint16(header[4])<<8 + uint16(header[5])
	answerNum := uint16(header[6])<<8 + uint16(header[7])
	data := res[12:]
	index := 0
	queriesBytes := make([][]byte, queryNum)
	answerBytes := make([][]byte, answerNum)

	for i := 0; i < int(queryNum); i++ {
		start := index
		l := 0
		for {
			l = int(data[index])
			if l == 0 {
				break
			}
			index += 1 + l
		}
		queriesBytes[i] = data[start:index]
		index += 5
	}
	if answerNum != 0 {
		for i := 0; i < int(answerNum); i++ {
			start := index
			nums := 2 + 2 + 2 + 4 + 2
			dataLenIndex := start + nums - 2
			dataLen := int(uint16(data[dataLenIndex])<<8 + uint16(data[dataLenIndex+1]))
			index = start + nums + dataLen - 1
			answerBytes[i] = data[start:index]
			index += 1
		}
	}
	queries = make([]string, queryNum)
	for i, bytes := range queriesBytes {
		queries[i] = getQuery(bytes)
	}
	answers = make([]string, answerNum)
	for i, bytes := range answerBytes {
		answers[i] = getAnswer(bytes)
	}
	return queries, answers
}

func getQuery(queryBytes []byte) string {
	return getDomain(queryBytes)
}

func getDomain(domainBytes []byte) string {
	domain := ""
	index := 0
	l := 0
	for {
		if index >= len(domainBytes) {
			break
		}
		l = int(domainBytes[index])
		if l == 0 {
			break
		}
		if index+1+l > len(domainBytes) {
			domain += string(domainBytes[index+1:]) + "."
		} else {
			domain += string(domainBytes[index+1:index+1+l]) + "."
		}
		index += 1 + l
	}
	domain = strings.Trim(domain, ".")
	return domain
}

func getAnswer(answerBytes []byte) string {
	typ := uint16(answerBytes[2])<<8 + uint16(answerBytes[3])
	dataLenIndex := 2 + 2 + 2 + 4
	dataLen := int(uint16(answerBytes[dataLenIndex])<<8 + uint16(answerBytes[dataLenIndex+1]))
	address := answerBytes[dataLenIndex+2 : dataLenIndex+2+dataLen]
	if typ == 1 {
		// IP
		return fmt.Sprintf("%d.%d.%d.%d", address[0], address[1], address[2], address[3])
	} else if typ == 5 {
		// CNAME
		return getDomain(address)
	}
	return ""
}

func main() {
	dnsServerAddr := "114.114.114.114:53"
	domain := "www.baidu.com"
	fmt.Println(DigDomain(dnsServerAddr, domain))
}
