package Beluga

import (
	"github.com/go-chi/render"
	"html/template"
	"net/http"
)

type PageHeader struct {
	Title string
}

type PageData struct {
	PageHeader PageHeader
}

type Templater struct {
	layout       string
	PageData     PageData
	partials     []string
	components   []string
	templateName string
	templateList []string
}

func NewTemplater(layout string) *Templater {
	return &Templater{
		layout:       "web/layouts/" + layout + ".go.html",
		partials:     []string{"web/partials/page.go.html"},
		components:   []string{},
		templateName: "Main",
	}
}

func (templater *Templater) AddPartial(partial string) *Templater {
	templater.partials = append(templater.partials, "web/partials/"+partial+".go.html")
	return templater
}

func (templater *Templater) AddComponent(component string) *Templater {
	templater.components = append(templater.components, "web/components/"+component+".go.html")
	return templater
}

func (templater *Templater) SetTemplateName(templateName string) *Templater {
	templater.templateName = templateName
	return templater
}

func (templater *Templater) SetPageData(pageData PageData) *Templater {
	templater.PageData = pageData
	return templater
}

func (templater *Templater) Freeze() *Templater {
	templater.templateList = append(templater.partials, templater.components...)
	templater.templateList = append(templater.templateList, templater.layout)
	return templater
}

func (templater *Templater) Render(writer http.ResponseWriter, request *http.Request, data interface{}) {
	tmpl, err := template.ParseFiles(templater.templateList...)
	if err != nil {
		render.Render(writer, request, ErrNotFound)
		return
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)

	if request.Header.Get("HX-Request") == "true" {
		err = tmpl.ExecuteTemplate(writer, templater.templateName, data)
	} else {
		templateData := struct {
			PageHeader PageHeader
			Data       interface{}
		}{
			PageHeader: templater.PageData.PageHeader,
			Data:       data,
		}

		err = tmpl.ExecuteTemplate(writer, templater.templateName, templateData)
	}

	if err != nil {
		render.Render(writer, request, ErrInternalServerError)
		return
	}
}
