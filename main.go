package main

import  (
	"github.com/johnhjwsosd/chatroom/chat_server"
	"golang.org/x/net/websocket"
	"net/http"
	"fmt"
	"sync"
	"time"
)
func main(){
	wait := &sync.WaitGroup{}
	go socketMain(wait)
	go webMain(wait)
	go webSocketMain(wait)
	time.Sleep(time.Second*1)
	wait.Wait()
}

func socketMain(w *sync.WaitGroup){
	w.Add(1)
	defer func(){
		w.Done()
		fmt.Println("socket occur fatal has closed")
	}()
	s:= chat_server.NewServer(":6000")
	s.Run()
}


func webSocketMain(w *sync.WaitGroup){
	w.Add(1)
	defer func(){
		fmt.Println("websocket occur fatal has closed")
		w.Done()
	}()
	fmt.Println("websocket run at 8090")
	http.Handle("/web/1",websocket.Handler(chat_server.HandlerWebSocket))
	http.ListenAndServe(":8090", nil)
}

func webMain(w *sync.WaitGroup){
	w.Add(1)
	defer func(){
		w.Done()
		fmt.Println("web server occur fatal has closed")
	}()
	mux:= http.NewServeMux()
	fmt.Println("web server run at 8088")
	mux.HandleFunc("/api/1",test)

	err:= http.ListenAndServe(":8088",middleWare(mux))
	if err !=nil{
		fmt.Println(err.Error())
	}
}

func middleWare(x http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter,req *http.Request){
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Access-Control-Request-Headers, APPID, Authorization, Authorization-Token")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		w.Header().Set("Content-Type", "application/json")
		x.ServeHTTP(w,req)
	})
}

func test(w http.ResponseWriter,req *http.Request){
	w.Write([]byte("{a:test}"))
}
