package sdk_params_callback

import (
	"github.com/soloohu/open_im_sdk/pkg/constant"
	"github.com/soloohu/open_im_sdk/pkg/db/model_struct"
	"github.com/soloohu/open_im_sdk/pkg/server_api_params"
)

//other user
type GetUsersInfoParam []string
type GetUsersInfoCallback []server_api_params.FullUserInfo

//type GetSelfUserInfoParam string
type GetSelfUserInfoCallback *model_struct.LocalUser

type SetSelfUserInfoParam server_api_params.ApiUserInfo

const SetSelfUserInfoCallback = constant.SuccessCallbackDefault
