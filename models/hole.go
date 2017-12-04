package models

import (
	"strings"

	"github.com/astaxie/beego/orm"
)

type HoleInfo struct {
	Id        int
	Logo      string // logo
	Name      string // 安卓口子名称
	Url       string // 安卓下载链接
	PhoneType int    // 1安卓，2iOS
}

// 新口子配置信息
type HoleConfig struct {
	Id               int
	ScorePrice       int     // 融豆价格
	MoneyPrice       float64 // 人民币价格
	DescribeTitle1   string  // 描述标题1
	DescribeTitle2   string  // 描述标题2
	DescribeContent1 string  // 描述内容1
	DescribeContent2 string  // 描述内容2
}

func GetAndroidHoleList(condition string, params interface{}, begin, end int) (holeInfo []*HoleInfo, err error) {
	sql := `SELECT id,name,shop_url AS url,logo FROM
			aso_android_data WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY id DESC LIMIT ?,? `
	_, err = orm.NewOrm().Raw(sql, params, begin, end).QueryRows(&holeInfo)
	return
}

func GetIosHoleList(condition string, params interface{}, begin, end int) (holeInfo []*HoleInfo, err error) {
	sql := `SELECT id,name,download_url AS url,logo FROM
			aso_data WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += ` ORDER BY id DESC LIMIT ?,? `
	_, err = orm.NewOrm().Raw(sql, params, begin, end).QueryRows(&holeInfo)
	return
}

//获取安卓口子总数
func GetAndroidHoleCount(condition string, params interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM aso_android_data WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//获取ios口子总数
func GetIosHoleCount(condition string, params interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM aso_data WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//保存口子信息
func SaveHole(hole HoleInfo) (err error) {
	o := orm.NewOrm()
	sql1 := `INSERT INTO aso_android_data (name,shop_url,logo)VALUES(?,?,?)`
	sql2 := `INSERT INTO aso_data (name,download_url,logo)VALUES(?,?,?)`
	if !strings.Contains(hole.Url, "http") {
		hole.Url = "http://" + hole.Url
	}
	if !strings.Contains(hole.Logo, "http") {
		hole.Logo = "http://" + hole.Logo
	}
	if hole.PhoneType == 1 {
		_, err = o.Raw(sql1, hole.Name, hole.Url, hole.Logo).Exec()
		if err != nil {
			return
		}
	}
	if hole.PhoneType == 2 {
		_, err = o.Raw(sql2, hole.Name, hole.Url, hole.Logo).Exec()
	}
	return
}

//修改口子信息
func UpdateHole(hole HoleInfo) (err error) {
	o := orm.NewOrm()
	sql1 := `UPDATE aso_android_data SET name=?,shop_url=?,logo=? WHERE id = ?`
	sql2 := `UPDATE aso_data SET name=?,download_url=?,logo=? WHERE id = ?`
	if !strings.Contains(hole.Url, "http") {
		hole.Url = "http://" + hole.Url
	}
	if !strings.Contains(hole.Logo, "http") {
		hole.Logo = "http://" + hole.Logo
	}
	if hole.PhoneType == 1 {
		_, err = o.Raw(sql1, hole.Name, hole.Url, hole.Logo, hole.Id).Exec()
		if err != nil {
			return
		}
	}
	if hole.PhoneType == 2 {
		_, err = o.Raw(sql2, hole.Name, hole.Url, hole.Logo, hole.Id).Exec()
	}
	return
}

//根据id删除口子
func DelHoleById(id, phoneType int) (err error) {
	o := orm.NewOrm()
	if phoneType == 1 {
		sql1 := `DELETE FROM aso_android_data WHERE id = ?`
		_, err = o.Raw(sql1, id).Exec()
		if err != nil {
			return
		}
	}
	if phoneType == 2 {
		sql2 := `DELETE FROM aso_data WHERE id = ?`
		_, err = o.Raw(sql2, id).Exec()
	}
	return
}

//根据id获取新口子配置
func GetHoleConfig(id int) (holeConfig *HoleConfig, err error) {
	sql := `SELECT id,score_price,money_price,describe_title1,describe_content1,describe_title2,describe_content2 FROM score_exchange_product WHERE id=? `
	err = orm.NewOrm().Raw(sql, id).QueryRow(&holeConfig)
	return
}

//保存新口子配置
func UpdateHoleConfig(holeConfig HoleConfig) error {
	sql := `UPDATE score_exchange_product SET score_price=?,money_price=?,describe_title1=?,describe_content1=?,describe_title2=?,describe_content2=?
	WHERE id=? `
	_, err := orm.NewOrm().Raw(sql, holeConfig.ScorePrice, holeConfig.MoneyPrice, holeConfig.DescribeTitle1,
		holeConfig.DescribeContent1, holeConfig.DescribeTitle2, holeConfig.DescribeContent2, holeConfig.Id).Exec()
	if err != nil {
		return err
	}
	return nil
}
