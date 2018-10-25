《Go语言四十二章经》第一章 Go安装与运行

作者：李骁

Go语言是一门全新的静态类型开发语言，具有自动垃圾回收，丰富的内置类型, 函数多返回值，错误处理，匿名函数, 并发编程，反射等特性，并具有简洁、安全、并行、开源等特性。

在接下来26章中，主要讲解Go语言的基本语法和特性。在掌握基础知识的前提下，需要对struct、interface、方法、通道以及锁和goroutine有彻底理解。

## 1.1 Go安装

Go语言支持以下系统：

* Linux
* FreeBSD
* Mac OS
* Windows

安装包下载地址为：https://golang.org/dl/
国内可以正常下载地址：https://golang.google.cn/dl/


UNIX/Linux/Mac OS X，和FreeBSD系统下使用源码安装方法：

1、下载源码包：go1.11.linux-amd64.tar.gz。
2、将下载的源码包解压至 /usr/local目录。
tar -C /usr/local -xzf go1.11.linux-amd64.tar.gz
3、将 /usr/local/go/bin 目录添加至PATH环境变量：
export PATH=$PATH:/usr/local/go/bin

>注意：MAC系统下你可以使用 .pkg 结尾的安装包直接双击来完成安装，安装目录在 /usr/local/go/ 下。

Windows系统下安装

你可以选择本地安装目录：D:\Go，把D:\Go\bin目录添加到 PATH 环境变量中。

设置$GOPATH=D:\goproject以及 $GOROOT=D:\Go\ 。

打开CMD，输入 go version，如下显示说明go运行环境已经安装成功：

![goversion.png](https://upload-images.jianshu.io/upload_images/6324013-5e10325d23d966c3.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

>在本书中，所有代码和标准库的讲解都基于go1.11，还没有升级的用户请及时升级。

$GOPATH允许多个目录，当有多个目录时，请注意分隔符，多个目录的时候Windows是分号;
当有多个$GOPATH时默认将go get获取的包存放在第一个目录下。

$GOPATH目录约定有三个子目录

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

LiteIDE是一款开源、跨平台的轻量级 Go 语言集成开发环境（IDE）。

支持的操作系统：
 Windows x86 (32-bit or 64-bit)
 Linux x86 (32-bit or 64-bit)

下载地址：http://sourceforge.net/projects/liteide/files/
源码地址：https://github.com/visualfc/liteide


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>本书《Go语言四十二章经》内容在简书同步地址：  https://www.jianshu.com/nb/29056963
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com

