package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"beeBlog2/models"
	"fmt"
	"path"
	"time"
	"math"
)
var OneUser string
type MainController struct {
	beego.Controller
}
//注册页面
func (c *MainController) Get() {
	//ok := false
	//c.Data["ok"]=ok
	c.TplName = "register.html"
}
func (c *MainController) Post() {
	//初始化user
	var user models.User
	//获取数据
	userName := c.GetString("userName")
	passWord := c.GetString("password")
	//user赋值
	user.Name =userName
	user.Pwd = passWord
	//NewOrm
	o :=orm.NewOrm()
	//判断是否注册过
	err := o.Read(&user,"Name","Pwd")
	if err==nil {
		ok :=false
		c.Data["ok"]=ok
		beego.Info("注册失败")
		c.TplName="register.html"
		return
	}
	//注册信息写入数据库
	_,err =o.Insert(&user)
	if err != nil {
		beego.Info("注册失败")
		c.Redirect("/reg",302)
		return
	}
	c.Redirect("/login",302)
}
//登录页面
func (c *MainController) ShowLogin()  {
	userName := c.Ctx.GetCookie("userName")
	pwd :=c.Ctx.GetCookie("pwd")
	if userName != ""{
		c.Data["userName"] = userName
		c.Data["pwd"]=pwd
		c.Data["checked"] = "checked"
		}
	c.TplName="login.html"
}
func (c *MainController) HandleLogin() {
	//拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("password")
	//判断
	if userName == ""|| pwd == "" {
		beego.Info("输入数据不合法")
		c.TplName="login.html"
		return
	}
	//查询账户和密码是否正确
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Pwd = pwd
	err := o.Read(&user,"Name","Pwd")
	if err != nil {
		beego.Info("查询失败")
		c.Redirect("/login",302)
		return
	}
	//记住密码设置cookie  200/-1
	remember := c.GetString("remember")
	if remember=="on"{
		c.Ctx.SetCookie("userName",userName,200)
		c.Ctx.SetCookie("pwd",pwd,200)
	}else {
		c.Ctx.SetCookie("userName",userName,-1)
		c.Ctx.SetCookie("pwd",pwd,-1)
	}
	//设置session
	c.SetSession("userName",userName)
	OneUser=userName
	c.Redirect("/article/index",302)
}

//显示主页
func (c *MainController)ShowIndex() {
	//newOrm
	o := orm.NewOrm()
	//获取显示数据类型id
	id,_ := c.GetInt("select")
	//查询文章数据
	qs := o.QueryTable("Article")
	//初始化文章
	var articles []models.Article
	var count int64
	//qs 筛选有文章类型字段的文章
	count,err := qs.RelatedSel("ArticleType").Count()
	if id != 0 {
		//筛选文章字段类型id=获得id的数量
		count,err = qs.RelatedSel("ArticleType").Filter("ArticleType__Id",id).Count()
	}
	if err != nil{
		beego.Info("查询错误")
		return
	}
	//分页
	pageSize := 2
	pageCount := math.Ceil(float64(count)/float64(pageSize))
	//获得显示第几页
	pageIndex,err := c.GetInt("pageIndex")
	if err != nil {
		pageIndex=1
	}
	//查询要求显示内容
	start:= pageSize*(pageIndex -1)
	if id == 0 {
		//id=0显示所有
		qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)
	}else{
		qs.Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__Id",id).All(&articles)
	}
	//优化首末页id
	FirstPage := false
	if pageIndex==1 {
		FirstPage=true
	}
	LastPage := false
	if pageIndex==int(pageCount) {
		LastPage=true
	}
	//显示主页的下拉选框
	var articleType []models.ArticleType
	_,err = o.QueryTable("ArticleType").All(&articleType)
	if err != nil {
		beego.Info("获取类型错误")
		return
	}

	//传递
	c.Data["oneuser"]=OneUser
	c.Data["articleType"]=articleType
	c.Data["FirstPage"]=FirstPage
	c.Data["LastPage"]=LastPage
	c.Data["pageIndex"]=pageIndex
	c.Data["pageCount"]=pageCount
	c.Data["count"]=count
	c.Data["articles"] = articles
	c.Data["typeid"] = id

	c.Layout ="layout.html"
	c.TplName = "index.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["contentHead"]="Hindex.html"
	}

//文章添加
func (c *MainController) ShowAdd()  {
	//neworm
	o := orm.NewOrm()
	//初始化文章类型切片
	var artiTypes  []models.ArticleType
	//查询数据库所有文章类型
	_,err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil {
		beego.Info("获取类型错误")
		return
	}
	//传递文章类型到view
	c.Data["oneuser"]=OneUser
	c.Data["articleType"]=artiTypes
	//网页结构化
	c.Layout ="layout.html"
	c.TplName="add.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["contentHead"]="Hadd.html"
}
func (c *MainController) HandleAdd()  {
	//拿到页面数据
	artiName :=c.GetString("articleName")
	artiContent :=c.GetString("content")
	//文件上传
	f,h,err := c.GetFile("uploadname")
	defer f.Close()
	//限定格式png，jpg
	fileecxt := path.Ext(h.Filename)
	if fileecxt != ".jpg" && fileecxt != ".png"{
		beego.Info("上传文件格式错误")
		return
	}
	//限制大小
	if h.Size > 400000 {
		beego.Info("上传文件过大")
		return
	}
	//对文件重新命名，防止重复
	filename :=time.Now().Format("2006-01-02") +fileecxt
	if err != nil {
		fmt.Println("getFile err=",err)
	}else {
		c.SaveToFile("uploadname","./static/img/"+filename)
	}
	if artiName == ""||artiContent =="" {
		beego.Info("添加文章数据错误")
		return
	}
	//获取文件类型通过页面下拉id传递
	var artiType models.ArticleType
	id,err := c.GetInt("select")
	if err != nil {
		beego.Info("err=",err)
	}
	o := orm.NewOrm()
	artiType = models.ArticleType{Id:id}
	o.Read(&artiType,"id")
	//赋值arti
	arti := models.Article{}
	arti.ArticleType=&artiType
	arti.Artiname = artiName
	arti.Acontent = artiContent
	arti.Aimg = filename
	arti.Atime= time.Now()
	//插入数据库
	_,err =o.Insert(&arti)
	if err!=nil {
		fmt.Println("插入数据库失败")
		return
	}

	c.Layout ="layout.html"
	c.Redirect("/article/index",302)
}

//显示详情
func (c *MainController)ShowContent()  {
	//获取查看人name
	userName := c.GetSession("userName")
	if userName == nil {
		c.Redirect("/login", 302)
		return
	}
	//获取文章id
	id,err := c.GetInt("id")
	if err != nil {
		beego.Info("获取文章Id错误",err)
		return
	}
	o:=orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询失败")
		return
	}
	//浏览量加一
	arti.Acount +=1
	//多表查询
	m2m := o.QueryM2M(&arti,"User")
	//赋值user为当前查看文章的人
	user := models.User{Name:userName.(string)}
	err =o.Read(&user,"Name")
	if err != nil {
		beego.Info("err=",err)
	}
	//添加数据
	m2m.Add(&user)
	o.Update(&arti)
	//初始化查看过用户
	var users  []models.User
	//查询
	o.QueryTable("User").Filter("Article__Article__id",id).Distinct().All(&users)
	//传递
	c.Data["oneuser"]=OneUser
	c.Data["users"]=users
	c.Data["article"]=arti

	c.Layout ="layout.html"
	c.TplName="content.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["contentHead"]="Hcontent.html"

}

//更新文章
func (c *MainController) ShowUpdate()  {
	//获取数据
	id ,err := c.GetInt("id")
	if err != nil {
		beego.Info("获取文章Id错误",err)
		return
	}
	o:=orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询失败",err)
		return
	}
	//传递数据给视图
	c.Data["oneuser"]=OneUser
	c.Data["article"]=arti

	c.Layout ="layout.html"
	c.TplName="update.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["contentHead"]="Hupdate.html"
}
func (c *MainController)HandleUpdate()  {

	id,_:=c.GetInt("id")
	artiname := c.GetString("articleName")
	content := c.GetString("content")


	f,h,err := c.GetFile("uploadname")
	if err !=nil {
		beego.Info("文章上传失败",err)
		return
	}else {
		defer f.Close()
	}
	//限定格式png，jpg
	fileecxt := path.Ext(h.Filename)
	if fileecxt != ".jpg" && fileecxt != ".png"{
		beego.Info("上传文件格式错误")
		return
	}
	//限制大小
	if h.Size > 400000 {
		beego.Info("上传文件过大")
		return
	}
	//对文件重新命名，防止重复
	filename :=time.Now().Format("2016-01-02") +fileecxt
	if err != nil {
		fmt.Println("getFile err=",err)
	}else {
		c.SaveToFile("uploadname","./static/img"+filename)
	}
	if artiname == ""||content =="" {
		beego.Info("添加文章数据错误")
		return
	}
	//更新数据
	o := orm.NewOrm()
	arti := models.Article{Id:id}
	err = o.Read(&arti)
	if err != nil {
		beego.Info("查询失败",err)
		return
	}
	arti.Artiname=artiname
	arti.Acontent=content
	arti.Aimg="./static/img/"+filename

	_,err = o.Update(&arti,"ArtiName","Acontent","Aimg")
	if err != nil {
		beego.Info("更新失败")
		return
	}

	c.Layout ="layout.html"
	c.Redirect("/article/index",302)
}

//删除文章
func (c *MainController)HandleDelete()  {
	//获取id
	id,err := c.GetInt("id")
	if err != nil {
		beego.Info("获取Id错误")
		return
	}
	//修改数据库
	o := orm.NewOrm()
	arti:= models.Article{Id:id}

	err = o.Read(&arti)
	if err!=nil {
		beego.Info("")
		return
	}
	o.Delete(&arti)
	//c.LayoutSections["titleContent"]="Hindex.html"
	c.Layout ="layout.html"
	c.Redirect("/article/index",302)
}

//添加文章
func (c *MainController)ShowAddType()  {
	//找到元素
	o := orm.NewOrm()
	var artiType []models.ArticleType
	//查询数据
	_,err := o.QueryTable("ArticleType").All(&artiType)
	if err != nil {
		beego.Info("未找到artiType",err)
	}
	c.Data["oneuser"]=OneUser
	c.Data["articleType"] = artiType
	c.Layout ="layout.html"
	c.TplName="addType.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["contentHead"]="HaddType.html"
}
func (c *MainController)HandleAddType()  {
	//获取数据
	typeName := c.GetString("typeName")
	//判断数据是否合法
	if typeName == ""{
		beego.Info("获取类型错误")
		return
	}
	//更新数据
	o := orm.NewOrm()
	artiType := models.ArticleType{}
	artiType.Tname=typeName
	_,err := o.Insert(&artiType)
	if err != nil {
		beego.Info("插入失败")
		return
	}

	c.Layout ="layout.html"
	c.Redirect("/article/addType",302)
}

//退出登录
func (c *MainController)LogOut()  {
	c.DelSession("userName")
	c.Redirect("/login",302)
}