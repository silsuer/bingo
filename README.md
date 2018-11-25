![](http://qiniu-cdn.zhiguanapp.com/629bfc026fdad3244dea2161ebb7e62f)

[![Build Status](https://travis-ci.org/silsuer/bingo.svg?branch=master)](https://travis-ci.org/silsuer/bingo)

`bingo`实际上是一个开发脚手架，使用它可以快速构建以 `bingo-router` 为核心的开发框架

受到`Laravel`的启发，将一些网站开发过程中必备的功能内置到了框架中，开箱即用

我致力于让它有着`Golang`的速度和`Laravel`的优雅

>  重新按照模块化改了一版，搭建了文档官网，目前正在备案中 ，请稍等 ...

## 模块列表

 - [x] [bingo脚手架](https://github.com/silsuer/bingo)

 - [x] [日志模块](https://github.com/silsuer/bingo-log)

 - [x] [路由模块](https://github.com/silsuer/bingo-router)

 - [x] [数据库模块](https://github.com/silsuer/bingo-orm)

 - 基于jwt-token的权限认证模块

 - 缓存模块

 - 队列模块

 - WebSocket模块



## 安装

> bingo 是集合了多个子模块得一个工具集合，提供了一个项目结构目录，规范开发，所以不提倡使用 `go get` 的方式安装框架包，因为那样结构就无法规范，在多人协作或接管前人代码时会造成不必要的困难，这里安装方式采用 git clone 方式

1. 下载

  ```shell
      git get -v https://github.com/silsuer/bingo.git
  ```

  > 如果出现错误，请先配置命令行科学上网，可以解决大部分错误问题

2. 创建项目

  ```
    bingo create bingo-demo
  ```

  使用该命令后将在命令行中出现如下显示:

  ![](http://qiniu-cdn.zhiguanapp.com/24a006d2c7f2f52d9a345e4c2454cd7b)

  在当前目录向将出现一个 `bingo-demo` 目录，里面放置着初始化好了的项目

3. 启动开发模式

  ```go
     cd bingo-demo
     bingo run dev
  ```

  将在命令行中显示如下:

  ![](http://qiniu-cdn.zhiguanapp.com/ca12fa181c4d494640a72055a7af4cf4)

  在浏览器中输入`http://localhost:8080`,若安装成功，会出现一个小狮纸...


## 更多内容，请查看 [wiki文档](https://github.com/silsuer/bingo/wiki)（已弃用，文档官网正在备案中）
