package fixedThreadPool

import (
	"fmt"
	"github.com/bilibili-convert-test/pkg/util"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// FixedThreadPool 线程池
type FixedThreadPool struct {
	*TotalProgress         // 总体进程
	Tasks          []*Task // 所有任务
	GoroutineNum   int     // 线程数
	IsVisible      bool    // 是否显示进度条
}

// NewFixedThreadPool 创建线程池
func NewFixedThreadPool(tasks []*Task, goroutineNum int, isVisible bool) *FixedThreadPool {
	return &FixedThreadPool{
		Tasks:         tasks,
		TotalProgress: NewTotalProgress(len(tasks), goroutineNum),
		GoroutineNum:  goroutineNum,
		IsVisible:     isVisible,
	}
}

// TotalProgress 总进度
type TotalProgress struct {
	FinishNum   int                  // 完成数
	TaskNum     int                  // 总任务数
	AllProgress []*GoroutineProgress // 所有goroutine的进度
}

// NewTotalProgress 初始化总进度
func NewTotalProgress(taskNum, goroutineNum int) *TotalProgress {
	return &TotalProgress{
		FinishNum:   0,
		TaskNum:     taskNum,
		AllProgress: make([]*GoroutineProgress, goroutineNum),
	}
}

// GoroutineProgress 每个goroutine的进度
type GoroutineProgress struct {
	GoroutineId int  // goroutine的ID
	*Task            // 正在执行的任务
	isIdle      bool // 当前线程是否空闲
}

// Exec 执行任务
func (goProgress *GoroutineProgress) Exec() {
	goProgress.isIdle = false
	goProgress.Run(goProgress.TaskProgress) // 执行任务
	goProgress.isIdle = true
}

// Task 每个goroutine所执行的任务
type Task struct {
	TaskProgress *TaskProgress
	Run          func(*TaskProgress) // 执行
}

func (t Task) GetProgress() int {
	return t.TaskProgress.CurrentProgress
}

func NewTask(f func(*TaskProgress)) *Task {
	return &Task{
		TaskProgress: &TaskProgress{CurrentProgress: 0},
		Run:          f,
	}
}

// TaskProgress 任务的进度
type TaskProgress struct {
	CurrentProgress int // 任务进度单位%
}

// Exec 线程池开始执行任务
func (pool *FixedThreadPool) Exec() {
	// 将所有任务放入到channel中
	taskChannel := make(chan *Task, pool.TaskNum)
	for _, task := range pool.Tasks {
		taskChannel <- task
	}
	// 同步主线程与创建的goroutine
	wg := &sync.WaitGroup{}
	wg.Add(pool.TaskNum)
	// 通知线程关闭的channel,+1是为了关闭打印进度的线程
	done := make(chan struct{}, pool.GoroutineNum+1)

	// 开启goroutine线程
	for i := 0; i < pool.GoroutineNum; i++ {
		// 监听每个线程的进度
		progress := &GoroutineProgress{
			GoroutineId: i,
		}
		pool.TotalProgress.AllProgress[i] = progress

		// 开启线程
		go func(done chan struct{}, wg *sync.WaitGroup, progress *GoroutineProgress, totalProgress *TotalProgress) {
		Label:
			for {
				select {
				case progress.Task = <-taskChannel:
					progress.Exec()
					wg.Done()
					totalProgress.FinishNum++
				case <-done:
					// 结束
					break Label
				}
			}
		}(done, wg, progress, pool.TotalProgress)
	}

	if pool.IsVisible {
		// 单独开启一个监听正在执行的任务的进度的goroutine
		go func(done chan struct{}, totalProgress *TotalProgress) {
		Label:
			for {
				select {
				case <-done:
					break Label
				default:
					// 输出任务的进度
					fmt.Printf("当前进度: %d/%d\n", totalProgress.FinishNum, totalProgress.TaskNum)
					for _, progress := range totalProgress.AllProgress {
						if !progress.isIdle && progress.Task != nil {
							fmt.Printf("任务进度: [%s%s] %d%%\n", strings.Repeat("=", progress.GetProgress()), strings.Repeat(" ", 100-progress.GetProgress()), progress.GetProgress())
						}
					}
					// 间隔1s
					time.Sleep(time.Second)
					// 清屏
					fmt.Println("\033c")
				}
			}
		}(done, pool.TotalProgress)
	}

	// 等待所有任务执行结束
	wg.Wait()
	for i := 0; i <= pool.GoroutineNum; i++ {
		done <- struct{}{}
	}
	log.Println("所有任务执行完毕")
}

type ThreadPoolWithProgressReader struct {
	*TaskProgress
	io.Reader
	FileName    string
	TotalSize   int64
	CurrentSize int64
}

func NewThreadPoolWithProgressReader(progress *TaskProgress) *ThreadPoolWithProgressReader {
	return &ThreadPoolWithProgressReader{
		TaskProgress: progress,
	}
}

func (reader *ThreadPoolWithProgressReader) Read(p []byte) (int, error) {
	n, err := reader.Reader.Read(p)
	if err == io.EOF {
		log.Println("下载完成")
		return n, err
	}
	if err != nil {
		log.Fatal(err)
	}
	// 读的时候所作的操作
	reader.updateCurrentSize(n)
	return n, err
}

func (reader *ThreadPoolWithProgressReader) Init(resp *http.Response, file *os.File) {
	reader.Reader = resp.Body
	reader.TotalSize = resp.ContentLength
	reader.FileName = file.Name()
}

func (reader *ThreadPoolWithProgressReader) updateCurrentSize(n int) {
	reader.CurrentSize = reader.CurrentSize + int64(n)
	reader.calculateProgress()
}

func (reader *ThreadPoolWithProgressReader) calculateProgress() {
	reader.CurrentProgress = int(reader.CurrentSize * 100 / reader.TotalSize)
}

func main() {
	urls := make([]string, 4)
	urls[0] = "https://sjbz-fd.zol-img.com.cn/t_s320x510c5/g2/M00/05/0C/ChMlWl1BWGKIa5b1AAkDHph43SoAAMQfgALVicACQM2533.jpg"
	urls[1] = "https://sjbz-fd.zol-img.com.cn/t_s320x510c5/g2/M00/05/0C/ChMlWl1BWGKIa5b1AAkDHph43SoAAMQfgALVicACQM2533.jpg"
	urls[2] = "https://sjbz-fd.zol-img.com.cn/t_s320x510c5/g2/M00/05/0C/ChMlWl1BWGKIa5b1AAkDHph43SoAAMQfgALVicACQM2533.jpg"
	urls[3] = "https://xiazai-fd.zol-img.com.cn/t_s960x600/g1/M01/03/06/Cg-4jVONmIiIa6NpAATdgesQtisAAN9YQLYqJcABN2Z899.jpg"
	tasks := make([]*Task, len(urls))
	for i := 0; i < len(urls); i++ {
		index := i
		filePath := fmt.Sprintf("./文件%d.jpg", index)
		url := urls[index]
		tasks[index] = NewTask(
			func(progress *TaskProgress) {
				util.CustomDownloadFile(filePath, url, NewThreadPoolWithProgressReader(progress))
			})
	}
	fixedThreadPool := NewFixedThreadPool(tasks, 2, true)
	fixedThreadPool.Exec()
}
