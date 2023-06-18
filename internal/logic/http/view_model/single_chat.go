package view_model

type PushKeysReq struct {
	Op   int32    `json:"op" validate:"required"`
	Keys []string `json:"keys" validate:"required"`
	Msg  string   `json:"msg" validate:"required"`
}

type PushMidsReq struct {
	Op   int32   `json:"op" validate:"required"`
	Mids []int64 `json:"mids" validate:"required"`
	Msg  string  `json:"msg" validate:"required"`
}
