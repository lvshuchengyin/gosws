// template
package gosws

import (
	"gosws/config"
	"gosws/context"
	"gosws/logger"
	"html/template"
	"path/filepath"
	"strings"
)

func unescaped(x string) interface{} {
	return template.HTML(x)
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

	tpl := template.New(filepath.Base(rawTpls[0]))
	tpl = tpl.Funcs(template.FuncMap{"unescaped": unescaped})
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
