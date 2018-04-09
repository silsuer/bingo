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
	where_sql    string
	limit_sql    string
	column_sql   string
	table_sql    string
	order_by_sql string
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
	m.table_sql = " " + tableName + " "
	m.where_sql = ""
	m.column_sql = " * "
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
	m.sql = "select" + m.column_sql + "from " + m.table_sql + m.where_sql + m.limit_sql + m.order_by_sql
	fmt.Println(m.sql)
	rows,err := m.connection.Query(m.sql)

	Check(err)
	return rows
}

func Check(err error)  {
	if err != nil{
		panic(err)
	}
}
