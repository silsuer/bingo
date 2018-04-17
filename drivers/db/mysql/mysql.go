package mysql

import (
	"database/sql"
	"sync"
	"errors"
	"strings"
)

// mysql操作类
// DB("table_name").Get()
// DB("table_name").Where().Where().Select().Get()
// 数据库链接结构体,单例模式
var conn *sql.DB
// mysql结构体，用来存储sql语句并执行
type Mysql struct {
	connection  *sql.DB
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

func (m *Mysql) Init(config string) {
	// 获取单例连接
	m.connection = GetInstanceConnection(config) // 获取数据库连接
}

func (m *Mysql) Table(tableName string) *Mysql {
	m.tableSql = " " + tableName + " "
	m.whereSql = ""
	m.columnSql = " * "
	m.tableName = tableName
	return m
}

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

func (m *Mysql) OrWhere(args ... interface{}) *Mysql {
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

	whereArr := strings.Fields(m.whereSql)
	var ifWhere bool
	if strings.ToUpper(whereArr[0]) == "WHERE" {
		ifWhere = true // 存在where
	} else {
		ifWhere = false // 不存在where
	}
	switch len(args) {
	case 2:
		// 有where 只需要写 a=b
		if ifWhere {
			// 如果是字符串类型，加引号
			if isString(cType) {
				m.whereSql = m.whereSql + ` OR ` + convertToString(args[0]) + `='` + convertToString(args[1]) + `'`
			} else {
				m.whereSql = m.whereSql + ` OR ` + convertToString(args[0]) + `=` + convertToString(args[1])
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
		m.Errors = append(m.Errors, errors.New("missing args length in orWhere function"))
	}
	return m
}

//
//func (m *Mysql) HasMany() *Mysql {
//	return m
//}
//
//func (m *Mysql) HasOne() *Mysql {
//	return m
//}

// 执行原生语句
func (m *Mysql) Query(sql string) (*sql.Rows, error) {
	return m.connection.Query(sql)
}

func (m *Mysql) Exec(sql string) (sql.Result, error) {
	return m.connection.Exec(sql)
}

func (m *Mysql) Limit(args ... int) *Mysql {
	// 传入一个或者两个参数
	if len(args) == 0 || (len(args) != 1 && len(args) != 2) {
		m.checkAppendError(errors.New(`the Limit function need 1 or 2 arguments `))
		return m // 终止程序
	}
	if len(args) == 1 {
		m.limitSql = ` limit ` + convertToString(args[0])
	} else {
		m.limitSql = ` limit ` + convertToString(args[0]) + `,` + convertToString(args[1])
	}
	return m
}

func (m *Mysql) OrderBy(args ... string) *Mysql {
	// 最多可以传入2个参数 ，第一个是字段名，第二个是排序规则
	// 判断现在的order by 语句中是否有order by，有的话，就在后面加 ,colName asc 
	if len(args) == 0 || (len(args) != 1 && len(args) != 2) {
		m.checkAppendError(errors.New(`the OrderBy function need 1 or 2 arguments `))
		return m
	}
	collate := "ASC"
	if len(args) == 2 {
		collate = args[1] // 如果传入了第二个参数，就赋值给排序规则
	}

	if odArr := strings.Fields(m.orderBySql); len(odArr) == 0 || strings.ToUpper(odArr[0]) != "ORDER" { // 没有order by
		m.orderBySql = ` ORDER BY ` + args[0] + ` ` + collate
	} else {
		m.orderBySql = m.orderBySql + `,` + args[0] + ` ` + collate
	}
	return m
}

func (m *Mysql) OrderByAsc(colName string) *Mysql {
	// 传入排序的字段名，升序排序
    m.OrderBy(colName,"ASC")
	return m
}

func (m *Mysql) OrderByDesc(colName string) *Mysql {
	//　传入排序的字段名，降序排序
    m.OrderBy(colName,"DESC")
	return m
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// 检查error
func (m *Mysql) checkAppendError(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}
