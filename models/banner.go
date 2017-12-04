package models

import (
	"github.com/astaxie/beego/orm"
)

//Banner
type Images struct {
	Id        int    `orm:'pk'` //图片ID
	Itype     int    //展示位置：1，首页 2，贷款页
	State     int    //跳转类型:1:跳转至产品详情 2:跳转至外放H5页面 3:不跳
	ImgUrl    string //图片路径
	ProductId int    //产品ID
	LinkUrl   string //跳转地址
	Title     string //跳转页标题
	Sort      int    //排列顺序
	IsUsed    int    //是否使用
	Name      string //产品名称
	ShowSec  int	 //展示时间
}

//保存Banner
func SaveBanner(image *Images) (err error) {
	sql := `INSERT INTO images
	        (itype,img_url,link_url,is_used,state,product_id,sort,title,show_sec)
	        VALUES
	        (?,?,?,?,?,?,?,?,?)`
	_, err = orm.NewOrm().Raw(sql, image.Itype, image.ImgUrl, image.LinkUrl, image.IsUsed, image.State, image.ProductId, image.Sort, image.Title,image.ShowSec).Exec()
	return
}

//根据ID删除Banner
func DeleteBanner(pid int) (err error) {
	sql := `DELETE FROM images WHERE id = ?`
	_, err = orm.NewOrm().Raw(sql, pid).Exec()
	return
}

//根据ID获取Banner
func GetBannerById(pid int) (image Images, err error) {
	sql := `SELECT i.id,i.itype,i.img_url,i.link_url,i.is_used,i.state,i.product_id,i.title,i.sort,i.show_sec,p.name AS name FROM images AS i INNER JOIN product AS p
	        ON i.product_id=p.id WHERE i.id=? `
	err = orm.NewOrm().Raw(sql, pid).QueryRow(&image)
	return
}

//修改banner
func UpdateBanner(image *Images) (err error) {
	sql := `UPDATE images SET itype=?,img_url=?,state=?,link_url=?,sort=?,is_used=?,product_id=?,title=?,show_sec=?
	        WHERE id=?`
	_, err = orm.NewOrm().Raw(sql, image.Itype, image.ImgUrl, image.State, image.LinkUrl, image.Sort, image.IsUsed, image.ProductId, image.Title,image.ShowSec, image.Id).Exec()
	return err
}

//查询缩略图信息
func GetAllBanner(condition string, params []int) (imageList []*Images, err error) {
	sql := `SELECT id,itype,img_url,is_used,sort FROM images WHERE 1 = 1 `
	if condition != "" {
		sql += condition
	}
	sql += " ORDER BY sort ASC "
	_, err = orm.NewOrm().Raw(sql, params).QueryRows(&imageList)
	return
}

//获取type
func GetTypeById() (t []int, err error) {
	sql := `SELECT itype FROM images GROUP BY itype`
	_, err = orm.NewOrm().Raw(sql).QueryRows(&t)
	return
}

//查询是否存在广告位
func QueryBannerType() (count int, err error) {
	sql := `SELECT COUNT(1) FROM images WHERE itype=3 AND is_used = 1 `
	err = orm.NewOrm().Raw(sql).QueryRow(&count)
	return
}
