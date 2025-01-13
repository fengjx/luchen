package marshal

import (
	"errors"
	"io"
)

var ErrUnsupportedCodec = errors.New("unsupported codec type")

// Marshaller 定义了序列化和反序列化数据的接口
// 实现此接口的类型可以处理特定格式的数据转换(如 JSON、Protobuf 等)
type Marshaller interface {
	// Marshal 将给定的值序列化为字节切片
	Marshal(v any) ([]byte, error)

	// Unmarshal 将字节切片反序列化为给定的值
	Unmarshal(data []byte, v any) error

	// NewDecoder 创建一个从 io.Reader 读取并解码数据的解码器
	NewDecoder(r io.Reader) Decoder

	// NewEncoder 创建一个将数据编码并写入 io.Writer 的编码器
	NewEncoder(w io.Writer) Encoder

	// ContentType 请求响应的 content-type
	ContentType() string
}

// Decoder 定义了数据解码接口
type Decoder interface {
	// Decode 将数据解码到指定的值中
	Decode(v any) error
}

// Encoder 定义了数据编码接口
type Encoder interface {
	// Encode 将指定的值编码为数据格式
	Encode(v any) error
}

// DecoderFunc 是实现 Decoder 接口的函数类型
type DecoderFunc func(v any) error

// Decode 实现 Decoder 接口
func (f DecoderFunc) Decode(v any) error { return f(v) }

// EncoderFunc 是实现 Encoder 接口的函数类型
type EncoderFunc func(v any) error

// Encode 实现 Encoder 接口
func (f EncoderFunc) Encode(v any) error { return f(v) }

// DefaultMarshaller 提供了一个简单的 Marshaler 实现，它总是返回不支持的错误
type DefaultMarshaller struct{}

// Marshal 返回不支持的编解码类型错误
func (d *DefaultMarshaller) Marshal(v any) ([]byte, error) {
	return nil, ErrUnsupportedCodec
}

// Unmarshal 返回不支持的编解码类型错误
func (d *DefaultMarshaller) Unmarshal(data []byte, v any) error {
	return ErrUnsupportedCodec
}

// NewDecoder 创建一个总是返回错误的解码器
func (d *DefaultMarshaller) NewDecoder(r io.Reader) Decoder {
	return DecoderFunc(func(v any) error {
		return ErrUnsupportedCodec
	})
}

// NewEncoder 创建一个总是返回错误的编码器
func (d *DefaultMarshaller) NewEncoder(w io.Writer) Encoder {
	return EncoderFunc(func(v any) error {
		return ErrUnsupportedCodec
	})
}

// ContentType 返回不支持的编解码类型错误
func (d *DefaultMarshaller) ContentType() string {
	return "text/plain"
}

// NewDefaultMarshaller 创建一个新的默认 Marshaller
func NewDefaultMarshaller() Marshaller {
	return &DefaultMarshaller{}
}
