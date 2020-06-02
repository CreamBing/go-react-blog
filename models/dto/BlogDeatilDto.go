package dto


import "github.com/kataras/iris"

type BlogDetailDto struct {
	Id int
}

func (u *BlogDetailDto) Bind(ctx iris.Context) error {
	ctx.ReadJSON(&u)
	return nil
}
