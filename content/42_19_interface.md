# 《Go语言四十二章经》第十九章 接口

作者：李骁

## 19.1 接口是什么

Go语言接口定义了一组方法集合，但是这些方法集合仅仅只是被定义，它们没有在接口中实现。接口(interface)类型是Go语言的一种数据类型。而因为所有的类型包括自定义类型都实现了空接口interface{}，所以空接口interface{}可以被当做任意类型的数值。

Go 语言中的所有类型包括自定义类型都实现了interface{}接口，这意味着所有的类型如string、 int、 int64甚至是自定义的结构体类型都拥有interface{}空接口，这一点interface{}和Java中的Object类比较相似。

接口类型的未初始化变量的值为nil。

```Go
var i interface{} = 99 // i可以是任何类型
i = 44.09
i = "All"  // i 可接受任意类型的赋值
```

接口是一组抽象方法的集合，它必须由其他非接口类型实现，不能自我实现。Go 语言通过它可以实现很多面向对象的特性。

通过如下格式定义接口：

```Go
type Namer interface {
    Method1(param_list) return_type
    Method2(param_list) return_type
    ...
}
```

上面的 Namer 是一个接口类型，按照惯例，单方法接口由方法名称加上-er后缀或类似修改来命名，以构造代理名词：Reader，Writer，Formatter，CloseNotifier等。还有一些不常用的方式（当后缀 er 不合适时），比如 Recoverable，此时接口名以 able 结尾，或者以 I 开头等。

Go 语言中的接口都很简短，通常它们会包含 0 个、最多 3 个方法。如标准库io包中定义了下面2个接口：

```Go
type Reader interface {
	Read(p []byte) (n int, err error)
}
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

在Go语言中，如果接口的所有方法在某个类型方法集中被实现，则认为该类型实现了这个接口。

类型不用显式声明实现了接口，只需要实现接口所有方法，这样的隐式实现解藕了实现接口的包和定义接口的包。

同一个接口可被多个类型可以实现，一个类型也可以实现多个接口。实现了某个接口的类型，还可以有其它的方法。有时我们甚至都不知道某个类型定义的方法集巧合地实现了某个接口。这种灵活性使我们不用像JAVA语言那样需要显式implement，一旦类型不需要实现某个接口，我们甚至可以不改动任何代码。

类型需要实现接口方法集中的所有方法，一定是接口方法集中所有方法。类型实现了这个接口，那么接口类型的变量也就可以存放该类型的值。

如下代码所示，结构体A和类型I都实现了接口B的方法f()，所有这两种类型也具有了接口B的一切特性，可以将该类型的值存储在接口B类型的变量中：

```Go
package main

import (
	"fmt"
)

type A struct {
	Books int
}

type B interface {
	f()
}

func (a A) f() {
	fmt.Println("A.f() ", a.Books)
}

type I int

func (i I) f() {
	fmt.Println("I.f() ", i)
}

func main() {
	var a A = A{Books: 9}
	a.f()

	var b B = A{Books: 99} // 接口类型可接受结构体A的值，因为结构体A实现了接口
	b.f()

	var i I = 199 // I是int类型引申出来的新类型
	i.f()

	var b2 B = I(299) // 接口类型可接受新类型I的值，因为新类型I实现了接口
	b2.f()
}

程序输出：
A.f()  9
A.f()  99
I.f()  199
I.f()  299
```

如果接口在类型之后才定义，或者二者处于不同的包中。但只要类型实现了接口中的所有方法，这个类型就实现了此接口。

因此Go语言中接口具有强大的灵活性。

注意：接口中的方法必须要全部实现，才能实现接口。

## 19.2 接口嵌入

一个接口可以包含一个或多个其他的接口，但是在接口内不能嵌入结构体，也不能嵌入接口自身，否则编译会出错。

下面这两种嵌入接口自身的方式都不能编译通过:

```Go
// 编译错误：invalid recursive type Bad
type Bad interface {
	Bad
}

// 编译错误：invalid recursive type Bad2
type Bad1 interface {
	Bad2
}
type Bad2 interface {
	Bad1
}
```

比如下面的接口 File 包含了 ReadWrite 和 Lock 的所有方法，它还额外有一个 Close() 方法。接口的嵌入方式和结构体的嵌入方式语法上差不多，直接写接口名即可。

```Go
type ReadWrite interface {
    Read(b Buffer) bool
    Write(b Buffer) bool
}

type Lock interface {
    Lock()
    Unlock()
}

type File interface {
    ReadWrite
    Lock
    Close()
}
```

## 19.3 类型断言

前面我们可以把实现了某个接口的类型值保存在接口变量中，但反过来某个接口变量属于哪个类型呢？如何检测接口变量的类型呢？这就是类型断言（Type Assertion）的作用。

接口类型I的变量 varI 中可以包含任何实现了这个接口的类型的值，如果多个类型都实现了这个接口，所以有时我们需要用一种动态方式来检测它的真实类型，即在运行时确定变量的实际类型。

通常我们可以使用类型断言（value, ok := element.(T)）来测试在某个时刻接口变量 varI 是否包含类型 T 的值：

```Go
value, ok := varI.(T)       // 类型断言
```
**varI 必须是一个接口变量**，否则编译器会报错：invalid type assertion: varI.(T) (non-interface type (type of I) on left) 。

类型断言可能是无效的，虽然编译器会尽力检查转换是否有效，但是它不可能预见所有的可能性。如果转换在程序运行时失败会导致错误发生。更安全的方式是使用以下形式来进行类型断言：

```Go
var varI I
varI = T("Tstring")
if v, ok := varI.(T); ok { // 类型断言
	fmt.Println("varI类型断言结果为：", v) // varI已经转为T类型
	varI.f()
}
```

如果断言成功，v 是 varI 转换到类型 T 的值，ok 会是 true；否则 v 是类型 T 的零值，ok 是 false，也没有运行时错误发生。

接口类型向普通类型转换有两种方式：Comma-ok断言和Type-switch测试。

**通过Type-switch做类型判断**

接口变量的类型可以使用一种特殊形式的 switch 做类型断言：

```Go
// Type-switch做类型判断
var value interface{}

switch str := value.(type) {
case string:
	fmt.Println("value类型断言结果为string:", str)

case Stringer:
	fmt.Println("value类型断言结果为Stringer:", str)

default:
	fmt.Println("value类型不在上述类型之中")
}
```

可以用 Type-switch 进行运行时类型分析，但是在 type-switch 时不允许有 fallthrough 。Type-switch让我们在处理未知类型的数据时，比如解析 json 等编码的数据，会非常方便。

**测试一个值是否实现了某个接口（Comma-ok断言）**

我们想测试它是否实现了 I 接口，可以这样做类型断言：

```Go
// Comma-ok断言
var varI I
varI = T("Tstring")
if v, ok := varI.(T); ok { // 类型断言
	fmt.Println("varI类型断言结果为：", v) // varI已经转为T类型
	varI.f()
}
```

接口描述了一系列的行为，规定可以做什么行为，“当一个东西，走起来像鸭子，叫起来也像鸭子，游泳也像鸭子，那么我们可以认为他就是一只鸭子”。类型实现不同的接口将拥有不同的行为方法集合，这就是多态的本质。

下面是上面几个代码片段的完整代码文件：

```Go
package main

import (
	"fmt"
)

type I interface {
	f()
}

type T string

func (t T) f() {
	fmt.Println("T Method")
}

type Stringer interface {
	String() string
}

func main() {

	// 类型断言
	var varI I
	varI = T("Tstring")
	if v, ok := varI.(T); ok { // 类型断言
		fmt.Println("varI类型断言结果为：", v) // varI已经转为T类型
		varI.f()
	}

	// Type-switch做类型判断
	var value interface{} // 默认为零值

	switch str := value.(type) {
	case string:
		fmt.Println("value类型断言结果为string:", str)

	case Stringer:
		fmt.Println("value类型断言结果为Stringer:", str)

	default:
		fmt.Println("value类型不在上述类型之中")
	}

	// Comma-ok断言
	value = "类型断言检查"
	str, ok := value.(string)
	if ok {
		fmt.Printf("value类型断言结果为：%T\n", str) // str已经转为string类型
	} else {
		fmt.Printf("value不是string类型 \n")
	}
}

程序输出：
varI类型断言结果为： Tstring
T Method
value类型不在上述类型之中
value类型断言结果为：string
```

使用接口使代码更具有普适性，例如函数的参数为接口变量。标准库中遵循了这个原则，但如果对接口概念没有良好的把握，是不能很好理解它是如何构建的。

那么为什么在Go语言中我们可以进行类型断言呢？我们可以在上面代码中看到，断言后的值 v, ok := varI.(T)，v值对应的是一个类型名：Tstring 。 因为在Go语言中，一个接口值(Interface Value)其实是由两部分组成：type :value 。所以在做类型断言时，变量只能是接口类型变量，断言得到的值其实是接口值中对应的类型名。这在后面讨论reflect反射包时将会有更深入的说明。

## 19.4 接口与动态类型

在经典的面向对象语言（像 C++，Java 和 C#）中，往往将数据和方法被封装为类的概念：类中包含它们两者，并且不能剥离。

Go 语言中没有类，数据（结构体或更一般的类型）和方法是一种松耦合的正交关系。Go 语言中的接口必须提供一个指定方法集的实现，但是更加灵活通用：任何提供了接口方法实现代码的类型都隐式地实现了该接口，而不用显式地声明。该特性允许我们在不改变已有的代码的情况下定义和使用新接口。

接收一个（或多个）接口类型作为参数的函数，其实参可以是任何实现了该接口的类型。 实现了某个接口的类型可以被传给任何以此接口为参数的函数 。

Go 语言动态类型的实现通常需要编译器静态检查的支持：当变量被赋值给一个接口类型的变量时，编译器会检查其是否实现了该接口的所有方法。我们也可以通过类型断言来检查接口变量是否实现了相应类型。

因此 Go 语言提供了动态语言的优点，却没有其他动态语言在运行时可能发生错误的缺点。Go 语言的接口提高了代码的分离度，改善了代码的复用性，使得代码开发过程中的设计模式更容易实现。

## 19.5 接口的提取

接口的提取，是非常有用的设计模式，良好的提取可以减少需要的类型和方法数量。而且在Go语言中不需要像传统的基于类的面向对象语言那样维护整个的类层次结构。

假设有一些拥有共同行为的对象，并且开发者想要抽象出这些行为，这时就可以创建一个接口来使用。在Go语言中这样操作甚至不会影响到前面开发的代码，所以我们不用提前设计出所有的接口，接口的设计可以不断演进，并且不用废弃之前的决定。而且类型要实现某个接口，类型本身不用改变，只需要在这个类型上实现新的接口方法集。

## 19.6 接口的继承

当一个类型包含（内嵌）另一个类型（实现了一个或多个接口）时，这个类型就可以使用（另一个类型）所有的接口方法。

类型可以通过继承多个接口来提供像多重继承一样的特性：

```Go
type ReaderWriter struct {
    io.Reader
    io.Writer
}
```



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第十八章 Struct 结构体](https://github.com/ffhelicopter/Go42/blob/master/content/42_18_struct.md)

[第二十章 方法](https://github.com/ffhelicopter/Go42/blob/master/content/42_20_method.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
