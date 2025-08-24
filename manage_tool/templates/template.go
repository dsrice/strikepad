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
	// レイアウトとコンポーネントを最初に読み込み、その後ページテンプレートを読み込む
	tmpl := template.New("")

	// レイアウトとコンポーネントを読み込み
	template.Must(tmpl.ParseFiles(
		"templates/layout.html",
		"templates/sidebar.html",
	))

	// 全てのテンプレートを読み込み
	template.Must(tmpl.ParseGlob("templates/*.html"))

	return &TemplateRenderer{
		templates: tmpl,
	}
}

// Render はテンプレートをレンダリング
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
