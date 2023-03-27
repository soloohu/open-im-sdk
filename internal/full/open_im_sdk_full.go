package full

import (
	"github.com/soloohu/open_im_sdk/open_im_sdk_callback"
	"github.com/soloohu/open_im_sdk/pkg/common"
	"github.com/soloohu/open_im_sdk/pkg/log"
	"github.com/soloohu/open_im_sdk/pkg/sdk_params_callback"
	"github.com/soloohu/open_im_sdk/pkg/utils"
)

func (u *Full) GetUsersInfo(callback open_im_sdk_callback.Base, userIDList string, operationID string) {
	fName := utils.GetSelfFuncName()
	go func() {
		log.NewInfo(operationID, fName, "args: ", userIDList)
		var unmarshalParam sdk_params_callback.GetUsersInfoParam
		common.JsonUnmarshalAndArgsValidate(userIDList, &unmarshalParam, callback, operationID)
		result := u.getUsersInfo(callback, unmarshalParam, operationID)
		callback.OnSuccess(utils.StructToJsonStringDefault(result))
		log.NewInfo(operationID, fName, "callback: ", utils.StructToJsonStringDefault(result))
	}()
}
