package models

import (
	//"fmt"

	"time"

	"github.com/astaxie/beego/orm"
)

/*
机构相关接口
*/
type Organization struct {
	Id          int       `form:"OrgId"`   //产品ID
	Name        string    `form:"OrgName"` //组织机构名称
	Address     string    //通讯录地址
	Email       string    `form:"OrgEmail"` //邮箱
	Create_time time.Time //创建时间
}

// 新增机构信息
func AddOrganization(organization *Organization) (id int, err error) {
	sql := `INSERT INTO organization (name, address, email, create_time) VALUES (?,?,?,now());`
	result, err := orm.NewOrm().Raw(sql, organization.Name, organization.Address, organization.Email).Exec()
	if err != nil {
		return 0, err
	}
	i, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(i), err
}

//根据id查询机构信息
func SearchOrgById(id int) (organization Organization, err error) {
	organization.Id = id
	err = orm.NewOrm().Read(&organization)
	return
}

//根据id更新机构信息
func UpdateOrgById(org *Organization) (err error) {
	_, err = orm.NewOrm().Update(org, "Name", "Address", "Email")

	return
}
