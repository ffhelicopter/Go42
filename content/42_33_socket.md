# 《Go语言四十二章经》第三十三章 Socket网络

作者：李骁

## 33.1 Socket基础知识

TCP/UDP、IP构成了网络通信的基石，TCP/IP是面向连接的通信协议，要求建立连接时进行3次握手确保连接已被建立，关闭连接时需要4次通信来保证客户端和服务端都已经关闭，也就是我们常说的三次握手，四次挥手。在通信过程中还有保证数据不丢失，在连接不畅通时还需要进行超时重试等等。

Socket就是封装了这一套基于TCP/UDP/IP协议细节，提供了一系列套接字接口进行通信。

我们知道Socket有两种：TCP Socket和UDP Socket，TCP和UDP是协议，而要确定一个进程的需要三元组，还需要IP地址和端口。

* IPv4地址

目前的全球因特网所采用的协议族是TCP/IP协议。IP是TCP/IP协议中网络层的协议，是TCP/IP协议族的核心协议。目前主要采用的IP协议的版本号是4(简称为IPv4)，IPv4的地址位数为32位，也就是最多有2的32次方的网络设备可以联到Internet上。

地址格式类似这样：127.0.0.1   

* IPv6地址

IPv6是新一版本的互联网协议，也可以说是新一代互联网的协议，它是为了解决IPv4在实施过程中遇到的各种问题而被提出的，IPv6采用128位地址长度，几乎可以不受限制地提供地址。在IPv6的设计过程中除了一劳永逸地解决了地址短缺问题以外，还考虑了在IPv4中解决不好的其它问题，主要有端到端IP连接、服务质量（QoS）、安全性、多播、移动性、即插即用等。

地址格式类似这样：2002:c0e8:82e7:0:0:0:c0e8:82e7

## 33.2 TCP 与 UDP 

Go是自带runtime的跨平台编程语言，Go中暴露给语言使用者的TCP socket api是建立OS原生TCP socket接口之上的，所以在使用上相对简单。

TCP Socket

建立网络连接过程：TCP连接的建立需要经历客户端和服务端的三次握手的过程。Go 语言net包封装了系列API，在TCP连接中，服务端是一个标准的Listen + Accept的结构，而在客户端Go语言使用net.Dial或DialTimeout进行连接建立：

在Go语言的net包中有一个类型TCPConn，这个类型可以用来作为客户端和服务器端交互的通道，他有两个主要的函数：

```Go
func (c *TCPConn) Write(b []byte) (n int, err os.Error)
func (c *TCPConn) Read(b []byte) (n int, err os.Error)
```

TCPConn可以用在客户端和服务器端来读写数据。

在Go语言中通过ResolveTCPAddr获取一个TCPAddr：
```Go
func ResolveTCPAddr(net, addr string) (*TCPAddr, os.Error)
```
net参数是"tcp4"、"tcp6"、"tcp"中的任意一个，分别表示TCP(IPv4-only), TCP(IPv6-only)或者TCP(IPv4, IPv6的任意一个)。

addr表示域名或者IP地址，例如"www.google.com:80" 或者"127.0.0.1:22"。

我们来看一个TCP 连接建立的具体代码：

```Go
// TCP server 服务端代码

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {

	var tcpAddr *net.TCPAddr

	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:999")

	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	defer tcpListener.Close()

	fmt.Println("Server ready to read ...")
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		fmt.Println("A client connected : " + tcpConn.RemoteAddr().String())
		go tcpPipe(tcpConn)
	}

}

func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()

	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	i := 0

	for {
		message, err := reader.ReadString('\n') //将数据按照换行符进行读取。
		if err != nil || err == io.EOF {
			break
		}

		fmt.Println(string(message))

		time.Sleep(time.Second * 3)

		msg := time.Now().String() + conn.RemoteAddr().String() + " Server Say hello! \n"
		b := []byte(msg)

		conn.Write(b)
		i++

		if i > 10 {
			break
		}
	}
}
```

服务端 tcpListener.AcceptTCP() 接受一个客户端连接请求，通过go tcpPipe(tcpConn) 开启一个新协程来管理这对连接。 在func tcpPipe(conn *net.TCPConn)  中，处理服务端和客户端数据的交换，在这段代码for中，通过 bufio.NewReader 读取客户端发送过来的数据。

客户端代码：
```Go
// TCP client

package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:999")

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Client connect error ! " + err.Error())
		return
	}

	defer conn.Close()

	fmt.Println(conn.LocalAddr().String() + " : Client connected!")

	onMessageRecived(conn)
}

func onMessageRecived(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	b := []byte(conn.LocalAddr().String() + " Say hello to Server... \n")
	conn.Write(b)
	for {
		msg, err := reader.ReadString('\n')
		fmt.Println("ReadString")
		fmt.Println(msg)

		if err != nil || err == io.EOF {
			fmt.Println(err)
			break
		}
		time.Sleep(time.Second * 2)

		fmt.Println("writing...")

		b := []byte(conn.LocalAddr().String() + " write data to Server... \n")
		_, err = conn.Write(b)

		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
```
客户端net.DialTCP("tcp", nil, tcpAddr) 向服务端发起一个连接请求，调用onMessageRecived(conn)，处理客户端和服务端数据的发送与接收。在func onMessageRecived(conn *net.TCPConn) 中，通过 bufio.NewReader 读取客户端发送过来的数据。

上面2个例子你可以试着运行一下，程序支持多个客户端同时运行。当然，这两个例子只是简单的TCP原始连接，在实际中，我们还可能需要定义协议。

用Socket进行通信，发送的数据包一定是有结构的，类似于：数据头+数据长度+数据内容+校验码+数据尾。而在TCP流传输的过程中，可能会出现分包与黏包的现象。我们为了解决这些问题，需要我们自定义通信协议进行封包与解包。对这方面内容如有兴趣可以去了解更多相关知识。


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第三十二章 fmt包与日志log包](https://github.com/ffhelicopter/Go42/blob/master/content/42_32_fmt.md)

[第三十四章 命令行flag包 ](https://github.com/ffhelicopter/Go42/blob/master/content/42_34_flag.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com