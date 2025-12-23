namespace go api

enum UserStatus{
    BANNED = 0 //封禁
    ACTIVE = 1 //正常活跃用户
    POWER  = 2 //管理员
}

struct User{
    1: i64 id
    2:string name
    3:string email
    4:string password
    5:string phone
    6:optional string avatar
    7:optional string bio
    8:optional i32 gender
    9:i64 createdAt
    10:i64 updatedAt
    11:UserStatus status = UserStatus.ACTIVE
    12:optional i64 lastLogin
}

struct SafeUser{
    1:i64 id
    2:string name
    3:string email
    4:string phone
    5:optional string avatar
    6:optional string bio
    7:optional i32 gender
    8:i64 createdAt
    9:i64 updatedAt
    10:UserStatus status
    11:optional i64 lastLogin
}

struct RegisterReq{
    1:string name
    2:string email
    3:string password
    4:string phone
}

struct LoginReq{
    1:string phone
    2:string password
    3:optional string email
}

struct UpdateUserReq{
    1:optional string name
    2:optional string avatar
    3:optional string bio
    4:optional i32 gender
}

struct ChangePasswordReq{
    1:string oldPassword
    2:string newPassword
}

struct ChangeEmailReq{
    1:string newEmail
    2:string password
    3:string code
}

struct ChangePhoneReq{
    1:string newPhone
    2:string password
    3:string code
}

struct RegisterResp{
    1:i64 id
    2:string token
    3:bool success
    4:optional string message
    5:i32 code = 0
}

struct LoginResp{
    1:i64 id
    2:string token
    3:bool success
    4:optional string message
    5:i32 code = 0
}

struct UpdateUserResp{
    1:optional string message
    2:bool success
    3:i32 code = 0
}

struct GetUserProfileReq{
    1:i64 id
}

struct GetUserProfileResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:optional SafeUser user
}

struct LogoutReq{}

struct LogoutResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
}

service UserService{
    RegisterResp Register(1:RegisterReq req)
    LoginResp Login(1:LoginReq req)
    UpdateUserResp UpdateUser(1:UpdateUserReq req)
    UpdateUserResp ChangePassword(1:ChangePasswordReq req)
    UpdateUserResp ChangeEmail(1:ChangeEmailReq req)
    UpdateUserResp ChangePhone(1:ChangePhoneReq req)
    GetUserProfileResp GetUserProfile(1:GetUserProfileReq req)
    LogoutResp Logout(1:LogoutReq req)
}