package pkg

type APIResult struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Result  interface{} `json:"result"`
}

type APIErrorResult struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message"`
}
