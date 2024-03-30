package luchen

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/fengjx/go-halo/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()
var validate *validator.Validate

func init() {
	decoder.IgnoreUnknownKeys(true)
	decoder.SetAliasTag("json")
	decoder.RegisterConverter([]string{}, func(s string) reflect.Value {
		return reflect.ValueOf(strings.Split(s, ","))
	})
	validate = validator.New()
	validate.SetTagName("binding")
}

// ShouldBind 从参数url参数和form表单解析参数
func ShouldBind(r *http.Request, obj any) error {
	values := r.URL.Query()
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		for key, val := range r.Form {
			values[key] = val
		}
	}
	err := decoder.Decode(obj, values)
	if err != nil {
		return err
	}
	return validate.Struct(obj)
}

// ShouldBindJSON 从body解析json
func ShouldBindJSON(r *http.Request, obj any) error {
	err := json.NewDecoder(r.Body).Decode(obj)
	if err != nil {
		return err
	}
	return validate.Struct(obj)
}
