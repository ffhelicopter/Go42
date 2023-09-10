
[![travis](https://travis-ci.org/ffhelicopter/go42.svg?branch=master)](https://travis-ci.org/ffhelicopter/go42)
[![Go Report Card](https://goreportcard.com/badge/github.com/ffhelicopter/go42)](https://goreportcard.com/report/github.com/ffhelicopter/go42)

# [《Go语言四十二章经》](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md "《Go语言四十二章经》")

作者：ffhelicopter（李骁）  完稿时间：2018-04-15


## 进阶阅读

《Go语言四十二章经》开源电子书升级版《深入学习Go语言》，在当当，天猫，京东有售，感谢各位对此书的支持与关注!

本书适合初次学习Go语言，以及对Go语言有初步了解的开发者，读者可通过本书努力在尽量短的时间内成长为一名合格的Go语言开发者。

![go.png](https://bkimg.cdn.bcebos.com/pic/77c6a7efce1b9d16940ab8bcfddeb48f8d546419?x-bce-process=image/resize,m_lfit,w_268,limit_1/format,f_jpg)



## 前言
写《Go语言四十二章经》，纯粹是因为开发过程中碰到过的一些问题，踩到过的一些坑，感觉在Go语言学习使用过程中，有必要深刻理解这门语言的核心思维、清晰掌握语言的细节规范以及反复琢磨标准包代码设计模式，于是才有了这本书。

Go语言以语法简单、门槛低、上手快著称。但入门后很多人发现要写出地道的、遵循 Go语言思维的代码却是不易。

在刚开始学习中，我带着比较强的面向对象编程思维惯性来写代码。但后来发现，带着面向对象的思路来写Go 语言代码会很难继续写下去，或者说看了系统源代码或其他知名开源包源代码后，围绕着struct和interface来写代码会更高效，代码更美观。虽然有人认为，Go语言的strcut 和 interface 一起，配合方法，也可以理解为面向对象，这点我姑且认可，但开发中不要过意考虑这些。因为在Go 语言中，interface接口的使用将更为灵活，刻意追求面向对象，会导致你很难理解接口在Go 语言中的妙处。

作为Go语言的爱好者，在阅读系统源代码或其他知名开源包源代码时，发现大牛对这门语言的了解之深入，代码实现之巧妙优美，所以我建议你有时间多多阅读这些代码。网上有说Go大神的标准是“能理解简洁和可组合性哲学”，的确Go语言追求代码简洁到极致，而组合思想可谓借助于struct和interface两者而成为Go的灵魂。

function，method，interface，type等词语是程序员们接触比较多的关键字，但在Go语言中，你会发现，其有了更强大，更灵活的用法。当你彻底理解了Go语言相关基本概念，以及对其特点有深入的认知，当然这也这本书的目的，再假以时日多练习和实践，我相信你应该很快就能彻底掌握这门语言，成为一名出色的Gopher。

这本书适合Go语言新手来细细阅读，对于有一定经验的开发人员，也可以根据自己的情况，选择一些章节来看。

第一章到第二十六章主要讲Go语言的基础知识，其中第十七章的type，第十八章的struct，第十九章的interface，以及第二十章的方法，都是Go语言中非常非常重要的部分。

而第二十一章的协程，第二十二章的通道以及第二十三章的同步与锁，这三章在并发处理中我们通常都需要用到，需要弄清楚他们的概念和彼此间联系。

从第二十七章开始，到第三十八章，讲述了Go标准包中比较重要的几个包，可以仔细看源代码来学习大师们的编程风格。

从第三十九章开始到结尾，主要讲述了比较常用的第三方包，但由于篇幅有限，也就不展开来讲述，有兴趣的朋友可直接到相关开源项目详细了解。

最后，希望更多的人了解和使用Go语言，也希望阅读本书的朋友们多多交流。


祝各位Gopher们工作开心，愉快编码！


本书内容在github更新：https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md<br>


#### [>>>开始阅读 第一章 Go安装与运行](https://github.com/ffhelicopter/Go42/blob/master/content/42_01_install.md)




## 交流

虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。


感谢以下网友对本书提出的修改建议： Joyboo 、林远鹏、Mr_RSI、magic-joker、3lackrush、Jacky2、tanjibo、wisecsj、eternal-flame-AD、isLishude、morya、adophper、ivanberry、xjl662750、huanglizhuo、xianyunyh、荣怡、pannz、yaaaaaaaan、sidbusy、NHibiki、awkj、yufy、lazyou、 liov 、飞翔不能的翔哥、橡_皮泥、刘冲_54ac、henng、slarsar



## 更新

因为各位热心朋友的支持与鼓励，让我有了动力不断持续更新完善本书。6年多来，看着Go语言被越来越多的人接受使用，而且本书开源以来也接到数千位读者的认可，本人非常开心。

但因为种种原因，很遗憾我需要停止本书的更新了，江湖不再有我，但我热爱的Go语言将会继续伴随我（大概率是我业余用唯一编程语言）。再次感谢给以支持的各位朋友。