package console

// 这个包用于各种bingo命令
/*
 * bingo create <project-name> 创建一个文件夹，文件夹中包括一个完整的项目
 * bingo init                  在当前目录下，初始化env文件、app目录、public目录
 * bingo init -env             在当前目录下，只初始化env文件，修改env文件后，再次执行bingo init即可根据env生成目录
 * bingo commonds list         输出所有控制台命令
 * bingo make:route /path -type=get/post -target=Controller@index -middleware=auth  添加一条路由
 * bingo make:controller WebController 添加一个控制器  并且注册这个控制器
 */
