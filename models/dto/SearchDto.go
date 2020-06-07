package dto


import "github.com/kataras/iris"

type SearchDto struct {
	Title string
}

func (u *SearchDto) Bind(ctx iris.Context) error {
	ctx.ReadJSON(&u)
	return nil
}
