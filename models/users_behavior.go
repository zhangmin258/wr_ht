package models

import (
	"github.com/astaxie/beego/orm"
)

type UsersBehavior struct {
	Uid        int    //用户编号
	Account    string //手机号
	Name       string //姓名
	Action     int    //意愿类型  1:页面停留时长；2：咨询客服；3：进入支付流程
	CreateTime string //创建时间
}

type BehaviorUsersExcel struct {
	CreateTime string //创建时间
	Account    string //手机号码
	Name       string //姓名
	Action     int    //意愿类型
}

//获取意愿用户列表
func GetBehaviorUsersList(condition string, params []string) (behaviorUsers []UsersBehavior, err error) {
	sql := `SELECT 
		ub.uid,
		um.account,
		um.verify_real_name AS name,
		ub.action,
		ub.create_time 
		FROM users_behavior AS ub 
		LEFT JOIN users_metadata AS um
		ON ub.uid=um.uid 
		WHERE ub.create_time >=CURDATE() `
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY ub.create_time DESC "
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&behaviorUsers)
	return
}

//获取意愿用户数量
func GetBehaviorUsersCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) count
			FROM users_behavior AS ub 
			LEFT JOIN users_metadata AS um
			ON ub.uid=um.uid 
			WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//导出意愿用户列表
func ExportBehaviorUsersList(condition string, params []string) (behaviorUsersExcel []BehaviorUsersExcel, err error) {
	sql := `SELECT 
		ub.create_time,
		um.account,
		um.verify_real_name AS name,
		ub.action
		FROM users_behavior AS ub 
		LEFT JOIN  users_metadata AS um 
		ON ub.uid=um.uid
		WHERE 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY ub.create_time DESC"
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&behaviorUsersExcel)
	return
}
