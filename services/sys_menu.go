package services

import (
	"encoding/json"
	"sort"
	"strconv"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"
)

//获取用户菜单
func GetSysMenuTreeByRoleId(role_id int) (models.SysMenuList, error) {
	t, _ := cache.GetSysMenuTreeByRoleId(role_id)
	if t != nil {
		return t, nil
	}
	//m, err := models.GetSysMenuTreeByRoleId(role_id)
	m, err := models.GetSysMenuTreeByStationId(role_id) //根据岗位ID,获取菜单信息
	l := len(m)
	var menu models.SysMenuList
	for i := 0; i < l; i++ {
		if m[i].ParentId == 0 {
			for j := 0; j < l; j++ {
				if m[j].ParentId == m[i].Id {
					m[i].ChildMenu = append(m[i].ChildMenu, &m[j])
				}
			}
			sort.Sort(m[i].ChildMenu)
			menu = append(menu, &m[i])
		}
	}
	sort.Sort(menu)
	if data, err2 := json.Marshal(menu); err == nil && err2 == nil && utils.Re == nil {
		utils.Rc.Put(utils.CacheKeyRoleMenuTreePrefix+strconv.Itoa(role_id), data, utils.RedisCacheTime_Role)
	}
	return menu, err
}

//获取用户菜单
// func GetSysMenuTreeByRoleId(role_id int) (models.SysMenuList, error) {
// 	t, _ := cache.GetSysMenuTreeByRoleId(role_id)
// 	if t != nil {
// 		return t, nil
// 	}
// 	//m, err := models.GetSysMenuTreeByRoleId(role_id)
// 	m, err := models.GetSysMenuTreeByStationId(role_id) //根据岗位ID,获取菜单信息
// 	l := len(m)
// 	var menu map[int]int
// 	for i := 0; i < l; i++ {
// 		menu[m[i].Id] = 1
// 		// if m[i].ParentId == 0 {
// 		// 	for j := 0; j < l; j++ {
// 		// 		if m[j].ParentId == m[i].Id {
// 		// 			m[i].ChildMenu = append(m[i].ChildMenu, &m[j])
// 		// 		}
// 		// 	}
// 		// 	sort.Sort(m[i].ChildMenu)
// 		// 	menu = append(menu, &m[i])
// 		// }
// 	}
// 	// sort.Sort(menu)
// 	if data, err2 := json.Marshal(menu); err == nil && err2 == nil && utils.Re == nil {
// 		utils.Rc.Put(utils.CacheKeyRoleMenuTreePrefix+strconv.Itoa(role_id), data, utils.RedisCacheTime_Role)
// 	}
// 	return menu, err
// }

func GetSysMenuZTree(list []models.SysMenu) []map[string]interface{} {
	var menu []map[string]interface{}
	l := len(list)
	if l == 0 {
		return []map[string]interface{}{}
	}
	for i := 0; i < l; i++ {
		menu = append(menu, map[string]interface{}{"id": list[i].Id, "pId": list[i].ParentId, "open": true, "name": list[i].DisplayName})
	}
	return menu
}
