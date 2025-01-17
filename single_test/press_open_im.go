package main

import (
	"flag"
	"github.com/soloohu/open_im_sdk/pkg/constant"
	"github.com/soloohu/open_im_sdk/pkg/log"
	"github.com/soloohu/open_im_sdk/test"
)

func main() {

	var senderNum *int          //Number of users sending messages
	var singleSenderMsgNum *int //Number of single user send messages
	var intervalTime *int       //Sending time interval, in millisecond
	senderNum = flag.Int("sn", 100, "sender num")
	singleSenderMsgNum = flag.Int("mn", 100, "single sender msg num")
	intervalTime = flag.Int("t", 100, "interval time mill second")
	flag.Parse()
	constant.OnlyForTest = 1
	log.NewPrivateLog("", test.LogLevel)
	log.Warn("", "press test begin, sender num: ", *senderNum, " single sender msg num: ", *singleSenderMsgNum, " send msg total num: ", *senderNum**singleSenderMsgNum)
	test.PressTest(*singleSenderMsgNum, *intervalTime, *senderNum)
	select {}
}
