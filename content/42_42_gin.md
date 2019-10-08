# 《Go语言四十二章经》第四十二章 WEB框架(Gin)

作者：李骁

## 42.1 有关于Gin

在Go语言开发的WEB框架中，有两款著名WEB框架的命名都以酒有关：Martini（ 马丁尼）和Gin（杜松子酒），由于我不擅于饮酒所以这两种酒的优劣暂不做评价，但说WEB框架相比较的话，Gin要比Martini强很多。

Gin是Go语言写的一个WEB框架，它具有运行速度快，分组的路由器，良好的崩溃捕获和错误处理，非常好的支持中间件和JSON。总之在Go语言开发领域是一款值得好好研究的WEB框架，开源网址：https://github.com/gin-gonic/gin

首先下载安装gin包：

```Go
go get -u github.com/gin-gonic/gin
```
一个简单的例子：

```Go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong", 
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
```

编译运行程序，打开浏览器，访问 http://localhost:8080/ping  
页面显示：

```Go
{"message":"pong"}
```

以JSON格式输出了数据。

gin的功能不只是简单输出JSON数据。它是一个轻量级的WEB框架，支持RESTful风格API，支持GET，POST，PUT，PATCH，DELETE，OPTIONS 等http方法，支持文件上传，分组路由，Multipart/Urlencoded FORM，以及支持JSONP，参数处理等等功能，这些都和WEB紧密相关，通过提供这些功能，使开发人员更方便地处理WEB业务。

## 42.2 Gin实际应用

接下来使用Gin作为框架来搭建一个拥有静态资源站点，动态WEB站点，以及RESTFull API接口站点（可专门作为手机APP应用提供服务使用）组成的，亦可根据情况分拆这套系统，每种功能独立出来单独提供服务。

下面按照一套系统但采用分站点来说明，首先是整个系统的目录结构，website目录下面static是资源类文件，为静态资源站点专用；photo目录是UGC上传图片目录，tpl是动态站点的模板。当然这个目录结构是一种约定，你可以根据情况来修改。整个项目已经开源，你可以访问来详细了解：https://github.com/ffhelicopter/tmm



具体每个站点的功能怎么实现呢？请看下面有关每个功能的讲述：

**一：静态资源站点**

一般网站开发中，我们会考虑把js，css，以及资源图片放在一起，作为静态站点部署在CDN，提升响应速度。采用Gin实现起来非常简单，当然也可以使用net/http包轻松实现，但使用Gin会更方便。

不管怎么样，使用Go开发，我们可以不用花太多时间在WEB服务环境搭建上，程序启动就直接可以提供WEB服务了。

```Go
package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 静态资源加载，本例为css,js以及资源图片
	router.StaticFS("/public", http.Dir("D:/goproject/src/github.com/ffhelicopter/tmm/website/static"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// Listen and serve on 0.0.0.0:80
	router.Run(":80")
}
```
首先需要是生成一个Engine ，这是gin的核心，默认带有Logger 和 Recovery 两个中间件。
```Go
router := gin.Default()
```
StaticFile 是加载单个文件，而StaticFS 是加载一个完整的目录资源：
```Go
func (group *RouterGroup) StaticFile(relativePath, filepath string) IRoutes
func (group *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes
```
这些目录下资源是可以随时更新，而不用重新启动程序。现在编译运行程序，静态站点就可以正常访问了。

访问http://localhost/public/images/logo.jpg 图片加载正常。每次请求响应都会在服务端有日志产生，包括响应时间，加载资源名称，响应状态值等等。

**二：动态站点**

如果需要动态交互的功能，比如发一段文字+图片上传。由于这些功能出来前端页面外，还需要服务端程序一起来实现，而且迭代需要经常需要修改代码和模板，所以把这些统一放在一个大目录下，姑且称动态站点。

tpl是动态站点所有模板的根目录，这些模板可调用静态资源站点的css，图片等；photo是图片上传后存放的目录。
```Go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ffhelicopter/tmm/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 静态资源加载，本例为css,js以及资源图片
	router.StaticFS("/public", http.Dir("D:/goproject/src/github.com/ffhelicopter/tmm/website/static"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 导入所有模板，多级目录结构需要这样写
	router.LoadHTMLGlob("website/tpl/*/*")

	// website分组
	v := router.Group("/")
	{

		v.GET("/index.html", handler.IndexHandler)
		v.GET("/add.html", handler.AddHandler)
		v.POST("/postme.html", handler.PostmeHandler)
	}

	// router.Run(":80") 
	// 这样写就可以了，下面所有代码（go1.8+）是为了优雅处理重启等动作。
	srv := &http.Server{
		Addr:         ":80",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		// 监听请求
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 优雅Shutdown（或重启）服务
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt) // syscall.SIGKILL
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
	}	
	log.Println("Server exiting")
}
```
在动态站点实现中，引入WEB分组以及优雅重启这两个功能。WEB分组功能可以通过不同的入口根路径来区别不同的模块，这里我们可以访问：http://localhost/index.html 。如果新增一个分组，比如：
```Go
v := router.Group("/login") 
```
我们可以访问：http://localhost/login/xxxx ，xxx是我们在v.GET方法或v.POST方法中的路径。

```Go
	// 导入所有模板，多级目录结构需要这样写
	router.LoadHTMLGlob("website/tpl/*/*")
	
	// website分组
	v := router.Group("/")
	{

		v.GET("/index.html", handler.IndexHandler)
		v.GET("/add.html", handler.AddHandler)
		v.POST("/postme.html", handler.PostmeHandler)
	}
```

通过router.LoadHTMLGlob("website/tpl/*/*") 导入模板根目录下所有的文件。在前面有讲过html/template 包的使用，这里模板文件中的语法和前面一致。

```Go
	router.LoadHTMLGlob("website/tpl/*/*")
```

比如v.GET("/index.html", handler.IndexHandler) ，通过访问http://localhost/index.html 这个URL，实际由handler.IndexHandler来处理。而在tmm目录下的handler存放了package handler 文件。在包里定义了IndexHandler函数，它使用了index.html模板。

```Go
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "作品欣赏",
	})
}
```
index.html模板：

```Go
<!DOCTYPE html>
<html>
<head>
{{template "header" .}}
</head>
<body>

<!--导航-->
<div class="feeds">
    <div class="top-nav">
    	<a href="/index.tml" class="active">欣赏</a>
    	<a href="/add.html" class="add-btn">
    		<svg class="icon" aria-hidden="true">
    			<use  xlink:href="#icon-add"></use>
    		</svg>
    		发布
    	</a>
    </div>
	<input type="hidden" id="showmore" value="{$showmore}">
	<input type="hidden" id="page" value="{$page}">
    <!--</div>-->
</div>
<script type="text/javascript">
	var done = true;
	$(window).scroll(function(){
        var scrollTop = $(window).scrollTop();
        var scrollHeight = $(document).height();
        var windowHeight = $(window).height();
        var showmore = $("#showmore").val();
        if(scrollTop + windowHeight + 300 >= scrollHeight && showmore == 1 && done){
        	var page = $("#page").val();
        	done = false;
	        $.get("{:U('Product/listsAjax')}", { page : page }, function(json) {
	        	if (json.rs != "") {
	        		$(".feeds").append(json.rs);
	        		$("#showmore").val(json.showmore);
	        		$("#page").val(json.page);
	        		done = true;
	        	}
	        },'json');
        }
    });
</script>
    <script src="//at.alicdn.com/t/font_ttszo9rnm0wwmi.js"></script>
</body>
</html>
```

在index.html模板中，通过{{template "header" .}}语句，嵌套了header.html模板。

header.html模板：

```Go
{{ define "header" }}
	<meta charset="UTF-8">	
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no, minimal-ui">
	<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
	<meta name="format-detection" content="telephone=no,email=no">
	<title>{{ .Title }}</title>
	<link rel="stylesheet" href="/public/css/common.css">
	<script src="/public/lib/jquery-3.1.1.min.js"></script>
    <script src="/public/lib/jquery.cookie.js"></script>
    <link href="/public/css/font-awesome.css?v=4.4.0" rel="stylesheet">
{{ end }}
```

{{ define "header" }} 让我们在模板嵌套时直接使用header名字，而在index.html中的{{template "header" .}} 注意“.”，可以使参数嵌套传递，否则不能传递，比如这里的Title。

现在我们访问 http://localhost/index.html  ，可以看到浏览器显示Title 是“作品欣赏”，这个Title是通过IndexHandler来指定的。

接下来点击“发布”按钮，我们进入发布页面，上传图片，点击“完成”提交，会提示我们成功上传图片。可以在photo目录中看到刚才上传的图片。

注意：
由于在本人在发布到github的代码中，在处理图片上传的代码中，除了服务器存储外，还实现了IPFS发布存储，如果不需要IPFS，请注释相关代码。

有关IPFS:
IPFS本质上是一种内容可寻址、版本化、点对点超媒体的分布式存储、传输协议，目标是补充甚至取代过去20年里使用的超文本媒体传输协议（HTTP），希望构建更快、更安全、更自由的互联网时代。

IPFS 不算严格意义上区块链项目，是一个去中心化存储解决方案，但有些区块链项目通过它来做存储。

IPFS项目有在github上开源，Go语言实现哦，可以关注并了解。

优雅重启在迭代中有较好的实际意义，每次版本发布，如果直接停服务在部署重启，对业务还是有蛮大的影响，而通过优雅重启，这方面的体验可以做得更好些。这里ctrl + c 后过5秒服务停止。


**三：中间件的使用，在API中可能使用限流，身份验证等**

Go 语言中net/http设计的一大特点就是特别容易构建中间件。 gin也提供了类似的中间件。需要注意的是在gin里面中间件只对注册过的路由函数起作用。

而对于分组路由，嵌套使用中间件，可以限定中间件的作用范围。大致分为全局中间件，单个路由中间件和分组中间件。

即使是全局中间件，其使用前的代码不受影响。 也可在handler中局部使用，具体见api.GetUser 。

在高并发场景中，有时候需要用到限流降速的功能，这里引入一个限流中间件。有关限流方法常见有两种，具体可自行研究，这里只讲使用。

导入 import  "github.com/didip/tollbooth/limiter" 包，在上面代码基础上增加如下语句：

```Go
	//rate-limit 限流中间件 
	lmt := tollbooth.NewLimiter(1, nil)
	lmt.SetMessage("服务繁忙，请稍后再试...")
```

并修改

```Go
v.GET("/index.html", LimitHandler(lmt), handler.IndexHandler) 
```

当F5刷新刷新http://localhost/index.html 页面时，浏览器会显示：服务繁忙，请稍后再试...

限流策略也可以为IP：

```Go
tollbooth.LimitByKeys(lmt, []string{"127.0.0.1", "/"})
```

更多限流策略的配置，可以进一步github.com/didip/tollbooth/limiter 了解。

**四：RESTful API接口**

前面说了在gin里面可以采用分组来组织访问URL，这里RESTful API需要给出不同的访问URL来和动态站点区分，所以新建了一个分组v1。

在浏览器中访问http://localhost/v1/user/1100000/ 


这里对v1.GET("/user/:id/*action", LimitHandler(lmt), api.GetUser) 进行了限流控制，所以如果频繁访问上面地址也将会有限制，这在API接口中非常有作用。

通过 api这个包，来实现所有有关API的代码。在GetUser函数中，通过读取mysql数据库，查找到对应userid的用户信息，并通过JSON格式返回给client。

在api.GetUser中，设置了一个局部中间件：

```Go
	//CORS 局部CORS，可在路由中设置全局的CORS
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
```
gin关于参数的处理，api包中api.go文件中有简单说明，限于篇幅原因，就不在此展开。这个项目的详细情况，请访问 https://github.com/ffhelicopter/tmm 了解。有关gin的更多信息，请访问 https://github.com/gin-gonic/gin ，该开源项目比较活跃，可以关注。

完整main.go代码：

```Go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/ffhelicopter/tmm/api"
	"github.com/ffhelicopter/tmm/handler"

	"github.com/gin-gonic/gin"
)

// 定义全局的CORS中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func LimitHandler(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			c.Data(httpError.StatusCode, lmt.GetMessageContentType(), []byte(httpError.Message))
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func main() {
gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 静态资源加载，本例为css,js以及资源图片
	router.StaticFS("/public", http.Dir("D:/goproject/src/github.com/ffhelicopter/tmm/website/static"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 导入所有模板，多级目录结构需要这样写
	router.LoadHTMLGlob("website/tpl/*/*")
	// 也可以根据handler，实时导入模板。

	// website分组
	v := router.Group("/")
	{

		v.GET("/index.html", handler.IndexHandler)
		v.GET("/add.html", handler.AddHandler)
		v.POST("/postme.html", handler.PostmeHandler)
	}

	// 中间件 golang的net/http设计的一大特点就是特别容易构建中间件。
	// gin也提供了类似的中间件。需要注意的是中间件只对注册过的路由函数起作用。
	// 对于分组路由，嵌套使用中间件，可以限定中间件的作用范围。
	// 大致分为全局中间件，单个路由中间件和群组中间件。

	// 使用全局CORS中间件。
	// router.Use(Cors())
	// 即使是全局中间件，在use前的代码不受影响
	// 也可在handler中局部使用，见api.GetUser

	//rate-limit 中间件
	lmt := tollbooth.NewLimiter(1, nil)
	lmt.SetMessage("服务繁忙，请稍后再试...")

	// API分组(RESTFULL)以及版本控制
	v1 := router.Group("/v1")
	{
		// 下面是群组中间的用法
		// v1.Use(Cors())

		// 单个中间件的用法
		// v1.GET("/user/:id/*action",Cors(), api.GetUser)

		// rate-limit
		v1.GET("/user/:id/*action", LimitHandler(lmt), api.GetUser)

		//v1.GET("/user/:id/*action", Cors(), api.GetUser)
		// AJAX OPTIONS ，下面是有关OPTIONS用法的示例
		// v1.OPTIONS("/users", OptionsUser)      // POST
		// v1.OPTIONS("/users/:id", OptionsUser)  // PUT, DELETE
	}

	srv := &http.Server{
		Addr:         ":80",
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 优雅Shutdown（或重启）服务
	// 5秒后优雅Shutdown服务
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt) //syscall.SIGKILL
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
	}
	log.Println("Server exiting")
}
```


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[前言](https://github.com/ffhelicopter/Go42/blob/master/README.md)

[第四十一章 网络爬虫](https://github.com/ffhelicopter/Go42/blob/master/content/42_41_crawler.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
