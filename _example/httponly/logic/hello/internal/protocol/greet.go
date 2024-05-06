package protocol

type HelloReq struct {
	Name string `json:"name"`
}

type HelloResp struct {
	Message string `json:"message"`
}
