# Bingo

Bingo是一款使用`httprouter`作为路由的Web全栈开发框架。

受到`Laravel`的启发，将一些网站开发过程中必备的功能内置到了框架中，开箱即用

我致力于让它有着`Golang`的速度和`Laravel`的优雅

目前正在开发中......

## 安装

-----题外话开始-------------

`Bingo`使用`glide`来管理依赖，但是`glide`在windows10 x64中存在bug，安装依赖时会报错，
网上对于这个bug的解决有好多不同的...我把这个bug改掉后也传到了项目里，如果需要的话，请把
`$GOPATH/src/github/silsuer/bingo`下的`glide.exe` 拷贝到你的项目工程目录下，在终端中即可
使用`glide get` 等命令


-----题外话结束---------

 1. 如果你像我一样正在使用`glide`作为包管理工具，在你的项目目录下，使用
 
    ```go
      glide init   // 初始化一个glide工程
      glide get github.com/silsuer/bingo  // 下载并安装bingo的源代码
      bingo init   // 初始化一个bingo项目
    ```
    
    此时你的项目中应该有`.env.yaml`,`start.go`这两个文件
    现在执行
    `go run start.go`,会出现`Bingo Running......` 字样
    
    在浏览器中输入：`localhost:12345`,看到欢迎界面，安装成功！
    
 2. 正常（我为什么要用这个词...）安装
 
    ```markdown
       go get github.com/silsuer/bingo  // 获取并安装bingo
     
       bingo init  // 初始化项目
   
       go run start.go // 运行初始化后的项目
     
     // 此时浏览器输入 localhost:12345 会出现Welcome to bingo字样。安装成功
    ```

## 路由

  `Bingo`的路由策略非常自由，基于`Httprouter`，性能强劲。
  
  随意建立一个go文件，或者就在start.go中，声明一个路由列表，然后使用`bingo.RegistRoute()`把这个路由注册进去即可
   
   ```go
       //示例：
        var Welcome = []bingo.Route{
        	{
        		Path:"/",
        		Method:bingo.GET,
        		Target: func(writer http.ResponseWriter, request *http.Request, params bingo.Params) {
        			fmt.Fprint(writer,"<h1>Welcome to Bingo!</h1>")
        		},
        	},
        	{
            		Path:"/admin", // 这是第二个路由
            		Method:bingo.POST,
            		Target:Admin,
            	},
        }
        
        // 上面注册的路由的Target的方法
       func Admin(w http.ResponseWriter,r *http.Request,_ bingo.Params)  {
       	
       }

        bingo.RegistRoute(Welcome)  // 调用这个方法，把我们上面定义的路由注册一下
   ```
   
 ####  未完待续