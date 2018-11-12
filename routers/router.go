package routers

import (
	"beeBlog2/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*",beego.BeforeRouter,beforExecFunc)
    beego.Router("/reg",&controllers.MainController{})
    beego.Router("/login",&controllers.MainController{},"Get:ShowLogin;Post:HandleLogin")
	beego.Router("/article/index", &controllers.MainController{},"get:ShowIndex")
	beego.Router("/article/addArticle", &controllers.MainController{},"get:ShowAdd;Post:HandleAdd")
	beego.Router("/article/content", &controllers.MainController{},"get:ShowContent")
	beego.Router("/article/update", &controllers.MainController{},"get:ShowUpdate;Post:HandleUpdate")
	beego.Router("/article/delete", &controllers.MainController{},"get:HandleDelete")
	beego.Router("/article/addType", &controllers.MainController{},"get:ShowAddType;Post:HandleAddType")
	beego.Router("/logout",&controllers.MainController{},"Get:LogOut")
}

var beforExecFunc = func(ctx *context.Context) {
		var userName = ctx.Input.Session("userName")
	if userName  == nil{
		ctx.Redirect(302,"/login")
		return
	}
}