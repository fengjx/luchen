package marshal

import (
	"errors"
	"io"

	"google.golang.org/protobuf/proto"
)

// ErrInvalidProtoMessage 当传入的值不是有效的 protobuf 消息时返回此错误
var ErrInvalidProtoMessage = errors.New("value is not a proto.Message")

// ProtoMarshaller 实现了 Marshaler 接口，提供 Protocol Buffers 格式的序列化和反序列化功能
type ProtoMarshaller struct{}

// Marshal 将给定的 protobuf 消息序列化为字节切片
func (p *ProtoMarshaller) Marshal(v any) ([]byte, error) {
	msg, ok := v.(proto.Message)
	if !ok {
		return nil, ErrInvalidProtoMessage
	}
	return p.Marshal(msg)
}

// Unmarshal 将字节切片反序列化为给定的 protobuf 消息
func (p *ProtoMarshaller) Unmarshal(data []byte, v any) error {
	msg, ok := v.(proto.Message)
	if !ok {
		return ErrInvalidProtoMessage
	}
	return p.Unmarshal(data, msg)
}

// NewDecoder 创建一个从 io.Reader 读取并解码 protobuf 数据的解码器
func (p *ProtoMarshaller) NewDecoder(r io.Reader) Decoder {
	return DecoderFunc(func(v any) error {
		msg, ok := v.(proto.Message)
		if !ok {
			return ErrInvalidProtoMessage
		}

		data, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		return p.Unmarshal(data, msg)
	})
}

// NewEncoder 创建一个将数据编码为 protobuf 并写入 io.Writer 的编码器
func (p *ProtoMarshaller) NewEncoder(w io.Writer) Encoder {
	return EncoderFunc(func(v any) error {
		msg, ok := v.(proto.Message)
		if !ok {
			return ErrInvalidProtoMessage
		}

		data, err := p.Marshal(msg)
		if err != nil {
			return err
		}

		_, err = w.Write(data)
		return err
	})
}

// ContentType 返回 Protocol Buffers 的 MIME 类型
func (p *ProtoMarshaller) ContentType() string {
	return ContentTypeProtobuf
}

// NewProtoMarshaller 创建一个新的 Protocol Buffers Marshaler
func NewProtoMarshaller() Marshaller {
	return &ProtoMarshaller{}
}
