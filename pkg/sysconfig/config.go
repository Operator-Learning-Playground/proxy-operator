package sysconfig

import (
	"fmt"
	proxyv1alpha1 "github.com/myoperator/proxyoperator/pkg/apis/proxy/v1alpha1"
	"github.com/myoperator/proxyoperator/pkg/common"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http/httputil"
	"os"
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

	err = yaml.Unmarshal(config, SysConfig1)
	if err != nil {
		return err
	}

	// 解析配置文件
	parseRule()

	return nil

}

var (
	ProxyMap     = make(map[string]*httputil.ReverseProxy)
	HostMap      = make(map[string]string)
	InitProxyMap = make(map[string]*httputil.ReverseProxy)
)

// parseRule 解析 sysConfig 文件
func parseRule() {
	for _, rule := range SysConfig1.Rules {
		splitUrl := strings.Split(rule.Path.Backend.Url, "://")
		res, _ := NewProxy(fmt.Sprintf("%s://%s", splitUrl[0], splitUrl[1]))
		ProxyMap[rule.Path.Backend.Prefix] = res
		HostMap[rule.Path.Backend.Prefix] = fmt.Sprintf("%s://%s", splitUrl[0], splitUrl[1])
		InitProxyMap[fmt.Sprintf("%s", splitUrl[1])] = res
		klog.Info(HostMap[rule.Path.Backend.Prefix], "", rule.Path.Backend.Prefix)
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
	Ip     string `yaml:"ip"`
	Port   int    `yaml:"port"`
	Params string `yaml:"params"` // FIXME: 在crd中还没实现
}

func CleanConfig() error {

	// 1. 把SysConfig1中的都删除
	// 清零后需要先更新app.yaml文件
	SysConfig1.Rules = make([]Rules, 0)
	if err := saveConfigToFile(); err != nil {
		return err
	}

	return ReloadConfig()
}

func AppConfig(proxy *proxyv1alpha1.Proxy) error {

	// 1. 需要先把SysConfig1中的都删除
	if len(SysConfig1.Rules) != len(proxy.Spec.Rules) {
		// 清零后需要先更新app.yaml文件
		SysConfig1.Rules = make([]Rules, len(proxy.Spec.Rules))
		if err := saveConfigToFile(); err != nil {
			return err
		}
	}

	// 2. 更新内存的配置
	for i, proxyPath := range proxy.Spec.Rules {
		SysConfig1.Rules[i].Path.Backend.Url = proxyPath.Path.Backend.Url
		SysConfig1.Rules[i].Path.Backend.Prefix = proxyPath.Path.Backend.Prefix
	}
	// 保存配置文件
	if err := saveConfigToFile(); err != nil {
		return err
	}

	return ReloadConfig()
}

// ReloadConfig 重载配置
func ReloadConfig() error {
	return InitConfig()
}

//saveConfigToFile 把config配置放入文件中
func saveConfigToFile() error {

	b, err := yaml.Marshal(SysConfig1)
	if err != nil {
		return err
	}
	// 读取文件
	path := common.GetWd()
	filePath := path + "/app.yaml"
	appYamlFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 644)
	if err != nil {
		return err
	}

	defer appYamlFile.Close()
	_, err = appYamlFile.Write(b)
	if err != nil {
		return err
	}

	return nil
}
