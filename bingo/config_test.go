package bingo

import (
	"sync"
	"testing"
)

func TestConfigMap_Set(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		data      map[string]string
		dataSlice map[string][]string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "case1",
			fields: struct {
				RWMutex   sync.RWMutex
				data      map[string]string
				dataSlice map[string][]string
			}{RWMutex: sync.RWMutex{}, data: make(map[string]string), dataSlice: make(map[string][]string)},
			args: struct {
				key   string
				value string
			}{key: "name", value: "value"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigMap{
				RWMutex:   tt.fields.RWMutex,
				data:      tt.fields.data,
				dataSlice: tt.fields.dataSlice,
			}
			if got := c.Set(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("ConfigMap.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigMap_SetSlice(t *testing.T) {
	type fields struct {
		RWMutex   sync.RWMutex
		data      map[string]string
		dataSlice map[string][]string
	}
	type args struct {
		key   string
		value []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "case1",
			fields: struct {
				RWMutex   sync.RWMutex
				data      map[string]string
				dataSlice map[string][]string
			}{RWMutex: sync.RWMutex{}, data: make(map[string]string), dataSlice: make(map[string][]string)},
			args: struct {
				key   string
				value []string
			}{key: "key", value: []string{"value1", "value2"}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigMap{
				RWMutex:   tt.fields.RWMutex,
				data:      tt.fields.data,
				dataSlice: tt.fields.dataSlice,
			}
			if got := c.SetSlice(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("ConfigMap.SetSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
