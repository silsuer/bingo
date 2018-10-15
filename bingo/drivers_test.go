package bingo

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDriverInit(t *testing.T) {

	// 测试一个不存在的驱动
	Env.Set("DB_DRIVER", "TEST")
	DriverInit()
	//fmt.Println(Driver)

	// 测试一个存在的驱动
	Env.Set("DB_DRIVER", "MYSQL")
	Env.Set("DB_USERNAME", "root")
	Env.Set("DB_PASSWORD", "bingo")
	Env.Set("DB_HOST", "127.0.0.1")
	Env.Set("DB_PORT", "3306")
	Env.Set("DB_NAME", "bingo")
	Env.Set("DB_CHARSET", "utf8")

	DriverInit()
	//fmt.Println(Driver)
}

