package mysql

import (
	"database/sql"
	"sync"
	"fmt"
)

// mysql操作类
// DB("table_name").Get()
// DB("table_name").Where().Where().Select().Get()
// 数据库链接结构体,单例模式
var conn *sql.DB
// mysql结构体，用来存储sql语句并执行
type Mysql struct {
	connection   *sql.DB
	sql          string
	whereSql    string
	limitSql    string
	columnSql   string
	tableSql    string
	orderBySql string
	tableName string
	TableSchema string  // 所在数据库，在缓存表结构的时候有用
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


func (m *Mysql) Where() *Mysql {
	return m
}

func (m *Mysql) OrWhere() *Mysql {
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
	rows,err := m.connection.Query(m.sql)

	Check(err)
	return rows
}

// 执行原生语句
func (m *Mysql) Query(sql string) (*sql.Rows,error)  {
	return m.connection.Query(sql)
}

func (m *Mysql) Exec(sql string ) (sql.Result,error)   {
	return m.connection.Exec(sql)
}



func Check(err error)  {
	if err != nil{
		panic(err)
	}
}
