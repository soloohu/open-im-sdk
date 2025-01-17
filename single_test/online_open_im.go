package main

import (
	"flag"
	"github.com/soloohu/open_im_sdk/pkg/log"
	"github.com/soloohu/open_im_sdk/test"
)

func main() {
	var onlineNum *int //Number of online users
	onlineNum = flag.Int("on", 10, "online num")
	flag.Parse()
	log.Warn("", "online test start, online num: ", *onlineNum)
	test.OnlineTest(*onlineNum)
	log.Warn("", "online test finish, online num: ", *onlineNum)
	select {}
}
