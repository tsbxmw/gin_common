package common

type (
    PageBaseRequest struct {
        PageSize int `json:"page_size"`
        PageIndex int `json:"page_index"`
    }

    PageFormBaseRequest struct {
        PageSize int `form:"page_size"`
        PageIndex int `form:"page_index"`
    }
)
