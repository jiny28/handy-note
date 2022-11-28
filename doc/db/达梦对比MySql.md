## 达梦特性
达梦有一个表空间概念，每个表空间绑定一个用户，用户可以直接使用自己表，使用其他用户的表要加上用户名

## 达梦和myslq区别

不支持 on duplicate key update
```
使用 merge into 代替 
```
不支持 ignore，达梦使用方法如下
```
insert ignore  into 
```
不支持 replace into
```
使用 merge into 代替 
```
不支持 auto_increment， 使用 identity 代替，代码建表用 auto_increment就不可以用图形化工具修改表结构
```
identity(1, 1)，自增长从 1 开始，每次增 1
```
创建表的时候，不支持在列的后面直接加 comment 注释，使用 COMMENT ON IS 代替，如：
```
COMMENT ON TABLE xxx IS xxx
COMMENT ON COLUMN xxx IS xxx
```
不支持 case-when-then-else，例子
```
select case  when id = 2 then "aaa" when id = 3 then "bbb" else "ccc" end as test
from (select id from person) tt;
```
不支持 longtext、TINYBLOB、MEDIUMBLOB、LONGBLOB类型
```
可用 CLOB 代替
```
不支持 date_sub 函数，使用 dateadd(datepart,n,date) 代替
```
其中，datepart可以为：
year(yy,yyyy)，quarter(qq,q)，month(mm,m)，dayofyear(dy,y)，day(dd,d)，week(wk,ww)，weekday(dw)，hour(hh) , minute(mi,n)， second(ss,s)， millisecond(ms)
例子：
select dateadd(month, -6, now());
select dateadd(month, 2, now()); 
```
不支持 substring_index 函数， 使用 substr / substring 代替，
```
substr(char[,m[,n]])
substring(char[from m[ for n]]) 
```
不支持 group_concat 函数，使用 wm_concat 代替
```
select wm_concat(id) as idstr from persion ORDER BY id ;
```

不支持 from_unixtime 函数，使用 round 代替
```
round(date[,format])
```
current_timestamp 的返回值带有时区
```
select current_timestamp();
 2018-12-17 14:34:18.433839 +08:00
```
convert(type, value) 函数
```
与 mysql 的 convert 一样，但是参数是反过来的，mysql 是 convert(value, type)
```
达梦ddl建表  int integer   无法指定大小，ddl指定就报错

达梦sql查询列信息和mysql不同，达梦是所有表的注释，主键，基础信息 各个存在一个表里
，查询时候需要进行联查，所以封装的列信息查询和mysql不同
