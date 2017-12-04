package cache

import (
	// "github.com/astaxie/beego/context"
	// "net/url"
	"encoding/json"
	"wr_v1/models"
	"wr_v1/utils"
)

//按天查看今日以前数据明细
func GetBeforeDailyDatasCache() (detailDatas []models.DetailData, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyDailyDatas) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeyDailyDatas); err1 == nil {
			err = json.Unmarshal(data, &detailDatas)
			if detailDatas != nil {
				return
			}
		}
	}
	return models.GetBeforeDailyDatas()
}

//今日以前活跃用户明细
func GetBeforeActiveDataCache() (detailDatas []models.ActiveData, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyActiveData) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeyActiveData); err1 == nil {
			err = json.Unmarshal(data, &detailDatas)
			if detailDatas != nil {
				return
			}
		}
	}
	return models.GetBeforeActiveData()
}

//获取今日以前注册/申请用户数
func GetBeforeRegisterUsersCache() (rus []models.RegisterUser, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyRegisterUsers) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeyRegisterUsers); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.GetBeforeRegisterUsers()
}

//获取今日以前注册/申请用户数
func GetBeforeIdentifyUsersCache() (rus []models.RegisterUser, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyIdentifyUsers) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeyIdentifyUsers); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.GetBeforeIdentifyUsers()
}

//获取今日以前申请贷款次数
func GetBeforeLoanUsersCache() (rus []models.RegisterUser, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyLoanUsers) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeyLoanUsers); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.GetBeforeLoanUsers()
}

//获取今日以前放款用户数
func GetBeforeCreditUsersCache() (rus []models.RegisterUser, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyCreditUsers) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeyCreditUsers); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.GetBeforeCreditUsers()
}

//获取今日以前活跃用户数
func GetBeforeActiveUsersCache() (rus []models.RegisterUser, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKeyActiveUsers) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKeyActiveUsers); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.GetBeforeActiveUsers()
}

//获取遥望今日之前注册用户数
func QueryYaoWangBeforeRegisterCountCache() (rus []models.RegisterUser, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKayYaoWangRegisterCount) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKayYaoWangRegisterCount); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.QueryYaoWangBeforeRegisterCount()
}

//获取遥望今日之前登录用户数
func QueryYaoWangBeforetLoginCountCache() (rus []models.RegisterUser, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKayYaoWangLoginCount) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKayYaoWangLoginCount); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.QueryYaoWangBeforetLoginCount()
}

//获取遥望今日之前总UV
func QueryYaoWangBeforeTotalClickCountCache() (rus []models.RegisterUser, err error) {
	if utils.Re == nil && utils.Rc.IsExist(utils.CacheKayYaoWangTotalClickCount) {
		if data, err1 := utils.Rc.RedisBytes(utils.CacheKayYaoWangTotalClickCount); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.QueryYaoWangBeforeTotalClickCount()
}

//获取遥望今日之前注册用户数
func QueryYWBeforeRegisterCountCache(source, sourceCondition string, sourceParam []string) (rus []models.RegisterUser, err error) {
	key := utils.CacheKayYaoWangRegisterCount
	if source != "" {
		key += source
	}
	if utils.Re == nil && utils.Rc.IsExist(key) {
		if data, err1 := utils.Rc.RedisBytes(key); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.QueryYWBeforeRegisterCount(key, sourceCondition, sourceParam)
}

//获取遥望今日之前登录用户数
func QueryYWBeforetLoginCountCache(source, sourceCondition string, sourceParam []string) (rus []models.RegisterUser, err error) {
	key := utils.CacheKayYaoWangLoginCount
	if source != "" {
		key += source
	}
	if utils.Re == nil && utils.Rc.IsExist(key) {
		if data, err1 := utils.Rc.RedisBytes(key); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.QueryYWBeforetLoginCount(key, sourceCondition, sourceParam)
}

//获取遥望今日之前总UV
func QueryYWBeforeTotalClickCountCache(source string) (rus []models.RegisterUser, err error) {
	key := utils.CacheKayYaoWangTotalClickCount
	if source != "" {
		key += source
	}
	if utils.Re == nil && utils.Rc.IsExist(key) {
		if data, err1 := utils.Rc.RedisBytes(key); err1 == nil {
			err = json.Unmarshal(data, &rus)
			if len(rus) != 0 {
				return
			}
		}
	}
	return models.QueryYWBeforeTotalClickCount(key, source)
}
