# 简介

[lucky](https://github.com/fengjx/lucky) 是基于 [amis](https://github.com/baidu/amis) + [luchen](https://github.com/fengjx/luchen) 开发的快速开发平台。采用前后端分离架构。支持通过命令生成通用crud代码和页面，类似于`django`。接私活利器。

前端页面使用 `amis` 开发，通过json配置，无需编写前端代码即可完成页面开发。

后端使用 `luchen` 开发，集成了项目实践中的常用组件，开箱即用。

## 项目特性

- 通过 cli 命令生成通用crud代码，缩短项目交付时间
- 通过json配置实现页面定制化开发，无需编写前端代码
- 统一工程规范，良好的分层设计，代码结构清晰、易扩展
- go 语言开发，节约服务器成本

## 适合场景

1. go 开发者，你在寻找一款基于go实现类似`django`的快速开发平台。
2. 没有专职前端开发，后端兼顾前端页面开发（管理后台）。
3. 前端页面使用 [amis](https://github.com/baidu/amis) 开发，你需要学习`amis`如何使用，[快速入门文档](https://baidu.github.io/amis/zh-CN/docs/start/getting-started)。
4. 虽然[amis](https://github.com/baidu/amis)通过json配置页面，但是如果你有一些定制化需求，还是需要懂一些基础的js知识。

## 工程说明

- [lucky](https://github.com/fengjx/lucky) 后端工程
- [lucky-web](https://github.com/fengjx/lucky-web) 前端工程
- [lc](https://github.com/fengjx/lc) cli 工具

## 快速体验

演示地址：<http://admin.luchen.fun>

账号密码：lucky / luchen

## 系统截图

![截图](./login.png)

![截图](./admin-user.png)