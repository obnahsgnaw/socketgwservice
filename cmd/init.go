package main

import (
	"bufio"
	"fmt"
	"github.com/obnahsgnaw/pbhttp/pkg/filerp"
	"os"
	"path/filepath"
	"strings"
)

var projectPackagePrefix = "github.com/efly/"                        // 项目包名前缀
var frameworkPackage = "github.com/obnahsgnaw/socketgwservice"       // framework 包名
var frameworkApiPackage = "github.com/obnahsgnaw/socketgwserviceapi" // framework api 包名， 分开是为了不被替换
var frameworkApiVersion = "v0.2.8"                                   // framework api 版本， 会替换go.mod中的信息，如果更新此处也需要调整
var projectName = ""                                                 // 此项目名称key，不设置则从当前目录名
var projectApiName = ""                                              // 此项目的api项目名称key，不设置则从当前目录名+api
var projectServiceName = ""                                          // 项目服务key，无则为projectName+Service

func main() {
	currentDir, err := filepath.Abs(".")
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	if projectName == "" {
		_, projectName = filepath.Split(currentDir)
	}
	if projectApiName == "" {
		projectApiName = projectName + "api"
	}
	if projectServiceName == "" {
		projectServiceName = ucfirst(projectName) + "Service"
	}

	projectPackage := projectPackagePrefix + projectName
	projectApiPackage := projectPackagePrefix + projectApiName

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("replaced to: \nproject package: %s, \nproject api package: %s \ncontinue?(y/n)", projectPackage, projectApiPackage)

	scanner.Scan()
	input := scanner.Text()
	input = strings.ToLower(input)
	if input != "y" {
		return
	}

	ignore := filerp.NewSet("README.md", "out", "cmd/gen.go", "g.sum")
	fromTo := map[string]string{
		frameworkPackage:    projectPackage,
		frameworkApiPackage: projectApiPackage,
		"/gen/tcpgateway_":  "/gen/" + projectName + "_",
	}

	if err = filerp.ReplaceDir(currentDir, fromTo, ignore); err != nil {
		println(err.Error())
		return
	}
	modFromTo := map[string]string{
		projectApiPackage + " " + frameworkApiVersion: projectApiPackage + " v0.0.0",
		"// customer replace":                         "replace " + projectApiPackage + " v0.0.0 => ../../apis/" + projectApiName,
	}
	if err = filerp.ReplaceFile(filepath.Join(currentDir, "go.mod"), modFromTo); err != nil {
		println(err.Error())
		return
	}
	if err = filerp.ReplaceFile(filepath.Join(currentDir, "makefile"), map[string]string{"framework": projectName}); err != nil {
		println(err.Error())
		return
	}
	if err = filerp.ReplaceFile(filepath.Join(currentDir, "config", "config.go"), map[string]string{"FrameworkService": projectServiceName, "FrameWork": ucfirst(projectName)}); err != nil {
		println(err.Error())
		return
	}
	if err = filerp.ReplaceFile(filepath.Join(currentDir, "config.example.yaml"), map[string]string{"framework": strings.ToLower(projectName)}); err != nil {
		println(err.Error())
		return
	}
}

func ucfirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}
