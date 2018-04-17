package mysql

import (
	"strings"
	"strconv"
	"errors"
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
		bp.colSql[bp.currentCol] = bp.colSql[bp.currentCol] + " default " + strings.TrimSpace(strconv.Itoa(def.(int)))
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

func (m *Mysql) CreateTableIfNotExist(tableName string, call func(table *Blueprint)) error {
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
	stmt, err := m.connection.Prepare(table.sql)
	_, err = stmt.Exec()
	return err
}

// 创建数据库并且指定字符集
func (m *Mysql) CreateDatabase(args ... string) *Mysql {
	// 最多可以传3个参数，第一个是名字，第二个是字符集,第三个是排序规则
	if len(args) == 0 {
		m.checkAppendError(errors.New(`the function CreateDatabase need 1 or 2 or 3 args. Use case: CreateDatabase("db_name","gbk","gbk_chinese_ci")`))
	}
	dbName := args[0] // 指定数据库名称
	//dbCharset := "utf8"            // 设置默认字符集
	//dbCollate := "utf8_general_ci" // 设置默认排序规则

	switch len(args) {
	case 1: // 默认是uft8的
		m.sql = `CREATE DATABASE  ` + dbName + ` DEFAULT CHARSET utf8 COLLATE utf8_general_ci`
		break
	case 2:
		if strings.ToUpper(args[1]) == "UTF8" || strings.ToUpper(args[1]) == "UTF-8" {
			m.sql = `CREATE DATABASE  ` + dbName + ` DEFAULT CHARSET utf8 COLLATE utf8_general_ci`
		} else if strings.ToUpper(args[1]) == "GBK" {
			m.sql = `CREATE DATABASE  ` + dbName + ` DEFAULT CHARSET gbk COLLATE gbk_chinese_ci`
		} else {
			// 除了这utf8和gbk这两种，其他的，必须指明排序规则
			m.checkAppendError(errors.New("you must set your collate when create a new database"))
			return m
		}
		break
	case 3:
		m.sql = `CREATE DATABASE  ` + dbName + ` DEFAULT CHARSET ` + args[1] + ` COLLATE ` + args[2]
		break
	default:
		m.checkAppendError(errors.New(`too many arguments in the CreateDatabase function`))
		break
	}
	res,err:= m.Exec(m.sql)
	m.Result = res
	m.checkAppendError(err)
	return m
}

// 创建数据库（如果不存在的话）
func (m *Mysql) CreateDatabaseIfNotExists(args ... string) *Mysql  {
	// 最多可以传3个参数，第一个是名字，第二个是字符集,第三个是排序规则
	if len(args) == 0 {
		m.checkAppendError(errors.New(`the function CreateDatabase need 1 or 2 or 3 args. Use case: CreateDatabase("db_name","gbk","gbk_chinese_ci")`))
	}
	dbName := args[0] // 指定数据库名称
	switch len(args) {
	case 1: // 默认是uft8的
		m.sql = `CREATE DATABASE IF NOT EXISTS  ` + dbName + ` DEFAULT CHARSET utf8 COLLATE utf8_general_ci`
		break
	case 2:
		if strings.ToUpper(args[1]) == "UTF8" || strings.ToUpper(args[1]) == "UTF-8" {
			m.sql = `CREATE DATABASE IF NOT EXISTS  ` + dbName + ` DEFAULT CHARSET utf8 COLLATE utf8_general_ci`
		} else if strings.ToUpper(args[1]) == "GBK" {
			m.sql = `CREATE DATABASE IF NOT EXISTS  ` + dbName + ` DEFAULT CHARSET gbk COLLATE gbk_chinese_ci`
		} else {
			// 除了这utf8和gbk这两种，其他的，必须指明排序规则
			m.checkAppendError(errors.New("you must set your collate when create a new database"))
			return m
		}
		break
	case 3:
		m.sql = `CREATE DATABASE IF NOT EXISTS  ` + dbName + ` DEFAULT CHARSET ` + args[1] + ` COLLATE ` + args[2]
		break
	default:
		m.checkAppendError(errors.New(`too many arguments in the CreateDatabase function`))
		break
	}
	res,err:= m.Exec(m.sql)
	m.Result = res
	m.checkAppendError(err)
	return m
}