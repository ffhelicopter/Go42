# 《Go语言四十二章经》第十二章 切片(slice)

作者：李骁

## 12.1 切片(slice)

**切片（slice）**是对数组一个连续片段的引用（该数组我们称之为相关数组，通常是匿名的），所以切片是一个引用类型（和数组不一样）。这个片段可以是整个数组，或者是由起始和终止索引标识的一些项的子集。需要注意的是，终止索引标识的项不包括在切片内。切片提供了一个相关数组的动态窗口（这里有关动态窗口的含义，可参考数据库窗口函数的解释）。

切片是可索引的，并且可以由 len() 函数获取长度。

给定项的切片索引可能比相关数组的相同元素的索引小。和数组不同的是，切片的长度可以在运行时修改，最小为 0 最大为相关数组的长度：切片是一个 长度可变的数组。

切片提供了计算容量的函数 cap() 可以测量切片最长可以达到多少：它等于切片的长度 + 数组除切片之外的长度。如果 s 是一个切片，cap(s) 就是从 s[0] 到数组末尾的数组长度。切片的长度永远不会超过它的容量，所以对于 切片 s 来说该不等式永远成立：0 <= len(s) <= cap(s)。

多个切片如果表示同一个数组的片段，它们可以共享数据；因此一个切片和相关数组的其他切片是共享存储的，相反，不同的数组总是代表不同的存储。数组实际上是切片的构建块。

**优点**

因为切片是引用，所以它们不需要使用额外的内存并且比使用数组更有效率，所以在 Go 代码中切片比数组更常用。

**注意**

绝对不要用指针指向 slice，切片本身已经是一个引用类型，所以它本身就是一个指针!

声明切片的格式是： var identifier []type（不需要说明长度）。

一个切片在未初始化之前默认为 nil，长度为 0。

切片的初始化格式是：

```Go
var slice1 []type = arr1[start:end]
```
这表示 slice1 是由数组 arr1 从 start 索引到 end-1 索引之间的元素构成的子集（切分数组，start:end 被称为 slice 表达式）。

切片也可以用类似数组的方式初始化：

```Go
var x = []int{2, 3, 5, 7, 11}
```
这样就创建了一个长度为 5 的数组并且创建了一个相关切片。

当相关数组还没有定义时，我们可以使用 make() 函数来创建一个切片 同时创建好相关数组：

```Go
var slice1 []type = make([]type, len)
```
也可以简写为 slice1 := make([]type, len)，这里 len 是数组的长度并且也是 slice 的初始长度。

make 的使用方式是：func make([]T, len, cap)，其中 cap 是可选参数。

```Go
v := make([]int, 10, 50)
```
这样分配一个有 50 个 int 值的数组，并且创建了一个长度为 10，容量为 50 的 切片 v，该 切片 指向数组的前 10 个元素。

以上我们列举了三种切片初始化方式，这三种方式都比较常用。

如果从数组或者切片中生成一个新的切片，我们可以使用下面的表达式：
```Go
a[low : high : max]
```
max-low的结果表示容量。

```Go
a := [5]int{1, 2, 3, 4, 5}
t := a[1:3:5]
```

这里t的容量（capacity）是5-1=4 ，长度是2。

## 12.2 切片重组(reslice)

```Go
slice1 := make([]type, start_length, capacity)
```
其中 start_length 作为切片初始长度而 capacity 作为相关数组的长度。

这么做的好处是我们的切片在达到容量上限后可以扩容。改变切片长度的过程称之为切片重组 reslicing，做法如下：slice1 = slice1[0:end]，其中 end 是新的末尾索引（即长度）。

当你重新划分一个slice时，新的slice将引用原有slice的数组。如果你忘了这个行为的话，在你的应用分配大量临时的slice用于创建新的slice来引用原有数据的一小部分时，会导致难以预期的内存使用。

```Go
package main

import "fmt"

func get() []byte {  
    raw := make([]byte, 10000)
    fmt.Println(len(raw), cap(raw), &raw[0]) // prints: 10000 10000 <byte_addr_x>
    return raw[:3]  // 10000个字节实际只需要引用3个，其他空间浪费
}

func main() {  
    data := get()
    fmt.Println(len(data), cap(data), &data[0]) // prints: 3 10000 <byte_addr_x>
}
```
为了避免这个陷阱，你需要从临时的slice中拷贝数据（而不是重新划分slice）。

```Go
package main

import "fmt"

func get() []byte {  
    raw := make([]byte, 10000)
    fmt.Println(len(raw), cap(raw), &raw[0]) // prints: 10000 10000 <byte_addr_x>
    res := make([]byte, 3)
    copy(res, raw[:3]) // 利用copy 函数复制，raw 可被GC释放
    return res
}

func main() {  
    data := get()
    fmt.Println(len(data), cap(data), &data[0]) // prints: 3 3 <byte_addr_y>
}
```
func append(s[]T, x ...T) []T 其中 append 方法将 0 个或多个具有相同类型 s 的元素追加到切片后面并且返回新的切片；追加的元素必须和原切片的元素同类型。如果 s 的容量不足以存储新增元素，append 会分配新的切片来保证已有切片元素和新增元素的存储。因此，返回的切片可能已经指向一个不同的相关数组了。append 方法总是返回成功，除非系统内存耗尽了。

append操作如果导致分配新的切片来保证已有切片元素和新增元素的存储，那么新的slice已经和原来slice没有任何关系，即使修改了数据也不会同步。append操作后，有没有生成新的slice需要看原有slice的容量是否足够，请见下面代码。

## 12.3 陈旧的(Stale)Slices
多个slice可以引用同一个底层数组。比如，当你从一个已有的slice创建一个新的slice时，这就会发生。如果你的应用功能需要这种行为，那么你将需要关注下“走味的”slice。

在某些情况下，在一个slice中添加新的数据，在原有数组无法保持更多新的数据时，将导致分配一个新的数组。而现在其他的slice还指向老的数组（和老的数据）。

```Go
package main

import "fmt"

func main() {
	s1 := []int{1, 2, 3}
	fmt.Println(len(s1), cap(s1), s1) // 输出 3 3 [1 2 3]
	s2 := s1[1:]
	fmt.Println(len(s2), cap(s2), s2) // 输出 2 2 [2 3]
	for i := range s2 {
		s2[i] += 20
	}
	// s2的修改会影响到数组数据，s1输出新数据
	fmt.Println(s1) // 输出 [1 22 23]
	fmt.Println(s2) // 输出 [22 23]

	s2 = append(s2, 4) // append  导致了slice 扩容

	for i := range s2 {
		s2[i] += 10
	}
	// s1 的数据现在是陈旧的老数据，而s2是新数据，他们的底层数组已经不是同一个了。
	fmt.Println(s1) // 输出[1 22 23]
	fmt.Println(s2) // 输出[32 33 14]
}
```

```Go
程序输出：
3 3 [1 2 3]
2 2 [2 3]
[1 22 23]
[22 23]
[1 22 23]
[32 33 14]
```


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>本书《Go语言四十二章经》内容在简书同步地址：  https://www.jianshu.com/nb/29056963
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com