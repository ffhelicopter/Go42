# 《Go语言四十二章经》第十六章 函数

作者：李骁

## 16.1 函数分类
Go 里面有三种类型的函数：

* 普通的带有名字的函数
* 匿名函数或者lambda函数 
* 方法（Methods）

除了main()、init()函数外，其它所有类型的函数都可以有参数与返回值。

函数参数、返回值以及它们的类型被统称为函数签名。

函数重载（function overloading）指的是可以编写多个同名函数，只要它们拥有不同的形参或者不同的返回值，在 Go 里面函数重载是不被允许的。

如果需要申明一个在外部定义的函数，你只需要给出函数名与函数签名，不需要给出函数体：

```Go
func flushICache(begin, end uintptr) 
```
函数也可以以申明的方式被使用，作为一个函数类型，就像：

```Go
type binOp func(int, int) int
```
在这里，不需要函数体 {}。

函数是一等值（first-class value）：它们可以赋值给变量，就像下面一样：

```Go
add := binOp 
```
这个变量知道自己指向的函数的签名，所以给它赋一个具有不同签名的函数值是不可能的。

函数值（functions value）之间可以相互比较：如果它们引用的是相同的函数或者都是 nil 的话，则认为它们是相同的函数。函数不能在其它函数里面声明（不能嵌套），不过我们可以通过使用匿名函数来破除这个限制。

没有参数的函数通常被称为 无参数函数（niladic function），就像 main.main()

## 16.2 函数调用

* 按值传递（call by value）
* 按引用传递（call by reference）

Go 默认使用按值传递来传递参数，也就是传递参数的副本。函数接收参数副本之后，在使用变量的过程中可能对副本的值进行更改，但不会影响到原来的变量，比如 Function(arg1)。

如果你希望函数可以直接修改参数的值，而不是对参数的副本进行操作，你需要将参数的地址（变量名前面添加&符号，比如 &variable）传递给函数，这就是按引用传递，比如 Function(&arg1)，此时传递给函数的是一个指针。如果传递给函数的是一个指针，指针的值（一个地址）会被复制，但指针的值所指向的地址上的值不会被复制；我们可以通过这个指针的值来修改这个值所指向的地址上的值。

在函数调用时，像切片（slice）、字典（map）、接口（interface）、通道（channel）这样的引用类型都是默认使用引用传递（即使没有显式的指出指针）

命名返回值作为结果形参（result parameters）被初始化为相应类型的零值，当需要返回的时候，我们只需要一条简单的不带参数的return语句。需要注意的是，即使只有一个命名返回值，也需要使用 () 括起来

如果函数的最后一个参数是采用 ...type 的形式，那么这个函数就可以处理一个变长的参数，这个长度可以为 0，这样的函数称为变参函数。

这个函数接受一个类似某个类型的 slice 的参数，该参数可以通过 for 循环结构迭代。

```Go
func min(s ...int) int {
    if len(s)==0 {
        return 0
    }
    min := s[0]
    for _, v := range s {
        if v < min {
            min = v
        }
    }
    return min
}
```

## 16.3 内置函数
Go 语言拥有一些不需要进行导入操作就可以使用的内置函数。它们有时可以针对不同的类型进行操作：

|名称|说明|
|:--|:--|
|close	|用于通道通信|
|len、cap	|len 用于返回某个类型的长度或数量（字符串、数组、切片、map 和通道）；cap 是容量的意思，用于返回某个类型的最大容量（只能用于切片和 map）|
|new、make	|new 和 make 均是用于分配内存：new 用于值类型和用户定义的类型，如自定义结构，make 用于内置引用类型（切片、map 和通道）。它们的用法就像是函数，但是将类型作为参数：new(type)、make(type)。new(T) 分配类型 T 的零值并返回其地址，也就是指向类型 T 的指针。它也可以被用于基本类型：v := new(int)。make(T) 返回类型 T 的初始化之后的值，因此它比 new 进行更多的工作。 new() 是一个函数，不要忘记它的括号。二者都是内存的分配（堆上），但是make只用于slice、map以及channel的初始化（非零值）；而new用于类型的内存分配，并且内存置为零。|
|copy、append	|用于复制和连接切片|
|panic、recover 	|两者均用于错误处理机制|

## 16.4 递归与回调
使用递归函数时经常会遇到的一个重要问题就是栈溢出：一般出现在大量的递归调用导致的程序栈内存分配耗尽。这个问题可以通过一个名为懒惰求值的技术解决，在 Go 语言中，我们可以使用通道（channel）和 goroutine。

Go 语言中也可以使用相互调用的递归函数：多个函数之间相互调用形成闭环。因为 Go 语言编译器的特殊性，这些函数的声明顺序可以是任意的。

函数可以作为其它函数的参数进行传递，然后在其它函数内调用执行，一般称之为回调。

```Go
package main
import (
    "fmt"
)
func main() {
    callback(1, Add)
}
func Add(a, b int) {
    fmt.Printf("The sum of %d and %d is: %d\n", a, b, a+b)
}
func callback(y int, f func(int, int)) {
    f(y, 2) // 实际上是 Add(1, 2)
}
```

## 16.5 匿名函数
当我们不希望给函数起名字的时候，可以使用匿名函数，例如：

```Go
func(x, y int) int { return x + y }
```
这样的一个函数不能够独立存在（编译器会返回错误：non-declaration statement outside function body），但可以被赋值于某个变量，即保存函数的地址到变量中：fplus := func(x, y int) int { return x + y }，然后通过变量名对函数进行调用：fplus(3, 4)。

当然，也可以直接对匿名函数进行调用：

```Go
func(x, y int) int { return x + y } (3, 4)
```
下面是一个计算从 1 到 1 百万整数的总和的匿名函数：

```Go
func() {
    sum := 0
    for i := 1; i <= 1e6; i++ {
        sum += i
    }
}()
```
表示参数列表的第一对括号必须紧挨着关键字 func，因为匿名函数没有名称。花括号 {} 涵盖着函数体，最后的一对括号表示对该匿名函数的调用。

## 16.6 闭包函数

匿名函数同样被称之为闭包（函数式语言的术语）：它们被允许调用定义在其它环境下的变量。闭包可使得某个函数捕捉到一些外部状态，例如：函数被创建时的状态。另一种表示方式为：一个闭包继承了函数所声明时的作用域。这种状态（作用域内的变量）都被共享到闭包的环境中，因此这些变量可以在闭包中被操作，直到被销毁。闭包经常被用作包装函数：它们会预先定义好 1 个或多个参数以用于包装。另一个不错的应用就是使用闭包来完成更加简洁的错误检查。

仅仅从形式上将闭包简单理解为匿名函数是不够的，还需要理解闭包实质上的含义。

实质上看，闭包是由函数及其相关引用环境组合而成的实体(即：闭包=函数+引用环境)。闭包在运行时可以有多个实例，不同的引用环境和相同的函数组合可以产生不同的实例。由闭包的实质含义，我们可以推论：闭包获取捕获变量相当于引用传递，而非值传递；对于闭包函数捕获的常量和变量，无论闭包何时何处被调用，闭包都可以使用这些常量和变量，而不用关心它们表面上的作用域。

应用闭包：将函数作为返回值，我们用一个例子来进行验证。

```Go
package main

import (
	"fmt"
)

func addNumber(x int) func(int) {
	fmt.Printf("x: %d, addr of x:%p\n", x, &x)
	return func(y int) {
		k := x + y
		x = k
		y = k
		fmt.Printf("x: %d, addr of x:%p\n", x, &x)
		fmt.Printf("y: %d, addr of y:%p\n", y, &y)
	}
}

func main() {
	addNum := addNumber(5)
	addNum(1)
	addNum(1)
	addNum(1)

	fmt.Println("---------------------")

	addNum1 := addNumber(5)
	addNum1(1)
	addNum1(1)
	addNum1(1)
}
```

```Go
程序输出：
x: 5, addr of x:0xc042054080
x: 6, addr of x:0xc042054080
y: 6, addr of y:0xc042054098
x: 7, addr of x:0xc042054080
y: 7, addr of y:0xc0420540d0
x: 8, addr of x:0xc042054080
y: 8, addr of y:0xc0420540e8
---------------------
x: 5, addr of x:0xc042054100
x: 6, addr of x:0xc042054100
y: 6, addr of y:0xc042054110
x: 7, addr of x:0xc042054100
y: 7, addr of y:0xc042054128
x: 8, addr of x:0xc042054100
y: 8, addr of y:0xc042054140
```

首先强调一点，x是闭包中被捕获的变量，y只是闭包内部的局部变量，而非被捕获的变量。因此，对于每一次引用，x的地址都是固定的，是同一个引用变量；y的地址则是变化的。另外，闭包被引用了两次，由此产生了两个闭包实例，即addNum := addNumber(5)和addNum1 :=addNumber(5)是两个不同实例，其中引用的两个x变量也来自两个不同的实例。

## 16.7 使用闭包调试
当您在分析和调试复杂的程序时，无数个函数在不同的代码文件中相互调用，如果这时候能够准确地知道哪个文件中的具体哪个函数正在执行，对于调试是十分有帮助的。您可以使用 runtime 或 log 包中的特殊函数来实现这样的功能。包 runtime 中的函数 Caller() 提供了相应的信息，因此可以在需要的时候实现一个 where() 闭包函数来打印函数执行的位置：

```Go
where := func() {
    _, file, line, _ := runtime.Caller(1)
    log.Printf("%s:%d", file, line)
}
where()
// some code
where()
// some more code
where()
```
或使用一个更加简短版本的 where 函数：

```Go
var where = log.Print
func func1() {
where()
... some code
where()
... some code
where()
}
```
## 16.8 高阶函数
在定义所需功能时我们可以利用函数可以作为（其它函数的）参数的事实来使用高阶函数

定义一个通用的 Process() 函数，它接收一个作用于每一辆 car 的 f 函数作参数：
```Go
// Process all cars with the given function f:
func (cs Cars) Process(f func(car *Car)) {
    for _, c := range cs {
        f(c)
    }
}
```


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>本书《Go语言四十二章经》内容在简书同步地址：  https://www.jianshu.com/nb/29056963
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
