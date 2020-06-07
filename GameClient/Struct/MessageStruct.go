package Struct

type TransData struct {
	Type string `json:"type"`
	Message []byte `json:"message"`
	Reciever string
}

type OptData struct {
	Type string `json:"type"`
	User string `json:"user"`
	Px string `json:"px"`
	Py string `json:"py"`
	Message string `json:"message"`
}

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}