package models

import "github.com/astaxie/beego/orm"

// 获取累计注册用户
func QueryUsersRegisterCount() (count int, err error) {
	sql := `SELECT COUNT(1) FROM users`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取累计认证用户
func QueryUsersAuthCount() (count int, err error) {
	sql := `SELECT COUNT(1)  count
			FROM users_auth
			WHERE is_real_name = 2  AND is_user_data = 2 AND is_zm_auth = 2 AND is_link_man = 2 `
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//累计放款用户
func QueryUsersLoanCount() (count int, err error) {
	sql := `SELECT COUNT(1)  count
			FROM business_loan 
			WHERE state = 'CONFIRM' OR state = 'FINISH' OR state = 'OVERDUE'`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

type UsersDataDetail struct {
	RegisterCount int
	AuthCount     int
	LoanCount     int
	ActiveCount   int
	CreateDate    string
}

//获取用户数据详情
func QueryUsersDataDetail(start, end int) (data []UsersDataDetail, err error) {
	sql := `SELECT * FROM  users_data_detail ORDER BY create_date DESC LIMIT ?,? `
	o := orm.NewOrm()
	o.Using("wr_backup")
	_, err = o.Raw(sql, start, end).QueryRows(&data)
	return
}

//获取今日数据明细
func QueryTodayDailyDatas() (detailData UsersDataDetail, err error) {
	sql := `SELECT 
			  CURDATE() create_date,
			  (SELECT COUNT(1) FROM users WHERE create_date = CURDATE()) register_count, 
			  (SELECT COUNT(1) FROM users_auth WHERE is_real_name = 2 AND is_user_data = 2 AND is_zm_auth = 2 AND is_link_man = 2 AND DATE(real_name_time)=CURDATE() AND DATE(link_man_time)=CURDATE() AND DATE(user_data_time)=CURDATE() AND DATE(zm_auth_time)=CURDATE()) auth_count, 
			  (SELECT COUNT(1) FROM business_loan WHERE (state = 'CONFIRM' OR state = 'FINISH' OR state = 'OVERDUE') AND real_time >= CURDATE()) loan_count 
			FROM DUAL ;`
	err = orm.NewOrm().Raw(sql).QueryRow(&detailData)
	return

}

//活跃今日用户明细
func QueryTodayActiveData() (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) count FROM users WHERE login_time >= CURDATE()`
	err = o.Raw(sql).QueryRow(&count)
	return
}

func QueryUsersDataDetailCount() (count int, err error) {
	sql := `SELECT COUNT(1) FROM  users_data_detail`
	o := orm.NewOrm()
	o.Using("wr_backup")
	err = o.Raw(sql).QueryRow(&count)
	return
}
