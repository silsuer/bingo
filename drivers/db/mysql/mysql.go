package mysql

import (
	"database/sql"
	"sync"
	"fmt"
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
	sql         string
	whereSql    string
	limitSql    string
	columnSql   string
	tableSql    string
	orderBySql  string
	tableName   string
	TableSchema string  // 所在数据库，在缓存表结构的时候有用
	Errors      []error // 保存在连贯操作中可能发生的错误
	Result      sql.Result   // 保存执行语句后的结果
	Rows        *sql.Rows  // 保存执行多行查询语句后的行
	Results     []sql.Result
	Row         *sql.Row   // 保存单行查询语句后的行
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
	ctype := ""
	if _, ok := cif[convertToString(args[0])]; ok {
		ctype = cif[convertToString(args[0])].Type
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
			if isString(ctype) {
				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + `='` + convertToString(args[1]) + `'`
			} else {
				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + `=` + convertToString(args[1])
			}
		} else {
			// 没有where 要写 where a=b
			if isString(ctype) {
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
			if isString(ctype) {
				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + convertToString(args[1]) + `'` + convertToString(args[2]) + `'`
			} else {
				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + convertToString(args[1]) + convertToString(args[2])
			}
		} else {
			// 没有where 要写 where a=b
			if isString(ctype) {
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

//func (m *Mysql) whereSql()  {
//
//}

func (m *Mysql) OrWhere(args ... interface{}) *Mysql {
	// 要根据传入字段的名字，获得字段类型，然后判断第三个参数是否要加引号
	// 最多传入三个参数  a=b
	// 先判断where sql 里有没有 where a=b
	cif := m.GetTableInfo().Info

	// 先判断这个字段在表中是否存在,如果存在，获取类型获取值等
	ctype := ""
	if _, ok := cif[convertToString(args[0])]; ok {
		ctype = cif[convertToString(args[0])].Type
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
			if isString(ctype) {
				m.whereSql = m.whereSql + ` OR ` + convertToString(args[0]) + `='` + convertToString(args[1]) + `'`
			} else {
				m.whereSql = m.whereSql + ` OR ` + convertToString(args[0]) + `=` + convertToString(args[1])
			}
		} else {
			// 没有where 要写 where a=b
			if isString(ctype) {
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
			if isString(ctype) {
				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + convertToString(args[1]) + `'` + convertToString(args[2]) + `'`
			} else {
				m.whereSql = m.whereSql + ` AND ` + convertToString(args[0]) + convertToString(args[1]) + convertToString(args[2])
			}
		} else {
			// 没有where 要写 where a=b
			if isString(ctype) {
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

func (m *Mysql) Limit() *Mysql {
	return m
}

func (m *Mysql) Select() *Mysql {
	return m
}

func (m *Mysql) HasMany() *Mysql {
	return m
}

func (m *Mysql) HasOne() *Mysql {
	return m
}

func (m *Mysql) Get() *sql.Rows {
	m.sql = "select" + m.columnSql + "from " + m.tableSql + m.whereSql + m.limitSql + m.orderBySql
	fmt.Println(m.sql)
	rows, err := m.connection.Query(m.sql)

	Check(err)
	return rows
}

// 执行原生语句
func (m *Mysql) Query(sql string) (*sql.Rows, error) {
	return m.connection.Query(sql)
}

func (m *Mysql) Exec(sql string) (sql.Result, error) {
	return m.connection.Exec(sql)
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
