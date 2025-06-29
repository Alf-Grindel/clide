namespace go clide.user
include "base.thrift"

struct UserRegisterReq {
    1: string user_account (api.vd = "len($) >= 4 && regexp('^[a-zA-Z0-9]+$')")
    2: string user_password (api.vd = "len($) >= 8")
}

struct UserRegisterResp {
    1: i64 id
    255: base.BaseResp resp
}

struct UserLoginReq {
    1: string user_account (api.vd = "len($) >= 4 && regexp('^[a-zA-Z0-9]+$')")
    2: string user_password (api.vd = "len($) >= 8")
}

struct UserLoginResp {
    1: base.UserVo user
    255: base.BaseResp resp
}

struct GetLoginUserReq{}

struct GetLoginUserResp {
    1: base.UserVo user
    255: base.BaseResp resp
}

struct UserLogoutReq{}

struct UserLogoutResp {
    255: base.BaseResp resp
}

struct UserEditReq {
    1: optional string user_password (api.vd = "$ == null || len($) >= 8")
    2: optional string user_avatar
    3: optional string user_profile
}

struct UserEditResp {
    1: base.UserVo user
    255: base.BaseResp resp
}

struct UserSearchReq {
    1: optional i64 id
    2: optional string user_account (api.vd = "$ == null || (len($) >= 4 && regexp('^[a-zA-Z0-9]+$'))")
    3: optional string user_profile
    4: i64 current_page
    5: i64 page_size
}

struct UserSearchResp {
    1: list<base.UserVo> users
    2: i64 total
    255: base.BaseResp resp
}

struct AddUserReq {
    1: string user_account (api.vd = "$ == null || (len($) >= 4 && regexp('^[a-zA-Z0-9]+$'))")
    2: optional string user_avatar
    3: optional string user_profile
    4: optional string user_role
}

struct AddUserResp {
    1: i64 id
    255: base.BaseResp resp
}

struct DeleteUserReq {
    1: i64 id
}

struct DeleteUserResp {
    255: base.BaseResp resp
}

struct UpdateUserReq {
    1: i64 id
    2: optional string user_password (api.vd = "$ == null || len($) >= 8")
    3: optional string user_avatar
    4: optional string user_profile
    5: optional string user_role
}

struct UpdateUserResp {
    1: base.UserVo user
    255: base.BaseResp resp
}

struct QueryUserReq {
    1: optional i64 id
    2: optional string user_account (api.vd = "$ == null || (len($) >= 4 && regexp('^[a-zA-Z0-9]+$'))")
    3: optional string user_role
    4: optional string user_profile
    5: i64 current_page
    6: i64 page_size
}

struct QueryUserResp {
    1: list<base.UserVo> users
    2: i64 total
    255: base.BaseResp resp
}

struct GetUserReq {
    1: i64 id
}

struct GetUserResp {
    1: base.User user
    255: base.BaseResp resp
}

service UserServices {

    ## public
    UserRegisterResp UserRegister(1:UserLoginReq req)
    UserLoginResp UserLogin(1:UserLoginReq req)

    ## auth
    GetLoginUserResp GetLoginUser(1: GetLoginUserReq req)
    UserLogoutResp LogoutUser(1: UserLogoutReq req)
    UserEditResp UserEdit(1: UserEditReq req)
    UserSearchResp UserSearches(1: UserSearchReq req)

    ## admin
    AddUserResp AddUser(1: AddUserReq req)
    DeleteUserResp DeleteUser(1: DeleteUserReq req)
    UpdateUserResp UpdateUser(1: UpdateUserReq req)
    QueryUserResp QueryUsers(1: QueryUserReq req)
    GetUserResp GetUser(1: GetUserReq req)
}