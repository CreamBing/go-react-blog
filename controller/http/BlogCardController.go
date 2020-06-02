package http

import (
	"github.com/kataras/iris"
	. "go-react-blog/listener"
	"go-react-blog/models/dto"
)

func mapToSlice(m map[int]Files,params dto.BlogCardDto) []Files {
	size := len(m)
	s := make([]Files, 0, size)
	if size>0{
		if params.Cur>0 && params.PageSize>0{
			start := params.Cur*(size - params.PageSize*params.Cur)+1
			end := size-(params.Cur-1)*params.PageSize
			if(end>0&&start>0&&end>=start){
				for i:=start;i<=end;i++{
					s = append(s, m[i])
				}
			}
			if(end>0&&start<=0){
				for i:=1;i<=end;i++{
					s = append(s, m[i])
				}
			}
		}
	}
	return s
}

type Articles struct {
	Total int
	Dirs []Files
}

func ActionBlogCards(ctx iris.Context) {
	var params dto.BlogCardDto
	params.Bind(ctx)
	filesmap := FILESMAP
	dirs := mapToSlice(filesmap,params);
	articles := Articles{Total: len(dirs),Dirs: dirs}
	ctx.JSON(iris.Map{"code": 200, "data": articles})
}



func ActionBlogDeatail(ctx iris.Context) {
	var params dto.BlogDetailDto
	params.Bind(ctx)
	filesmap := FILESMAP
	ctx.JSON(iris.Map{"code": 200, "data": filesmap[params.Id]})
}
