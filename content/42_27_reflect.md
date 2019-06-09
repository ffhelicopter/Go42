# 《Go语言四十二章经》第二十七章 反射(reflect)

作者：李骁

## 27.1 反射(reflect)

反射是应用程序检查其所拥有的结构，尤其是类型的一种能。每种语言的反射模型都不同，并且有些语言根本不支持反射。Go语言实现了反射，反射机制就是在运行时动态调用对象的方法和属性，即可从运行时态的示例对象反求其编码阶段的定义，标准库中reflect包提供了相关的功能。在reflect包中，通过reflect.TypeOf()，reflect.ValueOf()分别从类型、值的角度来描述一个Go对象。

```Go
func TypeOf(i interface{}) Type
type Type interface 

func ValueOf(i interface{}) Value
type Value struct 
```

在Go语言的实现中，一个interface类型的变量存储了2个信息, 一个<值，类型>对，<value,type> :

```Go
(value, type)
```

value是实际变量值，type是实际变量的类型。两个简单的函数，reflect.TypeOf 和 reflect.ValueOf，返回被检查对象的类型和值。

例如，x 被定义为：var x float64 = 3.4，那么 reflect.TypeOf(x) 返回 float64，reflect.ValueOf(x) 返回 3.4。实际上，反射是通过检查一个接口的值，变量首先被转换成空接口。这从下面两个函数签名能够很明显的看出来：

```Go
func TypeOf(i interface{}) Type
func ValueOf(i interface{}) Value
```

reflect.Type 和 reflect.Value 都有许多方法用于检查和操作它们。 

Type主要有：
Kind() 将返回一个常量，表示具体类型的底层类型
Elem()方法返回指针、数组、切片、字典、通道的基类型，这个方法要慎用，如果用在其他类型上面会出现panic

Value主要有：
Type() 将返回具体类型所对应的 reflect.Type（静态类型）
Kind() 将返回一个常量，表示具体类型的底层类型

反射可以在运行时检查类型和变量，例如它的大小、方法和 动态 的调用这些方法。这对于没有源代码的包尤其有用。

由于反射是一个强大的工具，但反射对性能有一定的影响，除非有必要，否则应当避免使用或小心使用。下面代码针对int、数组以及结构体分别使用反射机制，其中的差异请看注释。

```Go
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	name string
}

func main() {

	var a int = 50
	v := reflect.ValueOf(a) // 返回Value类型对象，值为50
	t := reflect.TypeOf(a)  // 返回Type类型对象，值为int
	fmt.Println(v, t, v.Type(), t.Kind())

	var b [5]int = [5]int{5, 6, 7, 8}
	fmt.Println(reflect.TypeOf(b), reflect.TypeOf(b).Kind(),reflect.TypeOf(b).Elem()) // [5]int array int

	var Pupil Student
	p := reflect.ValueOf(Pupil) // 使用ValueOf()获取到结构体的Value对象

	fmt.Println(p.Type()) // 输出:Student
	fmt.Println(p.Kind()) // 输出:struct

}
```

在Go语言中，类型包括 static type和concrete type. 简单说 static type是你在编码是看见的类型(如int、string)，concrete type是实际具体的类型，runtime系统看见的类型。

Type()返回的是静态类型，而kind()返回的是具体类型。上面代码中，在int，数组以及结构体三种类型情况中，可以看到kind()，type()返回值的差异。


**通过反射可以修改原对象**

d.CanAddr()方法：判断它是否可被取地址
d.CanSet()方法：判断它是否可被取地址并可被修改

通过一个settable的Value反射对象来访问、修改其对应的变量值：

```Go
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	name string
	Age  int
}

func main() {

	var a int = 50
	v := reflect.ValueOf(a) // 返回Value类型对象，值为50
	t := reflect.TypeOf(a)  // 返回Type类型对象，值为int
	fmt.Println(v, t, v.Type(), t.Kind(), reflect.ValueOf(&a).Elem())
	seta := reflect.ValueOf(&a).Elem() // 这样才能让seta保存a的值
	fmt.Println(seta, seta.CanSet())
	seta.SetInt(1000)
	fmt.Println(seta)

	var b [5]int = [5]int{5, 6, 7, 8}
	fmt.Println(reflect.TypeOf(b), reflect.TypeOf(b).Kind(), reflect.TypeOf(b).Elem())

	var Pupil Student = Student{"joke", 18}
	p := reflect.ValueOf(Pupil) // 使用ValueOf()获取到结构体的Value对象

	fmt.Println(p.Type()) // 输出:Student
	fmt.Println(p.Kind()) // 输出:struct

	setStudent := reflect.ValueOf(&Pupil).Elem()
	//setStudent.Field(0).SetString("Mike") // 未导出字段，不能修改，panic会发生
	setStudent.Field(1).SetInt(19)
	fmt.Println(setStudent)

}
```

虽然反射可以越过Go语言的导出规则的限制读取结构体中未导出的成员，但不能修改这些未导出的成员。因为一个结构体中只有被导出的字段才是可修改的。

在结构体中有tag标签，通过反射可获取结构体成员变量的tag信息。

```Go
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	name string
	Age  int `json:"years"`
}

func main() {
	var Pupil Student = Student{"joke", 18}
	setStudent := reflect.ValueOf(&Pupil).Elem()

	sSAge, _ := setStudent.Type().FieldByName("Age")
	fmt.Println(sSAge.Tag.Get("json")) // years
}

```

```Go
程序输出：
years
```

## 27.2 反射结构体

为了完整说明反射的情况，通过反射一个结构体类型，综合来说明。下面例子较为系统地利用一个结构体，来充分举例说明反射：

```Go
package main

import (
	"fmt"
	"reflect"
)

// 结构体
type ss struct {
	int
	string
	bool
	float64
}

func (s ss) Method1(i int) string  { return "结构体方法1" }
func (s *ss) Method2(i int) string { return "结构体方法2" }

var (
	structValue = ss{ // 结构体
		20, 
		"结构体", 
		false, 
		64.0, 
	}
)

// 复杂类型
var complexTypes = []interface{}{
	structValue, &structValue, // 结构体
	structValue.Method1, structValue.Method2, // 方法
}

func main() {
	// 测试复杂类型
	for i := 0; i < len(complexTypes); i++ {
		PrintInfo(complexTypes[i])
	}
}

func PrintInfo(i interface{}) {
	if i == nil {
		fmt.Println("--------------------")
		fmt.Printf("无效接口值：%v\n", i)
		fmt.Println("--------------------")
		return
	}
	v := reflect.ValueOf(i)
	PrintValue(v)
}

func PrintValue(v reflect.Value) {
	fmt.Println("--------------------")
	// ----- 通用方法 -----
	fmt.Println("String             :", v.String())  // 反射值的字符串形式
	fmt.Println("Type               :", v.Type())    // 反射值的类型
	fmt.Println("Kind               :", v.Kind())    // 反射值的类别
	fmt.Println("CanAddr            :", v.CanAddr()) // 是否可以获取地址
	fmt.Println("CanSet             :", v.CanSet())  // 是否可以修改
	if v.CanAddr() {
		fmt.Println("Addr               :", v.Addr())       // 获取地址
		fmt.Println("UnsafeAddr         :", v.UnsafeAddr()) // 获取自由地址
	}
	// 获取方法数量
	fmt.Println("NumMethod          :", v.NumMethod())
	if v.NumMethod() > 0 {
		// 遍历方法
		i := 0
		for ; i < v.NumMethod()-1; i++ {
			fmt.Printf("    ┣ %v\n", v.Method(i).String())
			//			if i >= 4 { // 只列举 5 个
			//				fmt.Println("    ┗ ...")
			//				break
			//			}
		}
		fmt.Printf("    ┗ %v\n", v.Method(i).String())
		// 通过名称获取方法
		fmt.Println("MethodByName       :", v.MethodByName("String").String())
	}

	switch v.Kind() {
	// 结构体：
	case reflect.Struct:
		fmt.Println("=== 结构体 ===")
		// 获取字段个数
		fmt.Println("NumField           :", v.NumField())
		if v.NumField() > 0 {
			var i int
			// 遍历结构体字段
			for i = 0; i < v.NumField()-1; i++ {
				field := v.Field(i) // 获取结构体字段
				fmt.Printf("    ├ %-8v %v\n", field.Type(), field.String())
			}
			field := v.Field(i) // 获取结构体字段
			fmt.Printf("    └ %-8v %v\n", field.Type(), field.String())
			// 通过名称查找字段
			if v := v.FieldByName("ptr"); v.IsValid() {
				fmt.Println("FieldByName(ptr)   :", v.Type().Name())
			}
			// 通过函数查找字段
			v := v.FieldByNameFunc(func(s string) bool { return len(s) > 3 })
			if v.IsValid() {
				fmt.Println("FieldByNameFunc    :", v.Type().Name())
			}
		}
	}
}
```

```Go
程序输出：
String             : <main.ss Value>
Type               : main.ss
Kind               : struct
CanAddr            : false
CanSet             : false
NumMethod          : 1
    ┗ <func(int) string Value>
MethodByName       : <invalid Value>
=== 结构体 ===
NumField           : 4
    ├ int      <int Value>
    ├ string   结构体
    ├ bool     <bool Value>
    └ float64  <float64 Value>
--------------------
String             : <*main.ss Value>
Type               : *main.ss
Kind               : ptr
CanAddr            : false
CanSet             : false
NumMethod          : 2
    ┣ <func(int) string Value>
    ┗ <func(int) string Value>
MethodByName       : <invalid Value>
--------------------
String             : <func(int) string Value>
Type               : func(int) string
Kind               : func
CanAddr            : false
CanSet             : false
NumMethod          : 0
--------------------
String             : <func(int) string Value>
Type               : func(int) string
Kind               : func
CanAddr            : false
CanSet             : false
NumMethod          : 0

```

细心的读者可能发现了上面代码中的一个有趣的问题，那就是structValue, &structValue的反射结果是不一样的，指针对象在这里有两个方法，而值对象只有一个方法，这是因为Method2()方法是指针方法，在值对象中是不能被反射到的。


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第二十六章 测试](https://github.com/ffhelicopter/Go42/blob/master/content/42_26_testing.md)

[第二十八章 unsafe包](https://github.com/ffhelicopter/Go42/blob/master/content/42_28_unsafe.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
