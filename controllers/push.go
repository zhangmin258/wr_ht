package controllers

import (
	"time"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"

	"github.com/astaxie/beego/toolbox"
)

//推送模块
type PushController struct {
	BaseController
}

// 跳转到推送页面
func (c *PushController) MessagePush() {
	c.IsNeedTemplate()
	c.TplName = "show-management/push_management.html"
}

// 执行推送任务并保存推送记录--根据用户推送
func (c *PushController) SendPushAndSaveRecord() {
	defer c.ServeJSON()
	//封装推送参数
	var jpushRecord models.JpushRecord
	err := c.ParseForm(&jpushRecord)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "表单解析失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "表单解析失败！"}
		return
	}
	jpushRecord.CreateTime = time.Now()
	// 获取是否定时字段
	isTimePush, err := c.GetBool("IsTimePush")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "IsTimePush解析失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "IsTimePush解析失败！"}
	}
	if isTimePush { // 定时执行推送
		PushTime := c.GetString("PushTimeForm") // 获取定时时间
		if PushTime != "" {
			jpushRecord.PushTime, err = time.Parse(utils.FormatDateTime, PushTime)
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "时间解析失败！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "时间解析失败！"}
				return
			}
		}
		var spec = utils.TimeToTaskSpec(jpushRecord.PushTime) //格式化时间
		if spec == "" {
			c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "请输入正确的推送时间！"}
			return
		}
		//定时任务
		taskMission := toolbox.NewTask("DelayPush", spec, func() error {
			go models.StartPush(&jpushRecord)
			return nil
		})
		toolbox.StopTask()
		toolbox.AddTask("DelayPush", taskMission)
		toolbox.StartTask()
	} else { // 立刻执行推送
		jpushRecord.PushTime = time.Now()
		err = models.StartPush(&jpushRecord)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "推送失败！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "推送失败！"}
			return
		}
	}
	//给用户添加消息
	message := &models.UsersMessage{
		ProductId:   jpushRecord.ProductId, // 产品id
		CreatedUser: c.User.Id,             //创建者
		Content:     jpushRecord.Content,   //内容
		UrlTitle:    jpushRecord.UrlTitle,  //H5标题
		UrlLink:     jpushRecord.Url,       //H5链接
		CreateTime:  time.Now(),            //创建时间
	}
	switch jpushRecord.JumpType { //消息类型
	case 2:
		message.MsgType = 1
	case 3:
		message.MsgType = 5
	}
	go models.AddMessage(message)
	//保存操作记录
	err = models.SavePushRecord(&jpushRecord)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存数据失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "保存数据失败！"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": "200", "msg": "推送成功！"}
	return
}

//获取推送历史列表
func (c *PushController) GetPushRecordList() {
	c.IsNeedTemplate()
	//读取分页信息`
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	params := []string{}
	//推送文本
	if content := c.GetString("Content"); content != "" {
		condition += `and Content like ?`
		params = append(params, "%"+content+"%")
	}
	//分页查询
	pushRecordList, err := models.GetPushRecordList(condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE), utils.PAGE_SIZE)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取推送历史记录失败！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetPushRecordListCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取推送历史总条数失败！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取需要的页数失败！", err.Error(), c.Ctx.Input)
	}
	c.Data["pushRecordList"] = pushRecordList
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.TplName = "show-management/push_history.html"

}

// 获取推送详情
func (c *PushController) GetPushInfo() {
	c.IsNeedTemplate()
	pushId, _ := c.GetInt("id")
	pushInfo, err := models.GetPushInfo(pushId)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取push信息失败！", err.Error(), c.Ctx.Input)
	}
	c.Data["pushInfo"] = pushInfo
	c.TplName = "show-management/check_push.html"
}

// TODO
//发送推送的方法
func (c *PushController) Test() {

	jpushid := "18171adc0337b394eec" //用户
	key := "a8888729bb9e4b68b957a903"
	secret := "b709dd50667cfd0d9639052a"
	user_msg := "内容内容内容内容内容内容内容" //内容
	title := "标题标题标题"            //标题
	urlTitle := "百度"             //标题
	url := "http://www.baidu.com/"
	typeInt := 3
	productIdInt := 8
	test := false

	utils.InitJPush(key, secret) //初始化

	//设置额外参数
	var extraData = make(map[string]interface{}, 0)
	var extraDataMap = make(map[string]interface{}, 0)
	var data = make(map[string]interface{}, 0)
	extraDataMap["Type"] = typeInt //类型5
	if typeInt == 2 {
		data["ProductId"] = productIdInt //产品id

	}
	if typeInt == 3 {
		data["url"] = url        //url
		data["title"] = urlTitle //标题
	}
	extraDataMap["Params"] = data
	extraData["ExtraData"] = extraDataMap
	j := utils.NewJPush(utils.Platform_ALL, utils.Audience_ID)
	j.SetApns(test)
	j.PushMessageWithExtra(extraData, title, user_msg, jpushid)
}

// 执行推送任务并保存推送记录--根据设备推送
func (c *PushController) SendPushAndSaveRecord2() {
	defer c.ServeJSON()
	//封装推送参数
	var jpushRecord models.JpushRecord
	err := c.ParseForm(&jpushRecord)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "表单解析失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "表单解析失败！"}
		return
	}
	jpushRecord.CreateTime = time.Now()
	// 获取是否定时字段
	isTimePush, err := c.GetBool("IsTimePush")
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "IsTimePush解析失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "IsTimePush解析失败！"}
	}
	if isTimePush { // 定时执行推送
		PushTime := c.GetString("PushTimeForm") // 获取定时时间
		if PushTime != "" {
			jpushRecord.PushTime, err = time.Parse(utils.FormatDateTime, PushTime)
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "时间解析失败！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "时间解析失败！"}
				return
			}
		}
		var spec = utils.TimeToTaskSpec(jpushRecord.PushTime) //格式化时间
		if spec == "" {
			c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "请输入正确的推送时间！"}
			return
		}
		//定时任务
		taskMission := toolbox.NewTask("DelayPush", spec, func() error {
			go models.StartPush(&jpushRecord)
			return nil
		})
		toolbox.AddTask("DelayPush", taskMission)
		toolbox.StartTask()
	} else { // 立刻执行推送
		jpushRecord.PushTime = time.Now()
		err = models.StartPush(&jpushRecord)
		if err != nil {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "推送失败！", err.Error(), c.Ctx.Input)
			c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "推送失败！"}
			return
		}
	}
	//给用户添加消息
	message := &models.UsersMessage{
		Content:    jpushRecord.Content,
		CreateTime: time.Now(),
	}
	go models.AddMessage(message)
	//保存操作记录
	err = models.SavePushRecord(&jpushRecord)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "保存数据失败！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": "400", "msg": "保存数据失败！"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": "200", "msg": "推送成功！"}
	return
}
