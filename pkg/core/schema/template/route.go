package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gooolib/errors"
	"github.com/version-1/gooo/pkg/core/schema/openapi/v3_0_0"
	"github.com/version-1/gooo/pkg/core/schema/template/route"
)

type RouteImplementsFile struct {
	Schema       *v3_0_0.RootSchema
	PackageName  string
	Dependencies []string
}

func (r RouteImplementsFile) Filename() string {
	return "internal/routes/routeimplments"
}

type RouteImplementsTemplateParams struct {
	Functions string
}

func (r RouteImplementsFile) Render() (string, error) {
	routes := extractRoutes(r.Schema)
	implString, err := renderRouteImplements(routes)
	if err != nil {
		return "", err
	}

	p := CommonPlainTemplateParams{
		Package:      r.PackageName,
		Dependencies: r.Dependencies,
		Content:      implString,
	}

	content, err := p.Render()

	res, err := pretify(r.Filename(), content)
	if err != nil {
		return "", errors.Wrap(err)
	}

	return string(res), nil
}

type RoutesFile struct {
	Schema       *v3_0_0.RootSchema
	PackageName  string
	Dependencies []string
}

type RoutesTemplateParams struct {
	Routes          string
	RouteHandlers   string
	RouteImplements string
	Dependencies    []string
}

func (r RoutesFile) Filename() string {
	return "internal/routes/routes"
}

func (r RoutesFile) Render() (string, error) {
	routes := extractRoutes(r.Schema)
	handlersString, err := renderRouteHandlers(routes)
	routeString := renderRoutes(routes)

	p := RoutesTemplateParams{
		Routes:        routeString,
		RouteHandlers: handlersString,
		Dependencies:  r.Dependencies,
	}

	var b bytes.Buffer
	tmpl := template.Must(template.New("routes").ParseFS(tmpl, "components/routes.go.tmpl"))
	if err := tmpl.ExecuteTemplate(&b, "routes.go.tmpl", p); err != nil {
		return "", errors.Wrap(err)
	}

	res, err := pretify(r.Filename(), b.String())
	if err != nil {
		return "", errors.Wrap(err)
	}

	return string(res), nil
}

type Route struct {
	Name       string
	InputType  string
	OutputType string
	Method     string
	Path       string
}

func renderRouteHandlers(routes []Route) (string, error) {
	var b bytes.Buffer
	for _, r := range routes {
		tmpl := template.Must(template.New("route").ParseFS(tmpl, "components/routehandler.go.tmpl"))
		p := struct {
			FuncName   string
			Path       string
			InputType  string
			OutputType string
		}{
			FuncName:   r.Name,
			Path:       r.Path,
			InputType:  r.InputType,
			OutputType: r.OutputType,
		}

		if err := tmpl.ExecuteTemplate(&b, "routehandler.go.tmpl", p); err != nil {
			return "", errors.Wrap(err)
		}
	}

	return b.String(), nil
}

func renderRouteImplements(routes []Route) (string, error) {
	var b bytes.Buffer
	for _, r := range routes {
		tmpl := template.Must(template.New("route").ParseFS(tmpl, "components/routeimpl.go.tmpl"))
		p := struct {
			FuncName   string
			Path       string
			InputType  string
			OutputType string
		}{
			FuncName:   r.Name,
			Path:       r.Path,
			InputType:  r.InputType,
			OutputType: r.OutputType,
		}

		if err := tmpl.ExecuteTemplate(&b, "routeimpl.go.tmpl", p); err != nil {
			return "", errors.Wrap(err)
		}
	}

	return b.String(), nil
}

func renderRoutes(routes []Route) string {
	var b bytes.Buffer
	for _, r := range routes {
		b.WriteString(fmt.Sprintf("%sHandler(),\n", r.Name))
	}
	return b.String()
}

func extractRoutes(r *v3_0_0.RootSchema) []Route {
	routes := []Route{}
	r.Paths.Each(func(path string, pathItem v3_0_0.PathItem) error {
		m := map[string]*v3_0_0.Operation{
			"Get":    pathItem.Get,
			"Post":   pathItem.Post,
			"Patch":  pathItem.Patch,
			"Put":    pathItem.Put,
			"Delete": pathItem.Delete,
		}
		for k, v := range m {
			if v == nil {
				continue
			}

			routeName := route.ResolveRouteName(k, path, v)
			if k == "Get" || k == "Delete" {
				route := Route{
					Name:       routeName,
					InputType:  "request.Void",
					OutputType: withSchemaPackageName(route.DetectOutputType(v, 200, "application/json")),
					Method:     k,
					Path:       path,
				}

				routes = append(routes, route)
			} else {
				statusCode := 200
				if k == "Post" {
					statusCode = 201
				}
				route := Route{
					Name:       routeName,
					InputType:  withSchemaPackageName(route.DetectInputType(v, "application/json")),
					OutputType: withSchemaPackageName(route.DetectOutputType(v, statusCode, "application/json")),
					Method:     k,
					Path:       path,
				}
				routes = append(routes, route)
			}
		}

		return nil
	})

	return routes
}
