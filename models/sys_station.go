package models

import (
	"github.com/astaxie/beego/orm"

	"strconv"
	"strings"
)

//1-通用岗位类型
//5-M3+
const (
	M1 = iota + 2 //逾期[1,29]
	M2            //逾期[30,59]
	M3            //逾期[60,89]
	M0 = 6        //当天到期
)

//岗位
type SysStation struct {
	Id     int    `orm:"column(id);pk"`
	Name   string `orm:"column(name);null"`    //岗位名称
	RoleId int    `orm:"column(role_id);null"` //角色ID
	OrgId  int    `orm:"column(org_id);null"`  //组织架构ID
}

//岗位类型
type SysStationType struct {
	Id        int    `orm:"column(id);pk"`
	StationId int    `orm:"column(station_id);null"` //岗位ID
	Type      int    `orm:"column(type);null"`       //类型
	TypeName  string `orm:"column(type_name);null"`  //岗位类型名称
}

//岗位-数据权限
type SysStationData struct {
	Id        int `orm:"column(id);pk"`
	StationId int `orm:"column(station_id);null"` //岗位ID
	OrgId     int `orm:"column(org_id);null"`     //组织机构ID
}

//根据名称获取岗位
func GetSysStationByName(stationName string, orgid int) (station *SysStation, err error) {
	sql := `SELECT DISTINCT a.id,a.name,a.role_id from sys_station a where a.name=? and org_id=?`
	err = orm.NewOrm().Raw(sql, stationName, orgid).QueryRow(&station)
	return
}

//根据组织架构ID获取岗位信息
func SysStationListByOrgId(orgId int) (list []*SysStation, err error) {
	sql := `SELECT a.id,a.name,a.role_id,a.org_id FROM sys_station AS a where a.org_id=?`
	_, err = orm.NewOrm().Raw(sql, orgId).QueryRows(&list)
	return
}

//根据岗位ID获取岗位信息
func SysStationById(stationId int) (map[string]interface{}, error) {
	var dataList []*SysStationData
	sql := `SELECT id,station_id,org_id FROM sys_station_data AS a where station_id=?`
	_, err := orm.NewOrm().Raw(sql, stationId).QueryRows(&dataList)
	var typeList []*SysStationType
	sql = `SELECT id,station_id,type FROM sys_station_type where station_id=?`
	_, err = orm.NewOrm().Raw(sql, stationId).QueryRows(&typeList)
	result := make(map[string]interface{})
	result["dataList"] = dataList
	result["typeList"] = typeList
	return result, err
}

//添加岗位信息
func (station *SysStation) Insert(stationType, stationData []string) error {
	o := orm.NewOrm()
	sql := `INSERT INTO sys_station(name,role_id,org_id)values(?,?,?)`
	o.Begin()
	res, err := o.Raw(sql, station.Name, station.RoleId, station.OrgId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	//获取岗位ID
	stationId, err := res.LastInsertId()
	if err != nil {
		o.Rollback()
		return err
	}
	//添加岗位类型
	if len(stationType) > 0 {
		sql = `INSERT INTO sys_station_type(station_id,type)values(?,?)`
		typeSql, err := o.Raw(sql).Prepare()
		for i := 0; i < len(stationType); i++ {
			typeInt, _ := strconv.Atoi(stationType[i])
			if typeInt > 0 {
				_, err = typeSql.Exec(stationId, typeInt)
				if err != nil {
					break
				}
			}
		}
		typeSql.Close()
		if err != nil {
			o.Rollback()
			return err
		}
	}
	//添加岗位权限
	if len(stationData) > 0 {
		sql = `INSERT INTO sys_station_data(station_id,org_id)values(?,?)`
		dataSql, err := o.Raw(sql).Prepare()
		defer dataSql.Close()
		for i := 0; i < len(stationData); i++ {
			dataInt, _ := strconv.Atoi(stationData[i])
			if dataInt > 0 {
				_, err = dataSql.Exec(stationId, dataInt)
				if err != nil {
					break
				}
			}
		}
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

//修改岗位信息
func (station *SysStation) Update(stationType, stationData []string) error {
	o := orm.NewOrm()
	sql := `UPDATE sys_station SET name=?,role_id=?,org_id=? WHERE id=?`
	o.Begin()
	_, err := o.Raw(sql, station.Name, station.RoleId, station.OrgId, station.Id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	//删除岗位类型
	sql = `DELETE FROM sys_station_type where station_id=?`
	_, err = o.Raw(sql, station.Id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	//添加岗位类型
	if len(stationType) > 0 {
		sql = `INSERT INTO sys_station_type(station_id,type)values(?,?)`
		stationSql, err := o.Raw(sql).Prepare()
		for i := 0; i < len(stationType); i++ {
			typeInt, _ := strconv.Atoi(stationType[i])
			_, err = stationSql.Exec(station.Id, typeInt)
		}
		stationSql.Close()
		if err != nil {
			o.Rollback()
			return err
		}
	}
	//删除岗位权限
	sql = `DELETE FROM sys_station_data where station_id=?`
	_, err = o.Raw(sql, station.Id).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	//添加岗位权限
	if len(stationData) > 0 {
		sql = `INSERT INTO sys_station_data(station_id,org_id)values(?,?)`
		insertSql, err := o.Raw(sql).Prepare()
		for i := 0; i < len(stationData); i++ {
			dataInt, _ := strconv.Atoi(stationData[i])
			_, err = insertSql.Exec(station.Id, dataInt)
		}
		insertSql.Close()
		if err != nil {
			o.Rollback()
			return err
		}
	}
	o.Commit()
	return nil
}

//删除岗位
func DelStation(stationId int) error {
	//删除岗位
	sql := `DELETE FROM sys_station WHERE id=?`
	o := orm.NewOrm()
	o.Begin()
	_, err := o.Raw(sql, stationId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	//删除岗位类型
	sql = `DELETE FROM sys_station_type WHERE station_id=?`
	_, err = o.Raw(sql, stationId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	//删除岗位权限
	sql = `DELETE FROM sys_station_data WHERE station_id=?`
	_, err = o.Raw(sql, stationId).Exec()
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

//根据岗位获取数据权限
func GetDataByStation(stationId int) (str string, err error) {
	var dataList []*SysStationData
	sql := `SELECT DISTINCT org_id FROM sys_station_data WHERE station_id=? ORDER BY org_id ASC`
	_, err = orm.NewOrm().Raw(sql, stationId).QueryRows(&dataList)
	ids := make([]string, 0, 10)
	if err == nil && len(dataList) > 0 {
		for _, v := range dataList {
			if v.OrgId > 0 {
				ids = append(ids, strconv.Itoa(v.OrgId))
			}
		}
	}
	str = strings.Join(ids, ",")
	return str, err
}

//根据岗位ID获取用户页面按钮权限
func GetPower(stationId int) bool {
	var count int
	sql := `SELECT count(1) FROM sys_station AS a
			INNER JOIN sys_station_data AS b
			ON a.id=b.station_id
			AND a.org_id=b.org_id
			WHERE a.id=? AND (A.name='贷后总监' or A.name='超级管理员')`
	err := orm.NewOrm().Raw(sql, stationId).QueryRow(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
