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

struct Picture {
    1: i64 id
    2: string url
    3: string picName
    4: string introduction
    5: string category
    6: list<string> tags
    7: i64 picSize
    8: i32 picWidth
    9: i32 picHeight
    10: double picScale
    11: string picFormat
    12: string editTime
    13: string createTime
    14: string updateTime
    15: string isDelete
    16: i64 userId
    17: User user
    18: string reviewStatus
    19: string reviewMessage
    20: i64 reviewId
    21: string reviewTime
}

struct PictureVo {
    1: i64 id
    2: string url
    3: string picName
    4: string introduction
    5: string category
    6: list<string> tags
    7: i64 picSize
    8: i32 picWidth
    9: i32 picHeight
    10: double picScale
    11: string picFormat
    12: string editTime
    13: string createTime
    14: i64 userId
    15: UserVo user
}