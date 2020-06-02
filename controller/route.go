package controller

import (
	"github.com/kataras/iris"
	"go-react-blog/controller/http"
)

// 定义500错误处理函数
func err500(ctx iris.Context) {
	ctx.WriteString("CUSTOM 500 ERROR")
}

// 定义404错误处理函数
func err404(ctx iris.Context) {
	ctx.WriteString("CUSTOM 404 ERROR")
}


func InnerRoute(app *iris.Application) {
	app.OnErrorCode(iris.StatusInternalServerError, err500)
	app.OnErrorCode(iris.StatusNotFound, err404)
	app.Post("/blogCards", http.ActionBlogCards)
	app.Post("/blogDetail", http.ActionBlogDeatail)
}
