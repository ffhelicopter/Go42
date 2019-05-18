《Go语言四十二章经》第五章 作用域

作者：李骁

## 5.1 作用域
* 局部变量
在函数体内或代码块内声明的变量称之为局部变量，它们的作用域只在代码块内，参数和返回值变量也是局部变量。

* 全局变量
作用域都是全局的（在本包范围内）
在函数体外声明的变量称之为全局变量，全局变量可以在整个包甚至外部包（被导出后）使用。
全局变量可以在任何函数中使用。

* 简式变量
使用 := 定义的变量，如果新变量Ga与那个同名已定义变量 (这里就是那个全局变量Ga)不在一个作用域中时，那么Go 语言会新定义这个变量Ga，遮盖住全局变量Ga。刚开始很容易在此犯错而茫然，解决方法是局部变量尽量不同名。

根据 Go语言的规范 ，Go的标识符作用域是基于代码块（code block）的。代码块就是包裹在一对大括号内部的声明和语句，并且是可嵌套的。在代码中直观可见的显式的(explicit)code block，比如：函数的函数体、for循环的循环体等；还有隐式的(implicit)code block。

我们使用最多的if语句类型就是 单if型 ，即:

```Go
if simplestmt; expression {
    ... ...
}
```

在这种类型的if语句中，有两个code block：一个隐式的code block和一个显式的code block。我们把上面的形式代码做一个等价变化，并加上code block起始和结束点的标注，结果如下：

```Go
{ // 隐式code block
    simplestmt
    if expression { // 显式的code block
            ... ...
    } 
} 
```

下面的代码综合了几种作用域的情况，很容易混淆。请各位仔细琢磨弄清楚。

```Go
package main

import (
	"fmt"
)

var (
	Ga int = 99
)

const (
	v int = 199
)

func GetGa() func() int {

	if Ga := 55; Ga < 60 {
		fmt.Println("GetGa if 中：", Ga)
	}

	for Ga := 2; ; {
		fmt.Println("GetGa循环中：", Ga)
		break
	}

	fmt.Println("GetGa函数中：", Ga)

	return func() int {
		Ga += 1
		return Ga
	}
}

func main() {
	Ga := "string"
	fmt.Println("main函数中：", Ga)

	b := GetGa()
	fmt.Println("main函数中：", b(), b(), b(), b())

	v := 1
	{
		v := 2
		fmt.Println(v)
		{
			v := 3
			fmt.Println(v)
		}
	}
	fmt.Println(v)
}


```

```Go
程序输出：

main函数中： string
GetGa if 中： 55
GetGa循环中： 2
GetGa函数中： 99
main函数中： 100 101 102 103
2
3
1

```

Ga作为全局变量纯在是int类型，值为99；而在main()中时，Ga通过简式声明 := 操作，是string类型，值为string。在main()中，v很典型地体现了在“{}”花括号中的作用域问题，每一层花括号，都是对上一层的屏蔽。而闭包函数，GetGa()返回的匿名函数，赋值给b，每次执行b()，Ga的值都被记忆在内存中，下次执行b()的时候，取b()上次执行后Ga的值，而不是全局变量Ga的值，这就是闭包函数可以使用包含它的函数内的变量，因为作为代码块一直存在，所以每次执行都是在上次基础上运行。

简单总结如下：

有花括号"{ }"一般都存在作用域的划分；
:= 简式声明会屏蔽所有上层代码块中的变量（常量），建议使用规则来规范，如对常量使用全部大写，而变量尽量小写；
在if等语句中存在隐式代码块，需要注意；
闭包函数可以理解为一个代码块，并且他可使用包含它的函数内的变量；

>注意，简式变量只能在函数内部声明使用，但是它可能会覆盖函数外全局同名变量。而且你不能在一个单独的声明中重复声明一个变量，但在多变量声明中这是允许的，而且其中至少要有一个新的声明变量。重复变量需要在相同的代码块内，否则你将得到一个隐藏变量。
>
>如果你在代码块中犯了这个错误，将不会出现编译错误，但应用运行结果可能不是你所期望。所以尽可能避免和全局变量同名。

思考：

```Go
func main() {
    if a := 1; false {
    } else if b := 2; false {
    } else if c := 3; false {
    } else {
        println(a, b, c)
    }
}

```

这段代码运行结果是什么，你能写出来吗？




[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第四章 常量](https://github.com/ffhelicopter/Go42/blob/master/content/42_04_const.md)

[第六章 约定和惯例](https://github.com/ffhelicopter/Go42/blob/master/content/42_06_convention.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42 
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com 