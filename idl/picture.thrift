namespace go picture

include "base.thrift"

struct PictureUploadReq {
    1: optional i64 id
}

struct PictureUploadResp {
    1: base.PictureVo picture
    255: base.BaseResp base
}


struct PictureDeleteReq {
    1: i64 id
}

struct PictureDeleteResp {
    255: base.BaseResp base
}

struct PictureUpdateReq {
    1: i64 id
    2: optional string pic_name
    3: optional string introduction (api.vd = "$ == null || len($) < 800")
    4: optional string category
    5: optional list<string> tags
} 

struct PictureUpdateResp {
    255: base.BaseResp base
}

struct PictureEditReq {
    1: i64 id
    2: optional string pic_name
    3: optional string introduction (api.vd = "$ == null || len($) < 800")
    4: optional string category
    5: optional list<string> tags
}

struct PictureEditResp {
    255: base.BaseResp base
}

struct PictureQueryReq {
    1: optional i64 id
    2: optional string pic_name
    3: optional string introduction (api.vd = "$ == null || len($) < 800")
    4: optional string category
    5: optional list<string> tags
    6: optional i64 pic_size
    8: optional i32 pic_width
    9: optional i32 pic_height
    10: optional double pic_scale
    11: optional string pic_format
    12: optional string search_text
    13: optional i64 user_id
    14: optional i64 current_page
} 

struct PictureQueryResp {
    1: list<base.PictureVo> pictures
    2: i64 total
    255: base.BaseResp base
}

struct PictureSearchReq {
    1: optional i64 id
    2: optional string pic_name
    3: optional string introduction (api.vd = "$ == null || len($) < 800")
    4: optional string category
    5: optional list<string> tags
    6: optional i64 pic_size
    8: optional i32 pic_width
    9: optional i32 pic_height
    10: optional double pic_scale
    11: optional string pic_format
    12: optional string search_text
    13: optional i64 user_id
    14: optional i64 current_page
}

struct PictureSearchResp {
    1: list<base.PictureVo> pictures
    2: i64 total
    255: base.BaseResp base
}

struct PictureQueryByIdReq {
    1: i64 id
}

struct PictureQueryByIdResp {
    1: base.PictureVo picture
    255: base.BaseResp base
}

struct PictureGetByIdReq {
    1: i64 id
}

struct PictureGetByIdResp {
    1: base.PictureVo picture
    255: base.BaseResp base
}

service PictureService {
    // admin
    PictureDeleteResp DeletePicture(1: PictureDeleteReq req)
    PictureUpdateResp UpdatePicture(1: PictureUpdateReq req)
    PictureQueryByIdResp GetPictureById(1: PictureQueryByIdReq req)
    PictureQueryResp ListPicture(1: PictureQueryReq req)

    // user
    PictureEditResp EditPicture(1: PictureEditReq req)
    PictureGetByIdResp GetPictureVoById(1: PictureGetByIdReq req)
    PictureSearchResp ListPictureVo(1: PictureSearchReq req)
}


