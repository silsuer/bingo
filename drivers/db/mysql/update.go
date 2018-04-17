package mysql

import (
	"database/sql"
	"errors"
	"strings"
)

// 对于mysql数据库的更新操作
// 1. 更新一条
// 2. 批量更新
// 3. 更新一条忽略多余值
// 4. 批量更新忽略多余值

// 更新一条数据
func (m *Mysql) UpdateOne(data map[string]interface{}) *Mysql {
	cif := m.GetTableInfo().Info // 获取表结构数据
	kvSql := ``                  // 初始化要拼接的sql语句
	for key, value := range data { // 遍历要传入的数据
		// 拼接set后的字符串  a=1,b='2',c=11
		if attr, ok := cif[key]; ok { // 如果表中存在这个字段
			// 存在这个字段
			if isString(attr.Type) {
				kvSql = kvSql + key + `='` + convertToString(value) + `',`
			} else {
				kvSql = kvSql + key + `=` + convertToString(value) + `,`
			}
		} else {
			// 如果不存在，返回错误信息
			m.Errors = append(m.Errors, errors.New("cannot find a column named "+key+" in "+m.tableName+" table "))
		}
	}
	kvSql = strings.TrimRight(kvSql, `,`) // 去掉最后的逗号
	m.sql = `UPDATE ` + m.tableName + ` SET ` + kvSql + ` ` + m.whereSql
	res, err := m.Exec(m.sql)
	m.Result = res
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
	return m
}

// 更新一条数据（过滤不存在的列）
func (m *Mysql) UpdateOneCasual(data map[string]interface{}) *Mysql {
	var d []map[string]interface{}
	d = append(d,data)
	res,err := m.updateSqlCasual(d)
	m.Result = res[0]
	if err[0]!=nil {
		m.Errors = append(m.Errors,err[0])
	}
	return m
}

//func (m *Mysql) updateSql(data []map[string]interface{})  {
//
//}

// 批量更新 太过麻烦，暂时先不写
//func (m *Mysql) Update(data []map[string]interface{}) (sql.Result, error) {
//	//UPDATE persondata SET age=age+1 where id=...
//	return nil, nil
//}

// 批量更新，过滤不存在的列，使用多次更新，比较影响性能
func (m *Mysql) UpdateCasual(data []map[string]interface{}) *Mysql {
    res,err:= m.updateSqlCasual(data)
    m.Results = res
    for _,v:= range err{
    	if v!=nil{
    		m.Errors = append(m.Errors,v)
		}
	}
	return m
}

//
func (m *Mysql) updateSqlCasual(data []map[string]interface{}) ([]sql.Result,[]error) {
	cif := m.GetTableInfo().Info // 获取表的列结构
	var res []sql.Result
	var ers []error
	var kvSql string
	for _,v := range data{
		// 遍历每一行，每一行都是一次更新
		kvSql = ``
		for key, value := range v {
			// 拼接set后的字符串  a=1,b='2',c=11
			if attr, ok := cif[key]; ok { // 如果表中存在这个字段
				// 存在这个字段
				if isString(attr.Type) {
					kvSql = kvSql + key + `='` + convertToString(value) + `',`
				} else {
					kvSql = kvSql + key + `=` + convertToString(value) + `,`
				}
			} else {
				// 如果不存在，跳过
				continue
			}
		}
		kvSql = strings.TrimRight(kvSql, `,`) // 去掉最后的逗号
		m.sql = `UPDATE ` + m.tableName + ` SET ` + kvSql + ` ` + m.whereSql   // 拼接更新sql
		r,err:= m.Exec(m.sql)
		res= append(res,r)
		ers = append(ers,err)
	}
	return res,ers
}