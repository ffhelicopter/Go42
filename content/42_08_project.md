# 《Go语言四十二章经》第八章 Go项目开发与编译

作者：李骁

## 8.1 项目结构
Go的工程项目管理非常简单，使用目录结构和包名来确定工程结构和构建顺序。

环境变量GOPATH在项目管理中非常重要，想要构建一个项目，必须确保项目目录在GOPATH中。而GOPATH可以有多个项目用";"分隔。

Go 项目目录下一般有三个子目录：

* src存放源代码
* pkg编译后生成的文件
* bin编译后生成的可执行文件

我们重点要关注的其实就是src文件夹中的目录结构。

为了进行一个项目，我们会在GOPATH目录下的src目录中，新建立一个项目的主要目录，比如我写的一个WEB项目《使用gin快速搭建WEB站点以及提供RESTful接口》。
https://github.com/ffhelicopter/tmm
项目主要目录“tmm”： GOPATH/src/github.com/ffhelicopter/tmm
在这个目录(tmm)下面还有其他目录，分别放置了其他代码，大概结构如下：

```Go
src/github.com/ffhelicopter/tmm  
                               /api  
                               /handler
                               /model
                               /task
                               /website
                               main.go
```
main.go 文件中定义了package main 。同时也在文件中import了

```Go
"github.com/ffhelicopter/tmm/api"
"github.com/ffhelicopter/tmm/handler"
```
2个自定义包。

上面的目录结构是一般项目的目录结构，基本上可以满足单个项目开发的需要。如果需要构建多个项目，可按照类似的结构，分别建立不同项目目录。

当我们运行go install main.go 会在GOPATH的bin 目录中生成可执行文件。

## 8.2 使用godoc

在程序中我们一般都会注释，如果我们按照一定规则，godoc工具会收集这些注释并产生一个技术文档。

```Go
// Copyright 2009 The Go Authors. All rights reserved.  
// Use of this source code is governed by a BSD-style  
// license that can be found in the LICENSE file.     

package zlib
....

// A Writer takes data written to it and writes the compressed
// form of that data to an underlying writer (see NewWriter).
type Writer struct {
    w           io.Writer
    level       int
    dict        []byte
    compressor  * flate.Writer
    digest      hash.Hash32
    err         error
    scratch     [4]byte
    wroteHeader bool
}

// NewWriter creates a new Writer.
// Writes to the returned Writer are compressed and written to w.
//
// It is the caller's responsibility to call Close on the WriteCloser when done.
// Writes may be buffered and not flushed until Close.
func NewWriter(w io.Writer) * Writer {
    z, _ := NewWriterLevelDict(w, DefaultCompression, nil)
    return z
}
```

命令行下进入目录下并输入命令： godoc -http=:6060 -goroot="."

然后在浏览器打开地址：http://localhost:6060

然后你会看到本地的 Godoc 页面，从左到右一次显示出目录中的包。
![godoc.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/godoc.png)


## 8.3 Go程序的编译

在Go语言中，和编译有关的命令主要是go run ,go build , go install这三个命令。

go run只能作用于main包文件，先运行compile 命令编译生成.a文件，然后 link 生成最终可执行文件并运行程序，这过程的产生的是临时文件，在go run 退出前会删除这些临时文件（含.a文件和可执行文件）。最后直接在命令行输出程序执行结果。go run 命令在第二次执行的时候，如果发现导入的代码包没有发生变化，那么 go run 不会再次编译这个导入的代码包，直接进行链接生成最终可执行文件并运行程序。

go install用于编译并安装指定的代码包及它们的依赖包，并且将编译后生成的可执行文件放到 bin 目录下（GOPATH/bin），编译后的包文件放到当前工作区的 pkg 的平台相关目录下。

go build用于编译指定的代码包以及它们的依赖包。如果用来编译非main包的源码，则只做检查性的编译，而不会输出任何结果文件。如果是一个可执行程序的源码（即是 main 包），这个过程与go run 大体相同，除了会在当前目录生成一个可执行文件外。

使用go build时有一个地方需要注意，对外发布编译文件如果不希望被人看到源代码，请使用go build -ldflags 命令，设置编译参数-ldflags "-w -s" 再编译后发布。避免使用gdb来调试而清楚看到源代码。

![ch5.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/ch5.png)

## 8.4 Go modules 包依赖管理

Go 1.11 新增了对模块的支持，希望借此解决“包依赖管理”。可以通过设置环境变量 GO111MODULE来开启或关闭模块支持，它有三个可选值： off、 on、 auto，默认值是 auto。

* GO111MODULE=off
    无模块支持，go 会从 GOPATH 和 vendor 文件夹寻找包。
* GO111MODULE=on
    模块支持，go 会忽略 GOPATH 和 vendor 文件夹，只根据 go.mod下载依赖。
* GO111MODULE=auto
在 GOPATH/src外面且根目录有 go.mod文件时，开启模块支持。

在使用模块的时候， GOPATH是无意义的，不过它还是会把下载的依赖储存在 GOPATH/pkg/mod 中。
 
运行命令，go help mod ，我们可以看到mod的操作子命令，主要是init、 edit、 tidy。

```Go
Go mod provides access to operations on modules.

Note that support for modules is built into all the go commands,
not just 'go mod'. For example, day-to-day adding, removing, upgrading,
and downgrading of dependencies should be done using 'go get'.
See 'go help modules' for an overview of module functionality.

Usage:

        go mod <command> [arguments]

The commands are:

        download    download modules to local cache
        edit        edit go.mod from tools or scripts
        graph       print module requirement graph
        init        initialize new module in current directory
        tidy        add missing and remove unused modules
        vendor      make vendored copy of dependencies
        verify      verify dependencies have expected content
        why         explain why packages or modules are needed

Use "go help mod <command>" for more information about a command.
```

命令含义：
download   下载依赖的module到本地cache
edit        编辑go.mod文件
graph      打印模块依赖图
init        在当前文件夹下初始化一个新的module, 创建go.mod文件
tidy       增加丢失的module，去掉未用的module
vendor     将依赖复制到vendor下
verify      校验依赖
why       解释为什么需要依赖

为了使用modules来管理项目，我们可以以下几个步骤来操作：

（1）首先需要设置GO111MODULE ，这里我们设置为auto。

（2）考虑和原来GOPATH有所隔离，新建立了一个目录D:\gomodules来存放modules管理的项目。

（3）在D:\gomodules下建立ind项目，建立对应的目录，D:\gomodules\ind

（4）在ind目录中，我们编写了该项目的主要文件main.go


```Go

package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit("http://go-colly.org/")
}
```

1、第一次需要我们运行init命令初始化：

D:\gomodules\ind>go mod init ind

go: creating new go.mod: module ind

可以在ind目录看到新生成了一个文件：go.mod ，这个modules名字叫ind。

2、接下来我们运行go mod tidy 命令，发现如下图一样出现报错，这主要是众所周知的网络原因，由于这里主要是golang.org/x下的包，所以可以简单使用replace命令来解决这个问题，如果是其他厂商的依赖包，还是优先解决网络问题。

然后重复运行go mod tidy ，如果出错在使用replace，直到能正常运行go mod tidy 命令完成。

go mod edit -replace=old[@v]=new[@v]

注意：replace版本号可以在错误信息中看到。

D:\gomodules\ind>go mod edit -replace=golang.org/x/net@v0.0.0-20181114220301-adae6a3d119a=github.com/golang/net@v0.0.0-20181114220301-adae6a3d119a


![go mod tidy 命令](https://github.com/ffhelicopter/Go42/blob/master/content/img/tidy.png)


3、我们看到在ind目录下面多了2个文件，分别是go.mod和go.sum。

go.mod文件

```Go
module ind

replace (
	golang.org/x/net v0.0.0-20180218175443-cbe0f9307d01 => github.com/golang/net v0.0.0-20180218175443-cbe0f9307d01
	golang.org/x/net v0.0.0-20181114220301-adae6a3d119a => github.com/golang/net v0.0.0-20181114220301-adae6a3d119a
)

require (
	github.com/PuerkitoBio/goquery v1.5.0 // indirect
	github.com/antchfx/htmlquery v0.0.0-20181207070731-9784ecda34b7 // indirect
	github.com/antchfx/xmlquery v0.0.0-20181204011708-431a9e9e7c44 // indirect
	github.com/antchfx/xpath v0.0.0-20181208024549-4bbdf6db12aa // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gocolly/colly v1.1.0
	github.com/kennygrant/sanitize v1.2.4 // indirect
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca // indirect
	github.com/temoto/robotstxt v0.0.0-20180810133444-97ee4a9ee6ea // indirect
)
```

go.mod文件可以通过require，replace和exclude语句使用的精确软件包集。

（1）require语句指定的依赖项模块

（2）replace语句可以替换依赖项模块

（3）exclude语句可以忽略依赖项模块


go.sum文件

```Go
github.com/PuerkitoBio/goquery v1.5.0 h1:uGvmFXOA73IKluu/F84Xd1tt/z07GYm8X49XKHP7EJk=
github.com/PuerkitoBio/goquery v1.5.0/go.mod h1:qD2PgZ9lccMbQlc7eEOjaeRlFQON7xY8kdmcsrnKqMg=
github.com/andybalholm/cascadia v1.0.0 h1:hOCXnnZ5A+3eVDX8pvgl4kofXv2ELss0bKcqRySc45o=
github.com/andybalholm/cascadia v1.0.0/go.mod h1:GsXiBklL0woXo1j/WYWtSYYC4ouU9PqHO0sqidkEA4Y=
github.com/antchfx/htmlquery v0.0.0-20181207070731-9784ecda34b7 h1:w7OFcAjjWOJ/Fp9/dlvikG46C44FV/B8G42Tj+KlFUk=
github.com/antchfx/htmlquery v0.0.0-20181207070731-9784ecda34b7/go.mod h1:MS9yksVSQXls00iXkiMqXr0J+umL/AmxXKuP28SUJM8=
github.com/antchfx/xmlquery v0.0.0-20181204011708-431a9e9e7c44 h1:utJNS82e0x9ZhwWvitDlUv2+0HgGYfyrSKX9hDf0uW0=
github.com/antchfx/xmlquery v0.0.0-20181204011708-431a9e9e7c44/go.mod h1:/+CnyD/DzHRnv2eRxrVbieRU/FIF6N0C+7oTtyUtCKk=
github.com/antchfx/xpath v0.0.0-20181208024549-4bbdf6db12aa h1:lL66YnJWy1tHlhjSx8fXnpgmv8kQVYnI4ilbYpNB6Zs=
github.com/antchfx/xpath v0.0.0-20181208024549-4bbdf6db12aa/go.mod h1:Yee4kTMuNiPYJ7nSNorELQMr1J33uOpXDMByNYhvtNk=
github.com/gobwas/glob v0.2.3 h1:A4xDbljILXROh+kObIiy5kIaPYD8e96x1tgBhUI5J+Y=
github.com/gobwas/glob v0.2.3/go.mod h1:d3Ez4x06l9bZtSvzIay5+Yzi0fmZzPgnTbPcKjJAkT8=
github.com/gocolly/colly v1.1.0 h1:B1M8NzjFpuhagut8f2ILUDlWMag+nTx+PWEmPy7RhrE=
github.com/gocolly/colly v1.1.0/go.mod h1:Hof5T3ZswNVsOHYmba1u03W65HDWgpV5HifSuueE0EA=
github.com/golang/net v0.0.0-20180218175443-cbe0f9307d01/go.mod h1:98y8FxUyMjTdJ5eOj/8vzuiVO14/dkJ98NYhEPG8QGY=
github.com/golang/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:98y8FxUyMjTdJ5eOj/8vzuiVO14/dkJ98NYhEPG8QGY=
github.com/kennygrant/sanitize v1.2.4 h1:gN25/otpP5vAsO2djbMhF/LQX6R7+O1TB4yv8NzpJ3o=
github.com/kennygrant/sanitize v1.2.4/go.mod h1:LGsjYYtgxbetdg5owWB2mpgUL6e2nfw2eObZ0u0qvak=
github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca h1:NugYot0LIVPxTvN8n+Kvkn6TrbMyxQiuvKdEwFdR9vI=
github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca/go.mod h1:uugorj2VCxiV1x+LzaIdVa9b4S4qGAcH6cbhh4qVxOU=
github.com/temoto/robotstxt v0.0.0-20180810133444-97ee4a9ee6ea h1:hH8P1IiDpzRU6ZDbDh/RDnVuezi2oOXJpApa06M0zyI=
github.com/temoto/robotstxt v0.0.0-20180810133444-97ee4a9ee6ea/go.mod h1:aOux3gHPCftJ3KHq6Pz/AlDjYJ7Y+yKfm1gU/3B0u04=
```

打开目录GOPATH/pkg/mod，我们可以看到这个项目下的依赖包都下载过来了。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第七章 代码结构化](https://github.com/ffhelicopter/Go42/blob/master/content/42_07_package.md)

[第九章 运算符](https://github.com/ffhelicopter/Go42/blob/master/content/42_09_operator.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com