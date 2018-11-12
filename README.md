//beego开发文章发布系统

1、用户注册和登录
    注册登录页面
    数据库生成-model
    注册post请求  -》获取页面数据
    路由处理
    beego处理判断 -》mysql数据库
    返回登录
    
    
2、添加文章
    设置表
    设计表
    	`orm:"pk;aoto"`
    	`orm:"size(20)"`
    	`orm:"auto_now"`,`orm:"auto_now_add"`
    	`orm:"default(0);null"`
    	`orm:"unique"`
    添加文章路由
        控制器
            获取页面数据
                图片：f,h.err :=c.GetFile()
                c.SaveToFile("from","to")
                限制文件格式为图片，判断后缀
                fileext := path.Ext(h.Filename)
                限制文件大小
                重新命名
        插入数据库
            添加文章页面
            form表单允许提交文件，表单属性enctype="multipart/from-data"
            
            
3、显示主页
    o.QueryTable("").All(&article)  //查询表所有内容
    view设计
        {{range .}}
        {{.}}
        {{end}}
        /content?id={{.Id}}
4、详情页面
    view
    路由
    控制获取view传递id，查询数据库，获取数据
    数据传递给view
    
    
5、编辑页面、
    orm操作数据库更新
    
    
6、删除页面
    orm操作数据库更新
    删除确认功能：
    
    
7、分页
        查询多少数据
        共有多少页  //Math.Ceil        Math.floor
        首末页              //qs.Limit(pageSize,start)
        上一页下一页
            beego模板函数  //beego.AddFuncMap()  {{.N | func }}
            优化
            {{if compare .FirstPage true}}
            
            
            
8、添加分类
    数据表
        `orm:"rel(fk)"`      `orm:"reverse(many)"`
        `orm:"rel(m2m)"`     `orm:"reverse(many)"`
    显示类别
        下拉框与后台数据
            选择的Id与类别Id比较
            选择单击触发submit
            qs.Limit().RelatedSel()Filter()All()   //RelatedSel关系查询，参数使用 expr
            
            
9、cookie
     c.Ctx.SetCookie("key","value",time)
     c.Ctx.GetCookie("key")
    userName := c.Ctx.GetCookie("userName")
    
    

10、session
    开启session
        配置文件  session=true
    c.SetSession("key","value")
    c.GetSession()
    c.DelSession()



11、退出



12、过滤路由
        beego.InsertFilter(pattern string, position int, filter FilterFunc, params ...bool)
        例：beego.InsertFilter("/article/*",beego.BeforeRouter,beforExecFunc)
        pattern 路由规则，可以根据一定的规则进行路由，如果你全匹配可以用 *
        position 执行 Filter 的地方，五个固定参数如下，分别表示不同的执行过程
            BeforeStatic 静态地址之前
            BeforeRouter 寻找路由之前
            BeforeExec 找到路由之后，开始执行相应的 Controller 之前
            AfterExec 执行完 Controller 逻辑之后执行的过滤器
            FinishRouter 执行完逻辑之后执行的过滤器
        filter filter 函数 type FilterFunc func(*context.Context)
        例：var beforExecFunc = func(ctx *context.Context) {
          		var userName = ctx.Input.Session("userName")
          	if userName  == nil{
          		ctx.Redirect(302,"/login")
          		return
          	}
          }



13、视图布局
    beego 支持 layout 设计，例如你在管理系统中，整个管理界面是固定的，只会变化中间的部分，那么你可以通过如下的设置：
    c.Layout="layout.html"
    在 layout.html 中你必须设置如下的变量
    {{.LayoutContent}}
    LayoutSection
            对于一个复杂的 LayoutContent，其中可能包括有javascript脚本、CSS 引用等，根据惯例，
            通常 css 会放到 Head 元素中，javascript 脚本需要放到 body 元素的末尾，而其它内容
            则根据需要放在合适的位置。在 Layout 页中仅有一个 LayoutContent 是不够的。所以在 
            Controller 中增加了一个 LayoutSections属性，可以允许 Layout 页中设置多个 section，
            然后每个 section 可以分别包含各自的子模板页。
  
  
  
14、细化
    浏览次数
    浏览人
        多对多查询
            m2m := o.QueryM2M(&arti,"User")
            m2m.Add(&user)
            o.QueryTable("User").Filter("Article__Article__id",id).Distinct().All(&users)  //去重
