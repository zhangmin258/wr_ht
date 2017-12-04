package models

import (
	"github.com/astaxie/beego/orm"
	"wr_v1/utils"
	"strconv"
)

type ProductUlr struct {
	Id  int
	Url string
}

//根据id查询商品的外放链接
func GetProductUrlList(id int) (url []*ProductUlr, err error) {
	sql := `SELECT id,url FROM product_agent_url WHERE is_used=1 AND has_own=1 AND pro_id=?`
	_, err = orm.NewOrm().Raw(sql, id).QueryRows(&url)
	return
}

//根据id获取URL
func GetUrlById(idList []string) (urlList []*ProductUlr, err error) {
	sql := `SELECT id,url FROM product_agent_url WHERE id=?`
	if len(idList)==0 {
		return urlList,err
	}
	ids:=[]int{}
	i:=0
	for k, v := range idList {
		if k!=len(idList)-1 {
			sql+=` OR id=?`
		}
		i,err=strconv.Atoi(v)
		ids=append(ids,i)
	}
	_,err = orm.NewOrm().Raw(sql,ids).QueryRows(&urlList)
	return
}

//根据产品id获取URL
func GetUrlByProId(proId int) (allUrlList []*ProductUlr, err error) {
	sql := `SELECT id,url FROM product_agent_url WHERE pro_id=? AND is_used = 1 AND has_own=1`
	_, err = orm.NewOrm().Raw(sql, proId).QueryRows(&allUrlList)
	return
}

// 根据id修改代理产品状态
func ChangeAgentProUrlById(urls []int) (err error) {

	sql := `UPDATE product_agent_url  SET  has_own = 2 WHERE id = ?`
	prepare, err := orm.NewOrm().Raw(sql).Prepare()
	defer prepare.Close()
	if err != nil {
		return
	}
	for _, url := range urls {
		_, err = prepare.Exec(url)
		if err != nil {
			return
		}
	}
	return
}

// 根据id修改代理产品状态
func ChangeOldAgentProUrlById(urls []int) (err error) {

	sql := `UPDATE product_agent_url  SET  has_own = 1 WHERE id = ?`
	prepare, err := orm.NewOrm().Raw(sql).Prepare()
	defer prepare.Close()
	if err != nil {
		return
	}
	for _, url := range urls {
		_, err = prepare.Exec(url)
		if err != nil {
			return
		}
	}
	return
}

// 批量插入代理微融url
func AddAgentProUrl(urls []string) (idList []int, err error) {
	sql := ` INSERT INTO wr.product_agent_url ( pro_id, url, is_used, has_own ) VALUES ( 0, ?, 1 , 2);`
	prepare, err := orm.NewOrm().Raw(sql).Prepare()
	defer prepare.Close()
	if err != nil {
		return
	}
	for _, url := range urls {
		reslut, err := prepare.Exec(utils.OutPutURL + url)
		if err != nil {
			break
		}
		id, err := reslut.LastInsertId()
		idList = append(idList, int(id))
	}
	return
}

// 批量修改代理微融url
func UpdateAgentProUrl(urls []ProductAgentUrl) ( err error) {
	sql := ` UPDATE product_agent_url SET url = ? WHERE id = ? ;`
	prepare, err := orm.NewOrm().Raw(sql).Prepare()
	defer prepare.Close()
	if err != nil {
		return
	}
	for _, url := range urls {
		_, err := prepare.Exec(utils.OutPutURL + url.Url,url.Id)
		if err != nil {
			break
		}
	}
	return
}
