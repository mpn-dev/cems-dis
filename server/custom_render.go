package server

import (
	"fmt"
	"html/template"
	"github.com/gin-gonic/gin/render"
)

type TemplateInitFunc func(*template.Template) *template.Template

type Render struct {
	initFunc			TemplateInitFunc
	templates    	map[string]*template.Template
	files        	map[string][]string
	debug        	bool
}

func (r *Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	r.templates[name] = tmpl
}

func (r *Render) AddFromFiles(name string, files ...string) *template.Template {
	tpl := template.New(name)
	if(r.initFunc != nil) {
		tpl = r.initFunc(tpl)
	}

	tpl, err := tpl.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	if r.debug {
		r.files[name] = files
	}

	r.Add(name, tpl)
	return tpl
}

// Instance implements gin's HTML render interface
func (r *Render) Instance(name string, data interface{}) render.Render {
	var tpl *template.Template

	if r.debug {
		tpl = r.loadTemplate(name)
	} else {
		tpl = r.templates[name]
	}

	return render.HTML{
		Template: tpl,
		Data:     data,
	}
}

func (r *Render) loadTemplate(name string) *template.Template {
	tpl := template.New(name)
	if(r.initFunc != nil) {
		tpl = r.initFunc(tpl)
	}
	tpl, err := tpl.ParseFiles(r.files[name]...)
	if err != nil {
		panic(fmt.Sprintf("Error loading template %s: %s", name, err.Error()))
	}
	return template.Must(tpl, err)
}

func NewCustomRender(initFunc TemplateInitFunc, debug bool) *Render {
	return &Render{
		initFunc:			initFunc, 
		templates:    make(map[string]*template.Template),
		files:        make(map[string][]string),
		debug:        debug,
	}
}
