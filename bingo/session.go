package bingo

import (
	"github.com/gorilla/sessions"
	"io"
	"encoding/base64"
	"crypto/rand"
	"net/http"
)

// 使用session
//var a Context

type Session struct {
	session *sessions.Session
	writer http.ResponseWriter
	req *http.Request
}

var globalSession *sessions.CookieStore

// 初始化全局session存储器
func init() {
	globalSession = sessions.NewCookieStore([]byte(randId()))
}

func randId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}


func (s *Session) Set(key string,value interface{}) error {
	// 设置值
	s.session.Values[key] = value
	// 保存session
	err:=s.session.Save(s.req,s.writer)
	return err
}

func (s *Session) Get(key string) interface{}  {
	return  s.session.Values[key]
}