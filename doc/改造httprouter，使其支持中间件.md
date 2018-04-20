# 改造httprouter使其支持中间件

## 写在前面

`httprouter`在业界广受好评，主要就是因为它的性能

`httprouter`项目地址：[httprouter](https://github.com/julienschmidt/httprouter)

`httprouter`的原理：[点这里点这里~](http://www.okyes.me/2016/05/08/httprouter.html)

而`httprouter`默认是不支持中间件等功能的

`README`中说：
`Where can I find Middleware X?
This package just provides a very efficient request router with a few extra features. The router is just a http.Handler, you can chain any http.Handler compatible middleware before the router, for example the Gorilla handlers. Or you could just write your own, it's very easy!`
            
而我目前在写的Bingo框架也是基于`httprouter`的，所以准备对它进行改造，让他支持中间件的功能

我的项目地址[silsuer/bingo](https://github.com/silsuer/bingo)

## 开始改造

 1. 原理

       查看httprouter的源代码，可以看到，它使用一个前缀树来管理注册的路由，而挂载到这棵树上的，是一个`Handle`类型的方法，
       
       这个方法长这样`type Handle func(http.ResponseWriter, *http.Request, Params)`
       
       使用`httprouter.Handle(args)`方法，会将路由根据路径放置在这棵树中，
       
       在每一个http请求进来的时候，会走到`ServeHttp`方法中，这个方法就是一个多路的路由器，代码如下：
       
       ```go
           
           func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
             // 判断panic函数
            if r.PanicHandler != nil {
                defer r.recv(w, req)
            }
            
            // 开始去查找注册的路由函数
            if root := r.trees[req.Method]; root != nil {
                path := req.URL.Path
                  
                if handle, ps, tsr := root.getValue(path); handle != nil {
                     // 查找到了，执行这个函数
                    handle(w, req, ps)
                    return
                } else if req.Method != "CONNECT" && path != "/" {
                     // 没有找到，开始进行重定向或者其他操作
                    code := 301 // Permanent redirect, request with GET method
                    if req.Method != "GET" {
                        // Temporary redirect, request with same method
                        // As of Go 1.3, Go does not support status code 308.
                        code = 307
                    }
           
                    if tsr && r.RedirectTrailingSlash {
                        if len(path) > 1 && path[len(path)-1] == '/' {
                            req.URL.Path = path[:len(path)-1]
                        } else {
                            req.URL.Path = path + "/"
                        }
                        http.Redirect(w, req, req.URL.String(), code)
                        return
                    }
           
                    // Try to fix the request path
                    if r.RedirectFixedPath {
                        fixedPath, found := root.findCaseInsensitivePath(
                            CleanPath(path),
                            r.RedirectTrailingSlash,
                        )
                        if found {
                            req.URL.Path = string(fixedPath)
                            http.Redirect(w, req, req.URL.String(), code)
                            return
                        }
                    }
                }
            }
           
            // Handle 405
            if r.HandleMethodNotAllowed {
                for method := range r.trees {
                    // Skip the requested method - we already tried this one
                    if method == req.Method {
                        continue
                    }
           
                    handle, _, _ := r.trees[method].getValue(req.URL.Path)
                    if handle != nil {
                        if r.MethodNotAllowed != nil {
                            r.MethodNotAllowed(w, req)
                        } else {
                            http.Error(w,
                                http.StatusText(http.StatusMethodNotAllowed),
                                http.StatusMethodNotAllowed,
                            )
                        }
                        return
                    }
                }
            }
           
            // 如果定义了NotFound函数的话，在查找不成功的时候会执行这个函数，否则执行默认的NotFound方法
            // Handle 404
            if r.NotFound != nil {
                r.NotFound(w, req)
            } else {
                http.NotFound(w, req)
            }
           }
    
    
       ```
    
      而所谓中间件，就是在查找成功之后，首先执行中间件的方法，然后再执行handle方法，那么我们的思路就有了
      
      不再在tree上挂载handle方法，而是挂载一个我们的自定义的结构体，当查找成功的时候，先查看这个结构体是否有中间件
      
      如果有执行，如果没有，直接执行handle方法
  
 2. 自定义结构体
 
   以前我写的两篇文章里
   
   [使用Go写一个简易的MVC的Web框架](https://studygolang.com/articles/12818)
   
   [使用Go封装一个便捷的ORM](https://studygolang.com/articles/12825)
   
   也介绍过，我们的路由结构体是这样的：
   
   ```go
     type Route struct {
     	Path       string   // 路径
     	Target     Handle   // 对应的控制器路径 Controller@index 这样的方法
     	Method     string   // 访问类型 是get post 或者其他
     	Alias      string   // 路由的别名，并没有什么卵用的样子.......
     	Middleware []Handle // 中间件名称
     }
   ```
   
   其中的Handle是使用上面`httprouter`定义的方法，
   
   接下来我们改造一下这个结构体
   
   ```go
       // 上下文结构体
       type Context struct {
       	Writer  http.ResponseWriter // 响应
       	Request *http.Request       // 请求
       	Params  Params              //参数
       }
       
       type TargetHandle func(context *Context)
       
       type MiddlewareHandle func(context *Context) *Context    // 中间件需要把上下文返回回来，用来传入TargetHandle中 
       
       type Route struct {
       	Path       string   // 路径
       	Target     TargetHandle   // 要执行的方法
       	Method     string   // 访问类型 是get post 或者其他
       	Alias      string   // 路由的别名，并没有什么卵用的样子.......
       	Middleware []MiddlewareHandle // 中间件名称，在执行TargetHandle之前执行的方法
       }
   ```
   
   我们将原来Handle中的参数都装入一个上下文结构体中，然后在Route结构体中指明函数的类型
   
 3. 把Route结构体挂载到Tree上
 
    查看`httprouter`的注册路由的方法：
    
    ```go
    
       func (r *Router) Handle(method, path string, handle Handle) {
         // 判断路径的格式是否正确
       	if path[0] != '/' {
       		panic("path must begin with '/' in path '" + path + "'")
       	}
         
   	    // 如果前缀树是空的，就新建一颗
       	if r.trees == nil {
       		r.trees = make(map[string]*node)
       	}
       
       	root := r.trees[method]
       	// 如果树的根节点为空，就新建一个根节点
       	if root == nil {
       		root = new(node)
       		r.trees[method] = root
       	}
        // 根据路径，把要执行的方法挂载到树上
       	root.addRoute(path, handle)
       }

    ```
   现在我们要把其中的Handle类型的数据都改成我们自己的Route类型，
   
   很简单，代码就不贴了，想看的请看[commit变更记录](https://github.com/silsuer/bingo/commit/b651757e328b9a711ad2fe274ece8326c954d762)
   
   接下来更改整个tree文件，实际上就是把`tree`的`node`结构体中的`handle`改为route
   
   然后将tree文件中用到`handle`的地方，都用`n.route.Target`代替，虽然是无脑操作，但是要改不少行
   
   也不贴代码了...  commit记录在这里...  [commit变更记录](https://github.com/silsuer/bingo/commit/b651757e328b9a711ad2fe274ece8326c954d762)
   
   改完之后的 `ServeHttp`是这样滴~
   
   ```go

     func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
     	if r.PanicHandler != nil {
     		defer r.recv(w, req)
     	}
     
     	// 在查找之前，要先看看是否存在中间件
     	// 注册路由的时候，应该把中间件也放在此处
     
     	// 开始去查找注册的路由函数
     	if root := r.trees[req.Method]; root != nil {
     		path := req.URL.Path
     
     		if route, ps, tsr := root.getValue(path); route.Target != nil {
     			// 封装上下文
     			context := &Context{w,req,ps}
     			// 执行目标函数
     			route.Target(context)
     			return
     		} else if req.Method != "CONNECT" && path != "/" {
     			code := 301 // Permanent redirect, request with GET method
     			if req.Method != "GET" {
     				// Temporary redirect, request with same method
     				// As of Go 1.3, Go does not support status code 308.
     				code = 307
     			}
     
     			if tsr && r.RedirectTrailingSlash {
     				if len(path) > 1 && path[len(path)-1] == '/' {
     					req.URL.Path = path[:len(path)-1]
     				} else {
     					req.URL.Path = path + "/"
     				}
     				http.Redirect(w, req, req.URL.String(), code)
     				return
     			}
     
     			// Try to fix the request path
     			if r.RedirectFixedPath {
     				fixedPath, found := root.findCaseInsensitivePath(
     					CleanPath(path),
     					r.RedirectTrailingSlash,
     				)
     				if found {
     					req.URL.Path = string(fixedPath)
     					http.Redirect(w, req, req.URL.String(), code)
     					return
     				}
     			}
     		}
     	}
     
     	// Handle 405
     	if r.HandleMethodNotAllowed {
     		for method := range r.trees {
     			// Skip the requested method - we already tried this one
     			if method == req.Method {
     				continue
     			}
     
     			route, _, _ := r.trees[method].getValue(req.URL.Path)
     			if route.Target != nil {
     				if r.MethodNotAllowed != nil {
     					r.MethodNotAllowed(w, req)
     				} else {
     					http.Error(w,
     						http.StatusText(http.StatusMethodNotAllowed),
     						http.StatusMethodNotAllowed,
     					)
     				}
     				return
     			}
     		}
     	}
     
     	// Handle 404
     	if r.NotFound != nil {
     		r.NotFound(w, req)
     	} else {
     		http.NotFound(w, req)
     	}
     }
   ```
   
   可以看到，在查找到节点后，我封装了一个`Context`，接下来执行了TargetHandle方法，
   
 4.更改代码，支持中间件：
   
   ```go

             	// 封装上下文
       			context := &Context{w,req,ps}
       
       			// 判断路由是否有中间件列表，如果有，就执行
       			if len(route.Middleware)!=0{
       				for _,middleHandle:= range route.Middleware{
       					context = middleHandle(context)   // 顺序执行中间件，得到的返回结果重新注入到上下文中
       				}
       			}
       			// 执行目标函数
       			route.Target(context)
       			return
    
   ```
   
  现在，我们定义如下一个路由：
  
  ```go
    var R = []bingo.Route{
        {
            Path:   "/home",
            Method: bingo.GET,
            Target: home.Index,
            Middleware: []bingo.MiddlewareHandle{
                home.M1, home.M2,
            },
        },
    }
  ```
  其中， Index,M1,M2定义如下：
  
  ```go
    func M1(c *bingo.Context) *bingo.Context  {
    	fmt.Fprintln(c.Writer,"这是中间件1")
    	return c
    }
    
    
    func M2(c *bingo.Context) *bingo.Context  {
    	fmt.Fprintln(c.Writer,"这是中间件2")
    	return c
    }
    
    func Index(c *bingo.Context) {
    	fmt.Fprint(c.Writer,"Hello World")
    }
  ```
   然后执行`go run start.go` ,浏览器访问`localhost:12345`,就可以看到中间件执行成功的痕迹了，
   
   改造成功
   
   Bingo!  欢迎star，欢迎PR~~~~