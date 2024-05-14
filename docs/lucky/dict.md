# 数据字典

在后台管理系统中，字典数据通常是指一组固定的、静态的数据，用于描述系统中的枚举值、选项或配置项。这些数据通常存储在数据库的表中，也可以在代码中以硬编码的方式定义。

字典数据通常用于以下几个方面：

- 枚举值的描述：字典数据用于描述系统中的枚举值，例如用户的状态（激活、禁用）、订单的状态（待处理、已完成）、性别（男、女）等。
- 选项的配置：字典数据用于配置系统中的选项，例如系统的设置、参数、配置项等，可以根据需要进行动态调整。
- 界面展示：字典数据可以用于在界面上展示选项的可选值，例如下拉框、单选框、复选框等，帮助用户在界面上进行选择。

## 字段定义

- 分组：对同一属性不同值的分组定义，例如：性别
- 数据值：用户代码判断
- 显示标签：用于显示

## 示例

以用户状态为例，分为：正常、封禁，可以定义如下：

| 分组              | 数据值     | 显示标签 | 
|-----------------|---------|------|
| sys_user.status | normal  | 正常   |
| sys_user.status | disable | 封禁   |

## 命名规范

为了代码生成方便，对于数据库字段的字典值，要求按照`${表名}.${字段名}`作为分组标识，如`示例`所示。

其他字典则可以自行定义，只要分组定义不重复即可。

## 使用

以用户表`sys_user`为例，

1. 需要在用户列表中显示用户状态
2. 新增、修改时需要下拉选择用户状态

如下图所示

![](/screenshot/lucky/user-list-status.png)

![](/screenshot/lucky/user-add-status.png)

页面中可以通过字典分组拉取数据

```json
// 列表显示
{
  "label": "状态",
  "name": "status",
  "width": 150,
  "quickEdit": {
    "mode": "inline",
    "type": "select",
    "source": "${options['sys_user.status']}"
  }
}
```


```json
// 新增下拉框
{
    "label": "状态",
    "type": "select",
    "name": "status",
    "source": "${options['sys_user.status']}"
}
```

完整页面代码：[/static/pages/sys/user/index.json](https://github.com/fengjx/lucky/blob/master/static/pages/sys/user/index.json)


::: tip
前端页面在页面加载时会拉取所以字典数据，并加载到数据域中。因此，页面中可以直接使用`${options['sys_user.status']}`来获取字典数据。

具体可以查看源码：<https://github.com/fengjx/lucky-web/blob/master/src/amis/index.js#L88>
:::
