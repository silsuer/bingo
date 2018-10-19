package bingo

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io"
	"io/ioutil"
	"sync"
	"syscall"
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

	res2, err := http.Get(s.URL + "/test")

	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	body2, err := ioutil.ReadAll(res2.Body)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(body))
	t.Log(string(body2))
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
				"name", "path=currentPath", "file=test=1",
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

func TestBingo_startServer(t *testing.T) {
	rr := NewRoute().Get("/").Target(func(c *Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "application/json")
		io.WriteString(c.Writer, `{"message":"Hello Bingo!"}`)
	}).Register()

	r := New()

	r.Handle("GET", "/", rr)

	b := Bingo{}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		syscall.Kill(12345, syscall.SIGHUP)
		go b.startServer([]string{"dev"}, ":12345", r)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		syscall.Kill(12343, syscall.SIGHUP)
		go b.startServer([]string{"daemon", "start"}, ":12343", r)
		wg.Done()
	}()
	wg.Wait()
}
