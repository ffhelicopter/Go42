# 《Go语言四十二章经》第十八章 Struct 结构体

作者：李骁

## 18.1结构体(struct)

Go 通过结构体的形式支持用户自定义类型，或者叫定制类型。

>Go 语言结构体是实现自定义类型的一种重要数据类型。
>
>结构体是复合类型（composite types），它由一系列属性组成，每个属性都有自己的类型和值的，结构体通过属性把数据聚集在一起。
>
>结构体类型和字段的命名遵循可见性规则。
>
>方法（Method）可以访问这些数据，就好像它们是这个独立实体的一部分。
>
>结构体是值类型，因此可以通过 new 函数来创建。

结构体是由一系列称为字段（fields）的命名元素组成，每个元素都有一个名称和一个类型。 字段名称可以显式指定（IdentifierList）或隐式指定（EmbeddedField），没有显式字段名称的字段称为匿名（内嵌）字段。在结构体中，非空字段名称必须是唯一的。

结构体定义的一般方式如下：

```Go
type identifier struct {
    field1 type1
    field2 type2
    ...
}
```

结构体里的字段一般都有名字，像 field1、field2 等，如果字段在代码中从来也不会被用到，那么可以命名它为 _。

空结构体如下所示：

```Go
struct {}
```

具有6个字段的结构体：

```Go
struct {
	x, y int
	u float32
	_ float32  // 填充
	A *[]int
	F func()
}
```

对于匿名字段，必须将匿名字段指定为类型名称T或指向非接口类型名称* T的指针，并且T本身可能不是指针类型。

```Go
struct {
	T1        // 字段名 T1
	*T2       // 字段名 T2
	P.T3      // 字段名 T3
	*P.T4     // f字段名T4
	x, y int    // 字段名 x 和 y
}
```

使用 new 函数给一个新的结构体变量分配内存，它返回指向已分配内存的指针：

```Go
type S struct { a int; b float64 }
new(S)
```

new(S)为S类型的变量分配内存，并初始化（a = 0，b = 0.0），返回包含该位置地址的类型* S的值。

我们一般的惯用方法是：t := new(T)，变量 t 是一个指向 T的指针，此时结构体字段的值是它们所属类型的零值。

也可以这样写：var t T ，也会给 t 分配内存，并零值化内存，但是这个时候 t 是类型T。

在这两种方式中，t 通常被称做类型 T 的一个实例（instance）或对象（object）。

使用点号符“.”可以获取结构体字段的值：structname.fieldname。在 Go 语言中“.”叫选择器（selector）。无论变量是一个结构体类型还是一个结构体类型指针，都使用同样的选择表示法来引用结构体的字段：

```Go
type myStruct struct { i int }
var v myStruct    // v是结构体类型变量
var p *myStruct   // p是指向一个结构体类型变量的指针
v.i
p.i

type Interval struct {
    start  int
    end   int
}
```

```Go
结构体变量有下面几种初始化方式，前面一种按照字段顺序，后面两种则对应字段名来初始化赋值：

intr := Interval{0, 3}            (A)
intr := Interval{end:5, start:1}    (B)
intr := Interval{end:5}           (C)
```

复合字面量是构造结构体，数组，切片和字典的值，并每次都创建新值。声明和初始化一个结构体实例（一个结构体字面量：struct-literal）方式如下：

定义结构体类型Point3D和Line：

```Go
type Point3D struct { x, y, z float64 }
type Line struct { p, q Point3D }
```

声明并初始化：

```Go
origin := Point3D{}                      //  Point3D 是零值
line := Line{origin, Point3D{y: -4, z: 12.3}}  //   line.q.x 是零值
```

这里 Point3D{}以及 Line{origin, Point3D{y: -4, z: 12.3}}都是结构体字面量。

表达式 new(Type) 和 &Type{} 是等价的。&struct1{a, b, c} 是一种简写，底层仍然会调用 new ()，这里值的顺序必须按照字段顺序来写。也可以通过在值的前面放上字段名来初始化字段的方式，这种方式就不必按照顺序来写了。

结构体类型和字段的命名遵循可见性规则，一个导出的结构体类型中有些字段是导出的，也即首字母大写字段会导出；另一些不可见，也即首字母小写为未导出，对外不可见。

## 18.2 结构体特性

* 结构体的内存布局
Go 语言中，结构体和它所包含的数据在内存中是以连续块的形式存在的，即使结构体中嵌套有其他的结构体，这在性能上带来了很大的优势。

* 递归结构体
结构体类型可以通过引用自身（指针类型）来定义。这在定义链表或二叉树的节点时特别有用，此时节点包含指向临近节点的链接。

```Go
type H struct {
	int
	*H
}
```

* 使用工厂方法
通过参考应用可见性规则，我们可以设定结构体名不能导出，就可以达到使用 new 函数，强制使用工厂方法的目的。

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

结构体中的字段除了有名字和类型外，还可以有一个可选的标签（tag）：它是一个附属于字段的字符串，可以是文档或其他的重要标记。标签的内容不可以在一般的编程中使用，只有 reflect 包能获取它。

reflect包可以在运行时自省类型、属性和方法，如变量是结构体类型，可以通过 Field 来索引结构体的字段，然后就可以使用 Tag 属性。

```Go
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	name string "学生名字"          // 结构体标签
	Age  int    "学生年龄"          // 结构体标签
	Room int    `json:"Roomid"` // 结构体标签
}

func main() {
	st := Student{"Titan", 14, 102}
	fmt.Println(reflect.TypeOf(st).Field(0).Tag)
	fmt.Println(reflect.TypeOf(st).Field(1).Tag)
	fmt.Println(reflect.TypeOf(st).Field(2).Tag)
}

程序输出：
学生名字
学生年龄
json:"Roomid"
```

从上面代码中可以看到，通过reflect我们很容易得到结构体字段的标签。

## 18.3 匿名成员

Go语言结构体中可以包含一个或多个匿名（内嵌）字段，即这些字段没有显式的名字，只有字段的类型是必须的，此时类型就是字段的名字（这一特征决定了在一个结构体中，每种数据类型只能有一个匿名字段）。

匿名（内嵌）字段本身也可以是一个结构体类型，即结构体可以包含内嵌结构体。

```Go
type Human struct {
	name string
}

type Student struct { // 含内嵌结构体Human
	Human // 匿名（内嵌）字段
	int   // 匿名（内嵌）字段
}
```

Go语言结构体中这种含匿名（内嵌）字段和内嵌结构体的结构，可近似地理解为面向对象语言中的继承概念。

Go 语言中的继承是通过内嵌或者说组合来实现的，所以可以说，在 Go 语言中，相比较于继承，组合更受青睐。

## 18.4 嵌入与聚合

结构体中包含匿名（内嵌）字段叫嵌入或者内嵌；而如果结构体中字段包含了类型名，还有字段名，则是聚合。聚合的在JAVA和C++都是常见的方式，而内嵌则是Go 的特有方式。

```Go
type Human struct {
	name string
}

type Person1 struct {           // 内嵌
	Human
}

type Person2 struct {           // 内嵌， 这种内嵌与上面内嵌有差异
	*Human
}

type Person3 struct{             // 聚合
	human Human
}
```

嵌入在结构体中广泛使用，在Go语言中如果只考虑结构体和接口的嵌入组合方式，一共有下面四种：

* 1.在接口中嵌入接口:

这里指的是在接口中定义中嵌入接口类型，而不是接口的一个实例，相当于合并了两个接口类型定义的全部函数。下面只有同时实现了Writer和 Reader 的接口，才可以说是实现了Teacher接口，即可以作为Teacher的实例。Teacher接口嵌入了Writer和 Reader 两个接口，在Teacher接口中，Writer和 Reader是两个匿名（内嵌）字段。

```Go
type Writer interface{
   Write()
}

type Reader interface{
   Read()
} 

type Teacher interface{
  Reader
  Writer
}
```

* 2.在接口中嵌入结构体: 

这种方式在Go语言中是不合法的，不能通过编译。

```Go
type Human struct {
	name string
}

type Writer interface {
	Write()
}

type Reader interface {
	Read()
}

type Teacher interface {
	Reader
	Writer
	Human
}
```

存在语法错误，并不具有实际的含义，编译报错: 
interface contains embedded non-interface Base

```Go
Interface 不能嵌入非interface的类型。
```

* 3.在结构体中内嵌接口:

初始化的时候，内嵌接口要用一个实现此接口的结构体赋值；或者定义一个新结构体，可以把新结构体作为receiver，实现接口的方法就实现了接口（先记住这句话，后面在讲述方法时会解释），这个新结构体可作为初始化时实现了内嵌接口的结构体来赋值。

```Go
package main

import (
	"fmt"
)

type Writer interface {
	Write()
}

type Author struct {
	name string
	Writer
}

// 定义新结构体，重点是实现接口方法Write()
type Other struct {
	i int
}

func (a Author) Write() {
	fmt.Println(a.name, "  Write.")
}

// 新结构体Other实现接口方法Write()，也就可以初始化时赋值给Writer 接口
func (o Other) Write() {
	fmt.Println(" Other Write.")
}

func main() {

	//  方法一：Other{99}作为Writer 接口赋值
	Ao := Author{"Other", Other{99}}
	Ao.Write()

	// 方法二：简易做法，对接口使用零值，可以完成初始化
	Au := Author{name: "Hawking"}
	Au.Write()
}

程序输出：
Other   Write.
Hawking   Write.
```

* 4.在结构体中嵌入结构体:

在结构体嵌入结构体很好理解，但不能嵌入自身值类型，可以嵌入自身的指针类型即递归嵌套。

在初始化时，内嵌结构体也进行赋值；外层结构自动获得内嵌结构体所有定义的字段和实现的方法。

下面代码完整演示了结构体中嵌入结构体，初始化以及字段的选择调用：

```Go
package main

import (
	"fmt"
)

type Human struct {
	name   string // 姓名
	Gender string // 性别
	Age    int    // 年龄
	string        // 匿名字段
}

type Student struct {
	Human     // 匿名字段
	Room  int // 教室
	int       // 匿名字段
}

func main() {
	//使用new方式
	stu := new(Student)
	stu.Room = 102
	stu.Human.name = "Titan"
	stu.Gender = "男"
	stu.Human.Age = 14
	stu.Human.string = "Student"

	fmt.Println("stu is:", stu)
	fmt.Printf("Student.Room is: %d\n", stu.Room)
	fmt.Printf("Student.int is: %d\n", stu.int) // 初始化时已自动给予零值：0
	fmt.Printf("Student.Human.name is: %s\n", stu.name) //  (*stu).name
	fmt.Printf("Student.Human.Gender is: %s\n", stu.Gender)
	fmt.Printf("Student.Human.Age is: %d\n", stu.Age)
	fmt.Printf("Student.Human.string is: %s\n", stu.string)

	// 使用结构体字面量赋值
	stud := Student{Room: 102, Human: Human{"Hawking", "男", 14, "Monitor"}}

	fmt.Println("stud is:", stud)
	fmt.Printf("Student.Room is: %d\n", stud.Room)
	fmt.Printf("Student.int is: %d\n", stud.int) // 初始化时已自动给予零值：0
	fmt.Printf("Student.Human.name is: %s\n", stud.Human.name)
	fmt.Printf("Student.Human.Gender is: %s\n", stud.Human.Gender)
	fmt.Printf("Student.Human.Age is: %d\n", stud.Human.Age)
	fmt.Printf("Student.Human.string is: %s\n", stud.Human.string)
}

程序输出：
stu is: &{ {Titan 男 14 Student} 102 0}
Student.Room is: 102
Student.int is: 0
Student.Human.name is: Titan
Student.Human.Gender is: 男
Student.Human.Age is: 14
Student.Human.string is: Student
stud is: { {Hawking 男 14 Monitor} 102 0}
Student.Room is: 102
Student.int is: 0
Student.Human.name is: Hawking
Student.Human.Gender is: 男
Student.Human.Age is: 14
Student.Human.string is: Monitor
```

内嵌结构体的字段，例如我们即可以stu.Human.name这样来选择使用，而如果外层结构体中没有同名的name字段，也可以stu.name直接来选择使用。对于嵌入和聚合结构体而言，我们在选择调用内部字段时，可以不用多层选择调用，在不同名情况下可直接调用。比如stu.name这样效果实际上与stu.Human.name一样。

我们通过对结构体使用new(T)，struct{filed:value}两种方式来声明初始化，这两种方式分别得到*T，和T。

我们从输出stu is: &{ {Titan 男 14 Student} 102 0} 可以得知，stu 是个指针，但是我们在随后调用字段时并没有使用指针，这是在Go语言中这里的 stu.name  相当于(*stu).name，这是一个语法糖，一般我们都使用stu.name方式来调用，但我们要知道有这个语法糖存在。

## 18.5 命名冲突
当两个字段拥有相同的名字（可能是继承来的名字）时该怎么办呢？外层名字会覆盖内层名字（但是两者的内存空间都保留）。

如果相同的名字在同一级别出现了两次，如果这个名字被程序使用了，将会引发一个错误，但不使用没关系。没有好办法来解决这种问题引起的二义性，一般由程序员完整写出来避免错误。

下面代码中如果写成 c.a 是错误的，因为我们不知道到底是要调用 c.A.a 还是 c.B.a。其实只要我们完整写出来（如:c.B.a）就不存在这个问题。

```Go
type A struct {a int}
type B struct {a, b int}

type C struct {A; B}
var c C
```


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>本书《Go语言四十二章经》内容在简书同步地址：  https://www.jianshu.com/nb/29056963
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
