package main

import (
	_ "beeBlog2/routers"
	"github.com/astaxie/beego"
)

func main() {
	//模板函数
	beego.AddFuncMap("addone",showNewAddOne)
	beego.AddFuncMap("prePage",prePage)
	beego.AddFuncMap("nextPage",nextPage)
	beego.Run()
}


func prePage(pageindex int)(preIndex int)  {
	preIndex = pageindex-1
	return
}
func nextPage(pageindex int)(nextIndex int)  {
	nextIndex = pageindex+1
	return
}

func showNewAddOne(num int)(newNum int){
	newNum = num + 1
	return
}