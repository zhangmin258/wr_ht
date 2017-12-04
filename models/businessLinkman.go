package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

/*
商务联系人相关接口
*/
type BusinessLinkman struct {
	Id          int       `form:"BusId"`   //id
	Name        string    `form:"BusName"` //商务联系人名称
	Phone       string    //手机号
	Wechat      string    //微信
	Qq          string    //qq
	Email       string    `form:"BusEmail"` //邮箱
	Create_time time.Time //创建时间
}

// 新增商务联系人信息
func AddBusinessLinkman(bus *BusinessLinkman) (id int, err error) {
	sql := `INSERT INTO business_linkman (name,phone,wechat,qq,email,create_time) VALUES (?,?,?,?,?, NOW());`
	result, err := orm.NewOrm().Raw(sql, bus.Name, bus.Phone, bus.Wechat, bus.Qq, bus.Email).Exec()
	if err != nil {
		return 0, err
	}
	i, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(i), err
}

//根据id查询商务联系人
func SearchBusinessLinkmanById(id int) (bus BusinessLinkman, err error) {
	bus.Id = id
	err = orm.NewOrm().Read(&bus)
	return
}

//根据id更新商务联系人
func UpdateBusById(bus *BusinessLinkman) (err error) {
	_, err = orm.NewOrm().Update(bus, "Name", "Phone", "Wechat", "Qq", "Email")
	return
}
