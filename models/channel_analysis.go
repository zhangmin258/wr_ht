package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type AnalysisChannel struct {
	Name                 string    //渠道名称
	HitsCount            int       //点击次数
	RegisteredUsersCount int       //注册用户数
	ActivatedUsersCount  int       //激活用户数
	Rac                  float64   //注册-激活转化率
	Hrc                  float64   //点击注册转化率
	Crpc                 float64   //人均注册平台数
	Cci                  float64   //人均创造收益
	PlatformsCount       int       //注册平台数
	Income_money         float64   //收益
	Consume_money        float64   //消耗
	UsersCount           int       //用户总数
	Date                 time.Time //日期
}

//收益消耗数
type ChannelIncomeAndConsume struct {
	Name         string
	IncomeMoney  float64
	ConsumeMoney float64
}

//渠道列表
type ChannelList struct {
	Name string //渠道名称
}

//用户注册平台数
type RegisteredPlatforms struct {
	Name  string
	Count int
	Date  time.Time
}

//渠道人数
type ChannelCount struct {
	Name  string
	Count int
}

//指定日期渠道人数
type ChannelCountDate struct {
	Name  string
	Count int
	Date  time.Time
}

//渠道收益消耗
type IncomeAndConsume struct {
	Name  string
	Money float64
}

//指定日期渠道收益消耗
type IncomeAndConsumeDate struct {
	Name    string
	Income  float64
	Consume float64
	Date    time.Time
}

//指定日期渠道点击注册激活人数
type ChannelHRA struct {
	Name                 string
	HitsCount            int
	RegisteredUsersCount int
	ActivatedUsersCount  int
	Date                 time.Time
}

//获取渠道列表
func GetAnalysisChannel(condition string) (channel []ChannelList, err error) {
	sql := `SELECT out_put_source AS name FROM users WHERE out_put_source !="" GROUP BY name	`
	sql += condition
	_, err = orm.NewOrm().Raw(sql).QueryRows(&channel)
	return
}

//获取渠道注册用户数
func GetH5RegisteredUsers_Count(condition string, params []string) (res []ChannelCount, err error) {
	sql := `SELECT out_put_source name,COUNT(1) count FROM users u  WHERE 1=1`
	sql += condition + " GROUP BY out_put_source"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取渠道按时间注册用户数
func GetH5RegisteredUsers_CountByTime(condition string, channel string, params []string) (res []ChannelHRA, err error) {
	sql := `SELECT  u.create_date AS date, COUNT(1) registered_users_count
	FROM users u WHERE  1=1`
	sql += condition + " GROUP BY u.create_date"
	_, err = orm.NewOrm().Raw(sql, channel, params).QueryRows(&res)
	return
}

//获取渠道激活用户数
func GetH5ActivatedUsers_Count(condition string, params []string) (res []ChannelCount, err error) {
	sql := `SELECT out_put_source name,COUNT(1) count
	FROM users u  WHERE active_time IS NOT NULL`
	sql += condition + " GROUP BY out_put_source"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取渠道按时间激活用户数
func GetH5ActivatedUsers_CountByTime(condition string, channel string, params []string) (res []ChannelHRA, err error) {
	sql := `SELECT create_date AS date,COUNT(1) activated_users_count
    FROM users u  WHERE active_time IS NOT NULL`
	sql += condition + " GROUP BY create_date"
	_, err = orm.NewOrm().Raw(sql, channel, params).QueryRows(&res)
	return
}

//获取收益
func GetH5IncomeMoney(condition string, params []string) (res []IncomeAndConsume, err error) {
	/*sql := `SELECT u.out_put_source name,SUM(income.money) money
	    FROM  users u
	    INNER JOIN (
		SELECT urp.uid,SUM(pc.cpa_price + pc.cps_price) money
		FROM users_register_platform urp
		LEFT JOIN product_cleaning pc ON urp.product_id = pc.product_id
		GROUP BY urp.uid
	)   income ON u.id =income.uid `*/
	sql := `SELECT u.out_put_source name,SUM(income.money) money
    FROM  users u
    INNER JOIN (
	SELECT urp.uid,SUM(pc.cpa_price + pc.cps_price) money
	FROM product_register_user urp
	LEFT JOIN product_cleaning pc ON urp.product_id = pc.product_id
	GROUP BY urp.uid
)   income ON u.id =income.uid WHERE 1=1`
	sql += condition + " GROUP BY u.out_put_source"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取渠道按时间收益
func GetH5IncomeMoneyByTime(condition string, channel string, params []string) (res []IncomeAndConsumeDate, err error) {
	sql := `SELECT
	income.create_date AS date,SUM(income.money) income
    FROM users u
    LEFT JOIN (
	SELECT urp.uid,SUM(pc.cpa_price + pc.cps_price) money,urp.create_date
	FROM users_register_platform urp
	LEFT JOIN product_cleaning pc ON urp.product_id = pc.product_id
	GROUP BY urp.uid
)   income ON u.id = income.uid WHERE 1=1`
	sql += condition + " GROUP BY income.create_date"
	_, err = orm.NewOrm().Raw(sql, channel, params).QueryRows(&res)
	return
}

//获取损耗
func GetH5ConsumeMoney(condition string, params []string) (res []IncomeAndConsume, err error) {
	/*sql := `SELECT u.out_put_source name,SUM(consume.money) money
		FROM users u
		INNER JOIN (
		SELECT urp.uid,sum(ap.cpa_price) money
		FROM users_register_platform urp
		INNER JOIN agent_product ap ON urp.product_id = ap.pro_id
		GROUP BY urp.uid
	)   consume ON u.id = consume.uid`*/
	sql := `SELECT u.out_put_source name,SUM(consume.money) money
	FROM users u
	INNER JOIN (
	SELECT urp.uid,sum(ap.cpa_price) money
	FROM users_register_platform urp
	INNER JOIN agent_product ap ON urp.product_id = ap.pro_id
	GROUP BY urp.uid
)   consume ON u.id = consume.uid WHERE 1=1`
	sql += condition + " GROUP BY u.out_put_source"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return

}

//获取渠道按时间损耗
func GetH5ConsumeMoneyByTime(condition string, channel string, params []string) (res []IncomeAndConsumeDate, err error) {
	/*sql := `SELECT DATE_FORMAT(consume.create_time,'%Y-%m-%d') date,SUM(consume.money) consume
	    FROM users u
	    LEFT JOIN (
		SELECT urp.uid,sum(ap.cpa_price) money,urp.create_time
		FROM users_register_platform urp
		LEFT JOIN agent_product ap ON urp.product_id = ap.pro_id
		GROUP BY urp.uid
	)   consume ON u.id = consume.uid`*/
	sql := `SELECT DATE_FORMAT(consume.create_time,'%Y-%m-%d') date,SUM(consume.money) consume
    FROM users u
    LEFT JOIN (
	SELECT urp.uid,sum(ap.cpa_price) money,urp.create_time
	FROM users_register_platform urp
	LEFT JOIN agent_product ap ON urp.product_id = ap.pro_id
	GROUP BY urp.uid
)   consume ON u.id = consume.uid WHERE 1=1`
	sql += condition + " GROUP BY DATE_FORMAT(consume.create_time,'%Y-%m-%d')"
	_, err = orm.NewOrm().Raw(sql, channel, params).QueryRows(&res)
	return
}

//获取渠道按时间点击人数
func GetHits_CountByTime(condition string, channel string, params []string) (res []ChannelHRA, err error) {
	sql := `SELECT plpr.create_date AS date,count(plpr.uid) hits_count
	FROM product_landing_page_record plpr
	INNER JOIN users u ON u.id = plpr.uid WHERE 1=1 `
	sql += condition + " GROUP BY plpr.create_date "
	_, err = orm.NewOrm().Raw(sql, channel, params).QueryRows(&res)
	return
}

//获取渠道注册平台数
func GetH5RegisteredPlatforms_Count(condition string, params []string) (res []RegisteredPlatforms, err error) {
	sql := `SELECT u.out_put_source name,COUNT(urp.id) count
	FROM users u
	LEFT JOIN product_register_user urp ON u.id=urp.uid  WHERE 1=1`
	sql += condition + " GROUP BY u.out_put_source"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&res)
	return
}

//获取渠道按时间注册平台数
func GetH5RegisteredPlatforms_CountByTime(condition string, channel string, params []string) (res []RegisteredPlatforms, err error) {
	sql := `SELECT DATE_FORMAT(urp.create_time,'%Y-%m-%d') date ,COUNT(urp.id) count
	FROM users u
	LEFT JOIN product_register_user urp ON u.id=urp.uid WHERE 1=1`
	sql += condition + " GROUP BY DATE_FORMAT(urp.create_time,'%Y-%m-%d')"
	_, err = orm.NewOrm().Raw(sql, channel, params).QueryRows(&res)
	return
}
