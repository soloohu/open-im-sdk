package full

import (
	"errors"
	"github.com/soloohu/open_im_sdk/internal/cache"
	"github.com/soloohu/open_im_sdk/internal/friend"
	"github.com/soloohu/open_im_sdk/internal/group"
	"github.com/soloohu/open_im_sdk/internal/super_group"
	"github.com/soloohu/open_im_sdk/internal/user"
	"github.com/soloohu/open_im_sdk/open_im_sdk_callback"
	"github.com/soloohu/open_im_sdk/pkg/common"
	"github.com/soloohu/open_im_sdk/pkg/constant"
	"github.com/soloohu/open_im_sdk/pkg/db"
	"github.com/soloohu/open_im_sdk/pkg/db/model_struct"
	sdk "github.com/soloohu/open_im_sdk/pkg/sdk_params_callback"
	api "github.com/soloohu/open_im_sdk/pkg/server_api_params"
	"github.com/soloohu/open_im_sdk/pkg/utils"
)

type Full struct {
	user       *user.User
	friend     *friend.Friend
	group      *group.Group
	ch         chan common.Cmd2Value
	userCache  *cache.Cache
	db         *db.DataBase
	SuperGroup *super_group.SuperGroup
}

func (u *Full) Group() *group.Group {
	return u.group
}

func NewFull(user *user.User, friend *friend.Friend, group *group.Group, ch chan common.Cmd2Value, userCache *cache.Cache, db *db.DataBase, superGroup *super_group.SuperGroup) *Full {
	return &Full{user: user, friend: friend, group: group, ch: ch, userCache: userCache, db: db, SuperGroup: superGroup}
}
func (u *Full) getUsersInfo(callback open_im_sdk_callback.Base, userIDList sdk.GetUsersInfoParam, operationID string) sdk.GetUsersInfoCallback {
	friendList := u.friend.GetDesignatedFriendListInfo(callback, userIDList, operationID)
	blackList := u.friend.GetDesignatedBlackListInfo(callback, userIDList, operationID)
	notIn := make([]string, 0)
	for _, v := range userIDList {
		inFriendList := 0
		for _, friend := range friendList {
			if v == friend.FriendUserID {
				inFriendList = 1
				break
			}
		}
		inBlackList := 0
		for _, black := range blackList {
			if v == black.BlockUserID {
				inBlackList = 1
				break
			}
		}
		if inFriendList == 0 && inBlackList == 0 {
			notIn = append(notIn, v)
		}
	}
	//from svr
	publicList := make([]*api.PublicUserInfo, 0)
	if len(notIn) > 0 {
		publicList = u.user.GetUsersInfoFromSvr(callback, notIn, operationID)
		go func() {
			for _, v := range publicList {
				u.userCache.Update(v.UserID, v.FaceURL, v.Nickname)
				//Update the faceURL and nickname information of the local chat history with non-friends
				_ = u.user.UpdateMsgSenderFaceURLAndSenderNickname(v.UserID, v.FaceURL, v.Nickname, constant.SingleChatType)
				//Update session information of local non-friends
				_ = common.TriggerCmdUpdateConversation(common.UpdateConNode{Action: constant.UpdateConFaceUrlAndNickName, Args: common.UpdateConInfo{UserID: v.UserID}}, u.ch)

			}
		}()
	}
	return common.MergeUserResult(publicList, friendList, blackList)
}

func (u *Full) GetGroupInfoFromLocal2Svr(groupID string, sessionType int32) (*model_struct.LocalGroup, error) {
	switch sessionType {
	case constant.GroupChatType:
		return u.group.GetGroupInfoFromLocal2Svr(groupID)
	case constant.SuperGroupChatType:
		return u.GetGroupInfoByGroupID(groupID)
	default:
		return nil, utils.Wrap(errors.New("err sessionType"), "")
	}
}
func (u *Full) GetReadDiffusionGroupIDList(operationID string) ([]string, error) {
	g1, err1 := u.group.GetJoinedDiffusionGroupIDListFromSvr(operationID)
	g2, err2 := u.SuperGroup.GetJoinedGroupIDListFromSvr(operationID)
	var groupIDList []string
	if err1 == nil {
		groupIDList = append(groupIDList, g1...)
	}
	if err2 == nil {
		groupIDList = append(groupIDList, g2...)
	}
	var err error
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	return groupIDList, err
}
