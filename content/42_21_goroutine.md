# 《Go语言四十二章经》第二十一章 协程(goroutine)

作者：李骁

>Concurrency is about dealing with lots of things at once. <br>
>Parallelism is about doing lots of things at once.<br>
>
>并发： 指的是程序的逻辑结构。如果程序代码结构中的某些函数逻辑上可以同时运行，但物理上未必会同时运行。<br>
>并行： 并行是指程序的运行状态。并行则指的就是在物理层面也就是使用了不同CPU在执行不同或者相同的任务。

## 21.1 并发

并发是在同一时间处理多件事情。并行是在同一时间做多件事情。并发的目的在于把当个 CPU 的利用率使用到最高。并行则需要多核 CPU 的支持。

Go 语言在语言层面上支持了并发，goroutine是Go语言提供的一种用户态线程，有时我们也称之为协程。所谓的协程，某种程度上也可以叫做轻量线程，它不由系统而由应用程序创建和管理，因此使用开销较低（一般为4K）。我们可以创建很多的协程，并且它们跑在同一个内核线程之上的时候，就需要一个调度器来维护这些协程，确保所有的协程都能使用CPU，并且是尽可能公平地使用CPU资源。

调度器的主要有4个重要部分，分别是M、G、P、Sched。

* M (work thread)  代表了系统线程内核线程，由操作系统管理。

* P (processor)    衔接M和G的调度上下文，它负责将等待执行的G与M对接。P的数量可以通过GOMAXPROCS()来设置，它其实也就代表了真正的并发度，即有多少个goroutine可以同时运行。

* G (goroutine)    协程的实体，包括了调用栈，重要的调度信息，例如channel等。

在操作系统的内核线程和编程语言的用户线程之间，实际上存在3种线程对应模型，也就是：1:1，1:N，M:N。

N:1 多个（N）用户线程始终在一个内核线程上跑，context上下文切换很快，但是无法真正的利用多核。 
1:1 一个用户线程就只在一个内核线程上跑，这时可以利用多核，但是上下文切换很慢，切换效率很低。 
M:N 多个协程在多个内核线程上跑，这个可以集齐上面两者的优势，但是无疑增加了调度的难度。

M:N 综合两种方式（N:1，1:1）的优势。多个协程可以在多个内核线程上处理。既能快速切换上下文，也能利用多核的优势，而Go正是选择这种实现方式。

Go 语言中的协程是运行在多核CPU中的(通过runtime.GOMAXPROCS(1)设定CPU核数)。 实际中运行的CPU核数未必会和实际物理CPU数相吻合。

每个协程都会被一个特定的P(某个CPU)选定维护，而M(物理计算资源)每次挑选一个有效P，然后执行P中的协程。

每个P会将自己所维护的协程放到一个G队列中，其中就包括了协程堆栈信息，是否可执行信息等等。

默认情况下，P的数量与实际物理CPU的数量相等。当我们通过循环来创建协程时，协程会被分配到不同的G队列中。 而M的数量又不是唯一的，当M随机挑选P时，也就等同随机挑选了协程。

所以，当我们碰到多个协程的执行顺序不是我们想象的顺序时就可以理解了，因为协程进入P管理的队列G是带有随机性的。

P的数量由runtime.GOMAXPROCS(1)所设定，通常来说它是和内核数对应，例如在4Core的服务器上会启动4个线程。G会有很多个，每个P会将协程从一个就绪的队列中做Pop操作，为了减小锁的竞争，通常情况下每个P会负责一个队列。

```Go
runtime.NumCPU()        // 返回当前CPU内核数
runtime.GOMAXPROCS(2)  // 设置运行时最大可执行CPU数
runtime.NumGoroutine() // 当前正在运行的协程 数
```

P维护着这个队列（称之为runqueue），Go语言里，启动一个协程很容易：go function 就行，所以每有一个go语句被执行，runqueue队列就在其末尾加入一个协程，在下一个调度点，就从runqueue中取出一个协程执行。

假如有两个M，即两个内核线程，分别对应一个P，每一个P调度一个G队列。如此一来，就组成的协程运行时的基本结构：

* 当有一个M返回时，它必须尝试取得一个P来运行协程，一般情况下，它会从其他的OS Thread线程那里窃取一个P过来，如果没有拿到，它就把协程放在一个global runqueue里，然后自己进入线程缓存里。

* 如果某个P所分配的任务G很快就执行完了，这会导致多个队列存在不平衡，会从其他队列中截取一部分协程到P上进行调度。一般来说，如果P从其他的P那里要取任务的话，一般就取run queue的一半，这就确保了每个内核线程都能充分的使用。

* 当一个内核线程被阻塞时，P可以转而投奔另一个内核线程。


我们可以运行下面代码体验下Go语言中通过设定runtime.GOMAXPROCS(2) ，也即手动指定CPU运行的核数，来体验多核CPU在并发处理时的威力。不得不提，递归函数的计算很费CPU和内存，运行时可以根据电脑配置修改循环或递归数量。

```Go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var quit chan int = make(chan int)

func loop() {
	for i := 0; i < 1000; i++ {
		Factorial(uint64(1000))
	}
	quit <- 1
}
func Factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * Factorial(n-1)
		return result
	}
	return 1
}

var wg1, wg2 sync.WaitGroup

func main() {
	fmt.Println("1:", time.Now())
	fmt.Println(runtime.NumCPU()) // 默认CPU核数
	a := 5000
	for i := 1; i <= a; i++ {
		wg1.Add(1)
		go loop()
	}

	for i := 0; i < a; i++ {
		select {
		case <-quit:
			wg1.Done()
		}
	}
	fmt.Println("2:", time.Now())
	wg1.Wait()

	fmt.Println("3:", time.Now())
	runtime.GOMAXPROCS(2) // 设置执行使用的核数
	a = 5000
	for i := 1; i <= a; i++ {
		wg2.Add(1)
		go loop()
	}

	for i := 0; i < a; i++ {
		select {
		case <-quit:
			wg2.Done()
		}
	}

	fmt.Println("4:", time.Now())
	wg2.Wait()
	fmt.Println("5:", time.Now())
}
```

我的测试电脑CPU默认是4核，对比手动设置CPU在2核时的运行耗时，4核耗时约8秒，2核约14秒，当然这是一种比较理想化的测试，因为阶乘很快导致unit64为0，所以这个测试并不严谨，但从中我们仍然可以体验到Go语言在处理并发时代码之简单，控制之方便。

在实际中运行速度延缓可能不一定仅仅是由于CPU的竞争，可能还有内存或者I/O的原因导致的，我们需要根据情况仔细分析。

最后，runtime.Gosched()用于让出CPU时间片，让出当前协程的执行权限，调度器安排其他等待的任务运行，并在下次某个时候从该位置恢复执行。

## 21.2 goroutine

在Go语言中，协程的使用很简单，直接在函数（代码块）前加上关键字 go 即可。go关键字就是用来创建一个协程的，后面的代码块就是这个协程需要执行的代码逻辑。

```Go
package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 1; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
	// 暂停一会，保证打印全部结束
	time.Sleep(1e9)
}
```

time.Sleep(1e9)让主程序不会马上退出，以便让协程运行完成，避免主程序退出时协程未处理完成甚至没有开始运行。

有关于协程之间的通信以及协程与主线程的控制以及多个协程的管理和控制，我们后续通过channel、context以及锁来进一步说明。




[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第二十章 方法](https://github.com/ffhelicopter/Go42/blob/master/content/42_20_method.md)

[第二十二章 通道(channel)](https://github.com/ffhelicopter/Go42/blob/master/content/42_22_channel.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
