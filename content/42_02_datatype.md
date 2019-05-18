# 《Go语言四十二章经》第二章 数据类型

作者：李骁

在 Go 语言中，数据类型可用于参数和变量声明。

## 2.1 基本数据类型

Go 语言按类别有以下几种数据类型：

* 布尔型：<br>
布尔型的值只可以是常量 true 或者 false。一个简单的例子：`var b bool = true`。

* 数字类型：<br>
整型 int 和浮点型 float32、float64，Go 语言支持整型和浮点型数字，并且原生支持复数，其中位的运算采用补码。

* 字符串类型：<br>
字符串就是一串固定长度的字符连接起来的字符序列。Go的字符串是由单个字节连接起来的。Go语言的字符串的字节使用UTF-8编码标识Unicode文本。

* 派生类型：<br>
包括：


    (a) 指针类型（Pointer）
    (b) 数组类型
    (c) 结构类型(struct)
    (d) Channel 类型
    (e) 函数类型
    (f) 切片类型
    (g) 接口类型（interface）
    (h) Map 类型

### 数字类型：

Go 也有基于架构的类型，例如：int、uint 和 uintptr，这些类型的长度都是根据运行程序所在的操作系统类型所决定的。


|类型|符号|长度范围|
|:--|:--|:--|
|uint8    |无符号  |8位整型 (0 到 255)|
|uint16   |无符号 |16位整型 (0 到 65535)|
|uint32   |无符号 |32位整型 (0 到 4294967295)|
|uint64   |无符号 |64位整型 (0 到 18446744073709551615)|
|int8     |有符号  |8位整型 (-128 到 127)|
|int16    |有符号 |16位整型 (-32768 到 32767)|
|int32    |有符号 |32位整型 (-2147483648 到 2147483647)|
|int64    |有符号 |64位整型 (-9223372036854775808 到 9223372036854775807)|


## 浮点型：

主要是为了表示小数，也可细分为float32和float64两种。浮点数能够表示的范围可以从很小到很巨大，这个极限值范围可以在math包中获取，math.MaxFloat32表示float32的最大值，大约是3.4e38，math.MaxFloat64大约是1.8e308，两个类型最小的非负值大约是1.4e-45和4.9e-324。


float32大约可以提供小数点后6位的精度，作为对比，float64可以提供小数点后15位的精度。通常情况应该优先选择float64，因此float32的精确度较低，在累积计算时误差扩散很快，而且float32能精确表达的最小正整数并不大，因为浮点数和整数的底层解释方式完全不同。

|类型|长度|
|:--|:--|
|float32  |IEEE-754   32位浮点型数|
|float64  |IEEE-754   64位浮点型数|

## 其他数字类型：

|类型|长度|
|:--|:--|
|byte      |类似 uint8|
|rune      |类似 int32|
|uint32     |或 64 位|
|int        |与 uint 一样大小|
|uintptr    |无符号整型，用于存放一个指针|

## 字符串：
只读的Unicode字节序列，Go语言使用UTF-8格式编码Unicode字符，每个字符对应一个rune类型。一旦字符串变量赋值之后，内部的字符就不能修改，英文是一个字节，中文是三个字节。

```Go
string转int：    int, err := strconv.Atoi(string)
string转int64：  int64, err := strconv.ParseInt(string, 10, 64)
int转string：    string := strconv.Itoa(int)
int64转string：  string := strconv.FormatInt(int64, 10)
```

而一个range循环会在每次迭代时，解码一个UTF-8编码的符文。每次循环时，循环的索引是当前文字的起始位置，以字节为单位，代码点是它的值（rune）。

使用range迭代字符串时，需要注意的是range迭代的是Unicode而不是字节。返回的两个值，第一个是被迭代的字符的UTF-8编码的第一个字节在字符串中的索引，第二个值的为对应的字符且类型为rune(实际就是表示unicode值的整形数据）。

```Go
const s = "Go语言"
for i, r := range s {
	fmt.Printf("%#U  ： %d\n", r, i)
}
```
程序输出：

U+0047 'G'   ： 0<br>
U+006F 'o'   ： 1<br>
U+8BED '语'  ： 2<br>
U+8A00 '言'  ： 5<br>

## 复数：
复数类型相对用的很少，主要是数学学科专业会用上。分为两种类型 complex64和complex128 前部分是实体后部分是虚体。

|类型|长度|
|:--|:--|
|complex64   |32位实数和虚数|
|complex128   |64位实数和虚数|

## 2.2 Unicode（UTF-8）

你可以通过增加前缀 0 来表示 8 进制数（如：077），增加前缀 0x 来表示 16 进制数（如：0xFF），以及使用 e 来表示 10 的连乘（如： 1e3 = 1000，或者 6.022e23 = 6.022 x 1e23）

不过 Go 同样支持 Unicode（UTF-8），因此字符同样称为 Unicode 代码点或者 runes，并在内存中使用 int 来表示。在文档中，一般使用格式 U+hhhh 来表示，其中 h 表示一个 16 进制数。其实 rune 也是 Go 当中的一个类型，并且是 int32 的别名。

在书写 Unicode 字符时，需要在 16 进制数之前加上前缀 \u 或者 \U。

因为 Unicode 至少占用 2 个字节，所以我们使用 int16 或者 int 类型来表示。如果需要使用到 4 字节，则会加上 \U 前缀；前缀 \u 则总是紧跟着长度为 4 的 16 进制数，前缀 \U 紧跟着长度为 8 的 16 进制数。

```Go
var ch int = '\u0041'
var ch2 int = '\u03B2'
var ch3 int = '\U00101234'
```

## 2.3 复数
Go 拥有以下复数类型：

complex64 (32 位实数和虚数)<br>
complex128 (64 位实数和虚数)<br>

复数使用 re+imi 来表示，其中 re 代表实数部分，im 代表虚数部分，i 为虚数单位。
示例：

```Go
var c1 complex64 = 5 + 10i
fmt.Printf("The value is: %v", c1)// 输出： 5 + 10i
```
如果 re 和 im 的类型均为 float32，那么类型为 complex64 的复数 c 可以通过以下方式来获得：

```Go
c = complex(re, im)
```
函数 real(c) 和 imag(c) 可以分别获得相应的实数和虚数部分。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第一章 Go安装与运行](https://github.com/ffhelicopter/Go42/blob/master/content/42_01_install.md)

[第三章 变量](https://github.com/ffhelicopter/Go42/blob/master/content/42_03_var.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com