# 《Go语言四十二章经》第一章 Go安装与运行

作者：李骁

Go语言是一门全新的静态类型开发语言，具有自动垃圾回收，丰富的内置类型, 函数多返回值，错误处理，匿名函数, 并发编程，反射，defer等关键特征，并具有简洁、安全、并行、开源等特性。从语言层面支持并发，可以充分的利用CPU多核，Go语言编译的程序可以媲美C或C++代码的速度，而且更加安全、支持并行进程。系统标准库功能完备，尤其是强大的网络库让建立Web服务成为再简单不过的事情。简单易学，内置runtime，支持继承、对象等，开发工具丰富，例如gofmt工具，自动格式化代码，让团队代码风格完美统一。同时Go非常适合用来进行服务器编程，网络编程，包括Web应用、API应用，分布式编程等等。

“Go让我体验到了从未有过的开发效率。”谷歌资深工程师罗布·派克(Rob Pike)如是说道，和C++或C一样，Go是一种系统语言，他表示，“使用它可以进行快速开发，同时它还是一个真正的编译语言，我们之所以现在将其开源，原因是我们认为它已经非常有用和强大。”

Go语言自2009年面世以来，已经有越来越多的公司开始转向Go语言开发，比如腾讯、百度、阿里、京东、小米以及360，而七牛云其技术栈基本上完全采用Go语言来开发。还有像今日头条、UBER这样的公司，他们也使用Go语言对自己的业务进行了彻底的重构。在全球范围内Go语言的使用不断增长，尤其是在云计算领域，用Go语言编写的几个主要云基础项目如Docker和Kubernetes，都取得了巨大成功。除此之外，还有各种有名的项目如etcd/consul/flannel等等，均使用Go语言实现。

Go语言有两快，一是编译运行快，还有一个是学习上手快。Go语言的学习曲线并不陡峭，无论是刚开始接触编程的朋友，还是有其他语言开发经验而打算学习Go语言的朋友，大家都可以放心大胆来学习和了解Go语言，“它值得拥有！”

让我们开始Go语言学习之旅吧！

## 1.1 Go安装

要用Go语言来进行开发，需要先搭建开发环境。Go 语言支持以下系统：

* Linux
* FreeBSD
* Mac OS
* Windows

首先需要下载Go语言安装包，Go语言的安装包下载地址为：https://golang.org/dl/ ， 国内可以正常下载地址：https://golang.google.cn/dl/

**源码编译安装**

Go语言是谷歌2009发布的第二款开源编程语言。经过几年的版本更迭，目前Go已经发布1.11版本，UNIX/Linux/Mac OS X，和 FreeBSD系统下可使用如下源码安装方法：

（1）下载源码包：https://golang.google.cn/dl/go1.11.1.linux-amd64.tar.gz
（2）将下载的源码包解压至 /usr/local目录：
tar -C /usr/local -xzf go1.11.1.linux-amd64.tar.gz
（3）将 /usr/local/go/bin 目录添加至PATH环境变量：
export PATH=$PATH:/usr/local/go/bin
（4）设置GOPATH，GOROOT环境变量：
GOPATH是工作目录，GOROOT为Go的安装目录，这里为/usr/local/go/

>注意：MAC系统下你可以使用 .pkg 结尾的安装包直接双击来完成安装，安装目录在 /usr/local/go/ 下。

**Windows系统下安装**

我们在Windows系统下一般采用直接安装，下载go1.11.1.windows-amd64.zip版本，直接解压到安装目录D:\Go，然后把D:\Go\bin目录添加到 PATH 环境变量中。

另外，还需要设置2个重要环境变量：

GOPATH=D:\goproject
GOROOT=D:\Go\

以上三个环境变量设置好后，我们就可以开始正式使用Go语言来开发了。

Windows系统也可以选择go1.11.1.windows-amd64.msi，双击运行程序根据提示来操作。


>
>GOPATH是我们的工作目录，可以有多个，用分号隔开。
>GOROOT为Go的安装目录。
>

Win+R打开CMD（注意：设置环境变量后需要重新打开CMD），输入 go ，如下显示说明Go语言运行环境已经安装成功：

```Go
D:\goproject\src>go
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report
        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         download and install packages and dependencies
        install     compile and install packages and dependencies
        list        list packages or modules
        mod         module maintenance
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildmode   build modes
        c           calling between Go and C
        cache       build and test caching
        environment environment variables
        filetype    file types
        go.mod      the go.mod file
        gopath      GOPATH environment variable
        gopath-get  legacy GOPATH go get
        goproxy     module proxy protocol
        importpath  import path syntax
        modules     modules, module versions, and more
        module-get  module-aware go get
        packages    package lists and patterns
        testflag    testing flags
        testfunc    testing functions

Use "go help <topic>" for more information about that topic.
```

另外，我们输入go version，可看到我们安装的Go版本，如图所示：

![gotool.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/gv.png)

**在本书中，所有代码编译运行和标准库的说明讲解都基于go1.11，还没有升级的用户请及时升级。**

GOPATH允许多个目录，当有多个目录时，请注意分隔符，多个目录的时候Windows是分号;

当有多个GOPATH时默认将go get获取的包存放在第一个目录下。

GOPATH目录约定有三个子目录

* src存放源代码(比如：.go .c .h .s等)   按照Go 默认约定，go run，go install等命令的当前工作路径（即在此路径下执行上述命令）。
* pkg编译时生成的中间文件（比如：.a）
* bin编译后生成的可执行文件，接下来就可以试试代码编译运行了。

文件名: test.go，代码如下：

```Go
package main

import "fmt"

func main() {
   fmt.Println("Hello, World!")
}
```

使用go命令执行以上代码输出结果如下：

D:\goproject>go run test.go

Hello，World!


## 1.2 Go语言开发工具

LiteIDE是一款开源、跨平台的轻量级 Go 语言集成开发环境（IDE）。在安装LiteIDE之前一定要先安装Go语言环境。LiteIDE支持以下的操作系统：
Windows x86 (32-bit or 64-bit)
Linux x86 (32-bit or 64-bit)

LiteIDE可以通过以下途径下载：

下载地址：https://sourceforge.net/projects/liteide/files/ 

源码地址：https://github.com/visualfc/liteide

golang中国：https://www.golangtc.com/download/liteide

也提供下载，国内下载速度可能会快一些，但版本更新较慢，建议还是选择官方地址下载。


Windows直接安装：

Windows下选择 liteidex35.1.windows-qt5.9.5.zip，下载之后解压，在liteide\bin文件夹下找到liteide.exe，双击运行。

如果不出意外，将会出现LiteIDE的运行界面。

有关LiteIDE 的使用相对来说比较简单，很容易上手，就不在此细说了。

源码编译安装：

LiteIDE源码位于https://github.com/visualfc/liteide上。需要使用Qt4/Qt5来编译源代码，Qt库可以从https://qt-project.org/downloads上获取。Mac OS X用户可以不从源代码编译Qt，直接在终端中运行brew update && brew install qt，节省大量时间。

有关LiteIDE 安装的更多说明请访问： http://liteide.org/cn/doc/install/

其他的开发工具还有Eclipse以及其集成goeclipse开发插件，以及Sublime text等，可以根据个人喜好情况选择使用。

现在Go 语言和开发工具我们都已经安装完成，接下来我们开始学习Go的基础知识，并实际使用他们来进行练习和开发。


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[前言](https://github.com/ffhelicopter/Go42/blob/master/README.md)

[第一章 Go安装与运行](https://github.com/ffhelicopter/Go42/blob/master/content/42_01_install.md)

[第二章 数据类型](https://github.com/ffhelicopter/Go42/blob/master/content/42_02_datatype.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com

