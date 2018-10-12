package bingo

import (
	"sync"
	"testing"
)

func Test_env_Set(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		submap  map[string]string
	}
	type args struct {
		k string
		v string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "case",
			fields: struct {
				RWMutex sync.RWMutex
				submap  map[string]string
			}{RWMutex: sync.RWMutex{}, submap: make(map[string]string)},
			args: struct {
				k string
				v string
			}{k: "key", v: "value"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &env{
				RWMutex: tt.fields.RWMutex,
				submap:  tt.fields.submap,
			}
			if got := e.Set(tt.args.k, tt.args.v); got != tt.want {
				t.Errorf("env.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_env_Get(t *testing.T) {
	type fields struct {
		RWMutex sync.RWMutex
		submap  map[string]string
	}
	type args struct {
		k string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "case",
			fields: struct {
				RWMutex sync.RWMutex
				submap  map[string]string
			}{RWMutex: sync.RWMutex{}, submap: make(map[string]string)},
			args: struct{ k string }{k: "key"},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &env{
				RWMutex: tt.fields.RWMutex,
				submap:  tt.fields.submap,
			}

			e.Set("key", "test")

			if got := e.Get(tt.args.k); got != tt.want {
				t.Errorf("env.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
