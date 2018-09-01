第十章 string
10.1 有关string
Go 语言中的string类型存储的字符串是不可变的， 如果要修改string内容需要将string转换为[]byte或[]rune，并且修改后的string内容是重新分配的。

那么byte和rune的区别是什么(下面写法是type别名):
type byte = uint8 
type rune = int32

string 类型的零值为长度为零的字符串，即空字符串 ""。

一般的比较运算符（==、!=、<、<=、>=、>）通过在内存中按字节比较来实现字符串的对比。你可以通过函数 len() 来获取字符串所占的字节长度，例如：len(str)。

字符串的内容（纯字节）可以通过标准索引法来获取，在中括号 [] 内写入索引，索引从 0 开始计数：
    字符串 str 的第 1 个字节：str[0]
    第 i 个字节：str[i - 1]
    最后 1 个字节：str[len(str)-1]

需要注意的是，这种转换方案只对纯 ASCII 码的字符串有效。

注意事项：

获取字符串中某个字节的地址的行为是非法的，例如：&str[i]。

10.2 字符串拼接
可以通过以下方式来对代码中多行的字符串进行拼接。
直接使用运算符

str := "Beginning of the string " +
"second part of the string"

由于编译器行尾自动补全分号的缘故，加号 + 必须放在第一行。
拼接的简写形式 += 也可以用于字符串：

s := "hel" + "lo, "
s += "world!"
fmt.Println(s) // 输出 “hello, world!”

里面的字符串都是不可变的，每次运算都会产生一个新的字符串，所以会产生很多临时的无用的字符串，不仅没有用，还会给 gc 带来额外的负担，所以性能比较差。

fmt.Sprintf()

fmt.Sprintf("%d:%s", 2018, "年")

内部使用 []byte 实现，不像直接运算符这种会产生很多临时的字符串，但是内部的逻辑比较复杂，有很多额外的判断，还用到了 interface，所以性能一般。

strings.Join()

strings.Join([]string{"hello", "world"}, ", ")

Join会先根据字符串数组的内容，计算出一个拼接之后的长度，然后申请对应大小的内存，一个一个字符串填入，在已有一个数组的情况下，这种效率会很高，但是本来没有，去构造这个数据的代价也不小。

bytes.Buffer

var buffer bytes.Buffer
	buffer.WriteString("hello")
	buffer.WriteString(", ")
	buffer.WriteString("world")

	fmt.Print(buffer.String())

这个比较理想，可以当成可变字符使用，对内存的增长也有优化，如果能预估字符串的长度，还可以用 buffer.Grow() 接口来设置 capacity。

strings.Builder

var b1 strings.Builder
	b1.WriteString("ABC")
	b1.WriteString("DEF")

	fmt.Print(b1.String())

strings.Builder 内部通过 slice 来保存和管理内容。slice 内部则是通过一个指针指向实际保存内容的数组。strings.Builder 同样也提供了 Grow() 来支持预定义容量。当我们可以预定义我们需要使用的容量时，strings.Builder 就能避免扩容而创建新的 slice 了。strings.Builder是非线程安全，性能上和bytes.Buffer 相差无几。
