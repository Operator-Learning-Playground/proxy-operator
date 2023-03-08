package sysconfig

import (
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http/httputil"
	"sigs.k8s.io/yaml"
	"strings"
)



var SysConfig1 = new(SysConfig)

func InitConfig() error {
	// 读取yaml配置
	config, err := ioutil.ReadFile("./app.yaml")
	if err != nil {
		return err
	}

	//SysConfig = NewSysConfig()

	err = yaml.Unmarshal(config, SysConfig1)
	if err != nil {
		return err
	}

	// 解析配置文件
	ParseRule()

	return nil

}

var (
	ProxyMap = make(map[string]*httputil.ReverseProxy)
	HostMap  = make(map[string]string)
	InitProxyMap = make(map[string]*httputil.ReverseProxy)
)

func ParseRule() {

	for _, rule := range SysConfig1.Rules {
		splitUrl := strings.Split(rule.Path.Backend.Url, "://")
		fmt.Printf("%s://%s\n", splitUrl[0], splitUrl[1])
		res, _ := NewProxy(fmt.Sprintf("%s://%s", splitUrl[0], splitUrl[1]))
		ProxyMap[rule.Path.Backend.Prefix] = res
		HostMap[rule.Path.Backend.Prefix] = fmt.Sprintf("%s://%s", splitUrl[0], splitUrl[1])
		InitProxyMap[fmt.Sprintf("%s", splitUrl[1])] = res
		klog.Info(rule.Path.Backend.Prefix, HostMap[rule.Path.Backend.Prefix])
	}

}

type SysConfig struct {
	Rules  []Rules `yaml:"rules"`
	Server Server  `yaml:"server"`
}

type Rules struct {
	Path Path `yaml:"path"`
}

type Path struct {
	Backend Backend `yaml:"backend"`
}

type Backend struct {
	Prefix string `yaml:"prefix"`
	Url    string `yaml:"url"`
}

type Server struct {
	Ip   string `yaml:"ip"`
	Port int 	`yaml:"port"`
}

//func AppConfig(proxy *proxyv1alpha1.Proxy) error {
//	isEdit := false
//
//	// 更新内存的配置
//	for i, config := range SysConfig.Ingress {
//		// 能在内存找到，代表是更新
//		if config.Name == ingress.Name && config.Namespace == ingress.Namespace {
//			SysConfig.Ingress[i] = *ingress
//			isEdit = true
//			break
//		}
//	}
//
//	// 新加入的
//	if !isEdit {
//		SysConfig.Ingress = append(SysConfig.Ingress, *ingress)
//
//	}
//
//	if err := saveConfigToFile(); err != nil {
//		return err
//	}
//
//
//	return ReloadConfig()
//}

// ReloadConfig 重载配置
//func ReloadConfig() error {
//	MyRouter = mux.NewRouter()
//	return InitConfig()
//
//}

//func DeleteConfig(name, namespace string) error {
//	isEdit := false
//	for i, config := range SysConfig.Ingress {
//		if config.Name == name && config.Namespace == namespace {
//			SysConfig.Ingress = append(SysConfig.Ingress[:i], )
//			isEdit = true
//			break
//		}
//	}
//	if isEdit {
//		if err := saveConfigToFile(); err != nil {
//			return err
//		}
//		return ReloadConfig()
//	}
//
//	return nil
//}

// saveConfigToFile 把config配置放入文件中
//func saveConfigToFile() error {
//
//	b, err := yaml.Marshal(SysConfig)
//	if err != nil {
//		return err
//	}
//	// 读取文件
//	path := common.GetWd()
//	filePath := path + "/app.yaml"
//	appYamlFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 644)
//	if err != nil {
//		return err
//	}
//
//	defer appYamlFile.Close()
//	_, err = appYamlFile.Write(b)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}