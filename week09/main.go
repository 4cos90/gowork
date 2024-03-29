package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"main/pkg"
	"net"
	"strconv"
	"time"
)

//1.总结几种 socket 粘包的解包方式：fix length/delimiter based/length field based frame decoder。尝试举例其应用。
//2.实现一个从 socket connection 中解码出 goim 协议的解码器。

var BufferSize int = 1024 //fix length 约定的消息长度
var MessageType int = 4   // 0-原生收发消息 1-fix length 2-delimiter based 3-length field based frame decoder 4-goim协议
var Delimiter byte = 00

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:10001")
	if err != nil {
		log.Fatalf("listen error: %v\n", err)
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Printf("accept error: %v\n", err)
				continue
			}
			switch MessageType {
			case 0:
				go handleConnbase(conn)
			case 1:
				go handleConnfixlength(conn)
			case 2:
				go handleConndelimiterbased(conn)
			case 3:
				go handleConnframedecoder(conn)
			case 4:
				go handleConngoim(conn)
			}

		}
	}()
	go func() {
		for {
			time.Sleep(time.Millisecond * 1000)
			switch MessageType {
			case 0:
				mocksendbase()
			case 1:
				mocksendfixlength()
			case 2:
				mocksenddelimiterbased()
			case 3:
				mocksendbaseframedecoder()
			case 4:
				mocksendgoim()
			}
		}
	}()
	select {}
}

//一个无任何处理的实现，仅接收消息，会出现拆包,粘包的问题。
func handleConnbase(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, BufferSize)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Printf("read error: %v\n", err)
				return
			}
		}
		fmt.Printf("Receive Success by base,Msg: %s\n", buf[0:n])
	}
}

//fix length 实现，通过发送和接收方约定消息长度。
func handleConnfixlength(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, BufferSize)
	result := bytes.NewBuffer(nil)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Printf("read error: %v\n", err)
				return
			}
		}
		result.Write(buf[0:n])
		fmt.Printf("Receive Success by fix length,Msg: %s\n", result.String())
		result.Reset()
	}
}

//delimiter based 在消息后加上分隔符
func handleConndelimiterbased(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, BufferSize)
	result := bytes.NewBuffer(nil)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				return
			} else {
				log.Printf("read error: %v\n", err)
				return
			}
		}
		result.Write(buf[0:n])

		var start int
		var end int
		for k, v := range result.Bytes() {
			if v == Delimiter {
				end = k
				fmt.Printf("Receive Success by delimiter based,Msg: %s\n", string(result.Bytes()[start:end]))
				start = end + 1
			}
		}
		result.Reset()
	}
}

//length field based frame decoder 在消息前加上消息的长度
func handleConnframedecoder(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		peek, err := reader.Peek(4)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Printf("read error: %v\n", err)
				break
			}
		}
		buffer := bytes.NewBuffer(peek)
		var size int32
		if err := binary.Read(buffer, binary.BigEndian, &size); err != nil {
			log.Println(err)
		}
		if int32(reader.Buffered()) < size+4 {
			continue
		}
		data := make([]byte, size+4)
		if _, err := reader.Read(data); err != nil {
			log.Println(err.Error())
			continue
		}
		log.Printf("Receive Success by length field based frame decoder,Msg: %s\n", string(data[4:]))
	}
}

//goim 消息处理
func handleConngoim(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		peek, err := reader.Peek(pkg.PackageLengthSize())
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Printf("read error: %v\n", err)
				break
			}
		}
		buffer := bytes.NewBuffer(peek)
		var size int32
		if err := binary.Read(buffer, binary.BigEndian, &size); err != nil {
			log.Println(err)
		}
		if int32(reader.Buffered()) < size {
			continue
		}
		data := make([]byte, size)
		if _, err := reader.Read(data); err != nil {
			log.Println(err.Error())
			continue
		}
		content, err := pkg.Decoder(data)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		log.Printf("Receive Success by goim,Msg: %s\n", string(content.Content))
	}
}

//无任何处理的发送
func mocksendbase() {
	conn, err := net.Dial("tcp", "127.0.0.1:10001")
	if err != nil {
		fmt.Println("dial failed, err\n", err)
		return
	}
	defer conn.Close()
	msg := "Hello World"
	fmt.Printf("Send Success by base,Msg:%s \n", msg)
	conn.Write([]byte(msg))
}

//fix length 发送方补齐消息长度
func mocksendfixlength() {
	conn, err := net.Dial("tcp", "127.0.0.1:10001")
	if err != nil {
		fmt.Println("dial failed, err\n", err)
		return
	}
	defer conn.Close()
	msg := "Hello World"
	patchmsg := fixlengthpatch(msg)
	fmt.Printf("Send Success by fix length,Msg:%s \n", patchmsg)
	conn.Write([]byte(patchmsg))
}

//fix length 补齐消息长度
func fixlengthpatch(message string) []byte {
	res := make([]byte, BufferSize)
	copy(res, []byte(message))
	return res
}

//delimiter based 在消息后加上分隔符
func mocksenddelimiterbased() {
	conn, err := net.Dial("tcp", "127.0.0.1:10001")
	if err != nil {
		fmt.Printf("dial failed, err %s \n", err)
		return
	}
	defer conn.Close()
	msg := "Hello World"
	patchmsg := delimiterbasedpatch(msg)
	fmt.Printf("Send Success by delimiter based,Msg:%s \n", patchmsg)
	conn.Write([]byte(patchmsg))
}

//delimiter based 在消息后加上分隔符
func delimiterbasedpatch(message string) []byte {
	data := []byte(message)
	data = append(data, Delimiter)
	return data
}

//length field based frame decoder 在消息前加上消息的长度
func mocksendbaseframedecoder() {
	conn, err := net.Dial("tcp", "127.0.0.1:10001")
	if err != nil {
		fmt.Println("dial failed, err\n", err)
		return
	}
	defer conn.Close()
	msg := "Hello World"
	patchmsg, err := Encode(msg)
	if err != nil {
		fmt.Printf("Encode Error: %s \n", err)
	}
	fmt.Printf("Send Success by length field based frame decoder,Msg:%s \n", string(patchmsg))
	conn.Write(patchmsg)
}

// Encode 将消息编码 在消息前加上消息的长度
func Encode(message string) ([]byte, error) {
	size := len(message)
	buf := bytes.NewBuffer(nil)
	if err := binary.Write(buf, binary.BigEndian, int32(size)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, []byte(message)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func mocksendgoim() {
	seq := 0
	version := 1
	code := 1
	conn, err := net.Dial("tcp", "127.0.0.1:10001")
	if err != nil {
		fmt.Println("dial failed, err\n", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 5; i++ {
		msg := "Hello World" + strconv.Itoa(seq)
		patchmsg := pkg.Encoder(pkg.NewPack(version, code, seq, []byte(msg)))
		seq++
		fmt.Printf("Send Success by goim,Msg:%s \n", string(patchmsg))
		conn.Write(patchmsg)
	}
}
