//Generated TestConnectionManager_add
//Generated TestConnectionManager_done
//Generated TestConnectionManager_close
//Generated TestConnectionManager_rmIdleConns
//Generated TestConnectionManager_addIdleConns
//Generated TestGracefulServe
//Generated Test_newConnectionManager
//Generated TestGraceful
//Generated Test_server_getListener
//Generated Test_server_fork
//Generated Test_server_shutdown
//Generated Test_isRunning
//Generated Test_savePid
//Generated Test_handleSignals
//Generated TestDaemonInit
//Generated Test_forkDaemon
//Generated Test_restart
//Generated Test_listener_File
package bingo

import (
	"net"
	"sync"
	"testing"
	"time"
	"reflect"
)

func TestConnectionManager_add(t *testing.T) {
	type fields struct {
		WaitGroup sync.WaitGroup
		Counter   int
		mux       sync.Mutex
		idleConns map[string]net.Conn
	}
	type args struct {
		delta int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "case",
			fields: struct {
				WaitGroup sync.WaitGroup
				Counter   int
				mux       sync.Mutex
				idleConns map[string]net.Conn
			}{WaitGroup: sync.WaitGroup{}, Counter: 0, mux: sync.Mutex{}, idleConns: make(map[string]net.Conn)},
			args: struct{ delta int }{delta: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ConnectionManager{
				WaitGroup: tt.fields.WaitGroup,
				Counter:   tt.fields.Counter,
				mux:       tt.fields.mux,
				idleConns: tt.fields.idleConns,
			}
			cm.add(tt.args.delta)
		})
	}
}

func TestConnectionManager_done(t *testing.T) {
	type fields struct {
		WaitGroup sync.WaitGroup
		Counter   int
		mux       sync.Mutex
		idleConns map[string]net.Conn
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "case",
			fields: struct {
				WaitGroup sync.WaitGroup
				Counter   int
				mux       sync.Mutex
				idleConns map[string]net.Conn
			}{WaitGroup: sync.WaitGroup{}, Counter: 0, mux: sync.Mutex{}, idleConns: make(map[string]net.Conn)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ConnectionManager{
				WaitGroup: tt.fields.WaitGroup,
				Counter:   tt.fields.Counter,
				mux:       tt.fields.mux,
				idleConns: tt.fields.idleConns,
			}

			cm.add(1)
			cm.done()
		})
	}
}

func TestConnectionManager_close(t *testing.T) {
	type fields struct {
		WaitGroup sync.WaitGroup
		Counter   int
		mux       sync.Mutex
		idleConns map[string]net.Conn
	}
	type args struct {
		t time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "case",
			fields: struct {
				WaitGroup sync.WaitGroup
				Counter   int
				mux       sync.Mutex
				idleConns map[string]net.Conn
			}{WaitGroup: sync.WaitGroup{}, Counter: 0, mux: sync.Mutex{}, idleConns: make(map[string]net.Conn)},
			args: struct{ t time.Duration }{t: time.Second},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ConnectionManager{
				WaitGroup: tt.fields.WaitGroup,
				Counter:   tt.fields.Counter,
				mux:       tt.fields.mux,
				idleConns: tt.fields.idleConns,
			}
			//cm.add(1)
			//cm.done()
			cm.close(tt.args.t)
		})
	}
}

func TestConnectionManager_rmIdleConns(t *testing.T) {
	type fields struct {
		WaitGroup sync.WaitGroup
		Counter   int
		mux       sync.Mutex
		idleConns map[string]net.Conn
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "case",
			fields: struct {
				WaitGroup sync.WaitGroup
				Counter   int
				mux       sync.Mutex
				idleConns map[string]net.Conn
			}{WaitGroup: sync.WaitGroup{}, Counter: 0, mux: sync.Mutex{}, idleConns: make(map[string]net.Conn)},
			args: struct{ key string }{key: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ConnectionManager{
				WaitGroup: tt.fields.WaitGroup,
				Counter:   tt.fields.Counter,
				mux:       tt.fields.mux,
				idleConns: tt.fields.idleConns,
			}
			cm.addIdleConns("test", nil)
			cm.rmIdleConns(tt.args.key)
		})
	}
}

func TestConnectionManager_addIdleConns(t *testing.T) {
	type fields struct {
		WaitGroup sync.WaitGroup
		Counter   int
		mux       sync.Mutex
		idleConns map[string]net.Conn
	}
	type args struct {
		key  string
		conn net.Conn
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "case",
			fields: struct {
				WaitGroup sync.WaitGroup
				Counter   int
				mux       sync.Mutex
				idleConns map[string]net.Conn
			}{WaitGroup: sync.WaitGroup{}, Counter: 0, mux: sync.Mutex{}, idleConns: make(map[string]net.Conn)},
			args: struct {
				key  string
				conn net.Conn
			}{key: "test", conn: nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &ConnectionManager{
				WaitGroup: tt.fields.WaitGroup,
				Counter:   tt.fields.Counter,
				mux:       tt.fields.mux,
				idleConns: tt.fields.idleConns,
			}
			cm.addIdleConns(tt.args.key, tt.args.conn)
		})
	}
}

func Test_newConnectionManager(t *testing.T) {
	tests := []struct {
		name string
		want *ConnectionManager
	}{
		{
			name: "case",
			want: &ConnectionManager{WaitGroup: sync.WaitGroup{}, Counter: 0, mux: sync.Mutex{}, idleConns: make(map[string]net.Conn)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newConnectionManager(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConnectionManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

