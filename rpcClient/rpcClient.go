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
func main() {
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:8096")
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}

	req := SearchRequest{"日"}
	var res SearchResponse
	for i:=0; i<10; i++{
		fmt.Println("time:", i)
		err = conn.Call("RPCEngine.Search", req, &res)
		if err != nil {
			log.Fatalln("arith error: ", err)
		}
		fmt.Println("request", req.Content)
		fmt.Println("response", res.Content)
		time.Sleep(1000000)
	}

}