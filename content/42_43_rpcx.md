# rpcx 框架

作者：李骁

## rpcx 框架简介

框架rpcx包含了服务发现、负载均衡、故障转移等服务治理能力，拥有较多的特性，例如无需定义.proto文件，支持跨语言的服务调用等。目前只支持Go语言，但性能良好，可以当作微服务框架来使用。

下面开始来了解下rpcx的使用，文中例子用户服务作为本篇全文的通用示例，看看利用rpcx框架来实现RPC难易程度如何。

首先安装 rpcx框架：

```Go
go get -u -v github.com/smallnest/rpcx/...
```

由于rpcx 后续服务注册中心的需要，还需要加上一些标签来安装，即使这些标签刚开始可能用不上，但建议最好都选择安装，或许最合适的安装命令是这样：

```Go
go get -u -v -tags "reuseport quic kcp zookeeper etcd consul ping" github.com/smallnest/rpcx/...
```


## rpcx 构建服务

由于rpcx对开发的目录结构并没有强制性的规定，所以首先需要为项目规划良好的工程目录结构，下面是用户服务在rpcx中的目录结构：

```Go
└─appservice
    └─member
        ├─cmd
        │  ├─client
        │  │      client.go
        │  │
        │  └─server
        │          server.go
        │
        ├─model
        │      member.go
        │
        └─service
                service.go
```

appservice作为所有服务的总目录入口，member目录是用户服务的目录，下面cmd作为客户端和服务端入口程序的目录，model目录专门用来定义数据结构，而service作为服务的主要实现目录，存放service.go文件，在该文件中定义了用户服务的所有方法和参数类型结构。这就是单个服务整体的目录结构，当然如果有配置项还可以建立conf目录。

该服务的目录结构建议在其他服务中也保持一致，这样在开发中对提升效率会有较大帮助，而且这样的约定也是在开发中非常有必要存在的。

在最关键的service目录中，定义了服务的主要实现。文件service.go主要代码如下：

```Go
type Args struct {
	Uid int
}

type Reply struct {
	model.User
}

type ServiceUser struct {
}

func (s *ServiceUser) UserInfo(ctx context.Context, args *Args, reply *Reply) error {
	fmt.Println("service:", args.Uid)
	reply.User.AddTime = 14990093
	reply.User.Uface = "http://image.xxxx.xxx/t.gif"
	reply.User.UID = int64(args.Uid)
	reply.User.UserName = "Joke"
	reply.User.UserType = 2
	return nil
}
```

ServiceUser作为服务结构体存在，UserInfo(ctx context.Context, args *Args, reply *Reply)方法是用户服务的方法，这个方法需要满足一定的约束：

* 服务方法是可导出的（首字母大写）
* 该方法必须有两个可导出或是内建类型的参数
* 第一个参数为context.Context，第二个参数是输入参数用来接收数据，第三个参数作为输出参数且必须是指针类型
* 方法返回类型为error

这些约束条件中除了第一个参数为context.Context，其他的条件大致与Go语言中定义的RPC方法需要满足一定的条件约束相一致。

在service.go文件中还分别定义了两个可导出的结构体Args和Reply，分别作为服务方法的第二个、第三个参数的类型。这两个参数类型可自定义或是内建类型，第二个参数也就是这里的Args是输入参数（接收），第三个参数也即Reply是输出参数。

对于方法UserInfo()，在实际中应该读取数据库或缓存，在这里不是讨论的重点，故直接赋值。有兴趣的读者可以进行拓展，可在model目录中来处理数据库的访问与处理。

在model目录中的文件member.go定义了用户结构体：

```Go
type User struct {
	UID      int64  `json:"id"`
	AddTime  int64  `json:"addtime"`
	UserType int32  `json:"utype"`
	Uface    string `json:"uface"`
	UserName string `json:"uname"`
}
```

接下来，通过服务端程序注册该服务以及方法，server.go文件在cmd目录下server目录中，主要代码如下：

```Go
var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	s := server.NewServer()
	//s.Register(new(service.ServiceUser), "")
	s.RegisterName("ServiceUser", new(service.ServiceUser), "")
	err := s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
}
```

首先使用 NewServer() 来创建一个服务实例，再通过RegisterName()或者Register()方法注册用户服务，方便客户端从服务注册中心查找并调用用户服务，然后调用 Serve 或者 ServeHTTP 来监听客户端的请求。

在rpcx框架中定义了一个非常重要和关键的结构体Server：

```Go
type Server struct {
	ln                 net.Listener
	readTimeout        time.Duration
	writeTimeout       time.Duration
	gatewayHTTPServer  *http.Server
	DisableHTTPGateway bool // 禁用http调用
	DisableJSONRPC     bool // 禁用json rpc
	serviceMapMu sync.RWMutex
	serviceMap   map[string]*service
	mu         sync.RWMutex
	activeConn map[net.Conn]struct{}
	doneChan   chan struct{}
	seq        uint64
	inShutdown int32
	onShutdown []func(s *Server)
	tlsConfig *tls.Config
	options map[string]interface{}
	// CORS 选项
	corsOptions *CORSOptions 
// 所有的插件
	Plugins PluginContainer
	// AuthFunc 用来鉴权
	AuthFunc func(ctx context.Context, req *protocol.Message, token string) error
	handlerMsgNum int32
}
```

## rpcx 启动选项

在rpcx 框架中，func NewServer(options ...OptionFn) *Server方法先实例化一个Server，然后再设置启动选项，一共提供了 3个 OptionFn 来设置启动选项：

```Go
    func WithReadTimeout(readTimeout time.Duration) OptionFn
    func WithTLSConfig(cfg *tls.Config) OptionFn
    func WithWriteTimeout(writeTimeout time.Duration) OptionFn
```

可以分别用来设置服务读超时、tls证书和写超时，也即设置结构体Server的readTimeout，tlsConfig，writeTimeout 这三个字段的值。当然这三个启动选项是可选的，可根据实际需要来决定。

OptionFn 的定义如下：

```Go
type OptionFn func(*Server)
```

是不是感觉很眼熟！没错，这里采用的就是功能选项设计模式，利用功能选项函数很方便地修改Server实例的字段，也可以做为函数NewServer()的参数来设定启动项的值。

服务注册（RegisterName()或者Register()）会通过反射机制，生成service结构体的实例，该结构体的字段中name为服务注册时的具体服务名，如没指定服务名则默认为该服务（本例中为service.ServiceUser）的类型名。

```Go
type service struct {
	name     string                   // 服务名字
	rcvr     reflect.Value            // 服务方法的接收器
	typ      reflect.Type             // 接收器的类型
	method   map[string]*methodType   // 注册的方法
	function map[string]*functionType // 注册的函数
}
```

最终所有注册的服务会生成serviceMap，也即在结构体Server中定义的字段```Go serviceMap   map[string]*service ```。

有关服务实例的生成和服务注册过程大致就这样。接下来完成客户端的实现，client.go文件在cmd目录下client目录中，主要代码如下：

```Go
var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("ServiceUser", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := service.Args{
		Uid: 999,
	}

	reply := &service.Reply{}
	err := xclient.Call(context.Background(), "UserInfo", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Println(args.Uid, ":", reply.User)

}
```

首先使用 NewPeer2PeerDiscovery() 来初始化点对点的服务发现（客户端直连每个服务节点），所谓的服务发现简单点说就是找到服务器列表。上面NewPeer2PeerDiscovery() 中的参数值tcp@ipaddress:port表示通过TCP通信。

在rpcx框架中可以通过TCP（tcp@ipaddress:port）、HTTP（http@ipaddress:port）、UnixDomain（unix@ipaddress:port）、QUIC（quic@ipaddress:port）和KCP（kcp@ipaddress:port）通信，而且http客户端可以通过网关或者http调用来访问rpcx服务。

在rpcx中使用network @ Host: port格式表示服务地址，network 可以为 tcp ， http ，unix ，quic或kcp，而Host可以是主机名或是IP地址。

接下来通过NewXClient()函数得到客户端的实例，这个客户端实例支持服务发现与服务治理，其结构体xClient如下：

```Go
type xClient struct {
	failMode     FailMode
	selectMode   SelectMode
	cachedClient map[string]RPCClient
	breakers     sync.Map
	servicePath  string
	option       Option
	mu        sync.RWMutex
	servers   map[string]string
	discovery ServiceDiscovery
	selector  Selector

	isShutdown bool
	auth string
	Plugins PluginContainer
	ch chan []*KVPair
	serverMessageChan chan<- *protocol.Message
}
```

再看 ```Go func NewXClient(servicePath string, failMode FailMode, selectMode SelectMode, discovery ServiceDiscovery, option Option) XClient ``` 函数的签名，servicePath 是前面服务端定义的服务名“ServiceUser”，在结构体xClient中对应字段为servicePath ，在客户端和服务端中，这两者对应的字符串需要一致才能正常调用。

上面用户服务的客户端在初始化xClient时选择使用 client.Failtry 错误模式，即调用远程方法失败后再次尝试当前服务器，客户端通过随机选择client.RandomSelect的方式来选择服务器，而服务发现在这里则使用了点对点的方式，也就是直接连接到服务器，可选项为client.DefaultOption，其中默认重试次数为3，默认的编码为MsgPack。大致如图1所示：

![rpcx-1.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/rpcx-1.png)

图1 xClient初始化

客户端在初始化时会根据参数（F.S.D.O，姑且这样称呼）来确定调用失败后处理模式、路由选择的模式、发现服务器列表以及可选配置项。

FailMode和SelectMode为服务治理 (失败模式与路由选择)的选项定义。在大规模的RPC系统中，许多服务节点在提供同一个服务。客户端如何选择最合适的节点来调用呢？如果调用失败，客户端应该选择另一个节点或者立即返回错误？可以通过NewXClient()来指定具体的模式。

失败处理模式FailMode仅对同步调用有效（xClient.Call），而异步调用（xClient.Go）无效，FailMode一共有下面几种值可选择：

```Go
type FailMode int

const (
	// 自动选择另一台服务器
	Failover FailMode = iota
	// 立即返回错误
	Failfast
	// 再次使用当前客户端
	Failtry
	// 如果第一台服务器在指定时间内没有快速响应，则选择另一台服务器
	Failbackup
)
```

Failfast模式：一旦调用服务节点失败，rpcx会立即返回错误。 注意这个错误可能是网络错误或者服务异常原因造成的。

Failover模式：rpcx如果遇到错误，它会尝试调用另外一个节点， 直到有服务节点能正常返回信息，或者达到最大的重试次数。 重试次数Retries在参数Option中设置， 缺省设置为3。

Failtry模式：rpcx调用一个服务节点出现错误，继续重试这个节点直到节点正常返回数据或者达到最大重试次数。

Failbackup模式： 如果服务节点在一定的时间内不返回结果， rpcx客户端会发送相同的请求到另外一个节点，只要在这两个节点中任一节点有返回，rpcx就算调用成功。

而路由选择模式SelectMode则有下面几种情况可选择：

```Go
// SelectMode 定义从候选者中选择服务的算法
type SelectMode int
const (
	// 随机选择
	RandomSelect SelectMode = iota
	// 轮询模式
	RoundRobin
	// 加权轮询模式
	WeightedRoundRobin
	// 加权网络质量优先
	WeightedICMP
	// 一致性Hash
	ConsistentHash
	// 最近的服务器
	Closest

	// 通过用户进行选择
	SelectByUser = 1000
)

```

注意，这里的路由是针对 ServicePath 和 ServiceMethod的路由。

随机模式：从服务节点中随机选择一个节点。由于节点是随机选择，所以并不能保证节点之间负载的均匀。

轮询模式：从服务节点列表中逐个选择依次使用，能保证每个节点均匀被访问，在节点服务能力相差不大时适用。

加权轮询模式：使用基于权重的轮询算法。

网络质量优先：客户端会基于ping(ICMP) 探测各个节点的网络质量，网络质量越好则节点的权重也就越高。

一致性哈希：使用 JumpConsistentHash 选择节点， 相同的servicePath, serviceMethod 和参数会路由到同一个节点上。 JumpConsistentHash 是一个快速计算一致性哈希的算法，但是有一个缺陷是它不能删除节点，如果删除节点，路由需要重新计算一致性哈希。

地理位置优先：它要求服务在注册的时候要设置它所在的地理经纬度。

在rpcx框架中，根据路由选择模式（SelectMode）并通过选择器（Selector）来确定具体的服务器。选择器是一个接口，其定义如下：

```Go
type Selector interface {
    Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string
    UpdateServer(servers map[string]string)
}
```

上述路由选择模式都已经实现选择器的接口。而且在rpcx框架中支持自定义选择器，如果上述路由选择模式不适合，可考虑实现自己的路由选择器。

另外，客户端的可选配置项结构如下：

```Go
type Option struct {
	Group string
	Retries int
	TLSConfig *tls.Config
	Block interface{}
	RPCPath string
	ConnectTimeout time.Duration
	ReadTimeout time.Duration
	WriteTimeout time.Duration
	BackupLatency time.Duration
	GenBreaker func() Breaker
	SerializeType protocol.SerializeType
	CompressType  protocol.CompressType
	Heartbeat         bool
	HeartbeatInterval time.Duration
}
```

在rpcx中，已经预设了一些可选项的值：

```Go
var DefaultOption = Option{
	Retries:        3,
	RPCPath:        share.DefaultRPCPath,
	ConnectTimeout: 10 * time.Second,
	SerializeType:  protocol.MsgPack,
	CompressType:   protocol.None,
	BackupLatency:  10 * time.Millisecond,
}
```

Retries ：重试次数。
ConnectTimeout：连接超时
SerializeType：默认通信协议

还可以设置自动的心跳来保持连接不断掉。客户端需要启用心跳选项，并且设置心跳间隔：

```Go
    option := client.DefaultOption
    option.Heartbeat = true
    option.HeartbeatInterval = time.Second
```

Call()方法是客户端同步远程调用的方法，而另外的Go()方法则是异步远程调用的方法。在这里Call()方法指定调用的RPC方法为用户服务的“UserInfo”方法。当执行Call()方法时，会根据选择器确定的算法（这里是随机）来选择通过服务发现找到的服务器列表，最终确定访问的服务器，远程调用时如果失败则根据失败模式来确定下一步动作，比如上面示例的代码选择Failtry失败模式会重试三次，消息的编码采用MsgPack。当然可以通过设置Option来确定采用其他的编码方式。

用户服务的客户端通过rpcx框架，使用RPC远程调用的方式来调用服务端的方法，现在分别运行服务端和客户端。

在命令行运行服务端程序：

```Go
>go run server.go
2019/07/26 20:50:22 server.go:174: INFO : server pid:724
```

然后在命令行运行客户端程序：

```Go
>go run client.go
2019/07/26 20:50:41 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

在客户端运行后，服务端接收到客户端请求并响应，控制台会显示：

```Go
service: 999
2019/07/26 20:50:41 server.go:358: INFO : client has closed this connection: 127.0.0.1:60186
```

当服务端停止服务后，再运行客户端程序，客户端发现调用远程方法失败，接下来会因为client.Failtry模式而重试，而可选项默认的配置是client.DefaultOption.Retries=3，表示重试的次数为三次，所以这里会重试三次而宣告失败。具体如下图2所示：

![rpcx-2.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/rpcx-2.png)

图2 服务端无响应Failtry

上图中第一次的失败是Call()方法调用失败时的信息。如果把失败模式改为Failfast，停止服务端运行，再运行客户端程序，则调用远程方法时程序会直接报错而不会去尝试重试，具体如下图3所示：

![rpcx-3.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/rpcx-3.png)

图3 服务端无响应Failfast模式

在上面用户服务中，服务注册是针对服务方法而言的，如s.RegisterName("ServiceUser", new(service.ServiceUser), "")，就是将service.ServiceUser用户服务这个结构体的所有方法注册到服务中心。rpcx 也支持将纯函数注册为服务，函数必须满足的条件和前面用户服务中对方法的要求一样：

* 该函数是可导出的（首字母大写）
* 该该函数必须有两个可导出或是内建类型的参数
* 第一个参数为context.Context，第二个参数是输入参数用来接收数据，第三个参数作为输出参数且必须是指针类型
* 函数返回类型为error

接下来在用户服务的service.go文件中增加一个函数，该函数要按照上面要求定义，否则不能注册成功：

```Go
func UserReply(ctx context.Context, args *Args, reply *Reply) error {
	reply.User.AddTime = 10000999
	reply.User.Uface = "http://image.xxxx.xxx/reply.gif"
	reply.User.UID = int64(args.Uid)
	reply.User.UserName = "Reply"
	reply.User.UserType = 3
	return nil
}
```

在服务端，即server.go文件中增加关键的一行，注册该函数到服务中心：

```Go
s.RegisterFunction("ServiceUserFn", service.UserReply, "")
```

上面方法的第一个参数为该函数的自定义服务名，第二个参数为函数名。

接下来在客户端远程调用这个函数，在clientfn.go中主要代码如下：

```Go
func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("ServiceUserFn", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args := service.Args{
		Uid: 888,
	}

	reply := &service.Reply{}
	err := xclient.Call(context.Background(), "UserReply", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Println(args.Uid, ":", reply.User)
}
```

函数的远程调用和服务方法的远程调用基本一样，接下来在命令行运行服务端程序：

```Go
>go run server.go
2019/07/28 21:10:36 server.go:174: INFO : server pid:739
```

然后在命令行运行客户端程序：

```Go
>go run clientfn.go
2019/07/28 21:11:12 888 : {888 10000999 3 http://image.xxxx.xxx/reply.gif Reply}
```

在客户端运行后，服务端接收到客户端请求并响应，控制台会显示：

```Go
2019/07/28 21:11:12 server.go:358: INFO : client has closed this connection: 127.0.0.1:60209
```

通过上面用户服务的例子可以看到，rpcx框架使用非常方便，RPC调用过程整体透明，而服务发现以及治理上只需要简单做好配置即可。这些方面对开发者而言，实在是非常的贴心。

当然，rpcx框架不止上面这些特征，还有其他非常值得了解的特性，下面继续来更深入了解和熟悉这款优秀的RPC框架。

## 服务注册中心

在rpcx框架中，服务注册中心用来实现服务发现和服务元数据的存储。在rpcx框架中支持多种服务注册中心， 并且支持进程内的注册中心，方便开发与测试。rpcx框架会自动将服务的服务名，监听地址，监听协议，权重等信息登记到注册中心，也会定时将服务的吞吐率更新到注册中心。

如果服务意外中断或者宕机，服务注册中心能够监测到事件发生，服务注册中心会通知客户端该服务当前不可用，在服务调用的时候不要再选择这个服务器。

客户端初始化的时候从服务注册中心得到服务器列表，然后根据不同的负载均衡模式选择合适的服务器进行服务调用，同时注册中心会通知客户端某个服务暂时不可用。

服务注册中心与客户端和服务端之间的关系可见下图4：

![rpcx-4.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/rpcx-4.png)

图4 服务注册中心

在rpcx框架中有几种不同的服务注册中心：

### 一、点对点

点对点使用 NewPeer2PeerDiscovery() 来初始化服务发现，由客户端直连服务节点，客户端根据唯一服务器的地址直接连接到服务器，事实上它并没有注册中心。而由于只有一个服务节点，函数func NewXClient()在生成xClient实例时，选择器Selector的selectMode实际上并没有什么作用，因为只有一个节点什么规则最终都只会而且只能选择这个节点。

上面的用户服务中，使用的就是点对点的服务注册中心，最简单直接的方式：

```Go
d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
```

### 二、点对多

点对多顾名思义同一服务部署在多台服务器，可以使用NewMultipleServersDiscovery()来发现部署服务的多台服务器。
    
为了测试这种服务注册，在前面用户服务的基础上建立新的服务，在appservice目录下建立membermultiple目录，暂且称为membermultiple服务。该服务业务逻辑和用户服务一样，只有服务发现有变化。

假设有两台服务器来部署这个服务，为了方便测试，这里需要通过不同的端口来模拟不同的服务器，服务端文件server.go主要代码为：

```Go
var (
	addr0 = flag.String("addr0", "localhost:8972", "server0 address")
	addr1 = flag.String("addr1", "localhost:8973", "server1 address")
)

func main() {
	flag.Parse()

	go createServer(*addr0)
	go createServer(*addr1)

	select {}
}

func createServer(addr string) {
	s := server.NewServer()

	s.RegisterName("ServiceUser", new(service.ServiceUser), "")
	err := s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}
```

上面代码相当于membermultiple服务同时在两台服务器上运行，而客户端采用NewMultipleServersDiscovery()方式来得到服务器信息，这里客户端采用编码的方式来配置服务器地址。

```Go
var (
	addr  = flag.String("addr0", "tcp@localhost:8972", "server0 address")
	addr1 = flag.String("addr1", "tcp@localhost:8973", "server1 address")
)
......
	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr}, {Key: *addr1}})
	xclient := client.NewXClient("ServiceUser", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()
......
	err := xclient.Call(context.Background(), "UserInfo", args, reply)
```

为了更好地观察这种方式的运行，这里稍微修改下rpcx框架中的xclient.go文件，在方法selectClient()中添加	fmt.Println("===", k)语句，这里K的值是访问服务器的协议以及地址，以便观察选择器最终选择服务器的结果。

接下来运行服务端程序：

```Go
>go run server.go
2019/07/30 20:44:23 server.go:174: INFO : server pid:11444
2019/07/30 20:44:23 server.go:174: INFO : server pid:11444
```

模拟的两个服务端已经正常运行，下面运行客户端：

```Go
>go run client.go
=== tcp@localhost:8972
2019/07/30 20:46:07 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

再次运行客户端：

```Go
 >go run client.go
=== tcp@localhost:8973
2019/07/30 20:46:41 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

可以看到两次运行客户端时，远程调用的结果一样，但访问的服务器不一样。经过多次测试发现，现在RandomSelect模式下服务器连接是随机的，并不是轮换，上面测试结果两次服务器不一样只是一种巧合。

现在修改客户端代码中的NewMultipleServersDiscovery()和NewXClient()的参数值为如下所示：

```Go
	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr, Value: "weight=1"}, {Key: *addr1, Value: "weight=9"}})
	xclient := client.NewXClient("ServiceUser", client.Failtry, client.WeightedRoundRobin, d, client.DefaultOption)
```

改为WeightedRoundRobin按照权重轮询后，再多次运行客户端会发现，访问的服务器大多数情况都是addr1，因为它的权重是9，所以基本上都会连接到这台服务器上。

### 三、Etcd

说到Etcd，它是一个强一致的分布式键值存储存储系统，主要用于配置和服务发现，用它来做rpcx框架的服务注册中心是非常合适的选择。下面来具体了解怎样利用Etcd做服务注册中心。

首先需要确定已经安装好Etcd，如果没有则请先安装Etcd。

为了测试Etcd服务注册，在前面用户服务的基础上建立新的服务，在appservice目录下建立memberetcd目录，暂且称为memberetcd服务。该服务业务逻辑和用户服务一样，服务注册在Etcd上面。

下面开始搭建memberetcd服务，服务端文件server.go主要代码为：

```Go
var (
	addr     = flag.String("addr", "localhost:8972", "server address")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {
	flag.Parse()

	s := server.NewServer()
	addRegistryPlugin(s)
	//s.Register(new(service.ServiceUser), "")

	s.RegisterName("ServiceUser", new(service.ServiceUser), "")

	err := s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
}

func addRegistryPlugin(s *server.Server) {
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		EtcdServers:    []string{*etcdAddr},
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
```

这里主要实现了把服务注册到Etcd，主要通过addRegistryPlugin()函数来实现。接下来实现客户端的主要代码如下：

```Go
var (
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)
......
	d := client.NewEtcdDiscovery(*basePath, "ServiceUser", []string{*etcdAddr}, nil)
	xclient := client.NewXClient("ServiceUser", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer xclient.Close()
	......
```

这里客户端采用NewEtcdDiscovery()方式发现服务。前面在安装rpcx时建议加上标签：-tags etcd，在这里也需要用到这个编译标签。在rpcx的etcd_discovery.go文件中带有编译标签：// +build etcd ，所以使用在运行或者编译时需要注意用上这个标签。

首先启动Etcd服务，接下来运行服务端程序：

```Go
>go run -tags etcd server.go
2019/08/05 21:22:38 server.go:174: INFO : server pid:11444
2019/08/05 21:22:38 server.go:174: INFO : server pid:11444
```

模拟服务端已经正常运行，下面运行客户端：

```Go
>go run -tags etcd client.go
=== tcp@localhost:8972
2019/08/05 21:25:16 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

此时停掉Etcd服务，再次运行客户端，则服务端和客户端都会发生错误导致程序不能正常运行。

而再次启动Etcd服务，此时再运行客户端可以得到正常结果：

```Go
>go run -tags etcd client.go
=== tcp@localhost:8972
2019/08/05 21:32:06 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

通过Etcd客户端可以看到注册的服务信息：

/rpcx_test/ServiceUser/tcp@localhost:8972

Etcd作为服务注册中心是可靠的，类似像ZooKeeper、Consul都可以作为可靠的服务注册中心，由于rpcx框架已经封装好了其作为服务注册中心的使用方法，因此Etcd和它们在使用上相差无几，这里就不再列举例子说明。 需要注意的是使用run命令运行或者构建应用程序时需要带上编译标签，如上面例子中 -tags etcd。

## rpcx调用

在rpcx框架中，调用有下面几种方式：

```Go
func (c *xClient) Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
func (c *xClient) Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{},done chan *Call) (*Call, error)
func (c *xClient) Fork(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
func (c *xClient) Broadcast(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
```

Go()方法是异步调用，在异步调用中，失败模式FailMode将会不起作用，它即时返回一个Call结构体实例。

Call()方法是同步调用，也是最常用的调用方式，它会根据选择器确定服务器，支持失败模式FailMode，可以设置Option可选项，来进行远程调用，直到服务器返回数据或者超时。

Broadcast()方法将请求发送到该服务的所有节点。如果所有的节点都正常返回才算成功。只有在所有节点没有错误的情况下， Broadcast()方法将返回其中的一个节点的返回信息。 如果有节点返回错误的话，Broadcast()方法将返回这些错误信息中的一个。失败模式FailMode和路由选择SelectMode在该方法中都不会生效，最好设置超时避免程序挂起。

Fork()方法将请求发送到该服务的所有节点。如果有任何一个节点正常返回，则成功，Fork()方法将返回其中的一个节点的返回结果。 如果所有节点返回错误的话，Fork()方法将返回这些错误信息中的一个。失败模式FailMode和路由选择SelectMode在该方法中都不会生效。

还是在用户服务的基础上来看看Fork()方法的实际运行情况，在用户服务目录下cmd目录中新建fork目录，作为Fork()方法的测试目录。

服务端模拟两个服务器：

```Go
var (
	addr0 = flag.String("addr0", "localhost:8972", "server0 address")
	addr1 = flag.String("addr1", "localhost:8973", "server1 address")
)
```

客户端使用多点服务发现，再使用Fork()方法：

```Go
var (
	addr  = flag.String("addr0", "tcp@localhost:8972", "server0 address")
	addr1 = flag.String("addr1", "tcp@localhost:8973", "server1 address")
)

func main() {
	flag.Parse()

	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr, Value: "weight=1"}, {Key: *addr1, Value: "weight=9"}})
	xclient := client.NewXClient("ServiceUser", client.Failtry, client.WeightedRoundRobin, d, client.DefaultOption)

	defer xclient.Close()

	args := service.Args{
		Uid: 999,
	}

	reply := &service.Reply{}
	err := xclient.Fork(context.Background(), "UserInfo", args, reply)
	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}

	log.Println(args.Uid, ":", reply.User)

}
```

先在命令行运行服务端：

```Go
>go run server.go
2019/08/17 15:47:36 server.go:174: INFO : server pid:9956
2019/08/17 15:47:36 server.go:174: INFO : server pid:9956
```

然后运行客户端：

```Go
>go run client.go
2019/08/17 15:47:50 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

而此时服务端控制台显示：

```Go
service: 999
service: 999
2019/08/17 15:47:50 server.go:358: INFO : client has closed this connection: 127.0.0.1:58487
2019/08/17 15:47:50 server.go:358: INFO : client has closed this connection: 127.0.0.1:58489
```

表明Fork()方法将请求发送到了这两个服务器，得到正常结果后返回。而当服务器全部停止服务时，则Fork()方法直接报ErrXClientNoServer错误信息。

对于RPC来说，序列化对于远程调用的响应速度、吞吐量、网络带宽消耗等也起着至关重要的作用，是分布式系统性能提升的关键因素之一。在rpcx框架中，默认使用 msgpack 编解码器，一共有下面几种编解码器：

SerializeNone：不会对数据进行编解码，要求数据为 []byte 类型。
JSON：通用的数据交换的格式，常规情况下可使用这种编解码。
protocol buffers：一种高性能的编解码器。
MsgPack：默认的编解码器，一种高性能的编解码器，是跨语言的编解码器。
Thrift：一种高性能的编解码器。

开发中可以设置Option.SerializeType来指定合适的编解码器。对于有特殊要求的场景，还可以定制新的编解码器。

编解码也即序列化/反序列化，在rpcx中需要将消息结构体序列化为二进制数据，同时也需要将网络流数据反序列化为内部使用的消息结构体。大致如下图5所示：

![rpcx-5.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/rpcx-5.png)

图5 rpcx编解码

由于在gRPC中只能使用ProtoBuf，因此看看在rpcx中怎样来使用ProtoBuf，下面基于前面的用户服务来实现，新建服务memberproto。

由于有ProtoBuf，建立pb目录来存放member.proto文件：

```Go
syntax = "proto3";

package pb;

message Args {
  int64 Id = 1;
}

message Reply {
   int64	UID =1;	
  int64	AddTime=2;
  int32	UserType=3;
  string Uface =4;
  string UserName=5;
}

message ProtoArgs { 
    int32 A = 1;
    int32 B = 2;
}

message ProtoReply { 
    int32 C = 1;
}
```

运行命令：

```Go
protoc --go_out=. member.proto
```

得到member.pb.go文件，接下来修改service.g文件代码，这次增加了一个方法Mul()：

```Go
func (s *ServiceUser) UserInfo(ctx context.Context, args *pb.Args, reply *pb.Reply) error {
	fmt.Println("service:", args.Id)
	reply.AddTime = 14990093
	reply.Uface = "http://image.xxxx.xxx/t.gif"
	reply.UID = int64(args.Id)
	reply.UserName = "Joke"
	reply.UserType = 2
	return nil
}

func (t *ServiceUser) Mul(ctx context.Context, args *pb.ProtoArgs, reply *pb.ProtoReply) error {
	reply.C = args.A * args.B
	fmt.Printf("call: %d * %d = %d\n", args.A, args.B, reply.C)
	return nil
}
```

客户端需要修改Option可选项，把默认的编解码器改为protocol.ProtoBuffer：

```Go
func main() {
	flag.Parse()

	// register customized codec
	option := client.DefaultOption
	option.SerializeType = protocol.ProtoBuffer

	d := client.NewPeer2PeerDiscovery("tcp@"+*addr, "")
	xclient := client.NewXClient("ServiceUser", client.Failover, client.RandomSelect, d, client.DefaultOption)
	defer xclient.Close()

	args1 := &pb.ProtoArgs{
		A: 10,
		B: 20,
	}

	reply1 := &pb.ProtoReply{}
	err := xclient.Call(context.Background(), "Mul", args1, reply1)

	args := &pb.Args{
		Id: 999,
	}

	reply := &pb.Reply{}

	err = xclient.Call(context.Background(), "UserInfo", args, reply)

	if err != nil {
		log.Fatalf("failed to call: %v", err)
	}
	log.Printf("%d * %d = %d", args1.A, args1.B, reply1.C)
	log.Println(args.Id, ":", reply)
}
```

接下来运行服务端程序：

```Go
>go run server.go
2019/08/17 23:05:22 server.go:174: INFO : server pid:7124
```

然后运行客户端：

```Go
>go run client.go
2019/08/17 23:05:26 10 * 20 = 200
2019/08/17 23:05:26 999 : UID:999 AddTime:14990093 UserType:2 Uface:"http://image.xxxx.xxx/t.gif" UserName:"Joke"
```

很明显新设置的protocol.ProtoBuffer编解码器生效了。由于ProtoBuf使用上更加麻烦，而且和MsgPack相比反倒是MsgPack更有优势，所以在rpcx中默认使用MsgPack也就有了很好的理由。

在rpcx中还可以定制新的编解码器，下面以gob作为新的编解码器，新建服务membergob来实验一下。

首先修改service.go文件，在里面加入gob编解码两个方法：

```Go
type GobCodec struct {
}

func (c *GobCodec) Decode(data []byte, i interface{}) error {
	enc := gob.NewDecoder(bytes.NewBuffer(data))
	err := enc.Decode(i)
	return err
}

func (c *GobCodec) Encode(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(i)
	return buf.Bytes(), err
}
```

在服务端server.go文件中增加新的gob编解码器：

```Go
share.Codecs[protocol.SerializeType(5)] = &service.GobCodec{}
```

在客户端client.go文件中增加新的gob编解码器，同时修改Option选项中的SerializeType为新增的gob编解码器：

```Go
share.Codecs[protocol.SerializeType(5)] = &service.GobCodec{}

option := client.DefaultOption
option.SerializeType = protocol.SerializeType(5)
```

现在可以运行服务端程序：

```Go
>go run server.go
2019/08/17 21:42:57 server.go:174: INFO : server pid:2588
```

然后运行客户端：

```Go
>go run client.go
2019/08/17 21:44:48 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

而此时服务端控制台显示：

```Go
service: 999
2019/08/17 21:47:50 server.go:358: INFO : client has closed this connection: 127.0.0.1:62004
```

通过抓包可以看到端口62004与8972之间存在通信，如下图6所示：

![rpcx-6.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/rpcx-6.png)

图6 gob编解码

如果需要新增其他的编解码器，只需要先实现编解码器，主要实现了Encode和Decode方法，就相当于有了新的编解码器，把新的编解码加入到rpcx中就非常容易了。所以，有兴趣的话读者可以自己尝试一下。


##  超时设置

随着对rpcx框架有了更多的了解，现在开始对框架进行更深一些入细致的了解，更多的细节弄清楚有助于我们更全面了解rpcx框架。比如在客户端和服务端，可以设置超时。

超时机制可认为是一种保护机制，避免服务陷入无限的等待之中。在给定的时间没有响应则服务调用就进入下一个状态，或者重试、或者立即返回错误。

在服务端，主要通过OptionFn来设置两种超时，分别是读超时readTimeout和写超时writeTimeout：

```Go
func WithReadTimeout(readTimeout time.Duration) OptionFn
func WithWriteTimeout(writeTimeout time.Duration) OptionFn
```

既可以在实例化服务时使用NewServer(options ...OptionFn)，也可以直接使用WithReadTimeout()等函数来直接设置。

在客户端可在Option中设置几个超时值：

```Go
type Option struct {
    ……
    // 连接超时
    ConnectTimeout time.Duration
    // 读超时
    ReadTimeout time.Duration
    // 写超时
    WriteTimeout time.Duration
    ……
}
```

在DefaultOption 中设置了连接超时值为 10 秒，但并没有设置 ReadTimeout 和 WriteTimeout。没有设置的超时字段，可以根据情况来设置，但一般默认就可以了。

在客户端中，使用context.Context也可以来控制超时，如使用context.WithTimeout 来设置超时时间，这是在客户端推荐的一种设置超时方式。

```Go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

##  元数据

在上面的远程RPC调用中，远程调用的方法可以通过参数来传递数据。在rpcx中，客户端和服务端可以互相传递元数据。元数据指的也具体业务无关的基础数据。在rpcx中元数据是一个键值队的列表，键和值都是字符串。

在服务器端注册服务的方法中有个设置元数据的参数metadata：

```Go
func (s *Server) RegisterName(name string, rcvr interface{}, metadata string) error
```

参数metadata一般为空字符串，但可以为服务增加一些元数据。比如服务状态（state）就是其中一类元数据，如果在元数据中设置了state=inactive，客户端将不能访问这些服务。这可以帮忙程序员临时禁用一些服务，大致实现服务的降级。

下面实际来看看具体的例子，基于前面的用户服务来实现，新建服务membermeta，在cmd目录中建立state目录，该目录下分别建立server和client目录。大致结构如下：

```Go
membermeta
├─cmd
│  └─state
│      ├─client
│      │      client.go
│      │
│      └─server
│              server.go
│
├─model
│      member.go
│
└─service
        service.go
```

修改server.go文件，在RegisterName()方法中设置元数据state=inactive：

```Go
s.RegisterName("ServiceUser", new(service.ServiceUser), "state=inactive")
```

修改client.go文件中代码：

```Go
d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr, Value: "state=inactive"}, {Key: *addr1, Value: "state=inactive"}})	
```

运行服务端程序：

```Go
>go run server.go
2019/08/18 22:18:01 server.go:174: INFO : server pid:17696
2019/08/18 22:18:01 server.go:174: INFO : server pid:17696
```

然后运行客户端：

```Go
>go run client.go
2019/08/18 22:18:06 connection.go:91: WARN : failed to dial server: dial tcp: missing address
2019/08/18 22:18:06 connection.go:91: WARN : failed to dial server: dial tcp: missing address
2019/08/18 22:18:06 connection.go:91: WARN : failed to dial server: dial tcp: missing address
2019/08/18 22:18:06 failed to call: can not found any server
```

可以看到客户都已经不能正常与服务端通信了，说明设置元数据state=inactive生效。这种通过设置元数据来修改服务状态的方法，在客户端如果不设置，即使服务端设置了也不会生效。

例如把客户端分别修改为下面三种情况：

```Go
d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr}, {Key: *addr1}})	
d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr}, {Key: *addr1, Value: "state=inactive"}})	
d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr, Value: "state=inactive"}, {Key: *addr1}})	
```

在服务端不改变情况下，分别运行三次客户端：

```Go
>go run client.go
=== tcp@localhost:8972
2019/08/18 22:23:38 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

```Go
>go run client.go
=== tcp@localhost:8972
2019/08/18 22:23:55 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

```Go
>go run client.go
=== tcp@localhost:8973
2019/08/18 22:24:13 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

第三种情况很明显，第一个服务器服务状态为inactive，只有第二个服务器正常，所以由它提供服务。

而甚至如果服务端不设置元数据state=inactive，客户端设置了元数据state=inactive也依然生效。如服务端不设置元数据state=inactive，而客户端设置为：

d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr, Value: "state=inactive"}, {Key: *addr1, Value: "state=inactive"}})	

则客户端运行结果依然是能正常与服务端通信。虽然服务端可不用设置，但建议还是设置这个参数。但如果客户端不设置这个元数据，则服务端无论怎样设置都不会起作用，这个必须要注意。

和服务状态类似的还有分组（Group）元数据，在初始化客户端实例时，NewXClient()函数中会调用下面的函数filterByStateAndGroup()，这个函数会检查服务状态和分组两个元数据信息，根据情况来把对应服务器从列表中删除：

```Go
func filterByStateAndGroup(group string, servers map[string]string) {
	for k, v := range servers {
		if values, err := url.ParseQuery(v); err == nil {
			if state := values.Get("state"); state == "inactive" {
				delete(servers, k)
			}
			if group != "" && group != values.Get("group") {
				delete(servers, k)
			}
		}
	}
}
```

为了测试分组元数据的使用，在上面的服务membermeta的cmd目录中新建group目录，分别建立server目录和client目录。

对server.go做简单修改，仅仅只添加metadata参数：

```Go
s.RegisterName("ServiceUser", new(service.ServiceUser), "group=Member")
```

对client.go文件主要修改如下：

```Go
	......
d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr, Value: ""}, {Key: *addr1, Value: "group=Member"}})
	option := client.DefaultOption
	option.Group = "Me"
	xclient := client.NewXClient("ServiceUser", client.Failover, client.RoundRobin, d, option)
	defer xclient.Close()
......
````

运行服务端程序：

```Go
>go run server.go
2019/08/19 00:02:06 server.go:174: INFO : server pid:7400
2019/08/19 00:02:06 server.go:174: INFO : server pid:7400
```

然后运行客户端：

```Go
>go run client.go
2019/08/19 00:14:06 failed to call: can not found any server
```

由于在服务端设置的元数据“group=Member”，而在客户端设置的option.Group = "Me"，客户端和服务端不在一组，所以客户端不能访问到服务器。

此时再修改下client.go文件，修改客户端设置的分组为服务端设置的分组值：

```Go
option.Group = "Member"
```

再运行客户端：

```Go
>go run client.go
=== tcp@localhost:8973
2019/08/19 00:21:38 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

和预想中基本一样，由于服务端和客户端都设置在同一分组，而且在服务发现中指定8973端口的服务器为分组Member，所以应该只有这台服务器可以访问，而且运行结果也证实了这一点。

如果此时把服务发现的分组元数据指定为其他值：

```Go
d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr, Value: ""}, {Key: *addr1, Value: "group=009"}}
```

则客户端运行结果：

```Go
>go run client.go
2019/08/19 00:22:52 failed to call: can not found any server
```

可以看到，运行结果失败，原因是由于该分组没有服务器，所以调用服务失败。所以如果在服务发现中指定分组值，与服务端也需要保持一致。

和服务状态一样，分组也可以在客户端避开。如果在客户端不设置Group这个可选项，其实分组的限制是不起作用的：

```Go
option := client.DefaultOption
option.Group = "Member"
```

把客户端的这个设置取消掉，改为：

```Go
option := client.DefaultOption
```

然后再运行客户端程序：

```Go
go run client.go
=== tcp@localhost:8972
2019/08/19 00:33:08 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
````

虽然服务发现中指定的元数据和服务端不一样，但是由于option.Group 没有设置，分组的限制没有生效。

## 网关

在rpcx框架中，可以通过网关（Gateway）的方式来实现跨语言的调用，比如Java、Python、C#、Node.js、Php、C\C++、Rust等来调用 rpcx 服务。如图7所示：

![rpcx-7.png](https://github.com/ffhelicopter/Go42/blob/master/content/img/rpcx-7.png)

图7 网关

使用网关程序有两种部署模式Gateway和Agent。

1.Gateway：网关模式需要将网关程序独立部署，所有http请求都将先发送给网关，网关将其转换为rpcx请求，再来调用相关rpcx服务，并将rpcx的返回结果转换成http响应, 最终返回给用户。

2.Agent：代理模式是将网关程序和客户端程序一起部署，代理作为一个后台服务部署。客户端发送http请求到本地的代理, 本地的代理将请求转为rpcx请求，并转发到相应的rpcx服务，最后将rpcx的返回结果转换为http响应返回给客户端，类似于Istio中的Sidecar。

下面来实际演示一下网关，在发送http请求时，需要额外设置一些Header信息：

```Go
    X-RPCX-Version: rpcx 版本
    X-RPCX-MesssageType: 设置为0,代表请求
    X-RPCX-Heartbeat: 是否是心跳请求, 缺省false
    X-RPCX-Oneway: 是否是单向请求, 缺省false.
    X-RPCX-SerializeType: 0 as raw bytes, 1 as JSON, 2 as protobuf, 3 as msgpack
    X-RPCX-MessageID: 消息id, uint64 类型
    X-RPCX-ServicePath: service path
    X-RPCX-ServiceMethod: service method
    X-RPCX-Meta: 额外的元数据
```

而对于 http 响应，也有对应的Header信息：

```Go
    X-RPCX-Version: rpcx 版本
    X-RPCX-MesssageType: 1 ,代表响应
    X-RPCX-Heartbeat: 是否是heartbeat请求
    X-RPCX-MessageStatusType: Error 还是正常返回结果
    X-RPCX-SerializeType: 0 as raw bytes, 1 as JSON, 2 as protobuf, 3 as msgpack
    X-RPCX-MessageID: 消息id, uint64 类型
    X-RPCX-ServicePath: service path
    X-RPCX-ServiceMethod: service method
    X-RPCX-Meta: 额外的元数据
    X-RPCX-ErrorMessage: 错误信息
```

在用户服务的基础上来看看网关的实际运行情况，在用户服务目录下cmd目录中新建gateway目录，作为测试目录。

在client目录中增加一个文件httpclient.go，通过这个文件来测试网关连接rpcx服务：

```Go
func main() {
	cc := &codec.MsgpackCodec{}

	args := service.Args{
		Uid: 999,
	}
	data, _ := cc.Encode(args)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8972/", bytes.NewReader(data))
	if err != nil {
		log.Fatal("failed to create request: ", err)
		return
	}

	h := req.Header
	h.Set(gateway.XMessageID, "10000")
	h.Set(gateway.XMessageType, "0")
	h.Set(gateway.XSerializeType, "3")
	h.Set(gateway.XServicePath, "ServiceUser")
	h.Set(gateway.XServiceMethod, "UserInfo")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to call: ", err)
	}
	defer res.Body.Close()

	// handle http response
	replyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("failed to read response: ", err)
	}

	reply := &service.Reply{}
	err = cc.Decode(replyData, reply)
	if err != nil {
		log.Fatal("failed to decode reply: ", err)
	}

	log.Println(args.Uid, ":", reply.User)
}
```

运行服务端程序：

```Go
>go run server.go
2019/08/19 20:53:51 server.go:174: INFO : server pid:16136
```

然后运行含有网关的客户端：

```Go
>go run httpclient.go
2019/08/19 21:02:43 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
```

可以看到http响应正常。

## 断路器（Circuit Breaker）

在rpcx中，可能出现一个节点宕机的情况，可以使用断路器（Circuit Breaker）模式来避免这个错误影响到其他服务，防止出现雪崩情况。

客户端通过断路器调用服务， 一旦连续的错误达到一个阈值，断路器就会断开进行保护，如果继续调用这个节点，系统直接返回错误。经过一段时间，断路器会处于半开的状态，允许一定数量的请求尝试发送到这个节点，如果正常访问，断路器就处于全开的状态，否则又进入短路的状态。

在rpcx 这样定义了断路器 Breaker 接口：

```Go
// Breaker is a CircuitBreaker interface.
type Breaker interface {
	Call(func() error, time.Duration) error
	Fail()
	Success()
	Ready() bool
}
```

在rpcx 中只提供了一个简单的断路器 ConsecCircuitBreaker，实现代码保存在circuit_breaker.go文件中，它在连续发生规定数量的故障或超时后跳闸，再经过一段时间后打开。

```Go
option := client.DefaultOption
option.GenBreaker = func() client.Breaker { return client.NewConsecCircuitBreaker(5, 30*time.Second) }
```

还是在用户服务的cmd目录中建立breaker目录来演示断路器在rpcx中的作用。通过在服务端两个端口只启动一个服务的简单模拟故障发生：

```Go
var (
	addr0 = flag.String("addr0", "localhost:8972", "server0 address")
	addr1 = flag.String("addr1", "localhost:8973", "server1 address")
)

func main() {
	flag.Parse()

	go createServer(*addr0)
	//go createServer(*addr1)

	select {}
}

func createServer(addr string) {
	s := server.NewServer()

	s.RegisterName("ServiceUser", new(service.ServiceUser), "")
	err := s.Serve("tcp", addr)
	if err != nil {
		panic(err)
	}
}
```

客户端需要在Option中指定GenBreaker，这里指定rpcx自带的断路器，这个断路器只有一种触发的条件，即连续多次触发断路，这里设定为2次：

```Go
var (
	addr  = flag.String("addr", "localhost:8972", "server address")
	addr1 = flag.String("addr1", "localhost:8973", "server1 address")
)

func main() {
	flag.Parse()

	option := client.DefaultOption
	option.GenBreaker = func() client.Breaker { return client.NewConsecCircuitBreaker(2, 30*time.Second) }

	d := client.NewMultipleServersDiscovery([]*client.KVPair{{Key: *addr}, {Key: *addr1}})
	option.Retries = 5
	xclient := client.NewXClient("ServiceUser", client.Failtry, client.RandomSelect, d, option)
	defer xclient.Close()

	args := service.Args{
		Uid: 999,
	}
	for i := 0; i < 100; i++ {
		reply := &service.Reply{}
		err := xclient.Call(context.Background(), "UserInfo", args, reply)
		if err != nil {
			log.Printf("failed to call: %v", err)
		}

		log.Println(args.Uid, ":", reply.User)
		time.Sleep(time.Second)
	}
}
```

运行服务端程序：

```Go
>go run server.go
2019/08/20 16:51:30 server.go:174: INFO : server pid:20496
```

然后运行含有断路器的客户端：

```Go
>go run client.go
=== localhost:8973
2019/08/20 18:29:42 connection.go:91: WARN : failed to dial server: dial tcp [::1]:8973: connectex: No connection could be made because the target machine actively refused it.
2019/08/20 18:29:43 connection.go:91: WARN : failed to dial server: dial tcp [::1]:8973: connectex: No connection could be made because the target machine actively refused it.
2019/08/20 18:29:43 failed to call: dial tcp [::1]:8973: connectex: No connection could be made because the target machine actively refused it.
2019/08/20 18:29:43 999 : {0 0 0  }
=== localhost:8973
2019/08/20 18:29:44 failed to call: breaker open
2019/08/20 18:29:44 999 : {0 0 0  }
=== localhost:8973
2019/08/20 18:29:45 failed to call: breaker open
2019/08/20 18:29:45 999 : {0 0 0  }
=== localhost:8972
2019/08/20 18:29:46 999 : {999 14990093 2 http://image.xxxx.xxx/t.gif Joke}
=== localhost:8973
2019/08/20 18:29:47 failed to call: breaker open
2019/08/20 18:29:47 999 : {0 0 0  }
```

客户端设置的失败模式是Failtry，所以重试两次后触发断路器，在断路器生效期间，再次调用则显示breaker open，表明断路器在有效期，这期间不会继续调用已经出问题的服务， 从而达到保护的目的，整个系统不会出现因为超时而产生的雪崩式连锁反应。

另外有开源包：github.com/rubyist/circuitbreaker，提供更多的断路器：

```Go
func NewBreaker() *Breaker                         // 空断路器
func NewThresholdBreaker(threshold int64) *Breaker     // 失败次数
func NewConsecutiveBreaker(threshold int64) *Breaker   // 连续失败次数
func NewRateBreaker(rate float64, minSamples int64) *Breaker // 根据失败比率
```

把客户端代码稍微修改下，改为使用circuitbreaker来做断路器：

```Go
option := client.DefaultOption
	option.GenBreaker = func() client.Breaker {
		return circuit.NewBreakerWithOptions(&circuit.Options{
			ShouldTrip: circuit.ThresholdTripFunc(2),
			WindowTime: 30 * time.Second,
		})
	}
```

运行客户端，也可以实现断路器模式：

```Go
>go run client.go
......
2019/08/20 19:30:06 failed to call: dial tcp [::1]:8973: connectex: No connection could be made because the target machine actively refused it.
2019/08/20 19:30:06 999 : {0 0 0  }
=== localhost:8973
2019/08/20 19:30:07 failed to call: breaker open
2019/08/20 19:30:07 999 : {0 0 0  }
```
