第八章 Go项目开发与编译
8.1 项目结构
Go的工程项目管理非常简单，使用目录结构和package名来确定工程结构和构建顺序。

环境变量GOPATH在项目管理中非常重要，想要构建一个项目，必须确保项目目录在GOPATH中。而GOPATH可以有多个项目用";"分隔。

Go 项目目录下一般有三个子目录：

src存放源代码
pkg编译后生成的文件
bin编译后生成的可执行文件

我们重点要关注的其实就是src文件夹中的目录结构。

为了进行一个项目，我们会在$GoPATH目录下的src目录中，新建立一个项目的主要目录，比如我写的一个WEB项目《使用gin快速搭建WEB站点以及提供RestFull接口》。
https://github.com/ffhelicopter/tmm
项目主要目录“tmm”： $GoPATH/src/github.com/ffhelicopter/tmm
在这个目录(tmm)下面还有其他目录，分别放置了其他代码，大概结构如下：

src/github.com/ffhelicopter/tmm  
                                     /api  
                                     /handler
                                     /model
                                     /task
                                     /website
                                     main.go

main.go 文件中定义了package main 。同时也在文件中import了

"github.com/ffhelicopter/tmm/api"
"github.com/ffhelicopter/tmm/handler"

2个自定义包。

上面的目录结构是一般项目的目录结构，基本上可以满足单个项目开发的需要。如果需要构建多个项目，可按照类似的结构，分别建立不同项目目录。

当我们运行go install main.go 会在GOPATH的bin 目录中生成可执行文件。
8.2 使用 Godoc
在程序中我们一般都会注释，如果我们按照一定规则，Godoc 工具会收集这些注释并产生一个技术文档。

Godoc工具在显示自定义包中的注释也有很好的效果：
注释必须以 // 开始并无空行放在声明（包，类型，函数）前。Godoc 会为每个文件生成一系列的网页。

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

一般采用行注释，每个// 后面采用一个空格隔开。注意函数或方法的整体注释以函数名或方法名开始，空换行也用//注释。整个注释和函数签名之间不要空行。

执行：

在 doc_examples 目录下我们有 Go 文件，文件中有一些注释（文件需要未编译） 命令行下进入目录下并输入命令： Godoc -http=:6060 -goroot="."

（. 是指当前目录，-goroot 参数可以是 /path/to/my/package1 这样的形式指出 package1 在你源码中的位置或接受用冒号形式分隔的路径，无根目录的路径为相对于当前目录的相对路径）

在浏览器打开地址：http://localhost:6060

然后你会看到本地的 Godoc 页面，从左到右一次显示出目录中的包。

8.3 Go程序的编译
如果想要构建一个程序，则包和包内的文件都必须以正确的顺序进行编译。包的依赖关系决定了其构建顺序。

属于同一个包的源文件必须全部被一起编译，一个包即是编译时的一个单元，因此根据惯例，每个目录都只包含一个包。

如果对一个包进行更改或重新编译，所有引用了这个包的客户端程序都必须全部重新编译。
Go 中的包模型采用了显式依赖关系的机制来达到快速编译的目的，编译器会从后缀名为 .go 的对象文件（需要且只需要这个文件）中提取传递依赖类型的信息。

如果 A.go 依赖 B.go，而 B.go 又依赖 C.go：

    编译 C.go，B.go，然后是 A.go。
为了编译 A.go，编译器读取的是 B.go 而不是 C.go。

这种机制对于编译大型的项目时可以显著地提升编译速度，每一段代码只会被编译一次。