# 使用Go封装一个便捷的ORM

最近在用Go写一个Web开发框架，看了一些ORM,大部分都需要自己拼接SQL，对我这种用惯了
`Laravel`的人来说，确实有点别扭，所以想自己写一个ORM，可以方便的对数据库进行连贯操作

由于代码太多，不贴了，只讲思路，具体代码在这里[silsuer/bingo](https://github.com/silsuer/bingo)

## 思路

1. 确定最后要做出的效果
   
   我想要做成类似`Laravel`那种，操作数据库大概是这样`DB::table(dbName)->Where('id',1)->get()`

2. 连贯操作原理   
   
   做出这种连贯操作的效果，除了结尾的方法，中间连续调用的那些方法都必须返回一个相同的对象
   
3. 定义数据库对象

   既然如此，那么先定义一个`Mysql`的结构体
   
   ```go

    // mysql结构体，用来存储sql语句并执行
    type Mysql struct {
        connection  *sql.DB    // 数据库连接对象
        sql         string     // 最终要执行的语句
        whereSql    string     // where语句
        limitSql    string     // limit语句
        columnSql   string     // select 中的列语句
        tableSql    string     // 表名语句
        orderBySql  string     // order by 语句
        groupBySql  string     // group by语句
        havingSql   string     // having语句
        tableName   string     // 表名
        TableSchema string     // 所在数据库，在缓存表结构的时候有用
        Errors      []error    // 保存在连贯操作中可能发生的错误
        Result      sql.Result // 保存执行语句后的结果
        Rows        *sql.Rows  // 保存执行多行查询语句后的行
        Results     []sql.Result
        Row         *sql.Row // 保存单行查询语句后的行
        //Res         []map[string]interface{}  // 把查询结果生成一个map
    }

   ```
   
   这样，我们每次连续调用的时候，只需要返回当前结构体的指针就可以了.
   
4. 连接数据库

   为了方便，先在公共函数中定义一个`DB()`函数，这个函数通过创建数据库驱动，来返回一个数据库对象（指针）
   
   ```go
        func DB() interface{} {
            // 返回一个驱动的操作类
            // 传入db配置  env的driver，实例化不同的类,单例模式，获取唯一驱动
            if Driver == nil {
                DriverInit()
            }
            con := Driver.GetConnection()  // 获取数据库连接
            return con
        }
   ```
   
   可以看到，为了节省资源，我们的数据库连接和驱动连接都是单例模式，只允许存在一次
   
   下面是初始化数据库驱动的方法`DriverInit()`
   
   ```go

   
    // 数据库驱动
    type driver struct {
        name     string // 驱动名
        dbConfig string // 配置
    }
    
    // 驱动初始化
    func DriverInit() {
        Driver = &driver{}
        Driver.name = strings.ToUpper(Env.Get("DB_DRIVER"))
        switch Driver.name {
        case "MYSQL":
            // 初始化了驱动之后，开始初始化数据库连接
            Driver.dbConfig = Env.Get("DB_USERNAME") + ":" + Env.Get("DB_PASSWORD") + "@tcp(" + Env.Get("DB_HOST") + ":" + Env.Get("DB_PORT") + ")" + "/" + Env.Get("DB_NAME") + "?" + "charset=" + Env.Get("DB_CHARSET")
            break
        default:
            break
        }
    
    }
   ```
   
   我们在初始化驱动的时候，会判断配置文件中用的是哪种数据库，然后根据数据库去拼接连接数据库的字符串
   
   这样为我们的ORM可以支持多个数据库提供了可能
   
   配置的字符串做好了，接下来就要开始获取数据库的单例连接了
   
   ```go
       // 根据数据库驱动，获取数据库连接
       func (d *driver) GetConnection() interface{} {
       	switch d.name {
       	case "MYSQL":
       		m := mysql.Mysql{}                // 实例化结构体
       		m.TableSchema = Env.Get("DB_NAME") // 数据库名
       		m.Init(Driver.dbConfig) // 设置表名和数据库连接
       		return &m                         // 返回实例
       		break
       	default:
       		break
       	}
       	return nil
       }
   ```
   
   获取连接时首先创建了一个数据库的对象，然后把一些基本信息赋给了这个对象的对应属性，`Init`方法
   
   获取了这个人数据库的唯一连接:
   
   ```go
       func (m *Mysql) Init(config string) {
        // 获取单例连接
        m.connection = GetInstanceConnection(config) // 获取数据库连接
       }

        var once sync.Once
        // 设置单例的数据库连接
        func GetInstanceConnection(config string) *sql.DB {
            once.Do(func() { // 只做一次
                db, err := sql.Open("mysql", config)
                if err != nil {
                    panic(err)
                }
                conn = db
            })
            return conn
        }
   ```
   
   `Once.Do`是`sync`包中提供的只允许代码执行一次的方法
   
   这样我们就获取到了数据库的连接，然后把这个连接放到`Mysql`对象的 `connection`中，返回即可
   
5. 执行连贯操作

   要注意由于我们使用的数据库不同，`DB()`返回的是一个`interface{}`,所以需要我们重新转成`*mysql.Mysql`:
   
   ```go
     bingo.DB().(*mysql.Mysql)
   ```
   
   接下来，只需要定义各种`Mysql`结构体的方法，然后返回这个结构体的指针就可以了，
   
   篇幅问题，只写一个最简单的在`test`表中查询`id>1`的数据,代码效果是这样:
  
   ```go
       res:= bingo.DB().(*mysql.Mysql).Table("test").Where("id",">",1).Get()
   ```
   
   首先，写`Table`方法，只是给Mysql对象赋个值而已:
   
   ```go
        // 初始化一些sql的值
       func (m *Mysql) Table(tableName string) *Mysql {
        m.tableSql = " " + tableName + " "
        m.whereSql = ""
        m.columnSql = " * "
        m.tableName = tableName
        return m
       }
   ```
   
   然后写`Where`方法，这个稍微复杂一点，如果传入两个参数，那么他们之间是等于的关系，如果是三个参数
   
   那么第二个参数就是他们直接的关系
   
   ```go

     func (m *Mysql) Where(args ... interface{}) *Mysql {
     	// 要根据传入字段的名字，获得字段类型，然后判断第三个参数是否要加引号
     	// 最多传入三个参数  a=b
     	// 先判断where sql 里有没有 where a=b
     	cif := m.GetTableInfo().Info
     
     	// 先判断这个字段在表中是否存在,如果存在，获取类型获取值等
     	cType := ""
     	if _, ok := cif[convertToString(args[0])]; ok {
     		cType = cif[convertToString(args[0])].Type
     	} else {
     		m.Errors = append(m.Errors, errors.New("cannot find a column named "+convertToString(args[0])+" in "+m.tableName+" table"))
     		return m // 终止这个函数
     	}
     	var ifWhere bool
     	whereArr := strings.Fields(m.whereSql)
     	if len(whereArr) == 0 { // 空的
     		ifWhere = false
     	} else {
     		if strings.ToUpper(whereArr[0]) == "WHERE" {
     			ifWhere = true // 存在where
     		} else {
     			ifWhere = false // 不存在where
     		}
     	}
     
     	switch len(args) {
     	case 2:
     		// 有where 只需要写 a=b
     		if ifWhere {
     			// 如果是字符串类型，加引号
     			if isString(cType) {
     				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + `='` + convertToString(args[1]) + `'`
     			} else {
     				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + `=` + convertToString(args[1])
     			}
     		} else {
     			// 没有where 要写 where a=b
     			if isString(cType) {
     				m.whereSql = ` WHERE ` + convertToString(args[0]) + `='` + convertToString(args[1]) + `'`
     			} else {
     				m.whereSql = ` WHERE ` + convertToString(args[0]) + `=` + convertToString(args[1])
     			}
     		}
     		break
     	case 3:
     		// 有where 只需要写 a=b
     		if ifWhere {
     			// 如果是字符串类型，加引号
     			if isString(cType) {
     				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + convertToString(args[1]) + `'` + convertToString(args[2]) + `'`
     			} else {
     				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + convertToString(args[1]) + convertToString(args[2])
     			}
     		} else {
     			// 没有where 要写 where a=b
     			if isString(cType) {
     				m.whereSql = ` WHERE ` + convertToString(args[0]) + convertToString(args[1]) + `'` + convertToString(args[2]) + `'`
     			} else {
     				m.whereSql = ` WHERE ` + convertToString(args[0]) + convertToString(args[1]) + convertToString(args[2])
     			}
     		}
     		break
     	default:
     		m.Errors = append(m.Errors, errors.New("missing args length in where function"))
     	}
     	return m
     }

   ```
   
   上面的代码中我用到了`GetTableInfo()` 方法，这个方法是获取当前表的结构，为了快速的对表执行操作，或者过滤掉表不存在的字段
   
   需要缓存表的结构，把列名和列类型缓存起来，具体代码不贴了，在这里 [缓存表结构](https://github.com/silsuer/bingo/blob/master/drivers/db/mysql/insert.go)
   
   然后会根据`Mysql`结构体的`whereSql`字段，拼接新的`whereSql`,并重新赋值
   
   接下来，写`Get()`方法，执行最后的SQL语句即可
   
   ```go
   
    // 查询数据
    func (m *Mysql) Get() *Mysql {
    	m.sql = "select" + m.columnSql + "from " + m.tableSql + m.whereSql + m.limitSql + m.orderBySql
    	rows, err := m.connection.Query(m.sql)
    	m.checkAppendError(err)   
    	m.Rows = rows
    	return m
    }
   ```
   
   首先拼接语句，然后执行原生的Query方法，把执行结果放置在结构体中即可
   
   这样我们就拿到了结构体，可以把它打出来看看
   
   ```go
    // 获取test表中id大于1的所有数据
    res:= bingo.DB().(*mysql.Mysql).Table("test").Where("id",">",1).Get()
    // 判断执行是否出错
    if len(res.Errors)!=0{
        	fmt.Fprintln(w,"执行出错！")
    	}
 	 // 执行成功，开始遍历
        for res.Rows.Next() {
            var id,age int
            var name string
            res.Rows.Scan(&id,&name,&age)
            fmt.Fprintln(w,"id:"+strconv.Itoa(id)+" name:"+name+" age:"+strconv.Itoa(age))
        }
   ```
   
6. 小结
   
   由于篇幅问题，我只写了最简单的条件查询，其他的可以去[GitHub](https://github.com/silsuer/bingo)上看
   
   目前实现的有：
    
        1. 创建数据库
        2. 创建数据表
        3. 插入数据（4种方法：插入一条数据，插入一条数据并过滤多余字段，批量插入数据，批量插入数据并过滤多余字段）
        4. 更新数据（3种方法：更新一条数据，更新一条数据并过滤多余字段，批量更新数据并过滤多余字段）
        5. 查询数据（条件查询、limit、order by等）
        6. 删除数据
        7. 删除数据表（2种方法：delete删除，truncate清空）
        8. 清空数据库中的所有表信息而不删除表
        9. 清空数据库，删除所有表
        10. 删除数据表
   
   将要实现的有：
       
        1. 关联查询
        2. join查询
        3. 分组查询
        4. 快速分页
        5. 数据库迁移 
        
   写的时候也遇到过一些小坑，比如说批量插入的时候，由于`map`是无序的，所以在遍历的时候会遇到字段名和要插入的值对不上的情况
   
   解决办法就是首先遍历一次，将字段的顺序固定下来，然后根据字段的顺序去设置后面要插入的值的顺序就可以了
   
7. 具体用法请去看看README，时间问题，写的有点乱，等这个玩意儿开发完了会重新整理一份文档的

   最后一句，求star，求翻牌子~~~ヾ(*´▽‘*)ﾉ

   