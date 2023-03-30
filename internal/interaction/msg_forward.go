package interaction

import (
	"github.com/soloohu/open_im_sdk/pkg/common"
	"github.com/soloohu/open_im_sdk/pkg/constant"
	"github.com/soloohu/open_im_sdk/pkg/log"
	"github.com/soloohu/open_im_sdk/pkg/utils"
	"github.com/soloohu/open_im_sdk/sdk_struct"
	"github.com/soloohu/open_im_sdk/open_im_sdk_callback"
)

type MsgForward struct {
	*Ws
	LoginUserID        string
	PushMsgAndMaxSeqCh chan common.Cmd2Value

	msgListener          open_im_sdk_callback.OnAdvancedMsgListener
	batchMsgListener     open_im_sdk_callback.OnBatchMsgListener
}

func (m *MsgForward) SetMsgListener(msgListener open_im_sdk_callback.OnAdvancedMsgListener) {
	m.msgListener = msgListener
}

func (m *MsgForward) SetBatchMsgListener(batchMsgListener open_im_sdk_callback.OnBatchMsgListener) {
	m.batchMsgListener = batchMsgListener
}

func (m *MsgForward) doPushMsg(cmd common.Cmd2Value) {
	msg := cmd.Value.(sdk_struct.CmdPushMsgToMsgSync)
	switch msg.Msg.SessionType {
	case constant.SingleChatType:
		if m.batchMsgListener != nil {
			msgstr := utils.StructToJsonString(msg)
			m.batchMsgListener.OnRecvNewMessages(msgstr)
		} else {
			log.Info("Recv push msg.", utils.StructToJsonString(msg))
		}
	default:
		log.Info("Recv other msg.", utils.StructToJsonString(msg))
	}
}

func (m *MsgForward) Work(cmd common.Cmd2Value) {
	switch cmd.Cmd {
	case constant.CmdPushMsg:
		m.doPushMsg(cmd)
	case constant.CmdMaxSeq:
		msg := cmd.Value.(sdk_struct.CmdMaxSeqToMsgSync)
		log.Info("Recv CmdMaxSeq msg.", utils.StructToJsonString(msg))
	default:
		log.Error("", "cmd failed ", cmd.Cmd)
	}
}

func (m *MsgForward) GetCh() chan common.Cmd2Value {
	return m.PushMsgAndMaxSeqCh
}

func NewMsgForward(ws *Ws, loginUserID string, pushMsgAndMaxSeqCh chan common.Cmd2Value) *MsgForward {
	p := &MsgForward{Ws: ws, LoginUserID: loginUserID, PushMsgAndMaxSeqCh: pushMsgAndMaxSeqCh}
	return p
}
