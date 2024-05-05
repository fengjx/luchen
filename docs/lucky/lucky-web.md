# lucky-web 前端工程

## 项目地址

<https://github.com/fengjx/lucky-web>

前端工程是一个页面框架，包括了登录和页面渲染。绝大多数情况下，你都不需要对这个工程做修改。除非你需要定制化一些功能，例如：接入sso企业账号登录。

## 下载

### 通过 git 下载

```bash
git clone https://github.com/fengjx/lucky-web
```

### 通过release页面下载

<https://github.com/fengjx/lucky-web/releases>

## 修改环境

你可以根据自己的情况修改接口地址

## 打包部署

### 打包

前端工程打包需要安装nodejs环境，执行一下命令打包。

```bash
git clone https://github.com/fengjx/lucky-web
cd lucky-web
pnpm i
# 或者 npm i
pnmm run build
```
打包后的文件在`dist`目录下。

### 通过二进制文件运行


### 部署到文件服务器，如 nginx


### 通过 node server 运行


