package spyun

// PrintReply 对应 POST /v1/printer/print 的 JSON 响应（字段以官方文档为准）。
type PrintReply struct {
	ErrorCode  int    `json:"errorcode"`
	ErrorMsg   string `json:"errormsg"`
	ID         string `json:"id"`
	CreateTime string `json:"create_time"`
}
