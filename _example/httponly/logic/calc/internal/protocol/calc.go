package protocol

type AddReq struct {
	A int32 `json:"a,omitempty"`
	B int32 `json:"b,omitempty"`
}

type AddResp struct {
	Result int32 `json:"result,omitempty"`
}
