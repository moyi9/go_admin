package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "title": "go_web API",
    "description": "Gin business backend starter API",
    "version": "1.0.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "schemes": ["http"],
  "paths": {}
}`

type s struct{}

func (s *s) ReadDoc() string {
	return docTemplate
}

func init() {
	swag.Register(swag.Name, &s{})
}
