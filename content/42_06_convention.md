# 《Go语言四十二章经》第六章 约定和惯例

作者：李骁

## 6.1 可见性规则

在Go语言中，标识符必须以一个大写字母开头，这样才可以被外部包的代码所使用，这被称为导出。标识符如果以小写字母开头，则对包外是不可见的，但是他们在整个包的内部是可见并且可用的。但是包名不管在什么情况下都必须小写。

在设计Go语言时，设计者们也希望确保它不是过于以ASCII为中心，这意味着需要从7位ASCII的范围来扩展标识符的空间。 所以Go语言标识符规定必须是Unicode定义的字母或数字，标识符是一个或多个Unicode字母和数字的序列， 标识符中的第一个字符必须是Unicode字母。

这条规则还有另外一个不幸的后果。由于导出的标识符必须以大写字母开头，因此根据定义，从某些语言的字符创建的标识符不能导出。目前唯一的解决方案是使用像“A语言”这样的东西，但这显然不能令人满意。

总而言之，为了确保我们的标识符能正常导出，我们建议在开发中还是尽量使用ASCII 码来作为标识符，虽然设计者们在避免以ASCII 码为中心，但出于习惯我们还是服从于这个现实。

>那么问题来了，使用中文命名的标识符能够正常导出吗？希望大家在了解后面的知识后，可以尝试一下试试。

## 6.2 命名规范以及语法惯例

当某个函数需要被外部包调用的时候需要使用大写字母开头，并遵循 Pascal 命名法（“大驼峰式命名法”）；否则就遵循“小驼峰式命名法”，即第一个单词的首字母小写，其余单词的首字母大写。

单词之间不以空格断开或连接号（-）、底线（_）连结，第一个单词首字母采用大写字母；后续单词的首字母亦用大写字母，例如：FirstName、LastName。每一个单词的首字母都采用大写字母的命名格式，被称为“Pascal命名法”，源自于Pascal语言的命名惯例，也有人称之为“大驼峰式命名法”（Upper Camel Case），为驼峰式大小写的子集。

当二个或二个以上单词连结在一起时，用驼峰式命名法可以增加变量和函数名称的可读性。

Go 语言追求简洁的代码风格，并通过 gofmt 强制实现风格统一。

Go 语言也使用分号作为语句的结束，但一般会省略分号。像在标识符后面；整数、浮点、复数、Rune或字符串等字面量后面；关键字break、continue、fallthrough、或者return后面；操作符或标点符号++、--、)、]或}之后等等都可以使用分号，但是往往会省略掉，像LiteIDE编辑器会在保存.go文件时自动过滤掉这些分号，所以在Go语言开发中一般不用过多关注分号的使用。

左大括号 { 不能单独一行，这是编译器的强制规定，否则你在使用 gofmt 时就会出现错误提示“ expected declaration, found '{' ”。右大括号 } 需要单独一行。

```Go
func functionName) () {
   …
}

if mod > 0 {
	div++
}
```

在定义接口名时也有惯例，一般单方法接口由方法名称加上-er后缀来命名。

## 6.3 注释

在Go语言中，注释有两种形式：

1.行注释：使用双斜线//开始，一般后面紧跟一个空格。行注释是Go语言中最常见的注释形式，在标准包中，一般都采用行注释，建议采用这种方式。
2.块注释：使用 /\* \*/，块注释不能嵌套。块注释一般用于包描述或注释成块的代码片段。

一般而言，注释文字尽量每行长度接近一致，过长的行应该换行以方便在编辑器阅读。注释可以是单行，多行，甚至可以使用doc.go文件来专门保存包注释。每个包只需要在一个go文件的package关键字上面注释，两者之间没有空行。对于变量，函数，结构体，接口等的注释直接加在声明前，注释与声明之间没有空行。例如：

```Go
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run genzfunc.go

// Package sort provides primitives for sorting slices and user-defined
// collections.
package sort

// A type, typically a collection, that satisfies sort.Interface can be
// sorted by the routines in this package. The methods require that the
// elements of the collection be enumerated by an integer index.
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

// Insertion sort
func insertionSort(data Interface, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}
```

函数或方法的注释需要以函数名开始，且两者之间没有空行，示例如下：

```Go
// ContainsRune reports whether the rune is contained in the UTF-8-encoded byte slice b.
func ContainsRune(b []byte, r rune) bool {
	return IndexRune(b, r) >= 0
}
```

需要预格式化的部分，直接加空格缩进即可，示例如下：

```Go
// For example, flags Ldate | Ltime (or LstdFlags) produce,
//	2009/01/23 01:23:23 message
// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
```

在方法，结构体或者包注释前面加上“Deprecated:”表示不建议使用，示例如下：

```Go
// Deprecated: Old 老旧方法，不建议使用
func Old(a int)(int){
    return a
}
```

在注释中，还可以插入空行，示例如下：

```Go
// Search calls f(i) only for i in the range [0, n).
//
// A common use of Search is to find the index i for a value x in
// a sorted, indexable data structure such as an array or slice.
// In this case, the argument f, typically a closure, captures the value
// to be searched for, and how the data structure is indexed and
// ordered.
//
// For instance, given a slice data sorted in ascending order,
// the call Search(len(data), func(i int) bool { return data[i] >= 23 })
// returns the smallest index i such that data[i] >= 23. If the caller
// wants to find whether 23 is in the slice, it must test data[i] == 23
// separately.
```



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第五章 作用域](https://github.com/ffhelicopter/Go42/blob/master/content/42_05_scope.md)

[第七章 代码结构化](https://github.com/ffhelicopter/Go42/blob/master/content/42_07_package.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
