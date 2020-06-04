package listener

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	. "go-react-blog/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Files struct {
	Id int `yaml:"id"`
	Alt string `yaml:"alt"`
	Image string `yaml:"image"`
	Title string `yaml:"title"`
	Desc string `yaml:"desc"`
	Comments bool `yaml:"comments"`
	Date string `yaml:"date"`
	Author string `yaml:"author"`
	Tags []string `yaml:"tags"`
	Path string
	Content string
}

var FILESMAP = make(map[int]Files)

func putMap(fileConf Files){
	if _, ok := FILESMAP[fileConf.Id]; ok {
		fmt.Printf("exist");
	}else{
		FILESMAP[fileConf.Id] = fileConf;
	}
}

func removeMap(id int){
	if _, ok := FILESMAP[id]; ok {
		delete(FILESMAP, id)
	}else{
		fmt.Printf("not exist");
	}
}

func readFiles(dir string){
	files,err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, fi := range files {
		if !fi.IsDir() { // 不是目录，开始读取其中的<!-- more -->
			readFile(path.Join(dir,fi.Name()))
		}
	}
}

func readFile(filePath string){
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
	}
	str := string(b)
	if len(str)>0{
		moreSplit := strings.Split(str,"<!-- more -->")
		if len(moreSplit)>1{
			confYaml := moreSplit[0]
			FileConf := &Files{}
			if err := yaml.Unmarshal([]byte(confYaml), FileConf); err != nil {
				fmt.Println("error:", err)
				return
			}
			FileConf.Path = filePath
			FileConf.Content = moreSplit[1]
			putMap(*FileConf);
		}
	}
}

func deleteFile(filePath string,createFileMap map[int]string){
	_, fileName := filepath.Split(filePath)
	if strings.HasPrefix(fileName,"."){
		fmt.Printf("hidden files")
		return
	}
	fileNameSplit := strings.Split(fileName,"-");
	if len(fileNameSplit)>1{
		if IsNum(fileNameSplit[0]){
			id,err:=strconv.Atoi(fileNameSplit[0])
			if(err!=nil){
				fmt.Println("ERROR", err)
				return
			}
			removeMap(id)
			delete(createFileMap, id)
		}
	}
}

func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func putCreateFileMap(createFileMap map[int]string,filePath string){
	_, fileName := filepath.Split(filePath)
	if strings.HasPrefix(fileName,"."){
		fmt.Printf("hidden files")
		return
	}
	fileNameSplit := strings.Split(fileName,"-");
	if len(fileNameSplit)>1{
		if IsNum(fileNameSplit[0]){
			id,err:=strconv.Atoi(fileNameSplit[0])
			if(err!=nil){
				fmt.Println("ERROR", err)
				return
			}
			if _, ok := createFileMap[id]; ok {
				fmt.Printf("exist")
			}else{
				createFileMap[id] = filePath;
			}
		}
	}
}

func SingleDirListener(){
	var watchDir = Conf.MonitorPath;
	readFiles(watchDir)
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()
	//
	done := make(chan bool)
	//
	go func() {
		createFileMap := make(map[int]string)
		for {
			timer := time.NewTimer(1*time.Second)
			select {
			// watch for events
			case event := <-watcher.Events:
				if(event.Op==fsnotify.Create){
					//这里不能直接读文件，因为新增只是一个事件，文件还没有写入完成
					putCreateFileMap(createFileMap,event.Name);
				}else if(event.Op==fsnotify.Remove){
					deleteFile(event.Name,createFileMap)
				}else if(event.Op==fsnotify.Rename){
					deleteFile(event.Name,createFileMap)
				}else if(event.Op==fsnotify.Write){
					putCreateFileMap(createFileMap,event.Name);
				}
				fmt.Printf("EVENT -> %s:%s\n", event.Op.String(), event.Name)
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			case <-timer.C:
				for key, v := range createFileMap {
					readFile(v)
					delete(createFileMap, key)
				}
			}
			timer.Stop()
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(watchDir); err != nil {
		fmt.Println("ERROR", err)
	}
	fmt.Println("monitor_path:", watchDir)
	<-done
}
