# 《Go语言四十二章经》第七章 代码结构化

作者：李骁

## 7.1 包的概念

Go语言使用包（package）的概念来组织管理代码，包是结构化代码的一种方式。和其他语言如JAVA类似，Go语言中包的主要作用是把功能相似或相关的代码组织在同一个包中，以方便查找和使用。在Go语言中，每个.go文件都必须归属于某一个包，每个文件都可有init()函数。包名在源文件中第一行通过关键字package指定，包名要小写。如下所示：

```Go
package fmt
```

每个目录下面可以有多个.go文件，这些文件只能属于同一个包，否则编译时会报错。同一个包下的不同.go文件相互之间可以直接引用变量和函数，所以这些文件中定义的全局变量和函数不能重名。

Go语言的可执行应用程序必须有main包，而且在main包中必须且只能有一个main()函数，main()函数是应用程序运行开始入口。在main包中也可以使用init()函数。

Go语言不强制要求包的名称和文件所在目录名称相同，但是这两者最好保持相同，否则很容易引起歧义。因为导入包的时候，会使用目录名作为包的路径，而在代码中使用时，却要使用包的名称。


## 7.2 包的导入

一个Go程序通过import关键字将一组包链接在一起。import其实是导入目录，而不是定义的包名称，实际应用中我们一般都会保持一致。

例如标准包中定义的big包：package big，import  "math/big" ，源代码其实是在GOROOT下src中的src/math/big目录。在代码中使用big.Int时，big指的才是.go文件中定义的包名称。

当导入多个包时，一般按照字母顺序排列包名称，像LiteIDE会在保存文件时自动完成这个动作。所谓导入包即等同于包含了这个包的所有的代码对象。

为避免名称冲突，同一包中所有对象的标识符必须要求唯一。但是相同的标识符可以在不同的包中使用，因为可以使用包名来区分它们。

import语句一般放在包名定义的下一行，导入包示例如下：

```Go
package main

import  "context"  //加载context包
```

导入多个包的常见的方式是：

```Go
import  (
"fmt"
"net/http"
 )
```

调用导入的包函数的一般方式：

```Go
fmt.Println("Hello World!")
```

下面介绍三种特殊的import方式。

点操作的含义是某个包导入之后，在调用这个包的函数时，可以省略前缀的包名，如这里可以写成Println("Hello World!")，而不是fmt.Println("Hello World!")。例如：
```Go
import( . "fmt" ) 

```

别名操作就是可以把包命名成另一个容易记忆的名字。例如：
```Go
import(
    f "fmt"
)
```
别名操作调用包函数时，前缀变成了别名，即f.Println("Hello World!")。在实际项目中有时这样使用，但请谨慎使用，不要不加节制地采用这种形式。


\_ 操作是引入某个包，但不直接使用包里的函数，而是调用该包里面的init函数，比如下面的mysql包的导入。此外在开发中，由于某种原因某个原来导入的包现在不再使用，也可以采用这种方式处理，比如下面fmt的包。代码示例如下：
```Go
import (
	_ "fmt"
	_ "github.com/go-sql-driver/mysql"
)
```

## 7.3 标准库

在 Go 的安装文件里包含了一些可以直接使用的标准库。在GOROOT/src中可以看到源码，也可以根据情况自行重新编译。

完整列表可以访问GoWalker（https://gowalker.org/）查看。

```Go
    unsafe: 包含了一些打破 Go 语言“类型安全”的命令，一般的程序中不会被使用，可用在 C/C++ 程序的调用中。
    syscall-os-os/exec:
    	os: 提供给我们一个平台无关性的操作系统功能接口，采用类UNIX设计，隐藏了不同操作系统间差异，让不同的文件系统和操作系统对象表现一致。
    	os/exec: 提供我们运行外部操作系统命令和程序的方式。
    	syscall: 底层的外部包，提供了操作系统底层调用的基本接口。
    archive/tar 和 /zip-compress：压缩(解压缩)文件功能。
    fmt-io-bufio-path/filepath-flag:
    	fmt: 提供了格式化输入输出功能。
    	io: 提供了基本输入输出功能，大多数是围绕系统功能的封装。
    	bufio: 缓冲输入输出功能的封装。
    	path/filepath: 用来操作在当前系统中的目标文件名路径。
    	flag: 对命令行参数的操作。　　
    strings-strconv-unicode-regexp-bytes:
    	strings: 提供对字符串的操作。
    	strconv: 提供将字符串转换为基础类型的功能。
    	unicode: 为 unicode 型的字符串提供特殊的功能。
    	regexp: 正则表达式功能。
    	bytes: 提供对字符型分片的操作。
    math-math/cmath-math/big-math/rand-sort:
    	math: 基本的数学函数。
    	math/cmath: 对复数的操作。
    	math/rand: 伪随机数生成。
    	sort: 为数组排序和自定义集合。
    	math/big: 大数的实现和计算。 　　
    container-/list-ring-heap: 实现对集合的操作。
    	list: 双链表。
    	ring: 环形链表。
   time-log:
        time: 日期和时间的基本操作。
        log: 记录程序运行时产生的日志。
    encoding/Json-encoding/xml-text/template:
        encoding/Json: 读取并解码和写入并编码 Json 数据。
        encoding/xml:简单的 XML1.0 解析器。
        text/template:生成像 HTML 一样的数据与文本混合的数据驱动模板。
    net-net/http-html:
        net: 网络数据的基本操作。
        http: 提供了一个可扩展的 HTTP 服务器和客户端，解析 HTTP 请求和回复。
        html: HTML5 解析器。
    runtime: Go 程序运行时的交互操作，例如垃圾回收和协程创建。
    reflect: 实现通过程序运行时反射，让程序操作任意类型的变量。
```

## 7.4 从 GitHub 安装包
如果有人想安装您的远端项目到本地机器，打开终端并执行（ffhelicopter是我在GitHub上的用户名）：

```Go
go get -u github.com/ffhelicopter/tmm
```
这样现在这台机器上的其他 Go 应用程序也可以通过导入路径："github.com/ffhelicopter/tmm" 来使用。 开发中一般这样操作：

```Go
import "github.com/ffhelicopter/tmm"
```
Go 对包的版本管理做的不是很友好，不过现在有些第三方项目做的不错，有兴趣的同学可以了解下（glide、godep、govendor）。

## 7.5 导入外部安装包
如果你要在你的应用中使用一个或多个外部包，你可以使用go install在你的本地机器上安装它们。go install 是Go语言中自动包安装工具：如需要将包安装到本地它会从远端仓库下载包：检出、编译和安装一气呵成。

在包安装前的先决条件是要自动处理包自身依赖关系的安装。被依赖的包也会安装到子目录下，但是没有文档和示例：可以到网上浏览。

**go install 使用了 GOPATH 变量**

假设你想使用https://github.com/gocolly/colly 这种托管在 Google Code、GitHub 和 Launchpad 等代码网站上的包。

你可以通过如下命令安装： go install github.com/gocolly/colly 将一个名为 github.com/gocolly/colly   安装在GOPATH/pkg/ 目录下。

go install/build都是用来编译包和其依赖的包。

区别： go build只对main包有效，在当前目录编译生成一个可执行的二进制文件（依赖包生成的静态库文件放在GOPATH/pkg）。

go install一般生成静态库文件放在GOPATH/pkg目录下，文件扩展名a。

>如果为main包，运行Go build则会在GOPATH/bin 生成一个可执行的二进制文件。

## 7.6 包的初始化


可执行应用程序的初始化和执行都起始于main包。如果main包的源代码中没有包含main()函数，则会引发构建错误 undefined: main.main。main()函数既没有参数，也没有返回类型，init()函数和main()函数在这一点上两者一样。

如果main包还导入了其它的包，那么就会在编译时将它们依次导入。有时某个包会被多个包同时导入，那么它只会被导入一次（例如很多包可能都会用到fmt包，但它只会被导入一次，因为没有必要导入多次）。

当某个包被导入时，如果该包还导入了其它的包，那么会先将其它包导入进来，然后再对这些包中的包级常量和变量进行初始化，接着执行init()函数（如果有的话），依次类推。

等所有被导入的包都加载完毕了，就会开始对main包中的包级常量和变量进行初始化，然后执行main包中的init()函数，最后执行main()函数。

Go语言中init()函数常用于包的初始化，该函数是Go语言的一个重要特性，有下面的特征：

* init函数是用于程序执行前做包的初始化的函数，比如初始化包里的变量等
* 每个包可以拥有多个init函数
* 包的每个源文件也可以拥有多个init函数
* 同一个包中多个init()函数的执行顺序不定
* 不同包的init()函数按照包导入的依赖关系决定该函数的执行顺序
* init()函数不能被其他函数调用，其在main函数执行之前，自动被调用




[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第六章 约定和惯例](https://github.com/ffhelicopter/Go42/blob/master/content/42_06_convention.md)

[第八章 Go项目开发与编译](https://github.com/ffhelicopter/Go42/blob/master/content/42_08_project.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com