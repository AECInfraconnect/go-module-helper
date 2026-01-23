package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joncalhoun/qson"
)

const (
	// MiddleWareJWT is a constant identifier for JWT middleware type.
	MiddleWareJWT = "jwt"
)

// GoMiddlewareInf defines the interface for request parsing middlewares.
type GoMiddlewareInf interface {
	InitContextIfNotExists() gin.HandlerFunc
	InputForm() gin.HandlerFunc
}

// GoMiddleware handles request parsing and context initialization.
type GoMiddleware struct {
	ctx       context.Context
	jwtSecret string
}

// InputForm parses request body from various content types.
//
// Supports application/json, multipart/form-data, and application/x-www-form-urlencoded.
// Parsed data is stored in context with key "params".
// Uploaded files are stored with keys "files" or "part_file".
//
// Example:
//
//	m := middleware.InitMiddleware("jwt-secret")
//	r.Use(m.InputForm())
func (m *GoMiddleware) InputForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := Form(c); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Next()
	}
}

// InitContextIfNotExists ensures the request context is initialized.
//
// Creates a background context if none exists.
// Should be applied early in the middleware chain.
//
// Example:
//
//	m := middleware.InitMiddleware("jwt-secret")
//	r.Use(m.InitContextIfNotExists())
func (m *GoMiddleware) InitContextIfNotExists() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if ctx == nil {
			bgCtx := context.Background()
			c.Request = c.Request.WithContext(bgCtx)
		}
		c.Next()
	}
}

// InitMiddleware initializes the middleware with a JWT secret key.
//
// Returns a GoMiddlewareInf interface providing InputForm and InitContextIfNotExists middlewares.
//
// Example:
//
//	m := middleware.InitMiddleware("your-jwt-secret")
//	r.Use(m.InitContextIfNotExists())
//	r.Use(m.InputForm())
func InitMiddleware(key string) GoMiddlewareInf {
	return &GoMiddleware{
		ctx:       context.TODO(),
		jwtSecret: key,
	}
}

// Form parses the request body based on Content-Type header.
//
// Supports JSON, multipart form data, and URL-encoded forms.
// Stores parsed parameters in context with key "params".
func Form(c *gin.Context) error {
	var data = map[string]any{}
	reqMethod := c.Request.Method
	Header := c.Request.Header

	if reqMethod == http.MethodPost || reqMethod == http.MethodPut || reqMethod == http.MethodDelete {
		contentType := Header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			form, err := c.MultipartForm()
			if err != nil {
				return fmt.Errorf("%s or has not any parameter", http.ErrMissingBoundary.Error())
			}
			bu, _ := qson.ToJSON(url.Values(form.Value).Encode())
			json.Unmarshal(bu, &data)

			data, err = parseOnKeyData(data)
			if err != nil {
				return err
			}

			/* รูปสำหรับ ใช้งานทั่วไป */
			if val, ok := form.File["files"]; ok {
				c.Set("files", val)
			}
			if val, ok := form.File["part_file"]; ok {
				c.Set("part_file", val)
			}
		} else if strings.Contains(contentType, "application/json") {
			var err error
			if err := c.ShouldBindJSON(&data); err != nil && err != io.EOF {
				return err
			}
			data, err = parseOnKeyData(data)
			if err != nil {
				return err
			}

		} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			postForm := c.Request.PostForm
			var err error
			if reqMethod == http.MethodDelete {
				buf := bytes.Buffer{}
				io.Copy(&buf, c.Request.Body)
				postForm, _ = url.ParseQuery(buf.String())
			}
			if len(postForm) > 0 {
				bu, _ := qson.ToJSON(postForm.Encode())
				json.Unmarshal(bu, &data)
			}
			data, err = parseOnKeyData(data)
			if err != nil {
				return err
			}
		}
	}

	if len(data) > 0 {
		c.Set("params", data)
	}
	return nil
}

func parseOnKeyData(data map[string]any) (map[string]any, error) {
	if len(data) == 1 {
		/*
			support on data from json format
		*/
		if v, ok := data["data"]; ok {
			valueType := reflect.ValueOf(v).Kind()
			switch valueType {
			case reflect.Map:
				data = v.(map[string]any)
			case reflect.String:
				data = map[string]any{}
				if err := json.Unmarshal([]byte(v.(string)), &data); err != nil {
					return data, err
				}
			}
		}
	}

	return data, nil
}
