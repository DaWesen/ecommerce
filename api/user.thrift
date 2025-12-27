namespace go api

enum UserStatus{
    BANNED = 0 //封禁
    ACTIVE = 1 //正常活跃用户
    POWER  = 2 //管理员
    Deleted= 3 //被注销
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

//获取用户状态请求
struct GetUserStatusReq{
    1:i64 userId
}

//获取用户状态响应
struct GetUserStatusResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:UserStatus status
    5:bool isBanned
    6:bool isDeleted
    7:optional i64 bannedAt
    8:optional string banReason
}

//封禁用户请求
struct BanUserReq{
    1:i64 userId
    2:string reason
}

//封禁用户响应
struct BanUserResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:optional i64 bannedAt
    5:optional string banReason
}

//解封用户请求
struct UnbanUserReq{
    1:i64 userId
}

//解封用户响应
struct UnbanUserResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
}

//删除用户请求（软删除）
struct DeleteUserReq{
    1:i64 userId
    2:optional string reason
}

//删除用户响应
struct DeleteUserResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:optional i64 deletedAt
}

//恢复用户请求
struct RestoreUserReq{
    1:i64 userId
}

//恢复用户响应
struct RestoreUserResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
}

//更新用户状态请求
struct UpdateUserStatusReq{
    1:i64 userId
    2:UserStatus status
    3:optional string reason
}

//更新用户状态响应
struct UpdateUserStatusResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:UserStatus oldStatus
    5:UserStatus newStatus
    6:i64 updatedAt
}

//用户列表请求
struct ListUsersReq{
    1:optional UserStatus status
    2:optional string keyword
    3:optional i64 minCreatedAt
    4:optional i64 maxCreatedAt
    5:optional i64 minLastLogin
    6:optional i64 maxLastLogin
    7:i32 page = 1
    8:i32 pageSize = 20
    9:optional string orderBy
    10:optional bool desc = false
}

//用户列表响应
struct ListUsersResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:i32 total
    5:i32 page
    6:i32 pageSize
    7:list<SafeUser> users
}

//搜索用户请求
struct SearchUsersReq{
    1:string keyword
    2:optional UserStatus status
    3:i32 page = 1
    4:i32 pageSize = 20
}

//搜索用户响应
struct SearchUsersResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:i32 total
    5:i32 page
    6:i32 pageSize
    7:list<SafeUser> users
}

//用户统计请求
struct CountUsersReq{
    1:optional UserStatus status
}

//用户统计响应
struct CountUsersResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:i64 count
}

//按状态统计响应
struct CountByStatusResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:map<UserStatus, i64> counts
}

//更新密码请求
struct UpdatePasswordReq{
    1:i64 userId
    2:string newPassword
}

//更新密码响应
struct UpdatePasswordResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
}

//更新邮箱请求
struct UpdateEmailReq{
    1:i64 userId
    2:string newEmail
}

//更新邮箱响应
struct UpdateEmailResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
}

//更新手机号请求
struct UpdatePhoneReq{
    1:i64 userId
    2:string newPhone
}

//更新手机号响应
struct UpdatePhoneResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
}

//更新用户资料请求
struct UpdateUserProfileReq{
    1:i64 userId
    2:optional string name
    3:optional string avatar
    4:optional string bio
    5:optional i32 gender
}

//更新用户资料响应
struct UpdateUserProfileResp{
    1:bool success
    2:i32 code = 0
    3:optional string message
    4:optional SafeUser user
}

service UserService {
    RegisterResp Register(1:RegisterReq req)
    LoginResp Login(1:LoginReq req)
    UpdateUserResp UpdateUser(1:UpdateUserReq req)
    UpdateUserResp ChangePassword(1:ChangePasswordReq req)
    UpdateUserResp ChangeEmail(1:ChangeEmailReq req)
    UpdateUserResp ChangePhone(1:ChangePhoneReq req)
    GetUserProfileResp GetUserProfile(1:GetUserProfileReq req)
    LogoutResp Logout(1:LogoutReq req)
    //用户状态管理
    GetUserStatusResp GetUserStatus(1:GetUserStatusReq req)
    BanUserResp BanUser(1:BanUserReq req)
    UnbanUserResp UnbanUser(1:UnbanUserReq req)
    DeleteUserResp DeleteUser(1:DeleteUserReq req)
    RestoreUserResp RestoreUser(1:RestoreUserReq req)
    UpdateUserStatusResp UpdateUserStatus(1:UpdateUserStatusReq req)
    //用户查询
    ListUsersResp ListUsers(1:ListUsersReq req)
    SearchUsersResp SearchUsers(1:SearchUsersReq req)
    //用户统计
    CountUsersResp CountUsers(1:CountUsersReq req)
    CountByStatusResp CountByStatus()
    //管理员操作用户信息
    UpdatePasswordResp AdminUpdatePassword(1:UpdatePasswordReq req)
    UpdateEmailResp AdminUpdateEmail(1:UpdateEmailReq req)
    UpdatePhoneResp AdminUpdatePhone(1:UpdatePhoneReq req)
    UpdateUserProfileResp AdminUpdateUserProfile(1:UpdateUserProfileReq req)
}