package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ProductLog struct {
	Id              int
	UserId          int    //操作人id
	ProId           int    //产品id
	UserName        string //操作人名称
	UserDisplayName string
	UserIp          string
	Action          string
	Logger          string
	UrlPath         string `orm:"column(urlpath)"`
	Message         string
	FromParams      string
	QueryStrings    string
	CreateTime      time.Time
	BusinessId      int //业务code
	OperateType     int //1:上线 2：下线 3：弹窗4：不弹窗5：首页6：贷款页7：导量上限8：展示顺序9：产品推广
	IsIndexShow     int
	LoanSort        int
	Sort            int
	LargeLoanSort   int //大额贷款页顺序
	UseTime         int64
	IsNew           int
}

//写入产品日志
func SaveProductLog(productLog *ProductLog) error {
	o := orm.NewOrm()
	sql := `INSERT INTO sys_product_log (
	user_id,pro_id,user_name,user_display_name,
	user_ip,action,logger,
	urlpath,message,from_params,
	query_strings,business_id,operate_type,create_time,is_index_show,loan_sort,sort,large_loan_sort,use_time,is_new
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	o.Using("wr_log")
	_, err := o.Raw(sql,
		productLog.UserId, productLog.ProId, productLog.UserName, productLog.UserDisplayName,
		productLog.UserIp, productLog.Action, productLog.Logger,
		productLog.UrlPath, productLog.Message, productLog.FromParams,
		productLog.QueryStrings, productLog.BusinessId,
		productLog.OperateType, productLog.CreateTime, productLog.IsIndexShow,
		productLog.LoanSort, productLog.Sort, productLog.LargeLoanSort, productLog.UseTime, productLog.IsNew).Exec()
	return err
}

//查出最近的上线或者下线操作的时间
func GetOperateTime(proId int) (dataTime time.Time, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `select MAX(create_time) from sys_product_log s where s.operate_type=1 and s.pro_id=? `
	err = o.Raw(sql, proId).QueryRow(&dataTime)
	if err != nil {
		beego.Debug(err.Error())
	}
	return
}

//获取最近的上线平台不同位置的时间
func GetNotEquTime(condition string, proId, isIndexShow, loanSort, sort, largeLoanSort int) (dataTime time.Time, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `select MAX(create_time) from sys_product_log s where (s.operate_type=1 OR s.operate_type=2 OR s.operate_type=5 OR s.operate_type=6 OR s.operate_type=8 OR s.operate_type=10) and s.pro_id=? `
	sql += condition
	err = o.Raw(sql, proId, isIndexShow, loanSort, sort, largeLoanSort).QueryRow(&dataTime)
	return
}
