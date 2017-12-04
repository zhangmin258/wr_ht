package controllers

import (
	"os"
	"strings"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
	"zcm_tools/uuid"
)

/*
产品相关接口
*/
type ProductController struct {
	BaseController
}

//保存产品的图片
//@router /uploadIcon [post]
func (c *ProductController) Post() {
	f, h, err := c.GetFile("proIcon")
	h.Filename = uuid.NewUUID().Hex() + ".png"
	defer func() {
		c.ServeJSON()
		f.Close()
	}()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "解析上传文件异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/uploadIcon|解析上传文件异常!"}
		return
	}
	_, err = os.Stat("static/upload/")
	if os.IsNotExist(err) {
		os.Mkdir("static/upload/", 0)
	}
	err = c.SaveToFile("proIcon", "static/upload/"+h.Filename) // 保存位置在 static/upload, 没有文件夹要先创建
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存上传文件异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/uploadIcon|保存上传文件异常!"}
		return
	}
	filePath := "static/upload/" + h.Filename
	err, url := utils.UploadAliyun(h.Filename, filePath)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "阿里云保存文件异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/uploadIcon|阿里云保存文件异常!"}
		return
	}
	err = os.Remove("static/upload/" + h.Filename)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除本地文件异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/uploadIcon|删除本地文件异常!"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "url": url}
	return

}

//添加产品
//@router /addProduct [post]
func (c *ProductController) AddProduct() {
	defer c.ServeJSON()
	//接受参数
	var product models.ProductForAdd
	err := c.ParseForm(&product)
	if err != nil {
		//解析参数异常
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/addProduc|参数解析异常!"}
		return
	}
	//插入产品数据
	var proId int
	//新增小额H5
	if product.LoanProductType == 0 && product.CooperationType == 1 {
		product.Code = uuid.NewUUID().Hex()
		proId, err = models.AddProductNew(&product)
		if err != nil && proId < 0 {
			//插入商品异常
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存小额H5商品数据异常！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/addProduc|保存小额商品数据异常!"}
			return
		}
	}
	//增加小额API
	if product.LoanProductType == 0 && product.CooperationType == 0 {
		product.Code = uuid.NewUUID().Hex()
		proId, err = models.AddProductNew(&product)
		if err != nil && proId == -1 {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "新增小额API合作产品异常！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/addProduc|插入小额API商品数据异常!"}
			return
		}
	}

	//增加大额API
	if product.LoanProductType == 1 && product.CooperationType == 0 {
		product.Code = uuid.NewUUID().Hex()
		proId, err = models.AddBigAPIProductNew(&product)
		if err != nil && proId < 0 {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存大额API产品数据异常！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/addProduc|保存大额API商品数据异常!"}
			return
		}
	}

	//新增大额H5
	if product.LoanProductType == 1 && product.CooperationType == 1 {
		product.Code = uuid.NewUUID().Hex()
		proId, err = models.AddBigProductNew(&product)
		if err != nil && proId == -1 {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入大额商品数据异常！", err.Error(), c.Ctx.Input)
			//插入商品异常
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/addProduct|插入大额商品数据异常!"}
			return
		}
	}
	if product.CooperationType == 1 {
		//如果添加的为H5产品默认绑定本公司代理商
		var agentId = 0 //代理商Id
		err = models.BindAgent(proId, agentId)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "绑定代理商异常！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/addProduct|绑定代理商异常!"}
			return
		}
	}

	//初始化business数据
	err = models.InsertBusiness(product.ProName)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "初始化business数据异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/addProduct|初始化business数据异常!"}
		return
	}
	//插入关联表数据
	err = models.AddProductJob(proId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入商品关联数据异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/addProduct|插入商品关联数据异常!"}
		return
	}
	v1Log.Println("请求地址：", c.Ctx.Input.URI(), "用户信息：", c.User.Id, "RequestBody：", c.Ctx.Request.Body, "IP：", c.Ctx.Input.IP())
	c.Data["json"] = map[string]interface{}{"ret": 200, "proId": proId}
}

//条件分页查询产品列表
//@router /getProductList [get]
func (c *ProductController) GetProductList() {
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//产品类型
	if loanProducType := c.GetString("loanProducType"); loanProducType != "" {
		condition += " AND p.loan_product_type=?"
		params = append(params, loanProducType)
	}
	//合作类型
	if cooperationType := c.GetString("cooperationType"); cooperationType != "" {
		condition += " AND p.cooperation_type = ?"
		params = append(params, cooperationType)
	}
	//产品名称或者机构名称
	if name := c.GetString("name"); name != "" {
		condition += " AND p.name like ? OR o.name LIKE ? "
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
	}
	//上下线状态
	if isUse := c.GetString("productType"); isUse != "" && isUse != "2" {
		condition += " AND p.is_use = ? "
		params = append(params, isUse)
	}
	//查询
	productList, err := models.GetProductList(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询产品数据异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetProductListCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询所有商品数量异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询需要页数异常！", err.Error(), c.Ctx.Input)
	}
	//查询冻结和使用产品数量
	UsedCount, err := models.GetUsedProductCount()
	NotUsedCount, err := models.GetNotUsedProductCount()
	c.Data["UsedCount"] = UsedCount
	c.Data["NotUsedCount"] = NotUsedCount
	c.Data["productList"] = productList
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.TplName = "agent-products/productList.html"
}

//跳转到添加产品界面:H5
//@router /jumpToAddH5 [get]
func (c *ProductController) JumpToAddH5() {
	c.IsNeedTemplate()
	loanProductType := c.GetString("loan_product_type")
	if loanProductType == "1" {
		c.TplName = "agent-products/addBigH5AgentProduct.html"
		return
	}
	c.TplName = "agent-products/addH5AgentProduct.html"
}

//跳转到添加产品界面:API
//@router /jumpToAddAPI [get]
func (c *ProductController) JumpToAddAPI() {
	c.IsNeedTemplate()
	c.TplName = "agent-products/addAPIAgentproduct.html"
}

//跳转到大额贷款API
//@router /jumptoEditBigAPI [get]
func (c *ProductController) JumptoEditBigAPI() {
	c.IsNeedTemplate()
	c.TplName = "agent-products/addBigAPIAgentproduct.html"
}

//H5
//@router /editH5Product [get]
func (c *ProductController) EditProduct() {
	c.IsNeedTemplate()
	id, err := c.GetInt("id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取id失败", err.Error(), c.Ctx.Input)
	}
	//查询该产品信息
	product, err := models.SearchProductById(id)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查找商品信息异常！", err.Error(), c.Ctx.Input)
	}
	//查询机构信息
	org, err := models.SearchOrgById(product.OrgId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询机构信息异常！", err.Error(), c.Ctx.Input)
	}
	//查询商务联系人
	bus, err := models.SearchBusinessLinkmanById(product.BusId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查询商务联系人异常！", err.Error(), c.Ctx.Input)
	}
	//查询清算信息
	cleaning, err := models.SearchCleaningById(id)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据产品的id查询清算信息数据异常！", err.Error(), c.Ctx.Input)
	}
	//查询主链接信息
	mainUrl, err := models.SearchMainUrlById(id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据产品id查询主链接异常！", err.Error(), c.Ctx.Input)
	}
	var proUrl models.ProductUrls
	for _, v := range mainUrl {
		proUrl.ProductId = v.ProductId
		switch v.UrlType {
		case 1:
			proUrl.UrlId = v.Id
			proUrl.Url = v.Url
		case 4:
			proUrl.BeforeId = v.Id
			proUrl.RegsteUrlBefore = v.Url
		case 5:
			proUrl.AfterId = v.Id
			proUrl.RegsteUrlAfter = v.Url
		}
	}
	//查询代理链接信息
	agentUrlList, err := models.SearchAgentUrlById(id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据产品id查询代理链接异常！", err.Error(), c.Ctx.Input)
	}
	//处理期限
	str := product.LoanTermCount
	if strings.Contains(str, "个月") {
		product.LoanTermCountType = "月"
		product.LoanTermCount = strings.Replace(str, "个月", "", -1)
	}
	if strings.Contains(str, "天") {
		product.LoanTermCountType = "天"
		product.LoanTermCount = strings.Replace(str, "天", "", -1)
	}
	//手动给利率赋值
	if product.FeeType == 1 && product.FeeMethod == 1 {
		product.LoanFee = product.LoanTaxFee
	}
	if product.FeeType == 1 && product.FeeMethod == 2 {
		product.LoanFee = product.LoanDailyFee
	}
	c.Data["product"] = product
	c.Data["org"] = org
	c.Data["bus"] = bus
	c.Data["cleaning"] = cleaning
	c.Data["mainUrl"] = proUrl
	c.Data["agentUrlList"] = agentUrlList
	c.Data["agentUrlCount"] = len(agentUrlList) + 1
	if product.LoanProductType == 1 {
		c.TplName = "agent-products/reviseBigH5AgentProduct.html"
		return
	}
	if product.LoanProductType == 0 {
		c.TplName = "agent-products/reviseH5AgentProduct.html"
		return
	}
	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询产品信息异常,产品类型错误！", err.Error(), c.Ctx.Input)
}

//回显API数据
//@router /editAPIProduct [get]
func (this *ProductController) ReturnShow() {
	this.IsNeedTemplate()
	proid, err := this.GetInt("id")
	if err != nil {
	}
	//查询该产品信息
	product, err := models.SearchProductById(proid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询该产品信息异常！", err.Error(), this.Ctx.Input)
	}
	org, err := models.SearchOrgById(product.OrgId)
	//查询机构信息
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询机构信息异常！", err.Error(), this.Ctx.Input)
	}
	//fmt.Println("--BusId--", product.BusId)

	//查询商务联系人
	bus, err := models.SearchBusinessLinkmanById(product.BusId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询商务联系人异常！", err.Error(), this.Ctx.Input)
	}
	//查询清算信息
	cleaning, err := models.SearchCleaningById(proid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询清算信息异常！", err.Error(), this.Ctx.Input)
	}

	//查询认证信息
	productAuthList, err := models.QueryProAuthBYId(proid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询认证信息异常！", err.Error(), this.Ctx.Input)
	}
	//统计同一产品的认证信息数目
	count, err := models.StaProAuthNumByProid(proid)
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "查询认证信息的数目异常！", err.Error(), this.Ctx.Input)
	}

	//处理期限
	str := product.LoanTermCount
	if strings.Contains(str, "个月") {
		product.LoanTermCountType = "月"
		product.LoanTermCount = strings.Replace(str, "个月", "", -1)
	}
	if strings.Contains(str, "天") {
		product.LoanTermCountType = "天"
		product.LoanTermCount = strings.Replace(str, "天", "", -1)
	}
	//手动给利率赋值
	if product.FeeType == 1 && product.FeeMethod == 1 {
		product.LoanFee = product.LoanTaxFee
	}
	if product.FeeType == 1 && product.FeeMethod == 2 {
		product.LoanFee = product.LoanDailyFee
	}

	this.Data["product"] = product
	this.Data["org"] = org
	this.Data["bus"] = bus
	this.Data["cleaning"] = cleaning
	this.Data["productAuthList"] = productAuthList
	this.Data["count"] = count

	if product.LoanProductType == 1 {
		this.TplName = "agent-products/reviseBigAPIAgentproduct.html"
		return
	}
	if product.LoanProductType == 0 {
		this.TplName = "agent-products/reviseAPIAgentProduct.html"
		return
	}
}

//需求更改
/*//根据id删除产品
//@router /deleteProduct [post]
func (c *ProductController) DeleteProduct() {
	defer c.ServeJSON()
	//接受参数
	pid, err := c.GetInt("productId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/deleteProduct|参数解析异常!"}
		return
	}
	err = models.DeleteProduct(pid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id删除产品信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/deleteProduct|删除产品异常!"}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}*/

//更新产品
//@router /update [post]
func (c *ProductController) UpdateProduct() {
	defer c.ServeJSON()
	//接受参数
	var product models.ProductForAdd
	if err := c.ParseForm(&product); err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/update|参数解析异常!"}
		return
	}
	if product.LoanProductType == 1 {
		if product.CooperationType == 0 {
			//fmt.Println("--product-222-", product)
			if err := models.UpdateBigAPIProduct(&product); err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "更新大额API产品信息异常！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/update|编辑大额API商品数据异常!"}
				return
			}
		}
		if product.CooperationType == 1 {
			if err := models.UpdateBigProduct(&product); err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "编辑大额商品数据异常！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/update|编辑大额商品数据异常!"}
				return
			}
			//c.Data["json"] = map[string]interface{}{"ret": 200}
		}

	} else if product.LoanProductType == 0 {
		if product.CooperationType == 0 {
			if err := models.UpdateProductNew(&product); err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "编辑小额API合作产品异常！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/update|编辑小额API合作产品异常!"}
				return
			}
		}
		if product.CooperationType == 1 {
			if err := models.UpdateProductNew(&product); err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "编辑小额H5合作产品异常！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "product/update|编辑小额H5合作产品异常!"}
				return
			}
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

//保存产品认证
//@router /saveProductAuth [post]
func (this *ProductController) SaveProductAuth() {
	defer this.ServeJSON()
	ProductId, err := this.GetInt("ProductId")
	AuthFieldList := this.GetString("AuthFieldList")
	AuthFieldArr := strings.Split(AuthFieldList, ",")

	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "参数解析异常！", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "/product/saveProductAuth|参数解析异常！"}
		return
	}
	if AuthFieldList != "" {
		err = models.SaveProductAuth(ProductId, AuthFieldArr)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存产品认证异常！", err.Error(), this.Ctx.Input)
			//保存产品认证信息异常
			this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "/product/saveProductAuth|保存产品认证信息异常!"}
			return
		}
		//保存日志记录
		v1Log.Println("请求地址：", this.Ctx.Input.URL(), "用户信息：", this.User.Id, "RequestBody：", this.Ctx.Request.Body, "IP：", this.Ctx.Input.IP())
	}
	this.Data["json"] = map[string]interface{}{"ret": 200}
}

//修改产品认证
//@router /updateProductAuth [post]
func (this *ProductController) UpdateProductAuth() {
	defer this.ServeJSON()
	ProductId, err := this.GetInt("ProductId")
	AuthFieldList := this.GetString("AuthFieldList")
	AuthFieldArr := strings.Split(AuthFieldList, ",")
	if err != nil {
		cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "参数解析异常！", err.Error(), this.Ctx.Input)
		this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "/product/updateProductAuth|参数解析异常！"}
		return
	}
	//删除proid匹配的所有对象
	err = models.DeleteProductAuthByProId(ProductId)
	//新增
	if AuthFieldList != "" {
		err = models.SaveProductAuth(ProductId, AuthFieldArr)
		if err != nil {
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存大额API产品认证信息异常", err.Error(), this.Ctx.Input)
			//保存产品认证信息异常
			this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "/product/saveProductAuth|保存大额API产品认证信息异常!"}
			return
		}
	}
	this.Data["json"] = map[string]interface{}{"ret": 200}
}

//保存小额API合作产品——认证信息
//@router /saveAPIProductAuth [post]
func (this *ProductController) SaveAPIProductAuth() {
	defer this.ServeJSON()
	var proAuth models.ProductAuth
	ProductId, err := this.GetInt("ProductId")
	AuthFieldList := this.GetString("AuthFieldList")
	AuthFieldArr := strings.Split(AuthFieldList, ",")

	if err := this.ParseForm(&proAuth); err != nil {
		this.Data["json"] = map[string]interface{}{"ret": 500, "err": err.Error(), "msg": "/product/saveAPIProductAuth|参数解析异常！"}
		return
	}
	if AuthFieldList != "" {
		err = models.SaveProductAuth(ProductId, AuthFieldArr)
		if err != nil {
			//保存失败，将失败信息存入日志
			cache.RecordLogs(this.User.Id, 0, this.User.Name, this.User.DisplayName, "", "保存小额API产品认证信息异常", err.Error(), this.Ctx.Input)
			//
			this.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "/product/saveAPIProductAuth|保存小额API产品认证信息异常!"}
			return
		}
	}
	this.Data["json"] = map[string]interface{}{"ret": 200}
}
