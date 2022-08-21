package res

type Response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type DataList struct {
	Item  interface{} `json:"item"`
	Total uint        `json:"total"`
}

type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
