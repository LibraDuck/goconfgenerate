# gcg

[![Build Status](https://camo.githubusercontent.com/6295f2374bc4cfaab78fb2a9bf1992fcdbb3c7b9/68747470733a2f2f7472617669732d63692e6f72672f66617465646965722f6672702e7376673f6272616e63683d6d6173746572)](https://travis-ci.org/fatedier/frp)

goconfgenerate (gcg) 是一个配置文件解析代码的自动生成工具。可根据配置文件，自动生成配置文件的 golang 解析代码，解放生产力。

## 使用方法

``` golang
./gcg --path ./example.conf
```

会根据`example.conf`生成解析代码，生成目录：`/config/config.go`

## 版本说明

* 目前只支持配置文件中 int 与 string 类型，使用频率较低的类型请手动更改或添加
* 计划新增 float/int32/int64	
* 请勿出现父级结构命名重复，如`volume`字段（`msg`字段的重复是允许的）：

	"docker": {
		"container_num": 10,
		"msg":"",
		"volume": {
		    "img": "/xcdata/www/Images/"
		}
	}
	
	"demo": {
		"msg":"",
		"volume": {
			"test":""
		}
	}
