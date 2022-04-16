package queue

import (
	"fmt"
	"testing"

	"github.com/andewx/dieselfluid/compute/common"
	"github.com/andewx/dieselfluid/math/mgl"
)

func TestQueue(t *testing.T) {
	job := common.ComputeFunction{}
	queue := New(10)
	job.Evaluate = func(a mgl.Vec) mgl.Vec {
		fmt.Printf("Eval %s\n", a.ToString())
		return a.Norm()
	}

	fmt.Printf("------------Testing Queue--------\n")
	queue.Push(job)
	if queue.Empty() {
		t.Errorf("Queue has jobs -- returns Empty:true\n")
	}

	avec := mgl.Vec{1.0, 2.0, 3.0}
	f := queue.Pop()
	b := f.Evaluate(avec)
	fmt.Printf("%s\n", b.ToString())
	queue.Pop()
	f.Evaluate(avec)
	queue.Pop()
	f.Evaluate(avec)
	if !queue.Empty() {
		t.Errorf("Queue not empty\n")
	}

}
