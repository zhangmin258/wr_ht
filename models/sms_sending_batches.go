package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type SendBatches struct {
	SmsId         string    //sms_id
	Plateform     int       //短信平台名称 ：1微融短信；2空间畅想；3云融正通
	Content       string    //短信文本
	PushCount     int       //短信推送总次数
	PushTime      time.Time //发送时间
	Source        string    //
	AccountSource string    //来源
	Name          string    //发送人名称
}
type SourceLoadCount struct {
	Source    string //source
	LoadCount int    //加载次数
}
type ProductClickCount struct {
	Proname    string //产品名称
	ClickCount int    //点击次数
}

//发送报表分页查询
func GetSmsMarketingData(condition string, params []interface{}, begin int, size int) (sendBatches []SendBatches, err error) {
	sql := `SELECT m.sms_id , m.push_time , m.plateform, m.content, SUM(m.push_count) as push_count,m.source,u.name,m.account_source FROM sms_management m LEFT JOIN sys_user u
			ON m.sys_user_id=u.id
			WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY m.flag ORDER BY m.push_time DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&sendBatches)
	if err != nil {
		return nil, err
	}
	return

}

//统计短信发送批次总数
func GetCountForSms(condition string, params []interface{}) (count int, err error) {
	sql := ` SELECT COUNT(1) AS count FROM (SELECT flag from sms_management WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += " GROUP BY flag) a"
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	if err != nil {
		return 0, err
	}
	return
}

//根据名称查询id
func GetPidByPname(contentName string) (productId int, err error) {
	sql := `SELECT id FROM product WHERE name=?`
	err = orm.NewOrm().Raw(sql, contentName).QueryRow(&productId)
	if err != nil {
		return 0, err
	}
	return
}

//根据source查询点击次数
func GetClickCount(source string) (count int, err error) {
	sql := `SELECT SUM(a.click_count) AS count
		FROM (
		SELECT COUNT(DISTINCT(pr.ip)) click_count
		FROM product_recommend_click_record pr
		WHERE pr.source = ? AND pr.type =1
		GROUP BY pr.productid) a `
	err = orm.NewOrm().Raw(sql, source).QueryRow(&count)
	return
}

//查询加载次数
func GetLoadCount(source string) (count int, err error) {
	sql := `SELECT COUNT(ip) AS count FROM product_recommend_click_record WHERE type=2 AND source = ? `
	err = orm.NewOrm().Raw(sql, source).QueryRow(&count)
	return
}

// 获取到达成功数量
func GetSuccessCount(sms_ids []string) (count int, err error) {
	o := orm.NewOrm()
	statusReport := "DELIVRD"
	var c int
	sql := ` SELECT COUNT(1) FROM sms_report WHERE status_report= ? AND sms_id = ?  `
	for _, k := range sms_ids {
		err = o.Raw(sql, statusReport, k).QueryRow(&c)
		if err != nil {
			return
		}
		count = count + c
	}
	return
}

func GetSMSFlag(smsId string) (smsIds []string, err error) {
	o := orm.NewOrm()
	o.Begin()
	defer func() {
		if err != nil {
			o.Rollback()
			return
		}
		o.Commit()
	}()
	//根据smsId找到flag
	var flag string
	sql := ` SELECT flag FROM sms_management WHERE sms_id = ?  `
	err = o.Raw(sql, smsId).QueryRow(&flag)
	if err != nil {
		return
	}
	sql = "SELECT sms_id FROM sms_management WHERE flag = ? "
	_, err = o.Raw(sql, flag).QueryRows(&smsIds)
	if err != nil {
		return
	}
	return
}

//  根据source 获取连接详情  产品名称，页面加载次数
func GetLinkStaLoadCount(source string, condition string, params []string) (ld_load SourceLoadCount, err error) {
	o := orm.NewOrm()
	sql := ` SELECT source , p.name proname,count(distinct(pr.ip)) load_count
			FROM  product_recommend_click_record pr INNER JOIN product p  ON pr.productid = p.id
			WHERE pr.source = ?  AND pr.type =2   `
	if condition != "" {
		sql += condition
	}
	err = o.Raw(sql, source, params).QueryRow(&ld_load)
	return
}

//  根据source 获取连接详情  产品名称，点击次数
func GetLinkStaClickCount(source string, condition string, params []string) (ld_click []ProductClickCount, err error) {
	o := orm.NewOrm()
	sql := ` SELECT source , p.name proname,count(distinct(pr.ip)) click_count
			FROM  product_recommend_click_record pr INNER JOIN product p  ON pr.productid = p.id
			WHERE pr.source = ?  AND pr.type =1  `
	if condition != "" {
		sql += condition
	}
	sql += `  GROUP BY  pr.productid `
	_, err = o.Raw(sql, source, params).QueryRows(&ld_click)
	return
}
