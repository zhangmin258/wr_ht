package controllers

import (
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

/*
businessLoan(订单)相关接口
*/
type BusinessLoanController struct {
	BaseController
}

//条件分页订单列表
//@router /getBusinessLoanList [get]
func (c *BusinessLoanController) GetBusinessLoanList() {
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//得到日期参数
	date := c.GetString("date")
	c.Data["date"] = date
	//用户姓名
	if userName := c.GetString("userName"); userName != "" {
		condition += " and b.verify_real_name = ? "
		params = append(params, userName)
	}
	//手机号
	if account := c.GetString("account"); account != "" {
		condition += " and  b.account = ?"
		params = append(params, account)
	}
	//订单状态
	if state := c.GetString("orderState"); state != "" {
		condition += " and a.state = ? "
		params = append(params, state)
	}
	// 产品id
	productId, err := c.GetInt("productId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取代理产品Id失败", err.Error(), c.Ctx.Input)
	}
	cooperationType, err := c.GetInt("cooperationType")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取代理产品类型失败", err.Error(), c.Ctx.Input)
	}
	//查询列表
	var orderList []models.BusinessLoan
	var count int
	if cooperationType == 0 {
		orderList, err = models.GetBusinessLoanPageApi(productId, date, condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE), utils.PAGE_SIZE)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询贷款人异常！", err.Error(), c.Ctx.Input)
		}
		//查询总记录数
		count, err = models.GetBusinessLoanListApi(productId, date, condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询所有贷款人的数量异常！", err.Error(), c.Ctx.Input)
		}
	} else if cooperationType == 1 {
		orderList, err = models.GetBusinessLoanPageH5(productId, date, condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE), utils.PAGE_SIZE)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询贷款人异常！", err.Error(), c.Ctx.Input)
		}
		//查询总记录数
		count, err = models.GetBusinessLoanListH5(productId, date, condition, params)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询所有贷款人的数量异常！", err.Error(), c.Ctx.Input)
		}
	}

	//查询总页数
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询页数失败", err.Error(), c.Ctx.Input)
	}
	c.Data["cooperationType"] = cooperationType
	c.Data["productId"] = productId
	c.Data["pageNum"] = pageNum
	c.Data["orderList"] = orderList
	c.Data["count"] = count
	c.Data["pageCount"] = pageCount
	c.TplName = "agent-products/orderList.html"

}

//查看订单详情
//@router /getBusinessLoan [get]
func (c *BusinessLoanController) GetBusinessLoan() {
	c.IsNeedTemplate()
	id, err := c.GetInt("id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取id失败", err.Error(), c.Ctx.Input)
	}
	// pageNum, _ := c.GetInt("page", 1)
	// if pageNum < 1 {
	// 	pageNum = 1
	// }
	cooperationType, err := c.GetInt("cooperationType")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取代理产品类型失败", err.Error(), c.Ctx.Input)
	}
	var businessLoan models.BusinessLoan
	if cooperationType == 0 {
		// condition := ""
		// params := []string{}
		//根据id查看详细信息
		businessLoan, err = models.GetBusinessLoanById(id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查看详细信息异常！", err.Error(), c.Ctx.Input)
		}
		// //查询某一个贷款人所有订单
		// count, err := models.GetBusinessLoanListById(condition, params, id)
		// if err != nil {
		// 	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询某一个贷款人所有订单异常！", err.Error(), c.Ctx.Input)
		// }
		// //查询总页数
		// pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE)
		// if err != nil {
		// 	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询总页数失败", err.Error(), c.Ctx.Input)
		// }
		// c.Data["pageNum"] = pageNum
		// c.Data["pageCount"] = pageCount
	} else if cooperationType == 1 {
		businessLoan, err = models.GetBusinessLoanH5ById(id)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查看详细信息异常！", err.Error(), c.Ctx.Input)
		}
	}
	c.Data["count"] = 1
	c.Data["order"] = businessLoan
	c.TplName = "agent-products/order_details.html"
}
