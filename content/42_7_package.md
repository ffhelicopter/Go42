# 《Go语言四十二章经》第七章 代码结构化

作者：李骁

## 7.1 包的概念
包是结构化代码的一种方式：每个程序都由包（通常简称为 pkg）的概念组成，可以使用自身的包或者从其它包中导入内容。

如同其它一些编程语言中的类库或命名空间的概念，每个 Go 文件都属于且仅属于一个包。一个包可以由许多以 .go 为扩展名的源文件组成，因此文件名和包名一般来说都是不相同的。

你必须在源文件中非注释的第一行指明这个文件属于哪个包，如：package main 。

package main表示一个可独立执行的程序，每个 Go 应用程序都包含一个名为 main 的包。package main包下可以有多个文件，但所有文件中只能有一个main()方法，main()方法代表程序入口。

一个应用程序可以包含不同的包，而且即使你只使用 main 包也不必把所有的代码都写在一个巨大的文件里：你可以用一些较小的文件，并且在每个文件非注释的第一行都使用 package main 来指明这些文件都属于 main 包。如果你打算编译包名不是为 main 的源文件，如 pack1，编译后产生的对象文件将会是 pack1.a 而不是可执行程序。另外要注意的是，所有的包名都应该使用小写字母。当然，main包是不能在其他文档import的，编译器会报错：

```Go
import "xx/xx" is a program, not an importable package。
```
简单地说，在含有mian包的目录下，你可以写多个文件，每个文件非注释的第一行都使用 package main 来指明这些文件都属于这个应用的 main 包，只有一个文件能有mian() 方法，也就是应用程序的入口。main包不是必须的，只有在可执行的应用程序中需要。

## 7.2 包的导入
一个 Go 程序是通过 import 关键字将一组包链接在一起。

import "fmt" 告诉 Go 编译器这个程序需要使用 fmt 包（的函数，或其他元素），fmt 包实现了格式化 IO（输入/输出）的函数。包名被封闭在半角双引号 "" 中。如果你打算从已编译的包中导入并加载公开声明的方法，不需要插入已编译包的源代码。

<b>import 其实是导入目录</b>，而不是定义的package名字，虽然我们一般都会保持一致，但其实是可以随便定义目录名，只是使用时会很容易混乱，不建议这么做。

例如：package big ，我们import  "math/big" ，其实是在src中的src/math目录。在代码中使用big.Int时，big指的才是Go文件中定义的package名字。

当你导入多个包时，最好按照字母顺序排列包名，这样做更加清晰易读。

如果包名不是以 ./ ，如 "fmt" 或者 "container/list"，则 Go 会在全局文件进行查找；如果包名以 ./ 开头，则 Go 会在相对目录中查找。

导入包即等同于包含了这个包的所有的代码对象。

除了符号 \_，包中所有代码对象的标识符必须是唯一的，以避免名称冲突。但是相同的标识符可以在不同的包中使用，因为可以使用包名来区分它们。

导入包的路径的几种情况：

* 第一种方式相对路径

```Go
import   "./module"   //当前文件同一目录的module目录， 此方式没什么用容易出错，不建议用
```
* 第二种方式绝对路径

```Go
import  "LearnGo/init"  //加载Gopath/src/LearnGo/init模块，一般建议这样使用""
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

下面展示一些特殊的import方式
* 点操作
import( . "fmt" )  
这个点操作的含义就是这个包导入之后在你调用这个包的函数时，你可以省略前缀的包名，如可以省略的写成Println("hello world!")

* 别名操作
别名操作顾名思义我们可以把包命名成另一个我们用起来容易记忆的名字。

```Go
import(
    f "fmt"
)
```
别名操作调用包函数时前缀变成了我们的前缀，即f.Println("hello world")。在实际项目中有这样使用，但请谨慎使用，不要广泛采用这种形式。

* \_操作
\_操作其实是引入该包，而不直接使用包里面的函数，而是调用了该包里面的init函数。

```Go
import (
	_ "github.com/revel/modules/testrunner/app"
	_ "guild_website/app"
)
```

## 7.3 标准库

在 Go 的安装文件里包含了一些可以直接使用的包，即标准库。在 Windows 下，标准库的位置在 Go 根目录下的子目录 pkg\windows_386 中；在 Linux 下，标准库在 Go 根目录下的子目录 pkg\linux_amd64 中（如果是安装的是 32 位，则在 linux_386 目录中）。Go 的标准库包含了大量的包（如：fmt 和 os）, 在$GoROOT/src中可以看到源码，也可以根据情况自行重新编译。

完整列表可以在 Go Walker 查看。

```Go
    unsafe: 包含了一些打破 Go 语言“类型安全”的命令，一般的程序中不会被使用，可用在 C/C++ 程序的调用中。
    syscall-os-os/exec:
    	os: 提供给我们一个平台无关性的操作系统功能接口，采用类Unix设计，
隐藏了不同操作系统间差异，让不同的文件系统和操作系统对象表现一致。
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
如果有人想安装您的远端项目到本地机器，打开终端并执行（ffhelicopter是我在 GitHub 上的用户名）：

```Go
go get -u github.com/ffhelicopter/tmm
```
这样现在这台机器上的其他 Go 应用程序也可以通过导入路径："github.com/ffhelicopter/tmm" 代替 "./ffhelicopter/tmm" 来使用。 也可以将其缩写为：import ind "github.com/ffhelicopter/tmm"；开发中一般这样操作：

```Go
import "github.com/ffhelicopter/tmm"
```
Go 对包的版本管理做的不是很友好，不过现在有些第三方项目做的不错，有兴趣的同学可以了解下（glide、godep、govendor）。

## 7.5 导入外部安装包
如果你要在你的应用中使用一个或多个外部包，你可以使用 Go install在你的本地机器上安装它们。Go install 是 Go 中自动包安装工具：如需要将包安装到本地它会从远端仓库下载包：检出、编译和安装一气呵成。

在包安装前的先决条件是要自动处理包自身依赖关系的安装。被依赖的包也会安装到子目录下，但是没有文档和示例：可以到网上浏览。

**Go install 使用了 GoPATH 变量**

假设你想使用https://github.com/gocolly/colly 这种托管在 Google Code、GitHub 和 Launchpad 等代码网站上的包。

你可以通过如下命令安装： Go install github.com/gocolly/colly 将一个名为 github.com/gocolly/colly   安装在$GoPATH/pkg/ 目录下。

Go install/build都是用来编译包和其依赖的包。

区别： Go build只对main包有效，在当前目录编译生成一个可执行的二进制文件（依赖包生成的静态库文件放在$GoPATH/pkg）。

Go install一般生成静态库文件放在$GoPATH/pkg目录下，文件扩展名a。

>如果为main包，运行Go buil则会在$GoPATH/bin 生成一个可执行的二进制文件。

## 7.6 包的分级声明和初始化
你可以在使用 import 导入包之后定义或声明 0 个或多个常量（const）、变量（var）和类型（type），这些对象的作用域都是全局的（在本包范围内），所以可以被本包中所有的函数调用，然后声明一个或多个函数（func）。

如果存在 init 函数的话，则对该函数进行定义（这是一个特殊的函数，每个含有该函数的包都会首先执行这个函数）。

程序开始执行并完成初始化后，第一个调用（程序的入口点）的函数是 main.main()（如果有 init() 函数则会先执行该函数）。

如果你的 main 包的源代码没有包含 main 函数，则会引发构建错误 undefined: main.main。main 函数既没有参数，也没有返回类型，这一点上 init 函数和 main 函数一样。

main函数一旦返回就表示程序已成功执行并立即退出。

Go 程序的执行（程序启动）顺序如下：
程序的初始化和执行都起始于main包。如果main包还导入了其它的包，那么就会在编译时将它们依次导入。有时一个包会被多个包同时导入，那么它只会被导入一次（例如很多包可能都会用到fmt包，但它只会被导入一次，因为没有必要导入多次）。当一个包被导入时，如果该包还导入了其它的包，那么会先将其它包导入进来，然后再对这些包中的包级常量和变量进行初始化，接着执行init函数（如果有的话），依次类推。等所有被导入的包都加载完毕了，就会开始对main包中的包级常量和变量进行初始化，然后执行main包中的init函数（如果存在的话），最后执行main函数。



Go语言中init函数用于包(package)的初始化，该函数是Go语言的一个重要特性，有下面的特征：

* init函数是用于程序执行前做包的初始化的函数，比如初始化包里的变量等
* 每个包可以拥有多个init函数
* 包的每个源文件也可以拥有多个init函数
* 同一个包中多个init函数的执行顺序Go语言没有明确的定义(说明)
* 不同包的init函数按照包导入的依赖关系决定该初始化函数的执行顺序
* init函数不能被其他函数调用，而是在main函数执行之前，自动被调用




>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>本书《Go语言四十二章经》内容在简书同步地址：  https://www.jianshu.com/nb/29056963
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com