package web

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"time"
)

func (api *API) parseTemplates() *template.Template {
	funcmap := template.FuncMap{
		"render": api.renderFunc(context.Background()),

		"timestamp": func(t time.Time) template.HTML {
			dateFmt := "2006-01-02 15:04"
			rfc := t.Format(time.RFC3339)
			return template.HTML(
				fmt.Sprintf(`<time datetime="%s">%s</time>`, rfc, t.Format(dateFmt)))
		},
		"date": func(t time.Time) template.HTML {
			dateFmt := "2006-01-02"
			rfc := t.Format(time.RFC3339)
			return template.HTML(
				fmt.Sprintf(`<time datetime="%s">%s</time>`, rfc, t.Format(dateFmt)))
		},
		"shortdate": func(t time.Time) string {
			layout := "01-02"
			now := time.Now()
			if t.Year() != now.Year() && now.Sub(t) >= 4*30*24*time.Hour {
				layout = "2006-01-02"
			}
			return t.Format(layout)
		},
		"days": func(d time.Duration) string {
			days := d.Round(24*time.Hour).Hours() / 24
			return fmt.Sprintf("%g", days)
		},
	}

	t := template.New("root").Funcs(funcmap)
	root := "templates"
	err := fs.WalkDir(api.templateFS, root, func(path string, info fs.DirEntry, err error) error {

		fmt.Println("path:", path)

		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// if !strings.HasSuffix(path, ".gohtml") {
		// 	return nil
		// }

		text, err := fs.ReadFile(api.templateFS, path)
		if err != nil {
			return err
		}

		name := path
		tmpl, err := template.New(name).Funcs(funcmap).Parse(string(text))
		if err != nil {
			return err
		}
		t.AddParseTree(name, tmpl.Tree)

		return nil
	})
	if err != nil {
		panic(err)
	}

	return t
}

func (api *API) getTemplate(ctx context.Context, name string) *template.Template {
	templates := api.templates.Load().(*template.Template)
	t := templates.Lookup("templates/" + name + ".gohtml")
	if t == nil {
		panic("no template: " + name)
	}

	t, err := t.Clone()
	if err != nil {
		panic(err)
	}
	t.Funcs(api.templateFuncs(ctx))
	return t
}

func (api *API) templateFuncs(ctx context.Context) template.FuncMap {
	m := make(template.FuncMap)

	m["render"] = api.renderFunc(ctx)
	return m
}

func (api *API) renderFunc(ctx context.Context) func(string, interface{}) (template.HTML, error) {
	return func(name string, data interface{}) (template.HTML, error) {
		target := api.getTemplate(ctx, name)
		if target == nil {
			return "", fmt.Errorf("render: missing template: %s", name)
		}
		var buf bytes.Buffer
		err := target.Execute(&buf, data)
		if err != nil {
			fmt.Println("ERR!!", err)
			return "", err
		}
		return template.HTML(buf.String()), nil
	}
}

func (api *API) renderTemplateString(ctx context.Context, name string, data interface{}) (string, error) {
	var buf bytes.Buffer
	tmpl := api.getTemplate(ctx, name)
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
