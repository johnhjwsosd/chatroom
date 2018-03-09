package main

import (
	"github.com/johnhjwsosd/chatroom/chat_client/client"
)

func main(){

	go client.Test_read()

	for i:=0 ;i<10;i++ {
		client.Test_normalSend()
	}
	select{

	}
}
