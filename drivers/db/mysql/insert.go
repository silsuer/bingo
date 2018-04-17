package mysql

import (
	"errors"
	"strings"
	"strconv"
	"database/sql"
)

// mysql数据库的插入操作
// 插入时要缓存表结构,略微影响表的性能
// 1. 插入
// 2. 插入返回id
// 3. 插入忽略多余的值

type colInfo struct {
	Name string // 列名
	Type string // 列类型
}

type tableInfo struct {
	Name string             // 表名
	Info map[string]colInfo // 表的所有列信息
}

var TableInfo map[string]tableInfo // 所有表的结构信息

// 向数据库中插入一条或多条数据,data键是列名，值是列数据
func (m *Mysql) InsertOne(data map[string]interface{}) *Mysql {
	var d []map[string]interface{}
	d = append(d, data)
	m.insertSql(d)
	res, err := m.Exec(m.sql)
	 m.Result = res
	if err != nil {
		m.Errors = append(m.Errors,err)
	}
	return m
}

// 插入多行
func (m *Mysql) Insert(data []map[string]interface{}) *Mysql {
	m.insertSql(data) // 拼接插入语句
	// 获取表结构
	res,err:= m.Exec(m.sql) // 执行语句并返回
	m.Result = res
	if err != nil {
		m.Errors = append(m.Errors,err)
	}
	return m
}

// 插入时如果有不存在的列，会忽略掉 #随便插插 :)
func (m *Mysql) InsertOneCasual(data map[string]interface{}) *Mysql {
	var d []map[string]interface{}
	d = append(d, data)
	res, err := m.insertCasualSql(d)
	m.Result = res[0]
	if err[0]!=nil {
		m.Errors = append(m.Errors,err[0])
	}
	m.Errors = append(m.Errors,err[0]) // 只插入一条，返回第一个
	return m
}

func (m *Mysql) InsertCasual(data []map[string]interface{}) *Mysql {
	res,err:= m.insertCasualSql(data)
	m.Results = res
	if len(err)!=0{
		for _,v := range err{
			m.Errors = append(m.Errors,v)
		}
	}
	return m
}

// 拼接随便插插的语句....
func (m *Mysql) insertCasualSql(data []map[string]interface{}) ([]sql.Result, []error) {
	cif := m.GetTableInfo().Info // 获取表的列结构
	var res []sql.Result
	var ers []error
	var keySql, valueSql string
	for _, v := range data {
		// 置空
		keySql = ``
		valueSql = ``
		// 再进行一次遍历，拼凑key，由于每个map中映射的个数不同（即插入的列不同），这里只能采用循环插入的办法
		for key, value := range v {
			if attr, ok := cif[key]; ok {
				keySql = keySql + key + `,`
				// 存在这个列，开始判断这个列的属性
				if isString(attr.Type) {
					// 如果是字符串类型，插入时要加引号,目前只能插入int和string类型的值
					valueSql = valueSql + `'` + convertToString(value) + `',`
				} else {
					valueSql = valueSql + convertToString(value) + `,`
				}
			} else {
				continue
			}
		}
		keySql = strings.TrimRight(keySql, `,`)
		valueSql = strings.TrimRight(valueSql, `,`)
		sql := `INSERT INTO ` + m.tableName + ` (` + keySql + ` ) VALUES ( ` + valueSql + `)`
		re, err := m.Exec(sql)
		res = append(res, re)
		ers = append(ers, err)
	}
	return res, ers
}

func (m *Mysql) insertSql(data []map[string]interface{}) error {
	cif := m.GetTableInfo().Info // 获取表的列结构
	var keySql, valueSql string
	for _, v := range data {
		// 置空
		keySql = ``
		// 再进行一次遍历，拼凑key
		for key, value := range v {
			keySql = keySql + key + `,`
			valueSql = valueSql + `(`
			if attr, ok := cif[key]; ok {

				// 存在这个列，开始判断这个列的属性
				if isString(attr.Type) {
					// 如果是字符串类型，插入时要加引号,目前只能插入int和string类型的值
					valueSql = valueSql + `'` + convertToString(value) + `',`
				} else {
					valueSql = valueSql + convertToString(value) + `,`
				}
				valueSql = strings.TrimRight(valueSql, `,`)
			} else {
				return errors.New("cannot find a column named " + key + " in " + m.tableName + " table ")
			}
			valueSql = valueSql + `),`
		}
	}
	keySql = strings.TrimRight(keySql, `,`)
	valueSql = strings.TrimRight(valueSql, `,`)
	// 拼接sql字符串
	m.sql = `INSERT INTO ` + m.tableName + ` (` + keySql + ` ) VALUES  ` + valueSql
	return nil
}

// 获取表结构信息并缓存
func (m *Mysql) GetTableInfo() *tableInfo {
	if info, ok := TableInfo[m.tableName]; ok {
		// 这张表的信息已经缓存
		return &info
	} else {
		// 未缓存，从数据库中获取表信息，并缓存
		m.UpdateTableSchema()
		inf := TableInfo[m.tableName]
		return &inf
	}
	return nil
}

func (m *Mysql) UpdateTableSchema() {
	// 初始化
	if TableInfo == nil {
		TableInfo = make(map[string]tableInfo)
	}

	// 未缓存，从数据库中获取表信息，并缓存
	sql := `SELECT  COLUMN_NAME,DATA_TYPE  FROM information_schema.COLUMNS where TABLE_NAME='` + m.tableName + `' AND TABLE_SCHEMA='` + m.TableSchema + `'`
	inf, err := m.Query(sql)
	Check(err)
	var tf tableInfo
	tf.Name = m.tableName
	tfCols := make(map[string]colInfo)
	for inf.Next() {
		var cname, ctype string
		inf.Scan(&cname, &ctype) // 获取字段名称和类型
		var newCol = &colInfo{Name: cname, Type: ctype}
		tfCols[cname] = *newCol
	}
	tf.Info = tfCols
	TableInfo[m.tableName] = tf // 更新并赋值
}

// 判断是否是字符串类型
func isString(str string) bool {
	switch strings.ToUpper(str) {
	case "TINYINT":
	case "SMALLINT":
	case "MEDIUMINT":
	case "INT":
	case "INTEGER":
	case "BIGINT":
	case "FLOAT":
	case "DOUBLE":
	case "DECIMAL":
		return false
	default:
		return true
	}
	return false
}

// 把数据转换为字符串
func convertToString(m interface{}) string {
	switch m.(type) {
	case int64:
	case int:
		return strconv.Itoa(m.(int))
		break
	default:
		return m.(string)
	}
	return "111"
}
