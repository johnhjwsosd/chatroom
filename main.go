package main

import  (
	"git.oschina.net/Ljohn/chatroom/chat_server"
)
func main(){
	s:= chat_server.NewServer(":6000")
	s.Run()

}
