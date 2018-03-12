package client

import (
	"net"
	"fmt"
	"bytes"
	"encoding/binary"
	"time"
)

type head struct{
	CmdType uint8
	PacketLen uint32
	Version uint8
}

var clientConn net.Conn

func init(){
	conn, err := net.Dial("tcp", "127.0.0.1:6000")
	if err != nil {
		fmt.Printf("Dial error: %s\n", err)
		return
	}
	clientConn = conn

	fmt.Printf("Client: %s\n", clientConn.LocalAddr())
}

func Test_normalSend(){
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	message := []byte("thisisatest")
	msglen := len(message)
	headInfo := &head{1,uint32(msglen),101}

	err := binary.Write(msgbuf,binary.BigEndian,headInfo)
	if err !=nil{
		fmt.Printf("write buffer occur fatal %s ",err)
		return
	}
	_,err= msgbuf.Write(message)
	if err!=nil{
		fmt.Printf("write byte occur fatal %s ",err)
		return
	}
	clientConn.Write(msgbuf.Next(6+msglen))
}

func Test_stickPacketSend(){
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	message := []byte("thisisatest")
	msglen := len(message)
	headInfo := &head{1,uint32(msglen),101}

	err := binary.Write(msgbuf,binary.BigEndian,headInfo)
	if err !=nil{
		fmt.Printf("write buffer occur fatal %s ",err)
		return
	}
	_,err= msgbuf.Write(message)
	if err!=nil{
		fmt.Printf("write byte occur fatal %s ",err)
		return
	}


	message = []byte("hiworld")
	secondLen:= len(message)
	allLen := secondLen + msglen

	headInfo = &head{1,uint32(secondLen),101}
	err = binary.Write(msgbuf,binary.BigEndian,headInfo)
	if err !=nil{
		fmt.Printf("write buffer occur fatal %s ",err)
		return
	}

	_,err= msgbuf.Write(message)
	if err!=nil{
		fmt.Printf("write byte occur fatal %s ",err)
		return
	}
	clientConn.Write(msgbuf.Next(allLen+3))


	time.Sleep(time.Second*3)
	clientConn.Write(msgbuf.Next(6))

	clientConn.Write(msgbuf.Next(3))
}


func Test_errorHeadSend(){
	msgbuf := bytes.NewBuffer(make([]byte, 0, 1024))
	message := []byte("thisisatest")
	msglen := len(message)
	headInfo := &head{1,uint32(msglen),101}

	err := binary.Write(msgbuf,binary.BigEndian,headInfo)
	if err !=nil{
		fmt.Printf("write buffer occur fatal %s ",err)
		return
	}
	_,err= msgbuf.Write(message)
	if err!=nil{
		fmt.Printf("write byte occur fatal %s ",err)
		return
	}
	fmt.Println(msglen)
	clientConn.Write(msgbuf.Next(msglen-10))
	clientConn.Write(msgbuf.Next(10))
}


func Test_read(){
	msg := make([]byte,10240)
	for {
		n, err := clientConn.Read(msg)
		if err != nil {
			fmt.Printf("read conn occur fatal %s", err)
			return
		}
		fmt.Println("Server Send Msg :",string(msg[:n]))
	}
}


