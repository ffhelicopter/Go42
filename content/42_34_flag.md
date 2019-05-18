# 《Go语言四十二章经》第三十四章 命令行flag包 

作者：李骁

## 34.1 命令行

写命令行程序时需要对命令参数进行解析，这时我们可以使用os库。os库可以通过变量Args来获取命令参数，os.Args返回一个字符串数组，其中第一个参数就是执行文件本身。

```Go
package main
 
import (
    "fmt"
    "os"
)
 
func main() {
    fmt.Println(os.Args)
}
```

编译执行后执行

```Go
$ ./cmd -user="root"
 [./cmd -user=root]
```

这种方式对于简单的参数格式还能使用，一旦面对复杂的参数格式，比较费时费劲，所以这时我们会选择flag库。


## 34.2 flag包

Go提供了flag包，可以很方便的操作命名行参数，下面介绍下flag的用法。

几个概念：

1）命令行参数（或参数）：是指运行程序提供的参数

2）已定义命令行参数：是指程序中通过flag.Xxx等这种形式定义了的参数

3）非flag（non-flag）命令行参数（或保留的命令行参数）：先可以简单理解为flag包不能解析的参数

```Go
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	h, H bool

	v bool
	q *bool

	D    string
	Conf string
)

func init() {
	flag.BoolVar(&h, "h", false, "帮助信息")
	flag.BoolVar(&h, "H", false, "帮助信息")

	flag.BoolVar(&v, "v", false, "显示版本号")

	//
	flag.StringVar(&D, "D", "deamon", "set descripton ")
	flag.StringVar(&Conf, "Conf", "/dev/conf/cli.conf", "set Conf filename ")

	// 另一种绑定方式
	q = flag.Bool("q", false, "退出程序")

	// 像flag.Xxx函数格式都是一样的，第一个参数表示参数名称，
	// 第二个参数表示默认值，第三个参数表示使用说明和描述。
	// flag.XxxVar这样的函数第一个参数换成了变量地址，
        // 后面的参数和flag.Xxx是一样的。

	// 改变默认的 Usage

	flag.Usage = usage

	flag.Parse()

	var cmd string = flag.Arg(0)

	fmt.Printf("-----------------------\n")
	fmt.Printf("cli non=flags      : %s\n", cmd)

	fmt.Printf("q: %b\n", *q)

	fmt.Printf("descripton:  %s\n", D)
	fmt.Printf("Conf filename : %s\n", Conf)

	fmt.Printf("-----------------------\n")
	fmt.Printf("there are %d non-flag input param\n", flag.NArg())
	for i, param := range flag.Args() {
		fmt.Printf("#%d    :%s\n", i, param)
	}

}

func main() {
	flag.Parse()

	if h || H {
		flag.Usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `CLI: 8.0
Usage: Cli [-hvq] [-D descripton] [-Conf filename] 

`)
	flag.PrintDefaults()
}
```

flag包实现了命令行参数的解析，大致需要几个步骤：

一：flag参数定义或绑定

定义flags有两种方式：

1）flag.Xxx()，其中Xxx可以是Int、String等；返回一个相应类型的指针，如：

```Go
var ip = flag.Int("flagname", 1234, "help message for flagname")
```

2）flag.XxxVar()，将flag绑定到一个变量上，如：

```Go
var flagvar int
flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
```

另外，还可以创建自定义flag，只要实现flag.Value接口即可（要求receiver是指针），这时候可以通过如下方式定义该flag：

flag.Var(&flagVal, "name", "help message for flagname")

命令行flag的语法有如下三种形式：
-flag // 只支持bool类型
-flag=x
-flag x // 只支持非bool类型

二：flag参数解析

在所有的flag定义完成之后，可以通过调用flag.Parse()进行解析。

根据Parse()中for循环终止的条件，当parseOne返回false，nil时，Parse解析终止。
```Go
s := f.args[0]
if len(s) == 0 || s[0] != '-' || len(s) == 1 {
    return false, nil
}
```
当遇到单独的一个“-”或不是“-”开始时，会停止解析。比如：./cli – -f 或 ./cli -f

这两种情况，-f都不会被正确解析。像这些参数，我们称之为non-flag参数

parseOne方法中接下来是处理-flag=x，然后是-flag（bool类型）（这里对bool进行了特殊处理），接着是-flag x这种形式，最后，将解析成功的Flag实例存入FlagSet的actual map中。

Arg(i int)和Args()、NArg()、NFlag()
Arg(i int)和Args()这两个方法就是获取non-flag参数的；NArg()获得non-flag个数；NFlag()获得FlagSet中actual长度（即被设置了的参数个数）。

flag解析遇到non-flag参数就停止了。所以如果我们将non-flag参数放在最前面，flag什么也不会解析，因为flag遇到了这个就停止解析了。

三：分支程序

根据参数值，代码进入分支程序，执行相关功能。上面代码提供了 -h 参数的功能执行。
```Go
if h || H {
		flag.Usage()
	}
```
总体而言，从例子上看，flag package很有用，但是并没有强大到解析一切的程度。如果你的入参解析非常复杂，flag可能捉襟见肘。

Cobra是一个用来创建强大的现代CLI命令行的Go开源库。开源包可能比较合适构建更为复杂的命令行程序。开源地址：https://github.com/spf13/cobra


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第三十三章 Socket网络](https://github.com/ffhelicopter/Go42/blob/master/content/42_33_socket.md)

[第三十五章 模板](https://github.com/ffhelicopter/Go42/blob/master/content/42_35_template.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com