package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"net/url"
	"wr_v1/models"
	"wr_v1/utils"
	"zcm_tools/log"
)

var v1Log *log.Log

func init() {
	v1Log = log.Init()
}

// BaseController 基础controller
type BaseController struct {
	beego.Controller
	User *models.SysUser
}

// Prepare 验证用户登录信息
func (c *BaseController) Prepare() {
	verify := false
	uid := c.Ctx.GetCookie("uid")
	pid := c.Ctx.GetCookie("pid")
	//判断
	if uid != "" && pid != "" {
		if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyUserPrefix+uid) {
			if data, err := utils.Rc.RedisBytes(utils.CacheKeyUserPrefix + uid); err == nil {
				err = json.Unmarshal(data, &c.User)
				if c.User != nil && utils.MD5(c.User.Password+utils.PasswordEncryptKey) == pid {
					// if c.Ctx.Input.URL() != "/user/modifypassword" && c.Ctx.Input.URL() != "/user/postmodifypassword" && c.User.Password == utils.MD5("123456") {
					// 	c.Ctx.Redirect(302, "/user/modifypassword")
					// 	c.StopRun()
					// 	return
					// }
					// beego.Notice(string(data))
					ip := c.Ctx.Input.IP()
					requestBody, _ := url.QueryUnescape(string(c.Ctx.Input.RequestBody))
					v1Log.Println("请求地址：", c.Ctx.Input.URI(), "用户信息：", string(data), "RequestBody：", requestBody, "IP：", ip)
					//重新保存用户状态
					if ip == "127.0.0.1" || ip == "60.191.125.34" || ip == "60.191.37.251" {
						utils.Rc.Put(utils.CacheKeyUserPrefix+c.User.Name, data, utils.RedisCacheTime_TwoHour)
					} else {
						utils.Rc.Put(utils.CacheKeyUserPrefix+c.User.Name, data, utils.RedisCacheTime_User)
					}

					verify = true
				}
			}
		}
	} else if c.Ctx.Input.IsUpload() { //上传文件跳过验证
		verify = true
	}
	if !verify {
		if c.Ctx.Input.IsAjax() {
			c.Ctx.Output.JSON(map[string]interface{}{"ret": 408, "msg": "timeout"}, false, false)
			c.StopRun()
		} else {
			c.Ctx.Redirect(302, "/login")
			c.StopRun()
		}
	}
}

// func (c *BaseController),Finish() {
// 	beego.Info(string(c.Ctx.ResponseWriter.Data))
// }

//是否需要模板
func (c *BaseController) IsNeedTemplate() {
	pushstate := c.GetString("pushstate")
	if pushstate != "1" {
		c.Data["DisplayName"] = c.User.DisplayName
		c.Layout = "layout/layout.html"
	}
}
