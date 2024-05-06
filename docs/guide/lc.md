# cli 命令

## 安装

```bash
go install github.com/fengjx/lc@latest
```

`lc -h`可以查看帮助文档

## 创建项目

```bash
lc start -m github.com/fengjx/lc-demo -t httponly
```

帮助文档

```bash
$ lc start -h
NAME:
   lc start - 开始一个新项目

USAGE:
   lc start [command options] [arguments...]

OPTIONS:
   --gomod value, -m value     指定 go.mod module
   --out value, -o value       文件生成目录，默认从 gomod 读取
   --template value, -t value  使用模板，可选参数：lucky, httponly, micro (default: "httponly")
   --help, -h                  show help
```

模板说明

| 模板名称     | 说明                    |
|----------|-----------------------|
| httponly | 仅支持http协议模板           |
| micro    | 支持http、grpc协议的微服务模板   |
| lucky    | lucky 快速开发模板，附带管理后台功能 |

## 代码生成

配置文件 `gen.yml`

```yml
ds: # 数据库连接
  type: mysql
  dsn: root:1234@tcp(localhost:3306)/lca?charset=utf8mb4
target:
  custom: # 自定义配置
    tag-name: json                            # 生成 entity tag
    out-dir: ./                               # 代码输出目录
    use-admin: true                           # 是否生成管理后台页面
    gomod: github.com/fengjx/demo             # go.mod 中的 module
    var: # 自定义变量（key-value）
      foo: bar
  tables: # 数据库表  table_name: {} 自定义表配置
    sys_user:
      module: sys
      simple-name: user
```

生成代码命令

```bash
lc migrate -c gen.yml
```

配置说明

| 参数                                      | 默认值   | 说明                                                    |
|-----------------------------------------|-------|-------------------------------------------------------|
| ds.type                                 | mysql | 目前只支持 mysql                                           |
| ds.dsn                                  | -     | 数据连接                                                  |
| target.custom.tag-name                  | json  | 生成 entity 的 tag name                                  |
| target.custom.template-dir              | -     | 自定义模板路径，不指定则使用内置模板                                    |
| target.custom.out-dir                   | ./    | 生成代码目录，默认为当前路径                                        |
| target.custom.out-dir                   | ./    | 生成代码目录，默认为当前路径                                        |
| target.custom.use-admin                 | false | 是否生成基于[lucky](https://github.com/fengjx/lucky)的管理后台代码 |
| target.custom.gomod                     | -     | 项目 go module path                                     |
| target.custom.var                       | -     | 自定义参数，key-value格式，在自定义模板中可以通过`.Var.xxx`获取             |
| target.tables                           | -     | 数据库表                                                  |
| target.tables.${table_name}.use-admin   | false | 与参数`target.custom.use-admin`相同，优先级更高                  |
| target.tables.${table_name}.module      | -     | 模块名称                                                  |
| target.tables.${table_name}.simple-name | -     | 数据库表在模块内的简称，如：sys_user表，simple-name可以使用user，默认与表名相同   |
| target.tables.${table_name}.Var         | -     | i自定义表参数，key-value格式，在自定义模板中可以通过`.TableOpt.Var.xxx`获取  |

自定义函数，可以在自定义模板中使用，源码：<https://github.com/fengjx/lc/blob/dev/commands/migrate/migrate.go#L226>

| 函数名         | 说明                      |
|-------------|-------------------------|
| FirstUpper  | 首字母转大写                  |
| FirstLower  | 首字母转小写                  |
| SnakeCase   | 驼峰转下划线                  |
| TitleCase   | 下划线转驼峰                  |
| GonicCase   | go 风格的驼峰命名，如：HTTPServer |
| IsLastIndex | 是否是最后一个元素               |
| Add         | 数字相加                    |
| Sub         | 数字相减                    |

自定义函数如何使用可以参考内置模板：<https://github.com/fengjx/lc/tree/dev/commands/migrate/template>
