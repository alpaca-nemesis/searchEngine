package main

import (
	"fmt"
	"github.com/go-ego/riot/types"
	"log"
	"net/rpc/jsonrpc"
	"time"
)

// 算数运算请求结构体
type RPCEngine struct {

}

type SearchRequest struct {
	Content string
}

type SearchResponse struct {
	Content types.SearchResp
}


type AddRequest struct {
	Compulsive bool
	Content string
}

type AddResponse struct {
	Content string
}

func search() {
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8096")
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}

	req := SearchRequest{"吃"}
	//req := SearchRequest{"日"}
	var res SearchResponse
	for i:=0; i<1; i++{
		fmt.Println("time:", i)
		err = conn.Call("RPCEngine.Search", req, &res)
		if err != nil {
			log.Fatalln("search error: ", err)
		}
		fmt.Println("request", req.Content)
		fmt.Println("response", res.Content)
		time.Sleep(1000000)
	}

}
func add() {
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8096")
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}

	req := AddRequest{false,"爬上乐山破我想唱歌"}
	var res AddResponse

	err = conn.Call("RPCEngine.AddContent", req, &res)
	fmt.Println("request", req.Content)
	fmt.Println("response", res.Content)
	if err != nil {
		log.Fatalln("add error: ", err)
	}

}

func main(){
	//for i:=0; i<100; i++{
	//	time.Sleep(time.Duration(100)*time.Millisecond)
	//	add()
	//}
	search()
}