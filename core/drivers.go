package core

import (
	"strings"
	"bingo/drivers/db/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// 数据库驱动
type driver struct {
	name     string // 驱动名
	dbConfig string // 配置
}

var Driver *driver

// 驱动初始化
func DriverInit() {
	Driver = &driver{}
	Driver.name = strings.ToUpper(Env.Get("DB_DRIVER"))
	switch Driver.name {
	case "MYSQL":
		// 初始化了驱动之后，开始初始化数据库连接
		Driver.dbConfig = Env.Get("DB_USERNAME") + ":" + Env.Get("DB_PASSWORD") + "@tcp(" + Env.Get("DB_HOST") + ":" + Env.Get("DB_PORT") + ")" + "/" + Env.Get("DB_NAME") + "?" + "charset=" + Env.Get("DB_CHARSET")
		break
	default:
		break
	}

}

// 根据数据库驱动，获取数据库连接
func (d *driver) GetConnection() interface{} {
	switch d.name {
	case "MYSQL":
		m := mysql.Mysql{}                // 实例化结构体
		m.Init(Driver.dbConfig) // 设置表名和数据库连接
		return &m                         // 返回实例
		break
	default:
		break
	}
	return nil
}
