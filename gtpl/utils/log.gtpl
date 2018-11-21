package utils

import "github.com/silsuer/bingo-log"

// 日志
var Log *bingo_log.Log

func init() {
	// 日志
	Log = bingo_log.NewLog(bingo_log.LogSyncMode)
	config := make(map[string]string)
	config["root"] = "."
	config["format"] = "2006_01_02"
	conn := bingo_log.NewKirinConnector(config)
	Log.LoadConnector(conn)
}


