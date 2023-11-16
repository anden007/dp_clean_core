package request

type OrderReq struct {
	PageIndex    int    `json:"pageIndex"`              // [必填] 页码
	PageSize     int    `json:"pageSize"`               // [必填] 每页包含条数，上限为500
	Type         int    `json:"type"`                   // [必填] 订单时间查询类型(1：下单时间，2：完成时间（购买用户确认收货时间），3：更新时间
	StartTime    string `json:"startTime"`              // [必填] 开始时间 格式yyyy-MM-dd HH:mm:ss，与endTime间隔不超过1小时
	EndTime      string `json:"endTime"`                // [必填] 结束时间 格式yyyy-MM-dd HH:mm:ss，与startTime间隔不超过1小时
	ChildUnionId int    `json:"childUnionId,omitempty"` // 子推客unionID，传入该值可查询子推客的订单，注意不可和key同时传入。（需联系运营开通PID权限才能拿到数据）
	Key          string `json:"key,omitempty"`          // 工具商传入推客的授权key，可帮助该推客查询订单，注意不可和childUnionid同时传入。（需联系运营开通工具商权限才能拿到数据）
	Fields       string `json:"fields,omitempty"`       // 支持出参数据筛选，逗号','分隔，目前可用：goodsInfo（商品信息）,categoryInfo(类目信息）
	OrderId      int    `json:"orderId,omitempty"`      // 订单号，当orderId不为空时，其他参数非必填
}

func NewOrderReq(timeType int, startTime string, endTime string, pageIndex int, pageSize int) *OrderReq {
	if timeType == 0 {
		timeType = 1 // 下单时间
	}
	if pageSize > 500 {
		pageSize = 500
	}
	return &OrderReq{
		PageIndex:    pageIndex,
		PageSize:     pageSize,
		Type:         timeType,
		StartTime:    startTime,
		EndTime:      endTime,
		ChildUnionId: 0,
		Key:          "",
		Fields:       "goodsInfo,categoryInfo",
		OrderId:      0,
	}
}
