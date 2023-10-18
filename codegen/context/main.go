package main

import (
	"io/fs"
	"os"
	"regexp"
	"strings"

	nano "github.com/fumiama/NanoBot"
)

const apirestr = `(\n//\s\w+\s.+\n(//.*\n)*)func\s\(bot\s\*Bot\)\s(.*)\s\{`

var apire = regexp.MustCompile(apirestr)

func main() {
	f, err := os.Create("api_generated.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(`// Code generated by codegen/context. DO NOT EDIT.

package nano

`)
	err = fs.WalkDir(os.DirFS("./"), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasPrefix(path, "openapi_") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		f.WriteString("// 生成自文件 ")
		f.WriteString(path)
		f.WriteString("\n")
		for _, define := range apire.FindAllStringSubmatch(nano.BytesToString(data), -1) {
			f.WriteString(define[1])          // 注释
			f.WriteString("func (ctx *Ctx) ") // 函数声明
			f.WriteString(define[3])
			f.WriteString(" {\n")
			// 函数调用
			f.WriteString("\treturn ctx.caller.")
			funcname, after, _ := strings.Cut(define[3], "(")
			f.WriteString(funcname)
			f.WriteString("(")
			after, _, _ = strings.Cut(after, ")")
			paras := strings.Split(after, ", ")
			switch len(paras) {
			case 0:
			case 1:
				name, def, _ := strings.Cut(paras[0], " ")
				f.WriteString(name)
				if strings.Contains(def, "...") {
					f.WriteString("...")
				}
			default:
				for _, para := range paras[:len(paras)-1] {
					name, _, _ := strings.Cut(para, " ")
					f.WriteString(name)
					f.WriteString(", ")
				}
				name, def, _ := strings.Cut(paras[len(paras)-1], " ")
				f.WriteString(name)
				if strings.Contains(def, "...") {
					f.WriteString("...")
				}
			}
			f.WriteString(")\n}\n")
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
