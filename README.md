![](http://qiniu-cdn.zhiguanapp.com/629bfc026fdad3244dea2161ebb7e62f)

[![Build Status](https://travis-ci.org/silsuer/bingo.svg?branch=master)](https://travis-ci.org/silsuer/bingo)

一款使用`httprouter`作为路由的Web全栈开发框架。

受到`Laravel`的启发，将一些网站开发过程中必备的功能内置到了框架中，开箱即用

我致力于让它有着`Golang`的速度和`Laravel`的优雅


> 我们使用分支开发，主干发布的工作流，但是并不保证master一定可用，如发现问题，可提issue
> 如果有想法，欢迎联系我 Email: silsuer.liu@gmail.com

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

 1. 如果你像我一样正在使用`glide`作为包管理工具，在你的项目目录下，使用
 
    ```go
      glide init   // 初始化一个glide工程
      glide get github.com/silsuer/bingo  // 下载并安装bingo的源代码
      bingo init   // 初始化一个bingo项目
    ```
    
    此时你的项目中应该有`.env.yaml`,`start.go`这两个文件以及一些文件夹
    现在执行
    `bingo run dev`, 在浏览器中输入：`localhost:12345`,看到欢迎界面，安装成功！
    （默认使用`12345`端口，如需更改，在`start.go`中指定）
    
 2. 正常（我为什么要用这个词...）安装
 
    ```markdown
       go get github.com/silsuer/bingo  // 获取并安装bingo
     
       bingo init  // 初始化项目
   
       bingo run dev // 运行初始化后的项目
     
     // 此时浏览器输入 localhost:12345 会出现Welcome to bingo字样。安装成功
    ```

## 更多内容，请查看 [wiki文档](https://github.com/silsuer/bingo/wiki)

----------------

####  未完待续