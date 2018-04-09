package mysql

import (
	"strings"
	"fmt"
	"strconv"
)

// 建表

type Blueprint struct {
	sql            string
	engine         string
	colSql         map[string]string
	currentCol     string
	currentColType string
}

func (bp *Blueprint) Increments(colName string) *Blueprint {
	//Id int primary key auto_increment,#部门编号 整形 主键 自增长
	if len(bp.colSql) == 0 {
		bp.colSql = make(map[string]string)
	}
	bp.colSql[colName] = colName + " int(11) not null primary key auto_increment"
	bp.currentCol = colName
	bp.currentColType = "integer"
	return bp
}

func (bp *Blueprint) Nullable() *Blueprint {
	bp.colSql[bp.currentCol] = strings.Replace(bp.colSql[bp.currentCol], "not null", "", -1) // 替换not null字符串为空
	return bp
}

func (bp *Blueprint) Comment(comment string) *Blueprint {
	bp.colSql[bp.currentCol] = bp.colSql[bp.currentCol] + " comment '" + comment + "'"
	return bp
}

// 给建表语句添加默认值
func (bp *Blueprint) Default(def interface{}) *Blueprint {
	if bp.currentColType == "integer" {
		bp.colSql[bp.currentCol] = bp.colSql[bp.currentCol] + " default " + strings.TrimSpace(def.(string))
	} else {
		bp.colSql[bp.currentCol] = bp.colSql[bp.currentCol] + " default '" + strings.TrimSpace(def.(string)) + "'"
	}
	return bp
}

// 字符串形式的数据
func (bp *Blueprint) String(colName string) *Blueprint {
	return bp.StringWithLength(colName, 255)
}

func (bp *Blueprint) StringWithLength(colName string, length int) *Blueprint {
	if len(bp.colSql) == 0 {
		bp.colSql = make(map[string]string)
	}
	//Name varchar(18)
	bp.colSql[colName] = colName + " varchar(" + strconv.Itoa(length) + ") not null "
	bp.currentCol = colName
	bp.currentColType = "varchar"
	return bp
}

func (bp *Blueprint) Integer(colName string) *Blueprint {
	return bp.IntegerWithLength(colName, 11)
}

func (bp *Blueprint) IntegerWithLength(colName string, length int) *Blueprint {
	if len(bp.colSql) == 0 {
		bp.colSql = make(map[string]string)
	}
	bp.colSql[colName] = colName + " int(" + strconv.Itoa(length) + ") not null"
	bp.currentCol = colName
	bp.currentColType = "integer"
	return bp
}

func (m *Mysql) CreateTableIfNotExist(tableName string, call func(table *Blueprint)) bool {
	table := Blueprint{}
	call(&table)
	// 判断是否指定了引擎，默认是innodb
	if table.engine == "" {
		table.engine = "innodb"
	}
	// 把所有数据拼合成建表语句，然后执行这个语句
	table.sql = "create table " + tableName + "("
	var cols []string
	for _, v := range table.colSql {
		cols = append(cols, v)
	}
	table.sql = table.sql + strings.Join(cols, ",")
	table.sql = table.sql + ") engine=" + table.engine
	fmt.Println(table.sql)
	stmt, err := m.connection.Prepare(table.sql)
	Check(err)
	_, err = stmt.Exec()
	Check(err)
	return true
}
