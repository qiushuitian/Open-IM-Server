package constant

import "errors"

// key = errCode, string = errMsg
type ErrInfo struct {
	ErrCode int32
	ErrMsg  string
}

var (
	OK = ErrInfo{0, ""}

	ErrMysql             = ErrInfo{100, ""}
	ErrMongo             = ErrInfo{110, ""}
	ErrRedis             = ErrInfo{120, ""}
	ErrParseToken        = ErrInfo{200, "Parse token failed"}
	ErrCreateToken       = ErrInfo{201, "Create token failed"}
	ErrAppServerKey      = ErrInfo{300, "key error"}
	ErrTencentCredential = ErrInfo{400, ""}

	ErrorUserRegister             = ErrInfo{600, "User registration failed"}
	ErrAccountExists              = ErrInfo{601, "The account is already registered and cannot be registered again"}
	ErrUserPassword               = ErrInfo{602, "User password error"}
	ErrRefreshToken               = ErrInfo{605, "Failed to refresh token"}
	ErrAddFriend                  = ErrInfo{606, "Failed to add friends"}
	ErrAgreeToAddFriend           = ErrInfo{607, "Failed to agree application"}
	ErrAddFriendToBlack           = ErrInfo{608, "Failed to add friends to the blacklist"}
	ErrGetBlackList               = ErrInfo{609, "Failed to get blacklist"}
	ErrDeleteFriend               = ErrInfo{610, "Failed to delete friend"}
	ErrGetFriendApplyList         = ErrInfo{611, "Failed to get friend application list"}
	ErrGetFriendList              = ErrInfo{612, "Failed to get friend list"}
	ErrRemoveBlackList            = ErrInfo{613, "Failed to remove blacklist"}
	ErrSearchUserInfo             = ErrInfo{614, "Can't find the user information"}
	ErrDelAppleDeviceToken        = ErrInfo{615, ""}
	ErrModifyUserInfo             = ErrInfo{616, "update user some attribute failed"}
	ErrSetFriendComment           = ErrInfo{617, "set friend comment failed"}
	ErrSearchUserInfoFromTheGroup = ErrInfo{618, "There is no such group or the user not in the group"}
	ErrCreateGroup                = ErrInfo{619, "create group chat failed"}
	ErrJoinGroupApplication       = ErrInfo{620, "Failed to apply to join the group"}
	ErrQuitGroup                  = ErrInfo{621, "Failed to quit the group"}
	ErrSetGroupInfo               = ErrInfo{622, "Failed to set group info"}
	ErrParam                      = ErrInfo{700, "param failed"}
	ErrTokenExpired               = ErrInfo{701, TokenExpired.Error()}
	ErrTokenInvalid               = ErrInfo{702, TokenInvalid.Error()}
	ErrTokenMalformed             = ErrInfo{703, TokenMalformed.Error()}
	ErrTokenNotValidYet           = ErrInfo{704, TokenNotValidYet.Error()}
	ErrTokenUnknown               = ErrInfo{705, TokenUnknown.Error()}

	ErrAccess = ErrInfo{ErrCode: 800, ErrMsg: "no permission"}

	ErrDb = ErrInfo{ErrCode: 900, ErrMsg: "db failed"}
)

var (
	TokenExpired     = errors.New("token is timed out, please log in again")
	TokenInvalid     = errors.New("token has been invalidated")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenUnknown     = errors.New("couldn't handle this token")
)

const (
	NoError          = 0
	FormattingError  = 10001
	DatabaseError    = 10002
	LogicalError     = 10003
	ServerError      = 10004
	HttpError        = 10005
	IoErrot          = 10006
	IntentionalError = 10007
	ContentIllegal   = 10101
)

func (e *ErrInfo) Error() string {
	return e.ErrMsg
}
