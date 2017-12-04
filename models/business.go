package models

import (
	"github.com/astaxie/beego/orm"
)

type Business struct {
	Id         int
	AppId      string //商务号
	PublicKey  string //第三方公钥
	PrivateKey string //第三方私钥
	Name       string //商务名称
	SignType   string //签名类型
	Version    string //版本
	Charset    string //字符集
}

func InsertBusiness(proName string) (err error) {
	sql := `INSERT INTO business (name) VALUES (?) `
	_, err = orm.NewOrm().Raw(sql, proName).Exec()
	return
}
