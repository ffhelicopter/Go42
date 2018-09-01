# <center>第二十六章 测试</center>

## 26.1 单元测试

首先所有的包都应该有一定的必要文档，然后同样重要的是对包的测试。

名为 testing 的包被专门用来进行自动化测试，日志和错误报告。并且还包含一些基准测试函数的功能。

对一个包做（单元）测试，需要写一些可以频繁（每次更新后）执行的小块测试单元来检查代码的正确性。于是我们必须写一些 Go 源文件来测试代码。测试程序必须属于被测试的包，并且文件名满足这种形式 *_test.Go，所以测试代码和包中的业务代码是分开的。

_test 程序不会被普通的 Go 编译器编译，所以当放应用部署到生产环境时它们不会被部署；只有 Gotest 会编译所有的程序：普通程序和测试程序。

测试文件中必须导入 "testing" 包，并写一些名字以 TestZzz 打头的全局函数，这里的 Zzz 是被测试函数的字母描述，如 TestFmtInterface，TestPayEmployees 等。

测试函数必须有这种形式的头部：

```Go
func TestAbcde(t *testing.T)
```
T 是传给测试函数的结构类型，用来管理测试状态，支持格式化测试日志，如 t.Log，t.Error，t.ErrorF 等。在函数的结尾把输出跟想要的结果对比，如果不等就打印一个错误。成功的测试则直接返回。

用下面这些函数来通知测试失败：

1）func (t *T) Fail()

    标记测试函数为失败，然后继续执行（剩下的测试）。

2）func (t *T) FailNow()

    标记测试函数为失败并中止执行；文件中别的测试也被略过，继续执行下一个文件。

3）func (t *T) Log(args ...interface{})

    args 被用默认的格式格式化并打印到错误日志中。

4）func (t *T) Fatal(args ...interface{})

    结合 先执行 3），然后执行 2）的效果。

运行 Go test 来编译测试程序，并执行程序中所有的 TestZZZ 函数。如果所有的测试都通过会打印出 PASS。

对不能导出的函数不能进行单元或者基准测试。

Gotest 可以接收一个或多个函数程序作为参数，并指定一些选项。

在系统包中，有很多 _test.go 结尾的程序，大家可以用来测试，这里我就不写具体例子了。

## 26.2 基准测试

testing 包中有一些类型和函数可以用来做简单的基准测试；测试代码中必须包含以 BenchmarkZzz 打头的函数并接收一个 *testing.B 类型的参数，比如：

```Go
func BenchmarkReverse(b *testing.B) {
    ...
}
```
命令 Go test –test.bench=.* 会运行所有的基准测试函数；代码中的函数会被调用 N 次（N是非常大的数，如 N = 1000000），并展示 N 的值和函数执行的平均时间，单位为 ns（纳秒，ns/op）。如果是用 testing.Benchmark 调用这些函数，直接运行程序即可。

测试的具体例子

```Go
package even
func Even(i int) bool {     // Exported function
    return i%2 == 0
}
func Odd(i int) bool {      // Exported function
    return i%2 != 0
}
```
在 even 包的路径下，我们创建一个名为 oddeven_test.go 的测试程序：

```Go
package even

import "testing"

func TestEven(t *testing.T) {
	if !Even(10) {
		t.Log(" 10 must be even!")
		t.Fail()
	}
	if Even(7) {
		t.Log(" 7 is not even!")
		t.Fail()
	}
}

func BenchmarkOdd(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Odd(115555555555555)
	}
}
```
现在我们可以用命令：Go test  -test.bench=.*（或 make test）来测试 even 包。

```Go
输出：

goos: windows
goarch: amd64
pkg: ind/even
BenchmarkOdd-4   	2000000000	         0.39 ns/op
PASS
ok  	ind/even	4.785s
```

## 26.3 分析并优化 Go 程序

时间和内存消耗

如果代码使用了 Go 中 testing 包的基准测试功能，我们可以用 Gotest 标准的 -cpuprofile 和 -memprofile 标志向指定文件写入 CPU 或 内存使用情况报告。

使用方式：

```Go
go test -x -v -test.cpuprofile=pprof.out
```
执行上面代码，执行结果 pprof.out 文件中写入 cpu 性能分析信息。

## 26.4 用 pprof 调试

要监控Go程序的堆栈，cpu的耗时等性能信息，我们可以通过使用pprof包来实现。
pprof包有两种方式导入：

```Go
"net/http/pprof"
"runtime/prof"
```
其实net/http/pprof中只是使用runtime/pprof包来进行封装了一下，并在http端口上暴露出来，让我们可以在浏览器查看程序的性能分析。我们可以自行查看net/http/pprof中代码，只有一个文件pprof.go。

下面我们具体说说怎么使用pprof，首先我们讲讲在开发中取得pprof信息的三种方式：

一：web 服务器程序

如果我们的Go程序是用http包启动的web服务器，你想查看自己的web服务器的状态。这个时候就可以选择net/http/pprof。你只需要引入包_"net/http/pprof"，然后就可以在浏览器中使用http://localhost:port/debug/pprof/直接看到当前web服务的状态，包括CPU占用情况和内存使用情况等。

这里port是8080，也就是我们web服务器监听的端口。

```Go
package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"  // 为什么用_ , 在讲解http包时有解释。
)

func myfunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hi")
}

func main() {
	http.HandleFunc("/", myfunc)
	http.ListenAndServe(":8080", nil)
}
```

二：服务进程

如果你的Go程序不是web服务器，而是一个服务进程，可以选择使用net/http/pprof包，然后开启一个goroutine来监听相应端口。

```Go
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"time"
)

func main() {
	// 开启pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()
	go hello()
	select {}
}
func hello() {
	for {
		go func() {
			fmt.Println("hello word")
		}()
		time.Sleep(time.Millisecond * 1)
	}
}
```

在前面这两种方式中，我们在命令行分别运行以下命令：

利用这个命令查看堆栈信息：

go tool pprof http://localhost:8080/debug/pprof/heap

利用这个命令可以查看程序CPU使用情况信息：

go tool pprof http://localhost:8080/debug/pprof/profile

使用这个命令可以查看block信息：

go tool pprof http://localhost:8080/debug/pprof/block

![gotool.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/gotool.png)

这里需要先安装graphviz，http://www.graphviz.org/download/ ，windows平台直接下载zip包，解压缩后把bin目录放到$path中。我们可以通过命令 png 产生图片，还有svg，gif，pdf等命令，生成的图片自动命名存放在当前目录下，我们这里生成了png。其他命令使用可通过help查看。

三：应用程序

如果你的go程序只是一个应用程序，那么你就不能使用net/http/pprof包了，你就需要使用到runtime/pprof。比如下面的例子：

```Go
package main

import (
	"flag"
	"fmt"
	"log"

	"os"
	"runtime/pprof"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	go hello()
	time.Sleep(10 * time.Second)
}
func hello() {
	for {
		go func() {
			fmt.Println("hello word")
		}()
		time.Sleep(time.Millisecond * 1)
	}
}
```
编译后运行：

```Go
study.exe --cpuprofile=cpu.prof
```
这里我们编译后可执行程序是study.exe , 程序运行完后的cpu信息就会记录到cpu.prof中。

现在有了cpu.prof 文件，我们就可以通过go tool pprof 来看相应的信息了。在命令行运行：

```Go
go tool pprof study.exe cpu.prof   
```
这里要注意的是需要带上可执行的程序名以及prof信息文件。

命令执行后会进入到：

![gotooblock.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/gotooblock.png)

界面和前面两种使用net/http/pprof包 一样。我们可以通过go tool pprof 生svg，png或者是pdf文件。

这是生成的png文件，和前面生成的png类似，前面我们生成的是block信息：



通过上面这三种情况的分析，我们可以知道，其实就是两种情况：go tool pprof 后面可以使用http://localhost:8080/debug/pprof/profile 这种url方式，也可以使用study.exe cpu.prof  这种文件方式来进行分析。可以根据你的项目情况灵活使用。

有关pprof，我们就讲这么多，在实际项目有，我们多使用就会发现这个工具还是蛮有用处的。
