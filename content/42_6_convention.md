# 《Go语言四十二章经》第六章 约定和惯例

作者：李骁

## 6.1 可见性规则
包通过下面这个被编译器强制执行的规则来决定是否将自身的代码对象暴露给外部文件：

当标识符（包括常量、变量、类型、函数名、结构字段等等）以一个大写字母开头，如：Group1，那么使用这种形式的标识符的对象就可以被外部包的代码所使用（客户端程序需要先导入这个包），这被称为导出（像面向对象语言中的 public）；标识符如果以小写字母开头，则对包外是不可见的，但是他们在整个包的内部是可见并且可用的（像面向对象语言中的 private ）。

（大写字母可以使用任何 Unicode 编码的字符，比如希腊文，不仅仅是 ASCII 码中的大写字母）。

因此，在导入一个外部包后，能够且只能够访问该包中导出的对象。
假设在包 pack1 中我们有一个变量或函数叫做 Thing（以 T 开头，所以它能够被导出），那么在当前包中导入 pack1 包，Thing 就可以像面向对象语言那样使用点标记来调用：

```Go
pack1.Thing //（pack1 在这里是不可以省略的）
```
因此包也可以作为命名空间使用，帮助避免命名冲突（名称冲突）：两个包中的同名变量的区别在于他们的包名，例如 pack1.Thing 和 pack2.Thing。

> 注意事项：
如果你导入了一个包却没有使用它，则会在构建程序时引发错误，如 imported and not used: os，这正是遵循了 Go 的格言：“没有不必要的代码！”。

## 6.2 命名规范以及语法惯例

干净、可读的代码和简洁性是 Go 追求的主要目标。通过 Gofmt 来强制实现统一的代码风格。Go 语言中对象的命名也应该是简洁且有意义的。像 Java 和 Python 中那样使用混合着大小写和下划线的冗长的名称会严重降低代码的可读性。名称不需要指出自己所属的包，因为在调用的时候会使用包名作为限定符。返回某个对象的函数或方法的名称一般都是使用名词，没有 Get... 之类的字符，如果是用于修改某个对象，则使用 SetName。有必须要的话可以使用大小写混合的方式，如 MixedCaps 或 mixedCaps，而不是使用下划线来分割多个名称。

函数里的代码（函数体）使用大括号 {} 括起来。

左大括号 { 必须与方法的声明放在同一行，这是编译器的强制规定，否则你在使用 Gofmt 时就会出现错误提示：

Go 语言虽然看起来不使用分号作为语句的结束，但实际上这一过程是由编译器自动完成，因此才会引发像上面这样的错误。

右大括号 } 需要被放在紧接着函数体的下一行。如果你的函数非常简短，你也可以将它们放在同一行：

```Go
func Sum(a, b int) int { return a + b }
```
对于大括号 {} 的使用规则在任何时候都是相同的（如：if 语句等）。

因此符合规范的函数一般写成如下的形式：

```Go
func functionName(parameter_list) (return_value_list) {
   …
}
```
只有当某个函数需要被外部包调用的时候才使用大写字母开头，并遵循 Pascal 命名法；否则就遵循骆驼命名法，即第一个单词的首字母小写，其余单词的首字母大写。

单字之间不以空格断开或连接号（-）、底线（\_）连结，第一个单字首字母采用大写字母；后续单字的首字母亦用大写字母，例如：FirstName、LastName。每一个单字的首字母都采用大写字母的命名格式，被称为“Pascal命名法”，源自于Pascal语言的命名惯例，也有人称之为“大驼峰式命名法”（Upper Camel Case），为驼峰式大小写的子集。

帕斯卡命名法指当变量名和函式名称是由二个或二个以上单字连结在一起，而构成的唯一识别字时，用以增加变量和函式的可读性。

单行注释是最常见的注释形式，你可以在任何地方使用以 // 开头的单行注释。

多行注释也叫块注释，均以 /\* 开头，并以 \*/ 结尾，且不可以嵌套使用，多行注释一般用于包的文档描述或注释成块的代码片段。

```Go
// Cap returns the capacity of the buffer's underlying byte slice,
// that is, the total space allocated for the buffer's data.

/*
  Cap returns the capacity of the buffer's underlying byte slice,
  that is, the total space allocated for the buffer's data.
*/
```
在Go标准库中，一般都采用单行注释，建议采用官方标准方式。

每一个包应该有相关注释，在 package 语句之前的块注释将被默认认为是这个包的文档说明，其中应该提供一些相关信息并对整体功能做简要的介绍。一个包可以分散在多个文件中，但是只需要在其中一个进行注释说明即可。当开发人员需要了解包的一些情况时，自然会用 Godoc 来显示包的文档说明，在首行的简要注释之后可以用成段的注释来进行更详细的说明，而不必拥挤在一起。

另外，在多段注释之间应以空行分隔加以区分，单行注释的//后面空一格，方便godoc生成标准文档。

```Go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hash provides interfaces for hash functions.
package hash
 ```

几乎所有全局作用域的类型、常量、变量、函数和被导出的对象都应该有一个合理的注释。如果这种注释（称为文档注释）出现在函数前面，例如函数 Abcd，则要以 "Abcd..." 作为开头。

```Go
// enterOrbit causes Superman to fly into low Earth orbit, a position
// that presents several possibilities for planet salvation.
func enterOrbit() error {
   ...
}
```


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>本书《Go语言四十二章经》内容在简书同步地址：  https://www.jianshu.com/nb/29056963
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com