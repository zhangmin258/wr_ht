package utils

import (
	"fmt"
	"zcm_tools/cache"

	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2"
)

var (
	RunMode               string //运行模式
	MYSQL_URL             string //数据库连接
	MYSQL_LOG_URL         string
	MYSQL_BACKUP_URL      string
	BEEGO_CACHE           string       //缓存地址
	Rc                    *cache.Cache //redis缓存
	Re                    error        //redis错误
	LoginVerifyCodePrefox string       // google验证码前缀
	//阿里云OSS参数
	AccessKeyId     string = "wU7F2aVGFB683FzB"
	AccessKeySecret string = "cPQJuzJno8RvKSwSbvKA0LlKFOLbZx"
	Endpoint        string = "oss-cn-hangzhou.aliyuncs.com" //阿里云OSS参数
	BucketName      string = "hzweirong"                    //存放目录名称
	Upload_dir      string = "v1/product/"                  //存放产品图片
	Imghost         string = "https://wrstatic.5ujr.cn/"

	MgoUri     string
	MgoDbName  string
	MgoSession *mgo.Session

	OutPutURL string = "http://weiyunjinrong.cn/user/register?source=" //外链的前缀

	JPushApns bool

	WR_ICON = "https://wrstatic.5ujr.cn/img/logo_share.png"

	WR_API_URL string //微融后台地址
)

var (
	PAGE_SIZE20         int = 20    //每页展示20条
	PAGE_SIZE           int = 10    //页面展示条数
	MESSAGE_INSERT_SIZE int = 10000 //给用户插入消息的时候每次插入的数量
)

func init() {

	RunMode = beego.AppConfig.String("run_mode")
	//RunMode = "release"

	config, err := beego.AppConfig.GetSection(RunMode)
	if err != nil {
		panic("配置文件读取错误 " + err.Error())
	}
	beego.Info(RunMode + "模式")
	MYSQL_URL = config["mysql_url"]
	MYSQL_LOG_URL = config["mysql_log_url"]
	MYSQL_BACKUP_URL = config["mysql_backup_url"]
	BEEGO_CACHE = config["beego_cache"]
	Rc, Re = cache.NewCache(BEEGO_CACHE) //初始化缓存

	WR_API_URL = "http://ang.s1.natapp.cc/" //微融测试后台地址
	if RunMode == "release" {
		WR_API_URL = "https://wradmin.5ujr.cn/" //微融后台地址
		//初始化Mongodb
		MgoDbName = "wjr"
		MgoUri = "mongodb://10.253.6.95:27017" // 生产阿里云
		MgoSession, err = mgo.Dial(MgoUri)
		if err != nil {
			beego.Error(err)
		}
	}

/*	//连连支付
	LL_oid_partner = "201708040000764008"
	LL_RSA_PRIVATE_KEY = `-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAKrjPDi+0rh2tUIA
VnQ8TPHOlweHcuDb4aBEkNUp8QapZFFvgFg0YG75ftUIXUFMhBYLemyLF93J7XW8
g94C6hoomGrCe1pBIHBLKudUpP8SDJYZHVz64+4e1J/qXNc1mnfKzIm2s4XMM/mz
M6UPdoyPQUdL9K1+Qr0e0XFsMdORAgMBAAECgYACSBRmwY14rwUOg4ij9qYkWFjF
3fYXsHfbtu6kGfJA34QaXj29b72V3bjmyTzNgWMGFFMnHdhMusRz3Pd5wFo1x5v/
P+ITSuWrU0g0/w537L2FJpmKxgSw3aWBf83GtWxWu29ycAULTnXoDpaT5fiMNAuD
hMaFFLC7wItGT7gkPQJBANnXpWREr5WDTuEfMlH2mBwVOBVww8716Uy2x5rxPWM3
aY7x5jW3+1qN+zWHkBh3C2ptfMql/S/r9zUs5nqoKAMCQQDI0hLvO+GujPlYoNZJ
WPoQKuRFAeuDRdL8I1x2QHFjfQyLYY0hwYQm+vsgvclSJl2RvyVpAz9P2y2xwPbn
dTPbAkEAmCakYCoRMS5rU5WEgfboOwUfDOqb+NuNPYWUWMYSCFBVq/+MuQxtxMvS
H4s1u8C5+nXKaYgSFPyMx1k7CYSVnQJAKjNbH0LqLhAZ5fIGletIwVUCGy5IG5H0
wF998qugKFQC6tdZHRrZdoePdlXrlIPTbelJJ0QzbciVVaFDQWhzuQJBAKgm3APe
eaKdQ18+nGlITRIfTu9ITcR43z8elFjH8r48FvI0xEacdDOafLpyH2ZVaUpQMQAy
miZOIfURLpR21II=
-----END PRIVATE KEY-----
`
	LL_RSA_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCq4zw4vtK4drVCAFZ0PEzxzpcHh3Lg2+GgRJDVKfEGqWRRb4BYNGBu+X7VCF1BTIQWC3psixfdye11vIPeAuoaKJhqwntaQSBwSyrnVKT/EgyWGR1c+uPuHtSf6lzXNZp3ysyJtrOFzDP5szOlD3aMj0FHS/StfkK9HtFxbDHTkQIDAQAB
-----END PUBLIC KEY-----
`
	LL_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSS/DiwdCf/aZsxxcacDnooGph3d2JOj5GXWi+
q3gznZauZjkNP8SKl3J2liP0O6rU/Y/29+IUe+GTMhMOFJuZm1htAtKiu5ekW0GlBMWxf4FPkYlQ
kPE0FtaoMP3gYfh+OwI+fIRrpW3ySn3mScnc6Z700nU/VYrRkfcSCbSnRwIDAQAB
-----END PUBLIC KEY-----
`*/

	// TODO 1.4上线时，使用如下账号。
	//连连快捷支付
	LL_oid_partner = "201710010000979003"
	LL_RSA_PRIVATE_KEY = `-----BEGIN PRIVATE KEY-----
MIICdQIBADANBgkqhkiG9w0BAQEFAASCAl8wggJbAgEAAoGBALvSN3AQKs0NIGbW
C0SL/aRiH+OF+IgLPWHFs3kMEnpW1+WjqCoJE9Nw99NNkZKEW552g8ejh3GISJYJ
3njNfqMCflKgsRzcncrZR3DSVgQM9GY5HLqDLzkAIrJ/EXCBfS/tbKzX8zJ+X5W/
MGYspn+w10YrqhjCV3wPO8kFNddxAgMBAAECgYBHsho5a+J6va0FtGU+uFWNP2u+
1XAmtmuq++XjqikPjEEDxvI1gZuQ1gm0HmMYU/AJUGJDffgA7a4PoBrNcFwLREfC
xqtrqnAfa5Ub+Xat/KPtd8hdQnC6JrEcpeLCKAjWsVLRAm/4/iinpa2xcueBv0Og
aC2iOYSRwotObdhPhQJBAN3poj8hf7l0JytBAviBmKHNOBLO4wY1RAyh2aP4bWnu
VZaJ5VnHra47ZrD2VqfvGyN2EfaprEtpPYEwtWodovcCQQDYq/0yNfAwrsMwufEQ
NjETpBVrWR49ztbHpuu3IneL0M7BAezcf5QuhjfkvRKsoC11vNBAr/enIWBC5r/j
RVbXAkBSRXDyaM/6iHaREawxR5K3weadCnieb5cH++U9ZjfiQwsWIY+XJnFcnAcp
alqcLghosDhes27+EklMISvQ6KXnAkBtEXatNdWoy/BZsOAWRxFBT9GwbfX5KwuX
CQGS+HixGvVY1v1Cqb4QBWRRcpPZ7e+0Ws2CIpJJwVVRmBJz902VAkBHmNtAe05P
xiV6/vxS53R12vVdaagCCt8sy/m0yMxZUVvzQoVMxqn3ABaG5qYaX+cNQjJ4CRI8
3IxpfMg5jWzP
-----END PRIVATE KEY-----
`
	LL_RSA_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC70jdwECrNDSBm1gtEi/2kYh/j
hfiICz1hxbN5DBJ6Vtflo6gqCRPTcPfTTZGShFuedoPHo4dxiEiWCd54zX6jAn5S
oLEc3J3K2Udw0lYEDPRmORy6gy85ACKyfxFwgX0v7Wys1/Myfl+VvzBmLKZ/sNdG
K6oYwld8DzvJBTXXcQIDAQAB
-----END PUBLIC KEY-----
`

	LL_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCSS/DiwdCf/aZsxxcacDnooGph3d2JOj5GXWi+
q3gznZauZjkNP8SKl3J2liP0O6rU/Y/29+IUe+GTMhMOFJuZm1htAtKiu5ekW0GlBMWxf4FPkYlQ
kPE0FtaoMP3gYfh+OwI+fIRrpW3ySn3mScnc6Z700nU/VYrRkfcSCbSnRwIDAQAB
-----END PUBLIC KEY-----
`
	LoginVerifyCodePrefox = "hi,dear" + RunMode[:4]
}

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func GetSession() *mgo.Session {
	if MgoSession == nil {
		var err error
		MgoSession, err = mgo.Dial(MgoUri)
		if err != nil {
			panic(err) //直接终止程序运行
		}
	} else {
		err := MgoSession.Ping()
		if err != nil {
			MgoSession, err = mgo.Dial(MgoUri)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	//最大连接池默认为4096
	return MgoSession.Clone()
}
