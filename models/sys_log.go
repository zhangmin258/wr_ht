package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"time"
	"wr_v1/utils"
)

type SysLog struct {
	Id              int       `xorm:"int pk autoincr 'id'"`
	UserId          int       `xorm:"int 'user_id'"`
	UserName        string    `xorm:"varchar(50) 'user_name'"`
	UserDisplayName string    `xorm:"varchar(50) 'user_display_name'"`
	UserIp          string    `xorm:"varchar(50) 'user_ip'"`
	Action          string    `xorm:"varchar(50) 'action'"`
	Logger          string    `xorm:"varchar(50) 'logger'"`
	UrlPath         string    `xorm:"varchar(50) 'urlpath'"`
	Message         string    `xorm:"text 'message'"`
	FromParams      string    `xorm:"text 'from_params'"`
	QueryStrings    string    `xorm:"text 'query_strings'"`
	CreateTime      time.Time `orm:"column(create_time);type(datetime);null"`
	BusinessId      int       `xorm:"int 'business_id'"`
}

var engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", utils.MYSQL_LOG_URL)
	if err != nil {
		panic(err.Error())
	}
}

// AddLogs insert a new Logs into database and returns
// last inserted Id on success.
func AddLogs(m *SysLog) (id int64, err error) {
	id, err = engine.InsertOne(m)
	return
}
