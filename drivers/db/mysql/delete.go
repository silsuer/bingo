package mysql

import (
	"strings"
	"errors"
)

// 对表进行删除,如果不写WHERE的话，不执行
func (m *Mysql) Delete(args ... bool) *Mysql {
	// 如果传入了true的话，允许清空整个表数据，如果没传入参数的话，不允许没有where的数据出现
	if len(args) == 0 || (len(args) >= 1 && !args[0]) {
		if iWhere := strings.Fields(m.whereSql); len(iWhere) == 0 {
			// 设置报错信息
			m.Errors = append(m.Errors, errors.New("cannot delete all the data in table '"+m.tableName+"'. You can set the arg 'true' to exec it"))
			return m
		}
	}

	// delete from
	m.sql = `DELETE FROM ` + m.tableName + ` ` + m.whereSql
	res, err := m.Exec(m.sql)
	m.checkAppendError(err)
	m.Result = res
	return m
}

// 清空表数据（使用truncate）
func (m *Mysql) Truncate() *Mysql {
	m.sql = `TRUNCATE TABLE ` + m.tableName
	res, err := m.Exec(m.sql)
	m.checkAppendError(err)
	m.Result = res
	return m
}

// 删除数据表
func (m *Mysql) DropTable() *Mysql {
	m.sql = `DROP TABLE ` + m.tableName
	res, err := m.Exec(m.sql)
	m.checkAppendError(err)
	m.Result = res
	return m
}

//  清空数据库，删除所有表
func (m *Mysql) TruncateDatabase() *Mysql {
	// 获取到数据库中的所有表名
	sql := `SELECT table_name FROM information_schema.TABLES WHERE table_schema='` + m.TableSchema + `'`
	rows, err := m.Query(sql)
	// 执行不成功，会终止执行
	if err != nil {
		m.checkAppendError(err)
		return m
	}
	for rows.Next() {
		// drop table table_name
		var name string
		rows.Scan(&name)
		dropSql := `DROP TABLE ` + name
		res, err := m.Exec(dropSql)
		m.checkAppendError(err)
		// 把结果保存在res中
		m.Results = append(m.Results, res)
	}
	return m
}

//  清空数据库中所有表的数据而不删除表
func (m *Mysql) TruncateDatabaseExceptTables() *Mysql {
	// 获取到数据库中的所有表名
	sql := `SELECT table_name FROM information_schema.TABLES WHERE table_schema='` + m.TableSchema + `'`
	rows, err := m.Query(sql)
	// 执行不成功，会终止执行
	if err != nil {
		m.checkAppendError(err)
		return m
	}
	for rows.Next() {
		// drop table table_name
		var name string
		rows.Scan(&name)
		dropSql := `TRUNCATE TABLE ` + name
		res, err := m.Exec(dropSql)
		m.checkAppendError(err)
		// 把结果保存在res中
		m.Results = append(m.Results, res)
	}
	return m
}

//  删除数据库
func (m *Mysql) DropDatabase() *Mysql {
	m.sql = `DROP DATABASE ` + m.TableSchema
	res, err := m.Exec(m.sql)
	m.checkAppendError(err)
	m.Result = res
	return m
}
