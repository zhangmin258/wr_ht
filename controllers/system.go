package controllers

import (
	"strings"

	"strconv"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/services"
	"wr_v1/utils"
)

type SystemController struct {
	BaseController
}

func (c *SystemController) UserList() {
	c.IsNeedTemplate()
	// fmt.Println()
	uidName := c.Ctx.GetCookie("uid")
	if uidName == "admin" {
		c.Data["isAdmin"] = true
	}
	page, _ := c.GetInt("page")
	if page < 1 {
		page = 1
	}

	condition := ""
	pars := []string{}
	if account := c.GetString("account"); account != "" {
		condition += " and su.name=?"
		pars = append(pars, account)
	}
	if username := c.GetString("username"); username != "" {
		condition += " and su.displayname=?"
		pars = append(pars, username)
	}
	if role := c.GetString("role"); role != "" {
		condition += " and sr.displayname=?"
		pars = append(pars, role)
	}
	if status := c.GetString("accountstate"); status != "" {
		condition += " and su.accountstatus=?"
		pars = append(pars, status)
	}
	list, err := models.SysUserList(condition, pars, utils.StartIndex(page, utils.PageSize20), utils.PageSize20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取系统用户列表失败", err.Error(), c.Ctx.Input)
	}
	count := models.SysUserCount(condition, pars)
	pagecount := utils.PageCount(count, utils.PageSize20)

	c.Data["currpage"] = page
	c.Data["pagecount"] = pagecount

	// 二维码
	if len(list) > 0 {
		for i := 0; i < len(list); i++ {
			list[i].Secret = utils.CreateXjdSecret(list[i].Id)
			list[i].AuthURL = utils.CreateXjdAuthURLEscape(list[i].Secret, list[i].Displayname)
		}
	}

	c.Data["list"] = list
	c.Data["count"] = count
	c.TplName = "system-management/system_account_list.html"
}

//@router /user [get]
func (c *SystemController) GetUser() {
	c.IsNeedTemplate()
	uid_str := c.GetString("uid")
	if uid_str == "" {
		c.Data["user"] = &models.SysUserMini{}
	} else {
		uid, err := c.GetInt("uid")
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "id错误", err.Error(), c.Ctx.Input)
			c.Ctx.WriteString("id错误")
			return
		}
		user, err := models.SysUserDetail(uid)
		if err != nil && err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取系统用户详情异常", err.Error(), c.Ctx.Input)
		}
		c.Data["user"] = user
	}
	list, _ := models.SysRoleList()
	c.Data["list"] = list
	c.TplName = "system-management/system_account_form.html"
}

//@router /user [post]
func (c *SystemController) UserAdd() {
	defer c.ServeJSON()
	uid, _ := c.GetInt("uid")
	var user *models.SysUserMini
	var err error
	if uid == 0 {
		user = &models.SysUserMini{}
		user.Name = c.GetString("account")
	} else {
		user, err = models.SysUserDetail(uid)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取系统用户详情异常", err.Error(), c.Ctx.Input)
			}
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
			return
		}
	}
	if c.GetString("password") != "" {
		user.Password = utils.MD5(c.GetString("password"))
	}

	user.Displayname = c.GetString("name")
	user.Email = c.GetString("email")
	user.Role_id, _ = c.GetInt("roleId")
	user.Accountstatus = c.GetString("account_status")
	user.Station_id, _ = c.GetInt("stationId")
	if uid == 0 {
		err = user.Insert()
	} else {
		err = user.Update()
	}
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "添加系统用户失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *SystemController) DelUser() {
	defer c.ServeJSON()
	uid, err := c.GetInt("uid")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取id失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	} else if uid == 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id错误"}
		return
	}
	if uid == 1 { // 系统管理员不给删
		c.Data["json"] = map[string]interface{}{"ret": 200}
		return
	}
	err = models.DeleteUser(uid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除系统用户失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

//修改密码
func (c *SystemController) ModifyPassword() {
	defer c.ServeJSON()
	originalPw := c.GetString("orgpwd")
	newPw := c.GetString("newpwd")
	confirmPw := c.GetString("newpwd2")

	originalPw = utils.MD5(originalPw)
	if originalPw != c.User.Password {
		c.Data["json"] = map[string]interface{}{"ret": 304, "err": "原始密码输入错误，请重新输入！"}
		return
	}
	if len(newPw) < 6 { //密码要求6位数以上
		c.Data["json"] = map[string]interface{}{"ret": 304, "err": "为了账户安全，请输入6位数及以上的新密码！"}
		return
	}

	if newPw != confirmPw {
		c.Data["json"] = map[string]interface{}{"ret": 304, "err": "新密码与确认密码不一致！"}
		return
	}

	flag, err := models.UpdatePassword(c.User.Id, utils.MD5(newPw), newPw)
	if flag && err == nil {
		c.Data["json"] = map[string]interface{}{"ret": 200}
		return
	}
	cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "密码修改失败！", err.Error(), c.Ctx.Input)
	c.Data["json"] = map[string]interface{}{"ret": 304, "err": "密码修改失败！"}
}

func (c *SystemController) RoleList() {
	c.IsNeedTemplate()
	rolelist, err := models.SysRoleList()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取系统角色失败！", err.Error(), c.Ctx.Input)
	}
	c.Data["list"] = rolelist
	c.TplName = "system-management/system_role_list.html"
}

func (c *SystemController) RoleEdit() {
	c.IsNeedTemplate()
	rid_str := c.GetString("rid")
	var role *models.SysRole
	if rid_str == "" {
		role = &models.SysRole{}
	} else {
		rid, err := c.GetInt("rid")
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "id错误！", err.Error(), c.Ctx.Input)
			c.Ctx.WriteString("id错误")
			return
		}
		role, err = models.SysRoleByRid(rid)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id获取系统角色失败！", err.Error(), c.Ctx.Input)
			}
			c.Ctx.WriteString(err.Error())
			return
		}
	}
	c.Data["role"] = role
	c.TplName = "system-management/system_role_edit.html"
}

func (c *SystemController) RoleAdd() {
	defer c.ServeJSON()
	var role *models.SysRole
	var err error
	rid, _ := c.GetInt("rid")
	if rid == 0 {
		role = &models.SysRole{}
	} else {
		role, err = models.SysRoleByRid(rid)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id获取系统角色失败！", err.Error(), c.Ctx.Input)
			}
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
			return
		}
	}
	ids_str := c.GetString("checkId") // 菜单权限
	ids := strings.Split(ids_str, ",")
	role.Displayname = c.GetString("account")
	if rid == 0 {
		err = role.Insert(ids)
	} else {
		err = role.Update(ids)
	}
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "新增或修改角色失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
}

func (c *SystemController) MenuData() {
	defer c.ServeJSON()
	rid_str := c.GetString("role_id")
	var list []models.SysMenu
	var err error
	if rid_str == "all" { // 所有菜单
		list, err = models.GetSysMenuTreeAll()
	} else {
		rid, err := c.GetInt("rid")
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "id错误", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id错误"}
			return
		} // 该角色有的菜单
		list, err = models.GetSysMenuTreeByRoleId(rid)
	}
	if err != nil && err.Error() != "<QuerySeter> no row found" {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取目录失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	m := services.GetSysMenuZTree(list)
	c.Data["json"] = map[string]interface{}{"ret": 200, "m": m}
}

//删除角色
func (c *SystemController) DelRole() {
	defer c.ServeJSON()
	rid, err := c.GetInt("rid")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "id错误", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id错误"}
		return
	}
	err = models.DelRole(rid)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除角色失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

func (c *SystemController) Organization() {
	c.IsNeedTemplate()
	c.TplName = "system-management/system_organization_list.html"
}

//获取组织架构信息
func (c *SystemController) GetOrganizationList() {
	o, _ := services.GetSysOrganizationZTree()
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationList": o}
	c.ServeJSON()
}

//获取组织架构 菜单 数据
func (c *SystemController) BaseData() {
	o, _ := services.GetSysOrganizationZTree()
	m, _ := services.GetAllSysMenuZTree()
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationData": o, "menuData": m}
	c.ServeJSON()
}

//添加组织架构
func (c *SystemController) AddOrganization() {
	defer c.ServeJSON()
	pid, _ := c.GetInt("pid")
	var err error
	var organization *models.SysOrganization
	if pid <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "上级组织选择错误!"}
		return
	} else {
		organization, err = models.GetSysOrganizationById(pid)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "上级组织获取异常", err.Error(), c.Ctx.Input)
			}
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "上级组织不存在!"}
			return
		}
	}
	name := c.GetString("organizationName")
	if name == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织名称不能为空!"}
		return
	}
	organization, err = models.GetSysOrganizationByName(name)
	if organization != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织机构已经存在!"}
		return
	}
	organization = &models.SysOrganization{}
	organization.ParentId = pid
	organization.Name = name
	err = organization.Insert()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "添加失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "添加成功!"}
}

//编辑组织架构
func (c *SystemController) EditOrganization() {
	defer c.ServeJSON()
	var err error
	name := c.GetString("name")
	id, err := strconv.Atoi(c.GetString("id"))
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	var organization *models.SysOrganization
	if id <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织选择错误!"}
		return
	} else {
		organization, err = models.GetSysOrganizationById(id)
		if err != nil {
			if err.Error() != utils.ErrNoRow() {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "组织机构获取异常", err.Error(), c.Ctx.Input)
			}
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织机构不存在!"}
			return
		}
	}
	if name == "" {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织名称不能为空!"}
		return
	}
	organization, err = models.GetSysOrganizationByName(name)
	if organization != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "组织机构已经存在!"}
		return
	}
	organization = &models.SysOrganization{}
	organization.Id = id
	organization.Name = name
	err = models.EditOrganization(name, id)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "编辑失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
		return
	}
	utils.Rc.Delete(utils.CacheKeySystemOrganization)
	c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "编辑成功!"}
}

//添加岗位
func (c *SystemController) AddStation() {
	defer c.ServeJSON()
	var err error
	var station *models.SysStation
	station = &models.SysStation{}
	stationName := c.GetString("stationName")
	orgId, _ := c.GetInt("orgId")
	if orgId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "岗位ID错误!"}
		return
	}
	station, _ = models.GetSysStationByName(stationName, orgId)
	if station != nil {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "岗位已经存在!"}
		return
	}
	roleId, _ := c.GetInt("roleId")
	if roleId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "角色ID错误!"}
		return
	} else {
		role, _ := models.SysRoleByRid(roleId)
		if role == nil {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "角色不存在!"}
			return
		}
	}
	typeStr := c.GetString("typeStr")
	typeArr := strings.Split(typeStr, ",")
	dataStr := c.GetString("dataStr")
	dataArr := strings.Split(dataStr, ",")
	if station == nil {
		station = &models.SysStation{}
	}

	station.Name = stationName
	station.RoleId = roleId
	station.OrgId = orgId
	err = station.Insert(typeArr, dataArr)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "新增失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "新增失败!"}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "新增成功!"}
	}
}

//编辑岗位
func (c *SystemController) UpdateStation() {
	defer c.ServeJSON()
	var err error
	stationId, _ := c.GetInt("stationId")
	if stationId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "岗位ID错误!"}
		return
	}
	var station *models.SysStation
	station = &models.SysStation{}
	stationName := c.GetString("stationName")
	roleId, _ := c.GetInt("roleId")
	if roleId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "角色ID错误!"}
		return
	} else {
		role, _ := models.SysRoleByRid(roleId)
		if role == nil {
			c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "角色不存在!"}
			return
		}
	}
	typeStr := c.GetString("typeStr")
	typeArr := strings.Split(typeStr, ",")
	dataStr := c.GetString("dataStr")
	dataArr := strings.Split(dataStr, ",")
	if station == nil {
		station = &models.SysStation{}
	}
	orgId, _ := c.GetInt("orgId")
	if orgId <= 0 {
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "岗位ID错误!"}
		return
	}
	station.Name = stationName
	station.RoleId = roleId
	station.Id = stationId
	station.OrgId = orgId
	err = station.Update(typeArr, dataArr)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "修改失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "修改失败!"}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "msg": "修改成功!"}
	}
}

//根据组织架构Id获取岗位信息
func (c *SystemController) GetStationData() {
	defer c.ServeJSON()
	orgId, _ := c.GetInt("orgId")
	m, _ := models.SysStationListByOrgId(orgId)
	c.Data["json"] = map[string]interface{}{"ret": 200, "stationData": m}
}

//获取角色列表
func (c *SystemController) GetRoleList() {
	defer c.ServeJSON()
	rolelist, err := models.SysRoleList()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取角色列表失败", err.Error(), c.Ctx.Input)
	}
	c.Data["json"] = map[string]interface{}{"ret": 200, "roleListData": rolelist}
}

//删除岗位
func (c *SystemController) DelStation() {
	defer c.ServeJSON()
	stationId, err := c.GetInt("stationId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "id错误", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id错误"}
		return
	}
	err = models.DelStation(stationId)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "删除岗位失败", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200}
	}
}

//根据岗位ID，获取岗位信息
func (c *SystemController) GetStationById() {
	defer c.ServeJSON()
	stationId, err := c.GetInt("stationId")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "id错误", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": "id错误"}
		return
	}
	result, err := models.SysStationById(stationId)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取岗位信息失败", err.Error(), c.Ctx.Input)
		}
		c.Data["json"] = map[string]interface{}{"ret": 304, "msg": err.Error()}
	} else {
		c.Data["json"] = map[string]interface{}{"ret": 200, "data": result}
	}
}

//获取组织架构和岗位信息
func (c *SystemController) GetOrganizationStation() {
	s := models.QueryDisplayQn()
	o, _ := models.GetOrganizationStations()
	c.Data["json"] = map[string]interface{}{"ret": 200, "organizationStationList": o, "station": s}
	c.ServeJSON()
}
