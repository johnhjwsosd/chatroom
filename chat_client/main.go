package main

import (
	"github.com/johnhjwsosd/chatroom/chat_client/client"
)

func main(){

	go client.Test_read()

	for i:=0 ;i<1;i++ {
		client.Test_errorHeadSend()
	}
	select{

	}
}
