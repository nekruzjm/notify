// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
            "email": "jamshedzodnekruz@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/notifications-external/v1/push": {
            "post": {
                "security": [
                    {
                        "SignatureAuth": []
                    }
                ],
                "description": "All fields except ` + "`" + `personExternalRef` + "`" + ` (crm_client_id) are required.\n- If you want to send push with ` + "`" + `personExternalRef` + "`" + `, do not provide ` + "`" + `phone` + "`" + `.\n- If ` + "`" + `showInFeed` + "`" + ` is true, the push will be shown in the feed; otherwise, it will be hidden.\n- If the users status is inactive or their push setting is disabled, the push will be saved in the feed but not sent to the device.\nIn that case, the payload will be ` + "`" + `inactive_user#fake_message_id` + "`" + ` or ` + "`" + `disabled_push#fake_message_id` + "`" + `.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "External"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Provide user ID created on the server side",
                        "name": "X-UserId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Provide unique request ID to build hash and track the request",
                        "name": "X-RequestId",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Provide current date in RFC1123 format (e.g., Mon, 02 Jan 2006 15:04:05 MST) to build hash",
                        "name": "X-Date",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Provide user action (push, sms) to send push",
                        "name": "X-UserAction",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Provide hash sum built with HMAC-SHA256 from the ` + "`" + `X-Date:X-RequestId` + "`" + ` using the secret key created on the server side",
                        "name": "X-RequestDigest",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Request payload",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/push.externalRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            }
        },
        "/notifications-internal/v1/events": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Get list of events",
                "parameters": [
                    {
                        "type": "string",
                        "description": "apply filter with id",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "apply filter with status",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "apply filter with topic",
                        "name": "topic",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "apply filter with limit, 10 settled by default",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "apply filter with offset, 0 settled by default",
                        "name": "offset",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/event.eventModel"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "List not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Create event",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/event.request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "$ref": "#/definitions/event.eventModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            }
        },
        "/notifications-internal/v1/events/{id}": {
            "put": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Update event",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/event.request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "$ref": "#/definitions/event.eventModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Delete event",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            }
        },
        "/notifications-internal/v1/events/{id}/image/{language}": {
            "post": {
                "description": "Provider event id - ` + "`" + `{id}` + "`" + ` and multi-lang image key - ` + "`" + `{language}` + "`" + ` as api route var to upload image for event",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Upload image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Load images with png, jpg, jpeg extensions",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "$ref": "#/definitions/event.eventModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Provider event id - ` + "`" + `{id}` + "`" + ` and multi-lang image key - ` + "`" + `{language}` + "`" + ` as api route var to remove specific image of event",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Remove image",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "$ref": "#/definitions/event.eventModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            }
        },
        "/notifications-internal/v1/events/{id}/load-all-users": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Subscribe all users to the event",
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "$ref": "#/definitions/event.eventModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            }
        },
        "/notifications-internal/v1/events/{id}/load-users": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Subscribe list of users to event",
                "parameters": [
                    {
                        "type": "file",
                        "description": "CSV file with ` + "`" + `userID` + "`" + ` header and data. Do not provide other headers",
                        "name": "users",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "$ref": "#/definitions/event.eventModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            }
        },
        "/notifications-internal/v1/events/{id}/run": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Run event manually",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/resp.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "payload": {
                                            "$ref": "#/definitions/event.eventModel"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "401": {
                        "description": "Invalid authorization data",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "404": {
                        "description": "Not found",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Error",
                        "schema": {
                            "$ref": "#/definitions/resp.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "event.eventModel": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/language.Language"
                },
                "category": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "extraData": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "$ref": "#/definitions/language.Language"
                },
                "link": {
                    "type": "string"
                },
                "scheduledAt": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "$ref": "#/definitions/language.Language"
                },
                "topic": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "event.request": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/language.Language"
                },
                "category": {
                    "type": "string"
                },
                "extraData": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "link": {
                    "type": "string"
                },
                "scheduledAt": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "$ref": "#/definitions/language.Language"
                },
                "topic": {
                    "type": "string"
                }
            }
        },
        "language.Language": {
            "type": "object",
            "properties": {
                "en": {
                    "type": "string"
                },
                "ru": {
                    "type": "string"
                },
                "tg": {
                    "type": "string"
                },
                "uz": {
                    "type": "string"
                }
            }
        },
        "push.externalRequest": {
            "type": "object",
            "required": [
                "body",
                "phone",
                "showInFeed",
                "title",
                "type"
            ],
            "properties": {
                "body": {
                    "$ref": "#/definitions/language.Language"
                },
                "personExternalRef": {
                    "type": "string",
                    "example": "123456"
                },
                "phone": {
                    "type": "string",
                    "example": "+992111111111"
                },
                "showInFeed": {
                    "type": "boolean"
                },
                "title": {
                    "$ref": "#/definitions/language.Language"
                },
                "type": {
                    "type": "string",
                    "example": "otp, push"
                }
            }
        },
        "resp.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "payload": {}
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "api-notifications.dev.my.cloud",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Notifications API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
