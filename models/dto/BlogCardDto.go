package dto

import "github.com/kataras/iris"

type BlogCardDto struct {
	Cur int
	PageSize int
}

func (u *BlogCardDto) Bind(ctx iris.Context) error {
	ctx.ReadJSON(&u)
	return nil
}

