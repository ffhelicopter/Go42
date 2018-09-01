# <center>第十五章 错误处理</center>

## 15.1 错误类型
任何时候当你需要一个新的错误类型，都可以用 errors（必须先 import）包的 errors.New 函数接收合适的错误信息来创建，像下面这样：

```Go
err := errors.New("math - square root of negative number")
func Sqrt(f float64) (float64, error) {
if f < 0 {
        return 0, errors.New ("math - square root of negative number")
    }
}
```
用 fmt 创建错误对象：

通常你想要返回包含错误参数的更有信息量的字符串，例如：可以用 fmt.Errorf() 来实现：它和 fmt.Printf() 完全一样，接收有一个或多个格式占位符的格式化字符串和相应数量的占位变量。和打印信息不同的是它用信息生成错误对象。
比如在前面的平方根例子中使用：

```Go
if f < 0 {
    return 0, fmt.Errorf("square root of negative number %g", f)
}
```

## 15.2 Panic
在多层嵌套的函数调用中调用 panic，可以马上中止当前函数的执行，所有的 defer 语句都会保证执行并把控制权交还给接收到 panic 的函数调用者。这样向上冒泡直到最顶层，并执行（每层的） defer，在栈顶处程序崩溃，并在命令行中用传给 panic 的值报告错误情况：这个终止过程就是 panicking。

标准库中有许多包含 Must 前缀的函数，像 regexp.MustComplie 和 template.Must；当正则表达式或模板中转入的转换字符串导致错误时，这些函数会 panic。

不能随意地用 panic 中止程序，必须尽力补救错误让程序能继续执行。

自定义包中的错误处理和 panicking，这是所有自定义包实现者应该遵守的最佳实践：

1）在包内部，总是应该从 panic 中 recover：不允许显式的超出包范围的 panic()

2）向包的调用者返回错误值（而不是 panic）。

recover() 的调用仅当它在 defer 函数中被直接调用时才有效。

下面主函数recover 了panic：

```Go
func Parse(input string) (numbers []int, err error) {
    defer func() {
        if r := recover(); r != nil {
            var ok bool
            err, ok = r.(error)
            if !ok {
                err = fmt.Errorf("pkg: %v", r)
            }
        }
    }()

    fields := strings.Fields(input)
    numbers = fields2numbers(fields)
    return
}

func fields2numbers(fields []string) (numbers []int) {
    if len(fields) == 0 {
        panic("no words to parse")
    }
    for idx, field := range fields {
        num, err := strconv.Atoi(field)
        if err != nil {
            panic(&ParseError{idx, field, err})
        }
        numbers = append(numbers, num)
    }
    return
}
```

## 15.3 Recover：从 panic 中恢复

正如名字一样，这个（recover）内建函数被用于从 panic 或 错误场景中恢复：让程序可以从 panicking 重新获得控制权，停止终止过程进而恢复正常执行。
recover 只能在 defer 修饰的函数中使用：用于取得 panic 调用中传递过来的错误值，如果是正常执行，调用 recover 会返回 nil，且没有其它效果。
总结：panic 会导致栈被展开直到 defer 修饰的 recover() 被调用或者程序中止。

```Go
func protect(g func()) {
    defer func() {
        log.Println("done")
        // 即使有panic，Println也正常执行。
        if err := recover(); err != nil {
        	log.Printf("run time panic: %v", err)
        }
    }()
    log.Println("start")
    g() //   可能发生运行时错误的地方
}
```

## 15.4 有关于defer

说到错误处理，就不得不提defer。先说说它的规则：

* 规则一 当defer被声明时，其参数就会被实时解析
* 规则二 defer执行顺序为先进后出
* 规则三 defer可以读取有名返回值

必须要先声明defer，否则不能捕获到panic异常。recover() 的调用仅当它在 defer 函数中被直接调用时才有效。

panic 是用来表示非常严重的不可恢复的错误的。在Go语言中这是一个内置函数，接收一个interface{}类型的值（也就是任何值了）作为参数。

函数执行的时候panic了，函数不往下走，开始运行defer，defer处理完再返回。这时候（defer的时候），recover内置函数可以捕获到当前的panic（如果有的话），被捕获到的panic就不会向上传递了。
recover之后，逻辑并不会恢复到panic那个点去，函数还是会在defer之后返回。

大致过程：

Panic--->defer-->recover

```Go
// 规则一，当defer被声明时，其参数就会被实时解析
package main

import "fmt"

func main() {
	var i int = 1

	defer fmt.Println("result =>", func() int { return i * 2 }())
	i++
	// 输出: result => 2 (而不是 4)
}
```

```Go
// 规则二 defer执行顺序为先进后出

package main

import "fmt"

func main() {

	defer fmt.Print(" !!! ")
	defer fmt.Print(" world ")
	fmt.Print(" hello ")

}
//输出:  hello  world  !!!
```

使用defer计算函数执行时间

```Go
package main
import(
        "fmt"
        "time"
)

func main(){
        defer timeCost(time.Now())
        fmt.Println("start program")
        time.Sleep(5*time.Second)
        fmt.Println("finish program")
}

func timeCost(start time.Time){
        terminal:=time.Since(start)
        fmt.Println(terminal)
}
```
另外一种计算函数执行时间方法：

在对比和基准测试中，我们需要知道一个计算执行消耗的时间。最简单的一个办法就是在计算开始之前设置一个起始时候，再由计算结束时的结束时间，最后取出它们的差值，就是这个计算所消耗的时间。想要实现这样的做法，可以使用 time 包中的 Now() 和 Sub 函数：
```Go
start := time.Now()
longCalculation()
end := time.Now()
delta := end.Sub(start)
fmt.Printf("longCalculation took this amount of time: %s\n", delta)
```