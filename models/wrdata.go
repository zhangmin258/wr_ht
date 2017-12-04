package models

import (
	"encoding/json"
	"time"
	"wr_v1/utils"

	"github.com/astaxie/beego/orm"
)

//微融信息
type DailyData struct {
	Date             time.Time //日期
	RegisterCount    int       //注册用户数量
	LoginCount       int       //激活用户数量
	OcrCount         int       //完成orc数量
	IdentifyCount    int       //完成认证用户数量(基础认证)
	ApplynowCount    int       //点击立即申请用户数量
	ProRegisterCount int       //第三方导流量
	UserPerCount     float64   //人均注册平台
	UserPerProfit    float64   //人均创造收益
	TotalProfit      float64   //总收益
}

//注册用户统计
type WrRegisterUser struct {
	CreateDate string
	Count      int
}

type WrDataCount struct {
	CreateDate time.Time
	Count      int
}

type WrDataFCount struct {
	CreateDate time.Time
	Count      float64
}

type MonGoQuery struct {
	Uid int
}

//导出到Excel
type ExportToExcel struct {
	CreateDate    string //日期
	RegisterCount int    //注册用户
	FirstCount    int    //激活用户
	OCRCount      int    //完成OCR用户
	IdentifyCount int    //认证完成用户
	ApplynowCount int    //点击立即申请用户数量
}

type WeiRongDataAll struct {
	RegisterCount int //注册用户
	FirstCount    int //激活用户
	OcrCount      int //完成OCR用户
	ApplynowCount int //点击立即申请用户数量
	LoanUser      int //放款用户数量
}
type DailyDataSort []DailyData

func (I DailyDataSort) Len() int {
	return len(I)
}
func (I DailyDataSort) Less(i, j int) bool {
	return I[i].Date.Unix() > I[j].Date.Unix()
}
func (I DailyDataSort) Swap(i, j int) {
	I[i], I[j] = I[j], I[i]
}

//截止今天凌晨累计微融用户统计
func GetWrRegisterAllCount(source string) (weiRongDataAll *WeiRongDataAll, err error) {
	sql := `SELECT a.register_count,b.first_count,c.ocr_count,d.applynow_count,e.loan_user,f.identify_count FROM `
	sql += `(SELECT COUNT(1) AS register_count FROM users u WHERE create_date<curdate() `
	if source != "" {
		sql += source
	}
	sql += `) a,(SELECT COUNT(1) AS first_count FROM users u WHERE active_time IS NOT NULL AND create_date<curdate() `
	if source != "" {
		sql += source
	}
	sql += ` ) b,(SELECT COUNT(uid) AS ocr_count FROM users u INNER JOIN users_auth ua ON u.id=ua.uid WHERE is_real_name=2 AND real_name_time<curdate() `
	if source != "" {
		sql += source
	}
	sql += ` ) c,(SELECT  COUNT(product_id) AS applynow_count FROM product_down_record AS p
		INNER JOIN users AS u
		ON p.user_id=u.id
 		WHERE DATE_FORMAT(p.create_time,"%Y-%m-%d") = u.create_date AND p.create_time<curdate()`
	if source != "" {
		sql += source
	}
	sql += ` ) d,(SELECT COUNT(1) AS loan_user FROM business_loan_h5 bl inner join users u on bl.uid=u.id WHERE  bl.create_time<curdate()`
	if source != "" {
		sql += source
	}
	sql += ` ) e,(SELECT COUNT(1) AS identify_count FROM users_auth ua
	INNER JOIN users u
	ON ua.uid=u.id
	WHERE ua.is_real_name = 2 AND ua.is_user_data = 2 AND ua.is_zm_auth = 2  AND ua.is_link_man = 2  AND GREATEST(ua.user_data_time,ua.zm_auth_time)<curdate() `
	if source != "" {
		sql += source
	}
	sql += ` )f`
	err = orm.NewOrm().Raw(sql).QueryRow(&weiRongDataAll)
	if data, err2 := json.Marshal(weiRongDataAll); err == nil && err2 == nil && utils.Re == nil {
		cache := ""
		if source == "  AND u.source !='' " {
			cache = utils.WEIRONGCOUNTSOURCE
		} else {
			cache = utils.WEIRONGCOUNTOUTPUTSOURCE
		}
		h := 24 - time.Now().Hour()
		utils.Rc.Put(cache, data, time.Duration(h)*time.Hour)
	}
	return
}

//历史数据明细
func GetChannelHistoryData(params []string) (dailyData []DailyData, err error) {
	sql := `SELECT data_date AS date,
			SUM(register_count) AS register_count,
			SUM(active_count) AS login_count,
			SUM(ocr_count) AS ocr_count,
			SUM(auth_count) AS identify_count,
			SUM(apply_count) AS applynow_count,
			SUM(pro_register_count) AS pro_register_count,
			SUM(total_profit) AS total_profit 
		FROM channel_data WHERE out_put_source!=''`
	if len(params) > 0 {
		sql += ` AND out_put_source=? `
	}
	sql += ` GROUP BY data_date ORDER BY data_date DESC`
	o := orm.NewOrm()
	o.Using("wr_backup")
	_, err = o.Raw(sql, params).QueryRows(&dailyData)
	return
}

//历史累计数据
func GetTotalHistoryData(outPutSource string, params []string) (weiRongDataAll WeiRongDataAll, err error) {
	sql := `SELECT SUM(register_count) AS register_count,SUM(active_count) AS first_count,SUM(ocr_count) AS ocr_count,SUM(apply_count) AS applynow_count FROM channel_data WHERE data_date<CURDATE() `
	if outPutSource != "" {
		sql += ` AND out_put_source=?`
	}
	o := orm.NewOrm()
	o.Using("wr_backup")
	err = o.Raw(sql, params).QueryRow(&weiRongDataAll)
	return
}

// 今天微融注册用户
func GetWrRegisterCountToday(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u WHERE out_put_source!='' AND create_date>=curdate() `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 今天激活（注册用户有过登录行为，通过users表里的登录时间是否为空来判断）
func GetWrFirstUserCountToday(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u WHERE active_time IS NOT NULL AND active_time>=CURDATE() AND out_put_source!='' `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 今天点击立即申请用户
func GetApplyNowUserCountToday(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users_first_loan_time t INNER JOIN users u ON t.uid=u.id  WHERE t.first_loan_date=CURDATE() AND u.out_put_source!=''`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//今天完成orc实名认证用户
func GetWrORCCountToday(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 AND a.real_name_time>=CURDATE() AND out_put_source!='' `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//今天完成认证用户
func GetWrIdentifyCountToday(condition string, params []string) (count int, err error) {
	sql := `SELECT
				COUNT(1) AS count
			FROM
				users u
			INNER JOIN users_auth ua ON u.id = ua.uid
			WHERE
				is_real_name = 2
			AND is_zm_auth = 2
			AND is_link_man = 2
			AND is_user_data = 2
			AND DATE(ua.real_name_time)= DATE(ua.link_man_time) 
			AND DATE(ua.real_name_time)= DATE(ua.user_data_time) 
			AND DATE(ua.real_name_time)= DATE(ua.zm_auth_time) 
			AND u.out_put_source!=''
			AND ua.user_data_time>=CURDATE()  `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//今日导流量
func GetProRegitsterCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM product_register_user p LEFT JOIN users u ON p.uid = u.id WHERE p.create_time >= CURDATE() AND u.out_put_source!=''`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//今日总收益
func GetTodayProfit(condition string, params []string) (profit float64, err error) {
	sql := `SELECT SUM(pc.cpa_price) AS count
			FROM product_register_user AS pru
			LEFT JOIN product_cleaning AS pc ON pru.product_id = pc.product_id
			LEFT JOIN users AS u ON u.id=pru.uid
			WHERE pru.create_time >= CURDATE() 
			AND u.out_put_source!='' `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&profit)
	return
}

//今天借款的用户
func GetTodayLoanUser(condition string) (count int, err error) {
	sql := `SELECT COUNT(1) AS loan_user FROM business_loan_h5 bl inner join users u on bl.uid=u.id WHERE  bl.create_time>=curdate() `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取每天的注册，登录，认证的人数  添加第三方导流量	人均注册平台	人均创造收益 总收益
func GetDailyDataCache(condition string) (dailyDataList []DailyData, err error) {
	sql := `select day1 as date,register_count,login_count,ocr_count,identify_count,applynow_count,pro_register_count,total_profit from
	(SELECT create_date AS day1,COUNT(1) AS register_count FROM users u WHERE create_date<curdate() `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day1 ) a
	left join
	(SELECT DATE_FORMAT(active_time,"%Y-%m-%d") AS day2,COUNT(1) AS login_count FROM users u WHERE active_time<curdate() `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day2 HAVING day2!="" ) b
	on a.day1=b.day2
	left join
	(SELECT DATE_FORMAT(real_name_time,"%Y-%m-%d") AS day3,COUNT(1) AS  ocr_count FROM users u INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 AND real_name_time<curdate() `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day3) c
	on a.day1=c.day3
	left join
	(SELECT DATE(GREATEST(user_data_time,zm_auth_time)) AS day4,COUNT(1) AS identify_count FROM users_auth ua INNER JOIN users u  ON ua.uid = u.id WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2  AND is_link_man = 2 AND GREATEST(user_data_time,zm_auth_time)<curdate() `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day4) d
	on a.day1=d.day4
	LEFT JOIN
	(SELECT  DATE_FORMAT(p.create_time,"%Y-%m-%d") AS day5,COUNT(product_id) AS applynow_count FROM product_down_record AS p
		INNER JOIN users AS u
		ON p.user_id=u.id
 		WHERE DATE_FORMAT(p.create_time,"%Y-%m-%d") = u.create_date AND p.create_time<curdate() `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day5) f
	ON a.day1=f.day5 
	LEFT JOIN  
	(SELECT DATE(pru.create_time) AS day6,SUM(pc.cpa_price) AS total_profit
	FROM product_register_user pru LEFT JOIN users u ON pru.uid = u.id LEFT JOIN product_cleaning pc ON pru.product_id = pc.product_id
	WHERE pru.create_time < CURDATE() `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day6) g
	ON a.day1=g.day6 
	LEFT JOIN  
	(SELECT DATE(pru.create_time) AS day7,COUNT(1) AS pro_register_count 
	FROM product_register_user pru
	LEFT JOIN users u ON pru.uid = u.id 
	WHERE pru.create_time < CURDATE() `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day7) h
	ON a.day1=h.day7 ORDER BY date DESC`
	// sql += ` GROUP BY day5) f
	// ON a.day1=f.day5  ORDER BY date DESC`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&dailyDataList)
	for i, _ := range dailyDataList {
		if dailyDataList[i].RegisterCount != 0 {
			dailyDataList[i].UserPerCount = float64(dailyDataList[i].ProRegisterCount) / float64(dailyDataList[i].RegisterCount)
			dailyDataList[i].UserPerCount = utils.SubFloatToFloat(dailyDataList[i].UserPerCount, 2)
			dailyDataList[i].UserPerProfit = dailyDataList[i].TotalProfit / float64(dailyDataList[i].RegisterCount)
			dailyDataList[i].UserPerProfit = utils.SubFloatToFloat(dailyDataList[i].UserPerProfit, 2)
		}
	}

	if data, err2 := json.Marshal(dailyDataList); err == nil && err2 == nil && utils.Re == nil {
		cache := ""
		if condition == "  AND u.source !='' " {
			cache = utils.WEIRONGEVERYDAYDATASOURCE
		} else {
			cache = utils.WEIRONGEVERYDAYDATAOUTPUTSOURCE
		}
		h := 24 - time.Now().Hour()
		utils.Rc.Put(cache, data, time.Duration(h)*time.Hour)
	}
	return
}

// 累计微融注册用户
func GetWrRegisterCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u WHERE u.out_put_source!=''`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 累计激活（注册用户有过登录行为，通过users表里的登录时间是否为空来判断）
func GetWrFirstUserCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u WHERE active_time IS NOT NULL AND u.out_put_source!='' `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 累计点击立即申请用户
func GetApplyNowUserCount(condition string, params []string) (count int, err error) {
	sql := ` SELECT
				COUNT(1) AS count
			FROM
				users_first_loan_time t
			INNER JOIN users u ON t.uid = u.id
			WHERE
				u.out_put_source!=''`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//完成orc实名认证用户
func GetWrORCCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 AND u.out_put_source!=''`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//累计放款用户
func GetWrLoanCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM business_loan_h5 bl INNER JOIN users u ON bl.uid = u.id WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//完成认证用户
func GetWrIdentifyCount(condition string, params []string) (count int, err error) {
	sql := `SELECT
				COUNT(1) AS count
			FROM
				users u
			INNER JOIN users_auth ua ON u.id = ua.uid
			WHERE
				is_real_name = 2
			AND is_zm_auth = 2
			AND is_link_man = 2
			AND is_user_data = 2
			AND DATE(ua.real_name_time)= DATE(ua.link_man_time) 
			AND DATE(ua.real_name_time)= DATE(ua.user_data_time) 
			AND DATE(ua.real_name_time)= DATE(ua.zm_auth_time) 
			AND u.out_put_source!=''`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//模糊查询获取所有渠道列表(按导流量排序)
func GetOutPutSourceListByName(condition string, params []string) (sourceList []string, err error) {
	sql := `SELECT a.out_put_source FROM 
            (SELECT u.out_put_source FROM users u GROUP BY u.out_put_source HAVING  u.out_put_source !="") a 
            LEFT JOIN 
            (SELECT u.out_put_source,COUNT(1) c FROM users u LEFT JOIN product_register_user p ON p.uid=u.id WHERE DATE(p.create_time)>=CURDATE() GROUP BY u.out_put_source HAVING u.out_put_source !="") b
            ON a.out_put_source = b.out_put_source 
            WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += `ORDER BY b.c DESC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&sourceList)
	return
}

//按名字判断是否有该渠道
func GetOutPutSourceByName(source string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users WHERE out_put_source = ? `
	err = orm.NewOrm().Raw(sql, source).QueryRow(&count)
	return
}

// 获取所有的渠道列表
func GetSourceList() (sourceList []string, err error) {
	sql := `SELECT u.source FROM users u GROUP BY source HAVING	u.source !=""`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&sourceList)
	return
}

//获取所有每天的注册，登录，认证的人数
func GetDailyDataAll(condition string, params []string) (dailyDataList []DailyData, err error) {
	sql := `select day1 as date,register_count,login_count,ocr_count,identify_count,applynow_count from
	(SELECT DATE_FORMAT(create_time,"%Y-%m-%d") AS day1,COUNT(1) AS register_count FROM users u WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day1 ) a
	left join
	(SELECT DATE_FORMAT(active_time,"%Y-%m-%d") AS day2,COUNT(1) AS login_count FROM users u WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day2 HAVING day2!="" ) b
	on a.day1=b.day2
	left join
	(SELECT DATE_FORMAT(real_name_time,"%Y-%m-%d") AS day3,COUNT(1) AS  ocr_count FROM users u INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day3) c
	on a.day1=c.day3
	left join
	(SELECT DATE_FORMAT(user_data_time,"%Y-%m-%d") AS day4,COUNT(1) AS identify_count FROM users_auth ua LEFT JOIN users u  ON ua.uid = u.id WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2  AND is_link_man = 2 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day4) d
	on a.day1=d.day4
	LEFT JOIN
    (SELECT  DATE_FORMAT(p.create_time,"%Y-%m-%d") AS day5,COUNT(product_id) AS applynow_count FROM product_down_record AS p
		INNER JOIN users AS u
		ON p.user_id=u.id
 		WHERE DATE_FORMAT(p.create_time,"%Y-%m-%d") = u.create_date `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day5) f  ON a.day1=f.day5 ORDER BY date DESC `
	_, err = orm.NewOrm().Raw(sql, params, params, params, params).QueryRows(&dailyDataList)
	return
}

//根据渠道获取每天的注册，登录，认证。点击立即申请的人数
func GetDailyDataBySource(types string, source string, begin, size int) (dailyDataList []DailyData, err error) {
	sql := `SELECT day1 AS date,register_count,login_count,ocr_count,identify_count , applynow_count FROM
	(SELECT create_date AS day1,COUNT(1) AS register_count FROM users WHERE ` + types + `=? GROUP BY day1 ) a
	LEFT JOIN
	(SELECT DATE_FORMAT(active_time,"%Y-%m-%d") AS day2,COUNT(1) AS login_count FROM users WHERE ` + types + `=? GROUP BY day2 HAVING day2!="" ) b
	ON a.day1=b.day2
	LEFT JOIN
	(SELECT DATE_FORMAT(real_name_time,"%Y-%m-%d") AS day3,COUNT(1) AS  ocr_count FROM users u  INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 AND u.` + types + `=? GROUP BY day3) c
	ON a.day1=c.day3
	LEFT JOIN
	(SELECT DATE_FORMAT(user_data_time,"%Y-%m-%d") AS day4,COUNT(1) AS identify_count FROM users_auth a LEFT JOIN users u ON u.id=a.uid  WHERE u.` + types + `=? AND  is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2  AND is_link_man = 2 GROUP BY day4) d
	ON a.day1=d.day4
	LEFT JOIN
	(SELECT  DATE_FORMAT(p.create_time,"%Y-%m-%d") AS day5,COUNT(product_id) AS applynow_count FROM product_down_record AS p
		INNER JOIN users AS u
		ON p.user_id=u.id
 		WHERE DATE_FORMAT(p.create_time,"%Y-%m-%d") = u.create_date AND u.` +
		types + `=? GROUP BY day5)f
	 ON a.day1=f.day5
	ORDER BY date DESC limit ?, ?`
	_, err = orm.NewOrm().Raw(sql, source, source, source, source, source, begin, size).QueryRows(&dailyDataList)
	return
}

//根据渠道获取所有每天的注册，登录，认证。点击立即申请的人数
func GetDailyDataBySourceAll(types string, source string) (dailyDataList []DailyData, err error) {
	sql := `SELECT day1 AS date,register_count,login_count,ocr_count,identify_count , applynow_count FROM
	(SELECT create_date AS day1,COUNT(1) AS register_count FROM users WHERE ` + types + `=? GROUP BY day1 ) a
	LEFT JOIN
	(SELECT DATE_FORMAT(active_time,"%Y-%m-%d") AS day2,COUNT(1) AS login_count FROM users WHERE ` + types + `=? GROUP BY day2 HAVING day2!="" ) b
	ON a.day1=b.day2
	LEFT JOIN
	(SELECT DATE_FORMAT(real_name_time,"%Y-%m-%d") AS day3,COUNT(1) AS  ocr_count FROM users u  INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 AND u.` + types + `=? GROUP BY day3) c
	ON a.day1=c.day3
	LEFT JOIN
	(SELECT DATE_FORMAT(user_data_time,"%Y-%m-%d") AS day4,COUNT(1) AS identify_count FROM users_auth a LEFT JOIN users u ON u.id=a.uid  WHERE u.` + types + `=? AND  is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2  AND is_link_man = 2 GROUP BY day4) d
	ON a.day1=d.day4
	LEFT JOIN
	(SELECT  DATE_FORMAT(p.create_time,"%Y-%m-%d") AS day5,COUNT(product_id) AS applynow_count FROM product_down_record AS p
		INNER JOIN users AS u
		ON p.user_id=u.id
 		WHERE DATE_FORMAT(p.create_time,"%Y-%m-%d") = u.create_date AND u.` +
		types + `=? GROUP BY day5)f
	 ON a.day1=f.day5
	ORDER BY date DESC `
	_, err = orm.NewOrm().Raw(sql, source, source, source, source, source).QueryRows(&dailyDataList)
	return
}

//获取数据总条数
func GetDailyDataCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM (SELECT DATE_FORMAT(u.create_time,"%Y-%m-%d") AS day1,COUNT(1) AS register_count FROM users u WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day1) AS a`
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//获取每日注册用户数量
func GetWrRegisterUsersByCondition(condition string, params []string) (rus []WrRegisterUser, err error) {
	sql := `SELECT  create_date,COUNT(1) AS count  FROM users u WHERE out_put_source!='' `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  create_date
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&rus)
	return
}

//获取每日激活用户数量
func GetWrFirstUsersByCondition(condition string, params []string) (fus []WrRegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(active_time,'%Y-%m-%d') AS create_date, COUNT(1) AS count FROM users u WHERE u.out_put_source!='' `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY DATE_FORMAT(active_time,'%Y-%m-%d')`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&fus)
	return
}

//获取每日立即点击申请用户数量
func GetApplyNowUsersByCondition(condition string, params []string) (fus []WrRegisterUser, err error) {
	sql := `SELECT
				COUNT(1) AS count,
				t.first_loan_date AS create_date
			FROM
				users_first_loan_time t
			INNER JOIN users u ON t.uid = u.id
			WHERE
				u.out_put_source !='' `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY t.first_loan_date `
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&fus)
	return
}

//获取每日ocr用户数量
func GetWrOcrUsersByCondition(condition string, params []string) (ous []WrRegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(real_name_time,"%Y-%m-%d")create_date ,COUNT(1) AS  count FROM users u  INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 AND  u.out_put_source!=''`
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE_FORMAT(real_name_time,'%Y-%m-%d')
			 ORDER BY  create_date ASC`

	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&ous)
	return
}

//获取每日完成认证用户数量
func GetWrIdentifyUsersByCondition(condition string, params []string) (ius []WrRegisterUser, err error) {
	sql := `SELECT DATE(GREATEST(user_data_time,zm_auth_time)) AS create_date ,COUNT(1) AS count FROM users_auth a INNER JOIN users u ON u.id=a.uid  WHERE   is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2  AND is_link_man = 2   AND DATE(a.real_name_time)= DATE(a.link_man_time) AND DATE(a.real_name_time)= DATE(a.user_data_time) AND DATE(a.real_name_time)= DATE(a.zm_auth_time) AND u.out_put_source!=''`
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE(GREATEST(user_data_time,zm_auth_time))`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&ius)
	return
}

//获取每天的注册，登录，认证的人数
func GetDailyData(condition string, params []string, begin, size int) (dailyDataList []DailyData, err error) {
	sql := `select day1 as date,register_count,login_count,ocr_count,identify_count,applynow_count from
	(SELECT create_date AS day1,COUNT(1) AS register_count FROM users u WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day1 ) a
	left join
	(SELECT DATE_FORMAT(active_time,"%Y-%m-%d") AS day2,COUNT(1) AS login_count FROM users u WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day2 HAVING day2!="" ) b
	on a.day1=b.day2
	left join
	(SELECT DATE_FORMAT(real_name_time,"%Y-%m-%d") AS day3,COUNT(1) AS  ocr_count FROM users u INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day3) c
	on a.day1=c.day3
	left join
	(SELECT DATE_FORMAT(user_data_time,"%Y-%m-%d") AS day4,COUNT(1) AS identify_count FROM users_auth ua LEFT JOIN users u  ON ua.uid = u.id WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2  AND is_link_man = 2 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day4) d
	on a.day1=d.day4 
	LEFT JOIN
    (SELECT  DATE_FORMAT(p.create_time,"%Y-%m-%d") AS day5,COUNT(product_id) AS applynow_count FROM product_down_record AS p
		INNER JOIN users AS u
		ON p.user_id=u.id
 		WHERE DATE_FORMAT(p.create_time,"%Y-%m-%d") = u.create_date  `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day5) f  ON a.day1=f.day5 ORDER BY date DESC limit ?, ?`
	_, err = orm.NewOrm().Raw(sql, params, params, params, params, params, begin, size).QueryRows(&dailyDataList)
	return
}

//按天获取注册用户
func GetWrRegisterEveryDay(source string, params []string) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT create_date, COUNT(1) AS count
		FROM users u
		WHERE create_date IS NOT NULL `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY create_date `
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&wrDataCount)
	return
}

//按天获取激活用户
func GetWrActiveEveryDay(source string, params []string) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE_FORMAT(active_time,"%Y-%m-%d") AS create_date, COUNT(1) AS count
			FROM users u
			WHERE active_time IS NOT NULL  `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE_FORMAT(active_time,"%Y-%m-%d")`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&wrDataCount)
	return
}

//按天获取orc用户
func GetWrOcrEveryDay(source string, params []string) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE_FORMAT(real_name_time,"%Y-%m-%d") AS create_date, COUNT(1) AS count
			FROM users u
			INNER JOIN users_auth a ON u.id=a.uid
			WHERE is_real_name=2 AND real_name_time IS NOT NULL `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE_FORMAT(real_name_time,"%Y-%m-%d")`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&wrDataCount)
	return
}

//按天获取完成认证用户
func GetWrIdentifyEveryDay(source string, params []string) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE(GREATEST(user_data_time,zm_auth_time)) AS create_date, COUNT(1) AS count
			FROM users_auth ua
			INNER JOIN users u ON ua.uid = u.id
			WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2 AND is_link_man = 2 AND user_data_time IS NOT NULL AND zm_auth_time IS NOT NULL `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE(GREATEST(user_data_time,zm_auth_time))`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&wrDataCount)
	return
}

//按天获取点击立即申请用户用户
func GetWrApplynowEveryDay(source string, params []string) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE_FORMAT(p.create_time,"%Y-%m-%d") AS create_date, COUNT(product_id) AS count
			FROM product_down_record AS p
			INNER JOIN users AS u ON p.user_id=u.id
			WHERE DATE_FORMAT(p.create_time,"%Y-%m-%d") = u.create_date AND p.create_time IS NOT NULL `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE_FORMAT(p.create_time,"%Y-%m-%d")`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&wrDataCount)
	return
}

//按天获取导流量
func GetWrProRegisterEveryDay(source string, params []string) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE(p.create_time) AS create_date,COUNT(1) AS count 
			FROM product_register_user p 
			LEFT JOIN users u 
			ON p.uid = u.id 
			WHERE 1 = 1 `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE(p.create_time)`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&wrDataCount)
	return
}

//按天获取总收益
func GetWrProfitEveryDay(source string, params []string) (wrDataFCount []WrDataFCount, err error) {
	sql := `SELECT SUM(pc.cpa_price) as count ,DATE(pru.create_time) create_date 
			FROM product_register_user pru 
			LEFT JOIN users u ON pru.uid = u.id 
			LEFT JOIN product_cleaning pc ON pru.product_id = pc.product_id 
			WHERE 1 = 1 `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE(pru.create_time)`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&wrDataFCount)
	return
}
