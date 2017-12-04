package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/tealeg/xlsx"
)

// 根据数据总数和每页显示数量，得到需要的页数
func GetPageCount(count int, size int) (pageCount int, err error) {
	if i := count % size; i != 0 {
		pageCount = count/size + 1
	} else {
		pageCount = count / size
	}
	return
}

func IdCradDispose1(idCard string) string {
	var err error
	var reg *regexp.Regexp
	var placeStr string
	if len(idCard) == 15 {
		reg, err = regexp.Compile("^(\\d{4})(\\d{9})(.*)")
		placeStr = "*********"
	} else if len(idCard) == 18 {
		reg, err = regexp.Compile("^(\\d{4})(\\d{12})(.*)")
		placeStr = "************"
	} else {
		return ""
	}
	if err != nil {
		return ""
	}
	if reg.MatchString(idCard) == true {
		submatch := reg.FindStringSubmatch(idCard)
		return submatch[1] + placeStr + submatch[3]
	}
	return ""
}

//隐藏手机号中间四位
func AccountDispose(account string) string {
	var err error
	var reg *regexp.Regexp
	var placeStr1 string
	if len(account) == 11 {
		reg, err = regexp.Compile("(\\d{3})\\d{4}(\\d{4})")
		placeStr1 = "****"
	} else {
		return ""
	}
	if err != nil {
		return ""
	}
	if reg.MatchString(account) == true {
		submatch := reg.FindStringSubmatch(account)
		return submatch[1] + placeStr1 + submatch[2]
	}
	return ""
}

//序列化
func ToString(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}

func IdCradDispose(idCard string) string {
	var err error
	var reg *regexp.Regexp
	var placeStr1, placeStr2 string
	if len(idCard) == 15 {
		reg, err = regexp.Compile("^(\\d{4})(\\d{2})(\\d{2})(\\d{5})(.*)")
		placeStr1 = "**"
		placeStr2 = "*****"
	} else if len(idCard) == 18 {
		reg, err = regexp.Compile("^(\\d{4})(\\d{2})(\\d{4})(\\d{6})(.*)")
		// reg, err = regexp.Compile("^(\\d{4})(\\d{12})(.*)")
		placeStr1 = "**"
		placeStr2 = "******"
	} else {
		return ""
	}
	if err != nil {
		return ""
	}
	if reg.MatchString(idCard) == true {
		submatch := reg.FindStringSubmatch(idCard)
		return submatch[1] + placeStr1 + submatch[3] + placeStr2 + submatch[5]
	}
	return ""
}

//md5加密
func MD5(data string) string {
	m := md5.Sum([]byte(data))
	return hex.EncodeToString(m[:])
}

func PageCount(count, pagesize int) int {
	if count%pagesize > 0 {
		return count/pagesize + 1
	} else {
		return count / pagesize
	}
}

func StartIndex(page, pagesize int) int {
	if page > 1 {
		return (page - 1) * pagesize
	}
	return 0
}

func Reduction(x, y float64) float64 {
	return x - y
}

//截取小数点后几位
func SubFloatToString(f float64, m int) string {
	n := strconv.FormatFloat(f, 'f', -1, 64)
	if n == "" {
		return ""
	}
	if m >= len(n) {
		return n
	}
	newn := strings.Split(n, ".")
	if m == 0 {
		return newn[0]
	}
	if len(newn) < 2 || m >= len(newn[1]) {
		return n
	}
	return newn[0] + "." + newn[1][:m]
}

//截取小数点后几位
func SubFloatToFloat(f float64, m int) float64 {
	newn := SubFloatToString(f, m)
	newf, _ := strconv.ParseFloat(newn, 64)
	return newf
}

//float64转字符串 保留2位小数
func Float64ToString(m float64) string {
	return fmt.Sprintf("%.2f", m)
}

func Abs(num int) int {
	if num < 0 {
		num = -num
	}
	return num
}

//模板函数
func AddOne(i int) (l int) {
	l = i + 1
	return
}

func ErrNoRow() string {
	return "<QuerySeter> no row found"
}

func GetSeriesMonth(start, end time.Time, flag int) map[string]interface{} {
	monthMap := make(map[string]interface{})
	if flag == 0 {
		monthMap[start.Format("2006-01-02")] = 0
		for start.Before(end) {
			start = start.AddDate(0, 0, 1)
			addDay := start.Format("2006-01-02")
			monthMap[addDay] = 0
		}
	} else if flag == 1 {
		monthMap[start.Format("2006-01-02")] = float32(0.0)
		for start.Before(end) {
			start = start.AddDate(0, 0, 1)
			addDay := start.Format("2006-01-02")
			monthMap[addDay] = float32(0.0)
		}
	}
	return monthMap
}

func GetSeriesDay(startDate, endDate string, flag int) map[string]interface{} {
	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)
	days := endTime.Sub(startTime).Hours()/24 + 1
	dayNum := int(days)
	if dayNum%7 != 0 {
		dayNum = (dayNum/7 + 1) * 7
	}
	dayMap := make(map[string]interface{})
	if flag == 0 {
		for i := 0; i < dayNum; i++ {
			addDay := startTime.AddDate(0, 0, i).Format("2006-01-02")
			dayMap[addDay] = 0
		}
	} else if flag == 1 {
		for i := 0; i < dayNum; i++ {
			addDay := startTime.AddDate(0, 0, i).Format("2006-01-02")
			dayMap[addDay] = float32(0.0)
		}
	}
	return dayMap
}

func GetSeriesMonths(start, end time.Time) map[string]int {
	monthMap := make(map[string]int)
	monthMap[start.Format("2006-01-02")] = 0
	for start.Before(end) {
		start = start.AddDate(0, 0, 1)
		addDay := start.Format("2006-01-02")
		monthMap[addDay] = 0
	}
	return monthMap
}

func GetSeriesWeek(startDate, endDate string) map[string]int {
	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)
	days := endTime.Sub(startTime).Hours()/24 + 1
	dayNum := int(days)
	if dayNum%7 != 0 {
		dayNum = (dayNum/7 + 1) * 7
	}
	dayMap := make(map[string]int)
	for i := 0; i < dayNum; i++ {
		addDay := startTime.AddDate(0, 0, i).Format("2006-01-02")
		dayMap[addDay] = 0
	}
	return dayMap
}

//手机运营商模版函数

//根据手机号返回运营商
func GetOperator(account string) (operator string) {
	a, _ := strconv.Atoi(account)
	switch a {
	case 139, 138, 137, 136, 135, 134, 159, 158, 157, 150, 151, 152, 147, 188, 187, 182, 183, 184, 178:
		return "移动"
	case 130, 131, 132, 156, 155, 186, 185, 145, 176:
		return "联通"
	case 133, 153, 189, 180, 181, 177, 173:
		return "电信"
	default:
		return "未知"
	}
}

//根据传入的值算出百分比---a/(a+b)
func GetOperatorPercentage(a string, b string) string {
	x, _ := strconv.Atoi(a)
	xx, _ := strconv.Atoi(b)
	par := (x * 1000) / (x + xx)
	if par%10 > 5 {
		par += 10
	}
	return ("0." + strconv.Itoa(par/10))
}

// 将时间格式转化为定时器使用的时间格式
func TimeToTaskSpec(datetime time.Time) (spec string) {
	month := strconv.Itoa(int(datetime.Month()))
	day := strconv.Itoa(datetime.Day())
	hour := strconv.Itoa(datetime.Hour())
	minute := strconv.Itoa(datetime.Minute())
	spec = "" + "0 " + minute + " " + hour + " " + day + " " + month + " *"
	return
}

//导出到excal下载
func ExportToExcel(data [][]string) (filename string, err error) {
	//遍历exportToExcel
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//遍历添加数据
	for _, v := range data {
		row = sheet.AddRow()
		for _, t := range v {
			cell = row.AddCell()
			cell.Value = t
		}
	}
	filename = "MyXLSXFile-" + time.Now().Format("2006-01-02 15-04-05") + ".xlsx"
	err = file.Save(filename)
	return filename, err
}

//根据产品导出数据报表
func ExportDataReportOne(data [][]string) (filename string, err error) {
	//遍历exportToExcel
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		return "", err
	}
	style := xlsx.NewStyle()
	style.Alignment.Horizontal = "center"
	style.Alignment.Vertical = "center"
	//遍历添加数据
	for k, v := range data {
		row = sheet.AddRow()
		for k1, t := range v {
			cell = row.AddCell()
			if k == 0 {
				cell.Merge(10, 0)
				cell.SetStyle(style)
			} else {
				sheet.Cols[k1].Width = 17
			}
			cell.Value = t
		}
	}
	filename = "数据报表-" + time.Now().Format("2006-01-02 15-04-05") + ".xlsx"
	err = file.Save(filename)
	return filename, err
}

func SendEmail(title, content, touser string) {
	// quexw@zcmlc.com
	host := "smtp.exmail.qq.com:25"
	to := strings.Split(touser, ";") //收件人  ;号隔开
	content_type := "Content-Type: text/html; charset=UTF-8"
	msg := []byte("To: " + touser + "\r\nFrom: quexw@zcmlc.com>\r\nSubject:" + title + "\r\n" + content_type + "\r\n\r\n" + content)
	err := smtp.SendMail(host, smtp.PlainAuth("", "quexw@zcmlc.com", "Qxw123", "smtp.exmail.qq.com"), "quexw@zcmlc.com", to, []byte(msg))
	if err != nil {
		fmt.Println(err)
	}
}

//判断月份的最大值
func MaxDay(date string) (maxDate string, ok bool) {
	ok = false
	if date == "" {
		return
	}
	d := strings.Split(date, "-")
	if len(d) > 1 {
		year, err := strconv.Atoi(d[0])
		if err != nil {
			return
		}
		month, err := strconv.Atoi(d[1])
		if err != nil {
			return
		}
		var day int
		switch month {
		case 1, 3, 5, 7, 8, 10, 12:
			day = 31
		case 2:
			if year%4 == 0 && year%100 != 0 || year%400 == 0 {
				day = 29
			} else {
				day = 28
			}
		case 4, 6, 9, 11:
			day = 30
		}
		ok = true
		maxDate = strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + strconv.Itoa(day)
		return
	} else {
		return
	}
	return
}

func SendEmailLwc(title, ip, requestStr, router, content, err, touser string) {
	if beego.AppConfig.String("run_mode") != "debug" {
		now := time.Now().Format("2006-01-02 15:04:05")
		host := "smtp.mxhichina.com:25"
		to := strings.Split(touser, ";") //收件人  ;号隔开
		content_type := "Content-Type: text/html; charset=UTF-8"
		msg := []byte("To: " + touser + "\r\nFrom: lwc@zcmlc.com>\r\nSubject:" + now + "," + title + "\r\n" + content_type + "\r\n\r\n" + "时间:" + now + "<br>IP:" + ip + "<br>请求参数:" + requestStr + "<br>路由:" + router + "<br>内容:" + content + "<br>错误信息:" + err)
		err := smtp.SendMail(host, smtp.PlainAuth("", "lwc@zcmlc.com", "dachang1234!", "smtp.mxhichina.com"), "lwc@zcmlc.com", to, []byte(msg))
		if err != nil {
			fmt.Println("344", err)
		}
	}
}

//格式化时间
func FormatTimeToString(t time.Time) string {
	t1 := t.Format(FormatDateTime)
	switch t1 {
	case "0001-01-01 00:00:00":
		return "-"
	default:
		return t1
	}
}

func SupportLoanExportToExcel(data [][]string, colWidth []float64, fileName string) (filename string, err error) {
	//遍历exportToExcel
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	style := xlsx.NewStyle()
	style.Alignment.Vertical = "center"
	style.Alignment.Horizontal = "center"
	style.Font.Size = 16

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	//遍历添加数据
	for _, v := range data {
		row = sheet.AddRow()
		for k, t := range v {
			cell = row.AddCell()
			cell.SetStyle(style)
			sheet.Cols[k].Width = colWidth[k]
			cell.Value = t
		}
	}
	filename = fileName + ".xlsx"
	err = file.Save(filename)
	return filename, err
}

//获取今天剩余秒数
func GetTodayLastSecond() time.Duration {
	today := time.Now().Format(FormatDate) + " 23:59:59"
	end, _ := time.ParseInLocation(FormatDateTime, today, time.Local)
	return time.Duration(end.Unix()-time.Now().Local().Unix()) * time.Second
}
