package mysql

import (
	"errors"
	"strings"
	"strconv"
)

// 查询数据
func (m *Mysql) Get() *Mysql {
	m.sql = "select" + m.columnSql + "from " + m.tableSql + m.whereSql + m.groupBySql + m.havingSql + m.limitSql + m.orderBySql
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
		m.columnSql = strings.TrimRight(m.columnSql, `,`) // 去掉最右侧的逗号
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
			if _, ok := cif[v]; ok {
				m.columnSql = m.columnSql + v + `,`
			}
		}
		m.columnSql = strings.TrimRight(m.columnSql, `,`) // 去掉最右侧的逗号
		break
	case string:
		for _, v := range args {
			if _, ok := cif[v.(string)]; ok {
				m.columnSql = m.columnSql + convertToString(v) + `,`
			}
		}
		m.columnSql = strings.TrimRight(m.columnSql, `,`) // 去掉最右侧的逗号
		break
	}
	return m
}

// First 获取第一行数据
func (m *Mysql) First() *Mysql {
	m.sql = "select" + m.columnSql + "from " + m.tableSql + m.whereSql + m.groupBySql + m.havingSql + m.limitSql + m.orderBySql
	m.Row = m.connection.QueryRow(m.sql) // 查询一行
	return m
}

// find
func (m *Mysql) Find(id int) *Mysql {
	// 只能传入一个id，然后返回一行数据
	m.Where("id", id).First() // 执行查询
	return m
}

// gruop by 对数据进行分组，传入多个字段值，一次查询只能调用一次GruopBy，调用多次后面的会覆盖掉前面的
func (m *Mysql) GroupBy(args ... string) *Mysql {
	var colNames []string
	for _, v := range args {
		colNames = append(colNames, v)
	}
	m.GroupByArr(colNames)
	return m
}

// 过滤掉不存在的列
func (m *Mysql) GroupByCasual(args ...  string) *Mysql {
	cif := m.GetTableInfo().Info
	var colNames []string
	for _, v := range args {
		// 如果存在这个列
		if _, ok := cif[v]; ok {
			colNames = append(colNames, v)
		}
	}
	m.GroupByArr(colNames)
	return m
}

func (m *Mysql) GroupByArr(colNames []string) *Mysql {
	// 拼接sql语句
	m.groupBySql = ` GROUP BY ` + strings.Join(colNames, `,`)
	return m
}

// 过滤掉不存在的列
func (m *Mysql) GroupByArrCasual(colNames []string) *Mysql {
	cif := m.GetTableInfo().Info
	var newColNames []string
	for _, v := range colNames {
		if _, ok := cif[v]; ok {
			newColNames = append(newColNames, v)
		}
	}
	m.GroupByArr(newColNames)
	return m
}

// having，由于having经常与聚组函数一起使用，所以这里不做过滤
// 可以传入 Having("sum(age)",">","10")
func (m *Mysql) Having(colName string, options string, value interface{}) *Mysql {
	subSql := ""
	switch value.(type) {
	case int:
		subSql = colName + options + strconv.Itoa(value.(int))
		break
	case string:
		subSql = colName + options + `'` + value.(string) + `'`
		break
	default:
		m.checkAppendError(errors.New(`the having function need int or string argument as the 3rd argument`))
		break
	}

	// 拼接having语句
	// 如果存在having
	if havingArr := strings.Fields(m.havingSql); len(havingArr) != 0 && strings.ToUpper(havingArr[0]) == "HAVING" {
		m.havingSql = m.havingSql + ` , ` + subSql
	} else {
		m.havingSql = " HAVING " + subSql
	}
	return m
}

// SetField  设置单独更新某个字段的值
func (m *Mysql) SetField(colName string, value interface{}) *Mysql {
	kvMap := make(map[string]interface{})
	kvMap[colName] = value
	m.UpdateOne(kvMap)
	return m
}

// 连接查询 Join


//// 关联查询 HasOne HasMany

// 事务
