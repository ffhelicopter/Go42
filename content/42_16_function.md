# 《Go语言四十二章经》第十六章 函数

作者：李骁

## 16.1 函数介绍

Go语言函数基本组成：关键字func、函数名、参数列表、返回值、函数体和返回语句。语法如下：

```Go
func 函数名(参数列表) (返回值列表) {
    // 函数体
return
}
```

除了main()、init()函数外，其它所有类型的函数都可以有参数与返回值。

对于函数，一般也可以这么认为："func" FunctionName Signature [ FunctionBody ] . 

"func" 为定义函数的关键字，FunctionName 为函数名，Signature 为函数签名，FunctionBody 为函数体。以下面定义的函数为例：

```Go
func FunctionName (a typea, b typeb) (t1 type1, t2 type2)
```

函数签名由函数参数、返回值以及它们的类型组成，被统称为函数签名。如：

```Go
(a typea, b typeb) (t1 type1, t2 type2)
```

如果两个函数的参数列表和返回值列表的变量类型能一一对应，那么这两个函数就有相同的签名，下面testa与testb具有相同的函数签名。

```Go
func testa  (a, b int, z float32) bool
func testb  (a, b int, z float32) (bool)
```

函数调用传入的参数必须按照参数声明的顺序。而且Go语言没有默认参数值的说法。函数签名中的最后传入参数可以具有前缀为....的类型（...int），这样的参数称为可变参数，并且可以使用零个或多个参数来调用该函数，这样的函数称为变参函数。

```Go
func doFix (prefix string, values ...int)
```

函数的参数和返回值列表始终带括号，但如果只有一个未命名的返回值（且只有此种情况），则可以将其写为未加括号的类型；一个函数也可以拥有多返回值，返回类型之间需要使用逗号分割，并使用小括号 () 将它们括起来。

```Go
func testa  (a, b int, z float32) bool
func swap  (a int, b int) (t1 int, t2 int)
```

在函数体中，参数是局部变量，被初始化为调用者传入的值。函数的参数和具名返回值是函数最外层的局部变量，它们的作用域就是整个函数。如果函数的签名声明了返回值，则函数体的语句列表必须以终止语句结束。

```Go
func IndexRune(s string, r rune) int {
	for i, c := range s {
		if c == r {
			return i
		}
	}
	return // 必须要有终止语句，如果这里没有return，则会编译错误：missing return at end of function
}
```

函数重载（function overloading）指的是可以编写多个同名函数，只要它们拥有不同的形参或者不同的返回值，在 Go 语言里面函数重载是不被允许的。

函数也可以作为函数类型被使用。函数类型也就是函数签名，函数类型表示具有相同参数和结果类型的所有函数的集合。函数类型的未初始化变量的值为nil。就像下面：

```Go
type  funcType func (int, int) int
```

上面通过type关键字，定义了一个新类型，函数类型 funcType 。

函数也可以在表达式中赋值给变量，这样作为表达式中右值出现，我们称之为函数值字面量（function literal），函数值字面量是一种表达式，它的值被称为匿名函数，就像下面一样：

```Go
f := func() int { return 7 }  
```

下面代码对以上2种情况都做了定义和调用：

```Go

package main

import (
	"fmt"
	"time"
)

type funcType func(time.Time)     // 定义函数类型funcType

func main() {
	f := func(t time.Time) time.Time { return t } // 方式一：直接赋值给变量
	fmt.Println(f(time.Now()))

	var timer funcType = CurrentTime // 方式二：定义函数类型funcType变量timer
	timer(time.Now())

	funcType(CurrentTime)(time.Now())  // 先把CurrentTime函数转为funcType类型，然后传入参数调用
// 这种处理方式在Go 中比较常见
}

func CurrentTime(start time.Time) {
	fmt.Println(start)
}
```

## 16.2 函数调用

Go 语言中函数默认使用按值传递来传递参数，也就是传递参数的副本。函数接收参数副本之后，在使用变量的过程中可能对副本的值进行更改，但不会影响到原来的变量。

如果我们希望函数可以直接修改参数的值，而不是对参数的副本进行操作，则需要将参数的地址传递给函数，这就是按引用传递，比如 Function(&arg1)，此时传递给函数的是一个指针。如果传递给函数的是一个指针，我们可以通过这个指针来修改对应地址上的变量值。

在函数调用时，像切片（slice）、字典（map）、接口（interface）、通道（channel）等这样的引用类型都是默认使用引用传递。

命名返回值被初始化为相应类型的零值，当需要返回的时候，我们只需要一条简单的不带参数的return语句。需要注意的是，即使只有一个命名返回值，也需要使用 () 括起来

前面说过，函数签名中的最后传入参数可以具有前缀为....的类型（...int），这样的函数称为变参函数。

变参函数可以接受某种类型的切片 slice 为参数：

```Go

package main

import (
	"fmt"
)

// 变参函数，参数不定长
func list(nums ...int) {
	fmt.Println(nums)
}

func main() {
	// 常规调用，参数可以多个
	list(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	// 在参数同类型时，可以组成slice使用 parms... 进行参数传递
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	list(numbers...) // slice时使用
}
```

## 16.3 内置函数

Go 语言拥有一些内置函数，内置函数是预先声明的，它们像任何其他函数一样被调用，内置函数没有标准的类型，因此它们只能出现在调用表达式中，它们不能用作函数值。它们有时可以针对不同的类型进行操作：

|内置函数|说明|
|:--|:--|
|close|	用于通道，对于通道c，内置函数close(c)将不再在通道c上发送值。 如果c是仅接收通道，则会出错。 发送或关闭已关闭的通道会导致运行时错误。 关闭nil通道也会导致运行时错误。|
|new、make|	new 和 make 均是用于分配内存：new用于值类型的内存分配，并且置为零值。make只用于slice、map以及channel这三种引用数据类型的内存分配和初始化。new(T) 分配类型 T 的零值并返回其地址，也就是指向类型 T 的指针。make(T) 它返回类型T的值（不是* T）。|

make()内置函数声明不同类型时的参数以及具体作用请见下面说明：

```Go
调用           T的类型     结果

make(T, n)       slice        T为切片类型，长度和容量都为n
make(T, n, m)     slice        T为切片类型，长度为n，容量为m （n<=m ，否则错误）

make(T)          map        T为字典类型
make(T, n)        map        T为字典类型，初始化n个元素的空间

make(T)          channel      T为通道类型，无缓冲区
make(T, n)        channel      T为通道类型，缓冲区长度为n
```

make()内置函数的实际使用举例见下面代码以及注释：

```Go
s := make([]int, 10, 100)       // slice with len(s) == 10, cap(s) == 100
s := make([]int, 1e3)           // slice with len(s) == cap(s) == 1000
s := make([]int, 1<<63)         // illegal: len(s) is not representable by a value of type int
s := make([]int, 10, 0)         // illegal: len(s) > cap(s)
c := make(chan int, 10)         // channel with a buffer size of 10
m := make(map[string]int, 100)  // map with initial space for approximately 100 elements
```

new(T)内置函数在运行时为该类型的变量分配内存，并返回指向它的类型* T的值。 并对变量初始化。

例如：
```Go
type S struct { a int; b float64 }
new(S)
```
new(S)为S类型的变量分配内存，并初始化（a = 0，b = 0.0），返回包含该位置地址的类型* S的值。

|内置函数	|参数类型	|结果|
|:--|:--|:--|
|len(s)	|string type ，[n]T, *[n]T ，[]T ，map[K]T ，chan T	|string的长度（按照字节计算），数组长度 ，切片长度 ，字典长度 ，通道缓冲区中排队的元素数|
|cap(s) |	[n]T, *[n]T ，[]T ，chan T	|数组长度 ，切片容量 ，通道缓冲区容量|

对于len(s)和cap(s)，如果s为nil值，则两个函数的取值都是0，我们还需要记住一个规则：

```Go
0 <= len(s) <= cap(s)
```

在Go语言中，常量在某些计算条件下也可以通过表达式计算得到。比如：如果s是字符串常量，则表达式len(s)是常量。 如果s的类型是数组或指向数组的指针而表达式不包含通道接收或（非常量）函数调用，则表达式len(s)和cap(s)是常量；否则len和cap的调用不是常量。

```Go
const (
	c1 = imag(2i)                  // imag(2i) = 2.0 是常量
	c2 = len([10]float64{2})         // [10]float64{2} 无函数调用
	c3 = len([10]float64{c1})        // [10]float64{c1} 无函数调用
	c4 = len([10]float64{imag(2i)})   // imag(2i)常量无函数调用
	c5 = len([10]float64{imag(z)})    // 无效: imag(z) 非常量函数调用
)
var z complex128
```

|内置函数|	说明|
|:--|:--|
|append	|用于附加连接切片|
|copy	|用于复制切片|
|delete	|从字典删除元素|

```Go
append(s S, x ...T) S  // T 是类型S的元素
```

append内置函数是变参函数，常常用来附加切片元素，将零或多个值x附加到S类型的切片s，它的可变参数必须是切片类型，并返回结果切片，也就是是S类型。值x传递给类型为...的参数T，其中T 是S的元素类型，并且适用相应的参数传递规则：

```Go
s0 := []int{0, 0}
s1 := append(s0, 2)            // append 附加连接单个元素   s1 == []int{0, 0, 2}
s2 := append(s1, 3, 5, 7)        // append 附加连接多个元素  s2 == []int{0, 0, 2, 3, 5, 7}
s3 := append(s2, s0...)         // append 附加连接切片s0  s3 == []int{0, 0, 2, 3, 5, 7, 0, 0}
s4 := append(s3[3:6], s3[2:]...)  // append 附加切片指定值 s4 == []int{3, 5, 7, 2, 3, 5, 7, 0, 0}

var t []interface{}
t = append(t, 42, 3.1415, "foo")  //  t == []interface{}{42, 3.1415, "foo"}

var b []byte
b = append(b, "bar"...)         // append 附加连接字符串内容  b == []byte{'b', 'a', 'r' }
```

```Go
copy(dst, src []T) int
copy(dst []byte, src string) int
```

copy内置函数常常将切片元素从源src复制到目标dst，并返回复制的元素数。 两个参数必须具有相同的元素类型T，并且必须可以分配给类型为[] T的切片。 复制的元素数量是len（src）和len（dst）的最小值。

作为特殊情况，copy函数还接受可分配给[] byte类型的目标参数，其中source参数为字符串类型。 此种情况将字符串中的字节复制到字节切片中。

```Go
var a = [...]int{0, 1, 2, 3, 4, 5, 6, 7}
var s = make([]int, 6)
var b = make([]byte, 5)
n1 := copy(s, a[0:])            // n1 == 6, s == []int{0, 1, 2, 3, 4, 5}
n2 := copy(s, s[2:])            // n2 == 4, s == []int{2, 3, 4, 5, 4, 5}
n3 := copy(b, "Hello, World!")  // n3 == 5, b == []byte("Hello")
```

```Go
delete(m, k)  //从字典m中删除元素 m[k] 
```

内置函数delete从字典m中删除带有键k的元素。

|内置函数	|说明|
|:--|:--|
|complex	|从浮点实部和虚部构造复数值|
|real	|提取复数值的实部|
|imag	|提取复数值的虚部|

```Go
complex(realPart, imaginaryPart floatT) complexT
real(complexT) floatT
imag(complexT) floatT
```

内置函数complex用浮点实部和虚部构造复数值，而real和imag则提取复数值的实部和虚部。

对于complex，两个参数必须是相同的浮点类型，返回类型是具有相应浮点组成的复数类型。float32用于complex64参数，float64用于complex128参数。如果其中一个参数求值为无类型常量，则首先将其转换为另一个参数的类型。如果两个参数都计算为无类型常量，则它们必须是非复数或其虚部必须为零，并且函数的返回值是无类型复数常量。

对于real和imag，参数必须是复数类型，返回类型是相应的浮点类型：float32一般为complex64返回类型，float64一般为complex128返回类型。如果参数求值为无类型常量，则它必须是数字，并且函数的返回值是无类型浮点常量。

real和imag函数一起形成复数的逆，因此对于复数类型Z的值z，z == Z(complex(real(z)，imag(z)))。

如果这些函数的操作数都是常量，则返回值是常量。

```Go
var a = complex(2, -2)             // complex128
const b = complex(1.0, -1.4)        // 无类型complex 常量 1 - 1.4i
x := float32(math.Cos(math.Pi/2))   // float32
var c64 = complex(5, -x)          // complex64
var s uint = complex(1, 0)         // 无类型 complex 常量 1 + 0i 可以转为uint
var rl = real(c64)                // float32
var im = imag(a)                // float64
const c = imag(b)               // 无类型常量 -1.4
```

|内置函数	|说明|
|:--|:--|
|panic	|用来表示非常严重的不可恢复的异常错误|
|recover	|用于从 panic 或 错误场景中恢复|

```Go
func panic(interface{})
func recover() interface{}
```

panic和recover两个内置函数，协助报告和处理运行时异常和程序定义的错误。

在执行函数F时，显式调用panic或者运行时发生panic都会终止F的执行。然后，由F延迟（defer）的任何函数都照常执行。 依此类推，直到执行goroutine中的顶级函数延迟。 此时，程序终止并报告错误条件，包括panic参数的值。

```Go
panic(42)
panic("unreachable")
panic(Error("cannot parse"))
```

recover函数允许程序管理发生panic的goroutine的行为。

另外，Go语言中提供了几个在引导期间有用的内置函数。 这些函数不保证会保留在Go语言中，一般不建议使用。

```Go
print      打印所有参数
println    打印所有参数并换行
```

## 16.4 递归与回调

函数直接或间接调用函数本身，则该函数称为递归函数。使用递归函数时经常会遇到的一个重要问题就是栈溢出：一般出现在大量的递归调用导致的内存分配耗尽。有时我们可以通过循环来解决：

```Go
package main

import "fmt"

// Factorial函数递归调用
func Factorial(n uint64)(result uint64) {
    if (n > 0) {
        result = n * Factorial(n-1)
        return result
    }
    return 1
}

// Fac2函数循环计算
func Fac2(n uint64) (result uint64) {
	result = 1
	var un uint64 = 1
	for i := un; i <= n; i++ {
		result *= i
	}
	return
}

func main() {  
    var i uint64= 7
    fmt.Printf("%d 的阶乘是 %d\n", i, Factorial(i)) 
    fmt.Printf("%d 的阶乘是 %d\n", i, Fac2(i))
}

程序输出：
7 的阶乘是 5040
7 的阶乘是 5040
```

Go 语言中也可以使用相互调用的递归函数：多个函数之间相互调用形成闭环。因为 Go 语言编译器的特殊性，这些函数的声明顺序可以是任意的。

Go语言中函数可以作为其它函数的参数进行传递，然后在其它函数内调用执行，一般称之为回调。

```Go
package main

import (
	"fmt"
)

func main() {
	callback(1, Add)
}

func Add(a, b int) {
	fmt.Printf("%d 与 %d 相加的和是: %d\n", a, b, a+b)
}

func callback(y int, f func(int, int)) {
	f(y, 2) // 回调函数f
}

程序输出：
1 与 2 相加的和是: 3
```

## 16.5 匿名函数

函数值字面量是一种表达式，它的值被称为匿名函数。从形式上看当我们不给函数起名字的时候，可以使用匿名函数，例如：

```Go
func(x, y int) int { return x + y }
```

这样的函数不能够独立存在，但可以被赋值于某个变量，即保存函数的地址到变量中：

```Go
fplus := func(x, y int) int { return x + y }
```

然后通过变量名对函数进行调用：

```Go
fplus(3, 4)
```

当然，也可以直接对匿名函数进行调用，注意匿名函数的最后面加上了括号并填入了参数值，如果没有参数，也需要加上空括号，代表直接调用：

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

参数列表的第一对括号必须紧挨着关键字 func，因为匿名函数没有名称。花括号 {} 涵盖着函数体，最后的一对括号表示对该匿名函数的调用。

下面代码演示了上面的几种情况：

```Go

package main

import (
	"fmt"
)

func main() {
	fn := func() {
		fmt.Println("hello")
	}
	fn()

	fmt.Println("匿名函数加法求和：", func(x, y int) int { return x + y }(3, 4))

	func() {
		sum := 0
		for i := 1; i <= 1e6; i++ {
			sum += i
		}
		fmt.Println("匿名函数加法循环求和：", sum)
	}()
}

程序输出：
hello
匿名函数加法求和： 7
匿名函数加法循环求和： 500000500000
```

# 16.6 闭包函数

匿名函数同样也被称之为闭包。

闭包可被允许调用定义在其环境下的变量，可以访问它们所在的外部函数中声明的所有局部变量、参数和声明的其他内部函数。闭包继承了函数所声明时的作用域，作用域内的变量都被共享到闭包的环境中，因此这些变量可以在闭包中被操作，直到被销毁。也可以理解为内层函数引用了外层函数中的变量或称为引用了自由变量。

实质上看，闭包是由函数及其相关引用环境组合而成的实体(即：闭包=函数+引用环境)。闭包在运行时可以有多个实例，不同的引用环境和相同的函数组合可以产生不同的实例。由闭包的实质含义，我们可以推论：闭包获取捕获变量相当于引用传递，而非值传递；对于闭包函数捕获的常量和变量，无论闭包何时何处被调用，闭包都可以使用这些常量和变量，而不用关心它们表面上的作用域。

换句话说闭包函数可以访问不是它自己内部的变量（这个变量在其它作用域内声明），且这个变量是未赋值的，它在闭包里面赋值。

我们通过下面代码来看看闭包的使用：
	
```Go

package main

import "fmt"

var G int = 7

func main() {
	// 影响全局变量G，代码块状态持续
	y := func() int {
		fmt.Printf("G: %d, G的地址:%p\n", G, &G)
		G += 1
		return G
	}
	fmt.Println(y(), y)
	fmt.Println(y(), y)
	fmt.Println(y(), y) //y的地址

	// 影响全局变量G，注意z的匿名函数是直接执行，所以结果不变
	z := func() int {
		G += 1
		return G
	}()
	fmt.Println(z, &z)
	fmt.Println(z, &z)
	fmt.Println(z, &z)

	// 影响外层（自由）变量i，代码块状态持续
	var f = N()
	fmt.Println(f(1), &f)
	fmt.Println(f(1), &f)
	fmt.Println(f(1), &f)

	var f1 = N()
	fmt.Println(f1(1), &f1)

}

func N() func(int) int {
	var i int
	return func(d int) int {
		fmt.Printf("i: %d, i的地址:%p\n", i, &i)
		i += d
		return i
	}
}


程序输出：
G: 7, G的地址:0x54b1e8
8 0x490340
G: 8, G的地址:0x54b1e8
9 0x490340
G: 9, G的地址:0x54b1e8
10 0x490340
11 0xc0000500c8
11 0xc0000500c8
11 0xc0000500c8
i: 0, i的地址:0xc0000500e8
1 0xc000078020
i: 1, i的地址:0xc0000500e8
2 0xc000078020
i: 2, i的地址:0xc0000500e8
3 0xc000078020
i: 0, i的地址:0xc000050118
1 0xc000078028
```

首先强调一点，G是闭包中被捕获的全局变量，因此，对于每一次引用，G的地址都是固定的，i是函数内部局部变量，地址也是固定的，他们都可以被闭包保持状态并修改。还要注意，f和f1是不同的实例，它们的地址是不一样的。

## 16.7 变参函数

可变参数也就是不定长参数，支持可变参数列表的函数可以支持任意个传入参数，比如fmt.Println函数就是一个支持可变长参数列表的函数。

```Go

package main

import "fmt"

func Greeting(who ...string) {
	for k, v := range who {

		fmt.Println(k, v)
	}
}

func main() {
	s := []string{"James", "Jasmine"}
	Greeting(s...)  // 注意这里切片s... ，把切片打散传入，与s具有相同底层数组的值。
}

程序输出：
0 James
1 Jasmine
```


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第十五章 错误处理](https://github.com/ffhelicopter/Go42/blob/master/content/42_15_errors.md)

[第十七章 Type关键字](https://github.com/ffhelicopter/Go42/blob/master/content/42_17_type.md)




>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42 
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com 


