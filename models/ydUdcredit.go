package models

//有盾用户档案

//银行卡信息
type Bankcards struct {
	Bank_name         string
	Bankcard_type     string
	Bankcard_brand    string
	Bankcard_currency string
}

//手机信息
type Mobiles struct {
	Mobile_detail Mobile_detail
	Mobile        string
}

type Mobile_detail struct {
	Province string
	City     string
	Isp      string
}

type User_detail struct {
	User_features                 []User_features
	Actual_loan_platform_count    string
	Repayment_last_date           string
	Last_modified_time            string
	Actual_loan_platform_count_3m string
	Loan_last_date                string
	Actual_loan_platform_count_1m string
	Actual_loan_platform_count_6m string
	Loan_platform_count           string
	Repayment_times_count         string
	Repayment_platform_count_3m   string
	Score                         string
	Loan_platform_count_3m        string
	Loan_platform_count_6m        string
	Names                         string
	Id_detail                     Id_detail
	Id_number_mask                string
	Name_credible                 string
	Repayment_platform_count_1m   string
	Repayment_platform_count      string
	Repayment_platform_count_6m   string
	Loan_platform_count_1m        string
	Risk_evaluation               string
}

type User_features struct {
	User_feature_type  string
	Last_modified_date string
}

type Id_detail struct {
	Birthday string
	Address  string
	Gender   string
	Nation   string
}

type Devices struct {
	Device_name          string
	Device_id            string
	Id_link_device_count string
	Device_link_id_count string
	Device_last_use_date string
	Device_detail        Device_detail
}

type Device_detail struct {
	App_instalment_count string
	Is_rooted            string
	Cheats_device        string
	Is_using_proxy_port  string
	Network_type         string
}

type Loan_industry struct {
	Actual_loan_platform_count string
	Repayment_count_times      string
	Repayment_last_date        string
	Loan_last_date             string
	Name                       string
	Loan_platform_count        string
	Repayment_times_count      string
	Repayment_platform_count   string
}

type YdudcreditResult struct {
	Uid        int
	Ydudcredit Ydudcredit
}

type Ydudcredit struct {
	Header struct {
		Resp_time string
		Ret_msg   string
		Version   string
		Ret_code  string
		Req_time  string
	}
	Body struct {
		Loan_industry            []Loan_industry
		Mobile_count             string
		Devices                  []Devices
		User_detail              User_detail
		Mobiles                  []Mobiles
		User_have_bankcard_count string
		Bankcard_count           string
		Bankcards                []Bankcards
		Ud_order_no              string
	}
}
