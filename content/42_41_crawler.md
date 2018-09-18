# 《Go语言四十二章经》第四十一章 网络爬虫

作者：李骁

## 41.1 go-colly

go-colly是用Go实现的网络爬虫框架。go-colly快速优雅，在单核上每秒可以发起1K以上请求；以回调函数的形式提供了一组接口，可以实现任意类型的爬虫。


Colly 特性：

清晰的API
快速（单个内核上的请求数大于1k）
管理每个域的请求延迟和最大并发数
自动cookie 和会话处理
同步/异步/并行抓取
高速缓存
自动处理非Unicode的编码
Robots.txt 支持

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

41.2 goquery 
colly框架配合goquery库，功能更加强大。goquery将jQuery的语法和特性引入到了Go语言中，可以更灵活地选择采集内容的数据项。
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
Colly + goquery 是抓取网络内容的利器，使用上极其方便。如今动态渲染的页面越来越多，爬虫们或多或少都需要用到headless browser来渲染待爬取的页面，这里推荐chromedp，开源网址：https://github.com/chromedp/chromedp