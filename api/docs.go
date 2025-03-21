// Package api Code generated by swaggo/swag. DO NOT EDIT
package api

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "SMSGate Support",
            "email": "support@sms-gate.app"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/csr": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CSR"
                ],
                "summary": "Submit CSR",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ca.PostCSRRequest"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted",
                        "schema": {
                            "$ref": "#/definitions/ca.PostCSRResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.JSONErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.JSONErrorResponse"
                        }
                    }
                }
            }
        },
        "/csr/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CSR"
                ],
                "summary": "Get CSR Status",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Request ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ca.GetCSRStatusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/http.JSONErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/http.JSONErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/http.JSONErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "ca.CSRStatus": {
            "type": "string",
            "enum": [
                "pending",
                "approved",
                "denied"
            ],
            "x-enum-varnames": [
                "CSRStatusPending",
                "CSRStatusApproved",
                "CSRStatusDenied"
            ]
        },
        "ca.CSRType": {
            "type": "string",
            "enum": [
                "webhook",
                "private_server"
            ],
            "x-enum-varnames": [
                "CSRTypeWebhook",
                "CSRTypePrivateServer"
            ]
        },
        "ca.GetCSRStatusResponse": {
            "type": "object",
            "properties": {
                "certificate": {
                    "description": "Certificate is the certificate issued by the CA. This field is only present\nif the status is ` + "`" + `approved` + "`" + `.",
                    "type": "string"
                },
                "message": {
                    "description": "Message is a human-readable description of the status.",
                    "type": "string"
                },
                "request_id": {
                    "description": "RequestID is the ID of the request. Can be used to request status.",
                    "type": "string"
                },
                "status": {
                    "description": "Status is the status of the requested certificate.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/ca.CSRStatus"
                        }
                    ]
                },
                "type": {
                    "description": "Type is the type of the requested certificate.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/ca.CSRType"
                        }
                    ]
                }
            }
        },
        "ca.PostCSRRequest": {
            "type": "object",
            "required": [
                "content"
            ],
            "properties": {
                "content": {
                    "description": "Content contains the CSR content and is required.",
                    "type": "string",
                    "maxLength": 16384
                },
                "metadata": {
                    "description": "Metadata includes additional metadata related to the CSR.",
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "type": {
                    "description": "Type is the type of the CSR. By default, it is set to \"webhook\".",
                    "default": "webhook",
                    "allOf": [
                        {
                            "$ref": "#/definitions/ca.CSRType"
                        }
                    ]
                }
            }
        },
        "ca.PostCSRResponse": {
            "type": "object",
            "properties": {
                "certificate": {
                    "description": "Certificate is the certificate issued by the CA. This field is only present\nif the status is ` + "`" + `approved` + "`" + `.",
                    "type": "string"
                },
                "message": {
                    "description": "Message is a human-readable description of the status.",
                    "type": "string"
                },
                "request_id": {
                    "description": "RequestID is the ID of the request. Can be used to request status.",
                    "type": "string"
                },
                "status": {
                    "description": "Status is the status of the requested certificate.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/ca.CSRStatus"
                        }
                    ]
                },
                "type": {
                    "description": "Type is the type of the requested certificate.",
                    "allOf": [
                        {
                            "$ref": "#/definitions/ca.CSRType"
                        }
                    ]
                }
            }
        },
        "http.JSONErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "properties": {
                        "code": {
                            "description": "Code",
                            "type": "integer"
                        },
                        "message": {
                            "description": "Message",
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "{{VERSION}}",
	Host:             "ca.sms-gate.app",
	BasePath:         "/api/v1",
	Schemes:          []string{"https"},
	Title:            "SMS Gate Certificate Authority API",
	Description:      "Provides methods to manage certificates",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
