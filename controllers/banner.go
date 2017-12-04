package controllers

import (
	"strconv"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

type BannerController struct {
	BaseController
}

//Banner管理页面
func (c *BannerController) ShowBannerList() {
	c.IsNeedTemplate()
	iType, _ := c.GetInt("imgType")
	condition := ""
	params := []int{}
	if iType != 0 && iType != 1 { //1:首页banner 2:贷款页banner 3:广告位 4:教程攻略 (1.4版本没有首页banner)
		condition += " AND itype = ? "
		params = append(params, iType)
	}
	imageList, err := models.GetAllBanner(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询缩略图信息异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["imageList"] = imageList
	c.TplName = "show-management/show_information_management.html"
}

//跳转到新增banner页面
func (c *BannerController) JumpToAddBanner() {
	c.IsNeedTemplate()
	c.TplName = "show-management/add_homePage_banner.html"
}

//保存Banner信息
func (c *BannerController) SaveBanner() {
	defer c.ServeJSON()
	var image models.Images
	err := c.ParseForm(&image)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ref": 400, "err": err.Error(), "msg": "control/saveHomePageBanner参数解析异常!"}
		return
	}
	if image.State == 1 { //指向产品
		proName := c.GetString("ProductName")
		if proName == "" {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": "产品名称解析异常!"}
			return
		}
		_, err := models.GetProductIdByName(proName)
		if err != nil && err.Error() != utils.ErrNoRow() {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
			return
		}
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": "没有该产品,请输入正确的产品名称!"}
			return
		}
	}

	//如果为广告位，判断是否存在，如果存在，不更新
	if image.Itype == 3 && image.IsUsed == 1 {
		count, err := models.QueryBannerType()
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询广告位异常！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
			return
		}
		if count > 0 {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": "已经存在广告位！"}
			return
		}
	}
	err = models.SaveBanner(&image)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入banner图异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ref": 400, "err": err.Error(), "msg": "control/saveHomePageBanner插入banner图异常!"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

//删除Banner
func (c *BannerController) DelBanner() {
	defer c.ServeJSON()
	imageId, err := c.GetInt("imageId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "control/deleteBanner|参数解析异常!"}
		return
	}
	err = models.DeleteBanner(imageId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除banner图异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "control/deleteBanner|删除产品异常!"}
		return
	}
	//获取type
	itype, err := models.GetTypeById()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取type异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "control/deleteBanner|获取type异常!"}
		return
	}
	//清除缓存
	for _, v := range itype {
		if utils.Re == nil && utils.Rc.IsExist(utils.CACHE_KEY_CAROUSE_IMGS+strconv.Itoa(v)) {
			utils.Rc.Delete(utils.CACHE_KEY_CAROUSE_IMGS + strconv.Itoa(v))
		}
	}

	c.Data["json"] = map[string]interface{}{"ret": 200}
}

//查询banner并跳转到修改页面
func (c *BannerController) JumpToUpdateBanner() {
	c.IsNeedTemplate()
	id, err := c.GetInt("Id")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除产品异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "control/deleteBanner|删除产品异常!"}
	}
	image, err := models.GetBannerById(id)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取数据分析-产品名称异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["image"] = image
	c.TplName = "show-management/edit_banner.html"
}

//修改banner
func (c *BannerController) UpdateBanner() {
	defer c.ServeJSON()
	var image models.Images
	err := c.ParseForm(&image)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
		return
	}

	if image.State == 1 { //指向产品
		proName := c.GetString("ProductName")
		if proName == "" {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": "产品名称解析异常!"}
			return
		}
		_, err := models.GetProductIdByName(proName)
		if err != nil && err.Error() != utils.ErrNoRow() {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
			return
		}
		if err != nil {
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": "没有该产品,请输入正确的产品名称!"}
			return
		}
	}
	//如果为广告位，判断是否存在，如果存在，不更新
	if image.Itype == 3 && image.IsUsed == 1 {
		img, err := models.GetBannerById(image.Id)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取banner信息失败", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
			return
		}
		if img.Itype != 3 || img.IsUsed != 1 {
			count, err := models.QueryBannerType()
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询广告位异常！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
				return
			}
			if count > 0 {
				c.Data["json"] = map[string]interface{}{"ret": 400, "err": "已经存在广告位！"}
				return
			}
		}
	}
	err = models.UpdateBanner(&image)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "更新banner异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error()}
		return
	}
	//获取type
	itype, err := models.GetTypeById()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取type异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取type异常!"}
		return
	}
	//清除缓存
	for _, v := range itype {
		if utils.Re == nil && utils.Rc.IsExist(utils.CACHE_KEY_CAROUSE_IMGS+strconv.Itoa(v)) {
			utils.Rc.Delete(utils.CACHE_KEY_CAROUSE_IMGS + strconv.Itoa(v))
		}
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
	return
}
