# 《Go语言四十二章经》第三十五章 模板

作者：李骁

标准库fmt包中，Printf()等方法可以做到输出格式化，当然，对于简单的例子来说足够了，但是我们有时候还是需要复杂的输出格式，甚至需要将格式化代码分离开来。这时，可以使用text/template和html/template。

Go 官方库提供了两个模板库： text/template 和 html/template 。这两个库类似，当需要输出html格式的代码时需要使用 html/template。

## 35.1 text/template

所谓模板引擎，则将模板和数据进行渲染的输出格式化后的字符程序。对于Go，执行这个流程大概需要三步。

* 创建模板对象
* 加载模板
* 执行渲染模板

其中最后一步就是把加载的字符和数据进行格式化。

```Go
package main

import (
	"log"
	"os"
	"text/template"
)

// printf "%6.2f" 表示6位宽度2位精度
const templ = ` 
{{range .}}----------------------------------------
Name:   {{.Name}}
Price:  {{.Price | printf "%6.2f"}}
{{end}}`

var report = template.Must(template.New("report").Parse(templ))

type Book struct {
	Name  string
	Price float64
}

func main() {
	Data := []Book{ {"《三国演义》", 19.82}, {"《儒林外史》", 99.09}, {"《史记》", 26.89} }
	if err := report.Execute(os.Stdout, Data); err != nil {
		log.Fatal(err)
	}
}
```

```Go
程序输出：
----------------------------------------
Name:   《三国演义》
Price:   19.82
----------------------------------------
Name:   《儒林外史》
Price:   99.09
----------------------------------------
Name:   《史记》
Price:   26.89
```

如果把模板的内容存在一个文本文件里tmp.txt：

```Go
{{range .}}----------------------------------------
Name:   {{.Name}}
Price:  {{.Price | printf "%6.2f"}}
{{end}}
```
我们可以这样处理：

```Go
package main

import (
	"log"
	"os"
	"text/template"
)

var report = template.Must(template.ParseFiles("tmp.txt"))

type Book struct {
	Name  string
	Price float64
}

func main() {
	Data := []Book{ {"《三国演义》", 19.82}, {"《儒林外史》", 99.09}, {"《史记》", 26.89} }
	if err := report.Execute(os.Stdout, Data); err != nil {
		log.Fatal(err)
	}
}
```

```Go
程序输出：

----------------------------------------
Name:   《三国演义》
Price:   19.82
----------------------------------------
Name:   《儒林外史》
Price:   99.09
----------------------------------------
Name:   《史记》
Price:   26.89
```


读取模板文件：

```Go
// 建立模板，自动 new("name")
// ParseFiles接受一个字符串，字符串的内容是一个模板文件的路径。
// ParseGlob是用glob的方式匹配多个文件。
Tmpl, err := template.ParseFiles("tmp.txt")  

// 假设一个目录里有a.txt b.txt c.txt的话，用ParseFiles需要写3行对应3个文件，
// 如果有更多文件，可以用ParseGlob。
// 写成template.ParseGlob("*.txt") 即可。
Tmpl, err :=template.ParseGlob("*.txt")

// 函数Must，它的作用是检测模板是否正确，例如大括号是否匹配，
// 注释是否正确的关闭，变量是否正确的书写。
var report = template.Must(template.ParseFiles("tmp.txt"))
```

## 35.2 html/template

和text、template类似，html/template主要在提供支持HTML的功能，所以基本使用上和上面差不多，我们来看看Go语言利用html/template怎样实现一个动态页面：

index.html.tmpl模板文件：

index.html.tmpl：
```Go
<!doctype html>
 <head>
  <meta charset="UTF-8">
  <meta name="Author" content="">
  <meta name="Keywords" content="">
  <meta name="Description" content="">
  <title>Go</title>
 </head>
 <body>
   {{ . }} 
 </body>
</html>
```

main.go：

```Go
package main

import (
	"net/http"
	"text/template"
)

func tHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("index.html.tmpl"))
	t.Execute(w, "Hello World!")
}

func main() {
	http.HandleFunc("/", tHandler)
	http.ListenAndServe(":8080", nil) // 启动web服务
}
```

运行程序，在浏览器打开：http://localhost:8080/    我们可以看到浏览器页面显示Hello World !即使模板文件这时有修改，刷新浏览器后页面也即更新，这个过程并不需要重启web服务。

```Go
func(t *Template) ParseFiles(filenames ...string) (*Template, error)
func(t *Template) ParseGlob(patternstring) (*Template, error)
```

从上面简单的代码中我们可以看到，通过ParseFile加载了单个Html模板文件，当然也可以使用ParseGlob加载多个模板文件。

如果最终的页面很可能是多个模板文件的嵌套结果，ParseFiles也支持加载多个模板文件，模板对象的名字则是第一个模板文件的文件名。

ExecuteTemplate()执行模板渲染的方法，这个方法可用于执行指定名字的模板，因为如果多个模板文件加载情况下，我们需要指定特定的模板渲染执行。面我们根据一段代码来看看：

Layout.html.tmpl模板：

```Go
{{ define "layout" }}

<!doctype html>
 <head>
  <meta charset="UTF-8">
  <meta name="Author" content="">
  <meta name="Keywords" content="">
  <meta name="Description" content="">
  <title>Go</title>
 </head>
 <body>
   {{ . }} 

   {{ template "index" }}
 </body>
</html>

{{ end }}
```

注意模板文件开头根据模板语法，定义了模板名字，{{define "layout"}}。

在模板layout中，通过 {{ template "index" }} 嵌入了模板index，也就是第二个模板文件index.html.tmpl，这个模板文件定义了模板名{{define "index"}}。

注意：通过将模板应用于一个数据结构(即该数据结构作为模板的参数)来执行数据渲染而获得输出。模板执行时会遍历结构并将指针表示为.(称之为dot)，指向运行过程中数据结构的当前位置的值。

{{template "header" .}}  嵌套模板中，加入.dot 代表在该模板中也可以使用该数据结构，否则不能显示。

Index.html.tmpl模板：


```Go
{{ define "index" }}

<div>
<b>Go 语言值得你拥有！</b>
</div>
{{ end }}
```

通过define定义模板名字，还可以通过template action引入模板，类似include。


```Go
package main

import (
	"net/http"
	"text/template"
)

func tHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("layout.html.tmpl", "index.html.tmpl")
	t.ExecuteTemplate(w, "layout", "Hello World!")
}

func main() {
	http.HandleFunc("/", tHandler)
	http.ListenAndServe(":8080", nil)
}
 ```
运行程序，在浏览器打开：http://localhost:8080/ 

```Go
Hello World!
Go 语言值得你拥有！
```


## 35.3 模板语法




（1）模板标签
```Go
{{  }}
```
模板标签用"{{"和"}}"括起来

（2）注释 
```Go
{{/* a comment */}}
```
使用{{/*和*/}}来包含注释内容。

（3）变量
```Go 
{{.}}
```
此标签输出当前对象的值。

```Go 
{{.Admpub}}
```
表示输出对象中字段或方法名称为Admpub的值。
当Admpub是匿名字段时，可以访问其内部字段或方法, 如"Com"：{{.Admpub.Com}} ，如果Com是一个方法并返回一个结构体对象，同样也可以访问其字段或方法：```{{.Admpub.Com.Field1}} ```。


```Go 
{{.Method1 "参数值1" "参数值2"}}
```
调用方法Method1，将后面的参数值依次传递给此方法，并输出其返回值。

```Go 
{{$admpub}}
```
此标签用于输出在模板中定义的名称为"admpub"的变量。当$admpub本身是一个结构体对象时，可访问其字段{{$admpub.Field1}}。

在模板中定义变量，变量名称用字母和数字组成，并带上$前缀，采用简式赋值。

例如：``` {{$x := "OK"}}``` 或 ```{{$x := pipeline}} ```。
 
（4）通道函数
```Go 
{{FuncName1}}
```
此标签将调用名称为"FuncName1"的模板函数（等同于执行"FuncName1()"，不传递任何参数）并输出其返回值。

```Go 
{{FuncName1 "参数值1" "参数值2"}}
```
此标签将调用FuncName1("参数值1", "参数值2")，并输出其返回值。


```Go 
{{.Admpub|FuncName1}}
```
此标签将调用名称为"FuncName1"的模板函数（等同于执行"FuncName1(this.Admpub)"，将竖线"|"左边的".Admpub"变量值作为函数参数传送）并输出其返回值。

（5）条件判断
```Go 
{{if pipeline}} T1 {{end}}
```
标签结构为``` {{if ...}} ... {{end}}```。

```Go 
{{if pipeline}} T1 {{else}} T0 {{end}}
```
标签结构为``` {{if ...}} ... {{else}} ... {{end}}```。

```Go 
{{if pipeline}} T1 {{else if pipeline}} T0 {{end}}
```
标签结构为``` {{if ...}} ... {{else if ...}} ... {{end}}```。
其中if后面可以是一个条件表达式（包括通道函数表达式），也可以是一个字符窜变量或布尔值变量。当为字符窜变量时，如为空字符串则判断为false，否则判断为true。
 
（6）循环遍历
```Go 
{{range $k, $v := .Var}} {{$k}} => {{$v}} {{end}}
```
range...end结构内部如要使用外部的变量，如.Var2，需要写为：$.Var2（即在外部变量名称前加符号"$"）。

```Go 
{{range .Var}} {{.}} {{end}}
```
将遍历值直接显示出来。

```Go 
{{range pipeline}} T1 {{else}} T0 {{end}}
```
当没有可遍历的值时，将执行else部分。

（7）嵌入子模板
```Go
 {{template "name"}}
```
嵌入名称为"name"的子模板。使用前请确保已经用``` {{define "name"}}子模板内容{{end}}```定义好了子模板内容。

```Go
{{template "name" pipeline}}
```
将通道的值赋给子模板中的"."（即"{{.}}"）。

（8）子模板嵌套
```Go
{{define "T1"}}ONE{{end}}
{{define "T2"}}TWO{{end}}
{{define "T3"}}{{template "T1"}} {{template "T2"}}{{end}}
{{template "T3"}}
```
输出如下：
ONE TWO

（9）定义局部变量
```Go
{{with pipeline}} T1 {{end}}
```
通道的值将赋给该标签内部的"."。（注：这里的“内部”一词是指被{{with pipeline}}...{{end}}包围起来的部分，即T1所在位置）

{{with pipeline}} T1 {{else}} T0 {{end}}
如果通道的值为空，"."不受影响并且执行T0，否则，将通道的值赋给"."并且执行T1。
说明：{{end}}标签是if、with、range的结束标签。

（10）输出字符串
```Go
{{"\"output\""}}
```
输出一个字符窜常量。


```Go
{{`"output"`}}
```
输出一个原始字符串常量。


```Go
{{printf "%q" "output"}}
```
函数调用，等同于``` printf("%q", "output")```。
 
```Go
{{"output" | printf "%q"}}
```
竖线"|"左边的结果作为函数最后一个参数，等同于``` printf("%q", "output")```。
 
```Go
{{printf "%q" (print "out" "put")}}
```
圆括号中表达式的整体结果作为printf函数的参数，等同于```printf("%q", print("out", "put"))```。
 
```Go
{{"put" | printf "%s%s" "out" | printf "%q"}}
```
一个更复杂的调用，等同于``` printf("%q", printf("%s%s", "out", "put"))```。
 
```Go
{{"output" | printf "%s" | printf "%q"}}
```
等同于``` printf("%q", printf("%s", "output"))```。
 
```Go
{{with "output"}}{{printf "%q" .}}{{end}}
```
一个使用点号"."的with操作，等同于：``` printf("%q", "output") ```。
 
```Go
{{with $x := "output" | printf "%q"}}{{$x}}{{end}}
```
with结构定义变量，值为执行通道函数之后的结果，等同于``` $x := printf("%q", "output") ```。
 
```Go
{{with $x := "output"}}{{printf "%q" $x}}{{end}}
```
with结构中，在其它动作中使用定义的变量。
 
```Go
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
```
with结构使用了通道，等同于``` printf("%q", "output") ```。

（11）预定义的模板全局函数
```Go
{{and x y}}
```
模板全局函数and，如果x为真，返回y，否则返回x。等同于Go中的x && y。

```Go
{{call .X.Y 1 2}}
```
模板全局函数call，后面的第一个参数的结果必须是一个函数（即这是一个函数类型的值），其余参数作为该函数的参数。
该函数必须返回一个或两个结果值，其中第二个结果值是error类型。
如果传递的参数与函数定义的不匹配或返回的error值不为nil，则停止执行。

```Go
{{html }}
```
模板全局函数html，转义文本中的html标签，如将"<"转义为"&lt;"，">"转义为"&gt;"等。

```Go
{{index x 1 2 3}}
```
模板全局函数index，返回index后面的第一个参数的某个索引对应的元素值，其余的参数为索引值。x必须是一个map、slice或数组。

```Go
{{js}}
```
模板全局函数js，返回用JavaScript的escape处理后的文本。
 
```Go
{{len x}}
```
模板全局函数len，返回参数的长度值（int类型）。

```Go
{{not x}}
```
模板全局函数not，返回单一参数的布尔否定值。

```Go
{{or x y}}
```
模板全局函数or，如果x为真返回x，否则返回y。等同于Go中的：x || y。

```Go
{{print }}
```
模板全局函数print，fmt.Sprint的别名。

```Go
{{printf }}
```
模板全局函数printf，fmt.Sprintf的别名。

```Go
{{println }}
```
模板全局函数println，fmt.Sprintln的别名。

```Go
{{urlquery }}
```
模板全局函数urlquery，返回适合在URL查询中嵌入到形参中的文本转义值。类似于PHP的urlencode。

（12）布尔函数
```Go
{{eq arg1 arg2}}
```
布尔函数eq，返回表达式"arg1 == arg2"的布尔值。

```Go
{{ne arg1 arg2}}
```
布尔函数ne，返回表达式"arg1 != arg2"的布尔值。

```Go
{{lt arg1 arg2}}
```
布尔函数lt，返回表达式"arg1 < arg2"的布尔值。

```Go
{{le arg1 arg2}}
```
布尔函数le，返回表达式"arg1 <= arg2"的布尔值。

```Go
{{gt arg1 arg2}}
```
布尔函数gt，返回表达式"arg1 > arg2"的布尔值。

```Go
{{ge arg1 arg2}}
```
布尔函数ge，返回表达式"arg1 >= arg2"的布尔值。

布尔函数对于任何零值返回false，非零值返回true。对于简单的多路相等测试，eq只接受两个参数进行比较，后面其它的参数将分别依次与第一个参数进行比较。

```Go
{{eq arg1 arg2 arg3 arg4}}
```

即只能作如下比较：
```Go
arg1==arg2 || arg1==arg3 || arg1==arg4 
```



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第三十四章 命令行flag包 ](https://github.com/ffhelicopter/Go42/blob/master/content/42_34_flag.md)

[第三十六章 net/http包](https://github.com/ffhelicopter/Go42/blob/master/content/42_36_http.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
