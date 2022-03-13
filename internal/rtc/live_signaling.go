package rtc

import (
	"errors"
	ws "open_im_sdk/internal/interaction"
	"open_im_sdk/open_im_sdk_callback"
	"open_im_sdk/pkg/common"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/sdk_params_callback"
	api "open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"

	"github.com/golang/protobuf/proto"
)

type LiveSignaling struct {
	*ws.Ws
	listener    open_im_sdk_callback.OnSignalingListener
	loginUserID string
}

func (s *LiveSignaling) DoNotification(msg *api.MsgData, conversationCh chan common.Cmd2Value) {
	var signalReq api.SignalReq
	err := proto.Unmarshal(msg.Content, &signalReq)
	if err != nil {
		log.Error("", "Unmarshal failed ", err.Error())
	}
	s.doSignalPush(&signalReq)
}

//invitee 被邀请者
func (s *LiveSignaling) invite(req *api.SignalInviteReq, callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.InviteCallback {
	var signalReq api.SignalReq
	*signalReq.GetInvite() = *req
	resp, err := s.SendSignalingReqWaitResp(&signalReq, 0, operationID)
	common.CheckAnyErrCallback(callback, 3001, err, operationID)
	switch payload := resp.Payload.(type) {
	case *api.SignalResp_Invite:
		go s.waitPush(req.Invitation.InviterUserID, req.Invitation.InviteeUserIDList[0], "invite", 100, operationID)
		return sdk_params_callback.InviteCallback(payload.Invite)
	default:
		log.Error(operationID, "resp payload type failed ", payload)
		common.CheckAnyErrCallback(callback, 3002, errors.New("resp payload type failed"), operationID)
		return nil
	}
}

func (s *LiveSignaling) waitPush(inviterUserID, inviteeUserID, event string, timeout int, operationID string) {
	req, err := s.SignalingWaitPush(inviterUserID, inviteeUserID, "invite", timeout, operationID)
	if err != nil {
		return
	}
	s.doSignalPush(req)
}

func (s *LiveSignaling) doSignalPush(req *api.SignalReq) {
	//payload.Accept
	switch payload := req.Payload.(type) {
	case *api.SignalReq_Invite:
		s.listener.OnReceiveNewInvitation(utils.StructToJsonString(payload.Invite))
	case *api.SignalReq_Accept:
		s.listener.OnInviteeAccepted(utils.StructToJsonString(payload.Accept))
	case *api.SignalReq_Reject:
		s.listener.OnInviteeRejected(utils.StructToJsonString(payload.Reject))
	case *api.SignalReq_Cancel:
		s.listener.OnInvitationCancelled(utils.StructToJsonString(payload.Cancel))
	default:
		log.Error("", "payload type failed ")
	}
}

func (s *LiveSignaling) inviteInGroup(groupID string, inviteeUserIDList []string, customData string, offlinePushInfo *api.OfflinePushInfo, timeout uint32, callback open_im_sdk_callback.Base, operationID string) sdk_params_callback.InviteInGroupCallback {
	return nil
}

func (s *LiveSignaling) SetListener(listener open_im_sdk_callback.OnSignalingListener, operationID string) {
	s.listener = listener
}

func (s *LiveSignaling) handleSignaling(req *api.SignalReq, callback open_im_sdk_callback.Base, operationID string) {
	resp, err := s.SendSignalingReqWaitResp(req, 0, operationID)
	common.CheckAnyErrCallback(callback, 3001, err, operationID)
	switch payload := resp.Payload.(type) {
	case *api.SignalResp_Accept:
		//return sdk_params_callback.AcceptCallback(payload.Accept)
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.AcceptCallback(payload.Accept)))
	case *api.SignalResp_Reject:
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.RejectCallback(payload.Reject)))
	case *api.SignalResp_HungUp:
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.HungUpCallback(payload.HungUp)))
	case *api.SignalResp_Cancel:
		callback.OnSuccess(utils.StructToJsonString(sdk_params_callback.CancelCallback(payload.Cancel)))
	default:
		log.Error(operationID, "resp payload type failed ", payload)
		common.CheckAnyErrCallback(callback, 3002, errors.New("resp payload type failed"), operationID)
	}
}