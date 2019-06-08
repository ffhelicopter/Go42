# 《Go语言四十二章经》第四十一章 网络爬虫

作者：李骁

## 41.1 Colly网络爬虫框架

Colly是用Go实现的网络爬虫框架。Colly快速优雅，在单核上每秒可以发起1K以上请求；以回调函数的形式提供了一组接口，可以实现任意类型的爬虫。

Colly 特性：

清晰的API
快速（单个内核上的请求数大于1k）
管理每个域的请求延迟和最大并发数
自动cookie 和会话处理
同步/异步/并行抓取
高速缓存
自动处理非Unicode的编码
支持Robots.txt
定制Agent信息
定制抓取频次

特性如此多，引无数程序员竞折腰。下面开始我们的Colly之旅：

首先，下载安装第三方包：go get -u github.com/gocolly/colly/...

接下来在代码中导入包：

```Go
import "github.com/gocolly/colly"
```

准备工作已经完成，接下来就看看Colly的使用方法和主要的用途。

colly的主体是Collector对象，管理网络通信和负责在作业运行时执行附加的回调函数。使用colly需要先初始化Collector：

```Go
c := colly.NewCollector() 
```
我们看看NewCollector，它也是变参函数，参数类型为函数类型func(*Collector)，主要是用来初始化一个&Collector{}对象。

而在Colly中有好些函数都返回这个函数类型func(*Collector)，如UserAgent(us string)用来设置UA。所以，这里其实是一种初始化对象，设置对象属性的一种模式。相比使用方法（set方法）这种传统方式来初始设置对象属性，采用回调函数的形式在Go语言中更灵活更方便：

```Go
NewCollector(options ...func(*Collector)) *Collector
UserAgent(ua string) func(*Collector)
```

一旦得到一个colly对象，可以向colly附加各种不同类型的回调函数（回调函数在Colly中广泛使用），来控制收集作业或获取信息，回调函数的调用顺序如下：

1. OnRequest
在发起请求前被调用

2. OnError
请求过程中如果发生错误被调用

3. OnResponse
收到回复后被调用

4. OnHTML
在OnResponse之后被调用，如果收到的内容是HTML

5. OnScraped
在OnHTML之后被调用

下面我们看一个例子：

```Go
package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	// NewCollector(options ...func(*Collector)) *Collector
	// 声明初始化NewCollector对象时可以指定Agent，连接递归深度，URL过滤以及domain限制等
	c := colly.NewCollector(
		//colly.AllowedDomains("news.baidu.com"),
		colly.UserAgent("Opera/9.80 (Windows NT 6.1; U; zh-cn) Presto/2.9.168 Version/11.50"))

	// 发出请求时附的回调
	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Host", "baidu.com")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", "")
		r.Headers.Set("Referer", "http://www.baidu.com")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")

		fmt.Println("Visiting", r.URL)
	})

	// 对响应的HTML元素处理
	c.OnHTML("title", func(e *colly.HTMLElement) {
		//e.Request.Visit(e.Attr("href"))
		fmt.Println("title:", e.Text)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// <div class="hotnews" alog-group="focustop-hotnews"> 下所有的a解析
		e.ForEach(".hotnews a", func(i int, el *colly.HTMLElement) {
			band := el.Attr("href")
			title := el.Text
			fmt.Printf("新闻 %d : %s - %s\n", i, title, band)
			// e.Request.Visit(band)
		})
	})

	// 发现并访问下一个连接
	//c.OnHTML(`.next a[href]`, func(e *colly.HTMLElement) {
	//	e.Request.Visit(e.Attr("href"))
	//})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("response received", r.StatusCode)
		// 设置context
		// fmt.Println(r.Ctx.Get("url"))
	})

	// 对visit的线程数做限制，visit可以同时运行多个
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		//Delay:      5 * time.Second,
	})

	c.Visit("http://news.baidu.com")
}
```

上面代码在开始处对Colly做了简单的初始化，增加UserAgent和域名限制，其他的设置可根据实际情况来设置，Url过滤，抓取深度等等都可以在此设置，也可以后运行时在具体设置。

该例只是简单说明了Colly在爬虫抓取，调度管理方面的优势，对此如有兴趣可更深入了解。大家在深入学习Colly时，可自行选择更合适的URL。

程序运行后，开始根据news.baidu.com抓取页面结果，通过OnHTML回调函数分析首页中的热点新闻标题及链接，并可不断地抓取更深层次的新链接进行访问，每个链接的访问结果我们可以通过OnHTML来进行分析，也可通过OnResponse来进行处理，例子中没有进一步展示深层链接的内容，有兴趣的朋友可以继续进一步研究。

我们来看看OnHTML这个方法的定义：

```Go
func (c *Collector) OnHTML(goquerySelector string, f HTMLCallback)
```

直接在参数中标明了 goquerySelector ，上例中我们有简单尝试。这和我们下面要介绍的goquery HTML解析框架有一定联系，我们也可以使用goquery，通过goquery 来更轻松分析HTML代码。

## 41.2 goquery HTML解析

Colly框架可以快速发起请求，接收服务器响应。但如果我们需要分析返回的HTML代码，这时候仅仅使用Colly就有点吃力。而goquery库是一个使用Go语言写成的HTML解析库，功能更加强大。

goquery将jQuery的语法和特性引入进来，所以可以更灵活地选择采集内容的数据项，就像jQuery那样的方式来操作DOM文档，使用起来非常的简便。

goquery主要的结构：

```Go
type Document struct {
	*Selection
	Url      *url.URL
	rootNode *html.Node
}
```

Document 嵌入了Selection 类型，因此，Document 可以直接使用 Selection 类型的方法。我们可以通过下面四种方式来初始化得到*Document对象。

```Go
func NewDocumentFromNode(root *html.Node) *Document 

func NewDocument(url string) (*Document, error) 

func NewDocumentFromReader(r io.Reader) (*Document, error) 

func NewDocumentFromResponse(res *http.Response) (*Document, error)
```

Selection 是重要的一个结构体，解析中最重要，最核心的方法方法都由它提供。

```Go
type Selection struct {
	Nodes    []*html.Node
	document *Document
	prevSel  *Selection
}
```

下面我们开始了解下怎么使用goquery：

首先，要确定已经下载安装这个第三方包：

go get github.com/PuerkitoBio/goquery

接下来在代码中导入包：

```Go
import "github.com/PuerkitoBio/goquery"
```

goquery的主要用法是选择器，需要借鉴jQuery的特性，多加练习就能很快掌握。限于篇幅，这里只能简单介绍了goquery的大概情况。

goquery可以直接发送url请求，获得响应后得到HTML代码。但goquery主要擅长于HTML代码分析，而Colly在爬虫抓取管理调度上有优势，所以下面以Colly作为爬虫框架，goquery作为HTML分析器，看看怎么抓取并分析页面内容：

```Go
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func main() {
	urlstr := "https://news.baidu.com"
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}
	c := colly.NewCollector()
	// 超时设定
	c.SetRequestTimeout(100 * time.Second)
	// 指定Agent信息
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		// Request头部设定
		r.Headers.Set("Host", u.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", u.Host)
		r.Headers.Set("Referer", urlstr)
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Println("title:", e.Text)
	})

	c.OnResponse(func(resp *colly.Response) {
		fmt.Println("response received", resp.StatusCode)

		// goquery直接读取resp.Body的内容
		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))

		// 读取url再传给goquery，访问url读取内容，此处不建议使用
		// htmlDoc, err := goquery.NewDocument(resp.Request.URL.String())

		if err != nil {
			log.Fatal(err)
		}

		// 找到抓取项 <div class="hotnews" alog-group="focustop-hotnews"> 下所有的a解析
		htmlDoc.Find(".hotnews a").Each(func(i int, s *goquery.Selection) {
			band, _ := s.Attr("href")
			title := s.Text()
			fmt.Printf("热点新闻 %d: %s - %s\n", i, title, band)
			c.Visit(band)
		})

	})

	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})

	err = c.Visit(urlstr)
}
```

上面代码中，goquery先通过 goquery.NewDocumentFromReader生成文档对象htmlDoc。有了htmlDoc就可以使用选择器，而选择器的目的主要是定位：htmlDoc.Find(".hotnews a").Each(func(i int, s *goquery.Selection)，找到文档中的&lt;div class="hotnews" alog-group="focustop-hotnews"&gt;。

有关选择器Find()方法的使用语法，是不是有些熟悉的感觉，没错就是jQuery的样子。

在goquery中，常用大概有以下选择器：

|选择器类型|说明|
|:--|:--|
|HTML Element |元素的选择器Find("a")|
|Element ID 选择器| Find(element#id)|
|Class选择器   |Find(".class")|
|属性选择器| |



|选择器 	|说明|
|:--|:--|
|Find(“div[lang]“) 	|筛选含有lang属性的div元素|
|Find(“div[lang=zh]“) 	|筛选lang属性为zh的div元素|
|Find(“div[lang!=zh]“) 	|筛选lang属性不等于zh的div元素|
|Find(“div[lang¦=zh]“) 	|筛选lang属性为zh或者zh-开头的div元素|
|Find(“div[lang*=zh]“) 	|筛选lang属性包含zh这个字符串的div元素|
|Find(“div[lang~=zh]“) 	|筛选lang属性包含zh这个单词的div元素，单词以空格分开的|
|Find(“div[lang$=zh]“) 	|筛选lang属性以zh结尾的div元素，区分大小写|
|Find(“div[lang^=zh]“) 	|筛选lang属性以zh开头的div元素，区分大小写|


parent>child选择器
如果我们想筛选出某个元素下符合条件的子元素，我们就可以使用子元素筛选器，它的语法为Find("parent>child"),表示筛选parent这个父元素下，符合child这个条件的最直接（一级）的子元素。

prev+next相邻选择器
假设我们要筛选的元素没有规律，但是该元素的上一个元素有规律，我们就可以使用这种下一个相邻选择器来进行选择。

prev~next选择器
有相邻就有兄弟，兄弟选择器就不一定要求相邻了，只要他们共有一个父元素就可以。

Colly + goquery 是抓取网络内容的利器，使用上极其方便。如今动态渲染的页面越来越多，爬虫们或多或少都需要用到headless browser来渲染待爬取的页面，这里推荐chromedp，开源网址：https://github.com/chromedp/chromedp


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第四十章 LevelDB与BoltDB](https://github.com/ffhelicopter/Go42/blob/master/content/42_40_kvdb.md)

[第四十二章 WEB框架(Gin)](https://github.com/ffhelicopter/Go42/blob/master/content/42_42_gin.md)


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
