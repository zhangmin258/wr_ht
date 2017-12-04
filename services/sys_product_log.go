package services

import (
	"fmt"
	"net/url"
	"time"
	"wr_v1/models"

	"github.com/astaxie/beego/context"
)

func ProductLogs(user_id, pro_id, business_id, operate_type, is_index_show, loan_sort, sort, largeLoanSort int, useTime int64, isNew int, username, displayname, action, logger, message string, input *context.BeegoInput) error {
	ip := input.IP()
	urlpath := input.URL()
	querystrings := input.URI()
	fromparams, _ := url.QueryUnescape(string(input.RequestBody))
	log := &models.ProductLog{UserId: user_id,
		UserName: username, UserDisplayName: displayname,
		UserIp: ip, Action: action, Logger: logger,
		UrlPath: urlpath, Message: message,
		FromParams:   fromparams,
		QueryStrings: querystrings,
		CreateTime:   time.Now(), BusinessId: business_id,
		ProId: pro_id, OperateType: operate_type,
		IsIndexShow: is_index_show, LoanSort: loan_sort, Sort: sort, LargeLoanSort: largeLoanSort,
		IsNew: isNew, UseTime: useTime}
	err := models.SaveProductLog(log)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}
