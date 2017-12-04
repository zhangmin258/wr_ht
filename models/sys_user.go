package models

import (
	"github.com/astaxie/beego/orm"
)

type SysUser struct {
	Id            int
	Code          string `orm:"column(code);null"`
	Name          string `orm:"column(name);null"`
	Password      string `orm:"column(password);null"`
	DisplayName   string `orm:"column(displayname);null"`
	Sex           string `orm:"column(sex);null"`
	Phone         string `orm:"column(phone);null"`
	Email         string `orm:"column(email);null"`
	AccountStatus string `orm:"column(accountstatus);null"`
	Remark        string `orm:"column(remark);null"`
	RoleId        int    `orm:"column(role_id);null"`
	StationId     int    `orm:"column(station_id);null"`
}

// type SysUserExpand struct {
// 	Id            int
// 	Code          string `orm:"column(code);null"`
// 	Name          string `orm:"column(name);null"`
// 	Password      string `orm:"column(password);null"`
// 	Password1     string `orm:"column(password1);null"`
// 	DisplayName   string `orm:"column(displayname);null"`
// 	Sex           string `orm:"column(sex);null"`
// 	Phone         string `orm:"column(phone);null"`
// 	Email         string `orm:"column(email);null"`
// 	AccountStatus string `orm:"column(accountstatus);null"`
// 	Remark        string `orm:"column(remark);null"`
// 	RoleId        string `orm:"column(role_id);null"`
// }

type SysUserMini struct {
	Id            int
	Name          string
	Password      string
	Displayname   string
	Email         string
	Role_id       int
	RoleName      string
	Accountstatus string
	Secret        string // 验证码密钥
	AuthURL       string
	Station_id    int    //岗位ID
	StationName   string //岗位名称
}

func Login(name, password string) (v *SysUser, err error) {
	o := orm.NewOrm()
	sql := `SELECT * FROM sys_user WHERE name= ? and password=? and accountStatus='启用'  `
	// fmt.Println(sql, name, password)
	err = o.Raw(sql, name, password).QueryRow(&v)
	return
}

// 系统用户列表
func SysUserList(condition string, pars []string, begin, count int) (list []SysUserMini, err error) {
	sql := `SELECT su.*, sr.displayname role_name,ss.name as station_name
			FROM sys_user su 
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			LEFT JOIN sys_station ss ON su.station_id=ss.id
			WHERE 1=1`
	sql += condition
	sql += " ORDER BY create_time DESC LIMIT ?, ?"
	_, err = orm.NewOrm().Raw(sql, pars, begin, count).QueryRows(&list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func SysUserCount(condition string, pars []string) int {
	sql := `SELECT count(1) 
			FROM sys_user su
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			WHERE 1=1`
	sql += condition
	var count int
	orm.NewOrm().Raw(sql, pars).QueryRow(&count)
	return count
}

func SysUserDetail(uid int) (user *SysUserMini, err error) {
	sql := `SELECT su.* , sr.displayname role_name
			FROM sys_user su
			LEFT JOIN sys_role sr ON su.role_id=sr.id
			WHERE su.id=?`
	err = orm.NewOrm().Raw(sql, uid).QueryRow(&user)
	return
}

func (u *SysUserMini) Insert() error {
	sql := `INSERT INTO sys_user (name, password,  displayname, email, accountstatus, role_id, create_time,station_id)
			values(?, ?, ?, ?, ?, ?, now(),?)`
	_, err := orm.NewOrm().Raw(sql, u.Name, u.Password, u.Displayname, u.Email, u.Accountstatus, u.Role_id, u.Station_id).Exec()
	return err
}

func (u *SysUserMini) Update() error {
	sql := `UPDATE sys_user SET password=?, displayname=?, email=?, accountstatus=?, role_id=?,station_id=?
			WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, u.Password, u.Displayname, u.Email, u.Accountstatus, u.Role_id, u.Station_id, u.Id).Exec()
	return err
}

func DeleteUser(uid int) error {
	sql := `DELETE FROM sys_user WHERE id=?`
	_, err := orm.NewOrm().Raw(sql, uid).Exec()
	return err
}

//修改用户密码
func UpdatePassword(id int, password, password1 string) (bool, error) {
	o := orm.NewOrm()
	sql := `update sys_user set password=? where id=?`
	_, err := o.Raw(sql, password, id).Exec()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
