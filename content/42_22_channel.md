# 《Go语言四十二章经》第二十二章 通道(channel)

作者：李骁

## 22.1 通道(channel)

Go 奉行通过通信来共享内存，而不是共享内存来通信。所以，**channel 是协程之间互相通信的通道**，协程之间可以通过它发送消息和接收消息。

通道是进程内的通信方式，因此通过通道传递对象的行为与函数调用时参数传递行为比较一致，比如也可以传递指针等。

通道消息传递与消息类型也有关系，一个通道只能传递（发送send或接收receive）类型的值，这需要在声明通道时指定。

默认情况下，通道是阻塞的 (叫做无缓冲的通道)。

使用make来建立一个通道：

```Go
var channel chan int = make(chan int)
// 或
channel := make(chan int)
```
Go中通道可以是发送（send）、接收（receive）、同时发送（send）和接收（receive）。

```Go
// 定义接收的通道
receive_only := make (<-chan int)
 
// 定义发送的通道
send_only := make (chan<- int)

// 可同时发送接收
send_receive := make (chan int)
```

* chan<- 表示数据进入通道，要把数据写进通道，对于调用者就是发送。
* <-chan 表示数据从通道出来，对于调用者就是得到通道的数据，当然就是接收。

定义只发送或只接收的通道意义不大，一般用于在参数传递中：

```Go
package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int) // 不使用带缓冲区的通道
	go send(c)
	go recv(c)
	time.Sleep(3 * time.Second)
close(c)
}

// 只能向chan里send数据
func send(c chan<- int) {
	for i := 0; i < 10; i++ {

		fmt.Println("send readey ", i)
		c <- i
		fmt.Println("send ", i)
	}
}

// 只能接收通道中的数据
func recv(c <-chan int) {
	for i := range c {
		fmt.Println("received ", i)
	}
}
```
```Go
程序输出：

send readey  0
send  0
send readey  1
received  0
received  1
send  1
send readey  2
send  2
send readey  3
received  2
received  3
send  3
send readey  4
send  4
send readey  5
received  4
received  5
send  5
send readey  6
send  6
send readey  7
received  6
received  7
send  7
send readey  8
send  8
send readey  9
received  8
received  9
send  9
```
运行结果上我们可以发现一个现象，往通道发送数据后，这个数据如果没有被取走，通道是阻塞的，也就是不能继续向通道里面发送数据。上面代码中，我们没有指定通道缓冲区的大小，默认情况下是阻塞的。

我们可以建立带缓冲区的通道：

```Go
c := make(chan int, 1024)
```
我们把前面的程序修改下：

```Go
package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 10) // 使用带缓冲区的通道
	go send(c)
	go recv(c)
	time.Sleep(3 * time.Second)
	close(c)
}

// 只能向chan里send发送数据
func send(c chan<- int) {
	for i := 0; i < 10; i++ {

		fmt.Println("send readey ", i)
		c <- i
		fmt.Println("send ", i)
	}
}

// 只能接收通道中的数据
func recv(c <-chan int) {
	for i := range c {
		fmt.Println("received ", i)
	}
}
```

```Go
程序输出：

send readey  0
send  0
send readey  1
send  1
send readey  2
send  2
send readey  3
send  3
send readey  4
send  4
send readey  5
received  0
received  1
received  2
received  3
received  4
received  5
send  5
send readey  6
send  6
send readey  7
send  7
send readey  8
send  8
send readey  9
send  9
received  6
received  7
received  8
received  9
```

从运行结果我们可以看到（每次执行顺序不一定相同，协程运行导致的原因），带有缓冲区的通道，在缓冲区有数据而未填满前，读取不会出现阻塞的情况。


* 无缓冲的通道（unbuffered channel）是指在接收前没有能力保存任何值的通道。

这种类型的通道要求发送协程和接收协程同时准备好，才能完成发送和接收操作。如果两个协程没有同时准备好，通道会导致先执行发送或接收操作的协程阻塞等待。

这种对通道进行发送和接收的交互行为本身就是同步的。

* 有缓冲的通道（buffered channel）是一种在被接收前能存储一个或者多个值的通道。

这种类型的通道并不强制要求协程之间必须同时完成发送和接收。通道会阻塞发送和接收动作的条件也会不同。只有在通道中没有要接收的值时，接收动作才会阻塞。只有在通道没有可用缓冲区容纳被发送的值时，发送动作才会阻塞。

这导致有缓冲的通道和无缓冲的通道之间的一个很大的不同：无缓冲的通道保证进行发送和接收的协程会在同一时间进行数据交换；有缓冲的通道没有这种保证。

如果给定了一个缓冲区容量，通道就是异步的。只要缓冲区有未使用空间用于发送数据，或还包含可以接收的数据，那么其通信就会无阻塞地进行。

可以通过内置的close函数来关闭通道实现。

* 通道不需要经常去关闭，只有当没有任何可发送数据时才去关闭通道；

* 关闭通道后，无法向通道再发送数据(引发panic 错误后导致接收立即返回零值)；

* 关闭通道后，可以继续向通道接收数据，不能继续发送数据；

* 对于nil 通道，无论收发都会被阻塞。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第二十一章 协程(goroutine)](https://github.com/ffhelicopter/Go42/blob/master/content/42_21_goroutine.md)

[第二十三章 同步与锁](https://github.com/ffhelicopter/Go42/blob/master/content/42_23_sync.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com