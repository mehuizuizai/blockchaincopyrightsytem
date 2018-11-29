package jobmgr

//	"fmt"

//通用的Job接口
type Job interface {
	DoWork()
}

//Job handler
type JobHandler func(Job)

type worker struct {
	id         int      //worker的id
	jobChannel chan Job //接收job的channnel
	jobmgr     *JobMgr  //job 管理器
	//handler    JobHandler //Job处理句柄
	quitChan chan bool //接收退出信号的channel
}

//Job管理器
type JobMgr struct {
	jobQueue    chan Job     //待处理的Job队列
	idleWorkers chan *worker //空闲的worker
	numWorker   int          //worker的数量
	//handler     JobHandler   //Job处理句柄
}

func (w *worker) start() {
	go func() {
		for {
			//fmt.Println("------------", w.id)
			w.jobmgr.idleWorkers <- w //将自己放入空闲的worker池里
			select {
			case job := <-w.jobChannel:
				job.DoWork() //调用job 处理句柄
			case <-w.quitChan:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *worker) stop() {
	w.quitChan <- true
}

func NewJobMgr(numCachedJob int, numWorker int) *JobMgr {
	mgr := &JobMgr{}
	mgr.jobQueue = make(chan Job, numCachedJob)
	mgr.idleWorkers = make(chan *worker, numWorker)
	mgr.numWorker = numWorker
	//mgr.handler = handler
	return mgr
}

//将job放入Job缓冲队列
func (mgr *JobMgr) PutJob(job Job) {
	mgr.jobQueue <- job
}

//启动Job Manager
func (mgr *JobMgr) Start() {
	for i := 1; i <= mgr.numWorker; i++ {
		worker := &worker{
			id:         i,
			jobChannel: make(chan Job),
			jobmgr:     mgr,
			quitChan:   make(chan bool),
			//handler:    mgr.handler,
		}
		worker.start()
	}
	go func(mgr *JobMgr) {
		for {
			select {
			case job := <-mgr.jobQueue:
				worker := <-mgr.idleWorkers
				worker.jobChannel <- job
			}
		}
	}(mgr)
}

//停止Job Manager
func (mgr *JobMgr) Stop() {

}
