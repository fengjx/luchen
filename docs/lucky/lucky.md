# lucky 后端工程

## 项目地址

<https://github.com/fengjx/lucky>

## 创建一个项目

```bash
lc start -m github.com/fengjx/lucky-demo --template lucky
```

## 工程规范

工程目录规范说明：<a href="/guide/specification" target="_blank">luchen推荐工程规范</a>

## 代码生成

以一个新闻表为例，演示通用crud代码生成。

```sql
# 表结构定义
create table if not exists cms_news
(
    `id`      bigint auto_increment primary key,
    `title`   varchar(32)  not null comment '标题',
    `content` text         not null comment '内容',
    `topic`   varchar(64)  not null default '' comment '主题',
    `status`  varchar(32)  not null default 'normal' comment '状态',
    `remark`  varchar(512) not null default '' comment '备注',
    `utime`   timestamp    not null default current_timestamp on update current_timestamp comment '更新时间',
    `ctime`   timestamp    not null default current_timestamp comment '创建时间',
    index idx_t (`topic`)
)
    engine = innodb
    default charset = utf8mb4 comment '新闻信息表';
```


配置文件`tools/gen/config.yml`
```yml
ds:
  type: mysql
  dsn: root:1234@tcp(192.168.1.200:3306)/lca2?charset=utf8mb4
target:
  custom:
    gomod: github.com/fengjx/lucky-demo
    use-admin: true
  tables:
    cms_news:
      module: cms
      simple-name: news
```
根据实际情况修改为你的数据库配置，更多参数<a href="/guide/lc" target="_blank">参考</a>

执行生成代码命令
```bash
lc migrate -c tools/gen.yml
# 或者
make migrate
```

生成文件
```bash
logic/cms/internal/dao/news.go  
logic/cms/internal/data/entity/cms_news.go 
logic/cms/internal/data/meta/cms_news.go
logic/cms/internal/endpoint/news_admin_endpoint.go
logic/cms/internal/endpoint/news_admin_http.go
logic/cms/internal/service/news_base.go
static/pages/cms/news/index.json
```




## 打包

### 编译二进制包

```bash
make build
```

默认编译为linux平台，如需要其他平台可以修改`Makefile`->`build-go`的`GOOS=linux`参数。

### docker 构建

```bash
# 构建 docker 镜像
docker build -t lucky-demo:v1.0.0 .
# 启动服务
docker run -d --name lucky  -p 8080:8080 -e APP_ENV=local lucky:v1.0.0
```

docker启动脚本在`deployments/docker/entrypoint.sh`，你可以按自己的需求调整。

