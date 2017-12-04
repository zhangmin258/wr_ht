package controllers

import (
	"strconv"
	"sync"
	"wr_v1/cache"
	"wr_v1/models"
	"wr_v1/utils"

	//"github.com/astaxie/beego"

	"time"
	"wr_v1/services"
)

/*
产品运营信息相关接口
*/
type ControlController struct {
	BaseController
}

//条件分页查询产品运营信息列表
//@router /getProductMange [get]
func (c *ControlController) GetProductMange() {
	c.IsNeedTemplate()
	//读取分页信息
	pageNum, _ := c.GetInt("page", 1)
	if pageNum < 1 {
		pageNum = 1
	}
	condition := ""
	order := ""
	params := []string{}
	//产品名字
	name := c.GetString("proName")
	if name != "" {
		condition += " and p.name like ? OR o.name LIKE ? "
		params = append(params, "%"+name+"%")
		params = append(params, "%"+name+"%")
	}
	//展示位置
	types := c.GetString("locationType")
	if types == "" || types == "3" {
		order = " ORDER BY p.id DESC"
	} else {
		condition += " and p.is_index_show =?"
		params = append(params, types)
		//排序
		if types == "1" { //首页
			order = " ORDER BY p.sort ASC"
		}
		if types == "0" { //贷款页
			order = " ORDER BY p.loan_sort ASC"
		}
		if types == "2" { //大额贷款页
			order = " ORDER BY p.large_loan_sort ASC "
		}
	}

	//查询
	products, err := models.GetProductMange(order, condition, params, utils.StartIndex(pageNum, utils.PAGE_SIZE20), utils.PAGE_SIZE20)
	//fmt.Println(products)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询产品运营信息异常！", err.Error(), c.Ctx.Input)
	}
	count, err := models.GetProductCount(condition, params)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询所有商品数量异常！", err.Error(), c.Ctx.Input)
	}
	pageCount, err := utils.GetPageCount(count, utils.PAGE_SIZE20)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询需要页数异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["products"] = products
	c.Data["pageNum"] = pageNum
	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.TplName = "show-management/product_management.html"

}

//产品管理页面->编辑页面
//@router /editProductMange [get]
func (c *ControlController) EditProductMange() {

	c.IsNeedTemplate()
	var product *models.ProductMange
	var findproaddress models.FindProAddress
	pid, _ := c.GetInt("id")
	product, err := models.ProductMangeById(pid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查找productmange异常！", err.Error(), c.Ctx.Input)
	}
	findproaddress, err = models.ProductAddressById(pid)
	if err != nil && err.Error() != utils.ErrNoRow() {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "根据id查找productaddress异常！", err.Error(), c.Ctx.Input)
	}
	c.Data["findproaddress"] = findproaddress
	c.Data["product"] = product
	c.TplName = "show-management/edit_product_management.html"
}

/*
//产品管理->编辑页面->保存按钮
//@router /saveProductMange [post]
func (c *ControlController) SaveProductMange() {
	defer c.ServeJSON()
	var pro models.ProductMange
	err := c.ParseForm(&pro)
	if err != nil {
		//解析参数异常
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "pro参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "pro参数解析异常!"}
		return
	}
	location := c.GetString("Location")
	i, err := strconv.Atoi(location)
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "pro参数解析异常!"}
	}
	if pro.IsIndexShow == 0 {
		pro.Sort = -1
		pro.LoanSort = i
	}
	if pro.IsIndexShow == 1 {
		pro.Sort = i
		pro.LoanSort = -1
	}
	var productAddress models.ProductAddress
	err = c.ParseForm(&productAddress)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "productAddress参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "productAddress参数解析异常!"}
		return
	}

	//改变product的sort值时 相应的其他的sort也需要重新排列
	//获取改变前的sort值
	proOld, err := models.GetProductMangeById(pro.Id)
	sortOld := proOld.Sort
	condition := "AND sort BETWEEN ? AND ? "
	param := []int{}
	//获取改变之前的sort到改变之后的sort所有数据
	if pro.Sort > sortOld {
		param = append(param, sortOld+1)
		param = append(param, pro.Sort)
	} else {
		param = append(param, pro.Sort)
		param = append(param, sortOld)
	}
	sortList, err := models.SelectProductSort(condition, param)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取sort列表异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取sort列表失败"}
		return
	}

	//获取最大sort值
	maxSort, maxLoanSort, err := models.GetMaxSort()
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取最大sort值异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取sort最大值失败"}
		return
	}
	//改变后的sort比之前的大并且比最大sort值小
	if sortOld < pro.Sort && pro.Sort <= maxSort {
		if sortOld == 0 {
			sortList, err = models.SelectProductSort(condition, []int{pro.Sort, maxSort + 1})
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取sort列表异常！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取sort列表失败"}
				return
			}
			for i := 0; i < len(sortList); i++ {
				sortList[i].Sort = sortList[i].Sort + 1
			}
		} else {
			for i := 0; i < len(sortList); i++ {
				sortList[i].Sort = sortList[i].Sort - 1
			}
		}
	}
	//改变后的sort比之前的大并且比最大sort值大
	if sortOld < pro.Sort && pro.Sort > maxSort {
		if sortOld == 0 {
			pro.Sort = maxSort + 1
		} else {
			for i := 0; i < len(sortList); i++ {
				sortList[i].Sort = sortList[i].Sort - 1
			}
			pro.Sort = maxSort
		}

	}
	//改变后的sort比之前的小
	if sortOld > pro.Sort {
		for i := 0; i < len(sortList); i++ {
			sortList[i].Sort = sortList[i].Sort + 1
		}
	}

	//改变product的loansort值时 相应的其他的loansort也需要重新排列
	//获取改变前的loansort值
	loanSortOld := proOld.LoanSort
	condition = "AND loan_sort BETWEEN ? AND ? "
	paramLoan := []int{}
	//获取改变之前的sort到改变之后的sort所有数据
	if pro.LoanSort > loanSortOld {
		paramLoan = append(paramLoan, loanSortOld+1)
		paramLoan = append(paramLoan, pro.LoanSort)
	} else {
		paramLoan = append(paramLoan, pro.LoanSort)
		paramLoan = append(paramLoan, loanSortOld)
	}
	loanSortList, err := models.SelectProductSort(condition, paramLoan)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取loansort列表异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取loansort列表失败"}
		return
	}

	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取loan_sort最大值异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取loan_sort最大值失败"}
		return
	}
	//改变后的loansort比之前的大并且比最大sort值小
	if loanSortOld < pro.LoanSort && pro.LoanSort <= maxLoanSort {
		if loanSortOld == 0 {
			loanSortList, err = models.SelectProductSort(condition, []int{pro.LoanSort, maxLoanSort + 1})
			if err != nil {
				cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取loan_sort列表异常！", err.Error(), c.Ctx.Input)
				c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取loan_sort列表失败"}
				return
			}
			for i := 0; i < len(loanSortList); i++ {
				loanSortList[i].LoanSort = loanSortList[i].LoanSort + 1
			}
		} else {
			for i := 0; i < len(loanSortList); i++ {
				loanSortList[i].LoanSort = loanSortList[i].LoanSort - 1
			}
		}
	}
	//改变后的loansort比之前的大并且比最大sort值大
	if loanSortOld < pro.LoanSort && pro.LoanSort > maxLoanSort {
		if loanSortOld == 0 {
			pro.LoanSort = maxLoanSort + 1
		} else {
			for i := 0; i < len(loanSortList); i++ {
				loanSortList[i].LoanSort = loanSortList[i].LoanSort - 1
			}
			pro.LoanSort = maxLoanSort
		}
	}
	//改变后的loansort比之前的小
	if loanSortOld > pro.LoanSort {
		for i := 0; i < len(loanSortList); i++ {
			loanSortList[i].LoanSort = loanSortList[i].LoanSort + 1
		}
	}
	err = models.UpdateProductMange(&pro, &productAddress, sortList, loanSortList) //插入或者修改数据
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "插入或者修改数据异常！", err.Error(), c.Ctx.Input)
		//插入申请条件异常
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "插入或者修改数据异常"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
	return
}
*/

//保存编辑过后的产品运营信息
//@router /editproductinfo [post]
func (c *ControlController) EditProduct() {
	var mutex sync.Mutex
	mutex.Lock()
	defer func() {
		c.ServeJSON()
		mutex.Unlock()
	}()
	var pro models.ProductMange
	err := c.ParseForm(&pro) //接收上线状态，展示顺序等信息
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "pro参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "pro参数解析异常!"}
		return
	}
	product, err := models.GetProductMangeInfo(pro.Id)
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取产品展示信息异常！", err.Error(), c.Ctx.Input)
		}
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取产品展示信息异常!"}
		return
	}
	sort, err := models.GetProductSortInfoById(pro.Id) //查询原信息
	if err != nil {
		if err.Error() != utils.ErrNoRow() {
			cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询原数据异常！", err.Error(), c.Ctx.Input)
		}
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "查询原数据异常!"}
	}
	maxSort, maxLoanSort, maxLargeLoanSort, err := models.GetMaxSort() //查询原最大sort和loan_sort
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询最大sort异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "查询最大sort异常!"}
	}
	maxUsedSort, maxUsedLoanSort, maxUsedLargeLoanSort, err := models.GetMaxUsedSortAndLoanSort() //查询使用中的最大的sort和loanSort
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "查询最大sort异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "查询最大sort异常!"}
	}
	bothSortInt, _ := c.GetInt("Location") //获取新的展示位置
	if err != nil {
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "参数解析异常!"}
	}
	/*
		if pro.IsIndexShow == 0 { //贷款页展示赋值
			pro.Sort = -1
			if bothSortInt > maxLoanSort+1 { //当值大于最大值+1时
				if pro.IsUse == 0 { //上线
					bothSortInt = maxUsedLoanSort //上线的最大值+1
				}
				if pro.IsUse == 1 {
					bothSortInt = maxLoanSort
				}
			}
			pro.LoanSort = bothSortInt
		}
		if pro.IsIndexShow == 1 { //首页展示赋值
			if bothSortInt > maxSort+1 {
				if pro.IsUse == 0 {
					bothSortInt = maxUsedSort

				}
				if pro.IsUse == 1 {
					bothSortInt = maxSort

				}
			}
			pro.Sort = bothSortInt
			pro.LoanSort = -1
		}
	*/
	pro.Sort = -1
	pro.LoanSort = -1
	pro.LargeLoanSort = -1
	if pro.IsIndexShow == 1 { //首页赋值
		pro.Sort = bothSortInt
	}
	if pro.IsIndexShow == 0 { //贷款页赋值
		pro.LoanSort = bothSortInt
	}
	if pro.IsIndexShow == 2 { //大额贷款页赋值
		pro.LargeLoanSort = bothSortInt
	}
	oldIsIndexShow := sort.IsIndexShow     //旧的展示位置
	oldSort := sort.Sort                   //旧的首页顺序
	oldLoanSort := sort.LoanSort           //旧的贷款页顺序
	oldLargeLoanSort := sort.LargeLoanSort //旧的大额贷款页顺序
	// oldIsuse := sort.IsUse                 //新的冻结情况
	newIsIndexShow := pro.IsIndexShow     //新的展示位置
	newSort := pro.Sort                   //新的首页顺序
	newLoanSort := pro.LoanSort           //新的贷款页顺序
	newLargeLoanSort := pro.LargeLoanSort //新的大额贷款页顺序
	newIsuse := pro.IsUse                 //新的冻结情况
	if oldIsIndexShow == 2 {              //大额贷款页产品 位置固定
		//若原来有顺序,原位置后的产品顺序依次降一
		if oldLargeLoanSort != 0 && oldLargeLoanSort != -1 {
			models.DownSort(oldLargeLoanSort, maxLargeLoanSort, "largeLoan")
			//获取新位置排序
			if newIsuse == 0 && newLargeLoanSort > maxUsedLargeLoanSort { //上线 并且 超出最大上线产品顺序时 取上线产品的最大顺序
				pro.LargeLoanSort = maxUsedLargeLoanSort
				newLargeLoanSort = maxUsedLargeLoanSort
			} else if newIsuse == 1 && newLargeLoanSort > maxLargeLoanSort { //下线 并且 超出所有大额产品最大顺序时 取所有大额产品的最大顺序
				pro.LargeLoanSort = maxLargeLoanSort
				newLargeLoanSort = maxLargeLoanSort
			}
		} else {
			//获取新位置排序
			if newIsuse == 0 && newLargeLoanSort > maxUsedLargeLoanSort { //上线 并且 超出最大上线产品顺序时 取上线产品的最大顺序
				pro.LargeLoanSort = maxUsedLargeLoanSort + 1
				newLargeLoanSort = maxUsedLargeLoanSort + 1
			} else if newIsuse == 1 && newLargeLoanSort > maxLargeLoanSort { //下线 并且 超出所有大额产品最大顺序时 取所有大额产品的最大顺序
				pro.LargeLoanSort = maxLargeLoanSort + 1
				newLargeLoanSort = maxLargeLoanSort + 1
			}
		}
		//若新顺序不为 0 或 -1,新位置后的产品顺序依次+1
		if newLargeLoanSort != -1 && newLargeLoanSort != 0 {
			models.UpSort(maxLargeLoanSort+1, newLargeLoanSort, "largeLoan") //若新位置排序不为 0 或 -1 大额贷款页所有在新位置后的排序 +1
		}
	} else {
		//除 0 和 -1 以外,原位置后的顺序依次 -1
		if oldIsIndexShow == 1 && oldSort != 0 && oldSort != -1 { //原位置在首页 且原位置排序不为 0 或 -1,则原位置后的产品降序
			models.DownSort(oldSort, maxSort, "index")
		}
		if oldIsIndexShow == 0 && oldLoanSort != 0 && oldLoanSort != -1 { //原位置在贷款页 且原位置排序不为 0 或 -1,则原位置后的产品降序
			models.DownSort(oldLoanSort, maxLoanSort, "loan")
		}
		if newIsIndexShow == 1 { //新位置在首页
			//获取新位置排序
			if newIsuse == 0 && newSort > maxUsedSort { //上线 并且 超出最大上线产品顺序时
				if oldSort != -1 && oldSort != 0 {
					pro.Sort = maxUsedSort
					newSort = maxUsedSort
				} else {
					pro.Sort = maxUsedSort + 1
					newSort = maxUsedSort + 1
				}
			} else if newIsuse == 1 && newSort > maxSort { //下线 并且 超出所有大额产品最大顺序时 取所有大额产品的最大顺序
				if oldSort != -1 && oldSort != 0 { //若原来产品在该页面有排序 取该页上线产品的最大顺序
					pro.Sort = maxSort
					newSort = maxSort
				} else { //若原来产品在该页面没有排序 取新页面上线产品的最大顺序 +1
					pro.Sort = maxSort + 1
					newSort = maxSort + 1
				}
			}
			//若新顺序不为 0 或 -1,新位置后的产品顺序依次 +1
			if newSort != -1 && newSort != 0 {
				models.UpSort(maxSort+1, newSort, "index") //若新位置排序不为 0 或 -1 贷款页所有在新位置后的排序 +1
			}
		}
		if newIsIndexShow == 0 { //新位置在贷款页
			//获取新位置排序
			if newIsuse == 0 && newLoanSort > maxUsedLoanSort { //上线 并且 超出最大上线产品顺序时
				if oldLoanSort != -1 && oldLoanSort != 0 {
					pro.LoanSort = maxUsedLoanSort
					newLoanSort = maxUsedLoanSort
				} else {
					pro.LoanSort = maxUsedLoanSort + 1
					newLoanSort = maxUsedLoanSort + 1
				}
			} else if newIsuse == 1 && newLoanSort > maxLoanSort { //下线 并且 超出所有大额产品最大顺序时 取所有大额产品的最大顺序
				if oldLoanSort != -1 && oldLoanSort != 0 { //若原来产品在该页面有排序 取该页上线产品的最大顺序
					pro.LoanSort = maxLoanSort
					newLoanSort = maxLoanSort
				} else { //若原来产品在该页面没有排序 取新页面上线产品的最大顺序 +1
					pro.LoanSort = maxLoanSort + 1
					newLoanSort = maxLoanSort + 1
				}
			}
			//若新顺序不为 0 或 -1,新位置后的产品顺序依次 +1
			if newLoanSort != -1 && newLoanSort != 0 {
				models.UpSort(maxLoanSort+1, newLoanSort, "loan") //若新位置排序不为 0 或 -1 贷款页所有在新位置后的排序 +1
			}
		}
		/*
			//当新的产品状态 >>>冻结      该产品位置设为贷款页，顺序设为最后一个
			if newIsuse == 1 {
				if oldIsIndexShow == 0 { //原位置在贷款页
					models.DownSort(oldLoanSort, maxLoanSort+1, "loan") //贷款页  >oldLoanSort的-1
					pro.IsIndexShow = 0
					pro.LoanSort = maxLoanSort
				}
				if oldIsIndexShow == 1 { //原位置在首页
					models.DownSort(oldSort, maxSort+1, "index") //首页  >oldSort的-1
					pro.IsIndexShow = 0                          //设为贷款页展示
					pro.LoanSort = maxLoanSort + 1
				}
				pro.Sort = -1
			}
			//当新的产品状态   冻结>>>上线
			if newIsuse == 0 && oldIsuse == 1 {
				//产品展示位置没有改动
				if newIsIndexShow == 0 {
					if newLoanSort > maxUsedLoanSort { //序号大于最大的使用序号
						pro.LoanSort = maxUsedLoanSort + 1 //序号设置为最大使用序号+1
						newLoanSort = pro.LoanSort
					}
				}
				//产品展示位置有改动
				if newIsIndexShow == 1 {
					if newSort > maxUsedSort { //序号大于最大的使用序号
						pro.Sort = maxUsedSort + 1                   //序号设置为最大使用序号+1
						models.UpSort(oldLoanSort, pro.Sort, "loan") //贷款页所有介于 oldLoanSort和newLoanSort的+1
					}
				}
			}
			//当新的产品状态   >>>上线
			if newIsuse == 0 {
				//产品展示位置没有改动
				if oldIsIndexShow == newIsIndexShow {
					//原位置在贷款页
					if oldIsIndexShow == 0 {
						if newLoanSort > oldLoanSort { //顺位排后
							models.DownSort(oldLoanSort, newLoanSort, "loan") //贷款页所有介于 oldLoanSort和newLoanSort的-1
						}
						if newLoanSort < oldLoanSort { //顺位前移
							models.UpSort(oldLoanSort, newLoanSort, "loan") //贷款页所有介于 oldLoanSort和newLoanSort的+1
						}
					}
					//原位置在首页
					if oldIsIndexShow == 1 {
						if newSort > oldSort { //顺位排后
							models.DownSort(oldSort, newSort, "index") //首页所有介于 oldLoanSort和newLoanSort的-1
						}
						if newSort < oldSort { //顺位前移
							models.UpSort(oldSort, newSort, "index") //首页所有介于 oldLoanSort和newLoanSort的+1
						}
					}
				}
				//产品展示位置有改动
				if oldIsIndexShow != newIsIndexShow {
					if oldIsIndexShow == 0 { //贷款页>>>>>首页
						models.DownSort(oldLoanSort, maxLoanSort, "loan") //贷款页所有大于oldLoanSort的-1
						models.UpSort(maxSort+1, newSort, "index")        //首页所有大于newSort的+1
					}
					if oldIsIndexShow == 1 { //首页>>>>>贷款页
						models.DownSort(oldSort, maxSort, "index")        //首页所有大于oldSort的-1
						models.UpSort(maxLoanSort+1, newLoanSort, "loan") //贷款页所有大于newLoanSort的+1
					}
				}
			}
		*/
	}

	if pro.IsUse == 1 {
		pro.RecommendSort = 0
	}
	var useTime int64
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	//查出最近的上线平台不同位置的时间
	con := " AND (is_index_show=? OR loan_sort=? OR sort=? OR large_loan_sort=?) "
	dataTime, err := models.GetNotEquTime(con, pro.Id, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "获取最近的上线平台不同位置的时间异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "获取最近的上线平台不同位置的时间异常!"}
		return
	}
	if dataTime.Before(tm1) {
		useTime = t.Unix() - tm1.Unix()
	} else {
		useTime = t.Unix() - dataTime.Unix()
	}
	if pro.IsUse != product.IsUse {
		if pro.IsUse == 0 { //上线
			useTime = 0
			services.ProductLogs(c.User.Id, pro.Id, 0, 1, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, 0, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为上线", c.Ctx.Input)
			services.ProductLogs(c.User.Id, pro.Id, 0, 1, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为上线", c.Ctx.Input)
		} else if pro.IsUse == 1 { //下线
			services.ProductLogs(c.User.Id, pro.Id, 0, 2, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, useTime, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为下线", c.Ctx.Input)
			services.ProductLogs(c.User.Id, pro.Id, 0, 2, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为下线", c.Ctx.Input)
		}
	}
	if pro.IsPopUp != product.IsPopUp {
		if pro.IsPopUp == 0 { //不弹窗
			services.ProductLogs(c.User.Id, pro.Id, 0, 4, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, 0, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为不弹窗", c.Ctx.Input)
			services.ProductLogs(c.User.Id, pro.Id, 0, 4, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为不弹窗", c.Ctx.Input)
		} else if pro.IsPopUp == 1 { //弹窗
			services.ProductLogs(c.User.Id, pro.Id, 0, 3, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, 0, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为弹窗", c.Ctx.Input)
			services.ProductLogs(c.User.Id, pro.Id, 0, 3, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为弹窗", c.Ctx.Input)
		}
	}

	if pro.IsIndexShow != product.IsIndexShow {
		if pro.IsIndexShow == 0 { //贷款页
			services.ProductLogs(c.User.Id, pro.Id, 0, 6, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, useTime, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为贷款页", c.Ctx.Input)
			services.ProductLogs(c.User.Id, pro.Id, 0, 6, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为贷款页", c.Ctx.Input)
		} else if pro.IsIndexShow == 1 { //首页
			services.ProductLogs(c.User.Id, pro.Id, 0, 5, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, useTime, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为首页", c.Ctx.Input)
			services.ProductLogs(c.User.Id, pro.Id, 0, 5, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为首页", c.Ctx.Input)
		}
	}
	if pro.FullGuide != product.FullGuide { //导量上限
		services.ProductLogs(c.User.Id, pro.Id, 0, 7, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, 0, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改为导量上限值", c.Ctx.Input)
		services.ProductLogs(c.User.Id, pro.Id, 0, 7, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改导量上限值为:"+strconv.Itoa(pro.FullGuide), c.Ctx.Input)
	}
	if pro.Sort != product.Sort { //首页排序
		if product.Sort == 0 {
			services.ProductLogs(c.User.Id, pro.Id, 0, 8, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, 0, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改首页排序", c.Ctx.Input)
		} else {
			services.ProductLogs(c.User.Id, pro.Id, 0, 8, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, useTime, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改首页排序", c.Ctx.Input)
		}
		services.ProductLogs(c.User.Id, pro.Id, 0, 8, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改首页排序值为:"+strconv.Itoa(pro.Sort), c.Ctx.Input)
	}
	if pro.LoanSort != product.LoanSort { //贷款页排序
		if product.LoanSort == 0 {
			services.ProductLogs(c.User.Id, pro.Id, 0, 8, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, 0, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改贷款页排序", c.Ctx.Input)
		} else {
			services.ProductLogs(c.User.Id, pro.Id, 0, 8, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, useTime, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改贷款页排序", c.Ctx.Input)
		}

		services.ProductLogs(c.User.Id, pro.Id, 0, 8, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改贷款页排序值为:"+strconv.Itoa(pro.LoanSort), c.Ctx.Input)
	}
	if pro.RecommendSort != product.RecommendSort { //产品推广
		services.ProductLogs(c.User.Id, pro.Id, 0, 9, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, 0, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改产品推广", c.Ctx.Input)
		services.ProductLogs(c.User.Id, pro.Id, 0, 9, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改产品推广值为:"+strconv.Itoa(pro.RecommendSort), c.Ctx.Input)
	}
	if pro.LargeLoanSort != product.LargeLoanSort { //大额贷款页排序
		if product.LargeLoanSort == 0 {
			services.ProductLogs(c.User.Id, pro.Id, 0, 10, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, 0, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改大额贷款页排序", c.Ctx.Input)
		} else {
			services.ProductLogs(c.User.Id, pro.Id, 0, 10, product.IsIndexShow, product.LoanSort, product.Sort, product.LargeLoanSort, useTime, 0, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改大额贷款页排序", c.Ctx.Input)
		}
		services.ProductLogs(c.User.Id, pro.Id, 0, 10, pro.IsIndexShow, pro.LoanSort, pro.Sort, pro.LargeLoanSort, 0, 1, c.User.Name, c.User.DisplayName, "", "用户操作产品", "修改大额贷款页排序值为:"+strconv.Itoa(pro.LargeLoanSort), c.Ctx.Input)
	}
	if pro.Tag == "" {
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "产品标签不得为空!"}
		return
	}
	err = models.UpdateProductInfo(&pro) //更新产品
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "更新产品信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "更新产品信息异常!"}
		return
	}
	var productAddress models.ProductAddress
	err = c.ParseForm(&productAddress)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "productAddress参数解析异常", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "productAddress参数解析异常!"}
		return
	}
	if productAddress.MaxAge != 0 && productAddress.MaxAge < productAddress.MinAge {
		c.Data["json"] = map[string]interface{}{"ret": 400, "msg": "最大年龄不能小于最小年龄!"}
		return
	}
	err = models.UpdateProductAddress(&productAddress)
	if err != nil {
		cache.RecordLogs(c.User.Id, 0, c.User.Name, c.User.DisplayName, "", "更新产品信息异常！", err.Error(), c.Ctx.Input)
		c.Data["json"] = map[string]interface{}{"ret": 400, "err": err.Error(), "msg": "更新产品信息异常!"}
		return
	}
	c.Data["json"] = map[string]interface{}{"ret": 200}
	return
}
