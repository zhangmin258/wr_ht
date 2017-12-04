package cache

import (
	"encoding/json"
	"wr_v1/models"
	"wr_v1/utils"
)

func GetWrDataCount(source, cacheStr string) (*models.WeiRongDataAll, error) {
	if utils.Re == nil && utils.Rc.IsExist(cacheStr) {
		if data, err := utils.Rc.RedisBytes(cacheStr); err == nil {
			var m *models.WeiRongDataAll
			err = json.Unmarshal(data, &m)
			if err == nil {
				return m, err
			}
		}
	}
	return models.GetWrRegisterAllCount(source)
}

func GetWrData(source, cacheStr string) ([]models.DailyData, error) {
	if utils.Re == nil && utils.Rc.IsExist(cacheStr) {
		if data, err := utils.Rc.RedisBytes(cacheStr); err == nil {
			var m []models.DailyData
			err = json.Unmarshal(data, &m)
			if err == nil {
				return m, err
			}
		}
	}
	return models.GetDailyDataCache(source)
}

func GetWrSCDataCount(source, cacheStr string, params []string, pkgParams []int) (*models.WeiRongDataAll, error) {
	if utils.Re == nil && utils.Rc.IsExist(cacheStr) {
		if data, err := utils.Rc.RedisBytes(cacheStr); err == nil {
			var m *models.WeiRongDataAll
			err = json.Unmarshal(data, &m)
			if err == nil {
				return m, err
			}
		}
	}
	return models.GetWrSCRegisterAllCount(source, params, pkgParams)
}

func GetWrSCData(source, cacheStr string, params []string, pkgParams []int) ([]models.DailyData, error) {
	if utils.Re == nil && utils.Rc.IsExist(cacheStr) {
		if data, err := utils.Rc.RedisBytes(cacheStr); err == nil {
			var m []models.DailyData
			err = json.Unmarshal(data, &m)
			if err == nil {
				return m, err
			}
		}
	}
	return models.GetSCDailyDataCache(source, params, pkgParams)
}
