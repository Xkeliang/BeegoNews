package models

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/astaxie/beego/orm"
)
//用户信息
type User struct {
	Id int `orm:"pk;auto"`
	Name string
	Pwd string

	Article []*Article `orm:"rel(m2m)"`
}
//文章内容
type Article struct {
	Id int  `orm:"pk;auto"`  //主键，自动增长
	Artiname string `orm:"size(20)"`   //ArtiName长度为20
	Atime time.Time   `orm:"auto_now"`//
	Acount  int			`orm:"default(0);null"`
	Acontent string
	Aimg string
	ArticleType *ArticleType `orm:"rel(fk)"`

	User []*User `orm:"reverse(many)"`
}
//文章类型
type ArticleType struct {
	Id int `orm:"pk;auto"`
	Tname string
	Articles []*Article   `orm:"reverse(many)"`
}


//注册数据表（驱动数据库数据模型）
func init()  {
	//orm.RegisterDriver("mysql", orm.MySQL)
	orm.RegisterDataBase("default","mysql","root:root@/test1?charset=utf8")
	orm.RegisterModel(new(User),new(Article),new(ArticleType))
	orm.RunSyncdb("default",false,true)
}