package models

import (
	"wr_v1/utils"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	orm.RegisterDataBase("default", "mysql", utils.MYSQL_URL)
	orm.RegisterDataBase("wr_log", "mysql", utils.MYSQL_LOG_URL)
	orm.RegisterDataBase("wr_backup", "mysql", utils.MYSQL_BACKUP_URL)
	orm.RegisterModel(
		new(ProductCleaning),
		new(ProductUrl),
		new(ProductForAdd),
		new(BusinessLinkman),
		new(Organization),
		new(ProductAgentUrl),
		new(Agent),
		new(AgentProduct),
		new(AgentDailyData),
		new(Images),
		new(JpushRecord),
	)
}
