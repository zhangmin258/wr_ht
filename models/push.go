package models

import (
	"time"
	"wr_v1/utils"

	"github.com/astaxie/beego/orm"
)

type JpushRecord struct {
	Id               int       `orm:"pk"` //id
	JumpType         int       //跳转位置:1.订单详情 2.h5 3.产品详情页
	Title            string    //标题
	UrlTitle         string    //跳转页标题
	Url              string    //跳转链接
	ProductId        int       //产品id
	Content          string    //内容
	PushTime         time.Time //推送时间
	PushLocationType string    //推送条件：1全部 2只贷给 3不贷给
	PushLocation     string    //推送条件：区域
	PushSys          int       //推送条件：1.全部：2.安卓 3.IOS
	PushMinZm        int       `orm:"column(min_zm)"` //推送条件：最小芝麻信用
	PushMaxZm        int       `orm:"column(max_zm)"` //推送条件：最大芝麻信用
	PushUserLable    string    //用户标签
	CreateTime       time.Time //创建时间
}

type PushHistory struct {
	Id            int       //id
	Title         string    //标题
	Content       string    //内容
	JumpType      int       //跳转位置: 1.h5 2.产品详情页
	ProductId     int       //产品id
	Url           string    //跳转链接
	PushTime      time.Time //推送时间
	PushLocation  string    //推送时间
	PushSys       int       //推送系统：1安卓 2IOS
	PushZm        int       //推送系统：1安卓 2IOS
	PushUserLable string    //用户标签
}

type PackageJpush struct {
	Id            int //id
	PackageId     int //包id
	IosKey        string
	IosSecret     string
	AndroidKey    string
	AndroidSecret string
	Source        string
}

type UserPushInfo struct {
	Id      int    //id
	Jpushid string //极光推送id
}

type JpushRecordShow struct {
	Id               int       `orm:"pk"` //id
	Name             string    // 产品名称
	JumpType         int       //跳转位置:1.订单详情 2.h5 3.产品详情页
	Title            string    //标题
	UrlTitle         string    //跳转页标题
	Url              string    //跳转链接
	ProductId        int       //产品id
	Content          string    //内容
	PushTime         time.Time //推送时间
	PushLocationType int64     //推送条件：1全部 2只贷给 3不贷给
	PushLocation     string    //推送条件：区域
	PushSys          int       //推送条件：1.全部：2.安卓 3.IOS
	PushMinZm        int       `orm:"column(min_zm)"` //推送条件：最小芝麻信用
	PushMaxZm        int       `orm:"column(max_zm)"` //推送条件：最大芝麻信用
	PushUserLable    string    //用户标签
	CreateTime       time.Time //创建时间
}

// 用户消息
type UsersMessage struct {
	Id          int       `orm:"pk"` //
	Uid         int       // 用户ID
	Title       string    // 标题
	Content     string    // 内容
	ProductId   int       // 该消息对应的产品id
	MsgType     int       // 该消息的类型：1.新口子推荐。2.审核未通过。3.审核通过。4.未放款成功。5.跳转url 。6注册成功信息
	IsRead      int       // 是否已读0未读1已读
	CreateTime  time.Time // 创建时间
	CreatedUser int       // 创建人
	Logo        string    // logo图片地址
	UrlLink     string    // 跳转的链接
	UrlTitle    string    //跳转的标题
}

//保存推送记录
func SavePushRecord(jpushRecord *JpushRecord) (err error) {
	_, err = orm.NewOrm().Insert(jpushRecord)
	return
}

//获取推送历史记录
func GetPushRecordList(condition string, params []string, begin, size int) (pushRecordList []*JpushRecord, err error) {
	sql := `SELECT 	id, title, content, jump_type, product_id, url_title, url, push_time, push_location_type,push_location, push_sys, min_zm, max_zm, push_user_lable, create_time FROM  jpush_record where 1=1 `
	if condition != "" {
		sql += condition
	}
	sql += " order by create_time DESC limit ?, ?"
	_, err = orm.NewOrm().Raw(sql, params, begin, size).QueryRows(&pushRecordList)
	return
}

//根据id获取push信息
func GetPushInfo(pushId int) (jpushRecord *JpushRecordShow, err error) {
	sql := `SELECT jr.title,jr.content,jr.jump_type,jr.url_title,jr.url,jr.push_time,jr.push_location_type,jr.push_location,
	 jr.push_sys,jr.min_zm,jr.max_zm,jr.push_user_lable,p.name
	 FROM jpush_record jr LEFT JOIN product p ON jr.product_id=p.id WHERE jr.id =?`
	err = orm.NewOrm().Raw(sql, pushId).QueryRow(&jpushRecord)
	return
}

// 获取推送历史总条数
func GetPushRecordListCount(condition string, params []string) (count int, err error) {
	sql := `SELECT COUNT(1) FROM  jpush_record where 1=1 `
	if condition != "" {
		sql += condition
	}
	err = orm.NewOrm().Raw(sql, params).QueryRow(&count)
	return
}

// 获取所有的包信息
func GetAllPackgeInfo() (packageJpushList []*PackageJpush, err error) {
	o := orm.NewOrm()
	sql := `SELECT package_id, ios_key, ios_secret, android_key, android_secret,source FROM package_jpush `
	_, err = o.Raw(sql).QueryRows(&packageJpushList)
	return
}

// 根据包Id和系统 分批获取对应所有用户Jpushid
func GetJpushIdByPackgeAndSysStep(packageId, app, start, step int) (pushIdList []string, hasNext bool, err error) {
	hasNext = true
	o := orm.NewOrm()
	sql := `SELECT jpushid FROM users WHERE jpushid!="" AND app=? AND register_source=? LIMIT ?,? `
	l, err := o.Raw(sql, app, packageId, start, step).QueryRows(&pushIdList)
	if int(l) < step {
		hasNext = false
	}
	return
}

// 根据包Id和系统 分批获取对应所有用户id
func GetIdByPackgeAndSysStep(packageId, app, start, step int) (idList []int, hasNext bool, err error) {
	hasNext = true
	o := orm.NewOrm()
	sql := `SELECT id FROM users WHERE jpushid!="" AND app=? AND register_source=? LIMIT ?,? `
	l, err := o.Raw(sql, app, packageId, start, step).QueryRows(&idList)
	if int(l) < step {
		hasNext = false
	}
	return
}

// 给每个用户添加消息
func AddMessage(usersMessage *UsersMessage) (err error) {
	IdList := []string{}
	o := orm.NewOrm()
	//先查询所有用户的id
	sql := `SELECT id FROM users`
	_, err = o.Raw(sql).QueryRows(&IdList)
	if err != nil {
		return
	}
	//每次插入10000条
	count := len(IdList) / utils.MESSAGE_INSERT_SIZE
	if len(IdList)%utils.MESSAGE_INSERT_SIZE != 0 {
		count = count + 1
	}
	if count == 0 {
		count = 1
	}
	//分批次添加message
	for i := 0; i < count; i++ {
		ids := []string{}
		if len(IdList)-i*utils.MESSAGE_INSERT_SIZE >= utils.MESSAGE_INSERT_SIZE { //当前批次大于等于10000个
			ids = IdList[i*utils.MESSAGE_INSERT_SIZE : i*utils.MESSAGE_INSERT_SIZE+utils.MESSAGE_INSERT_SIZE] //截取10000个用户
		} else { //当前批次人数小于10000个
			ids = IdList[i*utils.MESSAGE_INSERT_SIZE:]

		}
		sql = `INSERT INTO users_message ( uid,content, product_id, msg_type, create_time, created_user, logo, url_link, url_title) VALUES (?,?,?,?,?,?,?,?,?)`
		prepare, err := orm.NewOrm().Raw(sql).Prepare()
		if err != nil {
			break
		}
		for l := 0; l < len(ids); l++ {
			_, err := prepare.Exec(ids[l], usersMessage.Content, usersMessage.ProductId,
				usersMessage.MsgType, usersMessage.CreateTime, usersMessage.CreatedUser, utils.WR_ICON, usersMessage.UrlLink, usersMessage.UrlTitle)
			if err != nil {
				break
			}
		}
		prepare.Close()
	}
	return
}

//获取用户列表，并开始推送
func StartPush(jpushRecord *JpushRecord) (err error) {
	var packageJpushList []*PackageJpush
	packageJpushList, err = GetAllPackgeInfo() //获取所有包信息
	if err != nil {
		return
	}
	for _, packageJpush := range packageJpushList { //根据包id分批查询用户
		start1, step1 := 0, 1000
		iosHasNext := true
		for iosHasNext {
			IOSUsersPushId, next, err := GetJpushIdByPackgeAndSysStep(packageJpush.PackageId, 1, start1, step1) //分批查询IOS用户jpushId
			//IOSUsersId, next, err := GetIdByPackgeAndSysStep(packageJpush.PackageId, 1, start1, step1)          //分批查询IOS用户jpushId
			if err != nil && err.Error() != `<QuerySeter> no row found` {
				return err
			}
			if len(IOSUsersPushId) > 0 {
				SendJpush(jpushRecord, packageJpush, IOSUsersPushId, 1) //给查询出来的用户推送
				// fmt.Println("IOS:    ", packageJpush.PackageId, "------", len(IOSUsersPushId), "--", start1, "--", step1, "---", packageJpush.IosKey, "--", packageJpush.IosSecret)
			}
			/*if len(IOSUsersId) > 0 {
				//AddMessage(jpushRecord, IOSUsersId) //给查询出来的用户添加message
			}*/
			start1 += step1
			iosHasNext = next
		}
		start2, step2 := 0, 1000
		androidHasNext := true
		for androidHasNext {
			AndroidUsersPushId, next, err := GetJpushIdByPackgeAndSysStep(packageJpush.PackageId, 2, start2, step2) //分批查询安卓用户
			if err != nil && err.Error() != `<QuerySeter> no row found` {
				return err
			}
			if len(AndroidUsersPushId) > 0 {
				SendJpush(jpushRecord, packageJpush, AndroidUsersPushId, 2) //给查询出来的用户推送
				// fmt.Println("Android:", packageJpush.PackageId, "------", len(AndroidUsersPushId), "--", start2, "--", step2, "---", packageJpush.AndroidKey, "--", packageJpush.AndroidSecret)
			}
			start2 += step2
			androidHasNext = next
		}
	}
	return
}

//开始推送
func StartPush2(jpushRecord *JpushRecord) (err error) {
	return
}

//发送推送的方法--分用户推
func SendJpush(jpushRecord *JpushRecord, packageJpush *PackageJpush, users []string, app int) {
	user_msg := jpushRecord.Content //内容
	jpushid := users                //用户
	title := jpushRecord.Title      //标题
	key := ""
	secret := ""
	//设置key和secret
	if app == 1 { //ios
		key = packageJpush.IosKey
		secret = packageJpush.IosSecret
	}
	if app == 2 { //android
		key = packageJpush.AndroidKey
		secret = packageJpush.AndroidSecret
	}
	utils.InitJPush(key, secret) //初始化

	//设置额外参数
	var extraData = make(map[string]interface{}, 0)
	var extraDataMap = make(map[string]interface{}, 0)
	var data = make(map[string]interface{}, 0)

	if jpushRecord.JumpType == 3 { //跳转H5
		extraDataMap["Type"] = 3             //类型
		data["url"] = jpushRecord.Url        //url
		data["title"] = jpushRecord.UrlTitle //标题
	}
	if jpushRecord.JumpType == 2 { //跳转产品详情
		extraDataMap["Type"] = 2                  //类型
		data["ProductId"] = jpushRecord.ProductId //产品Id
	}
	extraDataMap["Params"] = data
	extraData["ExtraData"] = extraDataMap
	j := utils.NewJPush(utils.Platform_ALL, utils.Audience_ID)
	j.SetApns(true)
	j.PushMessageWithExtra(extraData, title, user_msg, jpushid...)
}

//发送推送的方法--分设备推
func SendJpush2(jpushRecord *JpushRecord, packageJpush *PackageJpush, users []string, app int) {
	user_msg := jpushRecord.Content //内容
	jpushid := users                //用户
	title := jpushRecord.Title      //标题
	key := ""
	secret := ""
	//设置key和secret
	if app == 1 { //ios
		key = packageJpush.IosKey
		secret = packageJpush.IosSecret
	}
	if app == 2 { //android
		key = packageJpush.AndroidKey
		secret = packageJpush.AndroidSecret
	}
	utils.InitJPush(key, secret) //初始化

	//设置额外参数
	var extraData = make(map[string]interface{}, 0)
	var extraDataMap = make(map[string]interface{}, 0)
	var data = make(map[string]interface{}, 0)

	if jpushRecord.JumpType == 3 { //跳转H5
		extraDataMap["Type"] = 3             //类型
		data["url"] = jpushRecord.Url        //url
		data["title"] = jpushRecord.UrlTitle //标题
	}
	if jpushRecord.JumpType == 2 { //跳转产品详情
		extraDataMap["Type"] = 2                  //类型
		data["ProductId"] = jpushRecord.ProductId //产品Id
	}
	extraDataMap["Params"] = data
	extraData["ExtraData"] = extraDataMap
	j := utils.NewJPush(utils.Platform_ALL, utils.Audience_ID)
	j.SetApns(true)
	j.PushMessageWithExtra(extraData, title, user_msg, jpushid...)
}
