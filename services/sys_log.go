package services

import (
	"encoding/json"
	"fmt"
	"wr_v1/models"
	"wr_v1/utils"
)

// the service for log
func AutoInsertLogToDB() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[AutoInsertLogToDB]", err)
		}
	}()
	for {
		utils.Rc.Brpop(utils.CacheKeySystemLogs, func(b []byte) {
			var log models.SysLog
			if err := json.Unmarshal(b, &log); err != nil {
				fmt.Println("json unmarshal wrong!")
			}
			if _, err := models.AddLogs(&log); err != nil {
				fmt.Println(err.Error(), log)
			}
		})
	}
}
