package view_model

type PushKeysReq struct {
	Op   int32    `json:"op" validate:"required"`
	Keys []string `json:"keys" validate:"required"`
	Msg  string   `json:"msg" validate:"required"`
}
