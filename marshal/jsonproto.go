package marshal

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONProtoMarshaller 实现了 Marshaler 接口，提供 JSON 和 Protobuf 之间的转换
// 将 JSON 输入转换为 Protobuf 消息，并将 Protobuf 消息转换为 JSON 输出
type JSONProtoMarshaller struct {
	MarshalOptions   *protojson.MarshalOptions
	UnmarshalOptions *protojson.UnmarshalOptions
}

// Marshal 将 Protobuf 消息序列化为 JSON 字节切片
func (jp *JSONProtoMarshaller) Marshal(v any) ([]byte, error) {
	if pb, ok := v.(proto.Message); ok {
		return jp.MarshalOptions.Marshal(pb)
	}

	return jsoniter.Marshal(v)
}

// Unmarshal 将 JSON 字节切片反序列化为 Protobuf 消息
func (jp *JSONProtoMarshaller) Unmarshal(data []byte, v any) error {
	if pb, ok := v.(proto.Message); ok {
		return jp.UnmarshalOptions.Unmarshal(data, pb)
	}

	return jsoniter.Unmarshal(data, v)
}

// NewDecoder 创建一个从 io.Reader 读取 JSON 并解码为 Protobuf 消息的解码器
func (jp *JSONProtoMarshaller) NewDecoder(r io.Reader) Decoder {
	return DecoderFunc(func(value any) error {
		buffer, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		return jp.Unmarshal(buffer, value)
	})
}

// NewEncoder 创建一个将 Protobuf 消息编码为 JSON 并写入 io.Writer 的编码器
func (jp *JSONProtoMarshaller) NewEncoder(w io.Writer) Encoder {
	return EncoderFunc(func(value any) error {
		buffer, err := jp.Marshal(value)
		if err != nil {
			return err
		}
		_, err = w.Write(buffer)
		if err != nil {
			return err
		}

		return nil
	})
}

// ContentType 返回 json 的 MIME 类型
func (p *JSONProtoMarshaller) ContentType() string {
	return ContentTypeJSON
}

// NewJSONProtoMarshaller 创建一个新的 JSON-Protobuf 转换 Marshaler
func NewJSONProtoMarshaller() Marshaller {
	return &JSONProtoMarshaller{
		MarshalOptions: &protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: &protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
}
