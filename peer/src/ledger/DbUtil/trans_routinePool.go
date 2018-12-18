package DbUtil

import (
	"fmt"
)

type TransRoutinePool struct {
	Queue         chan func() error
	Number        int // the number of thread
	Total         int
	result        chan error
	finshCallBack func()
}

func (self *TransRoutinePool) Init(number int, total int) {
	self.Queue = make(chan func() error, total)
	self.Number = number
	self.Total = total
	self.result = make(chan error, total)
}
func (self *TransRoutinePool) Add_transaction() {
	self.Total = self.Total + 1
	self.Queue = make(chan func() error, self.Total)
	self.result = make(chan error, self.Total)
}

func (self *TransRoutinePool) Start() {
	for i := 0; i < self.Number; i++ {
		go func() {
			for {
				fmt.Println("self.Queue", self.Queue)
				task, ok := <-self.Queue
				if !ok {
					break
				}
				err := task() //  Op task
				self.result <- err
			}
		}()
	}
	for j := 0; j < 1; j++ {
		//	for j := 0; j < self.Total; j++ {
		//	for j := ; j<receive_pool_total; j++{
		res, ok := <-self.result
		if !ok {
			break
		}
		if res != nil {
			fmt.Println(res) //  only for print the result
		}
	}
	if self.finshCallBack != nil {
		self.finshCallBack()
	}
}
func (self *TransRoutinePool) Stop() {
	close(self.Queue)
	close(self.result)
}
func (self *TransRoutinePool) AddTask(task func() error) {

	self.Queue <- task
}
func (self *TransRoutinePool) SetFinishCallback(callback func()) {
	self.finshCallBack = callback
}
