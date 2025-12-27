package model

//UserStatus
type UserStatus int32

const (
	UserStatusBANNED  UserStatus = 0 //封禁用户
	UserStatusACTIVE  UserStatus = 1 //活跃用户
	UserStatusPOWER   UserStatus = 2 //管理员
	UserStatusDELETED UserStatus = 3 //软删除用户
)
