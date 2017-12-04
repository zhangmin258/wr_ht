package models

import (
	"github.com/astaxie/beego/orm"
)

type SysMenuList []*SysMenu

func (list SysMenuList) Len() int {
	return len(list)
}

func (list SysMenuList) Less(i, j int) bool {
	return list[i].SortIndex < list[j].SortIndex
}

func (list SysMenuList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

//系统菜单
type SysMenu struct {
	Id          int         `orm:"column(id);pk"`
	DisplayName string      `orm:"column(displayname);null"`
	ControlUrl  string      `orm:"column(controlurl);null"`
	ParentId    int         `orm:"column(parent_id);null"`
	SortIndex   int         `orm:"column(sortindex);null"`
	ChildMenu   SysMenuList `orm:"-"`
}

func GetSysMenuTreeByRoleId(role_id int) ([]SysMenu, error) {
	o := orm.NewOrm()
	sql := "SELECT DISTINCT d.* from  sys_role_menu c INNER JOIN sys_menu d on c.menu_id=d.id WHERE d.isvisible=1 and c.role_id =? order by sortindex "
	res := []SysMenu{}
	_, err := o.Raw(sql, role_id).QueryRows(&res)
	return res, err
}

//根据岗位ID获取菜单信息
func GetSysMenuTreeByStationId(stationId int) ([]SysMenu, error) {
	o := orm.NewOrm()
	sql := `SELECT me.* FROM sys_station AS m
	INNER JOIN sys_role_menu AS r
	ON m.role_id=r.role_id
	INNER JOIN sys_menu AS me
	ON r.menu_id=me.id
	WHERE m.id=?
	AND isvisible=1
	ORDER BY sortindex`
	res := []SysMenu{}
	_, err := o.Raw(sql, stationId).QueryRows(&res)
	return res, err
}

func GetSysMenuTreeAll() (list []SysMenu, err error) {
	sql := `SELECT * FROM sys_menu WHERE isvisible=1 order by sortindex`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&list)
	return
}
