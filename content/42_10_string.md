# 《Go语言四十二章经》第十章 string

作者：李骁

## 10.1 字符串介绍

Go 语言中可以使用反引号或者双引号来定义字符串。反引号表示原生的字符串，即不进行转义。

1. 双引号：字符串使用双引号括起来，其中的相关的转义字符将被替换。例如：

```Go
str := "Hello World! \n Hello Gopher! \n"

输出：
Hello World! 
Hello Gopher!
```

2. 反引号：字符串使用反引号括起来，其中的相关的转义字符不会被替换。例如：

```Go
str :=  `Hello World! \n Hello Gopher! \n` 

输出：
Hello World! \nHello Gopher! \n
```

双引号中的转义字符被替换，而反引号中原生字符串中的 \n 会被原样输出。


Go 语言中的string类型是一种值类型，存储的字符串是不可变的，如果要修改string内容需要将string转换为[]byte或[]rune，并且修改后的string内容是重新分配的。

那么byte和rune的区别是什么(下面写法是type别名):

```Go
type byte = uint8
type rune = int32
```
从上面的定义中我们可清楚看到两者的区别。

而string类型的零值是为长度为零的字符串，即空字符串 ""。

一般的比较运算符（==、!=、<、<=、>=、>）通过在内存中按字节比较来实现字符串的对比。你可以通过函数 len() 来获取字符串所占的字节长度，例如：len(str)。

字符串的内容（纯字节）可以通过标准索引法来获取，在中括号 [] 内写入索引，索引从 0 开始计数：

字符串 str 的第 1 个字节：str[0]
第 i 个字节：str[i - 1]
最后 1 个字节：str[len(str)-1]

需要注意的是，在Go语言代码使用 UTF-8 编码，同时标识符也支持 Unicode 字符。在标准库 unicode 包中，提供了对 Unicode 相关编码、解码的支持。而UTF8编码由Go语言之父Ken Thompson和Rob Pike共同发明的，现在已经是Unicode的标准。

Go语言默认使用UTF-8编码，对Unicode的支持非常好。但这也带来一个问题，也就是很多资料中提到的“获取字符串长度”的问题。内置的len()函数获取的是每个字符的UTF-8编码的长度和，而不是直接的字符数量。

```Go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {

	s := "其实就是rune"
	fmt.Println(len(s))                    // "16"
	fmt.Println(utf8.RuneCountInString(s)) // "8"
}
```

如字符串含有中文等字符，我们可以看到每个中文字符的索引值相差3。下面代码同时说明了在for range循环处理字符时，不是按照字节的方式来处理的。v其实际上是一个rune类型值。实际上，Go语言的range循环在处理字符串的时候，会自动隐式解码UTF8字符串。

```Go
package main

import (
	"fmt"
)

func main() {
	s := "Go语言四十二章经"
	for k, v := range s {
		fmt.Printf("k：%d,v：%c == %d\n", k, v, v)
	}
}
```

```Go
程序输出：
k：0,v：G == 71
k：1,v：o == 111
k：2,v：语 == 35821
k：5,v：言 == 35328
k：8,v：四 == 22235
k：11,v：十 == 21313
k：14,v：二 == 20108
k：17,v：章 == 31456
k：20,v：经 == 32463
```



>注意事项：
>
>获取字符串中某个字节的地址的行为是非法的，例如：&str[i]。


## 10.2 字符串拼接

可以通过以下方式来对代码中多行的字符串进行拼接。
* 直接使用运算符

```Go
str := "Beginning of the string " +
"second part of the string"
```

由于编译器行尾自动补全分号的缘故，加号 + 必须放在第一行。
拼接的简写形式 += 也可以用于字符串：

```Go
s := "hel" + "lo, "
s += "world!"
fmt.Println(s) // 输出 “hello, world!”
```

里面的字符串都是不可变的，每次运算都会产生一个新的字符串，所以会产生很多临时的无用的字符串，不仅没有用，还会给 GC 带来额外的负担，所以性能比较差。

* fmt.Sprintf()

```Go
fmt.Sprintf("%d:%s", 2018, "年")
```

内部使用 []byte 实现，不像直接运算符这种会产生很多临时的字符串，但是内部的逻辑比较复杂，有很多额外的判断，还用到了 interface，所以性能一般。

* strings.Join()

```Go
strings.Join([]string{"hello", "world"}, ", ")
```

Join会先根据字符串数组的内容，计算出一个拼接之后的长度，然后申请对应大小的内存，一个一个字符串填入，在已有一个数组的情况下，这种效率会很高，但是本来没有，去构造这个数据的代价也不小。

* bytes.Buffer

```Go
var buffer bytes.Buffer
buffer.WriteString("hello")
buffer.WriteString(", ")
buffer.WriteString("world")

fmt.Print(buffer.String())
```

这个比较理想，可以当成可变字符使用，对内存的增长也有优化，如果能预估字符串的长度，还可以用 buffer.Grow() 接口来设置 capacity。

* strings.Builder

```Go
var b1 strings.Builder
b1.WriteString("ABC")
b1.WriteString("DEF")

fmt.Print(b1.String())
```

strings.Builder 内部通过 slice 来保存和管理内容。slice 内部则是通过一个指针指向实际保存内容的数组。strings.Builder 同样也提供了 Grow() 来支持预定义容量。当我们可以预定义我们需要使用的容量时，strings.Builder 就能避免扩容而创建新的 slice 了。strings.Builder是非线程安全，性能上和 bytes.Buffer 相差无几。


## 10.3 有关string处理

标准库中有四个包对字符串处理尤为重要：bytes、strings、strconv和unicode包。

strings包提供了许多如字符串的查询、替换、比较、截断、拆分和合并等功能。

bytes包也提供了很多类似功能的函数，但是针对和字符串有着相同结构的[]byte类型。因为字符串是只读的，因此逐步构建字符串会导致很多分配和复制。在这种情况下，使用bytes.Buffer类型将会更有效，稍后我们将展示。

strconv包提供了布尔型、整型数、浮点数和对应字符串的相互转换，还提供了双引号转义相关的转换。

unicode包提供了IsDigit、IsLetter、IsUpper和IsLower等类似功能，它们用于给字符分类。

strings 包提供了很多操作字符串的简单函数，通常一般的字符串操作需求都可以在这个包中找到。下面简单举几个例子：

判断是否以某字符串打头/结尾
strings.HasPrefix(s, prefix string) bool
strings.HasSuffix(s, suffix string) bool

字符串分割
strings.Split(s, sep string) []string

返回子串索引
strings.Index(s, substr string) int
strings.LastIndex 最后一个匹配索引

字符串连接
strings.Join(a []string, sep string) string
另外可以直接使用“+”来连接两个字符串

字符串替换
strings.Replace(s, old, new string, n int) string

字符串转化为大小写
strings.ToUpper(s string) string
strings.ToLower(s string) string

统计某个字符在字符串出现的次数
strings.Count(s, substr string) int

判断字符串的包含关系
strings.Contains(s, substr string) bool



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第九章 运算符](https://github.com/ffhelicopter/Go42/blob/master/content/42_09_operator.md)

[第十一章 数组(Array)](https://github.com/ffhelicopter/Go42/blob/master/content/42_11_array.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
