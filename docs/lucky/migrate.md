# 代码生成


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
  dsn: root:1234@tcp(192.168.1.200:3306)/lucky?charset=utf8mb4
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
logic/cms/internal/data/entity/cms_news.go // entity  
logic/cms/internal/data/meta/cms_news.go // 数据库表 metadata
logic/cms/internal/dao/news.go  // dao
logic/cms/internal/service/news_base.go // service
logic/cms/internal/endpoint/news_admin_endpoint.go // 端点
logic/cms/internal/endpoint/news_admin_http.go // http协议端点绑定
static/pages/cms/news/index.json // 后台页面
```

需要将端点手动注册到服务中，参考：[logic/cms/internal/endpoint/endpoint.go](https://github.com/fengjx/lucky/blob/master/logic/cms/internal/endpoint/endpoint.go)

```go
func Init(_ context.Context, httpServer *luchen.HTTPServer) {
    httpServer.Handler(
        &newsAdminHandler{},
        &userAdminHandler{},
        &subscriptionAdminHandler{},
    )
}
```
