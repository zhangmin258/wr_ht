package cache

import (
	"encoding/json"
	"strconv"
	"wr_v1/models"
	"wr_v1/utils"
)

//岗位所拥有的菜单权限-树结构
func GetSysMenuTreeByRoleId(role_id int) (m models.SysMenuList, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyRoleMenuTreePrefix+strconv.Itoa(role_id)) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeyRoleMenuTreePrefix + strconv.Itoa(role_id)); err == nil {
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return
}

//组织架构信息
func GetSysOrganization() ([]models.SysOrganization, error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeySystemOrganization) {
		if data, err := utils.Rc.RedisBytes(utils.CacheKeySystemOrganization); err == nil {
			var m []models.SysOrganization
			err = json.Unmarshal(data, &m)
			if m != nil {
				return m, err
			}
		}
	}
	return models.GetSysOrganization()
}
