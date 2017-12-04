package models

import "github.com/astaxie/beego/orm"

type CreditCard struct {
	Id              int
	ImgUrl          string
	LinkUrl         string
	Title           string
	Sort            int
	IsUsed          int8
}
//查询信用卡
func ShowCreditCardList(condition string, params []string,begin, size int) (this []CreditCard,err error){
	sql := `SELECT id,img_url, title, sort, is_used FROM credit_card_product WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	sql+=" ORDER BY id LIMIT ?, ?"
	_, err = orm.NewOrm().Raw(sql,params,begin,size).QueryRows(&this)
	return
}
// 查询所有信用卡数量
func GetCreditCardCount(condition string, params []string) (count int, err error) {
	sql := `SELECT count(1) FROM credit_card_product WHERE 1=1`
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql,params).QueryRow(&count)
	return
}
//根据ID查询信用卡信息
func GetCreditCardById(pid int)(card CreditCard, err error){
	sql := `SELECT id,img_url, link_url,title, sort, is_used FROM credit_card_product WHERE id=?`
	err = orm.NewOrm().Raw(sql,pid).QueryRow(&card)
	return
}
//新增信用卡信息
func AddCreditCard(card *CreditCard)(err error){
	sql := `INSERT INTO credit_card_product
	       (title,img_url,link_url,sort,is_used)
	       VALUES
	       (?,?,?,?,?)`
	_,err = orm.NewOrm().Raw(sql,card.Title,card.ImgUrl,card.LinkUrl,card.Sort,card.IsUsed).Exec()
	return
}
//修改信用卡信息
func UpdateCreditCard(card *CreditCard)(err error){
	sql := `UPDATE credit_card_product SET img_url=?,link_url=?,title=?,sort=?,is_used=?
	        WHERE id=?`
	_,err = orm.NewOrm().Raw(sql,card.ImgUrl,card.LinkUrl,card.Title,card.Sort,card.IsUsed,card.Id).Exec()
	return
}
//根据ID删除信用卡信息
func DeleteCreditCard(pid int)(err error){
	sql := `DELETE FROM credit_card_product WHERE id=?`
	_, err = orm.NewOrm().Raw(sql, pid).Exec()
	return
}
