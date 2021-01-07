package wcsph

import (
	"dslfluid.com/dsl/model"
	"dslfluid.com/dsl/model/sph"
)

type WCSPH struct {
	core sph.SPHCore
}

//---------SPHCore Run Methods------------------------//
func (p WCSPH) Run() {
	done := false

	for !done {
		p.core.ComputeDensity()
		p.core.AccumulatePressure()
		p.core.Accumulate_NonPressure()
		p.core.Collide()
		p.core.Update()
		p.core.UpdateTime()
	}

	return
}

//Run_ Executes SPH Loop in Thread Blocking I/O Manner. If an application
//Needs exclusive resource access to SPHCore data structures they should pass
//model.THREAD_WAIT to block thread execution. When access is no longer required
//the method should pass model.THREAD_GO to the specified channel
//Buffer access to SPHCore go slices should be read only access, other wise for thread safe
//Execution THREAD_WAIT should be called if modifying buffers or relying on temporal coherence
//for volatile data buffers
func (p WCSPH) Run_(t chan int) {
	done := false
	sync := true

	for !done {

		//Executes full in frame computation loop.core.
		if sync {
			p.core.ComputeDensity()
			p.core.AccumulatePressure()
			p.core.Accumulate_NonPressure()
			p.core.Collide()
			p.core.Update()
			p.core.UpdateTime()

			//Channel Monitor - Monitor Blocking I/O Request
			status := <-t
			if status == model.THREAD_WAIT {
				sync = false
				t <- model.SPH_THREAD_WAITING
				waitStatus := <-t
				if waitStatus == model.THREAD_GO {
					sync = true
				}
			}
			//End Thread Blocking
		}

		//Handle Thread Block Outside Initial Message
		waitStatus := <-t
		if waitStatus == model.THREAD_GO {
			sync = true
		}

	}

	return
}