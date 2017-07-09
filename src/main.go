/*
	只支持配置文件中 int 与 string 类型
	计划新增 float/int32/int64

	请勿出现父级结构命名重复，如：
	```
	"docker": {
		"container_num": 10,
		"volume": {
		    "img": "/xcdata/www/Wxbot/Images/img/"
		}
    	}
    	"demo": {
    		"msg":,
    		"volume": {
    			"test":
    		}
    	}
    	```
 */
package main

import (
	"log"
	"flag"
	"os"
	"io/ioutil"
	"container/list"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	var confPath string
	flag.StringVar(&confPath, "path", "", "Please set up a configuration file path. ")
	help := flag.Bool("help", false, "help")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	parse(confPath)
}

func parse(confPath string) {
	oriByteSlice, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Println(err)
	}
	var noblanksByteSlice []byte
	for _, v := range oriByteSlice {
		if v != '\n' && v != '\t' && v != ' ' {
			noblanksByteSlice = append(noblanksByteSlice, v)
		}
	}
	log.Println("The original data:\n"+string(noblanksByteSlice))	// 原始数据

	var res []byte

	// 一次解析 begin
	// 输出：剔除字段数据，提取层级信息及类型
	stack := struct {
		Value     *list.List
		State     int  // 状态机，1、判断符号配对，2、提取字段名
		Hierarchy byte // 层
		MaxHir    byte // 最高层
	}{}
	stack.Value = list.New()
	stack.Hierarchy = 'A' - 1
	for _, v := range noblanksByteSlice {
		// 状态机
		if stack.State == 1 {
			res = append(res, v)
		}

		// Hierarchy
		if v == '{' {
			stack.Value.PushBack(v)
			stack.Hierarchy++
			if stack.MaxHir < stack.Hierarchy {
				stack.MaxHir = stack.Hierarchy
			}
			if stack.Hierarchy > 'Z' {
				log.Println("too lager")
				return
			}
		}
		if v == '}' {
			if sb := stack.Value.Back(); sb.Value.(byte) == '{' {
				//log.Println("{ok")
				stack.Value.Remove(sb)
				stack.Hierarchy--
				if stack.Hierarchy < 'A' - 1 {
					log.Println("something error")
					return
				}
			} else {
				log.Println("{err")
				return
			}
		}
		if v == '"' {
			if sb := stack.Value.Back(); sb.Value.(byte) == '"' {
				stack.Value.Remove(sb)
				stack.State = 0
				res = append(res, stack.Hierarchy)
			} else {
				stack.Value.PushBack(v)
			}
		}
		if v == '{' || v == ',' {
			stack.State = 1
		}
	}
	log.Println("First result:\n"+string(res))//result
	maxHir := stack.MaxHir
	//// 一次解析 end

	//// 二次解析 begin
	// 输出：详细层级信息及类型
	// 符号约定：
	// %s string
	// %d int
	// %t child struct
	// %后的数字代表层级
	type structSlice []byte
	root := structSlice{}        // 二次解析结果
	bak := structSlice{}        // 缓存区
	var cardinal byte = 'A'        //
	var hir byte

	for i, v := range res {
		// 往bak丢缓存，当遇到数据时将缓存放入root中并清空缓存。
		// 检查数据类型和层级，与数据同时写入root。
		if v >= 'A' && v <= 'Z' {
			// hir: 层级递增，进入child
			for hir = cardinal; hir <= maxHir; hir++ {

				if v == hir && res[i - 1] == '"' {
					root = append(root, ExtractQuotesBytes(bak)...)
					if res[i + 1] == hir {
						root = append(root, "%" + ByteToString(hir) + "s"...)
					} else {
						root = append(root, "%" + ByteToString(hir) + "d"...)
					}

					// clear bak
					bak = append(bak[:0], []byte{}...)
				}

				if v == hir + 1 && res[i - 1] == '"' {
					if root[len(root) - 1] != 't' {
						root = append(root[:len(root) - 3], "%" + ByteToString(hir) + "t"...)
						cardinal++
					}
				}
			}

		}
		// reset层级，恢复parent层级
		if v == '"' && res[i + 1] >= 'A' && res[i + 1] < cardinal {
			cardinal = 'A'
		}

		bak = append(bak, v)
	}
	log.Println("Sec result:\n"+string(root))//result
	//// 二次解析 End

	// 三次解析
	Assembly()

	bak = append(bak[:0], []byte{}...)

	structs := make(map[int]string)
	structs[1] = ParseTpl(tplstruct, tplGlobal)
	for i := 0; i < len(root); i++ {
		if root[i] == '%' {
			layer := int(root[i + 1]) - 48

			var types string
			switch root[i + 2] {
			case 'd': types = tplint
			case 's': types = tplstring
			case 't': types = ParseTpl(tplfield, ParseField(string(bak)))
				log.Println(structs[layer+1])
				structs[layer+1] = ParseTpl(structs[layer+1], ``)
				AppendToFile(filepath, string(structs[layer+1]))
				delete(structs, layer+1)
				structs[layer+1] = ParseTpl(tplstruct, ParseField(string(bak)))
			}

			field := tplfieldval
			field = ParseTpl(field, string(bak))

			structs[layer] = ParseTpl(structs[layer], ParseField(string(bak)) + types + field + `
	?`)
			bak = append(bak[:0], []byte{}...)
			i += 2 //'%'后移两个字节
		} else {
			bak = append(bak, root[i])
		}
	}
	structs[1] = ParseTpl(structs[1], ``)
	AppendToFile(filepath, string(structs[1]))

	log.Println(structs[1])

	AppendToFile(filepath, tplparse)
}

