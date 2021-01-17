// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package api

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/accounts/login": {
            "post": {
                "description": "Get login token and role by providing username and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Account"
                ],
                "summary": "User login",
                "parameters": [
                    {
                        "description": "Account credentials",
                        "name": "account",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.AccountCredentialsDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.AuthenticationResponseDataDto"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/services/action/{serviceid}/{actioncommand}": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Apply a service specific action to a service instance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "Apply a service specific action to a service instance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of service",
                        "name": "serviceid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "action command",
                        "name": "actioncommand",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/services/create/{servicetype}": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Create an instance of a service from yaml",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "Create service instance from yaml",
                "parameters": [
                    {
                        "description": "Service-Yaml",
                        "name": "serviceyaml",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.ServiceYamlDto"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Type of service",
                        "name": "servicetype",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/services/info": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Get an overview over all service instances",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "Get an overview over all service instances",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ServiceInstanceDetailsOverviewDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/services/info/{serviceid}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Get details over a single service instance",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "Get details over a single service instance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of service",
                        "name": "serviceid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ServiceInstanceDetailsOverviewDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/services/update/{serviceid}": {
            "post": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Update an instance of a service from yaml",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "Update service instance from yaml",
                "parameters": [
                    {
                        "description": "Service-Yaml",
                        "name": "serviceyaml",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dtos.ServiceYamlDto"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Id of service",
                        "name": "serviceid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/services/yaml/{serviceid}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Get the yaml file for an specific service instance. Parameter serviceid has to be supplied.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "Get the yaml file for an instance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of service",
                        "name": "serviceid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ServiceYamlDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/services/{serviceid}": {
            "delete": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Delete an instance of a service",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Service"
                ],
                "summary": "Delete a service instance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of service",
                        "name": "serviceid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {},
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/servicestore/info": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Lists all possible deployable services",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Servicestore"
                ],
                "summary": "Lists all possible deployable services",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ServiceStoreOverviewDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        },
        "/servicestore/yaml/{servicetype}": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Get the default yaml file for a service-template",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Servicestore"
                ],
                "summary": "Get the default yaml for a service-template",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Type of service",
                        "name": "servicetype",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.ServiceStoreItemYamlDto"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/dtos.HTTPErrorDto"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.AccountCredentialsDto": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string",
                    "example": "masterkey"
                },
                "username": {
                    "type": "string",
                    "example": "admin"
                }
            }
        },
        "dtos.AuthenticationResponseDataDto": {
            "type": "object",
            "properties": {
                "isValid": {
                    "type": "boolean",
                    "example": true
                },
                "role": {
                    "type": "string",
                    "example": "admin"
                }
            }
        },
        "dtos.HTTPErrorDto": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "dtos.ServiceInstanceActionDto": {
            "type": "object",
            "properties": {
                "command": {
                    "type": "string",
                    "example": "cmdExpose"
                },
                "name": {
                    "type": "string",
                    "example": "Expose service"
                }
            }
        },
        "dtos.ServiceInstanceActionGroupDto": {
            "type": "object",
            "properties": {
                "actions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.ServiceInstanceActionDto"
                    }
                },
                "name": {
                    "type": "string",
                    "example": "Security"
                }
            }
        },
        "dtos.ServiceInstanceDetailsDto": {
            "type": "object",
            "properties": {
                "actionGroups": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.ServiceInstanceActionGroupDto"
                    }
                },
                "id": {
                    "type": "string",
                    "example": "936DA01F-9ABD-4D9D-80C7-02AF85C822A8"
                },
                "name": {
                    "type": "string",
                    "example": "my_kibana_instance_1"
                },
                "namespace": {
                    "type": "string",
                    "example": "user_namespace_42"
                },
                "status": {
                    "type": "string",
                    "example": "ok"
                },
                "type": {
                    "type": "string",
                    "example": "kibana"
                }
            }
        },
        "dtos.ServiceInstanceDetailsOverviewDto": {
            "type": "object",
            "properties": {
                "services": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.ServiceInstanceDetailsDto"
                    }
                }
            }
        },
        "dtos.ServiceStoreItemDto": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Dashboard for elasticsearch"
                },
                "imageBase64": {
                    "type": "string",
                    "example": "data:image/png;base64,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"
                },
                "name": {
                    "type": "string",
                    "example": "kibana"
                }
            }
        },
        "dtos.ServiceStoreItemYamlDto": {
            "type": "object",
            "properties": {
                "yaml": {
                    "type": "string",
                    "example": "item: text"
                }
            }
        },
        "dtos.ServiceStoreOverviewDto": {
            "type": "object",
            "properties": {
                "services": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.ServiceStoreItemDto"
                    }
                }
            }
        },
        "dtos.ServiceYamlDto": {
            "type": "object",
            "properties": {
                "yaml": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
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
	Version:     "1.0",
	Host:        "127.0.0.1:8080",
	BasePath:    "/api/v1",
	Schemes:     []string{"https"},
	Title:       "Operator Automation Backend API",
	Description: "Operator Automation Backend API overview.",
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
	swag.Register(swag.Name, &s{})
}
