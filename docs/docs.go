package docs

import (
	"github.com/swaggo/swag"
)

const docTemplate = `{
    "schemes": [
        "http"
    ],
    "swagger": "2.0",
    "info": {
        "description": "Uma API REST completa em Go com suporte aos verbos HTTP",
        "title": "API Go com Arquitetura",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api",
    "paths": {
        "/produtos": {
            "get": {
                "description": "Retorna uma lista de todos os produtos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "produtos"
                ],
                "summary": "Listar todos os produtos",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/Produto"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Cria um novo produto com os dados fornecidos",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "produtos"
                ],
                "summary": "Criar um novo produto",
                "parameters": [
                    {
                        "description": "Dados do produto",
                        "name": "produto",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Produto"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/Produto"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/produtos/{id}": {
            "get": {
                "description": "Retorna um produto pelo ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "produtos"
                ],
                "summary": "Obter um produto específico",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do produto",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Produto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Atualiza todos os campos de um produto existente",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "produtos"
                ],
                "summary": "Atualizar um produto completo",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do produto",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Dados atualizados do produto",
                        "name": "produto",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Produto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Produto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove um produto do sistema",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "produtos"
                ],
                "summary": "Deletar um produto",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do produto",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "patch": {
                "description": "Atualiza apenas os campos fornecidos de um produto",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "produtos"
                ],
                "summary": "Atualizar parcialmente um produto",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID do produto",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Campos a atualizar",
                        "name": "produto",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Produto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Produto"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Retorna o status da API",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Verificar saúde da API",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Produto": {
            "description": "Representa um produto da API",
            "type": "object",
            "properties": {
                "descricao": {
                    "type": "string",
                    "example": "Notebook de alta performance"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "nome": {
                    "type": "string",
                    "example": "Notebook"
                },
                "preco": {
                    "type": "number",
                    "example": 3500
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
var SwaggerInfo swaggerInfo

func init() {
	swag.Register(swag.Name, &swag.Spec{
		Version:          "1.0",
		Host:             "localhost:8080",
		BasePath:         "/api",
		Schemes:          []string{"http"},
		Title:            "API Go com Arquitetura",
		Description:      "Uma API REST completa em Go com suporte aos verbos HTTP",
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
		LeftDelim:        "{{",
		RightDelim:       "}}",
	})
}
