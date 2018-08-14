# Bingo

Bingo是一款使用`httprouter`作为路由的Web全栈开发框架。

受到`Laravel`的启发，将一些网站开发过程中必备的功能内置到了框架中，开箱即用

我致力于让它有着`Golang`的速度和`Laravel`的优雅

目前正在开发中......

> 开发过程中请勿使用master分支，如有需要，请使用release的版本
> 我们使用分支开发，主干发布的工作流，但是并不保证master一定可用，如发现问题，可提issue
> 如果有想法，欢迎联系我 Email: silsuer.liu@gmail.com

## 安装

-----题外话开始-------------

`Bingo`使用`glide`来管理依赖，但是`glide`在windows10 x64中存在bug，安装依赖时会报错，
网上对于这个bug的解决有好多不同的...我把这个bug改掉后也传到了项目里，如果需要的话，请把
`$GOPATH/src/github/silsuer/bingo`下的`glide.exe` 拷贝到你的项目工程目录下，在终端中即可
使用`glide get` 等命令


-----题外话结束---------

 1. 如果你像我一样正在使用`glide`作为包管理工具，在你的项目目录下，使用
 
    ```go
      glide init   // 初始化一个glide工程
      glide get github.com/silsuer/bingo  // 下载并安装bingo的源代码
      bingo init   // 初始化一个bingo项目
    ```
    
    此时你的项目中应该有`.env.yaml`,`start.go`这两个文件以及一些文件夹
    现在执行
    `bingo run dev`, 在浏览器中输入：`localhost:12345`,看到欢迎界面，安装成功！
    （默认使用`12345`端口，如需更改，在`start.go`中指定）
    
 2. 正常（我为什么要用这个词...）安装
 
    ```markdown
       go get github.com/silsuer/bingo  // 获取并安装bingo
     
       bingo init  // 初始化项目
   
       bingo run dev // 运行初始化后的项目
     
     // 此时浏览器输入 localhost:12345 会出现Welcome to bingo字样。安装成功
    ```

## 路由

  `Bingo`的路由策略非常自由，基于`Httprouter`，性能强劲。
  
  推荐在`routes/web.go`中注册路由
   
   ```go
      // 使用Get方法，注册'/'，对应路由为打印一个hello world字符串
      // 注意最后一定要使用Register()将该路由注册进去
      utils.Route().Get("/", func(c *bingo.Context) {
      		fmt.Fprintln(c.Writer,"Hello World")
      	}).Register()
   ```
 
## 数据库操作：

*如果你不愿意使用`DB()`函数的话（因为使用这个函数之后需要转换为数据库类型，比如`bingo.DB().(*mysql.Mysql)`）*

*你也可以直接使用另一个辅助函数`MysqlDB()`,这个函数替你做了类型转换的操作*

*即可以这样使用`bingo.MysqlDB().Table("test").Find(12)`(查找id为12的数据)*

### 创建数据库
```go
     // 创建数据库
	// 可以输入1到3个参数，分别是 数据库名 字符集，排序规则
	// 默认只需要输入数据库名即可，字符集为utf8,排序规则为 utf8_general_ci
	// 如果只想输入字符集的话，默认支持utf8和gbk字符集，会自动指明排序规则，
	// 如果是其他的字符集，就必须输入第三个参数指明排序规则
	res := bingo.DB().(*mysql.Mysql).CreateDatabase("bingo")
	// 这里相当于 create database if not exists...
	res := bingo.DB().(*mysql.Mysql).CreateDatabaseIfNotExists("bingo2","utf8")
	res := bingo.DB().(*mysql.Mysql).CreateDatabase("bingo","gbk","gbk_chinese_ci")
    fmt.Fprintln(w,res.Result)
```


### 创建数据表

```go
	    // 创建数据库
		// 第一个参数是表名，第二个参数是回调，在回调中指定每一个列的类型，默认值，和备注等数据
	res := bingo.DB().(*mysql.Mysql).CreateTableIfNotExist("test", func(table *mysql.Blueprint) {
		table.Increments("id").Comment("自增的id")
		table.String("name").Default("default").Comment("姓名")
		table.Integer("age").Default(18).Comment("年龄")
	})
	fmt.Fprintln(w,res)
```
   
### 插入数据

Bingo提供了4种方法向数据库中插入数据

  1.向表中插入一条数据
   ```go
     	// 向表中插入一条数据
     	// 1.InsertOne  向数据库中插入一条数据，传入的参数是map[string]interface{},其中，map的键是字段名，值是要插入的值 
     	insert := make(map[string]interface{})
     	insert["id"] = 1
     	insert["name"] = "silsuer"
     	insert["age"] = 18
     	res := bingo.DB().(*mysql.Mysql).Table("test").InsertOne(insert)
     	fmt.Fprintln(w, res)
   ```
 2. 向表中插入一条数据，如果map中有表中不存在的字段，将会被过滤掉
 ```go
    // 过滤字段，插入数据
    insert := make(map[string]interface{})
	insert["name"] = "silsuer"
	insert["age"] = 18
	res := bingo.DB().(*mysql.Mysql).Table("test4").InsertOneCasual(insert)
	fmt.Fprintln(w, res)
  ```
  
 3. 批量插入数据，并过滤不存在的字段
 ```go
	// 批量插入数据，InsertCasual方法是对每一行都执行一次插入操作，而不是一条语句全部插入
	// 因此插入大量数据时会十分耗时，慎用
	// 接收一个 map[string]interface 的切片，每一片代表一行要插入的数据
	var insertData []map[string]interface{}
	for i:=0;i<100 ; i++ {
		insert := make(map[string]interface{})
		insert["name"] = "test"+ strconv.Itoa(i)
		insert["age"] = 18+i
		insertData = append(insertData,insert)
	}
	res:=bingo.DB().(*mysql.Mysql).Table("test4").InsertCasual(insertData)
	fmt.Fprint(w,res)
  ```
 
 4. 批量插入数据
 ```go
     	// 批量插入数据
     	// 接收一个 map[string]interface 的切片，每一片代表一行要插入的数据
     	var insertData []map[string]interface{}
     	for i:=0;i<100 ; i++ {
     		insert := make(map[string]interface{})
     		insert["name"] = "test"+ strconv.Itoa(i)
     		insert["age"] = 18+i
     		insertData = append(insertData,insert)
     	}
     	res:=bingo.DB().(*mysql.Mysql).Table("test4").Insert(insertData)
     	fmt.Fprint(w,res)
 ```
 
 ### 更新数据
 
 Bingo提供了3种更新数据的方法
 
 1. 更新一条数据（如果传入多余字段会报错）
 ```go
         // 更新一条数据
     	// 接收一个map[string]interface{} 作为参数，如果存在表中没有的字段，将会报错
     	// Where 函数接收2或者3个参数，2个参数，中间为默认的 = 号， 三个参数即 （“id”,">=",2），可以连续调用
     	a:= make(map[string]interface{})
     	a["name"] = "test"
     	res := bingo.DB().(*mysql.Mysql).Table("test4").Where("id",1703).UpdateOne(a)
     	fmt.Fprintln(w,res)
 ```
 
 2. 更新一条数据（如果传入多余字段会被过滤掉）
 
 ```go
     // 更新一条数据，并且过滤多余字段
	// 接收一个map[string]interface{} 作为参数，如果存在表中没有的字段，将会报错
	// 请注意 Where和OrWhere的用法
	a:= make(map[string]interface{})
	a["name"] = "test"
	res := bingo.DB().(*mysql.Mysql).Table("test4").Where("name","test4").OrWhere("id",">",2000).UpdateOneCasual(a)
	fmt.Fprintln(w,res)
 ```
 
 3. 批量更新(将会多次调用更新语句，大量数据时十分耗时，慎用)
 ```go
    // 更新多条数据，并且过滤多余字段
	// 接收一个map[string]interface{} 的切片 作为参数，如果存在表中没有的字段，将会报错
	var d []map[string]interface{}
	a:= make(map[string]interface{})
	a["name"] = "test"
	b:= make(map[string]interface{})
	b["age"] = 19
	d = append(d,a)
	d = append(d,b)
	res := bingo.DB().(*mysql.Mysql).Table("test4").Where("name","test4").OrWhere("id",">",2000).UpdateCasual(d)
	fmt.Fprintln(w,res)
 ```
 
 
### 查询数据
1. 正常查询

    ```go
        // 查询数据
        // 查询所有
        res:= bingo.DB().(*mysql.Mysql).Table("test4").Get()
        for res.Rows.Next() {
            var id,age int
            var name string
            res.Rows.Scan(&id,&name,&age)
            fmt.Fprintln(w,"id:"+strconv.Itoa(id)+" name:"+name+" age:"+strconv.Itoa(age))
        }
        
               // 条件查询
        res := bingo.DB().(*mysql.Mysql).Table("test4").Where("name","test93").Where("name","test94").OrWhere("id",">",1920).OrWhere("id","<",10).Get()
    
          // 排序和分组
          // 默认是ASC排序，可以在第二个参数中传入排序规则，当然，也可以多次调用 
        res := bingo.DB().(*mysql.Mysql).Table("test4").OrderBy("age").Get()
        res := bingo.DB().(*mysql.Mysql).Table("test4").OrderBy("age","asc").OrderBy("name").Get()
        
         // 直接指明升序排列
        res := bingo.DB().(*mysql.Mysql).Table("test4").OrderByAsc("age").Get()
        
          // 直接指明降序排列 （这样做的好处是IDE会有代码提示~）
        res := bingo.DB().(*mysql.Mysql).Table("test4").OrderByDesc("age").Get()
        
          // 查询前五条数据
        res := bingo.DB().(*mysql.Mysql).Table("test4").Limit(5).Get()
        
          // 查询第6到第15条数据
        res := bingo.DB().(*mysql.Mysql).Table("test4").Limit(5,10).Get()
        
          // 单独设置某个字段的值
        res := bingo.DB().(*mysql.Mysql).Table("test").Where("id",1).SetField("age",14)
        res := bingo.DB().(*mysql.Mysql).Table("test").Where("id",1).SetField("name","silsuer")
        fmt.Fprintln(w,res)
        
        	// 分组查询
        res := bingo.DB().(*mysql.Mysql).Table("test").Limit(5).GroupBy("age","id").Having("id",">",190).Get()
        for res.Rows.Next() {
        	var id,age int
        	var name string
        	res.Rows.Scan(&id,&name,&age)
        	fmt.Fprintln(w,"id:"+strconv.Itoa(id)+" name:"+name+" age:"+strconv.Itoa(age))
        }
     
        //查询单条记录
     	res := bingo.DB().(*mysql.Mysql).Table("test").First()
     	res := bingo.DB().(*mysql.Mysql).Table("test").Find(1)
     	var id, age int
     	var name string
     	res.Row.Scan(id, name, age)
     	fmt.Fprintln(w, "id:"+strconv.Itoa(id)+" name:"+name+" age:"+strconv.Itoa(age))
    ```

2. 关联查询
3. 分页
4. 事务
    
    ```go
            // 使用事务,Transaction中传入一个回调，即可使用事务
             res := bingo.DB().(*mysql.Mysql).Transaction(func() {
                a:=make(map[string]interface{})
                a["name"] = "test"
                a["age"] = 18
                bingo.DB().(*mysql.Mysql).InsertOne(a)
                bingo.DB().(*mysql.Mysql).InsertOneCasual(a)
            })
             fmt.Fprintln(w,res)

     ```
    

### 删除数据

```go

        // 1. 删除数据
        res := bingo.DB().(*mysql.Mysql).Table("test4").Where("id",1).Delete()
        
        // 2. 删除表
        res := bingo.DB().(*mysql.Mysql).Table("test4").DropTable()
        
        // 3. 清空数据表,Delete禁止不使用Where直接删除，如果需要直接全部产出的话，必须在Delete方法中传入true
        res := bingo.DB().(*mysql.Mysql).Table("test4").Delete(true)
        // 4. 清空数据表，delete是一条一条的删除，直到数据表为空，而truncate 是直接清空数据表，速度比较快
        res := bingo.DB().(*mysql.Mysql).Table("test4").Truncate()
        // 5. 清空数据库,删除所有表
        res := bingo.DB().(*mysql.Mysql).TruncateDatabase()
        // 6. 清空数据库中的所有表的数据，而不删除表
        res := bingo.DB().(*mysql.Mysql).TruncateDatabaseExceptTables()
        
         // 7. 删除数据库
        res := bingo.DB().(mysql.Mysql).DropDatabase()

```

## 会话管理

### 使用session

 Bingo使用 `gorilla/sessions` 管理session，具体用法如下：
 
 ```go
    func Index(c *bingo.Context) {
    	// 设置一个session
    	c.Session.Set("name","silsuer")
    	// 读取一个session
    	fmt.Fprintln(c.Writer,c.Session.Get("name"))
    }
 ```


## 开发脚手架

## 1. bingo run 工具

  - 使用 `bingo run dev` 运行开发环境程序
      
  - 使用 `bingo run daemon` 以守护进程运行程序
      
  - 使用 `bingo run daemon start` 以守护进程运行程序
      
  - 使用 `bingo run daemon restart` 平滑重启守护进程
      
  - 使用 `bingo run daemon stop` 平滑关闭重启
  
  - 使用 `bingo run watch` 在开发时监听工程目录下文件变更，当发现文件变更时，自动重启服务

## 2. 自定义命令

### 简介

`Laravel` 的 `artisan`也是一大亮点，我经常使用它来洗数据～～～

在`bingo`中，我仿照`artisan`实现了 `sword` 命令，目前在`mac`上可以使用，并未在linux和windows上测试

### 创建一个命令

`bingo sword make:command --name=CommandName`

该命令将会在 `app/Console/Commands` 目录下生成一个 `CommandName.go` 文件

该文件的内容：

```go
package Commands

import (
	"github.com/silsuer/bingo/cli"
)

// 该命令结构体
type ExampleCommand struct {
	cli.Command
	Name        string
	Description string
	Args        map[string]string
}

// 设置命令名
func (m *ExampleCommand) SetName() {
	m.Name = "command:name"
}

// 设置命令所需参数
func (m *ExampleCommand) SetArgs() {
	m.Args = make(map[string]string)
	m.Args["name"] = ""
}

// 设置命令描述
func (m *ExampleCommand) SetDescription() {
	m.Description = "the command description."
}

// 设置命令实现的方法
func (m *ExampleCommand) Handle(input cli.Input, output cli.Output) {

}

```

当需要使用该命令的时候，请务必现在`app/Console/Kernel.go`文件中注册：

```go

var Commands = []interface{}{
	&Command.ExampleCommand{},  // 请务必传入该命令的地址
}
```


### 执行命令

当注册好命令后，需要使用时，使用 `bingo sword command:name --arg=value` 调用这个命令

将会调用对应命令的`Handle`方法


### 目前已经写好的命令

`bingo sword make:command` 创建一个命令

`bingo sword make:origin:command` 创建一个内置命令(将在`bingo/cli`目录下建立新的文件)

 ####  未完待续