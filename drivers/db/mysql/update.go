package mysql

import (
	"database/sql"
	"fmt"
)

// 对于mysql数据库的更新操作

// 更新一条数据
func (m *Mysql) UpdateOne(data map[string]interface{}) (sql.Result, error) {
	//UPDATE persondata SET age=age+1 where id=...
	//cif:=m.GetTableInfo().Info  // 获取表的数据结构
	fmt.Println(m)

	return nil, nil
}

//func (m *Mysql) updateSql(data []map[string]interface{}) error {
//	cif := m.GetTableInfo().Info
//	kvSql := ``
//	for _, v := range data {
//		for key, value := range v {
//			// 拼接set后的字符串  a=1,b='2',c=11
//			if attr, ok := cif[key]; ok {
//				// 存在这个字段
//				if isString(attr.Type) {
//					kvSql
//				}else{
//
//				}
//			} else {
//				return errors.New("cannot find a column named " + key + " in " + m.tableName + " table ")
//			}
//		}
//	}
//	m.sql = `UPDATE ` + m.tableName + ` SET ` + kvSql + ` ` + m.whereSql
//}

// 更新一条数据（过滤不存在的列）
func (m *Mysql) UpdateOneCasual(data map[string]interface{}) (sql.Result, error) {
	//UPDATE persondata SET age=age+1 where id=...
	return nil, nil
}

// 批量更新
func (m *Mysql) Update(data []map[string]interface{}) (sql.Result, error) {
	//UPDATE persondata SET age=age+1 where id=...
	return nil, nil
}

// 批量更新，过滤不存在的列
func (m *Mysql) UpdateCasual(data []map[string]interface{}) ([]sql.Result, []error) {
	//UPDATE persondata SET age=age+1 where id=...
	return nil, nil
}
