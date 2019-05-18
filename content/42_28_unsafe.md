# 《Go语言四十二章经》第二十八章 unsafe包

作者：李骁

## 28.1 unsafe 包

```Go
func Alignof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Sizeof(x ArbitraryType) uintptr
type ArbitraryType int
type Pointer *ArbitraryType
```
在unsafe包中，只提供了3个函数，两个类型。就这么少的量，却有着超级强悍的功能。一般我们在C语言中通过指针，在知道变量在内存中占用的字节数情况下，就可以通过指针加偏移量的操作，直接在地址中，修改，访问变量的值。在Go 语言中不支持指针运算，那怎么办呢？其实通过unsafe包，我们可以完成类似的操作。

ArbitraryType 是以int为基础定义的一个新类型，但是Go 语言unsafe包中，对ArbitraryType赋予了特殊的意义，通常，我们把interface{}看作是任意类型，那么ArbitraryType这个类型，在Go 语言系统中，比interface{}还要随意。

Pointer 是ArbitraryType指针类型为基础的新类型，在Go 语言系统中，可以把Pointer类型，理解成任何指针的亲爹。

Go 语言的指针类型长度与int类型长度，在内存中占用的字节数是一样的。ArbitraryType类型的变量也可以是指针。

```Go
func Alignof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Sizeof(x ArbitraryType) uintptr
```

通过分析发现，这三个函数的参数均是ArbitraryType类型。
1. Alignof返回变量对齐字节数量
2. Offsetof返回变量指定属性的偏移量，所以如果变量是一个struct类型，不能直接将这个struct类型的变量当作参数，只能将这个struct类型变量的属性当作参数。
3. Sizeof 返回变量在内存中占用的字节数，切记，如果是slice，则不会返回这个slice在内存中的实际占用长度。

unsafe中，通过ArbitraryType 、Pointer 这两个类型，可以将其他类型都转换过来，然后通过这三个函数，分别能取长度，偏移量，对齐字节数，就可以在内存地址映射中，来回游走。

## 28.2 指针运算

uintptr这个基础类型，在Go 语言中，字节长度是与int一致。通常Pointer不能参与指针运算，比如你要在某个指针地址上加上一个偏移量，Pointer是不能做这个运算的，那么谁可以呢？这里要靠uintptr类型了，只有将Pointer类型先转换成uintptr类型，做完地址加减法运算后，再转换成Pointer类型，通过*操作达到取值、修改值的目的。

unsafe.Pointer其实就是类似C的void *，在Go 语言中是用于各种指针相互转换的桥梁，也即是通用指针。它可以让任意类型的指针实现相互转换，也可以将任意类型的指针转换为 uintptr 进行指针运算。

uintptr是Go 语言的内置类型，是能存储指针的整型， uintptr 的底层类型是int，它和unsafe.Pointer可相互转换。

uintptr和unsafe.Pointer的区别就是：

* unsafe.Pointer只是单纯的通用指针类型，用于转换不同类型指针，它不可以参与指针运算；

* 而uintptr是用于指针运算的，GC 不把 uintptr 当指针，也就是说 uintptr 无法持有对象， uintptr 类型的目标会被回收；

* unsafe.Pointer 可以和 普通指针 进行相互转换；

* unsafe.Pointer 可以和 uintptr 进行相互转换。
 
Go 语言的unsafe包很强大，基本上很少会去用它。它可以像C一样去操作内存，但由于Go 语言不支持直接进行指针运算，所以用起来稍显麻烦。

uintptr和intptr是无符号和有符号的指针类型，并且确保在64位平台上是8个字节，在32位平台上是4个字节，uintptr主要用于Go 语言中的指针运算。

通过unsafe包来实现对V的成员i和j赋值，然后通过GetI()和GetJ()来打印观察输出结果。

以下是main.go源代码：

```Go
package main

import (
	"fmt"
	"unsafe"
)

type V struct {
	i int32
	j int64
}

func (v V) GetI() {
	fmt.Printf("i=%d\n", v.i)
}
func (v V) GetJ() {
	fmt.Printf("j=%d\n", v.j)
}

func main() {
	// 定义指针类型变量
	var v *V = &V{199, 299}

	// 取得v的指针并转为*int32的值，对应结构体的i。
	var i *int32 = (*int32)(unsafe.Pointer(v))

	fmt.Println("指针地址：", i)
	fmt.Println("指针uintptr值:", uintptr(unsafe.Pointer(i)))
	*i = int32(98)

	// 根据v的基准地址加上偏移量进行指针运算，运算后的值为j的地址，使用unsafe.Pointer转为指针
	var j *int64 = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + uintptr(unsafe.Sizeof(int64(0)))))

	*j = int64(763)

	v.GetI()
	v.GetJ()
}
```

```Go
指针地址： 0xc00000c180
指针uintptr值: 824633770368
i=98
j=763
```

要修改struct字段的值，需要提前知道结构体V的成员布局，然后根据字段计算偏移量，以及考虑对齐值，最后通过指针运算得到成员指针，利用指针达到修改成员值得目的。由于结构体的成员在内存中的分配是一段连续的内存，因此结构体中第一个成员的地址就是这个结构体的地址，我们也可以认为是相对于这个结构体偏移了0。相同的，这个结构体中的任一成员都可以相对于这个结构体的偏移来计算出它在内存中的绝对地址。

具体来讲解下main方法的实现：


```Go
var v *V = &V{199, 299}
```

通过&来分配一段内存(并按类型初始化)，返回一个指针。所以v就是类型为V的一个指针。和new函数的作用类似。

```Go
var i *int32 = (*int32)(unsafe.Pointer(v))
```

将指针v转成通用指针，再转成int32指针类型。这里就看到了unsafe.Pointer的作用了，您不能直接将v转成int32类型的指针，那样将会panic，但是unsafe.Pointer是可以转为任何指针。刚才说了v的地址其实就是它的第一个成员的地址，所以这个i就很显然指向了v的成员i，通过给i赋值就相当于给v.i赋值了，但是别忘了i只是个指针，要赋值得解引用。

```Go
*i = int32(98)
```

现在已经成功的改变了v的私有成员i的值。

但是对于v.j来说，怎么来得到它在内存中的地址呢？其实我们可以获取它相对于v的偏移量(unsafe.Sizeof可以为我们做这个事)，但上面的代码并没有这样去实现。各位别急，一步步来。


```Go
var j *int64 = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + uintptr(unsafe.Sizeof(int64(0)))))
```

其实我们已经知道v是有两个成员的，包括i和j，并且在定义中，i位于j的前面，而i是int32类型，也就是说i占4个字节。所以j是相对于v偏移了4个字节。您可以用uintptr(4)或uintptr(unsafe.Sizeof(int64(0)))来做这个事。unsafe.Sizeof方法用来得到一个值应该占用多少个字节空间。注意这里跟C的用法不一样，C是直接传入类型，而Go 语言是传入值。

之所以转成uintptr类型是因为需要做指针运算。v的地址加上j相对于v的偏移地址，也就得到了v.j在内存中的绝对地址，然后通过unsafe.Pointer转为指针，别忘了j的类型是int64，所以现在的j就是一个指向v.j的指针，接下来给它赋值：

```Go
*j = int64(763)
```

另外，我们可以看到两种地址表示上的差异：

```Go
指针地址： 0xc00000c180
指针uintptr值: 824633770368
```

上面结构体V中，定义了2个成员属性，如果我们定义一个byte类型的成员属性。我们来看下它的输出：

```Go
package main

import (
	"fmt"
	"unsafe"
)

type V struct {
	b byte
	i int32
	j int64
}

func (v V) GetI() {
	fmt.Printf("i=%d\n", v.i)
}
func (v V) GetJ() {
	fmt.Printf("j=%d\n", v.j)
}

func main() {
	// 定义指针类型变量
	var v *V = new(V)

	// v的长度
	fmt.Printf("size=%d\n", unsafe.Sizeof(*v))
	// 取得v的指针考虑对齐值计算偏移量，然后转为*int32的值，对应结构体的i。
	var i *int32 = (*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + uintptr(4*unsafe.Sizeof(byte(0)))))

	fmt.Println("指针地址：", i)
	fmt.Println("指针uintptr值:", uintptr(unsafe.Pointer(i)))
	*i = int32(98)

	// 根据v的基准地址加上偏移量进行指针运算，运算后的值为j的地址，使用unsafe.Pointer转为指针
	var j *int64 = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + uintptr(unsafe.Sizeof(int64(0)))))

	*j = int64(763)
	fmt.Println("指针uintptr值:", uintptr(unsafe.Pointer(&v.b)))
	fmt.Println("指针uintptr值:", uintptr(unsafe.Pointer(&v.i)))
	fmt.Println("指针uintptr值:", uintptr(unsafe.Pointer(&v.j)))
	v.GetI()
	v.GetJ()
}

```

```Go
程序输出：
size=16
指针地址： 0xc000050084
指针uintptr值: 824634048644
指针uintptr值: 824634048640
指针uintptr值: 824634048644
指针uintptr值: 824634048648
i=98
j=763
```

新结构体的长度为size=16，好像跟我们想像的不一致。我们计算一下：b是byte类型，占1个字节；i是int32类型，占4个字节；j是int64类型，占8个字节，1+4+8=13。这是怎么回事呢？

这是因为发生了对齐。在struct中，它的对齐值是它的成员中的最大对齐值。

每个成员类型都有它的对齐值，可以用unsafe.Alignof方法来计算，比如unsafe.Alignof(v.b)就可以得到b的对齐值为1 。但这个对齐值是其值类型的长度或引用的地址长度（32位或者64位），和其在结构体中的size不是简单相加的问题。经过在64位机器上测试，发现地址（uintptr）如下：

```Go
unsafe.Pointer(b): %s 824634048640
unsafe.Pointer(i): %s 824634048644
unsafe.Pointer(j): %s 824634048648
```

可以初步推断，也经过测试验证，取i值使用uintptr(4*unsafe.Sizeof(byte(0)))是准确的。至于size其实也和对齐值有关，也不是简单相加每个字段的长度。

unsafe.Offsetof 可以在实际中使用，如果改变私有的字段，需要程序员认真考虑后，按照上面的方法仔细确认好对齐值再进行操作。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第二十七章 反射(reflect)](https://github.com/ffhelicopter/Go42/blob/master/content/42_27_reflect.md)

[第二十九章 排序(sort)](https://github.com/ffhelicopter/Go42/blob/master/content/42_29_sort.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com