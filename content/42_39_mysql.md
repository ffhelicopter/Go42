# 《Go语言四十二章经》第三十九章 MySql数据库

作者：李骁

## 39.1 database/sql包

Go 提供了database/sql包用于对关系型数据库的访问，作为操作数据库的入口对象sql.DB，主要为我们提供了两个重要的功能：

* sql.DB 通过数据库驱动为我们提供管理底层数据库连接的打开和关闭操作.
* sql.DB 为我们管理数据库连接池

需要注意的是，sql.DB表示操作数据库的抽象访问接口, 而非一个数据库连接对象;它可以根据driver打开关闭数据库连接，管理连接池。正在使用的连接被标记为繁忙，用完后回到连接池等待下次使用。所以，如果你没有把连接释放回连接池，会导致过多连接使系统资源耗尽。

具体到某一类型的关系型数据库，需要导入对应的数据库驱动。下面以MySQL8.0为例，来讲讲怎么在Go语言中调用。

首先，需要下载第三方包：

go get github.com/go-sql-driver/mysql

在代码中导入mysql数据库驱动：

```Go
import (
   "database/sql"
   _ "github.com/go-sql-driver/mysql"
)
```

通常来说，不应该直接使用驱动所提供的方法，而是应该使用 sql.DB，因此在导入 mysql 驱动时，这里使用了匿名导入的方式(在包路径前添加 _)，当导入了一个数据库驱动后，此驱动会自行初始化并注册自己到Go的database/sql上下文中，因此我们就可以通过 database/sql 包提供的方法访问数据库了。


## 39.2 MySQL数据库操作

我们先建立表结构：

```Go
CREATE TABLE t_article_cate (
`cid` int(10) NOT NULL AUTO_INCREMENT, 
  `cname` varchar(60) NOT NULL, 
  `ename` varchar(100), 
  `cateimg` varchar(255), 
  `addtime` int(10) unsigned NOT NULL DEFAULT '0', 
  `publishtime` int(10) unsigned NOT NULL DEFAULT '0', 
  `scope` int(10) unsigned NOT NULL DEFAULT '10000', 
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0', 
  PRIMARY KEY (`cid`), 
  UNIQUE  KEY catename (`cname`)
) ENGINE=InnoDB AUTO_INCREMENT=99 DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
```

由于预编译语句(PreparedStatement)提供了诸多好处，可以实现自定义参数的查询，通常来说，比手动拼接字符串 SQL 语句高效，可以防止SQL注入攻击。

下面代码使用预编译的方式，来进行增删改查的操作，并通过事务来批量提交一批数据。

在Go语言中对数据类型要求很严格，一般查询数据时先定义数据类型，但是查询数据库中的数据存在三种可能:
存在值，存在零值，未赋值NULL 三种状态，因此可以将待查询的数据类型定义为sql.Nullxxx类型，可以通过判断Valid值来判断查询到的值是否为赋值状态还是未赋值NULL状态。如: sql.NullInt64 sql.NullString 

```Go
package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DbWorker struct {
	Dsn string
	Db  *sql.DB
}

type Cate struct {
	cid     int
	cname   string
	addtime int
	scope   int
}



func main() {
	dbw := DbWorker{Dsn: "root:123456@tcp(localhost:3306)/mydb?charset=utf8mb4"}
	// 支持下面几种DSN写法，具体看MySQL服务端配置，常见为第2种
	// user@unix(/path/to/socket)/dbname?charset=utf8
	// user:password@tcp(localhost:5555)/dbname?charset=utf8
	// user:password@/dbname
	// user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname

	dbtemp,  err := sql.Open("mysql",  dbw.Dsn)
	dbw.Db = dbtemp

	if err != nil {
		panic(err)
		return
	}
	defer dbw.Db.Close()

	// 插入数据测试
	dbw.insertData()

	// 删除数据测试
	dbw.deleteData()

	// 修改数据测试
	dbw.editData()

	// 查询数据测试
	dbw.queryData()

	// 事务操作测试
	dbw.transaction()
}
```

每次db.Query操作后，都建议调用rows.Close()。 因为 db.Query() 会从数据库连接池中获取一个连接，这个底层连接在结果集(rows)未关闭前会被标记为处于繁忙状态。当遍历读到最后一条记录时，会发生一个内部EOF错误，自动调用rows.Close(), 但如果提前退出循环，rows不会关闭，连接不会回到连接池中，连接也不会关闭，则此连接会一直被占用。 因此通常我们使用 defer rows.Close() 来确保数据库连接可以正确放回到连接池中。

插入数据：

```Go
// 插入数据，sql预编译
func (dbw *DbWorker) insertData() {
	stmt,  _ := dbw.Db.Prepare(`INSERT INTO t_article_cate (cname, addtime, scope) VALUES (?, ?, ?)`)
	defer stmt.Close()
	
	ret,  err := stmt.Exec("栏目1",  time.Now().UNIX(),  10)

	// 通过返回的ret可以进一步查询本次插入数据影响的行数
	// RowsAffected和最后插入的Id(如果数据库支持查询最后插入Id)
	if err != nil {
		fmt.Printf("insert data error: %v\n",  err)
		return
	}
	if LastInsertId,  err := ret.LastInsertId(); nil == err {
		fmt.Println("LastInsertId:",  LastInsertId)
	}
	if RowsAffected,  err := ret.RowsAffected(); nil == err {
		fmt.Println("RowsAffected:",  RowsAffected)
	}
}
```

删除数据：

```Go
// 删除数据，预编译
func (dbw *DbWorker) deleteData() {
	stmt,  err := dbw.Db.Prepare(`DELETE FROM t_article_cate WHERE cid=?`)
	ret,  err := stmt.Exec(122)
	// 通过返回的ret可以进一步查询本次插入数据影响的行数RowsAffected和
	// 最后插入的Id(如果数据库支持查询最后插入Id).
	if err != nil {
		fmt.Printf("insert data error: %v\n",  err)
		return
	}
	if RowsAffected,  err := ret.RowsAffected(); nil == err {
		fmt.Println("RowsAffected:",  RowsAffected)
	}
}
```

修改数据：

```Go
// 修改数据，预编译
func (dbw *DbWorker) editData() {
	stmt,  err := dbw.Db.Prepare(`UPDATE t_article_cate SET scope=? WHERE cid=?`)
	ret,  err := stmt.Exec(111,  123)
	// 通过返回的ret可以进一步查询本次插入数据影响的行数RowsAffected和
// 最后插入的Id(如果数据库支持查询最后插入Id).
	if err != nil {
		fmt.Printf("insert data error: %v\n",  err)
		return
	}
	if RowsAffected,  err := ret.RowsAffected(); nil == err {
		fmt.Println("RowsAffected:",  RowsAffected)
	}
}
```

查询数据：

```Go
// 查询数据，预编译
func (dbw *DbWorker) queryData() {
	// 如果方法包含Query，那么这个方法是用于查询并返回rows的。其他用Exec()
// 另外一种写法
	// rows, err := db.Query("select id, name from users where id = ?", 1) 
	stmt,  _ := dbw.Db.Prepare(`SELECT cid, cname, addtime, scope From t_article_cate where status=?`)
	//err = db.QueryRow("select name from users where id = ?", 1).Scan(&name) // 单行查询，直接处理
	defer stmt.Close()

	rows,  err := stmt.Query(0)
	defer rows.Close()
	if err != nil {
		fmt.Printf("insert data error: %v\n",  err)
		return
	}

	// 构造scanArgs、values两个slice，
// scanArgs的每个值指向values相应值的地址
	columns,  _ := rows.Columns()
	fmt.Println(columns)
	rowMaps := make([]map[string]string,  9)
	values := make([]sql.RawBytes,  len(columns))
	scans := make([]interface{},  len(columns))
	for i := range values {
		scans[i] = &values[i]
		scans[i] = &values[i]
	}
	i := 0
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scans...)

		each := make(map[string]string,  4)
    // 由于是map引用，放在上层for时，rowMaps最终返回值是最后一条。
		for i,  col := range values {
			each[columns[i]] = string(col)
		}

// 切片追加数据，索引位置有意思。不这样写就不是希望的样子。
		rowMaps = append(rowMaps[:i],  each) 
		fmt.Println(each)
		i++
	}
	fmt.Println(rowMaps)

	for i,  col := range rowMaps {
		fmt.Println(i,  col)
	}

	err = rows.Err()
	if err != nil {
		fmt.Printf(err.Error())
	}
}
```

事务处理：
db.Begin()开始事务，Commit() 或 Rollback()关闭事务。Tx从连接池中取出一个连接，在关闭之前都使用这个连接。Tx不能和DB层的BEGIN，COMMIT混合使用。

```Go
func (dbw *DbWorker) transaction() {
	tx,  err := dbw.Db.Begin()
	if err != nil {

		fmt.Printf("insert data error: %v\n",  err)
		return
	}
	defer tx.Rollback()
	stmt,  err := tx.Prepare(`INSERT INTO t_article_cate (cname, addtime, scope) VALUES (?, ?, ?)`)
	if err != nil {

		fmt.Printf("insert data error: %v\n",  err)
		return
	}

	for i := 100; i < 110; i++ {
		cname := strings.Join([]string{"栏目-",  string(i)},  "-")
		_,  err = stmt.Exec(cname,  time.Now().UNIX(),  i+20)
		if err != nil {
			fmt.Printf("insert data error: %v\n",  err)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Printf("insert data error: %v\n",  err)
		return
	}
	stmt.Close()
}
```


[目录](https://github.com/ffhelicopter/Go42/blob/master/SUMMARY.md)

[第三十八章 数据序列化](https://github.com/ffhelicopter/Go42/blob/master/content/42_38_json.md)

[第四十章 LevelDB与BoltDB](https://github.com/ffhelicopter/Go42/blob/master/content/42_40_kvdb.md)



>本书《Go语言四十二章经》内容在github上同步地址：https://github.com/ffhelicopter/Go42
>
>
>虽然本书中例子都经过实际运行，但难免出现错误和不足之处，烦请您指出；如有建议也欢迎交流。
>联系邮箱：roteman@163.com