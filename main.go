package main

import (
	"github.com/astaxie/beego"
	//"github.com/astaxie/beego/orm"
	"wr_v1/controllers"
	_ "wr_v1/routers"
	"wr_v1/utils"
	"wr_v1/services"
)

func main() {
	//orm.Debug = true
	services.Task()
	beego.Run()
}

func init() {
	beego.AddFuncMap("idCradDispose", utils.IdCradDispose)
	beego.AddFuncMap("addOne", utils.AddOne)
	beego.AddFuncMap("length", controllers.Length)
	beego.AddFuncMap("getOperator", utils.GetOperator)
	beego.AddFuncMap("getPercentage", utils.GetOperatorPercentage)
	beego.AddFuncMap("accountDispose", utils.AccountDispose)
	beego.AddFuncMap("formatTimeToString", utils.FormatTimeToString)

}
