package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"wr_v1/utils"
)

type Config struct {
	Id          int    //主键
	ConfigKey   string //配置键名
	ConfigValue string //配置键值
	ConfigDesc  string //配置描述
	Remark      string //备注
	ConfigUrl   string //
	UrlParam    int    //
	Title       string //
	Code        string //配置唯一标识
}

type ConfigLog struct {
	Id          int       //主键
	ConfigKey   string    //配置键名
	ConfigValue string    //配置键值
	ConfigDesc  string    //配置描述
	Remark      string    //备注
	ConfigUrl   string    //
	UrlParam    int       //
	Title       string    //
	Code        string    //配置唯一标识
	UserId      int       //操作人id
	UserName    string    //操作人姓名
	ModifyTime  time.Time //修改数据时间
	State       int       //0:修改前数据，1:修改后数据
}

// 分页查询配置数据
func GetConfigsList(condition string, params []string, begin, size int) (config []Config, err error) {
	sql := `SELECT id,config_key,
	 config_value,config_desc,
	 remark,config_url,
	 url_param,title,code
	 FROM config WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY id LIMIT ?, ?"
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&config)
	return
}

// 查询所有config数量
func GetConfigsCount(condition string, params []string) (count int, err error) {
	sql := `SELECT count(1) FROM config
			WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

//根据id获取config信息
func GetConfigById(id int) (config Config, err error) {
	sql := `SELECT * FROM config WHERE id=?`
	err = orm.NewOrm().Raw(sql, id).QueryRow(&config)
	return
}

// 保存config配置信息
func AddConfig(config *Config) error {
	sql := `INSERT INTO config
	(config_key,config_value,config_desc,remark,config_url,url_param,title,code)
	 VALUES
	 (?,?,?,?,?,?,?,?)`
	_, err := orm.NewOrm().Raw(sql, config.ConfigKey, config.ConfigValue, config.ConfigDesc,
		config.Remark, config.ConfigUrl, config.UrlParam, config.Title, config.Code).Exec()
	return err
}

//修改config配置信息
func UpdateConfig(configNew *Config, configOld Config, user *SysUser) error {
	o := orm.NewOrm()

	//备份config信息到config_log
	o.Using("wr_log")
	sql := `INSERT INTO config_log
	(config_key,config_value,config_desc,remark,config_url,url_param,title,code,user_id,modify_time,state)
	 VALUES
	 (?,?,?,?,?,?,?,?,?,?,?)`
	_, err := o.Raw(sql, configOld.ConfigKey, configOld.ConfigValue, configOld.ConfigDesc,
		configOld.Remark, configOld.ConfigUrl, configOld.UrlParam, configOld.Title, configOld.Code, user.Id, time.Now(), 0).Exec()

	//备份修改后的config信息到config_log
	sql = `INSERT INTO config_log
	(config_key,config_value,config_desc,remark,config_url,url_param,title,code,user_id,modify_time,state)
	 VALUES
	 (?,?,?,?,?,?,?,?,?,?,?)`
	_, err = o.Raw(sql, configNew.ConfigKey, configNew.ConfigValue, configNew.ConfigDesc, configNew.Remark,
		configNew.ConfigUrl, configNew.UrlParam, configNew.Title, configNew.Code, user.Id, time.Now(), 1).Exec()

	//update config信息
	o.Using("default")
	sql = `UPDATE config SET config_key = ?,config_value=?,config_desc=?,remark=?,config_url=?,url_param=?,title=? WHERE id=?`
	_, err = o.Raw(sql, configNew.ConfigKey, configNew.ConfigValue, configNew.ConfigDesc, configNew.Remark,
		configNew.ConfigUrl, configNew.UrlParam, configNew.Title, configNew.Id).Exec()

	//删除缓存中的config信息
	if utils.Re == nil && utils.Rc.IsExist(utils.WR_CACHE_KEY_CONFIG+configNew.Code) {
		utils.Rc.Delete(utils.WR_CACHE_KEY_CONFIG + configNew.Code)
	}
	return err
}

//删除config配置信息
func DelConfig(configOld Config, user *SysUser) error {
	o := orm.NewOrm()

	//备份config信息到config_log
	o.Using("wr_log")
	sql := `INSERT INTO config_log
	(config_key,config_value,config_desc,remark,config_url,url_param,title,code,user_id,modify_time,state)
	 VALUES
	 (?,?,?,?,?,?,?,?,?,?,?)`
	_, err := o.Raw(sql, configOld.ConfigKey, configOld.ConfigValue, configOld.ConfigDesc,
		configOld.Remark, configOld.ConfigUrl, configOld.UrlParam, configOld.Title, configOld.Code, user.Id, time.Now(), 0).Exec()

	//删除config信息
	o.Using("default")
	sql = `DELETE FROM config WHERE id=?`
	_, err = o.Raw(sql, configOld.Id).Exec()
	//删除缓存中的config信息
	if utils.Re == nil && utils.Rc.IsExist(utils.WR_CACHE_KEY_CONFIG+configOld.Code) {
		utils.Rc.Delete(utils.WR_CACHE_KEY_CONFIG + configOld.Code)
	}
	return err
}
