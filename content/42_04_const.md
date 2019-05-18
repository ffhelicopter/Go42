# 《Go语言四十二章经》第四章 常量

作者：李骁

## 4.1 常量以及iota

常量使用关键字 const 定义，用于存储不会改变的数据。常量不能被重新赋予任何值。 
存储在常量中的数据类型只可以是布尔型、数字型（整数型、浮点型和复数）和字符串型。
常量的定义格式：const identifier [type] = value，例如：

```Go
const Pi = 3.14159
```
在 Go 语言中，你可以省略类型说明符 [type]，因为编译器可以根据变量（常量）的值来推断其类型。


    显式类型定义： const b string = "abc"
    隐式类型定义： const b = "abc"

Go的常量定义可以限定常量类型，但不是必需的。如果定义常量时没有指定类型，那么它与字面常量一样，是无类型（untyped）常量。一个没有指定类型的常量被使用时，会根据其使用环境而推断出它所需要具备的类型。换句话说，未定义类型的常量会在必要时刻根据上下文来获得相关类型。

字面常量（literal），是指程序中硬编码的常量，如：-12。它们的值即为它们本身，是无法被改变的。 

常量的值必须是能够在编译时就能够确定的；你可以在其赋值表达式中涉及计算过程，但是所有用于计算的值必须在编译期间就能获得。

Go语言预定义了这些常量： true、 false和iota。布尔常量只包含两个值：true 和 false。iota比较特殊，可以被认为是一个可被编译器修改的常量，在每一个const关键字出现时被重置为0，然后在下一个const出现之前，每出现一次iota，其所代表的数字会自动增1。

在这个例子中，iota 可以被用作枚举值：

```Go
const (
    a = iota
    b = iota
    c = iota
)
```

第一个 iota 等于 0，每当 iota 在新的一行被使用时，它的值都会自动加 1；所以 a=0, b=1, c=2 可以简写为如下形式：

```Go
const (
    a = iota
    b
    c
)
```
注意：

```Go
const (
    a = iota
    b = 8
    c
)
```
a, b, c分别为0, 8, 8，新的常量b声明后，iota 不再向下赋值，后面常量如果没有赋值，则继承上一个常量值。

可以简单理解为在一个const块中，每换一行定义个常量，iota 都会自动+1。

（ 关于 iota 的使用涉及到非常复杂多样的情况 ，这里不展开来讲了，有兴趣可以查查资料研究）

iota 也可以用在表达式中，如：iota + 50。在每遇到一个新的常量块或单个常量声明时， iota 都会重置为 0（ **简单地讲，每遇到一次 const 关键字，iota 就重置为 0 ** ）。

使用位左移与 iota 计数配合可优雅地实现存储单位的常量枚举：

```Go
type ByteSize float64
const (
    _ = iota // 通过赋值给空白标识符来忽略值
    KB ByteSize = 1<<(10*iota)
    MB
    GB
    TB
    PB
    EB
    ZB
    YB
)
```

数值常量（Numeric constants）包括整数，浮点数以及复数常量。数值常量有一些微妙之处。

```Go
package main

import (
	"fmt"
)

func main() {
	const a = 5
	var intVar int = a
	var int32Var int32 = a
	var float64Var float64 = a
	var complex64Var complex64 = a
	fmt.Println("intVar", intVar, "\nint32Var", int32Var, "\nfloat64Var", float64Var, "\ncomplex64Var", complex64Var)
}
```

```Go
程序输出
intVar 5 
int32Var 5 
float64Var 5 
complex64Var (5+0i)

```

在这个程序中，a 的值是 5 并且 a 在语法上是泛化的（它既可以表示浮点数 5.0，也可以表示整数 5，甚至可以表示没有虚部的复数 5 + 0i），因此 a 可以赋值给任何与之类型兼容的变量。像 a 这种数值常量的默认类型可以想象成是通过上下文动态生成的。


当然，常量之所以为常量就是恒定不变的量，因此我们无法在程序运行过程中修改它的值；如果你在代码中试图修改常量的值则会引发编译错误。同时，在const 定义中，对常量名没有强制要求全部大写，不过我们一般都会全部字母大写，以便阅读。


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第三章 变量](https://github.com/ffhelicopter/Go42/blob/master/content/42_03_var.md)

[第五章 作用域](https://github.com/ffhelicopter/Go42/blob/master/content/42_05_scope.md)




>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com