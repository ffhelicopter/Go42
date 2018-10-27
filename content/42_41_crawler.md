# 《Go语言四十二章经》第四十一章 网络爬虫

作者：李骁

## 41.1 go-colly网络爬虫框架

go-colly是用Go实现的网络爬虫框架。go-colly快速优雅，在单核上每秒可以发起1K以上请求；以回调函数的形式提供了一组接口，可以实现任意类型的爬虫。

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

首先，下载安装第三方包：

go get -u github.com/gocolly/colly/...

Colly大致的使用说明：

在代码中导入包：

```Go
import "github.com/gocolly/colly"
```

colly的主体是Collector对象，管理网络通信和负责在作业运行时执行附加的回掉函数。使用colly需要先初始化Collector：

```Go
c := colly.NewCollector() 
```

可以向colly附加各种不同类型的回调函数，来控制收集作业或获取信息：
回调函数的调用顺序如下：

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

下面是官方提供的抓取例子：

```Go
package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://go-colly.org/")
}
```

程序输出：

```Go
Visiting http://go-colly.org/
Visiting http://go-colly.org/docs/
Visiting http://go-colly.org/articles/
Visiting http://go-colly.org/services/
Visiting http://go-colly.org/datasets/
......
```

## 41.2 goquery HTML解析

colly框架可以快速发起请求，接收服务器响应。但如果我们需要分析返回的HTML代码，这时候仅仅使用colly就有点吃力。而goquery库是一个使用go语言写成的HTML解析库，功能更加强大。goquery将jQuery的语法和特性引入进来，所以可以更灵活地选择采集内容的数据项，就像jQuery那样的方式来操作DOM文档，使用起来非常的简便。

goquery主要的结构：

```Go
type Document struct {
	*Selection
	Url      *url.URL
	rootNode *html.Node
}
```

Document 继承了Selection 类型，因此，Document 可以直接使用 Selection 类型的方法。我们可以通过下面四种方式来初始化得到*Document对象。

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

这些方法的具体使用，我们可以借鉴jQuery的特性，多加练习就能很快掌握。限于篇幅，这里只是简单介绍了goquery的大概情况。goquery可以直接发送url请求，获得响应后得到HTML代码，这里就不举例了。由于colly这方面的功能更为强大，下面以colly作为爬虫框架，goquery作为HTML分析器，看看怎么抓取并分析页面内容：

首先，要确定已经下载安装这个第三方包：

go get github.com/PuerkitoBio/goquery


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
	urlstr := "http://metalsucks.net"
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}
	c := colly.NewCollector()
	c.SetRequestTimeout(100 * time.Second)
	// 指定Agent信息
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.108 Safari/537.36"
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", u.Host)
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Origin", u.Host)
		r.Headers.Set("Referer", urlstr)
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN, zh;q=0.9")
	})

	c.OnResponse(func(resp *colly.Response) {
		// 读取url内容 colly读取的内容传入给goquery
		htmlDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
		if err != nil {
			log.Fatal(err)
		}

		// 找到抓取项
		htmlDoc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) {
			band := s.Find("a").Text()
			title := s.Find("i").Text()
			fmt.Printf("Review %d: %s - %s\n", i, band, title)
		})
	})
	c.OnError(func(resp *colly.Response, errHttp error) {
		err = errHttp
	})
	err = c.Visit(urlstr)
}
```

程序输出：

```Go
Review 0: Obscura - Diluvium
Review 1: Skeletonwitch - Devouring Radiant Light
Review 2: Deafheaven - Ordinary Corrupt Human Love
Review 3: Between the Buried and Me - Automata II
Review 4: Chelsea Grin - Eternal Nightmare
```

上面代码中，主要是：htmlDoc.Find(".sidebar-reviews article .content-block").Each(func(i int, s *goquery.Selection) ) {}
起到了分析作用，而Find()的使用语法，是不是有些熟悉的感觉，没错就是jQuery的样子。

在goquery中，大概有以下选择器:

HTML Element 元素的选择器
ID 选择器
Element ID 选择器 Find(element#id)
Class选择器   Find(".class")
属性选择器

选择器 	说明
Find(“div[lang]“) 	筛选含有lang属性的div元素
Find(“div[lang=zh]“) 	筛选lang属性为zh的div元素
Find(“div[lang!=zh]“) 	筛选lang属性不等于zh的div元素
Find(“div[lang¦=zh]“) 	筛选lang属性为zh或者zh-开头的div元素
Find(“div[lang*=zh]“) 	筛选lang属性包含zh这个字符串的div元素
Find(“div[lang~=zh]“) 	筛选lang属性包含zh这个单词的div元素，单词以空格分开的
Find(“div[lang$=zh]“) 	筛选lang属性以zh结尾的div元素，区分大小写
Find(“div[lang^=zh]“) 	筛选lang属性以zh开头的div元素，区分大小写

parent>child选择器
如果我们想筛选出某个元素下符合条件的子元素，我们就可以使用子元素筛选器，它的语法为Find("parent>child"),表示筛选parent这个父元素下，符合child这个条件的最直接（一级）的子元素。

prev+next相邻选择器
假设我们要筛选的元素没有规律，但是该元素的上一个元素有规律，我们就可以使用这种下一个相邻选择器来进行选择。

prev~next选择器
有相邻就有兄弟，兄弟选择器就不一定要求相邻了，只要他们共有一个父元素就可以。

Colly + goquery 是抓取网络内容的利器，使用上极其方便。如今动态渲染的页面越来越多，爬虫们或多或少都需要用到headless browser来渲染待爬取的页面，这里推荐chromedp，开源网址：https://github.com/chromedp/chromedp


>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>本书《Go语言四十二章经》内容在简书同步地址：  https://www.jianshu.com/nb/29056963
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com
