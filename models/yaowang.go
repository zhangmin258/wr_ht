package models

import (
	"encoding/json"
	"time"
	"wr_v1/utils"

	"github.com/astaxie/beego/orm"
)

type YaowangProduct struct {
	Id          int       //id
	Name        string    //名称
	Logo        string    //logo
	Url         string    //跳转链接
	IsUse       int       //是否使用 0:使用 1:不使用
	CreateTime  time.Time //创建时间
	Source      string    //来源
	Description string    //产品描述
	Tag         string    //产品标签
}

type ReturnYaowangProduct struct {
	ProductId   int       //id
	Name        string    //名称
	ImgUrl      string    //logo
	LinkUrl     string    //跳转链接
	IsUsed      int       //是否使用 0:使用 1:不使用
	CreateTime  time.Time //创建时间
	Source      string    //来源
	Description string    //产品描述
	Tag         string    //产品标签
}

type YaowangClickInfo struct {
	CreateTime  time.Time //时间
	ProductName string    //产品名称
	Count       int       //点击数量
}

type YaowangPageLoadInfo struct {
	CreateTime time.Time //时间
	LoadCount  int       //页面加载次数
}

type YaowangDataAnalysis struct {
	CreateTime  string                //时间
	LoadCount   int                   //页面加载次数
	ProductInfo []YaowangProductCount //产品点击记录
	Length      int                   //长度
}

type YaowangProductCount struct {
	ProductName string //产品名称
	ProductUV   int    //点击数量
}

type ProductDataCount struct {
	CreateDate time.Time //时间
	Count      int       //数量
}

type YaowangCountData struct {
	CreateDate  string //时间
	RegistCount int    //注册用户数量
	LoginCount  int    //登录用户数量
	UVCount     int    //UV点击数量
	ActiveCount int    //活跃用户数量
}

//获取页面加载次数
func GetPageLoadCount(condition string, params []string) (yaowangpageloadinfo []YaowangPageLoadInfo, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT DATE(create_time) create_time,COUNT(1) load_count FROM yaowang_load_record WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY DATE(create_time) DESC `
	_, err = o.Raw(sql, params).QueryRows(&yaowangpageloadinfo)
	return
}

//获取每个产品的点击数量
func GetYaowangProductUV(condition string, params []string) (yaowangclickinfos []YaowangClickInfo, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT * FROM 
			(SELECT cr.product_name,COUNT(1) AS count,DATE(cr.create_time) AS create_time 
			FROM yaowang_click_record AS cr 
			WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += ` GROUP BY cr.product_name,
			DATE(cr.create_time)) t 
		ORDER BY t.create_time DESC `
	_, err = o.Raw(sql, params).QueryRows(&yaowangclickinfos)
	return
}

//获取遥望产品列表
func GetYaowangProductRecommend(condition string, params []string, begin, size int) (products []YaowangProduct, err error) {
	o := orm.NewOrm()
	sql := `SELECT id,name,logo,url,is_use,create_time FROM yaowang_product WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY id DESC LIMIT ?, ?"
	_, err = o.Raw(sql, params, begin, size).QueryRows(&products)
	return
}

//获取遥望产品数量
func GetYaowangProductCount(condition string, params []string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(1) FROM yaowang_product WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY id DESC"
	err = o.Raw(sql, params).QueryRow(&count)
	return
}

//添加产品
func AddYaowangProduct(product ReturnYaowangProduct) (id int, err error) {
	o := orm.NewOrm()
	sql := `INSERT INTO yaowang_product
			(name,logo,url,is_use,create_time,description,tag) 
			VALUES(?,?,?,?,?,?,?)`
	result, err := o.Raw(sql, product.Name, product.ImgUrl, product.LinkUrl, product.IsUsed, product.CreateTime, product.Description, product.Tag).Exec()
	if err != nil {
		return
	}
	pid, err := result.LastInsertId()
	return int(pid), err
}

//修改产品
func UpdateYaowangProduct(product ReturnYaowangProduct) (err error) {
	o := orm.NewOrm()
	sql := `UPDATE yaowang_product SET name=?,logo=?,url=?,is_use=?,description=?,tag=? WHERE id = ?`
	_, err = o.Raw(sql, product.Name, product.ImgUrl, product.LinkUrl, product.IsUsed, product.Description, product.Tag, product.ProductId).Exec()
	return
}

//根据id查找遥望产品
func GetYaowangProductById(id int) (product YaowangProduct, err error) {
	o := orm.NewOrm()
	sql := `SELECT id,name,logo,url,is_use,create_time,description,tag FROM yaowang_product WHERE id =? `
	err = o.Raw(sql, id).QueryRow(&product)
	return
}

//获取每日总点击量
func GetTotalClickCount(begin, size int) (productdatacount []ProductDataCount, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT DATE(create_time) create_date,COUNT(1) count FROM yaowang_click_record  
			GROUP BY DATE(create_time) DESC LIMIT ?,?`
	_, err = o.Raw(sql, begin, size).QueryRows(&productdatacount)
	return
}

//获取每日总点击量
func GetYWTotalClickCount(param []string) (productdatacount []ProductDataCount, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT DATE(create_time) create_date,COUNT(1) count FROM yaowang_click_record WHERE 1 = 1 `
	if len(param) > 0 {
		sql += `AND source = ? `
	}
	sql += `GROUP BY DATE(create_time) DESC`
	_, err = o.Raw(sql, param).QueryRows(&productdatacount)
	return
}

//获取注册用户数据
func GetYaowangRegistCount(begin, size int) (productdatacount []ProductDataCount, err error) {
	o := orm.NewOrm()
	sql := `SELECT create_date,COUNT(1) count FROM users WHERE out_put_source = 'yw' 
			GROUP BY create_date ORDER BY create_date DESC LIMIT ?,? `
	_, err = o.Raw(sql, begin, size).QueryRows(&productdatacount)
	return
}

//获取注册用户数据
func GetYWRegistCount(sourceCondition string, sourceParam []string) (productdatacount []ProductDataCount, err error) {
	o := orm.NewOrm()
	sql := `SELECT create_date,COUNT(1) count FROM users WHERE 1 = 1 `
	if sourceCondition != "" {
		sql += sourceCondition
	}
	sql += `GROUP BY create_date ORDER BY create_date DESC  `
	_, err = o.Raw(sql, sourceParam).QueryRows(&productdatacount)
	return
}

//获取登录用户数据
func GetYaowangLoginCount(begin, size int) (productdatacount []ProductDataCount, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT  r.create_date,COUNT(1) count 
			FROM (SELECT uid,DATE(create_time) create_date FROM login_record WHERE out_put_source = 'yw' GROUP BY uid,DATE(create_time)) r
			GROUP BY r.create_date ORDER BY r.create_date DESC LIMIT ?,? `
	_, err = o.Raw(sql, begin, size).QueryRows(&productdatacount)
	return
}

//获取登录用户数据
func GetYWLoginCount(sourceCondition string, sourceParam []string) (productdatacount []ProductDataCount, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT  r.create_date,COUNT(1) count 
			FROM (SELECT uid,DATE(create_time) create_date FROM login_record WHERE 1 = 1 `
	if sourceCondition != "" {
		sql += sourceCondition
	}
	sql += `GROUP BY uid,DATE(create_time)) r
			GROUP BY r.create_date ORDER BY r.create_date DESC `
	_, err = o.Raw(sql, sourceParam).QueryRows(&productdatacount)
	return
}

//获取遥望数据明细条数
func GetYaowangDataCount() (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(DISTINCT(create_date)) count FROM users WHERE out_put_source = 'yw' `
	err = o.Raw(sql).QueryRow(&count)
	return
}

//获取遥望数据明细条数
func GetYWDataCount(sourceCondition string, sourceParam []string) (count int, err error) {
	o := orm.NewOrm()
	sql := `SELECT COUNT(DISTINCT(create_date)) count FROM users WHERE 1 = 1 `
	if sourceCondition != "" {
		sql += sourceCondition
	}
	err = o.Raw(sql, sourceParam).QueryRow(&count)
	return
}

//获取遥望今日之前注册用户数
func QueryYWBeforeRegisterCount(key, sourceCondition string, sourceParam []string) (rus []RegisterUser, err error) {
	sql := `SELECT create_date, COUNT(1) count
			FROM users 
			WHERE create_date < CURDATE() `
	if sourceCondition != "" {
		sql += sourceCondition
	}
	sql += `GROUP BY create_date 
			ORDER BY create_date DESC`
	_, err = orm.NewOrm().Raw(sql, sourceParam).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			utils.Rc.Put(key, data, utils.GetTodayLastSecond())
		}
	}
	return
}

//获取遥望今日之前登录用户数
func QueryYWBeforetLoginCount(key, sourceCondition string, sourceParam []string) (rus []RegisterUser, err error) {
	sql := `SELECT COUNT(DISTINCT (uid)) count, DATE(create_time) create_date 
			FROM login_record 
			WHERE 1 = 1 `
	if sourceCondition != "" {
		sql += sourceCondition
	}
	sql += `GROUP BY DATE(create_time)
			ORDER BY create_date DESC `
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err = o.Raw(sql, sourceParam).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			utils.Rc.Put(key, data, utils.GetTodayLastSecond())
		}
	}
	return
}

//获取遥望今日之前总uv
func QueryYWBeforeTotalClickCount(key, source string) (rus []RegisterUser, err error) {
	params := []string{}
	sql := `SELECT DATE(create_time) create_date, COUNT(1) count
			FROM yaowang_click_record WHERE 1 = 1 `
	if source != "" {
		sql += `AND source = ? `
		params = append(params, source)
	}
	sql += `GROUP BY DATE(create_time) 
			ORDER BY create_date DESC `
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err = o.Raw(sql, params).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			utils.Rc.Put(key, data, utils.GetTodayLastSecond())
		}
	}
	return
}

//获取遥望今日之前注册用户数
func QueryYaoWangBeforeRegisterCount() (rus []RegisterUser, err error) {
	sql := `SELECT 
				create_date,
				COUNT(1) count
			FROM
				users 
			WHERE create_date < CURDATE() 
				AND out_put_source = 'yw' 
			GROUP BY create_date 
			ORDER BY create_date DESC`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			utils.Rc.Put(utils.CacheKayYaoWangRegisterCount, data, utils.GetTodayLastSecond())
		}
	}
	return
}

//获取遥望今日之前登录用户数
func QueryYaoWangBeforetLoginCount() (rus []RegisterUser, err error) {
	sql := `SELECT 
				COUNT(DISTINCT (uid)) count,
				DATE(create_time) create_date 
			FROM
				login_record 
			WHERE out_put_source = 'yw' 
			GROUP BY DATE(create_time) 
			ORDER BY create_date DESC `
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err = o.Raw(sql).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			utils.Rc.Put(utils.CacheKayYaoWangLoginCount, data, utils.GetTodayLastSecond())
		}
	}
	return
}

//获取遥望今日之前总uv
func QueryYaoWangBeforeTotalClickCount() (rus []RegisterUser, err error) {
	sql := `SELECT 
				DATE(create_time) create_date,
				COUNT(1) count
			FROM
				yaowang_click_record 
			GROUP BY DATE(create_time) DESC `
	o := orm.NewOrm()
	o.Using("wr_log")
	_, err = o.Raw(sql).QueryRows(&rus)
	if err == nil {
		if data, err2 := json.Marshal(rus); err2 == nil && utils.Re == nil {
			utils.Rc.Put(utils.CacheKayYaoWangTotalClickCount, data, utils.GetTodayLastSecond())
		}
	}
	return
}

//获取遥望今日注册用户数
func QueryYaoWangTodayRegisterCount() (count int, err error) {
	sql := `SELECT 
				COUNT(1) COUNT
			FROM
				users 
			WHERE create_date = CURDATE() 
				AND out_put_source = 'yw' 
			GROUP BY create_date 
			ORDER BY create_date DESC 
			`
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}

//获取遥望今日注册用户数
func QueryYWTodayRegisterCount(sourceCondition string, sourceParam []string) (count int, err error) {
	sql := `SELECT COUNT(1) count
			FROM users 
			WHERE create_date = CURDATE() `
	if sourceCondition != "" {
		sql += sourceCondition
	}
	sql += `GROUP BY create_date 
			ORDER BY create_date DESC `
	err = orm.NewOrm().Raw(sql, sourceParam).QueryRow(&count)
	return
}

//获取遥望今日登录用户数
func QueryYaoWangTodayLoginCount() (count int, err error) {
	sql := `SELECT 
				COUNT(DISTINCT (uid)) count
			FROM
				login_record 
			WHERE out_put_source = 'yw' 
				AND DATE(create_time) = CURDATE()`
	o := orm.NewOrm()
	o.Using("wr_log")
	err = o.Raw(sql).QueryRow(&count)
	return
}

//获取遥望今日登录用户数
func QueryYWTodayLoginCount(sourceCondition string, sourceParam []string) (count int, err error) {
	sql := `SELECT COUNT(DISTINCT (uid)) count
			FROM login_record 
			WHERE DATE(create_time) = CURDATE()`
	if sourceCondition != "" {
		sql += sourceCondition
	}
	o := orm.NewOrm()
	o.Using("wr_log")
	err = o.Raw(sql, sourceParam).QueryRow(&count)
	return
}

//获取遥望今日总uv
func QueryYaoWangTodayTotalClickCount() (count int, err error) {
	sql := `SELECT 
				COUNT(1) COUNT
			FROM
				yaowang_click_record 
			WHERE DATE(create_time) = CURDATE()`
	o := orm.NewOrm()
	o.Using("wr_log")
	err = o.Raw(sql).QueryRow(&count)
	return
}

//获取遥望今日总uv
func QueryYWTodayTotalClickCount(source string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) COUNT
			FROM yaowang_click_record 
			WHERE DATE(create_time) = CURDATE()`
	if source != "" {
		sql += `AND source = ? `
	}
	o := orm.NewOrm()
	o.Using("wr_log")
	err = o.Raw(sql, params).QueryRow(&count)
	return
}

//获取所有遥望渠道
func GetYaowangSource() (sourceList []string, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT source FROM yaowang_load_record GROUP BY source HAVING source != '' `
	_, err = o.Raw(sql).QueryRows(&sourceList)
	return
}

//模糊查询获取渠道列表
func GetYaowangSourceList(condition string, params []string) (sourceList []string, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT source FROM yaowang_load_record WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += `GROUP BY source `
	_, err = o.Raw(sql, params).QueryRows(&sourceList)
	return
}

//判断是否有指定渠道
func GetYaowangSourceByName(source string) (count int, err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `SELECT COUNT(1) FROM yaowang_load_record WHERE source = ? `
	err = o.Raw(sql, source).QueryRow(&count)
	return
}

type YWProInfo struct {
	Id   int    //产品id
	Name string //产品名称
}

//遥望产品名称数据修复
func GetYWProductInfos() (ywProInfos []YWProInfo, err error) {
	o := orm.NewOrm()
	sql := `SELECT id,name FROM yaowang_product `
	_, err = o.Raw(sql).QueryRows(&ywProInfos)
	return
}

//修复遥望点击记录产品名称
func FixYWClickProductName(ywProInfos []YWProInfo) (err error) {
	o := orm.NewOrm()
	o.Using("wr_log")
	sql := `UPDATE yaowang_click_record SET product_name = ? WHERE yaowang_product_id = ? AND product_name = '' `
	p, err := o.Raw(sql).Prepare()
	if err != nil {
		return
	}
	defer p.Close()
	for _, v := range ywProInfos {
		_, err = p.Exec(v.Name, v.Id)
		if err != nil {
			return
		}
	}
	return
}
