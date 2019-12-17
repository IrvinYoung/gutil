package ds

//Result response data, using for restful API
type Result struct{
	Ok bool	`json:"ok"`
	Error string	`json:"error"`
	Data interface{} `json:"data"`
}