package marshal

// http header定义
const (
	ContentTypeJSON     = "application/json"
	ContentTypeProtobuf = "application/protobuf"
)

var (
	// content-type 到 Marshaller 的映射
	// 这里是 http 请求的 content-type 到 Marshaller 的映射
	// http 请求支持原生的 protobuf 协议和 json 协议
	// 如果是 json 协议，会把 json 转为 protobuf 协议，反序列化的时候再把 protobuf 转为 json
	marshallerMap = map[string]Marshaller{
		ContentTypeJSON:     NewJSONProtoMarshaller(),
		ContentTypeProtobuf: NewProtoMarshaller(),
	}

	// 默认的 Marshaller
	defaultMarshaller = NewDefaultMarshaller()
)

// GetMarshallerByContentType 获取根据 http content-type 获取对应 Marshaller
func GetMarshallerByContentType(contentType string) Marshaller {
	if m, ok := marshallerMap[contentType]; ok {
		return m
	}
	return defaultMarshaller
}

// UseDefaultMarshaller 设置默认的 Marshaller
func UseDefaultMarshaller(m Marshaller) {
	defaultMarshaller = m
}

// RegisterMarshaller 注册 Marshaller
// 不要在并发情况下调用
func RegisterMarshaller(contentType string, m Marshaller) {
	marshallerMap[contentType] = m
}
