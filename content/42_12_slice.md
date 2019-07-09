# 《Go语言四十二章经》第十二章 切片(slice)

作者：李骁

## 12.1 切片(slice)

**切片（slice）** 是对底层数组一个连续片段的引用，所以切片是一个引用类型。切片提供对该数组中编号的元素序列的访问。未初始化切片的值为nil。

与数组一样，切片是可索引的并且具有长度。切片s的长度可以通过内置函数len() 获取;与数组不同，切片的长度可能在执行期间发生变化。元素可以通过整数索引0到len（s）-1来寻址。我们可以把切片看成是一个长度可变的数组。

切片提供了计算容量的函数 cap() ，可以测量切片最大长度。切片的长度永远不会超过它的容量，所以对于切片 s 来说，这个不等式永远成立：0 <= len(s) <= cap(s)。

一旦初始化，切片始终与保存其元素的基础数组相关联。因此，切片会和与其拥有同一基础数组的其他切片共享存储;相比之下，不同的数组总是代表不同的存储。

切片下面的数组可以延伸超过切片的末端。容量是切片长度与切片之外的数组长度的总和。

使用内置函数make()可以给切片初始化，该函数指定切片类型和指定长度和可选容量的参数。

切片与数组相比较：

**优点**

因为切片是引用，所以它们不需要使用额外的内存并且比使用数组更有效率，所以在 Go 代码中切片比数组更常用。


声明切片的格式是： var identifier []type（不需要说明长度）。一个切片在未初始化之前默认为 nil，长度为 0。

切片的初始化格式是：

```Go
var slice1 []type = arr1[start:end]
```

这表示 slice1 是由数组 arr1 从 start 索引到 end-1 索引之间的元素构成的子集（切分数组，start:end 被称为切片表达式）。

切片也可以用类似数组的方式初始化：

```Go
var x = []int{2, 3, 5, 7, 11}
```

这样就创建了一个长度为 5 的数组并且创建了一个相关切片。

当相关数组还没有定义时，我们可以使用 make() 函数来创建一个切片，同时创建好相关数组：


```Go
var slice1 []type = make([]type, len,cap)
```

也可以简写为 slice1 := make([]type, len)，这里 len 是数组的长度并且也是切片的初始长度。cap是容量，其中cap是可选参数。

```Go
v := make([]int, 10, 50)
```

这样分配一个有 50 个int值的数组，并且创建了一个长度为10，容量为50的切片 v，该切片指向数组的前 10 个元素。

以上我们列举了三种切片初始化方式，这三种方式都比较常用。

如果从数组或者切片中生成一个新的切片，我们可以使用下面的表达式：

a[low : high : max]     max-low的结果表示容量，high-low的结果表示长度。

```Go
a := [5]int{1, 2, 3, 4, 5}
t := a[1:3:5]
```

这里t的容量（capacity）是5-1=4 ，长度是2。

如果切片取值时索引值大于长度会导致panic错误发生，即使容量远远大于长度也没有用，如下面代码所示：

```Go
package main

import "fmt"

func main() {
	sli := make([]int, 5, 10)
	fmt.Printf("切片sli长度和容量：%d, %d\n", len(sli), cap(sli))
	fmt.Println(sli)
	newsli := sli[:cap(sli)]
	fmt.Println(newsli)

	var x = []int{2, 3, 5, 7, 11}
	fmt.Printf("切片x长度和容量：%d, %d\n", len(x), cap(x))

	a := [5]int{1, 2, 3, 4, 5}
	t := a[1:3:5] // a[low : high : max]  max-low的结果表示容量  high-low为长度
	fmt.Printf("切片t长度和容量：%d, %d\n", len(t), cap(t))

	// fmt.Println(t[2]) // panic ，索引不能超过切片的长度
}

程序输出：
切片sli长度和容量：5, 10
[0 0 0 0 0]
[0 0 0 0 0 0 0 0 0 0]
切片x长度和容量：5, 5
切片t长度和容量：2, 4
```

## 12.2 切片重组(reslice)

```Go
slice1 := make([]type, start_length, capacity)
```

通过改变切片长度得到新切片的过程称之为切片重组 reslicing，做法如下：slice1 = slice1[0:end]，其中 end 是新的末尾索引（即长度）。

当我们在一个切片基础上重新划分一个切片时，新的切片会继续引用原有切片的数组。如果你忘了这个行为的话，在你的应用分配大量临时的切片用于创建新的切片来引用原有数据的一小部分时，会导致难以预期的内存使用。

```Go
package main

import "fmt"

func get() []byte {  
    raw := make([]byte, 10000)
    fmt.Println(len(raw), cap(raw), &raw[0]) // 显示: 10000 10000 数组首字节地址
    return raw[:3]  // 10000个字节实际只需要引用3个，其他空间浪费
}

func main() {  
    data := get()
    fmt.Println(len(data), cap(data), &data[0]) // 显示: 3 10000 数组首字节地址
}
```

为了避免这个陷阱，我们需要从临时的切片中使用内置函数copy()，拷贝数据（而不是重新划分切片）到新切片。

```Go
package main

import "fmt"

func get() []byte {
	raw := make([]byte, 10000)
	fmt.Println(len(raw), cap(raw), &raw[0]) // 显示: 10000 10000 数组首字节地址
	res := make([]byte, 3)
	copy(res, raw[:3]) // 利用copy 函数复制，raw 可被GC释放
	return res
}

func main() {
	data := get()
	fmt.Println(len(data), cap(data), &data[0]) // 显示: 3 3 数组首字节地址
}

程序输出：
10000 10000 0xc000086000
3 3 0xc000050098
```

append()内置函数：

```Go
func append(s S, x ...T) S  // T是S元素类型
```

append()函数将 0 个或多个具有相同类型S的元素追加到切片s后面并且返回新的切片；追加的元素必须和原切片的元素同类型。如果s的容量不足以存储新增元素，append()会分配新的切片来保证已有切片元素和新增元素的存储。

因此，append()函数返回的切片可能已经指向一个不同的相关数组了。append()函数总是返回成功，除非系统内存耗尽了。

```Go
s0 := []int{0, 0}
s1 := append(s0, 2)                // append 单个元素     s1 == []int{0, 0, 2}
s2 := append(s1, 3, 5, 7)          // append 多个元素    s2 == []int{0, 0, 2, 3, 5, 7}
s3 := append(s2, s0...)            // append 一个切片     s3 == []int{0, 0, 2, 3, 5, 7, 0, 0}
s4 := append(s3[3:6], s3[2:]...)   // append 切片片段    s4 == []int{3, 5, 7, 2, 3, 5, 7, 0, 0}
```

append()函数操作如果导致分配新的切片来保证已有切片元素和新增元素的存储，也就是返回的切片可能已经指向一个不同的相关数组了，那么新的切片已经和原来切片没有任何关系，即使修改了数据也不会同步。

append()函数操作后，有没有生成新的切片需要看原有切片的容量是否足够。

## 12.3 陈旧的切片(Stale Slices)

多个切片可以引用同一个底层数组。在某些情况下，在一个切片中添加新的数据，在原有数组无法保持更多新的数据时，将导致分配一个新的数组。而现在其他的切片还指向老的数组（和老的数据）。

上一节我们也说了：append()函数操作后，有没有生成新的切片需要看原有切片的容量是否足够。

下面，我们看看这个过程是怎么产生的：

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

	s2 = append(s2, 4) // append  s2容量为2，这个操作导致了切片 s2扩容，会生成新的底层数组。

	for i := range s2 {
		s2[i] += 10
	}
	// s1 的数据现在是老数据，而s2扩容了，复制数据到了新数组，他们的底层数组已经不是同一个了。
	fmt.Println(len(s1), cap(s1), s1) // 输出3 3 [1 22 23]
	fmt.Println(len(s2), cap(s2), s2) // 输出3 4 [32 33 14]
}


程序输出：
3 3 [1 2 3]
2 2 [2 3]
[1 22 23]
[22 23]
3 3 [1 22 23]
3 4 [32 33 14]
```


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第十一章 数组(Array)](https://github.com/ffhelicopter/Go42/blob/master/content/42_11_array.md)

[第十三章 字典(Map)](https://github.com/ffhelicopter/Go42/blob/master/content/42_13_map.md)




>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com