namespace go picture

include "base.thrift"

struct PictureUploadReq {
    1: optional i64 id
}

struct PictureUploadResp {
    1: base.PictureVo picture
    255: base.BaseResp base
}