package view_model

type PushRoomReq struct {
	Op   int32  `json:"op" validate:"required"`
	Typ  string `json:"typ" validate:"required"`
	Room string `json:"room" validate:"required"`
	Body string `json:"body" validate:"required"`
}
