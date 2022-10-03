package solver

type SPHMethod interface {
	Run()            //Standard Run
	Run_(t chan int) //Threaded Handler
}
