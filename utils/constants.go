package utils

import (
	"time"
)

//常量配置
const (
	key                = `h*.d;cy7x_12dkx?#j39fdl!` //api数据加密、解密key
	PasswordEncryptKey = "seeyoutomorrow"

	PageSize15     = 15 //列表页每页数据量
	PageSize10     = 10
	PageSize20     = 20
	PageSize40     = 40
	PageSize5      = 5
	PageSize2      = 2
	PageSize4      = 4
	FormatTime     = "15:04:05"            //时间格式
	FormatDate     = "2006-01-02"          //日期格式
	FormatDateTime = "2006-01-02 15:04:05" //完整时间格式

	KONGJIANSENDSMSCOUNT = 20000 //空间畅想每次发送短信的数量
)

//缓存key
const (
	CACHE_KEY_CAROUSE_IMGS           = "CACHE_KEY_CAROUSE_IMGS" //首页轮播图
	CacheKeyUserPrefix               = "weirong_CacheKeyUserPrefix_"
	CacheKeySystemLogs               = "weirong_CacheKeySystemLogs"
	CacheKeyRoleMenuTreePrefix       = "weirong_CacheKeyRoleMenuTreePrefix_"
	CacheKeyLoanApprovalPrefix       = "weirong_CacheKeyLoanApprovalPrefix_"
	CacheKeySystemOrganization       = "weirong_CacheKeySystemOrganization"                 //组织架构缓存key
	CacheKeySystemOrganizationHash   = "weirong_CachekeySystemOrganizationHash"             //组织架构hash缓存key
	CacheKeyConfigLoanVerfiy         = "weirong_CacheKeyConfig_LoanVerfiy_"                 //配置缓存
	WEIRONGEVERYDAYDATASOURCE        = "weirong_CacheKey_EveryDay_Data_Source"              //微融每日数据汇总
	WEIRONGEVERYDAYDATAOUTPUTSOURCE  = "weirong_CacheKey_EveryDay_Data_Out_Put_Source"      //微融每日数据汇总
	WEIRONGCOUNTSOURCE               = "weirong_CacheKey_EveryDay_Data_Count_Source"        //微融每日数据汇总累计总数
	WEIRONGCOUNTOUTPUTSOURCE         = "weirong_CacheKey_EveryDay_Data_Count_Output_Source" //微融每日数据汇总累计总数
	CACHE_KEY_Loan_Accept_State_Deal = "wr_CACHE_KEY_Loan_Accept_State_Deal_"               //微融放款处理

	CacheKeyUserRegisterCount = "weirong_CacheKeyUserRegisterCount" //今日以前累计注册用户数量
	CacheKeyUserIdentifyCount = "weirong_CacheKeyUserIdentifyCount" //今日以前累计认证用户数量
	CacheKeyUserLoanTotal     = "weirong_CacheKeyUserLoanTotal"     //今日以前累计申请借款用户数量
	CacheKeyUserCreditCount   = "weirong_CacheKeyUserCreditCount"   //今日以前累计放款用户数量
	CacheKeyDailyDatas        = "weirong_CacheKeyDailyDatas"        //今日以前数据明细
	CacheKeyActiveData        = "weirong_CacheKeyActiveData"        //今日以前活跃用户明细
	CacheKeyRegisterUsers     = "weirong_CacheKeyRegisterUsers"     //今日以前注册/申请用户数
	CacheKeyIdentifyUsers     = "weirong_CacheKeyIdentifyUsers"     //今日以前认证用户数
	CacheKeyLoanUsers         = "weirong_CacheKeyLoanUsers"         //今日以前申请贷款用户数
	CacheKeyCreditUsers       = "weirong_CacheKeyCreditUsers"       //今日以前放款用户数
	CacheKeyActiveUsers       = "weirong_CacheKeyActiveUsers"       //今日以前活跃用户数

	CacheKayYaoWangRegisterCount   = "weirong_CacheKayYaoWangUserRegisterCount"   //遥望今日以前累计注册用户数量
	CacheKayYaoWangLoginCount      = "weirong_CacheKayYaoWangUserLoginCount"      //遥望今日以前累计登录用户数量
	CacheKayYaoWangTotalClickCount = "weirong_CacheKayYaoWangUserTotalClickCount" //遥望今日以前总uv

	WR_CACHE_KEY_MONEY_SERVICE = "WR_CACHE_KEY_MONEY_SERVICE" // 增值服务详情

	CACHE_KEY_USERS_BEHAVIOR_LIST = "CACHE_KEY_USERS_BEHAVIOR_LIST" //今日之前微融贷款稳下意愿用户汇总
)

//缓存时间
const (
	RedisCacheTime_User         = 15 * time.Minute
	RedisCacheTime_TwoHour      = 2 * time.Hour
	RedisCacheTime_Role         = 15 * time.Second
	RedisCacheTime_Organization = 24 * time.Hour //24 * time.Hour //组织架构信息缓存时间
)

const (
	SMSURLWEIRONGCOUNT         = "http://services.zcmlc.com/v1/sms/balance"               //查询短信剩余数量
	SMSURLWEIRONGONE           = "https://services.zcmlc.com/v1/sms/sendsms"              //单条发送短信
	SMSURLWEIRONGMANY          = "https://services.zcmlc.com/v1/sms/sendsmsbatch"         //多条发送短信
	SMSURLKJCXSEND             = "http://160.19.212.218:8080/eums/utf8/send_strong.do"    //空间畅想发送短信
	SMSURLKJCXBALANCE          = "http://160.19.212.218:8080/eums/utf8/balance_strong.do" //空间畅想查询余额
	SMSURLYRZTMESSAGE          = "http://101.201.239.1/MessageTransferWebAppJs/servlet/messageTransferServiceServletByJson"
	CACHE_KEY_ChECKWITHSENDMSG = "cache_key_checkwith_send_msg"
)

//空间畅想用户密码
const (
	KJCXNAME = "hzcjwldk"   //用户名
	PASSWORD = "1314woaini" //密码
)

//云融正通用户密码
const (
	YRZTNAME     = "hzcjwlyx"
	YRZTPASSWORD = "888888"
)

var (
	LL_oid_partner     string                                                               //连连商户号
	LL_RSA_PRIVATE_KEY string                                                               //连连RSA私钥
	LL_RSA_PUBLIC_KEY  string                                                               //连连rsa公钥
	LL_PUBLIC_KEY      string                                                               //连连rsa公钥
	LL_Trade_Query     = "https://instantpay.lianlianpay.com/paymentapi/queryPayment.htm"   // 连连放款订单查询
	LL_Trade_API       = "https://instantpay.lianlianpay.com/paymentapi/payment.htm"        //连连实时付款API
	LL_Trade_CONFIRM   = "https://instantpay.lianlianpay.com/paymentapi/confirmPayment.htm" //连连确认付款地址
	LL_Trade_CallBack  = "notify/lltrade"                                                   //连连放款回调地址

)

//连连
const (
	LL_SignVerfy        = "https://traderapi.lianlianpay.com/bankcardbindverfy.htm" //连连签约验证
	LL_SignApply        = "https://traderapi.lianlianpay.com/bankcardbind.htm"      //连连签约申请
	LL_BANDCARD_BIN     = "https://queryapi.lianlianpay.com/bankcardbin.htm"
	WithdrawDepositCost = 4 //充值手续费
)

//有盾
const (
	YD_SECURITY_kEY     = "02de8d8f-979b-4c7d-8aaf-647bac07c2ef"
	YD_API_KEY          = "03b8919a-fd35-4c69-bd20-eb8450015112"
	YD_GET_DATA_URL     = "https://da.udcredit.com/frontserver/4.2/partner/get_data?signature="
	YD_HWY_API_KEY      = "HT1oh4pDajElf1cOazGi"
	YD_HWY_UDCREDIT_URL = "https://api4.udcredit.com/dsp-front/4.1/dsp-front/default/"
)

//连连银行卡
const (
	CACHE_KEY_All_BANKCARD_MAP        = "wr_CACHE_KEY_All_BANKCARD_MAP" //支持的银行map
	CACHE_KEY_CheckCardNumber         = "wr_CACHE_KEY_CheckCardNumber_" //验证银行卡号
	CACHE_KEY_ChECKWITHDRAWCASH_DEPID = "wr_CACHE_KEY_withdrawdeposit_depid_"
	CACHE_KEY_LLAPI_TOKEN_            = "wr_CACHE_KEY_LLAPI_TOKEN_" //连连API绑卡token
)

// send to users
const (
	ToUsers = "lwc@zcmlc.com;jgl@zcmlc.com;chenxn@zcmlc.com;angzx@zcmlc.com"
)

//excel表头匹配信息
const HEARDER = "日期,注册人数,认证完成人数,授信人数,申请借款人数,放款人数,放款金额,合作模式,CPA结算的有效事件,CPA的价格,CPS的首借百分比,CPS的复借百分比"

//新浪短链接转长链接
const (
	SHORTURLTOLANG = "https://api.weibo.com/2/short_url/expand.json"
	ACCESSTOKEN    = "2.00Lc7rvGH3sUwC0440fd10cblZoJGE"
)

//config缓存标识
const (
	WR_CACHE_KEY_CONFIG = "WR_CACHE_KEY_CONFIG_" //用户基础信息
)
