# <center>《Go语言四十二章经》第二十章 方法</center>

作者：李骁

在前面我们讲了结构（struct）和接口（interface），但关于这两种类型中非常重要的的方法以及方法调用一直没有具体讲解。那么在这一章里，我们来仔细看看方法有那些奇妙之处呢？

## 20.1 方法的定义

在 Go 语言中，结构体就像是类的一种简化形式，那么面向对象程序员可能会问：类的方法在哪里呢？在 Go 中有一个概念，它和方法有着同样的名字，并且大体上意思相近：

Go 方法是作用在接收器（receiver）上的一个函数，接收器是某种类型的变量。因此方法是一种特殊类型的函数。

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
上面代码中，我们定义了结构体 A ，注意f()就是 A 的方法，(a A)表示接收器。

a 是 receiver 的实例，f()是它的方法名，那么方法调用遵循传统的 object.name 选择器符号：a.f()。

如果 recv 一个指针，Go 会自动解引用。如果方法不需要使用 recv 的值，可以用 _ 替换它，比如：

```Go
func (_ receiver_type) methodName(parameter_list) (return_value_list) { ... }
```

* 接收器类型可以是（几乎）任何类型，不仅仅是结构体类型：任何类型都可以有方法，甚至可以是函数类型，可以是 int、bool、string 或数组的别名类型。

```Go
package main

import (
	"fmt"
)

type MyInt int

func (m MyInt) p() {
	fmt.Println("Now", m)
}

func main() {
	var pp MyInt = 8
	pp.p()
}
```

```Go
程序输出：
Now 8
```

* 接收器不能是一个接口类型，因为接口是一个抽象定义，但是方法却是具体实现；如果这样做会引发一个编译错误：invalid receiver type…。

* 接收器不能是一个指针类型，但是它可以是任何其他允许类型的指针。

>关于接收器的命名
>
>社区约定的接收器命名是类型的一个或两个字母的缩写(像 c 或者 cl 对于 Client)。不要使用泛指的名字像是 me，this 或者 self，也不要使用过度描述的名字，最后，如果你在一个地方使用了 c，那么就不要在别的地方使用 cl。

一个类型加上它的方法等价于面向对象中的一个类。一个重要的区别是：在 Go 中，类型的代码和绑定在它上面的方法的代码可以不放置在一起，它们可以存在在不同的源文件，唯一的要求是：它们必须是同一个包的。

类型 T（或 *T）上的所有方法的集合叫做类型 T（或 *T）的方法集。

因为方法是函数，所以同样的，不允许方法重载，即对于一个类型只能有一个给定名称的方法。但是如果基于接收器类型，是有重载的：具有同样名字的方法可以在 2 个或多个不同的接收器类型上存在，比如在同一个包里这么做是允许的：

```Go
func (a *denseMatrix) Add(b Matrix) Matrix
func (a *sparseMatrix) Add(b Matrix) Matrix
```
下面是非结构体类型上方法的例子：

```Go
type IntVector []int

func (v IntVector) Sum() (s int) {
    for _, x := range v {
        s += x
    }
    return
}
```
**类型和作用在它上面定义的方法必须在同一个包里定义，这就是为什么不能在 int、float 或类似这些的类型上定义方法。**

类型在其他的，或是非本地的包里定义，在它上面定义方法都会得到和上面同样的错误。

* 但是有一个间接的方式：可以先定义该类型（比如：int 或 float）的新的自定义类型，然后再为自定义类型定义方法。

```Go
package main

import (
	"fmt"
)

type MyInt int
type HeInt MyInt

func (m MyInt) p() {
	fmt.Println("Now", m)
}

func main() {
	var pp MyInt = 8
	pp.p()  

hh := HeInt(pp)
	hh.p()
}
```
程序运行结果：hh.p undefined (type HeInt has no field or method p)
因为hh 属于新的自定义类型 HeInt , 它没有定义p()方法，需要另外定义这个方法。

如果我们采用别名，Go 1.9上编译通过：

```Go
package main

import (
	"fmt"
)

type MyInt int
type HeInt = MyInt

func (m MyInt) p() {
	fmt.Println("Now", m)
}

func main() {
	var pp MyInt = 8
	pp.p()  

hh := HeInt(pp)
	hh.p()
}
```

```Go
程序输出：
Now 8
Now 8
```

* 或者像下面这样将它作为匿名类型嵌入在一个新的结构体中。当然方法只在这个自定义类型上有效。

```Go
package main

import (
    "fmt"
    "time"
)
type myTime struct {
time.Time 
}

func (t myTime) first3Chars() string {
    return t.Time.String()[0:3]
}
func main() {
    m := myTime{time.Now()}
    // 调用匿名Time上的String方法
    fmt.Println("Full time now:", m.String())
    // 调用myTime.first3Chars
    fmt.Println("First 3 chars:", m.first3Chars())
}
```
```Go
程序输出：

Full time now: 2018-08-28 20:36:47.1135231 +0800 CST m=+0.002990901
First 3 chars: 201
```

## 20.2 函数和方法的区别

函数将变量作为参数：Function1(recv)

方法在变量上被调用：recv.Method1()

在接收器是指针时，方法可以改变接收器的值（或状态），这点函数也可以做到（当参数作为指针传递，即通过引用调用时，函数也可以改变参数的状态）。

接收器必须有一个显式的名字，这个名字必须在方法中被使用。

receiver_type 叫做 （接收器）基本类型，这个类型必须在和方法同样的包中被声明。

在 Go 中，（接收器）类型关联的方法不写在类型结构里面，就像类那样；耦合更加宽松；类型和方法之间的关联由接收器来建立。

方法没有和数据定义（结构体）混在一起：

* 它们是正交的类型；
* 表示（数据）和行为（方法）是独立的。

## 20.3 指针或值方法

鉴于性能的原因，recv 最常见的是一个指向 receiver_type 的指针（因为我们不想要一个实例的拷贝，如果按值调用的话就会是这样），特别是在 receiver 类型是结构体时，就更是如此了。

>如果想要方法改变接收器的数据，就在接收器的指针类型上定义该方法；否则，就在普通的值类型上定义方法；分别叫做指针方法，值方法。

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
```
```Go
程序输出：
M1调用前： t1
M1调用后： t1
M2调用前： t1
M2调用后： name2
```
可见，t1.M2()修改了接收器数据。

>分析：
>
>我们姑且认为调用 t1.M1() 时相当于 M1(t1) ，实参和行参都是类型 T，可以接受。此时在M1()中的t只是t1的值拷贝，所以M1()的修改影响不到t1。
>
>同上， t1.M2() => M2(t1)，这是将 T 类型传给了 *T 类型，go可能会取 t1 的地址传进去： M2(&t1)。所以 M2() 的修改可以影响 t1 。

上面的例子同时也说明了：

```Go
 T 类型的变量这两个方法都是拥有的。
```
因为对于类型 T，如果在 *T 上存在方法 Meth()，并且 t 是这个类型的变量，那么 t.Meth() 会被自动转换为 (&t).Meth()。


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
```
```Go
程序输出：
M1调用前： t2
M1调用后： t2
M2调用前： t2
M2调用后： name2
```
>分析：
>
>t2.M1() => M1(t2)， t2 是指针类型， 取 t2 的值并拷贝一份传给 M1。
>
>t2.M2() => M2(t2)，都是指针类型，不需要转换。

```Go
*T 类型的变量也是拥有这两个方法的。
```

从上面我们可以得知：无论你声明方法的接收器是指针接收器还是值接收器，Go都可以帮你隐式转换为正确的值供方法使用。

无论是T类型变量还是*T类型变量，都拥有值方法或指针方法。但如果是接口变量呢，那么这两个方法都可以调用吗？

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


综合起来，接口类型的变量（实现了该接口）调用方法时，我们需要注意方法的接收器，是不是真正实现了接口。结合接口类型断言，我们做下测试：

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

// interface{}(t1) 先转为空接口（泛型），再使用接口断言
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

	T, ok4 := interface{}(&t1).(Intf)
	fmt.Println("&t1 => Intf", ok4)
	t.M1()
	t.M2()

	_, ok5 := interface{}(&t1).(T)
	fmt.Println("&t1 => T", ok5)

	_, ok6 := interface{}(&t1).(*T)
	fmt.Println("&t1 => *T", ok6)
	t1.M1()
	t1.M2()

}
```
```Go
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

执行结果表明，t1 没有实现Intf方法集，不是Intf接口类型；而&t1 则实现了Intf方法集，是Intf接口类型，可以调用相应方法。t1 这个结构体值变量本身则调用值方法或者指针方法都是可以的。

按照上面的两条规则，那究竟怎么选择是指针接收器还是值接收器呢？

* 何时使用值类型

1.如果接收器是一个 map，func 或者 chan，使用值类型(因为它们本身就是引用类型)。

2.如果接收器是一个 slice，并且方法不执行 reslice 操作，也不重新分配内存给 slice，使用值类型。

3.如果接收器是一个小的数组或者原生的值类型结构体类型(比如 time.Time 类型)，而且没有可修改的字段和指针，又或者接收器是一个简单地基本类型像是 int 和 string，使用值类型就好了。

一个值类型的接收器可以减少一定数量的垃圾生成，如果一个值被传入一个值类型接收器的方法，一个栈上的拷贝会替代在堆上分配内存(但不是保证一定成功)，所以在没搞明白代码想干什么之前，别因为这个原因而选择值类型接收器。

* 何时使用指针类型

1.如果方法需要修改接收器，接收器必须是指针类型。

2.如果接收器是一个包含了 sync.Mutex 或者类似同步字段的结构体，接收器必须是指针，这样可以避免拷贝。

3.如果接收器是一个大的结构体或者数组，那么指针类型接收器更有效率。(多大算大呢？假设把接收器的所有元素作为参数传给方法，如果你觉得参数有点多，那么它就是大)。

4.从此方法中并发的调用函数和方法时，接收器可以被修改吗？一个值类型的接收器当方法调用时会创建一份拷贝，所以外部的修改不能作用到这个接收器上。如果修改必须被原始的接收器可见，那么接收器必须是指针类型。

5.如果接收器是一个结构体，数组或者 slice，它们中任意一个元素是指针类型而且可能被修改，建议使用指针类型接收器，这样会增加程序的可读性

**当你看完这个还是有疑虑，还是不知道该使用哪种接收器，那么记住使用指针接收器。**

## 20.4 内嵌类型的方法提升

当一个匿名类型被内嵌在结构体中时，匿名类型的可见方法也同样被内嵌，这在效果上等同于外层类型继承了这些方法：将父类型放在子类型中来实现亚型。这个机制提供了一种简单的方式来模拟经典面向对象语言中的子类和继承相关的效果。

当我们嵌入一个类型，这个类型的方法就变成了外部类型的方法，但是当它被调用时，方法的接收器是内部类型(嵌入类型)，而非外部类型。
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
这儿我们通过类型名称来访问内部类型的字段和方法。然而，这些字段和方法也同样被提升到了外部类型：
```Go
OtherPeople.PeInfo()
```
下面是 Go 语言中内嵌类型方法集提升的规则：

给定一个结构体类型 S 和一个命名为 T 的类型，方法提升像下面规定的这样被包含在结构体方法集中：

* **如果 S 包含一个匿名字段 T，S 和 \*S 的方法集都包含接收器为 T 的方法提升**

    这条规则说的是当我们嵌入一个类型，嵌入类型的接收器为值类型的方法将被提升，可以被外部类型的值和指针调用。

* **如果 S 包含一个匿名字段 T， \*S 类型的方法集包含接收器为 \*T 的方法提升**

    这条规则说的是当我们嵌入一个类型，可以被外部类型的指针调用的方法集只有嵌入类型的接收器为指针类型的方法集，也就是说，当外部类型使用指针调用内部类型的方法时，只有接收器为指针类型的内部类型方法集将被提升。

* **如果 S 包含一个匿名字段 \*T，S 和 \*S 的方法集都包含接收器为 T 或者 \*T 的方法提升**

    这条规则说的是当我们嵌入一个类型的指针，嵌入类型的接收器为值类型或指针类型的方法将被提升，可以被外部类型的值或者指针调用。

这就是语言规范里方法提升中仅有的三条规则，根据这个推导出一条规则：

* **如果 S 包含一个匿名字段 T，S 的方法集不包含接收器为 \*T 的方法提升。**

这条规则说的是当我们嵌入一个类型，嵌入类型的接收器为指针的方法将不能被外部类型的值访问。这也是跟我们陈述的接口规则一致。

简单地说也是两条规则：

```Go
规则一：如果S包含嵌入字段T，则S和*S的方法集都包括具有接收方T的提升方法。*S的方法集还包括具有接收方*T的提升方法。

规则二：如果S包含嵌入字段*T，则S和*S的方法集都包括具有接收器T或*T的提升方法。
```
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
	for i, n := 0, t.NumMethod(); i < n; i++ {
		m := t.Method(i)
		fmt.Println(m.Name, m.Type)
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
```
```Go
程序输出：

People  张三 :  26 岁, 性别: Male
old name: 张三
new name: Joke
PeInfo func(main.OtherPeople)
PeInfo func(*main.OtherPeople)
PeName func(*main.OtherPeople, string)
NewPeople  李四 :  42 岁, 性别: Male
pold name: 李四
pnew name: Haw
PeInfo func(*main.NewPeople)
PeName func(*main.NewPeople, string)
```
我们可以从上面输出看到，*OtherPeople 下有两个方法，而OtherPeople只有一个方法。

但是在Go中存在一个语法糖，比如上面代码：
```Go
	p.PeInfo()
	p.PeName("Joke")

	methodSet(p) // T方法提升
```
虽然P 只有一个方法：PeInfo func(main.OtherPeople)，但我们依然可以调用p.PeName("Joke")。

这里Go自动转为(&p).PeName("Joke")，让我们以为p有两个方法，其实这里p只有一个方法。这就是上面所谓的内嵌方法集提升，初学者留意下这个规则。

结合前面的自定义类型赋值接口类型的规则，与内嵌类型的方法集提升规则这两个大规则一定要弄清楚，只有彻底弄清楚这些规则，我们在阅读和写代码时才能做到气闲神定。