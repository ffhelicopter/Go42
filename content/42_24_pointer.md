# 《Go语言四十二章经》第二十四章 指针和内存

作者：李骁

## 24.1 指针

一个指针变量可以指向任何一个值的内存地址。指针变量在 32 位计算机上占用 4B 内存，在 64 位计算机占用 8B内存，并且与它所指向的值的大小无关，因为指针变量只是地址的值而已。可以声明指针指向任何类型的值来表明它的原始性或结构性，也可以在指针类型前面加上\*号（前缀）来获取指针所指向的内容。

在Go语言中，指针类型表示指向给定类型（称为指针的基础类型）的变量的所有指针的集合。 符号 \* 可以放在一个类型前，如 \*T，那么它将以类型T为基础，生成指针类型\*T。未初始化指针的值为nil。例如：

```Go
type Point3D struct{ x, y, z float64 }
var pointer *Point3D
var i *[4]int
```

上面定义了两个指针类型变量。它们的值为nil，这时对它们的反向引用是不合法的，并且会使程序崩溃。

```Go
xx := (*pointer).x
panic: runtime error: invalid memory address or nil pointer dereference
```

符号 \* 可以放在一个指针前，如 (\*pointer)，那么它将得到这个指针指向地址上所存储的值，这称为反向引用。不过在Go语言中，(\*pointer).x可以简写为pointer.x。

对于任何一个变量 var， 表达式var == \*(&var)都是正确的。

注意：不能得到一个数字或常量的地址，下面的写法是错误的：

```Go
const i = 5
ptr := &i // error: cannot take the address of i
ptr2 := &10 // error: cannot take the address of 10
```

虽然Go 语言和 C、C++ 这些语言一样，都有指针的概念，但是指针运算在语法上是不允许的。这样做的目的是保证内存安全。从这一点看，Go 语言的指针基本就是一种引用。

指针的一个高级应用是可以传递一个变量的引用（如函数的参数），这样不会传递变量的副本。当调用函数时，如果参数为基础类型，传进去的是值，也就是另外复制了一份参数到当前的函数调用栈。参数为引用类型时，传进去的基本都是引用。而指针传递的成本很低，只占用 4B或 8B内存。

如果代码在运行中需要占用大量的内存，或很多变量，或者两者都有，这时使用指针会减少内存占用和提高运行效率。被指向的变量保存在内存中，直到没有任何指针指向它们。所以从它们被创建开始就具有相互独立的生命周期。
 
内存管理中的内存区域一般包括堆内存（heap）和栈内存（stack）， 栈内存主要用来存储当前调用栈用到的简单类型数据，如string，bool，int，float 等。这些类型基本上较少占用内存，容易回收，因此可以直接复制，进行垃圾回收时也比较容易做针对性的优化。 而复杂的复合类型占用的内存往往相对较大，存储在堆内存中，垃圾回收频率相对较低，代价也较大，因此传引用或指针可以避免进行成本较高的复制操作，并且节省内存，提高程序运行效率。
 
因此，在需要改变参数的值或者避免复制大批量数据而节省内存时（也会提高运行效率，毕竟大批量复制也耗费时间）都会选择使用指针。

另一方面，指针的频繁使用也会导致性能下降。指针也可以指向另一个指针，并且可以进行任意深度的嵌套，形成多级的间接引用，但会使代码结构不清晰。

在大多数情况下，Go 语言可以使程序员轻松创建指针，并且隐藏间接引用，如：自动反向引用。

**指针的使用方法：**

* 定义指针变量；

* 为指针变量赋值；

* 访问指针变量中指向地址的值；

* 在指针类型前面加上\*号来获取指针所指向的内容。

```Go
package main

import "fmt"

func main() {
	var a, b int = 20, 30 // 声明实际变量
	var ptra *int         // 声明指针变量
	var ptrb *int = &b

	ptra = &a // 指针变量的存储地址

	fmt.Printf("a  变量的地址是: %x\n", &a)
	fmt.Printf("b  变量的地址是: %x\n", &b)

	// 指针变量的存储地址
	fmt.Printf("ptra  变量的存储地址: %x\n", ptra)
	fmt.Printf("ptrb  变量的存储地址: %x\n", ptrb)

	// 使用指针访问值
	fmt.Printf("*ptra  变量的值: %d\n", *ptra)
	fmt.Printf("*ptrb  变量的值: %d\n", *ptrb)
}
```
## 24.2 new() 和 make() 的区别


new() 和 make() 都在堆上分配内存，但是它们的行为不同，适用于不同的类型。

new() 用于值类型的内存分配，并且置为零值。
make() 只用于切片、字典以及通道这三种引用数据类型的内存分配和初始化。

new(T) 分配类型 T 的零值并返回其地址，也就是指向类型 T 的指针。
make(T) 返回类型T的值（不是* T）。

然而在Go语言中，并不能准确判断变量是分配到栈还是堆上。在C++中，使用new()创建的变量总是在堆上。在Go中变量的位置是由编译器决定的。编译器根据变量的大小和泄露（逃逸）分析的结果来决定其位置。

如果想确切知道变量分配的位置，可在执行go build或go run时加上-m gc标志（即go run -gcflags -m app.go）。例如：


```Go
go run -gcflags -m main.go
# command-line-arguments
.\main.go:12:31: m.Alloc / 1024 escapes to heap
.\main.go:11:23: main &m does not escape
.\main.go:12:12: main ... argument does not escape
```

## 24.3 垃圾回收和 SetFinalizer

Go 语言开发者一般不需要写代码来释放不再使用的变量或结构体占用的内存，在 Go语言运行时有一个独立的进程，即垃圾收集器（GC），会专门处理这些事情，它搜索不再使用的变量然后释放它们占用的内存，这是自动垃圾回收；还有一种是主动垃圾回收，通过显式调用 runtime.GC()来实现。

通过调用 runtime.GC() 函数可以显式的触发 GC，这在某些的场景下非常有用，比如当内存资源不足时调用 runtime.GC()，它会在此函数执行的点上立即释放一大片内存，但此时程序可能会有短时的性能下降（因为 GC 进程在执行）。

下面代码中的func (p *Person) NewOpen()在某些情况下非常有必要这样处理，比如某些资源占用申请，开发人员可能忘记使用defer Close()来销毁处理，但通过SetFinalizer，如果GC自动运行或者手动运行GC，则都能及时销毁这些资源，释放占用的内存而避免内存泄漏。

GC过程中重要的函数func SetFinalizer(obj interface{}, finalizer interface{})有两个参数，参数一：obj必须是指针类型。参数二：finalizer是一个函数，其参数类型是obj的类型，其没有返回值。


```Go
package main

import (
	"log"
	"runtime"
	"time"
)

type Person struct {
	Name string
	Age  int
}

func (p *Person) Close() {
	p.Name = "NewName"
	log.Println(p)
	log.Println("Close")
}

func (p *Person) NewOpen() {
	log.Println("Init")
	runtime.SetFinalizer(p, (*Person).Close)
}

func Tt(p *Person) {
	p.Name = "NewName"
	log.Println(p)
	log.Println("Tt")
}

// 查看内存情况
func Mem(m *runtime.MemStats) {
	runtime.ReadMemStats(m)
	log.Printf("%d Kb\n", m.Alloc/1024)
}

func main() {
	var m runtime.MemStats
	Mem(&m)

	var p *Person = &Person{Name: "lee", Age: 4}
	p.NewOpen()
	log.Println("Gc完成第一次")
	log.Println("p:", p)
	runtime.GC()
	time.Sleep(time.Second * 5)
	Mem(&m)

	var p1 *Person = &Person{Name: "Goo", Age: 9}
	runtime.SetFinalizer(p1, Tt)
	log.Println("Gc完成第二次")
	time.Sleep(time.Second * 2)
	runtime.GC()
	time.Sleep(time.Second * 2)
	Mem(&m)

}

```


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第二十三章 同步与锁](https://github.com/ffhelicopter/Go42/blob/master/content/42_23_sync.md)

[第二十五章 面向对象](https://github.com/ffhelicopter/Go42/blob/master/content/42_25_oo.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com