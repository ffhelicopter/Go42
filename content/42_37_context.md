# 《Go语言四十二章经》第三十七章 context包

作者：李骁

## 37.1 context包

在Go中，每个请求的request在单独的协程中进行，处理一个request也可能涉及多个协程之间的交互。一个请求衍生出的各个协程之间需要满足一定的约束关系，以实现一些诸如有效期，中止routine树，传递请求全局变量之类的功能。于是Go为我们提供一个解决方案，标准context包。使用context可以使开发者方便的在这些协程之间传递request相关的数据、取消协程的signal或截止时间等。

每个协程在执行之前，都要先知道程序当前的执行状态，通常将这些执行状态封装在一个Context变量中，传递给要执行的协程中。上下文则几乎已经成为传递与请求同生存周期变量的标准方法。在网络编程下，当接收到一个网络请求Request，处理Request时，我们可能需要开启不同的协程来获取数据与逻辑处理，即一个请求Request，会在多个协程中处理。而这些协程可能需要共享Request的一些信息；同时当Request被取消或者超时的时候，所有从这个Request创建的所有协程也应该被结束。

context包不仅实现了在程序单元之间共享状态变量的方法，同时能通过简单的方法，使我们在被调用程序单元的外部，通过设置ctx变量值，将过期或撤销这些信号传递给被调用的程序单元。若存在A调用B的API，B再调用C的API，若A调用B取消，那也要取消B调用C，通过在A, B, C的API调用之间传递Context，以及判断其状态。

Context结构

```Go
// Context包含过期，取消信号，request值传递等，方法在多个协程中协程安全
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

* Deadline会返回一个超时时间，协程获得了超时时间后，例如可以对某些io操作设定超时时间。

* Done方法返回一个通道（channel），当Context被撤销或过期时，该通道关闭，即它是一个表示Context是否已关闭的信号。

* 当Done通道关闭后，Err方法表明Context被撤的原因。

* Value可以让协程共享一些数据，当然获得数据是协程安全的。但使用这些数据的时候要注意同步，比如返回了一个map，而这个map的读写则要加锁。


协程的创建和调用关系总是像层层调用进行的，就像人的辈分一样，而更靠顶部的协程应有办法主动关闭其下属的协程的执行（不然程序可能就失控了）。为了实现这种关系，Context结构也应该像一棵树，叶子节点须总是由根节点衍生出来的。

要创建Context树，第一步就是要得到根节点，context.Background函数的返回值就是根节点：

```Go
func Background() Context
```
该函数返回空的Context，该Context一般由接收请求的第一个协程创建，是与进入请求对应的Context根节点，它不能被取消、没有值、也没有过期时间。它常常作为处理Request的顶层context存在。

有了根节点，又该怎么创建其它的子节点，孙节点呢？context包为我们提供了多个函数来创建他们：

```Go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key interface{}, val interface{}) Context
```

函数都接收一个Context类型的参数parent，并返回一个Context类型的值，这样就层层创建出不同的节点。子节点是从复制父节点得到的，并且根据接收参数设定子节点的一些状态值，接着就可以将子节点传递给下层的协程了。

再回到之前的问题：该怎么通过Context传递改变后的状态呢？使用Context的协程无法取消某个操作，其实这也是符合常理的，因为这些协程是被某个父协程创建的，而理应只有父协程可以取消操作。在父协程中可以通过WithCancel方法获得一个cancel方法，从而获得cancel的权利。

第一个WithCancel函数，它是将父节点复制到子节点，并且还返回一个额外的CancelFunc函数类型变量，该函数类型的定义为：
```Go
type CancelFunc func()
```
调用CancelFunc对象将撤销对应的Context对象，这就是主动撤销Context的方法。在父节点的Context所对应的环境中，通过WithCancel函数不仅可创建子节点的Context，同时也获得了该节点Context的控制权，一旦执行该函数，则该节点Context就结束了，则子节点需要类似如下代码来判断是否已结束，并退出该协程：

```Go
select {
    case <-cxt.Done():
        // do some clean...
}
```

WithDeadline函数的作用也差不多，它返回的Context类型值同样是parent的副本，但其过期时间由deadline和parent的过期时间共同决定。这是因为父节点过期时，其所有的子孙节点必须同时关闭；反之，返回的父节点的过期时间则为deadline。

WithTimeout函数与WithDeadline类似，不过它传入的是从现在开始Context剩余的生命时长。他们都同样也都返回了所创建的子Context的控制权，一个CancelFunc类型的函数变量。

当顶层的Request请求函数结束后，我们就可以cancel掉某个context，从而再在对应协程中根据cxt.Done()来决定是否结束。

WithValue函数，它返回parent的一个副本，调用该副本的Value(key)方法将得到对应key的值。这样我们不光将根节点原有的值保留了，还可以在子孙节点中加入了新的值，注意若存在Key相同，则会被覆盖。

Context对象的生存周期一般仅为一个请求的处理周期。即针对一个请求创建一个Context变量（它为Context树结构的根）；在请求处理结束后，撤销此ctx变量，释放资源。

每次创建一个协程时，可以将原有的Context传递给这个子协程，或者新创建一个子Context传递给这个协程。

Context能灵活地存储不同类型、不同数目的值，并且使多个协程安全地读写其中的值。

当通过父Context对象创建子Context对象时，即可获得子Context的一个撤销函数，这样父Context对象的创建环境就获得了对子Context的撤销权。

注意：使用时遵循context规则

1. 不要把Context存在一个结构体当中，显式地传入函数。Context变量需要作为第一个参数使用，一般命名为ctx。

2. 即使方法允许，也不要传入一个nil的Context，如果你不确定你要用什么Context的时候传一个context.TODO。

3. 使用context的Value相关方法只应该用于在程序和接口中传递的和请求相关的元数据，不要用它来传递一些可选的参数。

4. 同样的Context可以用来传递到不同的协程中，Context在多个协程中是安全的。


在子Context被传递到的协程中，应该对该子Context的Done通道（channel）进行监控，一旦该通道被关闭（即上层运行环境撤销了本协程的执行），应主动终止对当前请求信息的处理，释放资源并返回。

## 37.2 context应用

前面介绍协程时，对协程的管理和控制我们并没有进行讨论。到目前我们已经清楚认识了channel、context以及sync包，通过这三者，我们完全可以达到完美控制协程运行的目的。

通过go关键字让我们很容易启动一个协程，但难的是很好的管理和控制他们的运行。有几种方法我们可以根据场景使用：

（1）使用sync.WaitGroup，它用于线程总同步，会等待一组线程集合完成，才会继续向下执行，这对监控所有子协程全部完成情况特别有用，但要控制某个协程就无能为力了；

（2）使用channel来传递消息，一个协程来发送channel信号，另一个协程通过select来得到channel信息，这种方式可以满足协程之间的通信，来控制协程运行。但如果协程数量达到一定程度，就很难把控了；或者这两个协程还和其他协程也有类似通信，比如A与B，B与C，如果A发信号B退出了，C有可能等不到B的channel信号而被遗忘；

（3）使用Context来传递消息，Context是层层传递机制，根节点完全控制了子节点，根节点（父节点）可以根据需要选择自动还是手动结束子节点。而每层节点所在的协程就可以根据信息来决定下一步的操作。

下面我们来看看具体使用Context怎么来控制协程的运行：

这里用Context同时控制2个协程，这2个协程都可以收到cancel()发出的信号，甚至doNothing这样不结束协程可反复接收cancel信息。

```Go
package main

import (
	"context"
	"log"
	"os"
	"time"
)

var logs *log.Logger

func doClearn(ctx context.Context) {
	// for 循环来每1秒work一下，判断ctx是否被取消了，如果是就退出
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			logs.Println("doClearn:收到Cancel，做好收尾工作后马上退出。")
			return
		default:
			logs.Println("doClearn:每隔1秒观察信号，继续观察...")
		}
	}
}

func doNothing(ctx context.Context) {
	for {
		time.Sleep(3 * time.Second)
		select {
		case <-ctx.Done():
			logs.Println("doNothing:收到Cancel，但不退出......")

			// 注释return可以观察到，ctx.Done()信号是可以一直接收到的，return不注释意味退出协程
			//return
		default:
			logs.Println("doNothing:每隔3秒观察信号，一直运行")
		}
	}
}

func main() {
	logs = log.New(os.Stdout, "", log.Ltime)

	// 新建一个ctx
	ctx, cancel := context.WithCancel(context.Background())

	// 传递ctx
	go doClearn(ctx)
	go doNothing(ctx)

	// 主程序阻塞20秒，留给协程来演示
	time.Sleep(20 * time.Second)
	logs.Println("cancel")

	// 调用cancel：context.WithCancel 返回的CancelFunc
	cancel()

	// 发出cancel 命令后，主程序阻塞10秒，再看协程的运行情况
	time.Sleep(10 * time.Second)
}

程序输出：
......
cancel
doClearn:收到Cancel，做好收尾工作后马上退出。
doNothing:收到Cancel，但不退出......
doNothing:收到Cancel，但不退出......
doNothing:收到Cancel，但不退出......
```

这里用Context嵌套控制3个协程，A，B，C。在主程序发出cancel信号后，每个协程都能接收根Context的Done()信号而退出。

```Go
package main

import (
	"context"
	"fmt"
	"time"
)

func A(ctx context.Context) int {
	ctx = context.WithValue(ctx, "AFunction", "Great")

	go B(ctx)

	select {
	// 监测自己上层的ctx ...
	case <-ctx.Done():
		fmt.Println("A Done")
		return -1
	}
	return 1
}

func B(ctx context.Context) int {
	fmt.Println("A value in B:", ctx.Value("AFunction"))
	ctx = context.WithValue(ctx, "BFunction", 999)

	go C(ctx)

	select {
	// 监测自己上层的ctx ...
	case <-ctx.Done():
		fmt.Println("B Done")
		return -2
	}
	return 2
}

func C(ctx context.Context) int {
	fmt.Println("B value in C:", ctx.Value("AFunction"))
	fmt.Println("B value in C:", ctx.Value("BFunction"))
	select {
	// 结束时候做点什么 ...
	case <-ctx.Done():
		fmt.Println("C Done")
		return -3
	}
	return 3
}

func main() {
	// 自动取消(定时取消)
	{
		timeout := 10 * time.Second
		ctx, _ := context.WithTimeout(context.Background(), timeout)

		fmt.Println("A 执行完成，返回：", A(ctx))
		select {
		case <-ctx.Done():
			fmt.Println("context Done")
			break
		}
	}
	time.Sleep(20 * time.Second)
}
```

最后我们看看Context在http 是怎么传递的：

```Go
package main

import (
	"context"
	"net/http"
	"time"
)

// ContextMiddle是http服务中间件，统一读取通行cookie并使用ctx传递
func ContextMiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("Check")
		if cookie != nil {
			ctx := context.WithValue(r.Context(), "Check", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// 强制设置通行cookie
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	expitation := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{Name: "Check", Value: "42", Expires: expitation}
	http.SetCookie(w, &cookie)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// 通过取中间件传过来的context值来判断是否放行通过
	if chk := r.Context().Value("Check"); chk == "42" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Let's go! \n"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No Pass!"))
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)

	// 人为设置通行cookie
	mux.HandleFunc("/chk", CheckHandler)

	ctxMux := ContextMiddle(mux)
	http.ListenAndServe(":8080", ctxMux)
}
```

我们打开浏览器访问：http://localhost:8080/chk ，然后再访问：http://localhost:8080/ ，将会看到我们正常通行后结果，否则将会看到没有正常通行下的信息。Context信息的传递主要靠中间件ContextMiddle来进行。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第三十六章 net/http包](https://github.com/ffhelicopter/Go42/blob/master/content/42_36_http.md)

[第三十八章 数据序列化](https://github.com/ffhelicopter/Go42/blob/master/content/42_38_json.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
