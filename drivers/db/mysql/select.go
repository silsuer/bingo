package mysql

import (
	"errors"
	"strings"
)

// 查询数据
func (m *Mysql) Get() *Mysql {
	m.sql = "select" + m.columnSql + "from " + m.tableSql + m.whereSql + m.limitSql + m.orderBySql
	rows, err := m.connection.Query(m.sql)
	m.checkAppendError(err)
	//defer rows.Close()
	m.Rows = rows
	return m
}

// 可以传入一个[]string数组，也可以传入多个string
func (m *Mysql) Select(args ... interface{}) *Mysql {
	if len(args) == 0 {
		m.Errors = append(m.Errors, errors.New(`the Select function need a arg([]string or multi-string)`))
	}
	m.columnSql = `` // 把原来的字段置空
	switch args[0].(type) {
	case []string:
		for _, v := range args[0].([]string) {
			m.columnSql = m.columnSql + v + `,`
		}
		m.columnSql = strings.TrimRight(m.columnSql, `,`) // 去掉最右侧的逗号
		break
	case string:
		for _, v := range args {
			m.columnSql = m.columnSql + convertToString(v) + `,`
		}
		m.columnSql = strings.TrimRight(m.columnSql,`,`) // 去掉最右侧的逗号
		break
	}

	return m
}

// 指定要查询的字段，过滤掉表中不存在的列
func (m *Mysql) SelectCasual(args ... interface{}) *Mysql {
	if len(args) == 0 {
		m.Errors = append(m.Errors, errors.New(`the Select function need a arg([]string or multi-string)`))
	}
	cif := m.GetTableInfo().Info
	m.columnSql = `` // 把原来的字段置空
	switch args[0].(type) {
	case []string:
		for _, v := range args[0].([]string) {
			// 判断是否存在这个字段
			if _,ok:=cif[v];ok {
				m.columnSql = m.columnSql + v + `,`
			}
		}
		m.columnSql = strings.TrimRight(m.columnSql, `,`) // 去掉最右侧的逗号
		break
	case string:
		for _, v := range args {
			if _,ok:=cif[v.(string)];ok{
				m.columnSql = m.columnSql + convertToString(v) + `,`
			}
		}
		m.columnSql = strings.TrimRight(m.columnSql,`,`) // 去掉最右侧的逗号
		break
	}
	return m
}

// First
// find
// SetField

