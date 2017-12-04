package controllers

import (
	"os"
	"sort"
	"strings"
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
	"zcm_tools/uuid"

	"github.com/astaxie/beego"
)

type YaowangProductController struct {
	BaseController
}

//条件分页查询产品列表
//@router /getyaowangproductlist [get]
func (c *YaowangProductController) GetYaowangProductlist() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "yw-cooperation/yw_product_management.html"
	}()
	//读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//产品名称
	if name := c.GetString("name"); name != "" {
		condition += " AND name like ? "
		params = append(params, "%"+name+"%")
	}
	//产品来源
	if source := c.GetString("Source"); source != "" {
		condition += " AND source like ? "
		params = append(params, "%"+source+"%")
	}
	//查询
	productlist, err := models.GetYaowangProductRecommend(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20) //产品名称列表
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询产品数据异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetYaowangProductCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询所有商品数量异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询需要页数异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["productList"] = productlist
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
}

//添加产品
//@router /addproduct [post]
func (c *YaowangProductController) AddProduct() {
	defer c.ServeJSON()
	//接受参数
	var product models.ReturnYaowangProduct
	err := c.ParseForm(&product)
	if err != nil {
		//解析参数异常
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "yaowangproduct/addproduct|参数解析异常!"}
		return
	}
	switch product.IsUsed {
	case 1:
		product.IsUsed = 0
	case 2:
		product.IsUsed = 1
	}
	product.CreateTime = time.Now()
	//插入产品数据
	var proId int
	proId, err = models.AddYaowangProduct(product)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入产品数据出错", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "yaowangproduct/addproduct|保存出错!"}

	}
	v1Log.Println("请求地址：", c.Ctx.Input.URI(), "用户信息：", c.User.Id, "RequestBody：", c.Ctx.Request.Body, "IP：", c.Ctx.Input.IP())
	c.Data["json"] = map[string]interface{}{"ret": 200, "proId": proId}
}

//更新产品
//@router /updateproduct [post]
func (c *YaowangProductController) UpdateProduct() {
	defer c.ServeJSON()
	//接受参数
	var product models.ReturnYaowangProduct
	err := c.ParseForm(&product)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "yaowangproduct/updateproduct|参数解析异常!"}
		return
	}
	err = models.UpdateYaowangProduct(product)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "更新大额API产品信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "yaowangproduct/updateproduct|修改产品失败!"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

//添加产品页面
//@router /addproductpage [get]
func (c *YaowangProductController) AddProductPage() {
	c.IsNeedTemplate()
	c.TplName = "yw-cooperation/add_yw_product.html"
}

//编辑页面
//@router /editproductmange [get]
func (c *YaowangProductController) EditProductMange() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "yw-cooperation/edit_yw_product.html"
	}()
	var product models.YaowangProduct
	id, _ := c.GetInt("id")
	product, err := models.GetYaowangProductById(id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查找产品异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["product"] = product
}

//保存产品的图片
//@router /uploadicon [post]
func (c *YaowangProductController) UpLoadIcon() {
	f, h, err := c.GetFile("proIcon")
	h.Filename = uuid.NewUUID().Hex() + ".png"
	defer func() {
		c.ServeJSON()
		f.Close()
	}()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析上传文件异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "yaowangproduct/uploadIcon|解析上传文件异常!"}
		return
	}
	_, err = os.Stat("static/upload/")
	if os.IsNotExist(err) {
		os.Mkdir("static/upload/", 0)
	}
	err = c.SaveToFile("proIcon", "static/upload/"+h.Filename) // 保存位置在 static/upload, 没有文件夹要先创建
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存上传文件异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "yaowangproduct/uploadIcon|保存上传文件异常!"}
		return
	}
	filePath := "static/upload/" + h.Filename
	err, url := utils.UploadAliyun(h.Filename, filePath)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "阿里云保存文件异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "yaowangproduct/uploadIcon|阿里云保存文件异常!"}
		return
	}
	err = os.Remove("static/upload/" + h.Filename)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除本地文件异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "yaowangproduct/uploadIcon|删除本地文件异常!"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "url": url}
	return
}

//模糊查询遥望渠道列表
////@router /getsourcelist [get]
func (c *YaowangProductController) GetSourceList() {
	defer c.ServeJSON()
	name := c.GetString("source")
	if name == "所有渠道" {
		name = ""
	}
	condition := ""
	var param []string
	if name != "" {
		condition = " AND source LIKE ? "
		param = append(param, "%"+name+"%")
	} else {
		condition = "AND source != ''"
	}
	//模糊查询渠道列表
	outPutSourceList, err := models.GetYaowangSourceList(condition, param)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": "获取遥望渠道列表失败!"}
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取遥望渠道列表异常", err.Error(), c.Ctx.Input)
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "sourceList": outPutSourceList}
	return
}

//数据明细
//@router /yaowangdatadetaillist [get]
func (c *YaowangProductController) YaowangDataDetailList() {
	c.IsNeedTemplate()
	defer func() {
		c.TplName = "yw-cooperation/yw_dataDetail.html"
	}()
	//点击—注册统计
	startDate := time.Now().AddDate(0, 0, -6).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")
	c.Data["startDate"] = startDate
	c.Data["endDate"] = endDate

	return
}

//获取遥望点击-注册统计
//@router /getyaowangdataanalysislist [get]
func (c *YaowangProductController) GetYaowangDataAnalysisList() {
	defer c.ServeJSON()
	//获取渠道
	source := c.GetString("source")
	if source == "所有渠道" {
		source = ""
	}
	//点击—注册统计
	startDate := c.GetString("startTime")
	endDate := c.GetString("endTime")
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -6).Format("2006-01-02")
		endDate = time.Now().Format("2006-01-02")
	}
	condition := ""
	condition1 := ""
	var params []string
	if startDate != "" {
		condition += " AND DATE(cr.create_time) >= ? "
		condition1 += " AND DATE(create_time) >= ? "
		params = append(params, startDate)
	}
	if endDate != "" {
		condition += " AND DATE(cr.create_time) <= ? "
		condition1 += " AND DATE(create_time) <= ? "
		params = append(params, endDate)
	}
	if source != "" {
		condition += " AND cr.source = ?"
		condition1 += " AND source = ?"
		params = append(params, source)
	}
	//加载数据
	loadInfos, err := models.GetPageLoadCount(condition1, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取遥望页面加载记录失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取遥望页面加载记录失败"}
		return
	}
	//产品点击数据
	clickInfos, err := models.GetYaowangProductUV(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取遥望点击记录失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取社保公积金信息失败"}
		return
	}
	var times []string
	timeMap := make(map[string]models.YaowangDataAnalysis)
	var dataAnalysisList []models.YaowangDataAnalysis
	for _, v := range loadInfos {
		var dataAnalysis models.YaowangDataAnalysis
		dataAnalysis.CreateTime = v.CreateTime.Format("2006-01-02")
		dataAnalysis.LoadCount = v.LoadCount
		timeMap[dataAnalysis.CreateTime] = dataAnalysis
		times = append(times, v.CreateTime.Format("2006-01-02"))
	}
	for _, v := range clickInfos {
		if m, ok := timeMap[v.CreateTime.Format("2006-01-02")]; ok {
			var productinfo models.YaowangProductCount
			productinfo.ProductName = v.ProductName
			productinfo.ProductUV = v.Count
			m.ProductInfo = append(m.ProductInfo, productinfo)
			timeMap[v.CreateTime.Format("2006-01-02")] = m
		}
	}
	for _, k := range times {
		v := timeMap[k]
		v.Length = len(v.ProductInfo) + 1
		dataAnalysisList = append(dataAnalysisList, v)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "data": dataAnalysisList}
	return
}

//获取遥望数据明细条目
//@router /getyaowangdatadetaillist [get]
func (c *YaowangProductController) GetYaowangDataDetailList() {
	defer c.ServeJSON()
	//获取渠道
	source := c.GetString("source")
	flag := false //用于判断是否有该渠道
	if source == "所有渠道" {
		source = ""
		flag = true
	}
	//获取所有渠道
	sourceList, err := models.GetYaowangSource()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取所有遥望渠道异常!", err.Error(), c.Ctx.Input)
	}
	if source != "" {
		for _, v := range sourceList {
			if source == v {
				flag = true
			}
		}
	} else {
		flag = true
	}
	if !flag { //没有该渠道,直接返回
		c.Data["json"] = map[string]interface{}{"ret": 403, "flag": flag, "err": "没有该渠道！"}
		return
	}
	sourceStr := strings.Join(sourceList, ",")
	sourceCondition := ""
	sourceParam := []string{}
	params := []string{}
	if source != "" {
		sourceCondition = " AND out_put_source = ? "
		sourceParam = append(sourceParam, source)
		params = append(params, source)
	} else {
		sourceCondition = " AND FIND_IN_SET(out_put_source,?) "
		sourceParam = append(sourceParam, sourceStr)
	}
	//数据明细
	pageNum, _ := c.GetInt("page", 1) //分页信息（第几页）
	if pageNum < 1 {
		pageNum = 1
	}
	var yaowangCountDataList []models.YaowangCountData
	//注册用户
	registInfo, err := models.GetYWRegistCount(sourceCondition, sourceParam)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取注册用户失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取注册用户失败"}
		return
	}
	//登录用户
	loginInfo, err := models.GetYWLoginCount(sourceCondition, sourceParam)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取登录用户失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取登录用户失败"}
		return
	}
	//总uv
	clickInfo, err := models.GetYWTotalClickCount(params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取遥望总点击数据失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取遥望总点击数据失败"}
		return
	}
	var yaowangCountData models.YaowangCountData
	timeNow := time.Now().Format("2006-01-02")
	timeStrs := []string{timeNow}
	yaowangCountData.CreateDate = timeNow
	dataMap := make(map[string]models.YaowangCountData)
	dataMap[timeNow] = yaowangCountData
	for _, v := range registInfo {
		if r, ok := dataMap[v.CreateDate.Format("2006-01-02")]; ok {
			r.RegistCount = v.Count
			r.ActiveCount = v.Count
			r.CreateDate = v.CreateDate.Format("2006-01-02")
			dataMap[v.CreateDate.Format("2006-01-02")] = r
		} else {
			timeStrs = append(timeStrs, v.CreateDate.Format("2006-01-02"))
			var m models.YaowangCountData
			m.CreateDate = v.CreateDate.Format("2006-01-02")
			m.RegistCount = v.Count
			m.ActiveCount = v.Count
			dataMap[v.CreateDate.Format("2006-01-02")] = m
		}
	}
	sort.Strings(timeStrs)
	//登录统计明细
	for _, v := range loginInfo {
		if r, ok := dataMap[v.CreateDate.Format("2006-01-02")]; ok {
			r.LoginCount = v.Count
			r.ActiveCount = v.Count + r.RegistCount
			r.CreateDate = v.CreateDate.Format("2006-01-02")
			dataMap[v.CreateDate.Format("2006-01-02")] = r
		} else {
			timeStrs = append(timeStrs, v.CreateDate.Format("2006-01-02"))
			var m models.YaowangCountData
			m.CreateDate = v.CreateDate.Format("2006-01-02")
			m.LoginCount = v.Count
			m.ActiveCount = v.Count
			dataMap[v.CreateDate.Format("2006-01-02")] = m
		}
	}
	sort.Strings(timeStrs)
	//uv明细
	for _, v := range clickInfo {
		if r, ok := dataMap[v.CreateDate.Format("2006-01-02")]; ok {
			r.UVCount = v.Count
			r.CreateDate = v.CreateDate.Format("2006-01-02")
			dataMap[v.CreateDate.Format("2006-01-02")] = r
		} else {
			timeStrs = append(timeStrs, v.CreateDate.Format("2006-01-02"))
			var m models.YaowangCountData
			m.CreateDate = v.CreateDate.Format("2006-01-02")
			m.UVCount = v.Count
			dataMap[v.CreateDate.Format("2006-01-02")] = m
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(timeStrs)))
	for _, v := range timeStrs {
		if r, ok := dataMap[v]; ok {
			yaowangCountDataList = append(yaowangCountDataList, r)
		}
	}
	count := len(yaowangCountDataList)
	//页数
	pageCount, err := utils.GetPageCount(count, utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取总页数失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 403, "err": "获取总页数失败"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "datalist": yaowangCountDataList, "pageCount": pageCount, "pageNum": pageNum, "pageSize": utils.PageSize20, "count": count, "flag": flag}
	return
}

//遥望数据修复
func FixYWProductName() {
	//获取遥望产品id和名称
	ywProInfos, err := models.GetYWProductInfos()
	if err != nil {
		beego.Info(err)
		return
	}
	//修复遥望点击记录中的产品名称
	err = models.FixYWClickProductName(ywProInfos)
	if err != nil {
		beego.Info(err)
	}
}
