package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type MakeConn struct {
	PlatformName string    //平台名称
	FullName     string    //姓名
	Phone        string    //联系电话
	Email        string    //联系邮箱
	Remark       time.Time //详细描述
	CreateTime   string    //创建时间
}

func GetAllMakeConn(begin, size int) (makeConn []MakeConn, err error) {
	sql := `SELECT platform_name,full_name,phone,email,remark,create_time
		 FROM make_connection 
		 ORDER BY create_time DESC 
		 limit ?,?`
	_, err = orm.NewOrm().Raw(sql, begin, size).QueryRows(&makeConn)
	return
}

// func GetMakeConnCount() (count int, err error) {
// 	sql := `SELECT COUNT(1) FROM make_connection`
// 	err = orm.NewOrm().Raw(sql).QueryRow(&count)
// 	return
// }
