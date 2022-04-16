package queue

import (
	"github.com/andewx/dieselfluid/compute/common"
)

type Queue interface {
	Push(common.Evaluatable)
	Pop() (common.Evaluatable, error)
	Empty() bool
}

type ComputeQueue struct {
	jobs  []common.ComputeFunction
	index int
	empty bool
}

//New ComputeQueue with queue array size = (size)
func New(size int) ComputeQueue {
	return ComputeQueue{make([]common.ComputeFunction, size), 0, true}
}

//Pushes evaluatbale element onto the jobs array and increments index
func (q *ComputeQueue) Push(m common.ComputeFunction) {

	if q.empty {
		q.empty = false
		q.jobs[0] = m
	} else {
		q.index++
		if q.index >= len(q.jobs) {
			q.jobs = append(q.jobs, m)
		} else {
			q.jobs[q.index] = m
		}
		q.empty = false
	}
}

//Pops current element off of array should always check error
func (q *ComputeQueue) Pop() common.ComputeFunction {
	if q.index == 0 {
		if q.empty {
			var f = common.ComputeFunction{}
			return f
		}
		m := q.jobs[0]
		q.jobs[0] = common.ComputeFunction{}
		q.empty = true
		return m
	}
	m := q.jobs[q.index]
	q.index -= 1
	return m
}

func (q *ComputeQueue) Empty() bool {
	return q.empty
}
