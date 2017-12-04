package models

import (
	"encoding/json"
	"time"
	"wr_v1/utils"

	"github.com/astaxie/beego/orm"
)

// 今天微融注册用户
func GetWrSCRegisterCountToday(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u WHERE 1=1 AND create_date>=curdate() AND u.out_put_source=""`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//今天完成orc实名认证用户
func GetWrSCORCCountToday(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 AND a.real_name_time>=curdate() AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//今天完成认证用户
func GetWrSCIdentifyCountToday(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT
				COUNT(1)
			FROM
				users_auth ua
			INNER JOIN users u ON ua.uid = u.id
			WHERE
				ua.is_real_name = 2
			AND ua.is_user_data = 2
			AND ua.is_zm_auth = 2
			AND ua.is_link_man = 2
			AND ua.user_data_time IS NOT NULL
			AND ua.zm_auth_time IS NOT NULL
			AND DATE(ua.real_name_time) = DATE(ua.link_man_time)
			AND DATE(ua.real_name_time) = DATE(ua.user_data_time)
			AND DATE(ua.real_name_time) = DATE(ua.zm_auth_time)
			AND u.register_source IS NOT NULL
			AND GREATEST(
				ua.user_data_time,
				ua.zm_auth_time
			) >= curdate()
			AND u.out_put_source = "" `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//今天借款的用户
func GetTodaySCLoanUser(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT COUNT(1) AS loan_user FROM business_loan_h5 bl inner join users u on bl.uid=u.id WHERE  bl.create_time>=curdate() AND u.out_put_source=""  `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

// 今天点击立即申请用户
func GetSCApplyNowUserCountToday(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT
				COUNT(1) AS count
			FROM
				users_first_loan_time t
			INNER JOIN users u ON t.uid = u.id
			WHERE
				u.register_source IS NOT NULL
			AND u.out_put_source = ""
			AND t.first_loan_date >= CURDATE()`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//今天第三方导流量
func GetTodayProRegesiterCount(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM product_register_user p INNER JOIN users u ON p.uid=u.id WHERE p.create_time>=CURDATE() AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//按天获取注册用户
func GetWrSCRegisterEveryDay(source string, params []string, pkgParams []int) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT create_date, COUNT(1) AS count
		FROM users u
		WHERE create_date IS NOT NULL AND create_date<CURDATE()`
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY create_date `
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&wrDataCount)
	return
}

//按天获取orc用户
func GetWrSCOcrEveryDay(source string, params []string, pkgParams []int) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE_FORMAT(real_name_time,"%Y-%m-%d") AS create_date, COUNT(1) AS count
			FROM users u
			INNER JOIN users_auth a ON u.id=a.uid
			WHERE is_real_name=2 AND real_name_time IS NOT NULL AND real_name_time<CURDATE()`
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE_FORMAT(real_name_time,"%Y-%m-%d")`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&wrDataCount)
	return
}

//按天获取完成认证用户
func GetWrSCIdentifyEveryDay(source string, params []string, pkgParams []int) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE(GREATEST(user_data_time,zm_auth_time)) AS create_date, COUNT(1) AS count
			FROM users_auth ua
			INNER JOIN users u ON ua.uid = u.id
			WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2 AND is_link_man = 2 AND user_data_time IS NOT NULL AND zm_auth_time IS NOT NULL AND GREATEST(ua.user_data_time,ua.zm_auth_time)<CURDATE() `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE(GREATEST(user_data_time,zm_auth_time))`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&wrDataCount)
	return
}

//按天获取点击立即申请用户用户
func GetWrSCApplynowEveryDay(source string, params []string, pkgParams []int) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE_FORMAT(p.create_time,"%Y-%m-%d") AS create_date, COUNT(product_id) AS count
			FROM product_down_record AS p
			INNER JOIN users AS u ON p.user_id=u.id
			WHERE DATE_FORMAT(p.create_time,"%Y-%m-%d") = u.create_date AND p.create_time IS NOT NULL AND u.create_time<CURDATE() `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE_FORMAT(p.create_time,"%Y-%m-%d")`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&wrDataCount)
	return
}

//按天获取第三方导流量
func GetProRegisterCount(source string, params []string, pkgParams []int) (wrDataCount []WrDataCount, err error) {
	sql := `SELECT DATE_FORMAT(p.create_time,'%Y-%m-%d') AS create_date,COUNT(product_id) AS count
			FROM product_register_user AS p 
			INNER JOIN users AS u ON p.uid=u.id
			WHERE DATE_FORMAT(p.create_time ,'%Y-%m-%d')=u.create_date AND p.create_time IS NOT NULL AND p.create_time<CURDATE()`
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE_FORMAT(p.create_time,'%Y-%m-%d')`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&wrDataCount)
	return
}

//今天本包创造的收益
// func GetTodayPackageMoney(condition string, params []string, pkgParams []int) (count float64, err error) {
// 	sql := `SELECT
// 			SUM(CASE
// 		WHEN dd.joint_mode = 1 THEN
// 			dd.cpa_price
// 		WHEN dd.joint_mode = 2 THEN
// 			dd.cps_first_per
// 		ELSE
// 			dd.cpa_price + dd.cps_first_per
// 		END) AS count
// 		FROM
// 			daily_data dd
// 		INNER JOIN agent_product ap ON dd.agent_product_id = ap.id
// 		INNER JOIN product_register_user pru ON ap.pro_id = pru.product_id
// 		INNER JOIN users u ON pru.uid=u.id
// 		WHERE pru.create_time>=CURDATE()  AND u.out_put_source="" `
// 	if condition != "" {
// 		sql += condition
// 	}
// 	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
// 	return
// }

//今天本包创造的收益
func GetTodayPackageMoney(condition string, params []string, pkgParams []int) (count float64, err error) {
	sql := `SELECT SUM(pc.cpa_price) AS count
			FROM product_register_user pru
			LEFT JOIN users u ON pru.uid = u.id
			LEFT JOIN product_cleaning pc ON pru.product_id=pc.product_id
			WHERE pru.create_time>=CURDATE()  AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//按天获取本包创造的收益
func GetPackageMoney(source string, params []string, pkgParams []int) (wrDataFCount []WrDataFCount, err error) {
	sql := `SELECT
			DATE_FORMAT(pru.create_time, '%Y-%m-%d') AS create_date,
			SUM(CASE
		WHEN dd.joint_mode = 1 THEN
			dd.cpa_price
		WHEN dd.joint_mode = 2 THEN
			dd.cps_first_per
		ELSE
			dd.cpa_price + dd.cps_first_per
		END) AS count
		FROM
			daily_data dd
		INNER JOIN agent_product ap ON dd.agent_product_id = ap.id
		INNER JOIN product_register_user pru ON ap.pro_id = pru.product_id 
		INNER JOIN users u ON pru.uid=u.id 
		WHERE pru.create_time<CURDATE() `
	if source != "" {
		sql += source
	}
	sql += ` GROUP BY DATE_FORMAT(pru.create_time, '%Y-%m-%d')`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&wrDataFCount)
	return
}

//获取当天注册用户数量
func GetWrSCRegisterUsersByCondition(condition string, params []string, pkgParams []int) (rus []WrRegisterUser, err error) {
	sql := `SELECT u.create_date,COUNT(1) AS count
		FROM
			users u
		WHERE
			 u.register_source IS NOT NULL 
		AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  create_date
			 ORDER BY  create_date ASC`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&rus)
	return
}

//获取当天ocr用户数量
func GetWrSCOcrUsersByCondition(condition string, params []string, pkgParams []int) (ous []WrRegisterUser, err error) {
	sql := `SELECT DATE_FORMAT(a.real_name_time,'%Y-%m-%d') AS create_date,COUNT(1) AS count
			FROM
				users u
			INNER JOIN users_auth a ON u.id = a.uid
			WHERE
				is_real_name = 2
			AND u.register_source IS NOT NULL 
			AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY DATE_FORMAT(a.real_name_time,'%Y-%m-%d')
			 ORDER BY  DATE_FORMAT(a.real_name_time,'%Y-%m-%d') ASC`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&ous)
	return
}

//获取当天完成认证用户数量
func GetWrSCIdentifyUsersByCondition(condition string, params []string, pkgParams []int) (ius []WrRegisterUser, err error) {
	sql := `SELECT DATE(GREATEST(user_data_time,zm_auth_time)) AS create_date, COUNT(1) AS count
			FROM users_auth ua
			INNER JOIN users u ON ua.uid = u.id
			WHERE ua.is_real_name = 2 
			AND ua.is_user_data = 2 
			AND ua.is_zm_auth = 2 
			AND ua.is_link_man = 2 
			AND ua.user_data_time IS NOT NULL 
			AND ua.zm_auth_time IS NOT NULL 
			AND DATE(ua.real_name_time)= DATE(ua.link_man_time) 
			AND DATE(ua.real_name_time)= DATE(ua.user_data_time) 
			AND DATE(ua.real_name_time)= DATE(ua.zm_auth_time) 
			AND u.register_source IS NOT NULL 
			AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY  DATE(GREATEST(user_data_time,zm_auth_time))
			 ORDER BY DATE(GREATEST(user_data_time,zm_auth_time)) `
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&ius)
	return
}

//获取当天立即点击申请用户数量
func GetSCApplyNowUsersByCondition(condition string, params []string, pkgParams []int) (fus []WrRegisterUser, err error) {
	sql := `SELECT t.first_loan_date AS create_date,COUNT(1) AS count
			FROM
				users_first_loan_time t
			INNER JOIN users u ON t.uid = u.id
			WHERE
				 u.register_source IS NOT NULL 
			AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY t.first_loan_date ORDER BY t.first_loan_date`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRows(&fus)
	return
}

type PackageName struct {
	PkgId   int
	PkgName string
}

//查询所有包名
func GetAllPackageName() (pkgName []int, err error) {
	sql := `SELECT register_source AS pkg_name FROM users WHERE register_source IS NOT NULL  GROUP BY pkg_name ORDER BY pkg_name `
	_, err = orm.NewOrm().Raw(sql).QueryRows(&pkgName)
	return
}

type SourceName struct {
	Source string
	Name   string
}

//根据包名查询市场
func GetAllSourceByPackage(condition string, params []int) (source []string, err error) {
	sql := `SELECT source FROM users WHERE source!='' `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY source`
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&source)
	return
}

func GetWrSCRegisterAllCount(source string, params []string, pkgParams []int) (weiRongDataAll *WeiRongDataAll, err error) {
	sql := `SELECT a.register_count,b.ocr_count,c.applynow_count,d.loan_user FROM `
	sql += `(SELECT COUNT(1) AS register_count FROM users u WHERE create_date<CURDATE() AND u.register_source IS NOT NULL AND out_put_source=""`
	if source != "" {
		sql += source
	}
	sql += `) a,(SELECT COUNT(uid) AS ocr_count FROM users u INNER JOIN users_auth ua ON u.id=ua.uid WHERE is_real_name=2 AND real_name_time<curdate() AND u.register_source IS NOT NULL  AND out_put_source=""`
	if source != "" {
		sql += source
	}
	sql += ` ) b,(SELECT COUNT(1) AS applynow_count
			FROM
				users_first_loan_time t
			INNER JOIN users u ON t.uid = u.id
			WHERE
				t.first_loan_date < CURDATE()
			AND u.register_source IS NOT NULL  AND out_put_source=""`
	if source != "" {
		sql += source
	}
	sql += ` ) c,(SELECT COUNT(1) AS loan_user FROM business_loan_h5 bl inner join users u on bl.uid=u.id WHERE  bl.create_time<curdate() AND  u.register_source IS NOT NULL  AND out_put_source=""`
	if source != "" {
		sql += source
	}
	sql += ` ) d`
	err = orm.NewOrm().Raw(sql, params, pkgParams, params, pkgParams, params, pkgParams, params, pkgParams).QueryRow(&weiRongDataAll)
	if source != "" {
		return
	}
	if data, err2 := json.Marshal(weiRongDataAll); err == nil && err2 == nil && utils.Re == nil {
		cache := utils.WEIRONGCOUNTSOURCE
		h := 24 - time.Now().Hour()
		utils.Rc.Put(cache, data, time.Duration(h)*time.Hour)
	}
	return
}

func GetSCDailyDataCache(condition string, params []string, pkgParams []int) (dailyDataList []DailyData, err error) {
	sql := `select day1 as date,register_count,login_count,ocr_count,identify_count,applynow_count,pro_register_count,CASE WHEN register_count!=0 THEN pro_register_count/register_count ELSE 0 END AS user_per_count,CASE WHEN register_count!=0 THEN h.money/register_count ELSE 0 END AS user_per_profit from
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
	(SELECT DATE_FORMAT(p.create_time,'%Y-%m-%d') AS day6,COUNT(product_id) AS pro_register_count
			FROM product_register_user AS p 
			INNER JOIN users AS u ON p.uid=u.id
			WHERE DATE_FORMAT(p.create_time ,'%Y-%m-%d')=u.create_date AND p.create_time<CURDATE() `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day6) g 
	ON a.day1=g.day6
	LEFT JOIN
	(SELECT DATE_FORMAT(pru.create_time, '%Y-%m-%d') AS day7,
		SUM(CASE
		WHEN dd.joint_mode = 1 THEN
			dd.cpa_price
		WHEN dd.joint_mode = 2 THEN
			dd.cps_first_per
		ELSE
			dd.cpa_price + dd.cps_first_per
		END) AS money
		FROM
			daily_data dd
		INNER JOIN agent_product ap ON dd.agent_product_id = ap.id
		INNER JOIN product_register_user pru ON ap.pro_id = pru.product_id
		INNER JOIN users u ON pru.uid=u.id
		WHERE pru.create_time<CURDATE()`
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY day7) h ON a.day1=h.day7 ORDER BY date DESC`
	_, err = orm.NewOrm().Raw(sql, params, pkgParams, params, pkgParams, params, pkgParams, params, pkgParams, params, pkgParams, params, pkgParams, params, pkgParams).QueryRows(&dailyDataList)
	if condition != "  AND u.source !='' " {
		return
	}
	if data, err2 := json.Marshal(dailyDataList); err == nil && err2 == nil && utils.Re == nil {
		cache := utils.WEIRONGEVERYDAYDATASOURCE
		h := 24 - time.Now().Hour()
		utils.Rc.Put(cache, data, time.Duration(h)*time.Hour)
	}
	return
}

type HistogramData struct {
	Name  string
	Count float64
}

// 累计微融注册用户
func GetWrSCRegisterCount(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u WHERE register_source IS NOT NULL AND out_put_source='' `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//累计完成orc实名认证用户
func GetWrSCORCCount(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT COUNT(1) FROM users u INNER JOIN users_auth a ON u.id=a.uid WHERE is_real_name=2 AND u.register_source IS NOT NULL AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//累计完成认证用户
func GetWrSCIdentifyCount(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT  COUNT(1) AS count
			FROM users_auth ua
			INNER JOIN users u ON ua.uid = u.id
			WHERE ua.is_real_name = 2 
			AND ua.is_user_data = 2 
			AND ua.is_zm_auth = 2 
			AND ua.is_link_man = 2 
			AND ua.user_data_time IS NOT NULL 
			AND ua.zm_auth_time IS NOT NULL 
			AND DATE(ua.real_name_time)= DATE(ua.link_man_time) 
			AND DATE(ua.real_name_time)= DATE(ua.user_data_time) 
			AND DATE(ua.real_name_time)= DATE(ua.zm_auth_time) 
			AND u.register_source IS NOT NULL 
			AND u.out_put_source=""  `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

// 累计点击立即申请用户
func GetSCApplyNowUserCount(condition string, params []string, pkgParams []int) (count int, err error) {
	sql := `SELECT COUNT(1) AS count
			FROM
				users_first_loan_time t
			INNER JOIN users u ON t.uid = u.id
			WHERE
				u.register_source IS NOT NULL 
			AND u.out_put_source="" `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params, pkgParams).QueryRow(&count)
	return
}

//查询微融市场历史数据
func GetWrScHistoryData(pkgName int, market string) (wrDailyData []DailyData, err error) {
	sql := `SELECT 
			date_date AS date,
			register_count,
			ocr_count,
			auth_count AS identify_count,
			apply_count AS applynow_count,
			pro_register_count,
			user_per_count,
			user_per_profit,
			total_profit
		FROM 
			channel_market_data
		WHERE 
			pkg_name=?
		AND 
			market=?`
	o := orm.NewOrm()
	o.Using("wr_backup")
	_, err = o.Raw(sql, pkgName, market).QueryRows(&wrDailyData)
	return
}
