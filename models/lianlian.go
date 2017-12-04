package models

//连连银行卡卡bin查询接口 返回数据
type LLBandCardBinResponse struct {
	Ret_code  string
	Ret_msg   string
	Sign_type string
	Sign      string
	Bank_code string
	Bank_name string
	Card_type string
}

type LLTradeApiResponse struct {
	Ret_code     string
	Ret_msg      string
	Sign_type    string
	Sign         string
	Oid_partner  string
	No_order     string
	Oid_paybill  string
	Confirm_code string
	Token        string
}

type LLTradeResponse struct {
	Ret_code     string
	Ret_msg      string
	Sign_type    string
	Sign         string
	Oid_partner  string
	No_order     string
	Oid_paybill  string
	Confirm_code string
}

