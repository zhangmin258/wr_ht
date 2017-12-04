package controllers

import (
	"strconv"
	"strings"
	"wr_v1/cache"
	"wr_v1/models"
)

/*
产品链接相关接口
*/
type UrlController struct {
	BaseController
}

//添加产品主链接和代理产品链接
//@router /addUrl [post]
func (c *UrlController) AddUrl() {
	defer c.ServeJSON()
	//接受参数
	var url models.ProductUrl
	err := c.ParseForm(&url)
	if url.Url == "" {
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "产品主URL不能为空！"}
		return
	}
	//接收代理链接
	agentUrl := c.GetString("AgentUrl")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析参数异常", err.Error(), c.Ctx.Input)
		//解析参数异常
		c.Data["json"] = map[string]interface{}{"ref": 400, "err": err.Error(), "msg": "UrlController/AddUrl|参数解析异常!"}
		return
	}
	//保存主链接
	err = models.AddMainUrl(&url)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存产品主url信息异常", err.Error(), c.Ctx.Input)
		//保存产品主url信息异常
		c.Data["json"] = map[string]interface{}{"ref": 400, "err": err.Error(), "msg": "UrlController/AddUrl|插入主url数据异常!"}
		return
	}
	//判断是否存在url,并保存代理链接
	if !strings.HasPrefix(agentUrl, ",") {
		urlList := strings.Split(agentUrl, ",")
		err = models.AddAgentUrl(url.ProductId, &urlList)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入代理url数据异常", err.Error(), c.Ctx.Input)
			//保存产品代理url信息异常
			c.Data["json"] = map[string]interface{}{"ref": 400, "err": err.Error(), "msg": "UrlController/AddUrl|插入代理url数据异常!"}
			return
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

//修改主url和代理url
//@router /update [post]
func (c *UrlController) UpdateUrl() {
	defer c.ServeJSON()
	//接受页面参数
	var url models.ProductUrls
	if err := c.ParseForm(&url); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ref": 400, "msg": "UrlController/update|参数解析异常"}
		return
	}
	//更新主链接
	if err := models.UpdateMianUrlById(&url); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "链接更新失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ref": 400, "msg": "链接更新失败"}
		return
	}

	var idList []string                        //id列表
	var urlList []string                       //代理url列表
	var i int                                  //id个数
	var l int                                  //代理url个数
	var AgentUrlList []*models.ProductAgentUrl //存放待更新的代理链接
	var NewUrlList []string                    //存放待新增的代理链接
	//获取代理url的id
	if agentUrlId := c.GetString("agentUrlId"); agentUrlId != "" {
		idList = strings.Split(agentUrlId, ",")
	}
	//获取代理url
	if agentUrl := c.GetString("agentUrl"); agentUrl != "" {
		urlList = strings.Split(agentUrl, ",")
	}
	i = len(idList)
	l = len(urlList)
	//处理已存在的代理链接id和url
	for a := 0; a < i; a++ {
		var productAgentUrl models.ProductAgentUrl
		id, err := strconv.Atoi(idList[a])
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ref": 400, "msg": "UrlController/update|代理链接id转化异常"}
			return
		}
		productAgentUrl.Id = id
		productAgentUrl.Url = urlList[a]
		AgentUrlList = append(AgentUrlList, &productAgentUrl)
	}
	//更新代理链接操作
	if err := models.UpdateAgentUrlList(AgentUrlList); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "代理链接更新失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ref": 400, "msg": "UrlController/update|代理链接更新失败"}
		return
	}

	//当存在新增的代理链接
	if l > i {
		//处理新添加的代理链接
		for b := i; b < l; b++ {
			if urlList[b] != "" {
				NewUrlList = append(NewUrlList, urlList[b])
			}
		}
		//新增代理链接操作
		if err := models.AddAgentUrl(url.ProductId, &NewUrlList); err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "代理链接更新失败", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ref": 400, "msg": "UrlController/update|代理链接更新失败"}
			return
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}
