// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "Get a list of songs with optional filtering and pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get all songs",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Page size",
                        "name": "pageSize",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Song name",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Group name",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Song text",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Release date",
                        "name": "releaseDate",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Link",
                        "name": "link",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Songs"
                        }
                    },
					"400": {
          				"description": "Bad request"
					},
        			"500": {
          				"description": "Internal server error"
					}
                }
            },
            "post": {
                "description": "Create a new song",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Create a new song",
                "parameters": [
                    {
                        "description": "Song details",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Songs"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
					"400": {
          				"description": "Bad request"
					},
        			"500": {
          				"description": "Internal server error"
					}
                }
            }
        },
        "/songs/{id}": {
            "get": {
                "description": "Get a song by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get a song by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 0,
                        "description": "Verse number",
                        "name": "vers",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Songs"
                        },
					"400": {
          				"description": "Bad request"
						},
        			"500": {
          				"description": "Internal server error"
						}
                    }
                }
            },
            "delete": {
                "description": "Delete a song by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Delete a song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
					"400": {
          				"description": "Bad request"
					},
        			"500": {
          				"description": "Internal server error"
					}
                }
            },
            "patch": {
                "description": "Update a song by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Update a song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Song details",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Songs"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
					"400": {
          				"description": "Bad request"
					},
        			"500": {
          				"description": "Internal server error"
					}
                }
            }
        }
    },
    "definitions": {
        "models.Songs": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "releaseDate": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                },
                "text": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a sample server Petstore server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}