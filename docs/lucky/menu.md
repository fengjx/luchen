# 菜单配置

新增的后台页面，需要在`菜单管理`手动进行配置，菜单接口按照[amis](https://github.com/baidu/amis)协议返回。

- amis菜单协议：<https://baidu.github.io/amis/zh-CN/components/app>
- 菜单拉取接口
  - 请求路径：`/admin/sys/menu/fetch`
  - 实现代码：[menu_admin_endpoint.go#makeFetchEndpoint](https://github.com/fengjx/lucky/blob/dev/logic/sys/internal/endpoint/menu_admin_endpoint.go)

可以参考系统菜单配置

![菜单管理](/screenshot/lucky/menu.png)

