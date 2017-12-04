package controllers

import (
	"strconv"
	"strings"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type AgentProController struct {
	BaseController
}

//查看下级代理商代理产品
func (c *AgentProController) GetAgentProduct() {
	c.IsNeedTemplate()
	agentId, _ := c.GetInt("AgentId")
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//代理产品名称
	if proName := c.GetString("proName"); proName != "" {
		condition += " and p.name =?"
		params = append(params, proName)
	}
	list, err := models.GetAgentProductList(agentId, condition, params, utils.StartIndex(pageNum, utils.PageSize5), utils.PageSize5)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询代理产品异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetAgentProductCount(agentId, condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询代理产品数量异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PageSize5)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取页数失败", err.Error(), c.Ctx.Input)
	}
	flag := false
	for _, v := range list { //判断是否有微融产品
		if v.AgentId == 0 {
			flag = true
			break
		}
	}
	c.Data["flag"] = flag
	c.Data["agentProductList"] = list
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["agentId"] = agentId

	c.TplName = "channel-management/product_agent.html"
}

//新增下级代理产品页面
func (c *AgentProController) AddAgentProduct() {
	c.IsNeedTemplate()
	agentId, err := c.GetInt("agentId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取代理商ID失败", err.Error(), c.Ctx.Input)
	}
	c.Data["agentId"] = agentId
	//获取标识
	wr := c.GetString("wr")
	if wr == "1" { //跳转到添加代理微融页面
		c.Data["weirongUrl"] = utils.OutPutURL
		c.TplName = "channel-management/add_wr_product.html"
		return
	}
	c.TplName = "channel-management/add_agent_product.html"
}

//編輯下级代理产品页面
func (c *AgentProController) EditAgentProduct() {
	c.IsNeedTemplate()
	var agentProduct models.AgentProduct
	c.ParseForm(&agentProduct)
	//查询代理信息
	err := models.GetAgentProductById(&agentProduct)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询代理产品异常！", err.Error(), c.Ctx.Input)
	}
	//查询代理产品
	productName, err := models.GetProductNameById(agentProduct.ProId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询产品名称异常！", err.Error(), c.Ctx.Input)
	}
	//查询当前代理当前产品所有已使用url
	urlIdList := strings.Split(agentProduct.UrlId, ",")
	urlList, err := models.GetUrlById(urlIdList)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询当前产品所有url异常！", err.Error(), c.Ctx.Input)
	}
	//查询当前产品所有可用url
	allUrlList, err := models.GetUrlByProId(agentProduct.ProId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "当前产品所有可用url异常！", err.Error(), c.Ctx.Input)
	}
	for _, url := range urlList {
		allUrlList = append(allUrlList, url)
	}
	c.Data["agentProduct"] = agentProduct
	c.Data["urlList"] = urlList
	c.Data["allUrlList"] = allUrlList
	c.Data["productName"] = productName
	c.TplName = "agent-products/edit_agent_product.html"

}

//模板函数
func Length(a []models.ProductUlr) (i int) {
	i = len(a)
	return
}

//保存代理商的代理产品信息
func (c *AgentProController) AgentProductSave() {
	defer c.ServeJSON()
	//接受参数
	var agentProduct models.AgentProduct
	err := c.ParseForm(&agentProduct)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "解析参数失败！"}
	}
	agentProduct.AgentTime = time.Now()
	//保存
	if agentProduct.UrlId == "-no more url can be used-" {
		agentProduct.UrlId = ""
	}
	err = models.AddAgentProduct(&agentProduct)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "当前产品所有可用url异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "插入数据失败！"}
	}
	//将所使用的url设为已分配
	var urlId int
	urlIds := []int{}
	if strings.Contains(agentProduct.UrlId, ",") {
		url := strings.Split(agentProduct.UrlId, ",")
		for _, u := range url {
			urlId, _ = strconv.Atoi(u)
			urlIds = append(urlIds, urlId)
		}
	} else {
		urlId, _ = strconv.Atoi(agentProduct.UrlId)
		urlIds = append(urlIds, urlId)
	}
	err = models.ChangeAgentProUrlById(urlIds)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改url状态失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "修改链接使用状态失败！"}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "操作成功"}
}

//保存代理微融产品的数据
func (c *AgentProController) AgenWrSave() {
	defer c.ServeJSON()
	//接受参数
	var agentProduct models.AgentProduct
	err := c.ParseForm(&agentProduct)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "解析参数失败！"}
	}
	agentProduct.AgentTime = time.Now()

	//保存新增的url
	urlList := []string{}
	if strings.Contains(agentProduct.UrlId, ",") {
		url := strings.Split(agentProduct.UrlId, ",")
		for _, u := range url {
			urlList = append(urlList, u)
		}
	} else {
		urlList = append(urlList, agentProduct.UrlId)
	}
	idList, err := models.AddAgentProUrl(urlList)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "添加代理商代理产品url异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "修改链接使用状态失败！"}
	}
	var urlIdStr string
	for index, id := range idList {
		var str string
		if index == len(idList)-1 {
			str = strconv.Itoa(id)
		} else {
			str = strconv.Itoa(id) + ","
		}
		urlIdStr += str
	}
	agentProduct.UrlId = urlIdStr
	agentProduct.ProId = 0
	err = models.AddAgentProduct(&agentProduct)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "添加代理商代理产品异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "插入数据失败！"}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "操作成功"}
}

//解除商品合作
func (c *AgentProController) DelAgentPro() {
	defer c.ServeJSON()
	id, _ := c.GetInt("Id")
	urls, err := models.GetUrlductById(id)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据Id查找urlId异常！", err.Error(), c.Ctx.Input)
		}
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "解析参数失败"}
	}
	//拆分urlId
	urlIds := []int{}
	if strings.Contains(urls, ",") {
		url := strings.Split(urls, ",")
		for _, u := range url {
			urlId, _ := strconv.Atoi(u)
			urlIds = append(urlIds, urlId)
		}
	} else {
		urlId, _ := strconv.Atoi(urls)
		urlIds = append(urlIds, urlId)
	}
	err = models.DelAgentProById(id, urlIds)
	if err != nil {

		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "代理商解除代理产品异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "解除产品异常"}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "解除成功"}
}

//更新下级代理商品
func (c *AgentProController) AgentProductUpdate() {
	defer c.ServeJSON()
	//接受参数
	var agentProduct models.AgentProduct
	err := c.ParseForm(&agentProduct)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "解析参数失败！"}
	}
	err = models.UpdateAgentProduct(&agentProduct)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改代理商代理产品异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "修改数据失败！"}
	}
	//将所使用旧的url设为未分配
	oldUrlIdStr := c.GetString("OldUrlId")
	oldUrlIds := []int{}
	if strings.Contains(oldUrlIdStr, ",") {
		url := strings.Split(oldUrlIdStr, ",")
		for _, u := range url {
			urlId, _ := strconv.Atoi(u)
			oldUrlIds = append(oldUrlIds, urlId)
		}
	} else {
		urlId, _ := strconv.Atoi(oldUrlIdStr)
		oldUrlIds = append(oldUrlIds, urlId)
	}
	err = models.ChangeOldAgentProUrlById(oldUrlIds)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改链接使用状态异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "修改链接使用状态失败！"}
	}
	//将新的url设为已分配
	urlIds := []int{}
	if strings.Contains(agentProduct.UrlId, ",") {
		url := strings.Split(agentProduct.UrlId, ",")
		for _, u := range url {
			urlId, _ := strconv.Atoi(u)
			urlIds = append(urlIds, urlId)
		}
	} else {
		urlId, _ := strconv.Atoi(agentProduct.UrlId)
		urlIds = append(urlIds, urlId)
	}
	err = models.ChangeAgentProUrlById(urlIds)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改链接使用状态异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "修改链接使用状态失败！"}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "操作成功"}
}

//編輯代理微融产品页面
func (c *AgentProController) EditWrProduct() {
	c.IsNeedTemplate()
	var agentProduct models.AgentProduct
	c.ParseForm(&agentProduct)
	//查询代理信息
	err := models.GetAgentProductById(&agentProduct)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询代理产品异常！", err.Error(), c.Ctx.Input)
	}
	//查询当前代理当前产品所有已使用url
	urlIdList := strings.Split(agentProduct.UrlId, ",")
	urlList, err := models.GetUrlById(urlIdList)
	for i := 0; i < len(urlList); i++ {
		source := strings.Split(urlList[i].Url, "source=")[1]
		urlList[i].Url = source
	}
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id获取Url失败", err.Error(), c.Ctx.Input)
	}
	c.Data["weirongUrl"] = utils.OutPutURL
	c.Data["agentProduct"] = agentProduct
	c.Data["urlList"] = urlList
	c.Data["weirongUrl"] = utils.OutPutURL
	c.TplName = "agent-products/edit_wr_product.html"

}

//保存代理微融产品的数据
func (c *AgentProController) UpdateWrAgent() {
	defer c.ServeJSON()
	//接受参数
	var agentProduct models.AgentProduct
	err := c.ParseForm(&agentProduct)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析参数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "解析参数失败！"}
		return
	}
	//插入新的url
	urlStr := agentProduct.UrlId
	var urlList []string
	idStr := c.GetString("IdList")
	var idList []string
	if strings.Contains(urlStr, ",") {
		url := strings.Split(urlStr, ",")
		for _, u := range url {
			urlList = append(urlList, u)
		}
	} else {
		urlList = append(urlList, agentProduct.UrlId)
	}
	if strings.Contains(idStr, ",") {
		url := strings.Split(idStr, ",")
		for _, u := range url {
			idList = append(idList, u)
		}
	} else {
		idList = append(idList, idStr)
	}
	var updateList []models.ProductAgentUrl //待更新的结构体
	var insertList []string                 //待插入的新的url
	for i := 0; i < len(urlList); i++ {
		if i < len(idList) { //有url，有id，做更新操作
			var productUpdate models.ProductAgentUrl
			id, err := strconv.Atoi(idList[i])
			if err != nil {
				c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "插入数据失败！"}
				return
			}
			productUpdate.Id = id
			productUpdate.Url = urlList[i]
			updateList = append(updateList, productUpdate)
		}
		if i >= len(idList) { //有url，没id，做插入操作
			insertList = append(insertList, urlList[i])
		}
	}
	//更新url
	err = models.UpdateAgentProUrl(updateList)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "批量修改代理微融url异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "插入数据失败！"}
		return
	}
	UrlIdList, err := models.AddAgentProUrl(insertList) //插入操作，得到id集合
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "批量插入代理微融url异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "批量插入代理微融url异常！"}
		return
	}
	//拼接id
	for _, id := range UrlIdList {
		var str string
		str = "," + strconv.Itoa(id)
		idStr += str
	}

	agentProduct.UrlId = idStr
	agentProduct.ProId = 0
	err = models.UpdateProductIdStr(&agentProduct) //更新id字符串
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改代理商代理微融的信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "插入数据失败！"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "操作成功"}
}
