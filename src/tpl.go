package main

import (
	"strings"
)

var filepath string = `config/config.go`

var tplpackage string = `package config
`
var tplimport string = `import (
	"encoding/json"
	"io/ioutil"
)

`
var tplstruct string = `type ?Config struct {
	?
}

`
var tplint string = `	int`
var tplstring string = `	string`
var tplfield string = `	*?Config`
var tplGlobal string = `Global`
var tplfieldval string = "	`json:\"?\"`"
var tplparse string = `var (
	config *GlobalConfig
)

func Config() *GlobalConfig {
	return config
}

//json格式的配置文件解析
func InitConfig(dir string) error {
	buf, err := ioutil.ReadFile(dir)
	if err != nil {
		return err
	}
	var c GlobalConfig
	if err := json.Unmarshal(buf, &c); err != nil {
		return err
	}
	config = &c
	return nil
}
`

// 将参数替代模板中的`?`
func ParseTpl(tpl string, args ...string) string{
	//log.Println(strings.Index(tpl, `?`))
	for i := range args {
		tpl = strings.Replace(tpl, `?`, args[i], 1)
	}
	return tpl
}