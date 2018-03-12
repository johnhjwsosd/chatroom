package chat_server

import (
	"fmt"
	"net"
	"sync/atomic"
	"io"
	"bytes"
	"encoding/binary"
)

type server struct{
	clientList []*clientInfo
}

type clientInfo struct{
	reader chan string
	writer chan string
	conn net.Conn
	connID uint32
}

type head struct{
	CmdType uint8
	PacketLen uint32
	Version uint8
}

var FLAGCLIENT uint32 = 1

var ListenPort string

var serverObj *server

func init(){
	serverObj =&server{}
}

func NewServer(port string) *server{
	ListenPort = port
	return serverObj
}

func (s *server) Run(){
	ln, err := net.Listen("tcp", ListenPort)
	if err != nil {
		fmt.Printf("Listen Error: %s\n", err)
		return
	}
	fmt.Println("socket server listening at  ",ListenPort,"  .... ")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Accept Error: %s\n", err)
			continue
		}
		go s.handlerClient(conn)
	}
}

func (s *server) handlerClient(conn net.Conn){
	num := atomic.AddUint32(&FLAGCLIENT,1)
	c :=&clientInfo{make(chan string,1024),make(chan string,1024),conn,num,}
	s.clientList = append(s.clientList,c)
	fmt.Println(" client connect : ",conn,"  client has :",len(s.clientList))

	go c.readMessage()
	go c.sendMessage()

}


func (c *clientInfo) sendMessage (){
	var msg string
	for{
		select{
		case msg =<-c.writer:
			c.conn.Write([]byte(msg))
		}
	}
}

func (c *clientInfo) readMessage(){
	dataPool := bytes.NewBuffer(make([]byte,0,65536))
	dataBuf := make([]byte,1024)

	isReadHead := 0
	msgLen := 0
	cmdType := 0
	var msg []byte

	var headInfo *head
	for {
		n, err := c.conn.Read(dataBuf)
		if err == io.EOF{
			c.clientLeft()
			fmt.Printf("client %d has left,server has client %d now \n",c.connID,len(serverObj.clientList))
			return
		}
		if err != nil {
			c.clientLeft()
			fmt.Println("conn occur falal: ", err, "       closeing connection : ", c.conn)
			return
		}
		n, err = dataPool.Write(dataBuf[:n])
		if err != nil {
			fmt.Printf("Buffer write error: %s\n", err)
			return
		}
		//处理缓存池内的数据
		for{
			//读头
			if isReadHead==0 && dataPool.Len()>=6{
				headInfo = &head{}
				err = binary.Read(dataPool,binary.BigEndian,headInfo)
				if err !=nil{
					fmt.Printf("read buffer occur fatal,%s  %d closed conn  \n",err,c.connID)
					c.clientLeft()
					return
				}
				fmt.Println("head ：",headInfo)
				if headInfo.Version == 101{
					isReadHead = 1
					msgLen = int(headInfo.PacketLen)
					cmdType = int(headInfo.CmdType)
					if msgLen >65535{
						fmt.Printf("illegal data,  %d cloesd conn \n",c.connID)
						c.clientLeft()
						return
					}
				}else {
					fmt.Printf("illegal client")
					c.clientLeft()
					return
				}
			}
			//读包
			if isReadHead ==1 && dataPool.Len() >=msgLen {
				isReadHead = 0
				switch cmdType {
				case 0:
					msg = dataPool.Next(msgLen)
					fmt.Println("个人消息：",string(msg))
					//todo:个人间文本消息转发
				case 1:
					msg = dataPool.Next(msgLen)
					fmt.Println("广播消息：",string(msg))
					c.broadcast(string(msg))
				case 2:
					//todo:其他类型
				}
			}else {  //判断既不读头，又不读报文的时候继续读取连接中内容存入缓存池。
				break
			}
		}
	}
}



func (c *clientInfo) clientLeft(){
	for k,v :=range serverObj.clientList{
		if v.connID == c.connID{
			serverObj.clientList = append(serverObj.clientList[0:k],serverObj.clientList[k+1:]...)
		}
	}
	c.conn.Close()
}

func (c *clientInfo) broadcast(content interface{}){
	if msg,ok :=content.(string);ok{
		for _,v :=range serverObj.clientList{
			v.writer <- msg
		}
	}
	//todo:广播其他类型
}

func(c *clientInfo) makePacket() []byte{
	return  nil
}


