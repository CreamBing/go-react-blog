package listener

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
)

type Files struct {
	FileName string
	Path string
}

var FILESMAP = make(map[string]Files)

func putMap(key string){
	if _, ok := FILESMAP[key]; ok {
		fmt.Printf("exist");
	}else{
		file := Files{key,key};
		FILESMAP[key] = file;
	}
}

func removeMap(key string){
	if _, ok := FILESMAP[key]; ok {
		delete(FILESMAP, key)
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
		if !fi.IsDir() { // 目录, 递归遍历
			putMap(fi.Name())
		}
	}
}

func SingleDirListener(){
	var watchDir = "/Users/zhaobing/dev/go/project/tmp"
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
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				if(event.Op==fsnotify.Create){
					putMap(event.Name)
				}else if(event.Op==fsnotify.Remove){
					removeMap(event.Name)
				}
				fmt.Printf("EVENT -> %s:%s\n", event.Op.String(), event.Name)
				fmt.Printf("EVENT -> %s\n", FILESMAP)
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(watchDir); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
