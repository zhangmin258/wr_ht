package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type LoginRecord struct {
	Id  int
	Uid int //用户id
	//Login_time    time.Time
	MobileType string // 手机型号
	DeviceCode string `orm:"column(devicecode)"` //设备标识
	//Account       string
	CreateTime     time.Time //创建时间
	App            int       //区分哪里登录
	MobileVersion  string    //登录手机app版本号
	AppVersion     string    //app版本号
	DeviceUniqueid string    //设备型号
	Location       string    //根据定位，获取到的定位市
	Address        string    //根据定位，获取到的定位详细地址
	Ip             string    //登录ip
	Stage          string    //loan:借款,auth:认证，login:登录,register:注册
	SilentLogin    int       //1代表是静默登录
}

// 登录历史
func GetLoginRecordList(uid int, condition string, pars []string, begin, count int) (list []LoginRecord, err error) {
	sql := `SELECT lr.*
			FROM login_record lr
			WHERE lr.uid=?`
	sql += condition
	sql += " ORDER BY lr.create_time DESC LIMIT ?, ?"
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err = o.Raw(sql, uid, pars, begin, count).QueryRows(&list)
	return
}

func GetLoginRecordCount(uid int, condition string, pars []string) int {
	sql := `SELECT count(1)
			FROM login_record lr
			WHERE lr.uid=?`
	sql += condition
	var count int
	o := orm.NewOrm()
	o.Using("wr_log")
	o.Raw(sql, uid, pars).QueryRow(&count)
	return count
}
