![](http://qiniu-cdn.zhiguanapp.com/629bfc026fdad3244dea2161ebb7e62f)

[![Build Status](https://travis-ci.org/silsuer/bingo.svg?branch=master)](https://travis-ci.org/silsuer/bingo)

一款使用`httprouter`作为路由的Web全栈开发框架。

受到`Laravel`的启发，将一些网站开发过程中必备的功能内置到了框架中，开箱即用

我致力于让它有着`Golang`的速度和`Laravel`的优雅



## 模块列表

 - [x] [bingo脚手架](https://github.com/silsuer/bingo)

 - [x] [日志模块](https://github.com/silsuer/bingo-log)

 - [x] [路由模块](https://github.com/silsuer/bingo-router)

 - [x] [数据库模块](https://github.com/silsuer/bingo-orm)

 - []基于jwt-token的权限认证模块

 - []缓存模块

 - []队列模块

 - []WebSocket模块



## 安装

> bingo 是集合了多个子模块得一个工具集合，提供了一个项目结构目录，规范开发，所以不提倡使用 `go get` 的方式安装框架包，因为那样结构就无法规范，在多人协作或接管前人代码时会造成不必要的困难，这里安装方式采用 git clone 方式

1. 下载

  ```shell
      git clone https://github.com/silsuer/bingo.git 
  ```
 
2. 安装依赖
  
  ```shell
     cd bingo
     glide install
  ```
  
  > 注: bingo采用 `glide` 管理依赖，在开发前请先安装 `glide`,点击此处查看 [安装方法](https://github.com/Masterminds/glide)

3. 启动开发模式

  ```go
     go run main.go
  ```
  
  访问浏览器 `http://localhost:12345`,可以看到 `hello bingo` 字样，安装成功

## 更多内容，请查看 [wiki文档](https://github.com/silsuer/bingo/wiki)

----------------

####  未完待续