package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type PlaCredit struct {
	Id                     int
	PlateformName          string //征信口子名称
	CreateTime             time.Time
	HasCreditInvestigation int64
}

//征信口子列表
func GetPlaCreditList(condition string, params interface{}, begin, end int) (placredit []*PlaCredit, err error) {
	sql := `SELECT id,plateform_name,has_credit_investigation FROM plateform_credit WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += `ORDER BY id DESC LIMIT ?,?`
	_, err = orm.NewOrm().Raw(sql, params, begin, end).QueryRows(&placredit)
	return
}

//征信口子总数
func PlaCreditCount(condition string, params interface{}) (count int, err error) {
	sql := `SELECT COUNT(1) FROM plateform_credit WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//新增征信口子
func SavePlaCredit(plaCredit *PlaCredit) error {
	sql := `INSERT INTO plateform_credit (plateform_name,has_credit_investigation,create_time)VALUES(?,?,NOW())`
	_, err := orm.NewOrm().Raw(sql, plaCredit.PlateformName, plaCredit.HasCreditInvestigation).Exec()
	if err != nil {
		return err
	}
	return nil
}

//修改征信口子
func UpdatePlaCredit(plaCredit *PlaCredit) error {
	sql := `UPDATE plateform_credit SET plateform_name=?,has_credit_investigation=? WHERE id = ? `
	_, err := orm.NewOrm().Raw(sql, plaCredit.PlateformName, plaCredit.HasCreditInvestigation, plaCredit.Id).Exec()
	if err != nil {
		return err
	}
	return nil
}

//根据id删除口子
func DelPlaCreditById(id int) (err error) {
	o := orm.NewOrm()
	sql := `DELETE FROM plateform_credit WHERE id = ?`
	_, err = o.Raw(sql, id).Exec()
	return
}

func GetAllPlaCredit() (plaCredit []PlaCredit, err error) {
	sql := `SELECT plateform_name FROM plateform_credit`
	orm.NewOrm().Raw(sql).QueryRows(&plaCredit)
	return
}
