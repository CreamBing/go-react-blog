package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	Conf        *conf        // 静态配置
)

//定义conf类型
//类型里的属性，全是配置文件里的属性
type conf struct {
	ServerPort   string `yaml:"server_port"`
	MonitorPath string `yaml:"monitor_path"`
}

// 优先从etcd中加载配置，没有则从配置文件中加载配置
func InitConfig() error {
	var err error
	var content []byte
	content, err = ioutil.ReadFile("reactblog.yml")
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return errors.New("not found nothing config")
	}
	Conf = &conf{}
	if err := yaml.Unmarshal(content, Conf); err != nil {
		return err
	}
	fmt.Printf("static config => [%#v]\n", Conf)
	return nil
}