package services

/*// 微融每日数据定时缓存
func WrData() error {
	//缓存累计数据
	source := " AND u.source !='' "               //微融市场数据
	outPutSource := " AND u.out_put_source !='' " //微融外链数据
	weirongCountSource, err := models.GetWrRegisterAllCount(source)
	if err != nil {
		return err
	}
	w, err := json.Marshal(weirongCountSource)
	if err != nil {
		return err
	}
	utils.Rc.Put(utils.WEIRONGCOUNTSOURCE, w, time.Hour*24)
	weirongCountOutPutSource, err := models.GetWrRegisterAllCount(outPutSource)
	if err != nil {
		return err
	}
	oPS, err := json.Marshal(weirongCountOutPutSource)
	if err != nil {
		return err
	}
	utils.Rc.Put(utils.WEIRONGCOUNTOUTPUTSOURCE, oPS, time.Hour*24)
	//缓存微融数据按每天的统计
	weirongDataCacheSource, err := models.GetDailyDataCache(source)
	if err != nil {
		return err
	}
	wc, err := json.Marshal(weirongDataCacheSource)
	if err != nil {
		return err
	}
	utils.Rc.Put(utils.WEIRONGEVERYDAYDATASOURCE, wc, time.Hour*24)

	weirongDataCacheOutPutSource, err := models.GetDailyDataCache(outPutSource)
	if err != nil {
		return err
	}
	wcOPS, err := json.Marshal(weirongDataCacheOutPutSource)
	if err != nil {
		return err
	}
	utils.Rc.Put(utils.WEIRONGEVERYDAYDATAOUTPUTSOURCE, wcOPS, time.Hour*24)
	return nil
}*/
