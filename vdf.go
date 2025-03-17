package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Jleagle/steam-go/steamvdf"
)

type (
	vdf struct {
		data map[string]interface{}
	}
)

func newVDF() *vdf {
	return &vdf{}
}

func (vdf *vdf) loadVDF(filename string) {
	data, err := steamvdf.ReadFile(filename)
	if err != nil {
		fmt.Printf("%s 文件读取出错！错误信息：%s", filename, err)
		os.Exit(1)
	}
	vdf.data = data.ToMapOuter()
}

func (vdf *vdf) indent(indent int, str string) string {
	out := ""
	for i := 0; i < indent; i++ {
		out += "\t"
	}
	out += str
	return out
}

func (vdf *vdf) dumpVDF(data map[string]interface{}, indent int) string {
	out := ""
	for k, v := range data {
		typ := reflect.TypeOf(v)
		if typ == reflect.TypeOf(map[string]interface{}(nil)) {
			out += vdf.indent(indent, "\""+k+"\"\n")
			out += vdf.indent(indent, "{\n")
			out += vdf.dumpVDF(v.(map[string]interface{}), indent+1)
			out += vdf.indent(indent, "}\n")
		} else {
			out += vdf.indent(indent, "\""+k+"\"\t\t\""+v.(string)+"\"\n")
		}
	}
	return out
}

func (vdf *vdf) saveVDF(filename string) {
	os.Rename(filename, filename+".bak")
	err := os.WriteFile(filename, []byte(vdf.dumpVDF(vdf.data, 0)), 0777)
	if err != nil {
		fmt.Printf("%s 文件写入失败！\n", filename)
		os.Rename(filename+".bak", filename)
		os.Exit(1)
	}
	fmt.Printf("%s 文件保存成功！原文件已备份至 %s！\n", filename, filename+".bak")
}
