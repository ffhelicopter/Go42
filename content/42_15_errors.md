# 《Go语言四十二章经》第十五章 错误处理

作者：李骁

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

在Go语言中 panic() 是一个内置函数，用来表示非常严重的不可恢复的错误。必须要先声明defer，否则不能捕获到异常。普通函数在执行的时候发生了异常，则开始运行defer（如有），defer处理完再返回。

在多层嵌套的函数调用中调用 panic()，可以马上中止当前函数的执行，所有的 defer 语句都会保证执行并把控制权交还给接收到异常的函数调用者。这样向上冒泡直到最顶层，并执行（每层的） defer，在栈顶处程序崩溃，并在命令行中用传给异常的值报告错误情况：这个终止过程就是 panicking。

一般不要随意用 panic() 中止程序，必须尽力补救错误让程序能继续执行。

自定义包中的错误处理和 panicking，这是所有自定义包实现者应该遵守的最佳实践：

1）在包内部，总是应该从异常中 recover：不允许显式的超出包范围的 panic()

2）向包的调用者返回错误值。

recover() 函数的调用仅当它在 defer 函数中被直接调用时才有效。

下面主函数捕获了异常：

```Go
package main

import (
	"fmt"
)

func div(a, b int) {

	defer func() {

		if r := recover(); r != nil {
			fmt.Printf("捕获到异常：%s\n", r)
		}
	}()

	if b < 0 {

		panic("除数需要大于0")
	}

	fmt.Println("余数为：", a/b)

}

func main() {
	// 捕捉内部的异常
	div(10, 0)

	// 捕捉主动的异常
	div(10, -1)
}

程序输出：

捕获到异常：runtime error: integer divide by zero
捕获到异常：除数需要大于0
```

## 15.3 Recover：从异常中恢复

recover() 这个内建函数被用于从异常或错误场景中恢复：让程序可以从 panicking 重新获得控制权，停止终止过程进而恢复正常执行。
recover() 只能在 defer 修饰的函数中使用：用于取得异常调用中传递过来的错误值，如果是正常执行，调用 recover() 会返回 nil，且没有其它效果。
总结：异常会导致栈被展开直到 defer 修饰的 recover() 被调用或者程序中止。

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
* 规则三 defer可以读取有名返回值，也就是可以改变有名返回参数的值。

这三个规则用起来需要注意下，避免出现代码陷阱，下面是具体代码：

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

上面讲了两条规则，第三条规则其实也不难理解，只要记住是可以改变有名返回值：

这是由于在Go语言中，return 语句不是原子操作，最先是所有结果值在进入函数时都会初始化为其类型的零值（姑且称为ret赋值），然后执行defer命令,最后才是return操作。如果是有名返回值，返回值变量其实可视为是引用赋值，可以能被defer修改。而在匿名返回值时，给ret的值相当于拷贝赋值，defer命令时不能直接修改。

```Go
func fun1() (i int)
```
上面函数签名中的 i 就是有名返回值，如果fun1()中定义了 defer 代码块，是可以改变返回值 i 的，函数返回语句return i 可以简写为 return 。

这里综合了一下，在下面这个例子里列举了几种情况，可以好好琢磨下：

```Go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("=========================")
	fmt.Println("return:", fun1())

	fmt.Println("=========================")
	fmt.Println("return:", fun2())
	fmt.Println("=========================")

	fmt.Println("return:", fun3())
	fmt.Println("=========================")

	fmt.Println("return:", fun4())
}

func fun1() (i int) {
	defer func() {
		i++
		fmt.Println("defer2:", i) // 打印结果为 defer2: 2
	}()

	// 规则二 defer执行顺序为先进后出

	defer func() {
		i++
		fmt.Println("defer1:", i) // 打印结果为 defer1: 1
	}()

	// 规则三 defer可以读取有名返回值（函数指定了返回参数名）

	return 0 //这里实际结果为2。如果是return 100呢
}

func fun2() int {
	var i int
	defer func() {
		i++
		fmt.Println("defer2:", i) // 打印结果为 defer2: 2
	}()

	defer func() {
		i++
		fmt.Println("defer1:", i) // 打印结果为 defer1: 1
	}()
	return i
}

func fun3() (r int) {
	t := 5
	defer func() {
		t = t + 5
		fmt.Println(t)
	}()
	return t
}

func fun4() int {
	i := 8
	// 规则一 当defer被声明时，其参数就会被实时解析
	defer func(i int) {
		i = 99
		fmt.Println(i)
	}(i)
	i = 19
	return i
}
```

在上面fun1() (i int)有名返回值情况下，return最终返回的实际值和期望的return 0有较大出入。因为在上面fun1() (i int) 中，如果return 100或return 0 ，这样的区别在于i的值实际上分别是100或0。而在上面中，如果return 100，则因为改变了有名返回值i，而defer可以读取有名返回值，所以返回值最终为102，而defer1打印101，defer打印102。因此我们一般直接写为return。

这点要注意，有时函数可能返回非我们希望的值，所以改为匿名返回也是一种办法。具体请看下面输出。



```Go
程序输出：
=========================
defer1: 1
defer2: 2
return: 2
=========================
defer1: 1
defer2: 2
return: 0
=========================
10
return: 5
=========================
99
return: 19

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


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第十四章 流程控制](https://github.com/ffhelicopter/Go42/blob/master/content/42_14_flow.md)

[第十六章 函数](https://github.com/ffhelicopter/Go42/blob/master/content/42_16_function.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com