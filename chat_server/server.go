package chat_server

import (
	"fmt"
	"net"
	"sync/atomic"
	"io"
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
	fmt.Println("server listening at  ",ListenPort,"  .... ")
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
	for{
		select{
		case <-c.writer:
			c.conn.Write([]byte("test"))
		}
	}
}

func (c *clientInfo) readMessage(){
	databuf := make([]byte, 200)
	for {
		n, err := c.conn.Read(databuf)
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
		msg :=  string(databuf[:n])
		if msg == "1"{
			c.broadcast()
			continue
		}
		fmt.Println(msg)
	}
}



func (c *clientInfo) clientLeft(){
	for k,v :=range serverObj.clientList{
		if v.connID == c.connID{
			serverObj.clientList = append(serverObj.clientList[0:k],serverObj.clientList[k+1:]...)
		}
	}
}

func (c *clientInfo) broadcast(){
	for _,v :=range serverObj.clientList{
		v.writer <- "123"
	}
}


