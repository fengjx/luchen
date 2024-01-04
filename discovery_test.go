package luchen_test

import (
	"testing"

	"github.com/fengjx/go-halo/json"
	"github.com/magiconair/properties/assert"

	"github.com/fengjx/luchen"
)

func TestEtcdV3Selector(t *testing.T) {
	serviceName := "hello"
	registrar := luchen.NewEtcdV3Registrar(
		newHelloHttpServer(serviceName, ":0"),
	)
	registrar.Register()
	defer registrar.Deregister()
	selector := luchen.GetEtcdV3Selector(serviceName)
	serviceInfo, err := selector.Next()
	if err != nil {
		t.Fatal(err)
	}
	jsonStr, _ := json.ToJson(serviceInfo)
	t.Log("serviceInfo", jsonStr)
	assert.Equal(t, serviceName, serviceInfo.Name)
}
