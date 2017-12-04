package models

import (
	"github.com/astaxie/beego/orm"
)

//产品url
type ProductUrl struct {
	Id              int    `orm:"pk"` //
	ProductId       int    //产品ID
	UrlType         int    //跳转页面类型1:小贷平台落地页2：账单详情 3：提现页面4：下载前链接 5：下载后链接----1
	Url             string //跳转URL
	RegsteUrlBefore string //下载前链接
	RegsteUrlAfter  string //下载后链接
}

//产品url
type ProductUrls struct {
	UrlId           int //
	BeforeId        int
	AfterId         int
	ProductId       int    //产品ID
	UrlType         int    //跳转页面类型1:小贷平台落地页2：账单详情 3：提现页面4：下载前链接 5：下载后链接----1
	Url             string //跳转URL
	RegsteUrlBefore string //下载前链接
	RegsteUrlAfter  string //下载后链接
}

//外放代理链接
type ProductAgentUrl struct {
	Id     int    `orm:"pk"` //主键
	ProId  int    //(11)产品id
	Url    string //外放链接
	IsUsed int    //是否使用：1，使用2，冻结
	HasOwn int    //是否分配：1 未分配 2已分配
}

// 保存主链接
func AddMainUrl(productUrl *ProductUrl) (err error) {
	o := orm.NewOrm()
	o.Begin()
	defer func() {
		if err != nil {
			o.Rollback()
			return
		}
		o.Commit()
	}()
	var count int
	sql := `INSERT INTO product_url (product_id,url_type,url) VALUES(?,?,?)`
	countSql := `SELECT COUNT(1) FROM product_url WHERE product_id=? AND url_type=? `
	err = o.Raw(countSql, productUrl.ProductId, 1).QueryRow(&count)
	if err != nil {
		return
	}
	if count == 0 {
		_, err = o.Raw(sql, productUrl.ProductId, 1, productUrl.Url).Exec()
		if err != nil {
			return
		}
	}

	err = o.Raw(countSql, productUrl.ProductId, 4).QueryRow(&count)
	if err != nil {
		return
	}
	if count == 0 {
		_, err = o.Raw(sql, productUrl.ProductId, 4, productUrl.RegsteUrlBefore).Exec()
		if err != nil {
			return
		}
	}
	err = o.Raw(countSql, productUrl.ProductId, 5).QueryRow(&count)
	if err != nil {
		return
	}
	if count == 0 {
		_, err = o.Raw(sql, productUrl.ProductId, 5, productUrl.RegsteUrlAfter).Exec()
		if err != nil {
			return
		}
	}
	return
}

// 保存代理链接
func AddAgentUrl(productId int, agentUrl *[]string) (err error) {
	sql := `INSERT INTO product_agent_url ( pro_id, url, is_used, has_own) VALUES ( ?, ?, 1, 1);` //默认未冻结，未分配
	prepare, err := orm.NewOrm().Raw(sql).Prepare()
	if err != nil {
		return
	}
	for _, url := range *agentUrl {
		_, err = prepare.Exec(productId, url)
		if err != nil {
			return
		}
	}
	defer prepare.Close()
	return
}

//根据产品id查询主链接
func SearchMainUrlById(proId int) (mainUrl []ProductUrl, err error) {
	sql := `SELECT id,product_id,url_type,url FROM product_url WHERE product_id=? `
	_, err = orm.NewOrm().Raw(sql, proId).QueryRows(&mainUrl)
	return
}

//根据产品id查询代理链接
func SearchAgentUrlById(proId int) (agentUrlList []*ProductAgentUrl, err error) {
	sql := `SELECT id, pro_id, url, is_used, has_own FROM product_agent_url WHERE pro_id=?`
	_, err = orm.NewOrm().Raw(sql, proId).QueryRows(&agentUrlList)
	return
}

//根据id更新主url
func UpdateMianUrlById(productUrl *ProductUrls) (err error) {
	o := orm.NewOrm()
	o.Begin()
	sqlUpdate := `UPDATE product_url SET product_id=?,url_type=?,url=? WHERE id = ?`
	sqlInsert := `INSERT INTO product_url (product_id,url_type,url) VALUES(?,?,?) `
	productUrl.UrlType = 1
	if productUrl.UrlId != 0 {
		_, err = o.Raw(sqlUpdate, productUrl.ProductId, productUrl.UrlType, productUrl.Url, productUrl.UrlId).Exec()
		if err != nil {
			o.Rollback()
			return
		}
	} else {
		_, err = o.Raw(sqlInsert, productUrl.ProductId, productUrl.UrlType, productUrl.Url).Exec()
		if err != nil {
			o.Rollback()
			return
		}
	}
	if productUrl.BeforeId != 0 {
		productUrl.UrlType = 4
		_, err = o.Raw(sqlUpdate, productUrl.ProductId, productUrl.UrlType, productUrl.RegsteUrlBefore, productUrl.BeforeId).Exec()
		if err != nil {
			o.Rollback()
			return
		}
	} else {
		productUrl.UrlType = 4
		_, err = o.Raw(sqlInsert, productUrl.ProductId, productUrl.UrlType, productUrl.RegsteUrlBefore).Exec()
		if err != nil {
			o.Rollback()
			return
		}
	}
	if productUrl.AfterId != 0 {
		productUrl.UrlType = 5
		_, err = o.Raw(sqlUpdate, productUrl.ProductId, productUrl.UrlType, productUrl.RegsteUrlAfter, productUrl.AfterId).Exec()
		if err != nil {
			o.Rollback()
			return
		}
	} else {
		productUrl.UrlType = 5
		_, err = o.Raw(sqlInsert, productUrl.ProductId, productUrl.UrlType, productUrl.RegsteUrlAfter).Exec()
		if err != nil {
			o.Rollback()
			return
		}
	}
	o.Commit()
	return
}

//更新代理链接
func UpdateAgentUrlList(agentUrlList []*ProductAgentUrl) (err error) {
	sql := `UPDATE wr.product_agent_url SET  url = ? WHERE id = ? `
	prepare, err := orm.NewOrm().Raw(sql).Prepare()
	defer prepare.Close()
	if err != nil {
		return
	}
	for _, agentUrl := range agentUrlList {
		_, err = prepare.Exec(agentUrl.Url, agentUrl.Id)
		if err != nil {
			return
		}
	}
	return
}
