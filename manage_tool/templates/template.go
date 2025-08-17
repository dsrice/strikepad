package templates

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer はEcho用のテンプレートレンダラー
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer は新しいテンプレートレンダラーを作成
func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

// Render はテンプレートをレンダリング
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}