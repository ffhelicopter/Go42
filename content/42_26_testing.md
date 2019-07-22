# 《Go语言四十二章经》第二十六章 测试

作者：李骁


在Go语言中，所有的包都应该有必要文档和注释，当然同样甚至更为重要的是对包进行必要的测试。

testing 包就是这样一个标准包，被专门用来进行单元测试以及进行自动化测试，打印日志和错误报告，方便程序员调试代码，并且还包含一些基准测试函数的功能。

testing 包含测试函数、测试辅助代码和示例函数；测试函数包括Test开头的单元测试函数和以Benchmark开头的基准测试函数两种，测试辅助代码是为测试函数服务的公共函数、初始化函数、测试数据等。

而示例函数则是名称以Example开头函数，通常保存在example_*_test.go文件中，示例函数检测实际输出与注释中的期望输出是否一致，一致则测试通过，不一致则测试失败。

## 26.1 单元测试

开发中经常需要对一个包做（单元）测试，写一些可以频繁（每次更新后）执行的小块测试单元来检查代码的正确性，于是我们必须写一些 Go 源文件来测试代码。

使用testing包，我们只需要遵守简单的规则，就可以很好地写出通用的测试程序。因为其他开发人员也会遵循这个包的规则来进行测试。

首先测试程序是独立的文件，他必须属于被测试的包，和这个包的其他程序放在一起，并且文件名满足这种形式 *_test.go。由于是独立的测试文件，所以测试代码和包中的业务代码是分开的。Go语言这样规定的好处是不言而喻的，因为在其他语言开发的程序中，我们经常可以看到代码中注释掉的测试代码，而且有把开发版作为生产版发布到线上导致异常的问题出现。

当然，好的规则需要我们遵守并严格执行。

_test 程序不会被普通的 Go 编译器编译，所以当放应用部署到生产环境时它们不会被部署；只有 Gotest 会编译所有的程序：普通程序和测试程序。

测试文件中必须导入 "testing" 包，测试函数名字是以 TestXxx 打头的全局函数，Xxx部分可以为任意的字母数字的组合，但是首字母不能是小写字母[a-z]，函数名我们可以以被测试函数的字母描述，如 TestFmtInterface，TestPayEmployees 等。测试用例会按照测试源代码中写的顺序依次执行。

测试函数一般都要求这种形式的头部：

```Go
func TestAbcde(t *testing.T)
```

*testing.T是传给测试函数的结构类型，用来管理测试状态，支持格式化测试日志，如 t.Log，t.Error，t.ErrorF 等。t.Log函数就像我们常常使用的fmt.Println一样，可以接受多个参数，方便输出调试结果。

用下面这些函数来通知测试失败：
1）func (t *T) Fail()
    标记测试函数为失败，然后继续执行剩下的测试。

2）func (t *T) FailNow()
    标记测试函数为失败并中止执行；文件中别的测试也被略过，继续执行下一个文件。

3）func (t *T) Log(args ...interface{})
    args 被用默认的格式格式化并打印到错误日志中。

4）func (t *T) Fatal(args ...interface{})
    结合 先执行 3），然后执行 2）的效果。

运行 go test 来编译测试程序，并执行程序中所有的 TestXxx 函数。如果所有的测试都通过会打印出 PASS。

当然，对于包中不能导出的函数不能进行单元或者基准测试。

gotest 可以接收一个或多个函数程序作为参数，并指定一些选项。

go test 常用参数
-cpu: 指定测试的GOMAXPROCS值，默认是GOMAXPROCS当前值
-count: 运行单元测试和基准测试n次（默认1）。如设置了-cpu，则为每个GOMAXPROCS运行n次，示例函数总运行一次。
-cover: 启用覆盖率分析
-run: 执行功能测试函数，支持正则匹配，可以选择测试函数或者测试文件来仅测试单个函数或者单个文件
-bench: 执行基准测试函数，支持正则匹配
-benchtime: 基准测试最大时间上限
-parallel: 允许并行执行的最大测试数，默认情况下设置为GOMAXPROCS的值
-v: 展示测试过程信息

在系统标准包中，有很多 _test.go 结尾的程序，大家可以用来测试，为节约篇幅这里我就不写具体例子了。

## 26.2 基准测试

testing 包中有一些类型和函数可以用来做简单的基准测试；测试代码中必须包含以 BenchmarkZzz 打头的函数并接收一个 *testing.B 类型的参数，比如：

```Go
func BenchmarkReverse(b *testing.B) {
    ...
}
```

命令 go test –test.bench=.* 会运行所有的基准测试函数；代码中的函数会被调用 N 次（N是非常大的数，如 N = 1000000），可以根据情况指定b.N的值，并展示 N 的值和函数执行的平均时间，单位为 ns（纳秒，ns/op）。如果是用 testing.Benchmark 调用这些函数，直接运行程序即可。

下面我们看一个测试的具体例子：

```Go
package even

func Loop(n uint64) (result uint64) {
	result = 1
	var i uint64 = 1
	for ; i <= n; i++ {
		result *= i
	}
	return result
}

func Factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * Factorial(n-1)
		return result
	}
	return 1
}
```

在 even 包的路径下，我们创建一个名为 even_test.go 的测试程序：

```Go
package even

import (
	"testing"
)

func TestLoop(t *testing.T) {
	t.Log("Loop:", Loop(uint64(32)))
}

func TestFactorial(t *testing.T) {
	t.Log("Factorial:", Factorial(uint64(32)))
}

func BenchmarkLoop(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Loop(uint64(40))
	}
}

func BenchmarkFactorial(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Factorial(uint64(40))
	}
}
```

现在我们可以在这个包的目录下使用命令：go test  -test.bench=.* 来测试 even 包。

输出：

```Go
输出：

goos: windows
goarch: amd64
pkg: go42/chapter-13/13.1/1
BenchmarkLoop-4        	50000000	        27.2 ns/op
BenchmarkFactorial-4   	10000000	       163 ns/op
PASS
ok  	go42/chapter-13/13.1/1	3.628s
```

递归函数的确是很耗费系统资源，而且运行也慢，不建议使用。

## 26.3 分析并优化 Go 程序

如果代码使用了 Go 中 testing 包的基准测试功能，我们可以用 gotest 标准的 -cpuprofile 和 -memprofile 标志向指定文件写入 CPU 或 内存使用情况报告。

使用方式：

```Go
go test -x -v -test.cpuprofile=pprof.out
```

运行上面代码，将会基于基准测试把执行结果中的 cpu 性能分析信息写到 pprof.out 文件中。我们可以根据这个文件做分析来详细了解性能情况。

## 26.4 用 pprof 调试

要监控Go程序的堆栈，cpu的耗时等性能信息，我们可以通过使用pprof包来实现。在代码中，pprof包有两种方式导入：

```Go
"net/http/pprof"
"runtime/prof"
```

其实net/http/pprof中只是使用runtime/pprof包来进行封装了一下，并在http端口上暴露出来，让我们可以在浏览器查看程序的性能分析。我们可以自行查看net/http/pprof中代码，只有一个文件pprof.go。

下面我们具体说说怎么使用pprof，首先我们讲讲在开发中取得pprof信息的三种方式：

一：web 服务器程序

如果我们的Go程序是web服务器，你想查看自己的web服务器的状态。这个时候就可以选择net/http/pprof。你只需要引入包_"net/http/pprof"，然后就可以在浏览器中使用http://localhost:port/debug/pprof/直接看到当前web服务的状态，包括CPU占用情况和内存使用情况等。

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

访问http://localhost:8080/debug/pprof/

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

访问http://localhost:8080/debug/pprof/
在前面这两种方式中，我们还可以在命令行分别运行以下命令：

利用这个命令查看堆栈信息：
go tool pprof http://localhost:8080/debug/pprof/heap
利用这个命令可以查看程序CPU使用情况信息：
go tool pprof http://localhost:8080/debug/pprof/profile
使用这个命令可以查看block信息：
go tool pprof http://localhost:8080/debug/pprof/block

![gotool.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/gotool.png)

这里需要先安装graphviz，http://www.graphviz.org/download/ ，windows平台直接下载zip包，解压缩后把bin目录放到$path中。我们可以通过执行命令 png 产生图片，还有svg，gif，pdf等命令，生成的图片自动命名存放在当前目录下，我们这里生成了png。其他命令使用可通过help查看。

三：应用程序

如果你的Go程序只是一个应用程序，那么你就不能使用net/http/pprof包了，你就需要使用到runtime/pprof。比如下面的例子：

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

func Factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * Factorial(n-1)
		return result
	}
	return 1
}

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

	go compute()
	time.Sleep(10 * time.Second)
}
func compute() {
	for i := 0; i < 100; i++ {
		go func() {
			fmt.Println(Factorial(uint64(40)))
		}()
		time.Sleep(time.Millisecond * 1)
	}
}
```

编译后生成3.exe文件并运行：

```Go
3.exe --cpuprofile=cpu.prof
```

这里我们编译后可执行程序是3.exe , 程序运行完后的cpu信息就会记录到cpu.prof中。

现在有了cpu.prof 文件，我们就可以通过go tool pprof 来看相应的信息了。在命令行运行：

```Go
go tool pprof 3.exe cpu.prof 
```

这里要注意的是需要带上可执行的程序名以及prof信息文件。

命令执行后会进入到：

![132.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/132.png)

命令界面和前面两种使用net/http/pprof包 一样。我们可以通过go tool pprof 生svg，png或者是pdf文件。

这是生成的png文件，和前面生成的png类似，前面我们生成的是block信息：

![profile001.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/profile001.png)

通过上面这三种情况的分析，我们可以知道，其实就是两种情况：
go tool pprof http://localhost:8080/debug/pprof/profile 这种url方式，或者
go tool pprof 3.exe cpu.prof   这种文件方式来进行分析。

我们可以根据项目情况灵活使用。有关pprof，我们就讲这么多，在实际项目中，我们多使用就会发现这个工具还是蛮有用处的。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第二十五章 面向对象](https://github.com/ffhelicopter/Go42/blob/master/content/42_25_oo.md)

[第二十七章 反射(reflect)](https://github.com/ffhelicopter/Go42/blob/master/content/42_27_reflect.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
