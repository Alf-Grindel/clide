namespace go base

struct BaseResp {
    1: i64 code
    2: string msg
}

struct User {
    1: i64 id
    2: string userAccount
    3: string userAvatar
    4: string userProfile
    5: string userRole
    6: string editTime
    7: string createTime
    8: string updateTime
    9: string isDelete
}

struct UserVo {
    1: i64 id
    2: string userAccount
    3: string userAvatar
    4: string userProfile
    5: string editTime
    6: string createTime
}

struct PictureVo {
    1: i64 id
    2: string url
    3: string pic_name
    4: string introduction
    5: string category
    6: list<string> tags
    7: i64 pic_size
    8: i32 pic_width
    9: i32 pic_height
    10: double pic_scale
    11: string pic_format
    12: string edit_time
    13: string create_time
    14: string update_time
    15: i64 user_id
    16: UserVo user
}