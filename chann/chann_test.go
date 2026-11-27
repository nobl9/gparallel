package chann_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/nobl9/gparallel/chann"
)

// TestCloseDeliversEveryValueToASlowReceiver verifies that an unbounded
// channel drops nothing when it is closed before the receiver is ready.
func TestCloseDeliversEveryValueToASlowReceiver(t *testing.T) {
	const count = 100
	const timeForCloseToTerminate = 100 * time.Millisecond

	ch := chann.New[int]()
	want := make([]int, count)
	for i := range want {
		want[i] = i
		ch.In() <- i
	}
	ch.Close()
	time.Sleep(timeForCloseToTerminate)

	var got []int
	for v := range ch.Out() {
		got = append(got, v)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("received %d of %d values: %v", len(got), count, got)
	}
}
