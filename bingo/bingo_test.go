//Generated TestBingo_Run
//Generated TestBingo_setGlobalParamFromArgs
//Generated TestBingo_startServer
//Generated Test_startDevServer
//Generated Test_startWatchServer
//Generated Test_restartDaemonServer
//Generated Test_listeningWatcherDir
//Generated Test_startDaemonServer
package bingo

import (
	"testing"
)
//
//func TestBingo_Run(t *testing.T) {
//
//	b := Bingo{}
//	//b.Run(":12345")
//
//	NewRoute().Get("/").Target(func(c *Context) {
//		c.Writer.WriteHeader(http.StatusOK)
//		c.Writer.Header().Set("Content-Type", "application/json")
//		io.WriteString(c.Writer, `{"message":"Hello Bingo!"}`)
//	}).Register()
//
//	// 开启服务，并且不阻塞主进程
//	go func() {
//		b.Run(":12345")
//	}()
//
//	//time.Sleep(1 * time.Second)
//	resp, err := http.Get("http://localhost:12345")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	t.Log(string(body))
//}

func TestBingo_setGlobalParamFromArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		b    *Bingo
		args args
	}{
		{
			name: "case1",
			b:    &Bingo{},
			args: struct{ args []string }{args:
			[]string{
				"name", "path=currentPath", "file=test",
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bingo{}
			b.setGlobalParamFromArgs(tt.args.args)
		})
	}
}

