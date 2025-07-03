namespace go clide.picture

include "base.thrift"

// public
struct PictureTagCategoryReq {}

struct PictureTagCategoryResp {
    1: list<string> tag_list
    2: list<string> category_list
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
    14: i64 current_page
    15: i64 page_size (api.vd = " $ <=  20")
}

struct PictureSearchResp {
    1: i64 total
    2: list<base.PictureVo> pictures
    255: base.BaseResp base
}

struct PictureGetByIdReq {
    1: i64 id
}

struct PictureGetByIdResp {
    1: base.PictureVo picture
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

struct UploadPictureReq {
    1: optional i64 id
    2: optional string file_url
    3: optional string pic_name
}

struct UploadPictureResp {
    1: i64 id
    255: base.BaseResp base
}

## admin
struct DeletePictureReq {
    1: i64 id
}

struct DeletePictureResp {
    255: base.BaseResp base
}

struct UpdatePictureReq {
    1: i64 id
    2: optional string pic_name
    3: optional string introduction (api.vd = "$ == null || len($) < 800")
    4: optional string category
    5: optional list<string> tags
} 

struct UpdatePictureResp {
    255: base.BaseResp base
}


struct QueryPictureReq {
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
    14: optional string review_status
    15: optional string review_message
    16: optional i64 review_id
    17: i64 current_page
    18: i64 page_size
} 

struct QueryPictureResp {
    1: i64 total
    2: list<base.Picture> pictures
    255: base.BaseResp base
}

struct QueryPictureByIdReq {
    1: i64 id
}

struct QueryPictureByIdResp {
    1: base.Picture picture
    255: base.BaseResp base
}

struct ReviewPictureReq {
    1: i64 id
    2: string review_status
    3: string review_message
}

struct ReviewPictureResp {
    255: base.BaseResp base
}


struct UploadPictureByBatchReq {
    1: string  search_text
    2: optional i64 upload_count (api.vd = " $ == null || $ < 30 ")
}

struct UploadPictureByBatchResp {
    1: i64 upload_count
    255: base.BaseResp base
}

service PictureService {

    ## public
    PictureTagCategoryResp PictureListTagCategory(1: PictureTagCategoryReq req)

    PictureSearchResp PictureSearch(1: PictureSearchReq req)
    PictureGetByIdResp PictureGetById(1: PictureGetByIdReq req)

    ## auth
    PictureEditResp PictureEdit (1: PictureEditReq req)
    UploadPictureResp UploadPicture(1: UploadPictureReq req)

    ## admin
    DeletePictureResp DeletePicture(1: DeletePictureReq req)
    UpdatePictureResp UpdatePicture(1: UpdatePictureReq req)
    QueryPictureResp QueryPicture(1: QueryPictureReq req)
    QueryPictureByIdResp QueryPictureById(1: QueryPictureByIdReq req)
    ReviewPictureResp ReviewPicture(1: ReviewPictureReq req)
    UploadPictureByBatchResp UploadPictureByBatch(1: UploadPictureByBatchReq req)
}

