// template
package gosws

import (
	"encoding/json"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/lvshuchengyin/gosws/config"
	"github.com/lvshuchengyin/gosws/context"
	"github.com/lvshuchengyin/gosws/logger"
)

func unescaped(x string) interface{} {
	return template.HTML(x)
}

func marshal(v interface{}) template.JS {
	a, _ := json.Marshal(v)
	return template.JS(a)
}

func TemplateRender(arg *context.Context, data interface{}, filePaths ...string) (err error) {
	templateDir := config.TemplatesDir()

	rawTpls := make([]string, 0, len(filePaths)+1)
	rawTpls = append(rawTpls, templateDir+"/base.html")

	for _, fp := range filePaths {
		if !strings.HasPrefix(fp, templateDir) {
			fp = templateDir + "/" + fp
		}
		rawTpls = append(rawTpls, fp)
	}

	tpl := template.New(filepath.Base(rawTpls[0])).Delims("[[", "]]")
	tpl = tpl.Funcs(template.FuncMap{"unescaped": unescaped, "marshal": marshal})
	tpl, err = tpl.ParseFiles(rawTpls...)
	if err != nil {
		logger.Error("template has err: %s, %s", filePaths, err.Error())
		return
	}

	arg.Res.Header().Set("Content-Type", "text/html")
	err = tpl.Execute(arg.Res, data)
	if err != nil {
		logger.Error("template execute err: %s, %s", filePaths, err.Error())
		return
	}

	return
}
