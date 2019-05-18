# 《Go语言四十二章经》第三十八章 数据序列化

作者：李骁

## 38.1 序列化与反序列化

我们的数据对象要在网络中传输或保存到文件，就需要对其编码和解码动作，目前存在很多编码格式：JSON，XML，Gob，Google Protocol Buffer等，Go 语言当然也支持所有这些编码格式。

序列化 (Serialization)是将对象的状态信息转换为可以存储或传输的形式的过程。在序列化期间，对象将其当前状态写入到临时或持久性存储区。通过从存储区中读取对象的状态，重新创建该对象，则为反序列化。

简单地说把某种数据结构转为指定数据格式为“序列化”或“编码”（传输之前）；而把“指定数据格式”转为某种数据结构则为“反序列化”或“解码”（传输之后）。

在Go语言中，encoding/json标准包处理JSON数据的序列化与反序列化问题。

JSON数据序列化函数主要有：

json.Marshal() 

```Go
func Marshal(v interface{}) ([]byte, error) {
	e := newEncodeState()

	err := e.marshal(v, encOpts{escapeHTML: true})
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)

	e.Reset()
	encodeStatePool.Put(e)

	return buf, nil
}
```

从上面的Marshal()函数我们可以看到，数据结构序列化后返回的是字节数组，而字节数组很容易通过网络传输或写入文件存储。而且在Go中，Marshal()默认是设置escapeHTML = true的，会自动把 '<', '>', 以及 '&' 等转化为"\u003c" ， "\u003e"以及 "\u0026"。

JSON数据反序列化函数主要有：

UnMarshal() 

```Go
func Unmarshal(data []byte, v interface{}) error // 把 JSON 解码为数据结构
```

从上面的UnMarshal()函数我们可以看到，反序列化是读取字节数组，进而解析为对应的数据结构。


注意：

不是所有的数据都可以编码为 JSON 格式,只有验证通过的数据结构才能被编码：

* json 对象只支持字符串类型的 key；要编码一个 Go map 类型，map 必须是 map[string]T（T是 json 包中支持的任何类型）
* channel，复杂类型和函数类型不能被编码
* 不支持循环数据结构；它将引起序列化进入一个无限循环
* 指针可以被编码，实际上是对指针指向的值进行编码（或者指针是 nil）


而在Go中，JSON 与 Go 类型对应如下：

* bool    对应 JSON 的 booleans
* float64 对应 JSON 的 numbers
* string  对应 JSON 的 strings
* nil     对应 JSON 的 null

在解析 JSON 格式数据时，若以 interface{} 接收数据，需要按照以上规则进行解析。


## 38.2 JSON数据格式

在Go语言中，利用encoding/json标准包将数据序列化为JSON数据格式这个过程简单直接，直接使用json.Marshal(v)来处理任意类型，序列化成功后得到一个字节数组。 

反过来我们将一个JSON数据来反序列化或解码，则就不那么容易了，下面我们一一来说明。

（一）将JSON数据反序列化到结构体：

这种需求是最常见的，在我们知道 JSON 的数据结构前提情况下，我们完全可以定义一个或几个适当的结构体并对 JSON 数据反序列化。例如：

```Go
package main

import (
	"encoding/json"
	"fmt"
)

type Human struct {
	name   string `json:"name"` // 姓名
	Gender  string `json:"s"`    // 性别，性别的tag表明在json中为s字段
	Age    int    `json:"Age"`  // 年龄
	Lesson
}

type Lesson struct {
	Lessons []string `json:"lessons"`
}

func main() {
	jsonStr := `{"Age": 18,"name": "Jim" ,"s": "男",
	"lessons":["English","History"],"Room":201,"n":null,"b":false}`

	var hu Human
	if err := json.Unmarshal([]byte(jsonStr), &hu); err == nil {
		fmt.Println("\n结构体Human")
		fmt.Println(hu)
	}

	var le Lesson
	if err := json.Unmarshal([]byte(jsonStr), &le); err == nil {
		fmt.Println("\n结构体Lesson")
		fmt.Println(le)
	}

	jsonStr = `["English","History"]`

	var str []string
	if err := json.Unmarshal([]byte(jsonStr), &str); err == nil {
		fmt.Println("\n字符串数组")
		fmt.Println(str)
	} else {
		fmt.Println(err)
	}
}

程序输出：
结构体Human
{ 男 18 {[English History]}}

结构体Lesson
{[English History]}

字符串数组
[English History 
```

我们定义了2个结构体Human和Lesson，结构体Human的Gender字段tag标签为：`json:"s"`，表明这个字段在JSON中的名字对应为s。而且结构体Human中嵌入了Lesson结构体。

jsonStr 我们可以认作为一个JSON数据，通过json.Unmarshal，我们可以把JSON中的数据反序列化到了对应结构体，由于结构体Human的name字段不能导出，所以并不能实际得到JSON中的值，这是我们在定义结构体时需要注意的，字段首字母大写。

对JSON中的Age，在结构体Human对应Age int，不能是string。另外，如果是JSON数组，可以把数据反序列化给一个字符串数组。

总之，知道JSON的数据结构很关键，有了这个前提做反序列化就容易多了。而且结构体的字段并不需要和JSON中所有数据都一一对应，定义的结构体字段可以是JSON中的一部分。

（二）反序列化任意JSON数据：

encoding/json 包使用 map[string]interface{} 和 []interface{} 储存任意的 JSON 对象和数组；其可以被反序列化为任何的 JSON blob 存储到接口值中。

直接使用 Unmarshal 把这个数据反序列化，并保存在map[string]interface{} 中，要访问这个数据，我们可以使用类型断言：

```Go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonStr := `{"Age": 18,"name": "Jim" ,"s": "男","Lessons":["English","History"],"Room":201,"n":null,"b":false}`

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err == nil {
		fmt.Println("map结构")
		fmt.Println(data)
	}

	for k, v := range data {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "是string", vv)
		case bool:
			fmt.Println(k, "是bool", vv)
		case float64:
			fmt.Println(k, "是float64", vv)
		case nil:
			fmt.Println(k, "是nil", vv)
		case []interface{}:
			fmt.Println(k, "是array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}
		default:
			fmt.Println(k, "未知数据类型")
		}
	}
}

程序输出：
map结构
map[n:<nil> b:false Age:18 name:Jim s:男 Lessons:[English History] Room:201]
name 是string Jim
s 是string 男
Lessons 是array:
0 English
1 History
Room 是float64 201
n 是nil <nil>
b 是bool false
Age 是float64 18
```

通过这种方式，即使是未知 JSON 数据结构，我们也可以反序列化，同时可以确保类型安全。在switch-type中，我们可以根据表16-3 JSON与Go数据类型对照表来做选择。比如Age是float64而不是int类型，另外JSON的booleans、null类型在JSON也常常出现，在这里都做了case。


（三）JSON数据编码和解码：

json 包提供 Decoder 和 Encoder 类型来支持常用 JSON 数据流读写。NewDecoder 和 NewEncoder 函数分别封装了 io.Reader 和 io.Writer 接口。

```Go
func NewDecoder(r io.Reader) *Decoder
func NewEncoder(w io.Writer) *Encoder
```

如果要想把 JSON 直接写入文件，可以使用 json.NewEncoder 初始化文件（或者任何实现 io.Writer 的类型），并调用 Encode()；反过来与其对应的是使用 json.Decoder 和 Decode() 函数：

```Go
func NewDecoder(r io.Reader) *Decoder
func (dec *Decoder) Decode(v interface{}) error
```

由于 Go 语言中很多标准包都实现了 Reader 和 Writer接口，因此 Encoder 和 Decoder 使用起来非常方便。

例如，下面例子使用 Decode方法解码一段JSON格式数据，同时使用Encode方法将我们的结构体数据保存到文件t.json中：

```Go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Human struct {
	name   string `json:"name"` // 姓名
	Gender string `json:"s"`    // 性别，性别的tag表明在json中为s字段
	Age    int    `json:"Age"`  // 年龄
	Lesson
}

type Lesson struct {
	Lessons []string `json:"lessons"`
}

func main() {
	// json数据的字符串
	jsonStr := `{"Age": 18,"name": "Jim" ,"s": "男",
	"lessons":["English","History"],"Room":201,"n":null,"b":false}`
	strR := strings.NewReader(jsonStr)
	h := &Human{}

	// Decode 解码json数据到结构体Human中
	err := json.NewDecoder(strR).Decode(h)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(h)

	// 定义Encode需要的Writer
	f, err := os.Create("./t.json")

	// 把保存数据的Human结构体对象编码为json保存到文件
	json.NewEncoder(f).Encode(h)

}

程序输出：
&{ 男 18 {[English History]}}
```

我们调用json.NewDecoder 函数构造了 Decoder 对象，使用这个对象的 Decode方法解码给定义好的结构体对象h。对于字符串，使用 strings.NewReader 方法，让字符串变成一个 Reader。

类似解码过程，我们通过json.NewEncoder()函数来构造Encoder对象，由于os中文件操作已经实现了Writer接口，所以可以直接使用，把h结构体对象编码为JSON数据格式保存在文件t.json中。

文件t.json中内容为：

{"s":"男","Age":18,"lessons":["English","History"]}

（四）JSON数据延迟解析

Human.Name字段，由于可以等到使用的时候，再根据具体数据类型来解析，因此我们可以延迟解析。当结构体Human的Name字段的类型设置为 json.RawMessage 时，它将在解码后继续以 byte 数组方式存在。

```Go
package main

import (
	"encoding/json"
	"fmt"
)

type Human struct {
	Name   json.RawMessage `json:"name"` // 姓名，json.RawMessage 类型不会进行解码
	Gender string          `json:"s"`    // 性别，性别的tag表明在json中为s字段
	Age    int             `json:"Age"`  // 年龄
	Lesson
}

type Lesson struct {
	Lessons []string `json:"lessons"`
}

func main() {
	jsonStr := `{"Age": 18,"name": "Jim" ,"s": "男",
	"lessons":["English","History"],"Room":201,"n":null,"b":false}`

	var hu Human
	if err := json.Unmarshal([]byte(jsonStr), &hu); err == nil {
		fmt.Printf("\n 结构体Human \n")
		fmt.Printf("%+v \n", hu) // 可以看到Name字段未解码，还是字节数组
	}

	// 对延迟解码的Human.Name进行反序列化
	var UName string
	if err := json.Unmarshal(hu.Name, &UName); err == nil {
		fmt.Printf("\n Human.Name: %s \n", UName)
	}
}

程序输出：

 结构体Human 
{Name:[34 74 105 109 34] Gender:男 Age:18 Lesson:{Lessons:[English History]}} 

 Human.Name: Jim 
```

在对JSON数据第一次解码后，保存在Human的hu.Name的值还是二进制数组，在后面对hu.Name进行解码后才真正发序列化为string类型的真实字符串对象。

除了Go标准库外，还有很多的第三方库也能较好解析JSON数据。这里我推荐一个第三方库：https://github.com/buger/jsonparser

如同 encoding/json 包一样，在Go语言中XML也有 Marshal() 和 UnMarshal() 从 XML 中编码和解码数据；也可以从文件中读取和写入（或者任何实现了 io.Reader 和 io.Writer 接口的类型）。和 JSON 的方式一样，XML 数据可以序列化为结构，或者从结构反序列化为 XML 数据。


## 38.3 Protocol Buffer数据格式

Protocol Buffer 简单称为protobuf(Pb)，是Google开发出来的一个语言无关、平台无关的数据序列化工具，在rpc或tcp通信等很多场景都可以使用。在服务端定义一个数据结构，通过protobuf转化为字节流，再传送到客户端解码，就可以得到对应的数据结构。它的通信效率极高，同一条消息数据，用protobuf序列化后的大小是JSON的10分之一左右。

为了正常使用protobuf，我们需要做一些准备工作。

1、下载protobuf的编译器protoc，地址：https://github.com/google/protobuf/releases

window用户下载: protoc-3.6.1-win32.zip，然后解压，把protoc.exe文件复制到GOPATH/bin下。
linux 用户下载：protoc-3.6.1-linux-x86_64.zip 或 protoc-3.6.1-linux-x86_32.zip，然后解压，把protoc文件复制到GOPATH/bin下。

2、获取protobuf的编译器插件protoc-gen-go。

在命令行运行 go get -u github.com/golang/protobuf/protoc-gen-go
会在GOPATH/bin下生成protoc-gen-go.exe文件，如果没有请自行build。GOPATH/bin目录建议加入path，以便后续操作方便。

接下来我们可以正式开始尝试怎么使用protobuf了。我们需要创建一个.proto 结尾的文件，这个文件需要按照一定规则编写。

具体请见官方说明：https://developers.google.com/protocol-buffers/docs/proto
也可参考：https://gowalker.org/github.com/golang/protobuf/proto

protobuf的使用方法是将数据结构写入到.proto文件中，使用protoc编译器编译（通过调用protoc-gen-go）得到一个新的go包，里面包含go中可以使用的数据结构和一些辅助方法。

下面我们先创建一个msg.proto文件

```Go
syntax = "proto3";

package learn;

message UserInfo {
    int32 UserType = 1;     //必选字段
    string UserName = 2;    //必选字段
    string UserInfo = 3;    //必选字段
}
```

运行如下命令

```Go
> protoc --go_out=.  msg.proto
```

会生成一个msg.pb.go的文件，代码如下。

```Go
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg.proto

package learn

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type UserInfo struct {
	UserType             int32    `protobuf:"varint,1,opt,name=UserType,proto3" json:"UserType,omitempty"`
	UserName             string   `protobuf:"bytes,2,opt,name=UserName,proto3" json:"UserName,omitempty"`
	UserInfo             string   `protobuf:"bytes,3,opt,name=UserInfo,proto3" json:"UserInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserInfo) Reset()         { *m = UserInfo{} }
func (m *UserInfo) String() string { return proto.CompactTextString(m) }
func (*UserInfo) ProtoMessage()    {}
func (*UserInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

func (m *UserInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserInfo.Unmarshal(m, b)
}
func (m *UserInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserInfo.Marshal(b, m, deterministic)
}
func (m *UserInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserInfo.Merge(m, src)
}
func (m *UserInfo) XXX_Size() int {
	return xxx_messageInfo_UserInfo.Size(m)
}
func (m *UserInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_UserInfo.DiscardUnknown(m)
}

var xxx_messageInfo_UserInfo proto.InternalMessageInfo

func (m *UserInfo) GetUserType() int32 {
	if m != nil {
		return m.UserType
	}
	return 0
}

func (m *UserInfo) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

func (m *UserInfo) GetUserInfo() string {
	if m != nil {
		return m.UserInfo
	}
	return ""
}

func init() {
	proto.RegisterType((*UserInfo)(nil), "learn.UserInfo")
}

func init() { proto.RegisterFile("msg.proto", fileDescriptor_c06e4cca6c2cc899) }

var fileDescriptor_c06e4cca6c2cc899 = []byte{
	// 100 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcc, 0x2d, 0x4e, 0xd7,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0x49, 0x4d, 0x2c, 0xca, 0x53, 0x8a, 0xe3, 0xe2,
	0x08, 0x2d, 0x4e, 0x2d, 0xf2, 0xcc, 0x4b, 0xcb, 0x17, 0x92, 0x82, 0xb0, 0x43, 0x2a, 0x0b, 0x52,
	0x25, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83, 0xe0, 0x7c, 0x98, 0x9c, 0x5f, 0x62, 0x6e, 0xaa, 0x04,
	0x93, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x9c, 0x0f, 0x93, 0x03, 0x99, 0x21, 0xc1, 0x8c, 0x90, 0x03,
	0xf1, 0x93, 0xd8, 0xc0, 0xb6, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x7e, 0x5d, 0xad, 0x78,
	0x7a, 0x00, 0x00, 0x00,
}
```


接下来，我们在Go语言程序中使用protobuf。

```Go
package main

import (
	"github.com/golang/protobuf/proto"

	"fmt"
	"ind/pb"
)

func main() {
	//初始化message UserInfo
	usermsg := &pb.UserInfo{
		UserType: 1,
		UserName: "Jok",
		UserInfo: "I am a woker!",
	}

	//序列化
	userdata, err := proto.Marshal(usermsg)
	if err != nil {
		fmt.Println("Marshaling error: ", err)
	}

	//反序列化
	encodingmsg := &pb.UserInfo{}
	err = proto.Unmarshal(userdata, encodingmsg)

	if err != nil {
		fmt.Println("Unmarshaling error: ", err)
	}

	fmt.Printf("GetUserType: %d\n", encodingmsg.GetUserType())
	fmt.Printf("GetUserName: %s\n", encodingmsg.GetUserName())
	fmt.Printf("GetUserInfo: %s\n", encodingmsg.GetUserInfo())
}
```

```Go
程序输出：

GetUserType: 1
GetUserName: Jok
GetUserInfo: I am a woker!
```

通过上面的介绍，我们已经学会了怎么使用protobuf来处理我们的数据。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第三十七章 context包](https://github.com/ffhelicopter/Go42/blob/master/content/42_37_context.md)

[第三十九章 Mysql数据库](https://github.com/ffhelicopter/Go42/blob/master/content/42_39_mysql.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com