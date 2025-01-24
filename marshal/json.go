package marshal

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

// JSONMarshaller 实现了 Marshaller 接口，提供 JSON 格式的序列化和反序列化功能
type JSONMarshaller struct{}

// Marshal 将给定的值序列化为 JSON 字节切片
func (j *JSONMarshaller) Marshal(v any) ([]byte, error) {
	return jsoniter.Marshal(v)
}

// Unmarshal 将 JSON 字节切片反序列化为给定的值
func (j *JSONMarshaller) Unmarshal(data []byte, v any) error {
	return jsoniter.Unmarshal(data, v)
}

// NewDecoder 创建一个从 io.Reader 读取并解码 JSON 数据的解码器
func (j *JSONMarshaller) NewDecoder(r io.Reader) Decoder {
	return jsoniter.NewDecoder(r)
}

// NewEncoder 创建一个将数据编码为 JSON 并写入 io.Writer 的编码器
func (j *JSONMarshaller) NewEncoder(w io.Writer) Encoder {
	return jsoniter.NewEncoder(w)
}

// ContentType 返回 json 的 MIME 类型
func (p *JSONMarshaller) ContentType() string {
	return ContentTypeJSON
}

// NewJSONMarshaller 创建一个新的 JSON Marshaller
func NewJSONMarshaller() Marshaller {
	return &JSONMarshaller{}
}
