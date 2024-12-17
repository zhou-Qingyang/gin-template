package request

type CommonParams struct {
	NickName string `json:"nick_name"`
	Key      string `json:"key"`
	Status   uint8  `json:"status"`
	PageInfo
}
