package models

import (
	"time"

	"github.com/astaxie/beego/orm"

	"encoding/json"
	"wr_v1/utils"
)

//用户列表展示
type UserListShow struct {
	Uid            int `orm:"column(id)"` //用户id
	Account        string                 //账号/手机号
	IdCard         string                 //身份证号码
	VerifyRealName string                 //IDnum
	CreateTime     time.Time              //创建时间
	CreditCode     float64                //信用评分
}

//用户信息展示
type UserInformation struct {
	Id             int       //用户id
	Account        string    //账号/手机号
	IdCard         string    //身份证号码
	VerifyRealName string    //IDnum
	CreateTime     time.Time //创建时间
	CreditCode     float64   //信用评分
	Source         string    //来源渠道
	Zm_score       int       //芝麻信用
	Sex            string
}

//用户收益
type UserMoney struct {
	H5Money  float64 //h5收益
	APIMoney float64 //api收益
}

//订单展示
type OrderList struct {
	Id               int
	Account          string    //用户手机号
	Verify_real_name string    //用户姓名
	Create_time      time.Time //借款时间
	Money            float64   //借款金额
	Loan_term_count  int       //借款期限
	State            string    //订单状态
}

//注册用户统计
type RegisterUser struct {
	CreateDate string
	Count      int
}

//借款金额统计
type CreditMoneyUser struct {
	CreateDate string
	Count      float32
}

//资信分析---性别
type UserGender struct {
	Num int    `json:"Count"`
	Sex string `json:"Data"`
}

//资信分析---芝麻分
type StaUser struct {
	Num   int    `json:"Count"`
	Score string `json:"Data"`
}

//资信分析---年龄
type UserAge struct {
	Num  int    `json:"Count"`
	Data string `json:"Data"`
}

//资信分析---信用评分
type UserCredit struct {
	Num    int    `json:"Count"`
	Credit string `json:"Data"`
}

//资信分析---职业身份
type UserJob struct {
	Num     int    `json:"Count"`
	Userjob string `json:"Data"`
}

//资信分析---操作系统
type UserOS struct {
	Num int    `json:"Count"`
	Os  string `json:"Data"`
}

//资信分析---运营商
type UserOperators struct {
	Num       int    `json:"Count"`
	Operators string `json:"Data"`
}

//用户贷款需求分析 --- 贷款金额
type BusinessLoanMoney struct {
	Num          int    `json:"Count"`
	MoneyAccount string `json:"Data"`
}

//用户贷款需求分析 --- 贷款期限
type BusinessLoanTermCount struct {
	Num       int    `json:"Count"`
	Termcount string `json:"Data"`
}

//用户贷款需求分析 --- 创造收益
//用户贷款需求分析 --- 贷款次数
type BusinessLoanTimes struct {
	Num   int    `json:"Count"`
	Times string `json:"Data"`
}

//用户地域分析---注册用户
type UserAddress struct {
	Provience string `json:"name"`
	Num       int    `json:"value"`
}

//获取有盾用户信息
type UserIdent struct {
	Uid          int
	Id_name      string
	Id_no        string
	Address      string
	Nation       string
	Gender       string
	Start_card   string
	Front_card   string
	Back_card    string
	Photo_get    string
	Photo_grid   string
	Photo_living string
	Fail_reason  string
	Be_idcard    string
}

//数据明细
type DetailData struct {
	CreateDate    string //日期
	AddedCount    int    //新增用户数
	RegisterCount int    //注册用户数
	ActiveCount   int    //活跃用户数
	IdentifyCount int    //认证用户数
	CreditCount   int    // 放款用户数
}

//注册用户明细
type RegisterData struct {
	CreateDate    string
	RegisterCount int
}

//活跃用户明细
type ActiveData struct {
	CreateDate  string
	ActiveCount int
}

//认证用户明细
type IdentifyData struct {
	CreateDate    string
	IdentifyCount int
}

//放款用户明细
type CreditData struct {
	CreateDate  string
	CreditCount int
}

//累计授信
func GetCreditExtensionCount(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM product_register_user  pru
	INNER JOIN users_business ub
	ON pru.uid=ub.uid
	WHERE pru.product_id=? AND ub.is_credit = 1`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

// 累计注册
func GetRegisterCount(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM product_register_user  WHERE product_id=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

//获取最早注册的时间
func GetEarliestDateFromUsersByRegister() (earlisetDate string, err error) {
	sql := ` SELECT create_time FROM users ORDER BY create_time LIMIT 1 `
	err = orm.NewOrm().Raw(sql).QueryRow(&earlisetDate)
	return
}

// 累计认证
func GetIdentifyCount() (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_auth ua
	INNER JOIN users u ON u.id=ua.uid
	WHERE ua.is_real_name = 2 AND ua.is_user_data = 2 AND ua.is_zm_auth = 2  AND ua.is_link_man=2 `
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

// 累计申请
func GetApplyCount() (count int, err error) {
	sql := `SELECT COUNT(1) FROM users`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

// 累计认证
func GetAuthCount(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM product_register_user pru INNER JOIN users_auth ua
	ON pru.uid=ua.uid
	WHERE ua.is_real_name = 2 AND ua.is_user_data = 2 AND ua.is_zm_auth = 2  AND ua.is_link_man=2 AND  pru.product_id=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

// 累计放款
func GetCreditCount(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM product_register_user pru
	INNER JOIN
			 business_loan bl
			ON bl.uid = pru.uid
			WHERE (bl.state = 'CONFIRM' OR bl.state = 'FINISH' OR bl.state = 'OVERDUE') AND pru.product_id=? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

// 分页查询用户数据
func GetUserList(condition string, params []string, begin, size int) (users []UserListShow, err error) {
	sql := `SELECT u.id, u.account, um.verify_real_name, u.create_time, um.id_card, um.credit_code FROM users_metadata um RIGHT JOIN users u ON u.id = um.uid WHERE 1 = 1`
	if condition != "" {
		sql += condition
	}
	sql += " order by create_time DESC limit ?, ?"
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&users)

	if err != nil {
		return nil, err
	}
	return
}

// 查询所有用户数量
func GetUserCount(condition string, params []string) (count int, err error) {
	sql := `SELECT count(1) FROM users u LEFT JOIN users_metadata um ON u.id=um.uid WHERE 1=1`
	if condition != "" {
		sql += condition
	} else {
		sql = `SELECT count(1) FROM users`
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	if err != nil {
		return 0, err
	}
	return
}

/**
根据ID查询用户
*/
func GetUserById(id int) (users *UserInformation, err error) {
	sql := `SELECT u.id, u.account, um.verify_real_name,um.sex, u.create_time, um.id_card, um.credit_code,u.source,um.zm_score FROM users_metadata um RIGHT JOIN users u ON u.id = um.uid WHERE u.id=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&users)
	return
}

func GetUserMoney(id, ctype int) (money float64, err error) {
	sql := `SELECT SUM(pc.cpa_price)
	FROM product_register_user pru
	LEFT JOIN product_cleaning pc ON pru.product_id = pc.product_id
	LEFT JOIN product p ON pru.product_id = p.id
	WHERE pru.uid= ? and p.cooperation_type= ?`
	err = orm.NewOrm().Raw(sql, id, ctype).QueryRow(&money)
	return
}

//获取累计申请借款用户数
func GetLoanTotal(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM product_register_user  pru
	INNER JOIN business_loan bl
	ON pru.uid=bl.uid WHERE pru.product_id=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

//注册用户明细
func GetRegisterData(condition string, params []interface{}, begin int, size int) (detailData []RegisterData, err error) {
	sql := `SELECT DATE_FORMAT(create_time,'%Y-%m-%d')create_date,COUNT(1) register_count FROM users WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY create_date ORDER BY create_date DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, params, params, params, begin, size).QueryRows(&detailData)
	return
}

//认证用户明细
func GetIdentifyData(condition string, params []interface{}, begin int, size int) (detailData []IdentifyData, err error) {
	sql := `SELECT DATE_FORMAT(u.create_time,'%Y-%m-%d')create_date,COUNT(1) identify_count FROM users AS u LEFT  JOIN users_auth AS ua on u.id=ua.uid
			WHERE ua.is_real_name = 2 AND ua.is_user_data = 2 AND ua.is_zm_auth = 2 AND ua.is_link_man=2 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY create_date ORDER BY create_date DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, params, params, params, begin, size).QueryRows(&detailData)
	return
}

//放款用户明细
func GetCreditData(condition string, params []interface{}, begin int, size int) (detailData []CreditData, err error) {
	sql := `SELECT DATE_FORMAT(u.create_time,'%Y-%m-%d')create_date,COUNT(1) credit_count FROM users AS u INNER JOIN business_loan AS bl ON u.id=bl.uid
		WHERE (bl.state='CONFIRM' OR bl.state='FINISH' OR bl.state='OVERDUE') `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY create_date ORDER BY create_date DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, params, params, params, begin, size).QueryRows(&detailData)
	return
}

//按天查看数据明细
func GetDailyDatas(begin int, size int) (detailData []DetailData, err error) {
	sql := `SELECT register.create_date,register_count,identify_count,credit_count FROM
		    (SELECT create_date ,COUNT(1) register_count FROM users GROUP BY create_date) register
			LEFT JOIN
			(SELECT u.create_date dayi,COUNT(1) identify_count FROM users AS u LEFT  JOIN users_auth AS ua on u.id=ua.uid
			WHERE ua.is_real_name = 2 AND ua.is_user_data = 2 AND ua.is_zm_auth = 2 AND ua.is_link_man=2 GROUP BY dayi) identify
			ON register.create_date=identify.dayi
			LEFT JOIN
			(SELECT DATE_FORMAT(bl.real_time,"%Y-%m-%d") dayc , COUNT(1) credit_count
			FROM business_loan bl
			LEFT JOIN users u ON bl.uid = u.id
			WHERE (bl.state = 'CONFIRM' OR bl.state = 'FINISH' OR bl.state = 'OVERDUE')
			GROUP BY dayc
			) credit
			ON register.create_date=credit.dayc
			ORDER BY register.create_date DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, begin, size).QueryRows(&detailData)
	return

}

//查询数据总数
func GetDailyDatasCount() (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT u.create_date AS day1,COUNT(1) AS register_count FROM users u GROUP BY day1) AS a`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取注册/申请用户数
func GetRegisterUsers() (rus []RegisterUser, err error) {
	sql := `SELECT  u.create_date,COUNT(1) count FROM users u
	        GROUP BY  u.create_date
        	ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	return
}

//获取活跃用户数量
func GetActiveUsers() (rus []RegisterUser, err error) {
	sql := `SELECT DATE(create_time)create_date, COUNT(DISTINCT uid) count
			FROM login_record
            GROUP BY  create_date
			ORDER BY  create_date ASC`
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err = o.Raw(sql).QueryRows(&rus)
	return
}

//获取认证用户数
func GetIdentifyUsers() (rus []RegisterUser, err error) {
	sql := `SELECT u.create_date,COUNT(1) count FROM users AS u LEFT  JOIN users_auth AS ua on u.id=ua.uid
			WHERE ua.is_real_name = 2 AND ua.is_user_data = 2 AND ua.is_zm_auth = 2 AND ua.is_link_man=2
			GROUP BY  u.create_date
			ORDER BY  u.create_date `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	return
}

//获取申请贷款用户数
func GetLoanUsers() (rus []RegisterUser, err error) {
	sql := ` SELECT u.create_date,COUNT(1) count FROM business_loan bl LEFT JOIN users u
             ON u.id=bl.uid
             GROUP BY  u.create_date
			 ORDER BY  u.create_date `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	return
}

//获取放款用户数 TODO 该语句在business_loan数据较大时可能会炸
func GetCreditUsers() (rus []RegisterUser, err error) {
	sql := `SELECT u.create_date,COUNT(1) count FROM users AS u INNER JOIN business_loan AS bl
			ON u.id=bl.uid
			WHERE bl.state='CONFIRM' OR bl.state='FINISH' OR bl.state='OVERDUE'
			GROUP BY  u.create_date
			ORDER BY  u.create_date`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	return
}

/**
根据条件获取注册用户
*/
func GetRegisterUsersByCondition(id, pageStart, pageSize int) (rus []RegisterUser, err error) {
	sql := `SELECT  DATE_FORMAT(create_time,'%Y-%m-%d') create_date,COUNT(1) count FROM product_register_user  WHERE product_id =? GROUP BY create_date
			 ORDER BY  create_date  DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, id, pageStart, pageSize).QueryRows(&rus)
	return
}

/**
根据条件获取注册用户
*/
func GetRegisterUsersByCondition2(id int, condition string, params []string) (rus []RegisterUser, err error) {
	sql := `SELECT  DATE_FORMAT(pru.create_time,'%Y-%m-%d') create_date,COUNT(1) count FROM product_register_user pru  WHERE product_id =? `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY create_date`
	_, err = orm.NewOrm().Raw(sql, id, params).QueryRows(&rus)
	return
}
func GetRegisterUsersCountByCondition(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT DATE_FORMAT(create_time,'%Y-%m-%d') create_date FROM product_register_user WHERE product_id =? GROUP BY create_date) a`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

/**
根据条件获取申请用户
*/
func GetLoanUsersByCondition(condition string, params []string) (lus []RegisterUser, err error) {
	sql := `SELECT  DATE_FORMAT(u.create_time,'%Y-%m-%d') create_date,COUNT(1) count FROM users u WHERE 1 = 1`
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(u.create_time,'%Y-%m-%d')
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&lus)
	return
}

/**
根据条件获取认证用户
*/
func GetAuthsersByCondition(id int, condition string, params []string) (lus []RegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(ua.link_man_time,'%Y-%m-%d') create_date,COUNT(1) count FROM product_register_user pru INNER JOIN users_auth ua
			ON pru.uid=ua.uid
			WHERE ua.is_real_name = 2 AND ua.is_user_data = 2 AND ua.is_zm_auth = 2  AND ua.is_link_man=2
			AND DATE(ua.real_name_time)= DATE(ua.link_man_time)
			AND DATE(ua.user_data_time)= DATE(ua.link_man_time)
			AND DATE(ua.zm_auth_time)= DATE(ua.link_man_time)
			AND  pru.product_id=? `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(ua.link_man_time,'%Y-%m-%d')`
	_, err = orm.NewOrm().Raw(sql, id, params).QueryRows(&lus)
	return
}

/**
根据条件获取授信用户
*/
func GetCreditExtensionUsersByCondition(id int, condition string, params []string) (ceus []RegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(ub.auth_time,'%Y-%m-%d') create_date,COUNT(1) count FROM product_register_user  pru
		INNER JOIN users_business ub
		ON pru.uid=ub.uid
		WHERE pru.product_id=? AND ub.is_credit = 1 `

	/*sql := `SELECT  DATE_FORMAT(u.create_time,'%Y-%m-%d') create_date,COUNT(1) count FROM users u
	INNER JOIN users_business ub ON u.id = ub.uid WHERE ub.is_credit = 1 `*/
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(ub.auth_time,'%Y-%m-%d')
			 ORDER BY  create_date DESC `
	_, err = orm.NewOrm().Raw(sql, id, params).QueryRows(&ceus)
	return
}

/**
根据条件获取申请借款用户
*/
func GetLoanTotalUsersByCondition(id int, condition string, params []string) (ltus []RegisterUser, err error) {
	sql := `SELECT  DATE_FORMAT(bl.loan_time,'%Y-%m-%d') create_date,COUNT(1) count FROM product_register_user  pru
		INNER JOIN business_loan bl
		ON pru.uid=bl.uid WHERE pru.product_id=? `

	/*sql := `SELECT  DATE_FORMAT(bl.loan_time,'%Y-%m-%d') create_date,COUNT(1) count FROM users u
	INNER JOIN business_loan bl ON u.id = bl.uid WHERE 1 = 1 `*/
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(bl.loan_time,'%Y-%m-%d')
			 ORDER BY  create_date ASC `
	_, err = orm.NewOrm().Raw(sql, id, params).QueryRows(&ltus)
	return
}

/**
根据条件获取放款用户
*/
func GetCreditUsersByCondition(id int, condition string, params []string) (cus []RegisterUser, err error) {
	sql := `SELECT  DATE_FORMAT(bl.real_time,'%Y-%m-%d') create_date,COUNT(1) count FROM product_register_user  pru
			INNER JOIN business_loan bl ON pru.uid = bl.uid WHERE (bl.state='CONFIRM' OR bl.state='OVERDUE' OR bl.state='FINISH') AND pru.product_id=? `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(bl.real_time,'%Y-%m-%d')
			 ORDER BY  create_date DESC `
	_, err = orm.NewOrm().Raw(sql, id, params).QueryRows(&cus)
	return
}

/**
根据条件获取放款金额
*/
func GetCreditMoneyUsersByCondition(id int, condition string, params []string) (cmu []CreditMoneyUser, err error) {
	sql := `SELECT  DATE_FORMAT(bl.real_time,'%Y-%m-%d') create_date,SUM(bl.real_money) count FROM product_register_user  pru
			INNER JOIN business_loan bl ON pru.uid = bl.uid WHERE (bl.state='CONFIRM' OR bl.state='OVERDUE' OR bl.state='FINISH') AND pru.product_id=? `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(bl.real_time,'%Y-%m-%d')
			 ORDER BY  create_date DESC `
	_, err = orm.NewOrm().Raw(sql, id, params).QueryRows(&cmu)
	return
}

/*
根据条件获取订单列表指定日期的数据
*/

func GetOrderlist(condition string, params []string, begin, size int) (list []OrderList, err error) {

	sql := `SELECT l.id,um.account,um.verify_real_name,l.create_time,l.money,l.loan_term_count,l.state FROM 
	business_loan l LEFT JOIN users_metadata um ON l.uid=um.uid left join product p on p.id=l.product_code WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += " order by l.create_time DESC limit ?, ?"
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&list)
	return
}

/*
根据条件获取用户订单详情
*/
func GetOrderlistCount(condition string, params []string) (c int) {

	sql := `SELECT count(1)  FROM business_loan l LEFT JOIN users_metadata um ON l.uid=um.uid left join product p on p.id=l.product_code WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	_ = orm.NewOrm().Raw(sql, params).QueryRow(&c)

	return
}

//用户数据分析——性别
func StaGender() (list []UserGender, err error) {
	o := orm.NewOrm()
	sql := ` SELECT COUNT(1) AS num, CASE 
			WHEN um.sex = '男' THEN "男" 
			WHEN um.sex = '女' THEN "女" 
			ELSE  "未认证"
			END AS sex
			FROM users_metadata  um
			GROUP BY sex `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

//用户数据分析——芝麻分
func StaZMScore() (list []StaUser, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num, CASE
            WHEN zm_score =0 THEN "未认证"
            WHEN zm_score >0 AND zm_score <=499 THEN "500分以下"
            WHEN zm_score>=500 AND zm_score<=549 THEN "500-549分"
            WHEN zm_score>=550 AND zm_score<=599 THEN "550-599分"
            WHEN zm_score>=600 AND zm_score<=649 THEN "600-649分"
            WHEN zm_score>=650 AND zm_score<=699 THEN "650-699分"
            WHEN zm_score>=700 THEN "700分以上" END AS score
            FROM users_metadata  um
            GROUP BY score  `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户数据分析——年龄
func StaAge() (list []UserAge, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num,
			CASE
			WHEN age > 0 AND age <18 THEN "18岁以下"
			WHEN age >= 18 AND age <=20 THEN "18-20岁"
			WHEN age >= 21 AND age <= 22 THEN "21-22岁"
			WHEN age >= 23 AND age <= 35 THEN "23-35岁"
			WHEN age >= 36 AND age <= 40 THEN "36-40岁"
			WHEN age >= 41 AND age <= 55 THEN "41-55岁"
			WHEN age >= 56 THEN "55以上" 
			WHEN age =0 THEN "未认证"
			END AS data
			FROM users_metadata
			GROUP BY data
			ORDER BY age `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户数据分析——信用评分
func StaCredit() (list []UserCredit, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num,
			CASE
			WHEN credit_code = 0.0 THEN "未认证"
			WHEN credit_code > 0 AND credit_code <45 THEN "信用评分低于45分" 
			WHEN credit_code >= 45 AND credit_code <=50 THEN "信用评分45-50分" 
			WHEN credit_code >= 51 AND credit_code <=60 THEN "信用评分51-60分" 
			WHEN credit_code >= 61 AND credit_code <=75 THEN "信用评分61-75分" 
			WHEN credit_code >= 76 THEN "信用评分高于75分" END AS credit
			from users_metadata 
			GROUP BY credit  `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户数据分析——职业身份
func StaJob() (list []UserJob, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num,
			CASE 
			WHEN ub.identity = 1 THEN "上班族" 
			WHEN ub.identity = 2 THEN "个体户" 
			WHEN ub.identity = 3 THEN "无固定职业" 
			WHEN ub.identity = 4 THEN "企业主" 
			WHEN ub.identity = 5 THEN "学生"
			ELSE "未认证" 
			END AS userjob 
			FROM 
			users_basedata ub
			LEFT JOIN job j ON ub.identity = j.id
			GROUP BY userjob `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户数据分析——操作系统
func StaOS() (list []UserOS, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num,
            CASE
            WHEN u.app = 1 THEN "ios"
            WHEN u.app = 2 THEN "android"
            WHEN u.app = 3 THEN "wx"
            WHEN u.app = 4 THEN "pc"
            WHEN u.app = 5 THEN "wd"
            WHEN u.app = 6 THEN "wap"
            WHEN u.app = 7 THEN "落地页"
            ELSE "未激活"
            END AS os
            FROM users u
            GROUP BY os
            ORDER BY u.app DESC  `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户数据分析——运营商
func StaOperators() (list []UserOperators, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num,
			CASE 
			WHEN TRIM(u.account) REGEXP "^1([3][4-9]|[4][7]|[5][0-27-9]|[8][2-478])[0-9]{8}$" THEN '移动' 
			WHEN TRIM(u.account) REGEXP "^1([3][0-2]|[4][5]|[5][5-6]|[7][6]|[8][5-6])[0-9]{8}$" THEN '联通' 
			WHEN TRIM(u.account) REGEXP "^1(3[3]|5[3]|7[37]|8[019])[0-9]{8}$" THEN '电信' 
			ELSE '其他运营商' 
			END AS operators
			FROM users u
			GROUP BY operators DESC `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

//获取所有用户手机号
func GetAccounts() (accounts []string, err error) {
	sql := `SELECT account FROM users `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&accounts)
	return
}

// 用户贷款需求分析---贷款金额
func StaLoanMoney() (list []BusinessLoanMoney, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num, 
			CASE 
			WHEN  ub.loan_amount < 1000 THEN "1000以下" 
			WHEN ub.loan_amount >= 1000 AND ub.loan_amount <= 2999 THEN "1000-2999" 
			WHEN ub.loan_amount >= 3000 AND ub.loan_amount <= 4999 THEN "3000-4999" 
			WHEN ub.loan_amount >= 5000 AND ub.loan_amount <= 7999 THEN "5000-7999" 
			WHEN ub.loan_amount >= 8000 AND ub.loan_amount <= 9999 THEN "8000-10000" 		
			WHEN ub.loan_amount >= 10000 AND ub.loan_amount <= 29999 THEN "1~3万" 
			WHEN ub.loan_amount >= 30000 THEN "3万以上"
			WHEN ub.loan_amount IS NULL THEN "暂无数据"
			END AS money_account
			FROM users_basedata ub
			GROUP BY money_account
			ORDER BY money_account `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户贷款需求分析---贷款期限
func StaLoanTermCount() (list []BusinessLoanTermCount, err error) {
	o := orm.NewOrm()
	sql := ` SELECT COUNT(1) AS num, CASE 
			WHEN ub.loan_term = 7 THEN "7天" 
			WHEN ub.loan_term = 14 THEN "14天" 
			WHEN ub.loan_term = 21 THEN "21天" 
			WHEN ub.loan_term = 28 THEN "28天" 
			WHEN ub.loan_term = 30 THEN "1个月" 
			WHEN ub.loan_term = 90 THEN "3个月" 
			WHEN ub.loan_term = 180 THEN "6个月" 
			END AS termcount
			FROM users_basedata ub
			GROUP BY termcount
			HAVING  termcount != ""
			ORDER BY termcount  `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户贷款需求分析---借款次数
func StaBusinessLoanTimes() (list []BusinessLoanTimes, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(id) AS num,'0次' AS times
         FROM users
         WHERE id NOT IN (SELECT DISTINCT uid FROM product_register_user)
			UNION 	
			(SELECT COUNT(1) count2,CASE
			WHEN ts = 1 THEN "1次"
			WHEN ts = 2 THEN "2次"
			WHEN ts = 3 THEN "3次"
			WHEN ts = 4 THEN "4次"
			WHEN ts = 5 THEN "5次"
			WHEN ts > 5 THEN "5次以上"
			END AS times
			FROM(
			SELECT COUNT(1) ts
			FROM product_register_user
			GROUP BY uid) temp
			GROUP BY times) `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户地域分析--- 注册成功
func StaRegisterUserAddress() (list []UserAddress, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num, CASE WHEN
			LEFT(u.address,4) REGEXP '北京' THEN '北京' WHEN
			LEFT(u.address,4) REGEXP '天津' THEN '天津' WHEN
			LEFT(u.address,4) REGEXP '上海' THEN '上海' WHEN
			LEFT(u.address,4) REGEXP '重庆' THEN '重庆' WHEN
			LEFT(u.address,4) REGEXP '河北' THEN '河北' WHEN
			LEFT(u.address,4) REGEXP '山西' THEN '山西' WHEN
			LEFT(u.address,4) REGEXP '台湾' THEN '台湾' WHEN
			LEFT(u.address,4) REGEXP '辽宁' THEN '辽宁' WHEN
			LEFT(u.address,4) REGEXP '吉林' THEN '吉林' WHEN
			LEFT(u.address,5) REGEXP '黑龙江' THEN '黑龙江' WHEN
			LEFT(u.address,4) REGEXP '江苏' THEN '江苏' WHEN
			LEFT(u.address,4) REGEXP '浙江' THEN '浙江' WHEN
			LEFT(u.address,4) REGEXP '安徽' THEN '安徽' WHEN
			LEFT(u.address,4) REGEXP '福建' THEN '福建' WHEN
			LEFT(u.address,4) REGEXP '江西' THEN '江西' WHEN
			LEFT(u.address,4) REGEXP '山东' THEN '山东' WHEN
			LEFT(u.address,4) REGEXP '河南' THEN '河南' WHEN
			LEFT(u.address,4) REGEXP '湖北' THEN '湖北' WHEN
			LEFT(u.address,4) REGEXP '湖南' THEN '湖南' WHEN
			LEFT(u.address,4) REGEXP '广东' THEN '广东' WHEN
			LEFT(u.address,4) REGEXP '甘肃' THEN '甘肃' WHEN
			LEFT(u.address,4) REGEXP '四川' THEN '四川' WHEN
			LEFT(u.address,4) REGEXP '贵州' THEN '贵州' WHEN
			LEFT(u.address,4) REGEXP '海南' THEN '海南' WHEN
			LEFT(u.address,4) REGEXP '云南' THEN '云南' WHEN
			LEFT(u.address,4) REGEXP '青海' THEN '青海' WHEN
			LEFT(u.address,4) REGEXP '陕西' THEN '陕西' WHEN
			LEFT(u.address,4) REGEXP '广西' THEN '广西' WHEN
			LEFT(u.address,4) REGEXP '西藏' THEN '西藏' WHEN
			LEFT(u.address,4) REGEXP '宁夏' THEN '宁夏' WHEN
			LEFT(u.address,4) REGEXP '新疆' THEN '新疆' WHEN
			LEFT(u.address,5) REGEXP '内蒙古' THEN '内蒙古' WHEN
			LEFT(u.address,4) REGEXP '澳门' THEN '澳门' WHEN
			LEFT(u.address,4) REGEXP '香港' THEN '香港' END AS provience
			FROM users u
			GROUP BY provience HAVING  provience IS NOT NULL `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户地域分析--- 认证成功
func StaApprovedUserAddress() (list []UserAddress, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num, CASE WHEN
LEFT(u.address,4) = '中国北京' or LEFT(u.address,2) = '北京' then '北京' WHEN
LEFT(u.address,4) = '中国天津' or LEFT(u.address,2) = '天津'then'天津' WHEN
LEFT(u.address,4) = '中国上海' or LEFT(u.address,2) = '上海' then'上海' WHEN
LEFT(u.address,4) = '中国重庆' or LEFT(u.address,2) = '重庆'then '重庆' WHEN
LEFT(u.address,4) = '中国河北' or LEFT(u.address,2) = '河北' then'河北' WHEN
LEFT(u.address,4) = '中国山西' or LEFT(u.address,2) = '山西' then'山西' WHEN
LEFT(u.address,4) = '中国台湾' or LEFT(u.address,2) = '台湾'then '台湾' WHEN
LEFT(u.address,4) = '中国辽宁' or LEFT(u.address,2) = '辽宁' then'辽宁' WHEN
LEFT(u.address,4) = '中国吉林' or LEFT(u.address,2) = '吉林' then'吉林' WHEN
LEFT(u.address,5) = '中国黑龙江' or LEFT(u.address,3) = '黑龙江' then'黑龙江' WHEN
LEFT(u.address,4) = '中国江苏' or LEFT(u.address,2) = '江苏' then'江苏' WHEN
LEFT(u.address,4) = '中国浙江' or LEFT(u.address,2) = '浙江'then '浙江' WHEN
LEFT(u.address,4) = '中国安徽' or LEFT(u.address,2) = '安徽' then'安徽' WHEN
LEFT(u.address,4) = '中国福建' or LEFT(u.address,2) = '福建'then '福建' WHEN
LEFT(u.address,4) = '中国江西' or LEFT(u.address,2) = '江西'then '江西' WHEN
LEFT(u.address,4) = '中国山东' or LEFT(u.address,2) = '山东' then'山东' WHEN
LEFT(u.address,4) = '中国河南' or LEFT(u.address,2) = '河南'then '河南' WHEN
LEFT(u.address,4) = '中国湖北' or LEFT(u.address,2) = '湖北' then'湖北' WHEN
LEFT(u.address,4) = '中国湖南' or LEFT(u.address,2) = '湖南'then '湖南' WHEN
LEFT(u.address,4) = '中国广东' or LEFT(u.address,2) = '广东'then '广东' WHEN
LEFT(u.address,4) = '中国甘肃' or LEFT(u.address,2) = '甘肃'then '甘肃' WHEN
LEFT(u.address,4) = '中国四川' or LEFT(u.address,2) = '四川' then'四川' WHEN
LEFT(u.address,4) = '中国贵州' or LEFT(u.address,2) = '贵州'then '贵州' WHEN
LEFT(u.address,4) = '中国海南' or LEFT(u.address,2) = '海南' then'海南' WHEN
LEFT(u.address,4) = '中国云南' or LEFT(u.address,2) = '云南' then'云南' WHEN
LEFT(u.address,4) = '中国青海' or LEFT(u.address,2) = '青海' then'青海' WHEN
LEFT(u.address,4) = '中国陕西' or LEFT(u.address,2) = '陕西'then '陕西' WHEN
LEFT(u.address,4) = '中国广西' or LEFT(u.address,2) = '广西' then'广西' WHEN
LEFT(u.address,4) = '中国西藏' or LEFT(u.address,2) = '西藏' then'西藏' WHEN
LEFT(u.address,4) = '中国宁夏' or LEFT(u.address,2) = '宁夏' then'宁夏' WHEN
LEFT(u.address,4) = '中国新疆'  or LEFT(u.address,2) = '新疆'then '新疆' WHEN
LEFT(u.address,5) = '中国内蒙古' or LEFT(u.address,3) = '内蒙古'then '内蒙古' WHEN
LEFT(u.address,4) = '中国澳门'  or LEFT(u.address,2) ='澳门'then '澳门' WHEN
LEFT(u.address,4) = '中国香港' '香港' or LEFT(u.address,2) = '香港'then  '香港' END AS provience
FROM users_auth ua
INNER JOIN users u ON ua.uid = u.id
WHERE ua.is_real_name = 2
GROUP BY provience HAVING  provience IS NOT NULL `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

// 用户地域分析--- 借款成功
func StaLoanedUserAddress() (list []UserAddress, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) AS num,b.uid, CASE WHEN
LEFT(u.address,4) = '中国北京' or LEFT(u.address,2) = '北京' then '北京' WHEN
LEFT(u.address,4) = '中国天津' or LEFT(u.address,2) = '天津'then'天津' WHEN
LEFT(u.address,4) = '中国上海' or LEFT(u.address,2) = '上海' then'上海' WHEN
LEFT(u.address,4) = '中国重庆' or LEFT(u.address,2) = '重庆'then '重庆' WHEN
LEFT(u.address,4) = '中国河北' or LEFT(u.address,2) = '河北' then'河北' WHEN
LEFT(u.address,4) = '中国山西' or LEFT(u.address,2) = '山西' then'山西' WHEN
LEFT(u.address,4) = '中国台湾' or LEFT(u.address,2) = '台湾'then '台湾' WHEN
LEFT(u.address,4) = '中国辽宁' or LEFT(u.address,2) = '辽宁' then'辽宁' WHEN
LEFT(u.address,4) = '中国吉林' or LEFT(u.address,2) = '吉林' then'吉林' WHEN
LEFT(u.address,5) = '中国黑龙江' or LEFT(u.address,3) = '黑龙江' then'黑龙江' WHEN
LEFT(u.address,4) = '中国江苏' or LEFT(u.address,2) = '江苏' then'江苏' WHEN
LEFT(u.address,4) = '中国浙江' or LEFT(u.address,2) = '浙江'then '浙江' WHEN
LEFT(u.address,4) = '中国安徽' or LEFT(u.address,2) = '安徽' then'安徽' WHEN
LEFT(u.address,4) = '中国福建' or LEFT(u.address,2) = '福建'then '福建' WHEN
LEFT(u.address,4) = '中国江西' or LEFT(u.address,2) = '江西'then '江西' WHEN
LEFT(u.address,4) = '中国山东' or LEFT(u.address,2) = '山东' then'山东' WHEN
LEFT(u.address,4) = '中国河南' or LEFT(u.address,2) = '河南'then '河南' WHEN
LEFT(u.address,4) = '中国湖北' or LEFT(u.address,2) = '湖北' then'湖北' WHEN
LEFT(u.address,4) = '中国湖南' or LEFT(u.address,2) = '湖南'then '湖南' WHEN
LEFT(u.address,4) = '中国广东' or LEFT(u.address,2) = '广东'then '广东' WHEN
LEFT(u.address,4) = '中国甘肃' or LEFT(u.address,2) = '甘肃'then '甘肃' WHEN
LEFT(u.address,4) = '中国四川' or LEFT(u.address,2) = '四川' then'四川' WHEN
LEFT(u.address,4) = '中国贵州' or LEFT(u.address,2) = '贵州'then '贵州' WHEN
LEFT(u.address,4) = '中国海南' or LEFT(u.address,2) = '海南' then'海南' WHEN
LEFT(u.address,4) = '中国云南' or LEFT(u.address,2) = '云南' then'云南' WHEN
LEFT(u.address,4) = '中国青海' or LEFT(u.address,2) = '青海' then'青海' WHEN
LEFT(u.address,4) = '中国陕西' or LEFT(u.address,2) = '陕西'then '陕西' WHEN
LEFT(u.address,4) = '中国广西' or LEFT(u.address,2) = '广西' then'广西' WHEN
LEFT(u.address,4) = '中国西藏' or LEFT(u.address,2) = '西藏' then'西藏' WHEN
LEFT(u.address,4) = '中国宁夏' or LEFT(u.address,2) = '宁夏' then'宁夏' WHEN
LEFT(u.address,4) = '中国新疆'  or LEFT(u.address,2) = '新疆'then '新疆' WHEN
LEFT(u.address,5) = '中国内蒙古' or LEFT(u.address,3) = '内蒙古'then '内蒙古' WHEN
LEFT(u.address,4) = '中国澳门'  or LEFT(u.address,2) ='澳门'then '澳门' WHEN
LEFT(u.address,4) = '中国香港' '香港' or LEFT(u.address,2) = '香港'then  '香港' END AS provience
FROM business_loan b
INNER JOIN users_auth ua ON ua.uid = b.uid
INNER JOIN users u ON b.uid = u.id
WHERE b.state REGEXP 'CONFIRM' OR b.state = 'FINISH' OR b.state = 'OVERDUE' AND ua.is_real_name=2
GROUP BY provience HAVING  provience IS NOT NULL `
	_, err = o.Raw(sql).QueryRows(&list)
	return
}

//获取有盾用户信息

func GetIdentificationUser(id int) (UserIdent *UserIdent) {
	sql := `SELECT i.uid,i.id_name,i.id_no,i.address,i.nation,i.gender,i.start_card,i.front_card,i.back_card,i.photo_get,i.photo_grid,i.photo_living,i.fail_reason,i.be_idcard 
	FROM identification i WHERE i.uid=?`

	_ = orm.NewOrm().Raw(sql, id).QueryRow(&UserIdent)
	return
}

//邀请好友
type InviteFriend struct {
	Id         int       // id
	Account    string    // 手机号
	Name       string    // 姓名
	IdCard     string    // 身份证号
	CreditCode int       // 信用积分
	InviteId   int       //邀请人id
	CreateTime time.Time // 注册时间
}

//查询 用户邀请好友列表
func GetInviteFriendList(id, start, size int) (inviteFriendList []*InviteFriend, err error) {
	sql := `SELECT i.old_uid AS id,
			u.account,
			um.verify_real_name AS name,
			um.id_card,
			um.credit_code,
			i.new_uid AS invite_id,
			u.create_time
			FROM invitation_record i
			INNER JOIN users u
			ON u.id = i.new_uid 
			LEFT JOIN users_metadata um
			ON i.new_uid = um.uid
			WHERE i.old_uid = ?
			ORDER BY u.create_time DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, id, start, size).QueryRows(&inviteFriendList)
	return
}

//查询 用户邀请好友列表的数量
func GetInviteFriendListCount(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM invitation_record i INNER JOIN users u ON  u.id = i.new_uid LEFT JOIN users_metadata um ON um.uid = i.new_uid WHERE i.old_uid = ? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

// 网贷记账
type LoanAccount struct {
	Id              int       // id
	Uid             int       // 用户id
	PlatformName    string    // 平台名称
	RepaymentDay    int       // 还款日
	TotalPeriods    int       // 总期数
	CurrentPeriod   int       // 当前期数
	RepayEachPeriod float64   //每期应还
	CreateTime      time.Time // 借款时间`orm:"type(date)"`
}

//查询 用户网贷记账列表
func GetLoanAccountList(id, start, size int) (loanAccountList []*LoanAccount, err error) {
	sql := `SELECT id,uid,platform_name,repayment_day,total_periods,current_period,repay_each_period, create_time FROM loan_account WHERE uid=? LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, id, start, size).QueryRows(&loanAccountList)
	return
}

//查询 用户网贷记账列表的数量
func GetLoanAccountListCount(id int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM loan_account WHERE uid=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&count)
	return
}

//紧急联系人
type EmergencyContact struct {
	Id                 int       //id
	Uid                int       //用户id
	Relation           string    //关系
	LinkmanName        string    //紧急联系人
	LinkmanPhoneNumber string    //紧急联系人电话号码
	IsNormal           int       //紧急联系人电话是否异常
	NormalReason       string    //异常原因
	CreateTime         time.Time //创建时间
}

//查询紧急联系人
func GetEmergencyContactList(uid, begin, size int) (contact []EmergencyContact, err error) {
	sql := `SELECT uid,relation,linkman_name,linkman_phone_number,is_normal,normal_reason,create_time
	      FROM users_linkman
	      WHERE uid=?
	      ORDER BY create_time DESC
	      LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, uid, begin, size).QueryRows(&contact)
	return
}

//查询紧急联系人数量
func GetContactCount(uid int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_linkman WHERE uid=?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&count)
	return
}

//活跃用户明细
func GetActiveData(begin int, size int) (detailData []ActiveData, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT DATE(create_time) as create_date, COUNT(DISTINCT uid) as active_count
			FROM login_record
			WHERE 1= 1  GROUP BY create_date ORDER BY create_date DESC LIMIT ?,?`
	_, err = o.Raw(sql, begin, size).QueryRows(&detailData)
	return
}

func ProvienceInit() (s []UserAddress) {
	bj := UserAddress{Provience: "北京"}
	tj := UserAddress{Provience: "天津"}
	sh := UserAddress{Provience: "上海"}
	cq := UserAddress{Provience: "重庆"}
	hb := UserAddress{Provience: "河北"}
	sx := UserAddress{Provience: "山西"}
	tw := UserAddress{Provience: "台湾"}
	ll := UserAddress{Provience: "辽宁"}
	jl := UserAddress{Provience: "吉林"}
	hlj := UserAddress{Provience: "黑龙江"}
	js := UserAddress{Provience: "江苏"}
	zj := UserAddress{Provience: "浙江"}
	ah := UserAddress{Provience: "安徽"}
	fj := UserAddress{Provience: "福建"}
	jx := UserAddress{Provience: "江西"}
	sd := UserAddress{Provience: "山东"}
	hn := UserAddress{Provience: "河南"}
	heb := UserAddress{Provience: "湖北"}
	hun := UserAddress{Provience: "湖南"}
	gd := UserAddress{Provience: "广东"}
	gs := UserAddress{Provience: "甘肃"}
	sc := UserAddress{Provience: "四川"}
	gz := UserAddress{Provience: "贵州"}
	hunan := UserAddress{Provience: "海南"}
	yn := UserAddress{Provience: "云南"}
	qh := UserAddress{Provience: "青海"}
	shx := UserAddress{Provience: "陕西"}
	gx := UserAddress{Provience: "广西"}
	xz := UserAddress{Provience: "西藏"}
	nx := UserAddress{Provience: "宁夏"}
	xj := UserAddress{Provience: "新疆"}
	nmg := UserAddress{Provience: "内蒙古"}
	am := UserAddress{Provience: "澳门"}
	xg := UserAddress{Provience: "香港"}
	s = []UserAddress{bj, tj, sh, cq, hb, sx, tw, ll, jl, hlj, js, zj, ah, fj, jx, sd, hn, heb, hun, gd, gs, sc, gz, hunan, yn, qh, shx, gx, xz, nx, xj, nmg, am, xg}
	return
}

func StaLoanMoneyInit() (s []BusinessLoanMoney) {
	s1 := BusinessLoanMoney{MoneyAccount: "1000以下"}
	s2 := BusinessLoanMoney{MoneyAccount: "1000-2999"}
	s3 := BusinessLoanMoney{MoneyAccount: "3000-4999"}
	s4 := BusinessLoanMoney{MoneyAccount: "5000-7999"}
	s5 := BusinessLoanMoney{MoneyAccount: "8000-10000"}
	s6 := BusinessLoanMoney{MoneyAccount: "1~3万"}
	s7 := BusinessLoanMoney{MoneyAccount: "3万以上"}
	s8 := BusinessLoanMoney{MoneyAccount: "暂无数据"}
	s = []BusinessLoanMoney{s1, s2, s3, s4, s5, s6, s7, s8}
	return
}

func BusinessLoanTermCountInit() (s []BusinessLoanTermCount) {
	dkqx1 := BusinessLoanTermCount{Termcount: "7天"}
	dkqx2 := BusinessLoanTermCount{Termcount: "14天"}
	dkqx3 := BusinessLoanTermCount{Termcount: "21天"}
	dkqx4 := BusinessLoanTermCount{Termcount: "28天"}
	dkqx5 := BusinessLoanTermCount{Termcount: "1个月"}
	dkqx6 := BusinessLoanTermCount{Termcount: "3个月"}
	dkqx7 := BusinessLoanTermCount{Termcount: "6个月"}
	s = []BusinessLoanTermCount{dkqx1, dkqx2, dkqx3, dkqx4, dkqx5, dkqx6, dkqx7}
	return
}

func BusinessLoanTimesInit() (s []BusinessLoanTimes) {
	t0 := BusinessLoanTimes{Times: "0次"}
	t1 := BusinessLoanTimes{Times: "1次"}
	t2 := BusinessLoanTimes{Times: "2次"}
	t3 := BusinessLoanTimes{Times: "3次"}
	t4 := BusinessLoanTimes{Times: "4次"}
	t5 := BusinessLoanTimes{Times: "5次"}
	t6 := BusinessLoanTimes{Times: "5次以上"}
	s = []BusinessLoanTimes{t0, t1, t2, t3, t4, t5, t6}
	return
}

func AgeInit() (s []UserGender) {
	a1 := UserGender{Sex: "女"}
	a2 := UserGender{Sex: "男"}
	a3 := UserGender{Sex: "未认证"}
	s = []UserGender{a1, a2, a3}
	return
}

func ZMScoreInit() (s []StaUser) {
	score1 := StaUser{Score: "500分以下"}
	score2 := StaUser{Score: "500-549分"}
	score3 := StaUser{Score: "550-599分"}
	score4 := StaUser{Score: "600-649分"}
	score5 := StaUser{Score: "650-699分"}
	score6 := StaUser{Score: "700分以上"}
	score7 := StaUser{Score: "未认证"}
	s = []StaUser{score1, score2, score3, score4, score5, score6, score7}
	return
}

func StaAgeInit() (s []UserAge) {
	age1 := UserAge{Data: "18岁以下"}
	age2 := UserAge{Data: "18-20岁"}
	age3 := UserAge{Data: "21-22岁"}
	age4 := UserAge{Data: "23-35岁"}
	age5 := UserAge{Data: "36-40岁"}
	age6 := UserAge{Data: "41-55岁"}
	age7 := UserAge{Data: "55以上"}
	age8 := UserAge{Data: "未认证"}
	s = []UserAge{age1, age2, age3, age4, age5, age6, age7, age8}
	return
}

func StaCreditInit() (s []UserCredit) {
	credit1 := UserCredit{Credit: "信用评分低于45分"}
	credit2 := UserCredit{Credit: "信用评分45-50分"}
	credit3 := UserCredit{Credit: "信用评分51-60分"}
	credit4 := UserCredit{Credit: "信用评分61-75分"}
	credit5 := UserCredit{Credit: "信用评分高于75分"}
	credit6 := UserCredit{Credit: "未认证"}
	s = []UserCredit{credit1, credit2, credit3, credit4, credit5, credit6}
	return
}

func StaJobInit() (s []UserJob) {
	job1 := UserJob{Userjob: "上班族"}
	job2 := UserJob{Userjob: "个体户"}
	job3 := UserJob{Userjob: "无固定职业"}
	job4 := UserJob{Userjob: "企业主"}
	job5 := UserJob{Userjob: "学生"}
	job6 := UserJob{Userjob: "未认证"}
	s = []UserJob{job1, job2, job3, job4, job5, job6}
	return
}

func StaOSInit() (s []UserOS) {
	os1 := UserOS{Os: "ios"}
	os2 := UserOS{Os: "android"}
	os3 := UserOS{Os: "wx"}
	os4 := UserOS{Os: "pc"}
	os5 := UserOS{Os: "wd"}
	os6 := UserOS{Os: "wap"}
	os7 := UserOS{Os: "落地页"}
	os8 := UserOS{Os: "未激活"}
	s = []UserOS{os1, os2, os3, os4, os5, os6, os7, os8}
	return
}

func UserOperatorsInit() (s []UserOperators) {
	op1 := UserOperators{Operators: "移动"}
	op2 := UserOperators{Operators: "联通"}
	op3 := UserOperators{Operators: "电信"}
	op4 := UserOperators{Operators: "其他运营商"}
	s = []UserOperators{op1, op2, op3, op4}
	return
}

//获取今日数据明细
func GetTodayDailyDatas() (detailData DetailData, err error) {
	sql := `SELECT 
			  CURDATE() create_date,
			  (SELECT COUNT(1) FROM users WHERE create_date = CURDATE()) register_count, 
			  (SELECT COUNT(1) FROM users_auth WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2 AND is_link_man = 2 AND GREATEST(user_data_time,zm_auth_time) >= CURDATE()) identify_count, 
			  (SELECT COUNT(1) FROM business_loan WHERE (state = 'CONFIRM' OR state = 'FINISH' OR state = 'OVERDUE') AND real_time >= CURDATE()) credit_count 
			FROM DUAL ;`
	err = orm.NewOrm().Raw(sql).QueryRow(&detailData)
	return

}

//活跃今日用户明细
func GetTodayActiveData() (detailData ActiveData, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT CURDATE() create_date,COUNT(DISTINCT uid) AS active_count 
			FROM login_record 
			WHERE create_time >= CURDATE() `
	err = o.Raw(sql).QueryRow(&detailData)
	return
}

//获取今日注册用户数
func GetTodayRegisterUsers() (ru RegisterUser, err error) {
	sql := `SELECT  CURDATE() create_date,COUNT(1) count 
        	FROM users 
	        WHERE create_date = CURDATE()`
	err = orm.NewOrm().Raw(sql).QueryRow(&ru)
	return
}

//获取今日认证用户数
func GetTodayIdentifyUsers() (ru RegisterUser, err error) {
	sql := `SELECT CURDATE() create_date,COUNT(1) count 
			FROM users_auth 
			WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2 AND is_link_man=2 AND DATE(real_name_time)=CURDATE() AND DATE(link_man_time)=CURDATE() AND DATE(user_data_time)=CURDATE() AND DATE(zm_auth_time)=CURDATE() `
	err = orm.NewOrm().Raw(sql).QueryRow(&ru)
	return
}

//获取今日申请贷款次数
func GetTodayLoanUsers() (ru RegisterUser, err error) {
	sql := `SELECT CURDATE() create_date,COUNT(1) count 
			FROM business_loan 
			WHERE create_time >= CURDATE() `
	err = orm.NewOrm().Raw(sql).QueryRow(&ru)
	return
}

//获取今日放款用户数
func GetTodayCreditUsers() (ru RegisterUser, err error) {
	sql := `SELECT CURDATE() create_date,COUNT(1) count 
			FROM business_loan 
			WHERE (state='CONFIRM' OR state='FINISH' OR state='OVERDUE') AND real_time >= CURDATE()`
	err = orm.NewOrm().Raw(sql).QueryRow(&ru)
	return
}

//获取今日活跃用户数量
func GetTodayActiveUsers() (ru RegisterUser, err error) {
	sql := `SELECT CURDATE() create_date,COUNT(1) count FROM users WHERE login_time >= CURDATE()`
	o := orm.NewOrm()
	err = o.Raw(sql).QueryRow(&ru)
	return
}

//按天查看今日以前数据明细
func GetBeforeDailyDatas() (detailData []DetailData, err error) {
	o := orm.NewOrm()
	sql := `SELECT register.create_date,register_count,identify_count,credit_count FROM
			(SELECT create_date ,COUNT(1) register_count FROM users WHERE create_date < CURDATE() GROUP BY create_date) register
			LEFT JOIN
			(SELECT DATE(GREATEST(user_data_time,zm_auth_time)) dayi,COUNT(1) identify_count FROM users_auth 
			WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2 AND is_link_man=2 AND GREATEST(user_data_time,zm_auth_time) < CURDATE() GROUP BY dayi) identify
			ON register.create_date=identify.dayi
			LEFT JOIN
			(SELECT DATE(real_time) dayc , COUNT(1) credit_count
			FROM business_loan 
			WHERE (state = 'CONFIRM' OR state = 'FINISH' OR state = 'OVERDUE') AND real_time < CURDATE()
			GROUP BY dayc
			) credit
			ON register.create_date=credit.dayc
			ORDER BY register.create_date DESC `
	_, err = o.Raw(sql).QueryRows(&detailData)
	if err == nil {
		if data, err2 := json.Marshal(detailData); err2 == nil && utils.Re == nil {
			h := 24 - time.Now().Hour()
			utils.Rc.Put(utils.CacheKeyDailyDatas, data, time.Duration(h)*time.Hour)
		}
	}
	return
}

//活跃用户明细
func GetBeforeActiveData() (detailData []ActiveData, err error) {
	o := orm.NewOrm()
	sql := `SELECT 
				DATE(login_time) create_date,
				COUNT(1) active_count 
			FROM
				users 
			WHERE login_time < CURDATE() 
			GROUP BY DATE(login_time) 
			ORDER BY DATE(login_time) DESC `
	_, err = o.Raw(sql).QueryRows(&detailData)
	if err == nil {
		if data, err2 := json.Marshal(detailData); err2 == nil && utils.Re == nil {
			h := 24 - time.Now().Hour()
			utils.Rc.Put(utils.CacheKeyActiveData, data, time.Duration(h)*time.Hour)
		}
	}
	return
}

//获取今日以前注册用户数
func GetBeforeRegisterUsers() (rus []RegisterUser, err error) {
	sql := `SELECT create_date,COUNT(1) count 
			FROM users 
			WHERE create_date < CURDATE()
	        GROUP BY  create_date
        	ORDER BY  create_date `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			h := 24 - time.Now().Hour()
			utils.Rc.Put(utils.CacheKeyRegisterUsers, data, time.Duration(h)*time.Hour)
		}
	}
	return
}

//获取认证用户数
func GetBeforeIdentifyUsers() (rus []RegisterUser, err error) {
	sql := `SELECT 
				DATE(user_data_time) create_date,
				COUNT(1) count
			FROM
				users_auth 
			WHERE is_real_name = 2 
				AND is_user_data = 2 
				AND is_zm_auth = 2 
				AND is_link_man = 2 
				AND DATE(real_name_time)= DATE(link_man_time)
				AND DATE(real_name_time)= DATE(user_data_time)
				AND DATE(real_name_time)= DATE(zm_auth_time) 
				AND DATE(real_name_time)< CURDATE()
			GROUP BY create_date 
			ORDER BY create_date`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			h := 24 - time.Now().Hour()
			utils.Rc.Put(utils.CacheKeyIdentifyUsers, data, time.Duration(h)*time.Hour)
		}
	}
	return
}

//获取今日以前申请贷款次数
func GetBeforeLoanUsers() (rus []RegisterUser, err error) {
	sql := ` SELECT DATE(create_time) create_date,COUNT(1) count FROM business_loan 
             WHERE create_time < CURDATE()
             GROUP BY  DATE(create_time)
			 ORDER BY  DATE(create_time) `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			h := 24 - time.Now().Hour()
			utils.Rc.Put(utils.CacheKeyLoanUsers, data, time.Duration(h)*time.Hour)
		}
	}
	return
}

//获取今日以前放款用户数
func GetBeforeCreditUsers() (rus []RegisterUser, err error) {
	sql := `SELECT DATE(real_time) create_date,COUNT(1) count FROM business_loan
			WHERE (state='CONFIRM' OR state='FINISH' OR state='OVERDUE') AND real_time < CURDATE()
			GROUP BY  DATE(real_time)
			ORDER BY  DATE(real_time)`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			h := 24 - time.Now().Hour()
			utils.Rc.Put(utils.CacheKeyCreditUsers, data, time.Duration(h)*time.Hour)
		}
	}
	return
}

//获取今日以前活跃用户数
func GetBeforeActiveUsers() (rus []RegisterUser, err error) {
	sql := `SELECT DATE(create_time) create_date, COUNT(DISTINCT uid) count
			FROM login_record WHERE create_time < CURDATE() 
            GROUP BY  DATE(create_time)
			ORDER BY  DATE(create_time) `
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err = o.Raw(sql).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			h := 24 - time.Now().Hour()
			utils.Rc.Put(utils.CacheKeyActiveUsers, data, time.Duration(h)*time.Hour)
		}
	}
	return
}

//获取累计认证用户数量
func GetIdentifyUserCount() (count int, err error) {
	sql := `SELECT COUNT(1)  count
			FROM users_auth
			WHERE is_real_name = 2  AND is_user_data = 2 AND is_zm_auth = 2 AND is_link_man = 2 `
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取累计申请借款次数
func GetLoanUserCount() (count int, err error) {
	sql := `SELECT COUNT(1) AS count FROM business_loan  `
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取累计放款用户数量
func GetCreditUserCount() (count int, err error) {
	sql := `SELECT COUNT(1)  count
			FROM business_loan 
			WHERE (state = 'CONFIRM' OR state = 'FINISH' OR state = 'OVERDUE') `
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取所有外链渠道名称
func GetOutPutSourceName(condition, params string) (name []string, err error) {
	sql := `SELECT out_put_source FROM users WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY out_put_source HAVING out_put_source !=""`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&name)
	return
}

//获取第一个外链渠道名称
func GetFristOutPutSourceName() (name string, err error) {
	sql := `SELECT out_put_source FROM users GROUP BY out_put_source HAVING out_put_source !="" LIMIT 1`
	err = orm.NewOrm().Raw(sql).QueryRow(&name)
	return
}

//根据名称获取外链渠道
func GetOutPutSourceNameByName(sourceName string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users WHERE out_put_source =? `
	err = orm.NewOrm().Raw(sql, sourceName).QueryRow(&count)
	return
}
