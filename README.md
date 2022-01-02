## 概览

news是一个用Golang编写的自带服务的内容爬虫及静态页面生成工具。简单易用，方便配置。news可以将爬取生成的静态页面通过GitHub action自动发布到GitHub pages，也可以通过自带的web服务自行托管。

## 使用

### 编译

```shell
$ git clone git@github.com:wenkechen/news.git
$ cd news
$ go build
```

### 配置

```shell
$ cp config.yaml.example config.yaml
```

项目配置
```yaml
app:
  port: 9999 # web服务端口
  debug: false # 是否开启log（gin+gorm）
  cacher: 'file' # 缓存器 file|redis
  database: 'sqlite3' # 数据库 sqlite3|mysql
  baseUrl: '/news' # 目录域名
  pages: true # 是否为静态pages

mysql:
  dsn: 'root:root@tcp(127.0.0.1:3306)/news?charset=utf8mb4&parseTime=True&loc=Local'

sqlite3:
  dsn: "./data/news.db"

redis:
  dsn: '127.0.0.1:6379'

log:
  file: '' # There must be empty on GitHub action
```

### web服务

```yaml
app:
  port: 9999 # web服务端口
  debug: false # 是否开启log（gin+gorm）
  cacher: 'file' # 缓存器 file|redis
  database: 'sqlite3' # 数据库 sqlite3|mysql
  baseUrl: '' # 目录域名
  pages: false # 是否为静态pages
...
```

运行
```shell
$ ./news
```
or
```shell
$ ./news serve
```
### CLI

爬取最新一篇文章
```shell
$ ./news latest
```

爬取最新一页（过早文章会关闭访问）
```shell
$ ./news page
```

爬取指定文章
```shell
$ ./news url :url
```

将数据库数据缓存到缓存器中
```shell
$ ./news cache
```

### GitHub pages

配置示例见`config.yaml.example`
