package routers

import (
	"wr_v1/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.AccountController{}, "get:Login;post:CheckPassword") //登录页
	beego.Router("/loginout", &controllers.AccountController{}, "get:LoginOut")              //退出登录
	beego.MyAutoRouter(&controllers.ProductController{})                                     // 产品管理
	beego.MyAutoRouter(&controllers.OrganizationController{})                                // 代理产品
	beego.MyAutoRouter(&controllers.CleaningController{})                                    // 产品清算
	beego.MyAutoRouter(&controllers.UrlController{})                                         // 产品清算
	beego.MyAutoRouter(&controllers.UserController{})                                        // 用户管理
	beego.MyAutoRouter(&controllers.SystemController{})                                      // 系统管理
	beego.MyAutoRouter(&controllers.ControlController{})                                     // 产品运营管理
	beego.MyAutoRouter(&controllers.UserMetadataController{})                                // 认证用户管理
	//beego.MyAutoRouter(&controllers.TestController{})                                        // 认证用户管理
	beego.MyAutoRouter(&controllers.AgentProController{})    //下级代理商品
	beego.MyAutoRouter(&controllers.AgentController{})       //下级代理商品
	beego.MyAutoRouter(&controllers.AgentProUrlController{}) //外放链接
	beego.MyAutoRouter(&controllers.ProductDataController{}) //产品数据
	beego.MyAutoRouter(&controllers.ChannelAnalysisController{})
	beego.MyAutoRouter(&controllers.ChannelController{}) //渠道管理
	beego.MyAutoRouter(&controllers.BusinessLoanController{})
	beego.MyAutoRouter(&controllers.ProductAnalysisController{}) //产品分析
	beego.MyAutoRouter(&controllers.ChannelAnalysisController{}) //渠道分析
	beego.MyAutoRouter(&controllers.ConfigsController{})         //配置管理
	beego.MyAutoRouter(&controllers.BannerController{})          //轮播图
	beego.MyAutoRouter(&controllers.CreditCardController{})      //信用卡列表
	beego.MyAutoRouter(&controllers.PushController{})            //推送历史信息管理
	beego.MyAutoRouter(&controllers.SmsManagementController{})   //推送短信
	beego.MyAutoRouter(&controllers.MsgtestController{})         //短信测试
	beego.MyAutoRouter(&controllers.PushController{})            //短信测试
	//beego.MyAutoRouter(&controllers.IdetifyController{})         //认证信息
	beego.MyAutoRouter(&controllers.MsgMarketingController{})      //短信营销
	beego.MyAutoRouter(&controllers.SMSSendingBatchesController{}) //发送报表
	beego.MyAutoRouter(&controllers.SmsReportController{})         //短信回调接口
	beego.MyAutoRouter(&controllers.ProCleaningController{})       //产品清算接口
	beego.MyAutoRouter(&controllers.CleanHisController{})          //清算历史接口
	beego.MyAutoRouter(&controllers.SettlementController{})        //结算产品接口
	beego.MyAutoRouter(&controllers.UserDepositController{})       //用户佣金提现接口
	beego.MyAutoRouter(&controllers.RegisteredController{})        //点击注册接口
	beego.MyAutoRouter(&controllers.NotifyController{})            //连连回调接口
	beego.MyAutoRouter(&controllers.HoleController{})              //新口子接口
	beego.MyAutoRouter(&controllers.PlateformHoleController{})     //征信口子接口
	beego.MyAutoRouter(&controllers.DataController{})              //数据运营
	beego.MyAutoRouter(&controllers.YaowangProductController{})    //遥望
	beego.MyAutoRouter(&controllers.MakeConnController{})          //合作顾问预约信息
	beego.MyAutoRouter(&controllers.UserDataController{})          //用户信息
	beego.MyAutoRouter(&controllers.SupportLoanController{})       //贷款稳下-付费用户
	beego.MyAutoRouter(&controllers.UsersBehaviorController{})     //贷款稳下-意愿用户
	beego.MyAutoRouter(&controllers.DataReportController{})        //数据报表
}
