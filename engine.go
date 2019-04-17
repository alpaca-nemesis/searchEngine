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
)

var (
	// searcher is coroutine safe
	searcher = riot.Engine{}
	path     = "/home/crowix/go/src/storeTest/store"
	text = "我日你个鬼鬼"
	text1 = "hehe"
	text2 = "阿基克里斯蒂发几阿拉开设的发觉了卡萨将大幅卡拉上来看打飞机拉卡撒"

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

type RPCEngine struct {

}

type SearchRequest struct {
	Content string
}

type SearchResponse struct {
	Content types.SearchResp
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

func initEngine() {
	// gob.Register(MyAttriStruct{})

	// var path = "./riot-index"
	searcher.Init(opts)
	defer searcher.Close()
	// os.MkdirAll(path, 0777)

	// Add the document to the index, docId starts at 1
	searcher.Index("21", types.DocData{Content: text})
	searcher.Index("22", types.DocData{Content: text1})
	searcher.Index("23", types.DocData{Content: text2})

	searcher.RemoveDoc("5")

	// Wait for the index to refresh
	searcher.Flush()

	log.Println("Created index number: ", searcher.NumDocsIndexed())
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
	initRPC()
	//restoreIndex()



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
