# <center>《Go语言四十二章经》第十八章 Struct 结构体</center>

作者：李骁

## 18.1结构体(struct)
Go 通过结构体的形式支持用户自定义类型，或者叫定制类型。

>一个带属性的结构体试图表示一个现实世界中的实体。
>
>结构体是复合类型（composite types），当需要定义一个类型，它由一系列属性组成，每个属性都有自己的类型和值的时候，就应该使用结构体，它把数据聚集在一起。
>
>然后（方法）可以访问这些数据，就好像它们是一个独立实体的一部分。
>
>结构体是值类型，因此可以通过 new 函数来创建。

组成结构体类型的那些数据称为字段（fields）。每个字段都有一个类型和一个名字；在一个结构体中，字段名字必须是唯一的。
结构体定义的一般方式如下：

```Go
type identifier struct {
    field1 type1
    field2 type2
    ...
}
```
结构体里的字段都有 名字，像 field1、field2 等，如果字段在代码中从来也不会被用到，那么可以命名它为 _。

**使用 new**
使用 new 函数给一个新的结构体变量分配内存，它返回指向已分配内存的指针：var t *T = new(T)，如果需要可以把这条语句放在不同的行（比如定义是包范围的，但是分配却没有必要在开始就做）。

```Go
var t *T
t = new(T)
```
写这条语句的惯用方法是：t := new(T)，变量 t 是一个指向 T的指针，此时结构体字段的值是它们所属类型的零值。

声明 var t T 也会给 t 分配内存，并零值化内存，但是这个时候 t 是类型T。在这两种方式中，t 通常被称做类型 T 的一个实例（instance）或对象（object）。

同样的，使用点号符可以获取结构体字段的值：structname.fieldname。
在 Go 语言中这叫 选择器（selector）。无论变量是一个结构体类型还是一个结构体类型指针，都使用同样的 选择器符（selector-notation） 来引用结构体的字段：

```Go
type myStruct struct { i int }
var v myStruct    // v是结构体类型变量
var p *myStruct   // p是指向一个结构体类型变量的指针
v.i
p.i

type Interval struct {
    start int
    end   int
}
```
初始化方式：

```Go
intr := Interval{0, 3}            (A)
intr := Interval{end:5, start:1}    (B)
intr := Interval{end:5}           (C)
```
初始化一个结构体实例（一个结构体字面量：struct-literal）的更简短和惯用的方式如下：

```Go
    ms := &struct1{10, 15.5, "Chris"}
    // 此时ms的类型是 *struct1
```
或者：

```Go
    var ms struct1
    ms = struct1{10, 15.5, "Chris"}
```

**混合字面量语法**（composite literal syntax）

&struct1{a, b, c} 是一种简写，底层仍然会调用 new ()，这里值的顺序必须按照字段顺序来写。在下面的例子中能看到可以通过在值的前面放上字段名来初始化字段的方式。

表达式 new(Type) 和 &Type{} 是等价的。

**结构体类型和字段的命名遵循可见性规则，一个导出的结构体类型中有些字段是导出的，另一些不可见。**

## 18.2 结构体特性

* 结构体的内存布局<br>
Go 语言中，结构体和它所包含的数据在内存中是以连续块的形式存在的，即使结构体中嵌套有其他的结构体，这在性能上带来了很大的优势。

* 递归结构体<br>
结构体类型可以通过引用自身来定义。这在定义链表或二叉树的元素（通常叫节点）时特别有用，此时节点包含指向临近节点的链接（地址）。如下所示，链表中的 su，树中的 ri 和 le 分别是指向别的节点的指针。

* 链表<br>
这块的 data 字段用于存放有效数据（比如 float64），su 指针指向后继节点。

Go 代码：

```Go
type Node struct {
    data    float64
    su      *Node
}
```
链表中的第一个元素叫 head，它指向第二个元素；最后一个元素叫 tail，它没有后继元素，所以它的 su 为 nil 值。当然真实的链接会有很多数据节点，并且链表可以动态增长或收缩。
同样地可以定义一个双向链表，它有一个前趋节点 pr 和一个后继节点 su：

```Go
type Node struct {
    pr      *Node
    data    float64
    su      *Node
}
```
* 二叉树<br>
二叉树中每个节点最多能链接至两个节点：左节点（le）和右节点（ri），这两个节点本身又可以有左右节点，依次类推。树的顶层节点叫根节点（root），底层没有子节点的节点叫叶子节点（leaves），叶子节点的 le 和 ri 指针为 nil 值。在 Go 中可以如下定义二叉树：

```Go
type Tree strcut {
    le      *Tree
    data    float64
    ri      *Tree
}
```
* 结构体工厂<br>

Go 语言不支持面向对象编程语言中那样的构造子方法，但是可以很容易的在 Go 中实现 “构造子工厂”方法。为了方便通常会为类型定义一个工厂，按惯例，工厂的名字以 new 或 New 开头。假设定义了如下的 File 结构体类型：
```Go
type File struct {
    fd      int     // 文件描述符
    name    string  // 文件名
}
```
下面是这个结构体类型对应的工厂方法，它返回一个指向结构体实例的指针：

```Go
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }

    return &File{fd, name}
}
```
然后这样调用它：

```Go
f := NewFile(10, "./test.txt")
```
在 Go 语言中常常像上面这样在工厂方法里使用初始化来简便的实现构造函数。

如果 File 是一个结构体类型，那么表达式 new(File) 和 &File{} 是等价的。
这可以和大多数面向对象编程语言中笨拙的初始化方式做个比较：File f = new File(...)。
我们可以说是工厂实例化了类型的一个对象，就像在基于类的OO语言中那样。
如果想知道结构体类型T的一个实例占用了多少内存，可以使用：size := unsafe.Sizeof(T{})。


* 如何强制使用工厂方法

通过应用可见性规则参考，就可以禁止使用 new 函数，强制用户使用工厂方法，从而使类型变成私有的，就像在面向对象语言中那样。

```Go
type matrix struct {
    ...
}

func NewMatrix(params) *matrix {
    m := new(matrix) // 初始化 m
    return m
}
```
在包外，只有通过NewMatrix函数才可以初始化matrix 结构。

* 带标签的结构体

结构体中的字段除了有名字和类型外，还可以有一个可选的标签（tag）：它是一个附属于字段的字符串，可以是文档或其他的重要标记。标签的内容不可以在一般的编程中使用，只有包 reflect 能获取它。reflect包可以在运行时自省类型、属性和方法，比如：在一个变量上调用 reflect.TypeOf() 可以获取变量的正确类型，如果变量是一个结构体类型，就可以通过 Field 来索引结构体的字段，然后就可以使用 Tag 属性。

```Go
package main

import (
    "fmt"
    "reflect"
)

type TagType struct { // 结构体标签
    field1 bool   "An important answer"
    field2 string "The name of the thing"
    field3 int    "How much there are"
}

func main() {
    tt := TagType{true, "Barak Obama", 1}
    for i := 0; i < 3; i++ {
        refTag(tt, i)
    }
}

func refTag(tt TagType, ix int) {
    ttType := reflect.TypeOf(tt)
    ixField := ttType.Field(ix)
    fmt.Printf("%v\n", ixField.Tag)
}
```

```Go
程序输出：

An important answer
The name of the thing
How much there are
```

## 18.3 匿名成员

Go语言有一个特性让我们只声明一个成员对应的数据类型而不指名成员的名字；这类成员就叫匿名成员。匿名成员的数据类型必须是命名的类型或指向一个命名的类型的指针。

结构体可以包含一个或多个 匿名（或内嵌）字段，即这些字段没有显式的名字，只有字段的类型是必须的，此时类型就是字段的名字（这决定了在一个结构体中对于每一种数据类型只能有一个匿名字段。）。匿名字段本身可以是一个结构体类型，即 结构体可以包含内嵌结构体。

```Go
type Base struct {
   basename string
}

type Derive struct { // 含内嵌结构体
   Base   // 匿名 
   int
}
```
可以粗略地将这个和面向对象语言中的继承概念相比较，随后将会看到它被用来模拟类似继承的行为。Go 语言中的继承是通过内嵌或组合来实现的，所以可以说，在 Go 语言中，相比较于继承，组合更受青睐。

## 18.4 内嵌(embeded)结构体

**内嵌与聚合：**  
外部类型只包含了内部类型的类型名， 而没有field 名， 则是内嵌。外部类型包含了内部类型的类型名，还有filed名，则是聚合。聚合的在JAVA和C++都是常见的方式。而内嵌则是Go 的特有方式。
 
```Go
type Base struct {
  basename string
}
 
type Derive struct {           // 内嵌
   Base
}

type Derive struct {           // 内嵌， 这种内嵌与上面内嵌有差异
  *Base
}

type Derive struct{             // 聚合
  base Base
}
```

>内嵌的方式： 
>主要是通过结构体和接口的组合，有四种。

* 接口中内嵌接口:

这里的做为内嵌接口的含义实际上还是指的一个定义，而不是接口的一个实例，相当于合并了两个接口定义的函数，只有同时了Writer和 Reader 接口，是可以说是实现了WRer接口，即才可以作为WRer的实例。

```Go
type Writer interface{
   Write()
}

type Reader interface{
   Read()
} 

type WRer interface{
  Reader
  Writer
}
```
* 在接口中内嵌struct : 

存在语法错误，并不具有实际的含义， 编译报错: 
```Go
interface contains embedded non-interface Person

Interface 不能嵌入非interface的类型。
```

* 在结构体（struct）中内嵌 接口（interface）

1，初始化的时候，内嵌接口要用一个实现此接口的结构体赋值。

2，外层结构体中，只能调用内层接口定义的函数。 这是由于编译时决定。

3，外层结构体，可以作为receiver，重新定义同名函数，这样可以覆盖内层内嵌结构中定义的函数。

4，如果上述第3条实现，那么可以用外层结构体引用内嵌接口的实例，并调用内嵌接口的函数。

```Go
package main

import (
	"fmt"
)

type Printer interface {
	Print()
}

type CanonPrinter struct {
	Printname string
}

func (printer CanonPrinter) Print() {
	fmt.Println("this is cannoprinter printing now ")
}

type PrintWorker struct {
	Printer
	name string
	age  int
}

// 如果没有下面实现，则
func (printworker PrintWorker) Print() {
	fmt.Println("this is printing from PrintWorker ")
	printworker.Printer.Print()
	// 这里 printworker 首先引用内部嵌入Printer接口的实例，
// 然后调用Printer 接口实例的Print()方法
}

func main() {
	canon := CanonPrinter{"canoprint_num_1"}
	printworker := PrintWorker{Printer: canon, name: "ansendong", age: 34}
	printworker.Print()
	// 如果没有上述部分Func (printworker PrintWorker) Print()的实现，
// 则这里只调用CanonPrinter实现的Print()方法。
}
```

* 结构体（struct）中内嵌 结构体（struct）

1，初始化，内嵌结构体要进行赋值。

2，外层结构自动获得内嵌结构体所有定义的field和实现的方法（method）。

3，同上述结构体中内嵌接口类似，同样外层结构体可以定义同名方法，这样覆盖内层结构体的定义的方法。 同样也可以定义同名变量，覆盖内层结构体的变量。

4，同样可以内层结构体引用，内层结构体重已经定义的方法和变量。

同样地结构体也是一种数据类型，所以它也可以作为一个匿名字段来使用，如同下面例子中那样。外层结构体通过 outer.in1 直接进入内层结构体的字段，内嵌结构体甚至可以来自其他包。内层结构体被简单的插入或者内嵌进外层结构体。这个简单的“继承”机制提供了一种方式，使得可以从另外一个或一些类型继承部分或全部实现。

```Go
package main

import "fmt"

type innerS struct {
    in1 int
    in2 int
}

type outerS struct {
    b    int
    c    float32
    int  // anonymous field
    innerS //anonymous field
}

func main() {
    outer := new(outerS)
    outer.b = 6
    outer.c = 7.5
    outer.int = 60
    outer.in1 = 5
    outer.in2 = 10

    fmt.Printf("outer.b is: %d\n", outer.b)
    fmt.Printf("outer.c is: %f\n", outer.c)
    fmt.Printf("outer.int is: %d\n", outer.int)
    fmt.Printf("outer.in1 is: %d\n", outer.in1)
    fmt.Printf("outer.in2 is: %d\n", outer.in2)

    // 使用结构体字面量
    outer2 := outerS{6, 7.5, 60, innerS{5, 10}}
    fmt.Println("outer2 is:", outer2)
}
```

```Go
程序输出：

outer.b is: 6
outer.c is: 7.500000
outer.int is: 60
outer.in1 is: 5
outer.in2 is: 10
outer2 is:{6 7.5 60 {5 10}}
```
## 18.5 命名冲突
当两个字段拥有相同的名字（可能是继承来的名字）时该怎么办呢？
外层名字会覆盖内层名字（但是两者的内存空间都保留），这提供了一种重载字段或方法的方式；
如果相同的名字在同一级别出现了两次，如果这个名字被程序使用了，将会引发一个错误（不使用没关系）。没有办法来解决这种问题引起的二义性，必须由程序员自己修正。

使用 c.a 是错误的，到底是 c.A.a 还是 c.B.a。但可以完整写出来避免错误。

```Go
type A struct {a int}
type B struct {a, b int}

type C struct {A; B}
var c C
```