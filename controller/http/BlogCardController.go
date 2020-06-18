package http

import (
	"github.com/kataras/iris"
	. "go-react-blog/listener"
	"go-react-blog/models/dto"
	"strings"
	"sync"
)

func mapToSlice(m sync.Map,params dto.BlogCardDto) ([]Files,int) {
	size := 0
	m.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	s := make([]Files, 0, size)
	if size>0{
		if params.Cur>0 && params.PageSize>0{
			start := params.Cur*(size - params.PageSize*params.Cur)+1
			end := size-(params.Cur-1)*params.PageSize
			if(end>0&&start>0&&end>=start){
				for i:=end;i>=start;i--{
					vv,_ := m.Load(i)
					if(vv!=nil){
						s = append(s, vv.(Files))
					}
				}
			}
			if(end>0&&start<=0){
				for i:=end;i>=1;i--{
					vv,_ := m.Load(i)
					if(vv!=nil){
						s = append(s, vv.(Files))
					}

				}
			}
		}
	}
	return s,size;
}

type Articles struct {
	Total int
	Dirs []Files
}

func ActionBlogCards(ctx iris.Context) {
	var params dto.BlogCardDto
	params.Bind(ctx)
	filesmap := FILESMAP
	dirs,size := mapToSlice(filesmap,params);
	articles := Articles{Total: size,Dirs: dirs}
	ctx.JSON(iris.Map{"code": 200, "data": articles})
}



func ActionBlogDeatail(ctx iris.Context) {
	var params dto.BlogDetailDto
	params.Bind(ctx)
	filesmap := FILESMAP
	vv,_ := filesmap.Load(params.Id)
	ctx.JSON(iris.Map{"code": 200, "data": vv})
}

func ActionSearch(ctx iris.Context) {
	var params dto.SearchDto
	params.Bind(ctx)
	filesmap := FILESMAP
	search(filesmap,params)
	ctx.JSON(iris.Map{"code": 200, "data": search(filesmap,params)})
}

type SearchVo struct {
	Id int
	Title string
}

func search(m sync.Map,params dto.SearchDto) []SearchVo{
	var result []SearchVo
	m.Range(func(k, v interface{}) bool {
		value := v.(Files);
		if(strings.Contains(value.Title,params.Title)){
			result=append(result, SearchVo{Id: value.Id,Title: value.Title});
		}
		return true
	})
	return result
}
