// Copyright 2016 ego authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"github.com/go-ego/riot"
	"github.com/go-ego/riot/types"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

var (
	// searcher is coroutine safe
	searcher = riot.Engine{}
	path     = "/home/crowix/go/src/searchEngine/store"

	opts = types.EngineOpts{
		Using: 1,
		IndexerOpts: &types.IndexerOpts{
			IndexType: types.DocIdsIndex,
		},
		UseStore:    true,
		StoreFolder: path,
		StoreEngine: "bolt", // bg: badger, lbd: leveldb, bolt: bolt
		// GseDict: "../../data/dict/dictionary.txt",
		GseDict:       "zh",
		StopTokenFile: "/home/crowix/go/src/data/dict/stop_tokens.txt",
	}
)


//--------------------search rpc engine--------------
type RPCEngine struct {
	flushCount int
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

func (this *RPCEngine)Search(req SearchRequest, res *SearchResponse) error{
	sea := searcher.Search(types.SearchReq{
		Text: req.Content,
		RankOpts: &types.RankOpts{
			OutputOffset: 0,
			MaxOutputs:   100,
		}})
	res.Content = sea
	fmt.Println("search response: ", sea)
	return nil
}

func (this *RPCEngine)AddContent(req AddRequest, res *AddResponse) error{
	index := int(searcher.NumDocsIndexed()) + this.flushCount

	searcher.Index( strconv.Itoa(index), types.DocData{Content: req.Content})
	//time.Sleep(time.Duration(1)*time.Second)
	searcher.Flush()
	//time.Sleep(time.Duration(1)*time.Second)
	res.Content = fmt.Sprintln("Created index number: ", searcher.NumDocsIndexed())
	log.Println("Created index number: ", searcher.NumDocsIndexed()-1)
	return nil
}


func initEngine() {
	// gob.Register(MyAttriStruct{})

	// var path = "./riot-index"
	searcher.Init(opts)
	defer searcher.Flush()

	save := searcher.NumDocsIndexed()
	var text string
	for i:=0;i<10;i++{
		save = save + 1
		text = fmt.Sprint("%s:%d","我日你个鬼鬼",  save)
		time.Sleep(time.Duration(1)*time.Second)
		searcher.Index(string(save), types.DocData{Content: text})
	}
	searcher.Flush()
	//localAdd()
	// os.MkdirAll(path, 0777)

	// Add the document to the index, docId starts at 1

	//
	////searcher.RemoveDoc("5")
	//
	////Wait for the index to refresh

	log.Println("Created index number: ", searcher.NumDocsIndexed())
}

func localAdd(){
	save := searcher.NumDocsIndexed()
	var text string
	for i:=0;i<10;i++{
		save = save + 1
		text = fmt.Sprint("%s:%d","我日你个鬼鬼",  save)
		searcher.Index(string(save), types.DocData{Content: text})
	}
	searcher.Flush()
}

func initRPC(){
	err := rpc.Register(new(RPCEngine)) // 注册rpc服务
	if err == nil{
		fmt.Println("RPC service has been registered")
	}

	lis, err := net.Listen("tcp", "127.0.0.1:8096")
	if err != nil {
		log.Fatalln("fatal error: ", err)
	}

	fmt.Println( "start connection")

	for {
		conn, err := lis.Accept() // 接收客户端连接请求
		if err != nil {
			continue
		}

		go func(conn net.Conn) { // 并发处理客户端请求
			fmt.Println( "new client in coming")
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}

func main() {
	initEngine()
	//localAdd()
	initRPC()
	//restoreIndex()
	defer searcher.Close()



	//sea := searcher.Search(types.SearchReq{
	//	Text: "日",
	//	RankOpts: &types.RankOpts{
	//		OutputOffset: 0,
	//		MaxOutputs:   100,
	//	}})
	//
	//fmt.Println("search response: ", sea)
	//fmt.Println("docs: ", sea.Docs)
	//
	//
	//sea = searcher.Search(types.SearchReq{
	//	Text: "hehe",
	//	RankOpts: &types.RankOpts{
	//		OutputOffset: 0,
	//		MaxOutputs:   100,
	//	}})
	//
	//fmt.Println("search response: ", sea)
	//fmt.Println("docs: ", sea.Docs)

	// os.RemoveAll("riot-index")
}
