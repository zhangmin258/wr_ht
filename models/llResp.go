package models

/*//连连付款异步回调
type NotifyLLTradeRequest struct {
	Oid_partner string
	Sign_type   string
	Sign        string
	Dt_order    string
	No_order    string
	Oid_paybill string
	Meney_order string
	Result_pay  string
	Settle_date string
	Info_order  string
	RetCode 	string
	RetMsg      string
}*/

// 连连回调参数
type RechargeFeedback struct {
	Ret_code    string
	Ret_msg     string
	Sign_type   string
	Sign        string
	Oid_partner string
	No_order    string
	Dt_order    string
	Money_order string
	Oid_paybill string
	Result_pay  string
	Settle_date string
	Info_order  string
}

//连连放款订单查询
type LLTradeQueryResponse struct {
	Ret_code    string
	Ret_msg     string
	Sign_type   string
	Sign        string
	Oid_partner string
	Result_pay  string
	Dt_order    string
	No_order    string
	Oid_paybill string
	Meney_order string
	Settle_date string
	Info_order  string
}