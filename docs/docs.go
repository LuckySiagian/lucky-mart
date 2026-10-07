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
        "/health/liveness": {
            "get": {
                "tags": ["System"],
                "summary": "Liveness check endpoint",
                "responses": {"200": {"description": "OK"}}
            }
        },
        "/health/readiness": {
            "get": {
                "tags": ["System"],
                "summary": "Readiness check endpoint",
                "responses": {"200": {"description": "OK"}}
            }
        },
        "/api/v1/auth/register": {
            "post": {
                "tags": ["Authentication"],
                "summary": "Register Customer / B2B Account",
                "responses": {"201": {"description": "User Registered Successfully"}}
            }
        },
        "/api/v1/auth/login": {
            "post": {
                "tags": ["Authentication"],
                "summary": "Login with Phone / Email & Password",
                "responses": {"200": {"description": "Login Successful"}}
            }
        },
        "/api/v1/auth/refresh": {
            "post": {
                "tags": ["Authentication"],
                "summary": "Refresh Access Token",
                "responses": {"200": {"description": "Token Refreshed"}}
            }
        },
        "/api/v1/users/me": {
            "get": {
                "tags": ["User Management"],
                "summary": "Get Profile Info",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Profile Info"}}
            }
        },
        "/api/v1/products": {
            "get": {
                "tags": ["Product Catalog"],
                "summary": "List Products",
                "responses": {"200": {"description": "Products List"}}
            },
            "post": {
                "tags": ["Product Catalog"],
                "summary": "Create Product",
                "security": [{"BearerAuth": []}],
                "responses": {"201": {"description": "Product Created"}}
            }
        },
        "/api/v1/inventory/adjust": {
            "post": {
                "tags": ["Inventory Management"],
                "summary": "Stock Adjustment (Stock In/Out Audit)",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Stock Adjusted"}}
            }
        },
        "/api/v1/cart": {
            "get": {
                "tags": ["Shopping Cart"],
                "summary": "Get Active Shopping Cart",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Cart Details"}}
            },
            "delete": {
                "tags": ["Shopping Cart"],
                "summary": "Clear Cart",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Cart Cleared"}}
            }
        },
        "/api/v1/cart/items": {
            "post": {
                "tags": ["Shopping Cart"],
                "summary": "Add Product Unit to Cart",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Item Added"}}
            }
        },
        "/api/v1/cart/checkout": {
            "post": {
                "tags": ["Checkout Engine"],
                "summary": "Generate Checkout Summary",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Checkout Summary Generated"}}
            }
        },
        "/api/v1/orders": {
            "get": {
                "tags": ["Order Engine"],
                "summary": "List Orders",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Orders List"}}
            },
            "post": {
                "tags": ["Order Engine"],
                "summary": "Create Order (ACID Tx & Row Lock)",
                "security": [{"BearerAuth": []}],
                "responses": {"201": {"description": "Order Created"}}
            }
        },
        "/api/v1/addresses": {
            "get": {
                "tags": ["Address Management"],
                "summary": "List Customer Addresses",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Addresses List"}}
            },
            "post": {
                "tags": ["Address Management"],
                "summary": "Create Address",
                "security": [{"BearerAuth": []}],
                "responses": {"201": {"description": "Address Created"}}
            }
        },
        "/api/v1/shipping/methods": {
            "get": {
                "tags": ["Shipping & Logistics"],
                "summary": "List Shipping Methods",
                "responses": {"200": {"description": "Shipping Methods List"}}
            }
        },
        "/api/v1/payments": {
            "get": {
                "tags": ["Payment Engine"],
                "summary": "List Payments",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Payments List"}}
            },
            "post": {
                "tags": ["Payment Engine"],
                "summary": "Create Payment for Invoice",
                "security": [{"BearerAuth": []}],
                "responses": {"201": {"description": "Payment Created"}}
            }
        },
        "/api/v1/payments/upload-proof": {
            "post": {
                "tags": ["Payment Engine"],
                "summary": "Upload Manual Transfer Proof",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Proof Uploaded"}}
            }
        },
        "/api/v1/credit-ledgers": {
            "get": {
                "tags": ["B2B Credit Ledger"],
                "summary": "List Credit Ledgers",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Credit Ledgers List"}}
            }
        },
        "/api/v1/notifications": {
            "get": {
                "tags": ["Notification Engine"],
                "summary": "List Notifications",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Notifications List"}}
            }
        },
        "/api/v1/events": {
            "get": {
                "tags": ["Event Management"],
                "summary": "List Domain Events Log",
                "security": [{"BearerAuth": []}],
                "responses": {"200": {"description": "Event Logs List"}}
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "LuckyMart API",
	Description:      "Omnichannel Grocery Ordering & Supply Chain Backend API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
