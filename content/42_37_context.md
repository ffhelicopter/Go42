# 《Go语言四十二章经》第三十七章 context包

作者：李骁

## 37.1 context包

在Go中，每个请求的request在单独的goroutine中进行，处理一个request也可能涉及多个goroutine之间的交互。一个请求衍生出的各个 goroutine 之间需要满足一定的约束关系，以实现一些诸如有效期，中止routine树，传递请求全局变量之类的功能。于是Go为我们提供一个解决方案，标准context包。使用context可以使开发者方便的在这些goroutine之间传递request相关的数据、取消goroutine的signal或截止时间等。

每个goroutine在执行之前，都要先知道程序当前的执行状态，通常将这些执行状态封装在一个Context变量中，传递给要执行的goroutine中。上下文则几乎已经成为传递与请求同生存周期变量的标准方法。在网络编程下，当接收到一个网络请求Request，处理Request时，我们可能需要开启不同的goroutine来获取数据与逻辑处理，即一个请求Request，会在多个goroutine中处理。而这些goroutine可能需要共享Request的一些信息；同时当Request被取消或者超时的时候，所有从这个Request创建的所有goroutine也应该被结束。

context包不仅实现了在程序单元之间共享状态变量的方法，同时能通过简单的方法，使我们在被调用程序单元的外部，通过设置ctx变量值，将过期或撤销这些信号传递给被调用的程序单元。若存在A调用B的API，B再调用C的API，若A调用B取消，那也要取消B调用C，通过在A, B, C的API调用之间传递Context，以及判断其状态。

Context结构

```Go
// Context包含过期，取消信号，request值传递等，方法在多个goroutine中协程安全
type Context interface {
    // Done 方法在context被取消或者超时返回一个close的channel
    Done() <-chan struct{}

    Err() error

    // Deadline 返回context超时时间
    Deadline() (deadline time.Time, ok bool)

    // Value 返回context相关key对应的值
    Value(key interface{}) interface{}
}
```

* Deadline会返回一个超时时间，goroutine获得了超时时间后，例如可以对某些io操作设定超时时间。

* Done方法返回一个通道（channel），当Context被撤销或过期时，该通道关闭，即它是一个表示Context是否已关闭的信号。

* 当Done通道关闭后，Err方法表明Context被撤的原因。

* Value可以让goroutine共享一些数据，当然获得数据是协程安全的。但使用这些数据的时候要注意同步，比如返回了一个map，而这个map的读写则要加锁。


goroutine的创建和调用关系总是像层层调用进行的，就像人的辈分一样，而更靠顶部的goroutine应有办法主动关闭其下属的goroutine的执行（不然程序可能就失控了）。为了实现这种关系，Context结构也应该像一棵树，叶子节点须总是由根节点衍生出来的。

要创建Context树，第一步就是要得到根节点，context.Background函数的返回值就是根节点：

```Go
func Background() Context
```
该函数返回空的Context，该Context一般由接收请求的第一个goroutine创建，是与进入请求对应的Context根节点，它不能被取消、没有值、也没有过期时间。它常常作为处理Request的顶层context存在。

有了根节点，又该怎么创建其它的子节点，孙节点呢？context包为我们提供了多个函数来创建他们：

```Go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key interface{}, val interface{}) Context
```

函数都接收一个Context类型的参数parent，并返回一个Context类型的值，这样就层层创建出不同的节点。子节点是从复制父节点得到的，并且根据接收参数设定子节点的一些状态值，接着就可以将子节点传递给下层的goroutine了。

再回到之前的问题：该怎么通过Context传递改变后的状态呢？使用Context的goroutine无法取消某个操作，其实这也是符合常理的，因为这些goroutine是被某个父goroutine创建的，而理应只有父goroutine可以取消操作。在父goroutine中可以通过WithCancel方法获得一个cancel方法，从而获得cancel的权利。

第一个WithCancel函数，它是将父节点复制到子节点，并且还返回一个额外的CancelFunc函数类型变量，该函数类型的定义为：
```Go
type CancelFunc func()
```
调用CancelFunc对象将撤销对应的Context对象，这就是主动撤销Context的方法。在父节点的Context所对应的环境中，通过WithCancel函数不仅可创建子节点的Context，同时也获得了该节点Context的控制权，一旦执行该函数，则该节点Context就结束了，则子节点需要类似如下代码来判断是否已结束，并退出该goroutine：

```Go
select {
    case <-cxt.Done():
        // do some clean...
}
```

WithDeadline函数的作用也差不多，它返回的Context类型值同样是parent的副本，但其过期时间由deadline和parent的过期时间共同决定。当parent的过期时间早于传入的deadline时间时，返回的过期时间应与parent相同。父节点过期时，其所有的子孙节点必须同时关闭；反之，返回的父节点的过期时间则为deadline。

WithTimeout函数与WithDeadline类似，只不过它传入的是从现在开始Context剩余的生命时长。他们都同样也都返回了所创建的子Context的控制权，一个CancelFunc类型的函数变量。

当顶层的Request请求函数结束后，我们就可以cancel掉某个context，从而层层goroutine根据判断cxt.Done()来结束。

WithValue函数，它返回parent的一个副本，调用该副本的Value(key)方法将得到val。这样我们不光将根节点原有的值保留了，还在子孙节点中加入了新的值，注意若存在Key相同，则会被覆盖。

context包通过构建树型关系的Context，来达到上一层goroutine能对传递给下一层goroutine的控制。对于处理一个Request请求操作，需要采用context来层层控制goroutine，以及传递一些变量来共享。

Context对象的生存周期一般仅为一个请求的处理周期。即针对一个请求创建一个Context变量（它为Context树结构的根）；在请求处理结束后，撤销此ctx变量，释放资源。

每次创建一个goroutine，要么将原有的Context传递给goroutine，要么创建一个子Context并传递给goroutine。

Context能灵活地存储不同类型、不同数目的值，并且使多个goroutine安全地读写其中的值。

当通过父Context对象创建子Context对象时，可同时获得子Context的一个撤销函数，这样父Context对象的创建环境就获得了对子Context将要被传递到的goroutine的撤销权。

注意：使用时遵循context规则

1. 不要把Context存在一个结构体当中，显式地传入函数。Context变量需要作为第一个参数使用，一般命名为ctx。

2. 即使方法允许，也不要传入一个nil的Context，如果你不确定你要用什么Context的时候传一个context.TODO。

3. 使用context的Value相关方法只应该用于在程序和接口中传递的和请求相关的元数据，不要用它来传递一些可选的参数。

4. 同样的Context可以用来传递到不同的goroutine中，Context在多个goroutine中是安全的。


在子Context被传递到的goroutine中，应该对该子Context的Done通道（channel）进行监控，一旦该通道被关闭（即上层运行环境撤销了本goroutine的执行），应主动终止对当前请求信息的处理，释放资源并返回。

## 37.2 context应用

```Go
package main

import (
	"context"
	"log"
	"os"
	"time"
)

var logg *log.Logger

func someHandler() {
	// 新建一个ctx
	ctx, cancel := context.WithCancel(context.Background())

	//传递ctx
	go doStuff(ctx)

	//10秒后取消doStuff
	time.Sleep(10 * time.Second)
	log.Println("cancel")

	//调用cancel：context.WithCancel 返回的CancelFunc
	cancel()

}

func doStuff(ctx context.Context) {

	// for 循环来每1秒work一下，判断ctx是否被取消了，如果是就退出

	for {
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			logg.Printf("done")
			return
		default:
			logg.Printf("work")
		}
	}
}

func main() {
	logg = log.New(os.Stdout, "", log.Ltime)
	someHandler()
	logg.Printf("down")
}

```

```Go
程序输出：
16:28:21 work
16:28:22 work
16:28:23 work
16:28:24 work
16:28:25 work
16:28:26 work
16:28:27 work
16:28:28 work
16:28:29 work
2018/08/22 16:28:30 cancel
16:28:30 down
```

someHandler() 作为顶层的Request请求函数，处理完主要任务后，主动cancel掉context，而子层goroutine  doStuff(ctx context.Context) 根据判断cxt.Done()来结束。


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>本书《Go语言四十二章经》内容在简书同步地址：  https://www.jianshu.com/nb/29056963
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
