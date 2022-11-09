package group

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestGroupData_Add_Get_1(t *testing.T) {
	tests := []struct {
		name  string
		group *GroupData[string]
		args  []string
	}{
		{
			name:  "test 1",
			group: NewGroupData[string](context.TODO()),
			args:  []string{"1"},
		},
		{
			name:  "test 2",
			group: NewGroupData[string](context.TODO()),
			args:  []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, arg := range tt.args {
				tt.group.Add(arg)
			}
			rs := tt.group.Get(time.Second)
			if len(rs) != len(tt.args) {
				t.Errorf("Get Value")
			}
		})
	}
}

func TestGroupData_Add_Get_2(t *testing.T) {
	tests := []struct {
		name  string
		group *GroupData[string]
		args  []string
	}{
		{
			name:  "test 3",
			group: NewGroupData[string](context.TODO()),
			args:  []string{"1"},
		},
		{
			name:  "test 4",
			group: NewGroupData[string](context.TODO()),
			args:  []string{"1", "2", "3"},
		},
	}
	wg := &sync.WaitGroup{}
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			go func(g *GroupData[string], args []string, wg *sync.WaitGroup) {
				defer wg.Done()
				rs := g.Get(time.Second)
				if len(rs) != len(args) {
					t.Errorf("Get Value")
				}
			}(tt.group, tt.args, wg)
			for _, arg := range tt.args {
				tt.group.Add(arg)
			}

		})
	}
	wg.Wait()
}

func TestGroupData_New_Stop(t *testing.T) {
	g := NewGroupData[string](context.TODO())
	g.Stop()
	if !isChanClosed(g.queue) {
		t.Errorf("Not Stop")
	}
}

func isChanClosed(ch interface{}) bool {
	if reflect.TypeOf(ch).Kind() != reflect.Chan {
		return false
	}

	// get interface value pointer, from cgo_export
	// typedef struct { void *t; void *v; } GoInterface;
	// then get channel real pointer
	cptr := *(*uintptr)(unsafe.Pointer(
		unsafe.Pointer(uintptr(unsafe.Pointer(&ch)) + unsafe.Sizeof(uint(0))),
	))

	// this function will return true if chan.closed > 0
	// see hchan on https://github.com/golang/go/blob/master/src/runtime/chan.go
	// type hchan struct {
	// qcount   uint           // total data in the queue
	// dataqsiz uint           // size of the circular queue
	// buf      unsafe.Pointer // points to an array of dataqsiz elements
	// elemsize uint16
	// closed   uint32
	// **

	cptr += unsafe.Sizeof(uint(0)) * 2
	cptr += unsafe.Sizeof(unsafe.Pointer(uintptr(0)))
	cptr += unsafe.Sizeof(uint16(0))
	return *(*uint32)(unsafe.Pointer(cptr)) > 0
}
