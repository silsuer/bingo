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
	"net/http/httptest"
	"net/http"
	"io"
	"io/ioutil"
)

func TestBingo_Run(t *testing.T) {

	rr := NewRoute().Get("/").Target(func(c *Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "application/json")
		io.WriteString(c.Writer, `{"message":"Hello Bingo!"}`)
	}).Register()

	r := New()

	r.Handle("GET", "/", rr)

	s := httptest.NewServer(r)
	defer s.Close()

	res, err := http.Get(s.URL)

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(body))
}

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
