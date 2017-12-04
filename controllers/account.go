package controllers

import (
	"encoding/json"
	"fmt"
	"wr_v1/models"
	"wr_v1/utils"

	"github.com/astaxie/beego"
)

// AccountController 登录控制器
type AccountController struct {
	beego.Controller
}

// Login 登录
func (c *AccountController) Login() {
	c.TplName = "login.html"
}

// CheckPassword 验证密码
func (c *AccountController) CheckPassword() {
	defer c.ServeJSON()

	name := c.GetString("username")
	password := c.GetString("password")
	password = utils.MD5(password)
	verify_code := c.GetString("verify_code") // 将军令
	m, _ := models.Login(name, password)
	if m != nil {
		// TODO 记得恢复判断
		ip := c.Ctx.Input.IP()
		fmt.Println(ip, verify_code)
		// && hostName != "60.191.37.251"
		/*if ip != "127.0.0.1" && ip != "60.191.125.34" && ip != "60.191.37.251" {
			if verify_code == "" {
				c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "请输入验证码"}
				return
			}
			result, _ := utils.Authenticate(m.Id, verify_code)
			if !result {
				c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "登录失败3:验证码错误"}
				return
			}
		}*/
		if data, err2 := json.Marshal(m); err2 == nil && utils.Re == nil {
			//添加日志记录
			utils.Rc.Put(utils.CacheKeyUserPrefix+m.Name, data, utils.RedisCacheTime_User)
			//保存用户缓存和cookie
			c.Ctx.SetCookie("uid", m.Name)
			password2 := utils.MD5(password + utils.PasswordEncryptKey)
			c.Ctx.SetCookie("pid", password2)
			// c.Ctx.SetCookie("userId", strconv.Itoa(m.Id))
			c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "登录成功"}
		} else {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "登录失败！用户信息解析失败或Redis服务器异常！"}
		}

	} else {
		c.Data["json"] = map[string]interface{}{"ret": 404, "msg": "登录失败！用户名或密码不正确"}
	}
}

// LoginOut 退出登录
func (c *AccountController) LoginOut() {
	uid := c.Ctx.GetCookie("uid")
	if uid != "" {
		//清除cookie
		c.Ctx.SetCookie("uid", "0", -1)
		c.Ctx.SetCookie("pid", "0", -1)
		//清除redis
		utils.Rc.Delete(utils.CacheKeyUserPrefix + uid)
	}
	c.Ctx.Redirect(302, "/login")
}
