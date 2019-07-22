# 《Go语言四十二章经》第二十章 方法

作者：李骁

在前面我们讲了结构体（struct）和接口（interface），在里面也提到过方法，但没有详细介绍方法（Method）。在这一章里，我们来仔细看看方法有那些奇妙之处呢？

## 20.1 方法的定义

在 Go 语言中，结构体就像是类的一种简化形式，那么面向对象程序员可能会问：类的方法在哪里呢？在 Go 语言中有一个概念，它和方法有着同样的名字，并且大体上意思相近。

Go 语言中方法和函数在形式上很像，它是作用在接收器（receiver）上的一个函数，接收器是某种类型的变量。因此方法是一种特殊类型的函数，方法只是比函数多了一个接收器（receiver），当然在接口中定义的函数我们也称为方法（因为最终还是要通过绑定到类型来实现）。

正是因为有了接收器，方法才可以作用于接收器的类型（变量）上，类似于面向对象中类的方法可以作用于类属性上。

定义方法的一般格式如下：

```Go
func (recv receiver_type) methodName(parameter_list) (return_value_list) { ... }
```
在方法名之前，func 关键字之后的括号中指定接收器 receiver。

```Go
type A struct {
	Face int
}

func (a A) f() {
	fmt.Println("hi ", a.Face)
}
```
上面代码中，我们定义了结构体 A ，注意f()就是 A 的方法，(a A)表示接收器。a 是 A的实例，f()是它的方法名，方法调用遵循传统的 object.name 即选择器符号：a.f()。

**接收器(receiver)**

* 接收器类型除了不能是指针类型或接口类型外，可以是其他任何类型，不仅仅是结构体类型，也可以是函数类型，还可以是 int、bool、string 等等为基础的自定义类型。

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

func (h Human) print() { // 值方法
	fmt.Println("Human:", h)
}

type MyInt int

func (m MyInt) print() { // 值方法
	fmt.Println("MyInt:", m)
}

func main() {
	//使用new方式
	hu := new(Human)
	hu.name = "Titan"
	hu.Gender = "男"
	hu.Age = 14
	hu.string = "Student"
	hu.print()

	// 指针变量
	mi := new(MyInt)
	mi.print()

	// 使用结构体字面量赋值
	hum := Human{"Hawking", "男", 14, "Monitor"}
	hum.print()

	// 值变量
	myi := MyInt(99)
	myi.print()
}

程序输出：
Human: {Titan 男 14 Student}
MyInt: 0
Human: {Hawking 男 14 Monitor}
MyInt: 99
```

* 接收器不能是一个接口类型，因为接口是一个抽象定义，但是方法却是具体实现；如果这样做会引发一个编译错误：invalid receiver type…。

```Go
package main

import (
	"fmt"
)

type printer interface {
	print()
}

func (p printer) print() { //  invalid receiver type printer (printer is an interface type)
	fmt.Println("printer:", p)
}
func main() {}
```

* 接收器不能是一个指针类型，但是它可以是任何其他允许类型的指针。

```Go
package main

import (
	"fmt"
)

type MyInt int

type Q *MyInt

func (q Q) print() { // invalid receiver type Q (Q is a pointer type)
	fmt.Println("Q:", q)
}

func main() {}
```

接收器不能是指针类型，但可以是类型的指针，有点绕口。下面我们看个例子：

```Go
package main

import (
	"fmt"
)

type MyInt int

func (mi *MyInt) print() { // 指针接收器，指针方法
	fmt.Println("MyInt:", *mi)
}
func (mi MyInt) echo() { // 值接收器，值方法
	fmt.Println("MyInt:", mi)
}
func main() {
	i := MyInt(9)
	i.print()
}
```

如果有类型T，方法的接收器为(t T)时我们称为值接收器，该方法称为值方法；方法的接收器为(t *T)时我们称为指针接收器，该方法称为指针方法。

类型 T（或 *T）上的所有方法的集合叫做类型 T（或 *T）的方法集。

>关于接收器的命名
>
>社区约定的接收器命名是类型的一个或两个字母的缩写(像 c 或者 cl 对于 Client)。不要使用泛指的名字像是 me，this 或者 self，也不要使用过度描述的名字，简短即可。

**方法表达式与方法值**

在Go语言中，方法调用的方式如下：如有类型X的变量x，m()是其方法，则方法有效调用方式是x.m()，如果x是指针变量，则x.m()实际上是(&x).m()的简写。所以我们看到指针方法的调用写成x.m()，这其实是一种语法糖。

这里我们了解下Go语言的选择器（selector），如：

```Go
x.f
```

上面代码表示如果x不是包名，则表示是x（或* x）的f（字段或方法）。标识符f（字段或方法）称为选择器(selector)，选择器不能是空白标识符。选择器表达式的类型是f的类型。

选择器f可以表示类型T的字段或方法，或者指嵌入字段T的字段或方法f。遍历到f的嵌入字段的层数被称为其在T中的深度。在T中声明的字段或方法f的深度为零。在T中的嵌入字段A中声明的字段或方法f的深度是A中的f的深度加1。

在Go语言中，我们认为方法的显式接收器(explicit receiver)x是方法x.m()的等效函数X.m()的第一个参数，所以x.m()和X.m(x)是等价的，下面我们看看具体例子：

```Go
package main

import (
	"fmt"
)

type T struct {
	a int
}

func (tv T) Mv(a int) int {
	fmt.Printf("Mv的值是: %d\n", a)
	return a
} // 值方法

func (tp *T) Mp(f float32) float32 {
	fmt.Printf("Mp: %f\n", f)
	return f
} // 指针方法

func main() {
	var t T
	// 下面几种调用方法是等价的
	t.Mv(1)    // 一般调用
	T.Mv(t, 1) // 显式接收器t可以当做为函数的第一个参数
	f0 := t.Mv // 通过选择器（selector）t.Mv将方法值赋值给一个变量 f0
	f0(2)
	T.Mv(t, 3)
	(T).Mv(t, 4)
	f1 := T.Mv // 利用方法表达式(Method Expression) T.Mv 取到函数值
	f1(t, 5)
	f2 := (T).Mv // 利用方法表达式(Method Expression) T.Mv 取到函数值
	f2(t, 6)
}

```

t.Mv(1)和T.Mv(t, 1)效果是一致的，这里显式接收器t可以当做为等效函数T.Mv()的第一个参数。而在Go语言中，我们可以利用选择器，将方法值(Method Value)取到，并可以将其赋值给其它变量。使用 t.Mv，就可以得到 Mv 方法的方法值，而且这个方法值绑定到了显式接收器（实参）t。

```Go
f0 := t.Mv // 通过选择器将方法值t.Mv赋值给一个变量 f0
```

除了使用选择器取到方法值外，还可以使用方法表达式(Method Expression) 取到函数值(Function Value)。方法表达式(Method Expression)产生的是一个函数值(Function Value)而不是方法值(Method Value)。

```Go
f1 := T.Mv // 利用方法表达式(Method Expression) T.Mv 取到函数值
f1(t, 5)
f2 := (T).Mv // 利用方法表达式(Method Expression) T.Mv 取到函数值
f2(t, 6)
```

这个函数值的第一个参数必须是一个接收器：

```Go
f1(t, 5)
f2(t, 6)
```

上面有关选择器，方法表达式，函数值，方法值等概念可以帮助我们更好理解方法，掌握他们可以更好地使用好方法。

在Go语言中不允许方法重载，因为方法是函数，所以对于一个类型只能有唯一一个特定名称的方法。但是如果基于接收器类型，我们可以通过一种变通的方法，达到这个目的：具有同样名字的方法可以在 2 个或多个不同的接收器类型上存在，比如在同一个包里这么做是允许的：

```Go
type MyInt1 int
type MyInt2 int

func (a *MyInt1) Add(b int) int { return 0 }
func (a *MyInt2) Add(b int) int { return 0 }
```

**自定义类型方法与匿名嵌入**

Go语言中类型加上它的方法集等价于面向对象中的类。但在 Go 语言中，类型的代码和绑定在它上面的方法集的代码可以不放置在同一个文件中，它们可以保存在同一个包下的其他源文件中。

下面是在非结构体类型上定义方法的例子：

```Go
type MyInt int

func (m MyInt) print() { // 值方法
	fmt.Println("MyInt:", m)
}
```

注意：类型和作用在它上面定义的方法必须在同一个包里定义，所以基础类型int、float 等上不能直接定义。

类型在其他的，或是非本地的包里定义，在它上面定义方法都会发生错误。

```Go
package main

import (
	"fmt"
)

func (i int) print() { // cannot define new methods on non-local type int
	fmt.Println("Int:", i)
}

func main() {
}

程序编译不通过，错误如下：
cannot define new methods on non-local type int
```

虽然我们不能直接为非同一包下的类型直接定义方法，但我们可以以这个类型（比如：int 或 float）为基础来自定义新类型，然后再为新类型定义方法。

```Go
package main

import (
	"fmt"
)

type MyInt int

func (m MyInt) print() { // 值方法
	fmt.Println("MyInt:", m)
}

func main() {
	myi := MyInt(99)
	myi.print()
}

程序输出：
MyInt: 99
```

MyInt类型由int 为基础自定义的，MyInt定义了一个方法print()。

下面我们再以这个代码为例看看在类型别名下的方法情况，类型别名情况下方法是保留的，但自定义的新类型方法是需要重新定义的，原方法不保留。

如果我们采用类型别名下面程序可正常运行，Go 1.9及以上版本编译通过：

```Go
package main

import (
	"fmt"
)

type MyInt int
type NewInt = MyInt

func (m MyInt) print() { // 值方法
	fmt.Println("MyInt:", m)
}

func main() {
	myi := MyInt(99)
	myi.print()

	Ni := NewInt(myi)
	Ni.print()
}

程序输出：
MyInt: 99
MyInt: 99
```

但上面代码我们稍微修改，把type NewInt = MyInt 改为type NewInt  MyInt 。一个符号“=”去掉使得NewInt 变为新类型，会报程序错误：

```Go
Ni.print undefined (type NewInt has no field or method print)
```

因为Ni 属于新的自定义类型 NewInt, 它没有定义print()方法，需要另外定义这个方法。

我们也可以像下面这样将定义好的类型作为匿名类型嵌入在一个新的结构体中。当然新方法只在这个自定义类型上有效。

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

func (h Human) String() { // 值方法
	fmt.Println("Human")
}

func (s Student) String() { // 值方法
	fmt.Println("Student")
}

func (s Student) Print() { // 值方法
	fmt.Println("Print")
}

func main() {
	stud := Student{Room: 102, Human: Human{"Hawking", "男", 14, "Monitor"}}
	stud.String()
	stud.Human.String()
}

程序输出：
Student
Human
```

## 20.2 函数和方法的区别

方法相对于函数多了接收器，这是他们之间最大的区别。

函数是直接调用，而方法是作用在接收器上，方法需要类型的实例来调用。方法接收器必须有一个显式的名字，这个名字必须在方法中被使用。

在接收器是指针时，方法可以改变接收器的值（或状态），这点函数也可以做到（当参数作为指针传递，即通过引用调用时，函数也可以改变参数的状态）。

在 Go 语言中，（接收器）类型关联的方法不写在类型结构里面，就像类那样；耦合更加宽松；类型和方法之间的关联由接收器来建立。

方法没有和定义的数据类型（结构体）混在一起，方法和数据是正交，而且数据和行为（方法）是相对独立的。

## 20.3 指针方法与值方法

有类型T，方法的接收器为(t T)时我们称为值接收器，该方法称为值方法；方法的接收器为(t *T)时我们称为指针接收器，该方法称为指针方法。

如果想要方法改变接收器的数据，就在接收器的指针上定义该方法；否则，就在普通的值类型上定义方法。这是指针方法和值方法最大的区别。

下面声明一个 T 类型的变量，并调用其方法 M1() 和 M2() 。

```Go
package main

import (
	"fmt"
)

type T struct {
	Name string
}

func (t T) M1() {
	t.Name = "name1"
}

func (t *T) M2() {
	t.Name = "name2"
}
func main() {

	t1 := T{"t1"}

	fmt.Println("M1调用前：", t1.Name)
	t1.M1()
	fmt.Println("M1调用后：", t1.Name)

	fmt.Println("M2调用前：", t1.Name)
	t1.M2()
	fmt.Println("M2调用后：", t1.Name)

}

程序输出：
M1调用前： t1
M1调用后： t1
M2调用前： t1
M2调用后： name2
```
可见，t1.M2()修改了接收器数据。

>分析：
>
>分析：
>由于调用 t1.M1() 时相当于T.M1(t1)，实参和形参都是类型 T。此时在M1()中的t只是t1的值拷贝，所以M1()的修改影响不到t1。
>
>同上， t1.M2() => M2(t1)，这是将 T 类型传给了 *T 类型，Go会取 t1 的地址传进去：M2(&t1)，所以M2()的修改可以影响 t1 。


上面的例子同时也说明了：

```Go
 T 类型的变量可以调用M1()和M2()这两个方法。
```

因为对于类型 T，如果在 *T 上存在方法 M2()，并且 t 是这个类型的变量，那么 t.M2() 会被自动转换为 (&t).M2()。

下面声明一个 *T 类型的变量，并调用方法 M1() 和 M2() 。

```Go
package main

import (
	"fmt"
)

type T struct {
	Name string
}

func (t T) M1() {
	t.Name = "name1"
}

func (t *T) M2() {
	t.Name = "name2"
}
func main() {

	t2 := &T{"t2"}

	fmt.Println("M1调用前：", t2.Name)
	t2.M1()
	fmt.Println("M1调用后：", t2.Name)

	fmt.Println("M2调用前：", t2.Name)
	t2.M2()
	fmt.Println("M2调用后：", t2.Name)

}

程序输出：
M1调用前： t2
M1调用后： t2
M2调用前： t2
M2调用后： name2
```

>分析：
>
>t2.M1() => M1(t2)，t2 是指针类型，取t2的值并拷贝一份传给M1()。
>
>t2.M2() => M2(t2)，都是指针类型，不需要转换。

```Go
*T 类型的变量也可以调用M1()和M2()这两个方法。
```

从上面调用我们可以得知：无论你声明方法的接收器是指针接收器还是值接收器，Go都可以帮你隐式转换为正确的方法使用。

但我们需要记住，值变量只拥有值方法集，而指针变量则同时拥有值方法集和指针方法集。

**接口变量上的指针方法与值方法**

无论是T类型变量还是*T类型变量，都可调用值方法或指针方法。但如果是接口变量呢，那么这两个方法都可以调用吗？

我们添加一个接口看看：

```Go
package main

type T struct {
	Name string
}
type Intf interface {
	M1()
	M2()
}

func (t T) M1() {
	t.Name = "name1"
}

func (t *T) M2() {
	t.Name = "name2"
}
func main() {
	var t1 T = T{"t1"}
	t1.M1()
	t1.M2()

	var t2 Intf = t1
	t2.M1()
	t2.M2()
}
```
编译不通过：

cannot use t1 (type T) as type Intf in assignment:
	T does not implement Intf (M2 method has pointer receiver)

上面代码中我们看到，var t2 Intf 中，t2是Intf接口类型变量，t1是T类型值变量。上面错误信息中已经明确了T没有实现接口Intf，所以不能直接赋值。这是为什么呢？

首先这是Go语言的一种规则，具体如下：

* **规则一：如果使用指针方法来实现一个接口，那么只有指向那个类型的指针才能够实现对应的接口。**
* **规则二：如果使用值方法来实现一个接口，那么那个类型的值和指针都能够实现对应的接口。**

按照上面两条规则的规则一，我们稍微修改下代码：

```Go
package main

type T struct {
	Name string
}
type Intf interface {
	M1()
	M2()
}

func (t T) M1() {
	t.Name = "name1"
}

func (t *T) M2() {
	t.Name = "name2"
}
func main() {

	var t1 T = T{"t1"}
	t1.M1()
	t1.M2()

	var t2 Intf = &t1
	t2.M1()
	t2.M2()
}
```
程序编译通过。

程序编译通过。综合起来看，接口类型的变量（实现了该接口的类型变量）调用方法时，我们需要注意方法的接收器，是不是真正实现了接口。结合接口类型断言，我们做下测试：

```Go
package main

import (
	"fmt"
)

type T struct {
	Name string
}
type Intf interface {
	M1()
	M2()
}

func (t T) M1() {
	t.Name = "name1"
	fmt.Println("M1")
}

func (t *T) M2() {
	t.Name = "name2"
	fmt.Println("M2")
}
func main() {

	var t1 T = T{"t1"}

	// interface{}(t1) 先转为空接口，再使用接口断言
	_, ok1 := interface{}(t1).(Intf)
	fmt.Println("t1 => Intf", ok1)

	_, ok2 := interface{}(t1).(T)
	fmt.Println("t1 => T", ok2)
	t1.M1()
	t1.M2()

	_, ok3 := interface{}(t1).(*T)
	fmt.Println("t1 => *T", ok3)
	t1.M1()
	t1.M2()

	_, ok4 := interface{}(&t1).(Intf)
	fmt.Println("&t1 => Intf", ok4)
	t1.M1()
	t1.M2()

	_, ok5 := interface{}(&t1).(T)
	fmt.Println("&t1 => T", ok5)

	_, ok6 := interface{}(&t1).(*T)
	fmt.Println("&t1 => *T", ok6)
	t1.M1()
	t1.M2()

}


程序输出：
t1 => Intf false
t1 => T true
M1
M2
t1 => *T false
M1
M2
&t1 => Intf true
M1
M2
&t1 => T false
&t1 => *T true
M1
M2
```

执行结果表明，t1 没有实现Intf方法集，不是Intf接口类型；而&t1 则实现了Intf方法集，是Intf接口类型，可以调用相应方法。t1 这个结构体值变量本身则调用值方法或者指针方法都是可以的，这是因为语法糖存在的原因。

按照上面的两条规则，那究竟怎么选择是指针接收器还是值接收器呢？

* 何时使用值类型

（1）如果接收器是一个 map，func 或者 chan，使用值类型（因为它们本身就是引用类型）。
（2）如果接收器是一个 slice，并且方法不执行 reslice 操作，也不重新分配内存给 slice，使用值类型。
（3）如果接收器是一个小的数组或者原生的值类型结构体类型(比如 time.Time 类型)，而且没有可修改的字段和指针，又或者接收器是一个简单地基本类型像是 int 和 string，使用值类型就好了。

值类型的接收器可以减少一定数量的内存垃圾生成，值类型接收器一般会在栈上分配到内存（但也不一定），在没搞明白代码想干什么之前，别为这个原因而选择值类型接收器。

* 何时使用指针类型

（1）如果方法需要修改接收器里的数据，则接收器必须是指针类型。
（2）如果接收器是一个包含了 sync.Mutex 或者类似同步字段的结构体，接收器必须是指针，这样可以避免拷贝。
（3）如果接收器是一个大的结构体或者数组，那么指针类型接收器更有效率。
（4）如果接收器是一个结构体，数组或者 slice，它们中任意一个元素是指针类型而且可能被修改，建议使用指针类型接收器，这样会增加程序的可读性。

最后如果实在还是不知道该使用哪种接收器，那么记住使用指针接收器是最靠谱的。

## 20.4 匿名类型的方法提升

当一个匿名类型被嵌入在结构体中时，匿名类型的可见方法也同样被内嵌，这在效果上等同于外层类型继承了这些方法：将父类型放在子类型中来实现亚型。这个机制提供了一种简单的方式来模拟经典面向对象语言中的子类和继承相关的效果。

当我们嵌入一个匿名类型，这个类型的方法就变成了外部类型的方法，但是当它的方法被调用时，方法的接收器是内部类型(嵌入的匿名类型)，而非外部类型。

```Go
type People struct {
	Age    int
	gender string
	Name   string
}

type OtherPeople struct {
	People
}

func (p People) PeInfo() {
	fmt.Println("People ", p.Name, ": ", p.Age, "岁, 性别:", p.gender)
}
```

因此嵌入类型的名字充当着字段名，同时嵌入类型作为内部类型存在，我们可以使用下面的调用方法：

```Go
OtherPeople.People.PeInfo()
```

这儿我们可以通过类型名称来访问内部类型的字段和方法。然而，这些字段和方法也同样被提升到了外部类型，我们可以直接访问：

```Go
OtherPeople.PeInfo()
```

前面我们看到了嵌入类型的方法提升，在 Go 语言中匿名嵌入类型方法集提升的规则：

给定一个结构体类型 S 和一个命名为 T 的类型，方法提升像下面规定的这样被包含在结构体方法集中：

**简单地说是两条规则：**

**规则一：如果S包含嵌入字段T，则S和\*S的方法集都包括具有接收器T的提升方法。\*S的方法集还包括具有接收器*T的提升方法。**

**规则二：如果S包含嵌入字段*T，则S和\*S的方法集都包括具有接收器T或\*T的提升方法。**

当嵌入一个类型，嵌入类型的接收器为指针的方法将不能被外部类型的值访问。这跟接口规则一致。

注意：以上规则在调用指针方法 t.M() 时会被自动转换为 (&t).M() ，由于这个语法糖，导致我们很容易误解上面的规则不起作用，而实际上规则是有效的，在实际应用中我们可以留意这个问题。

我们通过下面代码验证下：

```Go
package main

import (
	"fmt"
	"reflect"
)

type People struct {
	Age    int
	gender string
	Name   string
}

type OtherPeople struct {
	People
}

type NewPeople People

func (p *NewPeople) PeName(pname string) {
	fmt.Println("pold name:", p.Name)
	p.Name = pname
	fmt.Println("pnew name:", p.Name)
}

func (p NewPeople) PeInfo() {
	fmt.Println("NewPeople ", p.Name, ": ", p.Age, "岁, 性别:", p.gender)
}

func (p *People) PeName(pname string) {
	fmt.Println("old name:", p.Name)
	p.Name = pname
	fmt.Println("new name:", p.Name)
}

func (p People) PeInfo() {
	fmt.Println("People ", p.Name, ": ", p.Age, "岁, 性别:", p.gender)
}

func methodSet(a interface{}) {
	t := reflect.TypeOf(a)	
	fmt.Printf("%T\n", a)
	for i, n := 0, t.NumMethod(); i < n; i++ {
		m := t.Method(i)
		fmt.Println(i, ":", m.Name, m.Type)
	}
}

func main() {
	p := OtherPeople{People{26, "Male", "张三"}}
	p.PeInfo()
	p.PeName("Joke")

	methodSet(p) // T方法提升

	methodSet(&p) // *T和T方法提升

	pp := NewPeople{42, "Male", "李四"}
	pp.PeInfo()
	pp.PeName("Haw")

	methodSet(&pp)
}


程序输出：
People  张三 :  26 岁, 性别: Male
old name: 张三
new name: Joke
main.OtherPeople
0 : PeInfo func(main.OtherPeople)
*main.OtherPeople
0 : PeInfo func(*main.OtherPeople)
1 : PeName func(*main.OtherPeople, string)
NewPeople  李四 :  42 岁, 性别: Male
pold name: 李四
pnew name: Haw
*main.NewPeople
0 : PeInfo func(*main.NewPeople)
1 : PeName func(*main.NewPeople, string)
```

我们可以从上面输出看到，*OtherPeople 下有两个方法PeInfo(),PeName(string)可以调用，而OtherPeople只有一个方法PeInfo()可以调用。

但是在Go中存在一个语法糖：

```Go
	p.PeInfo()
	p.PeName("Joke")

	methodSet(p) // T方法提升
```

虽然P 只有一个方法：PeInfo func(main.OtherPeople)，但我们依然可以调用p.PeName("Joke")。

这里Go自动转为(&p).PeName("Joke")，其调用后结果让我们以为p有两个方法，其实这里p只有一个方法。

有关于内嵌字段方法集的提升，初学者需要好好留意下这个规则。

结合前面的自定义类型赋值接口类型的规则，与内嵌类型的方法集提升规则这两个大规则一定要弄清楚，只有彻底弄清楚这些规则，我们在阅读和写代码时才能做到气定神闲。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第十九章 接口](https://github.com/ffhelicopter/Go42/blob/master/content/42_19_interface.md)

[第二十一章 协程(goroutine)](https://github.com/ffhelicopter/Go42/blob/master/content/42_21_goroutine.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com