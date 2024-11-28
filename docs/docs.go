// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/concerts": {
            "get": {
                "description": "Returns list of 10 last concerts in db oredered by concert_id.\nUse last_id query parameter to select concerts before this last_id",
                "tags": [
                    "Concerts"
                ],
                "summary": "Returns list of concerts",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id before which to get last concerts",
                        "name": "last_id",
                        "in": "query"
                    }
                ],
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "Concerts"
                ],
                "summary": "Creates new concert",
                "parameters": [
                    {
                        "description": "Concert to create",
                        "name": "concert",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Concert"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/concerts/{id}": {
            "get": {
                "tags": [
                    "Concerts"
                ],
                "summary": "Returns concert info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Concert id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "Concerts"
                ],
                "summary": "Updates concert info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Concert id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated concert info",
                        "name": "concert",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Concert"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/tickets": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Returns list of 10 last tickets in db oredered by ticket_id.\nUse last_id query parameter to select tickets before this last_id",
                "tags": [
                    "Tickets"
                ],
                "summary": "Returns list of tickets",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id before which to get last tickets",
                        "name": "last_id",
                        "in": "query"
                    }
                ],
                "responses": {}
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Creates new ticket",
                "parameters": [
                    {
                        "description": "Ticket to create",
                        "name": "ticket",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Ticket"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/tickets/own": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Returns list of 10 last tickets that the user has in db oredered by ticket_id.\nUse last_id query parameter to select tickets before this last_id",
                "tags": [
                    "Tickets"
                ],
                "summary": "Returns list of tickets that user has",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "id before which to get last tickets",
                        "name": "last_id",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/tickets/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Returns ticket info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Ticket id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "tags": [
                    "Tickets"
                ],
                "summary": "Updates ticket info",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Ticket id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Updated ticket info",
                        "name": "ticket",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Ticket"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "domain.Concert": {
            "type": "object",
            "properties": {
                "author-id": {
                    "type": "string"
                },
                "create-date": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "location": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "update-date": {
                    "type": "string"
                }
            }
        },
        "domain.Ticket": {
            "type": "object",
            "properties": {
                "concert-id": {
                    "type": "integer"
                },
                "create-date": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "update-date": {
                    "type": "string"
                },
                "user-id": {
                    "type": "string"
                },
                "verification-token": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "description": "To obtain token go to [/auth/login](/auth/login). You also need to add Bearer before pasting it belove. It should look like: Bearer your-access-token",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Concerts api",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
