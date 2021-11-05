// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/compare": {
            "post": {
                "description": "compare two jpg image (max allowed size - 50 mb)",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "compare two jpg image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "image1 to compare",
                        "name": "image1",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "image2 to compare",
                        "name": "image2",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.CompareResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpext.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/get/{id}": {
            "get": {
                "description": "return image by its id",
                "produces": [
                    "image/jpeg"
                ],
                "summary": "get image by its id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "image id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "scale coeff",
                        "name": "scale",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/httpext.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpext.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/upload": {
            "post": {
                "description": "upload jpg image into MySQL database (max allowed size - 50 mb)",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Upload jpg image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "image to upload",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.UploadResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httpext.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/httpext.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.CompareResponse": {
            "type": "object",
            "properties": {
                "hammingDestance": {
                    "type": "integer"
                },
                "isSimilar": {
                    "type": "boolean"
                },
                "pHash1": {
                    "type": "string"
                },
                "pHash2": {
                    "type": "string"
                }
            }
        },
        "app.UploadResponse": {
            "type": "object",
            "properties": {
                "imageToken": {
                    "type": "string"
                },
                "pHash": {
                    "type": "string"
                }
            }
        },
        "httpext.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{"http", "https"},
	Title:       "Imloader Server API",
	Description: "API for images",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
