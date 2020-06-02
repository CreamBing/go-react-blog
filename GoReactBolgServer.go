package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	. "go-react-blog/config"
	. "go-react-blog/controller"
	. "go-react-blog/listener"
	"os"
)



func mapToSlice(m map[int]Files) []Files {
	s := make([]Files, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

func main() {
	fmt.Print("InitConfig...\r")
	checkErr("InitConfig", InitConfig())
	fmt.Print("InitConfig Success!!!\n")
	//读取遍历文件，将文件中的信息转为对象
	go SingleDirListener()
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	InnerRoute(app);
	//输出html
	// 请求方式: GET
	// 访问地址: http://localhost:8080/welcome
	//输出字符串
	// 类似于 app.Handle("GET", "/ping", [...])
	// 请求方式: GET
	// 请求地址: http://localhost:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})
	//输出json
	// 请求方式: GET
	// 请求地址: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello Iris!"})
	})
	app.Run(iris.Addr(":"+Conf.ServerPort))//8080 监听端口
}

// 检查错误
func checkErr(errMsg string, err error) {
	if err != nil {
		fmt.Printf("%s Error: %v\n", errMsg, err)
		os.Exit(1)
	}
}