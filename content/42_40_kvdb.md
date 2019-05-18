# 《Go语言四十二章经》第四十章 LevelDB与BoltDB

作者：李骁

LevelDB 和 BoltDB 都是k/v非关系型数据库。

LevelDB没有事务，LevelDB实现了一个日志结构化的merge tree。它将有序的key/value存储在不同文件的之中，通过db, _ := leveldb.OpenFile("db", nil)，在db目录下有很多数据文件，并通过“层级”把它们分开，并且周期性地将小的文件merge为更大的文件。这让其在随机写的时候会很快，但是读的时候却很慢。

这也让LevelDB的性能不可预知：但数据量很小的时候，它可能性能很好，但是当随着数据量的增加，性能只会越来越糟糕。而且做merge的线程也会在服务器上出现问题。

>LSM树而且通过批量存储技术规避磁盘随机写入问题。 LSM树的设计思想非常朴素，它的原理是把一颗大树拆分成N棵小树， 它首先写入到内存中（内存没有寻道速度的问题，随机写的性能得到大幅提升），在内存中构建一颗有序小树，随着小树越来越大，内存的小树会flush到磁盘上。磁盘中的树定期可以做merge操作，合并成一棵大树，以优化读性能。


BoltDB会在数据文件上获得一个文件锁，所以多个进程不能同时打开同一个数据库。BoltDB使用一个单独的内存映射的文件(.db)，实现一个写入时拷贝的B+树，这能让读取更快。而且，BoltDB的载入时间很快，特别是在从crash恢复的时候，因为它不需要去通过读log去找到上次成功的事务，它仅仅从两个B+树的根节点读取ID。

BoltDB支持完全可序列化的ACID事务，让应用程序可以更简单的处理复杂操作。

BoltDB设计源于LMDB，具有以下特点：

* 直接使用API存取数据，没有查询语句；
* 支持完全可序列化的ACID事务，这个特性比LevelDB强；
* 数据保存在内存映射的文件里。没有wal、线程压缩和垃圾回收；
* 通过COW技术，可实现无锁的读写并发，但是无法实现无锁的写写并发，这就注定了读性能超高，但写性能一般，适合与读多写少的场景。
* 最后，BoltDB使用Golang开发，而且被应用于influxDB项目作为底层存储。


>LMDB的全称是Lightning Memory-Mapped Database(快如闪电的内存映射数据库)，它的文件结构简单，包含一个数据文件和一个锁文件.
>LMDB文件可以同时由多个进程打开，具有极高的数据存取速度，访问简单，不需要运行单独的数据库管理进程，只要在访问数据的代码里引用LMDB库，访问时给文件路径即可。
>
>让系统访问大量小文件的开销很大，而LMDB使用内存映射的方式访问文件，使得文件内寻址的开销非常小，使用指针运算就能实现。数据库单文件还能减少数据集复制/传输过程的开销。

## 40.1 LevelDB

Go语言LevelDB的实现我们使用 github.com/syndtr/goleveldb/leveldb 包，通过go get命令下载该包后在程序中导入。

goleveldb主要有Get()，Put()等方法，可进行key/value的读取和写入，可进行事务批量Put()插入key，Delete()删除某个key。

```Go
package main

import (
	"fmt"
	"strconv"

	"crypto/md5"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var md = md5.New()

// 测试专用
func Read(db *leveldb.DB, num int) {
	var kStr string
	var haskKey string
	kStr = strconv.Itoa(num)
	md.Write([]byte(kStr))
	haskKey = fmt.Sprintf("%x", md.Sum(nil))
	md.Reset()

	db.Get([]byte(haskKey), nil)
}

// 测试专用
func Write(db *leveldb.DB, num int) {
	var kStr string
	var haskKey string
	kStr = strconv.Itoa(num)
	md.Write([]byte(kStr))
	haskKey = fmt.Sprintf("%x", md.Sum(nil))
	md.Reset()

	db.Put([]byte(haskKey), []byte(kStr), nil)
}

func main() {
	// 打开数据库文件 /path/to/db ,第一个参数为存放数据的目录，不是具体文件
	// o := &opt.Options{	Filter: filter.NewBloomFilter(10),}
	// OpenFile第2个参数这里指定为nil，在数据集大时可设置比如布隆过滤器。
	// *opt.Options 为nil默认为false ，true为只读模式ReadOnly
	db, _ := leveldb.OpenFile("levdb", nil)

	defer db.Close()

	// 读数据库:Get(key,nil)，写数据库:Put(key,value,nil)
	// Put第三个参数为nil，默认就好，默认时写的时候如果机器崩了数据会丢失。
	// key和value都是字节slice
	_ = db.Put([]byte("key1"), []byte("好好检查"), nil)
	_ = db.Put([]byte("key2"), []byte("天天向上"), nil)
	_ = db.Put([]byte("key:3"), []byte("就会一个本事"), nil)
	_ = db.Put([]byte("uname"), []byte("Jim"), nil)
	_ = db.Put([]byte("time"), []byte("1450932202"), nil)

	// 读数据库:Get(key,nil)，返回字节slice
	data, _ := db.Get([]byte("key1"), nil)
	fmt.Println("key1=>", string(data))

	// 删除某个key(key,nil)，key不存在时并不返回错误
	_ = db.Delete([]byte("key"), nil)

	//迭代数据库内容:
	iter := db.NewIterator(nil, nil)
	fmt.Println("迭代所有key/value")
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key), "=>", string(value))

	}
	iter.Release()
	iter.Error()

	//Seek()定位到比给定key值(字节值)要大的第一个key，可next迭代所有筛选出的key/value:
	iter = db.NewIterator(nil, nil)
	fmt.Println("\nSeek()按值筛选查找key")
	for ok := iter.Seek([]byte("t")); ok; ok = iter.Next() {
		// Use key/value.
		fmt.Println("Seek-then-Iterate:")
		fmt.Println(string(iter.Key()), "=>", string(iter.Value()))
	}
	iter.Release()

	//迭代内容子集:start表示key中包含有的字符串， Limit表示key不能包含有字符串
	fmt.Println("\n 按照指定（排除）条件筛选key")
	iter = db.NewIterator(&util.Range{Start: []byte("key"), Limit: []byte("no")}, nil)
	for iter.Next() {
		// Use key/value.
		fmt.Println("Iterate over subset of database content:")
		fmt.Println(string(iter.Key()), "=>", string(iter.Value()))
	}
	iter.Release()

	//迭代子集内容，key的前缀是指定字符串:
	fmt.Println("\n 查找指定前缀key")
	iter = db.NewIterator(util.BytesPrefix([]byte("key")), nil)
	for iter.Next() {
		// Use key/value.
		fmt.Println("Iterate over subset of database content with a particular prefix:")
		fmt.Println(string(iter.Key()), "=>", string(iter.Value()))
	}
	iter.Release()

	_ = iter.Error()

	//批量写:
	batch := new(leveldb.Batch)
	var kStr string
	var batchkey string
	for i := 0; i < 10; i++ {
		kStr = strconv.Itoa(i)
		md.Write([]byte(kStr))
		batchkey = fmt.Sprintf("%x", md.Sum(nil))
		batch.Put([]byte(batchkey), []byte(kStr))
	}
	md.Reset()
	batch.Delete([]byte("lazy"))
	_ = db.Write(batch, nil)
}
```

Leveldb比较突出的问题是在读操作上，在大量key的情况下可能成为性能的瓶颈，我们可以根据场景来选择使用。下面是我们进行的几种数量级别的基准测试数据：

BenchmarkWrite-4   	  100000	     14541 ns/op
BenchmarkRead-4    	  100000	     13094 ns/op

BenchmarkWrite-4   	  500000	     12724 ns/op
BenchmarkRead-4    	  500000	     17002 ns/op

BenchmarkWrite-4   	 1000000	     13355 ns/op
BenchmarkRead-4    	 1000000	     20610 ns/op

BenchmarkWrite-4   	 3000000	     15644 ns/op
BenchmarkRead-4    	 3000000	     22742 ns/op

我们可以看到随着key的数量的增加，读的性能明显地下降，而写的性能则不受影响。

## 40.2 BoltDB

Go语言BoltDB的实现我们使用 github.com/boltdb/bolt 包，通过go get命令下载该包后在程序中导入。

BoltDB中存储比较重要的概念是bucket，存取操作之前都需要指定bucket，如果读数据时指定bucket不存在，则会panic。

```Go
package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

func main() {
	Boltdb()
}
func Boltdb() error {
	// Bolt 会在数据文件上获得一个文件锁，所以多个进程不能同时打开同一个数据库。
	// 打开一个已经打开的 Bolt 数据库将导致它挂起，直到另一个进程关闭它。
	// 为防止无限期等待，您可以将超时选项传递给Open()函数：
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 10 * time.Second})
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	//	两种处理方式：读-写和只读操作，读-写方式开始于db.Update方法：
	//	Bolt 一次只允许一个读写事务，但是一次允许多个只读事务。
	// 每个事务处理都有一个始终如一的数据视图
	err = db.Update(func(tx *bolt.Tx) error {
		// 这里还有另外一层：k-v存储在bucket中，
		// 可以将bucket当做一个key的集合或者是数据库中的表。
		//（顺便提一句，buckets中可以包含其他的buckets，这将会相当有用）
		// Buckets 是键值对在数据库中的集合.所有在bucket中的key必须唯一，
		// 使用DB.CreateBucket() 函数建立buket
		// Tx.DeleteBucket() 删除bucket
		// b := tx.Bucket([]byte("MyBucket"))
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))

		//要将 key/value 对保存到 bucket，请使用 Bucket.Put() 函数：
		//这将在 MyBucket 的 bucket 中将 "answer" key的值设置为"42"。
		err = b.Put([]byte("answer"), []byte("42"))
		err = b.Put([]byte("why"), []byte("101010"))
		return err
	})

	// 可以看到，传入db.update函数一个参数，在函数内部你可以get/set数据和处理error。
	// 如返回为nil，事务就会从数据库得到一个commit，但如果返回一个实际的错误，则会做回滚，
	// 你在函数中做的事情都不会commit。这很自然，因为你不需要人为地去关心事务的回滚，
	// 只需要返回一个错误，其他的由Bolt去帮你完成。
	// 只读事务 只读事务和读写事务不应该相互依赖，一般不应该在同一个例程中同时打开。
	// 这可能会导致死锁，因为读写事务需要定期重新映射数据文件，
	// 但只有在只读事务处于打开状态时才能这样做。

	// 批量读写事务.每一次新的事物都需要等待上一次事物的结束，
	// 可以通过DB.Batch()批处理来完
	err = db.Batch(func(tx *bolt.Tx) error {
		return nil
	})

	//只读事务在db.View函数之中：在函数中可以读取，但是不能做修改。
	db.View(func(tx *bolt.Tx) error {
		//要检索这个value，我们可以使用 Bucket.Get() 函数：
		//由于Get是有安全保障的，所有不会返回错误，不存在的key返回nil
		b := tx.Bucket([]byte("MyBucket"))
		//tx.Bucket([]byte("MyBucket")).Cursor() 可这样写
		v := b.Get([]byte("answer"))
		id, _ := b.NextSequence()
		fmt.Printf("The answer is: %s %d \n", v, id)

		//游标遍历key
		c := b.Cursor()
		fmt.Println("\n游标遍历key")
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		//游标上有以下函数：
		//First()  移动到第一个健.
		//Last()   移动到最后一个健.
		//Seek()   移动到特定的一个健.
		//Next()   移动到下一个健.
		//Prev()   移动到上一个健.

		//Prefix 前缀扫描
		fmt.Println("\nPrefix 前缀扫描")
		prefix := []byte("a")
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
		}
		return nil
	})

	//如果你知道所在桶中拥有键，你也可以使用ForEach()来迭代：
	db.View(func(tx *bolt.Tx) error {
		fmt.Println("\nForEach()来迭代")
		b := tx.Bucket([]byte("MyBucket"))
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})

	//事务处理
	// 开始事务
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 使用事务...
	_, err = tx.CreateBucket([]byte("MyBucket"))
	if err != nil {
		return err
	}

	// 事务提交
	if err = tx.Commit(); err != nil {
		return err
	}
	return err

	//还可以在一个键中存储一个桶，以创建嵌套的桶：
	//func (*Bucket) CreateBucket(key []byte) (*Bucket, error)
	//func (*Bucket) CreateBucketIfNotExists(key []byte) (*Bucket, error)
	//func (*Bucket) DeleteBucket(key []byte) error
}
```

BoltDB的性能测试这里就不再做阐述，和LevelDB正好相反，它在写性能上存在瓶颈，而读性能上非常有优势，这两者我们需要根据场景来选择使用。



[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第三十九章 Mysql数据库](https://github.com/ffhelicopter/Go42/blob/master/content/42_39_mysql.md)

[第四十一章 网络爬虫](https://github.com/ffhelicopter/Go42/blob/master/content/42_41_crawler.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com