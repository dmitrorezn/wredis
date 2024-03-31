package wredis

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"sync"

	"google.golang.org/protobuf/proto"
)

type encoder interface {
	Encode(a any) ([]byte, error)
}

type decoder interface {
	Decode(b []byte, a any) error
}

type EncoderDecoder interface {
	encoder
	decoder
}

const (
	JSON  = "json"
	PROTO = "proto"
	GOB   = "gob"
)

var ErrNotFoundEncoder = errors.New("not found encoder")

type safeEncoders map[string]EncoderDecoder

func (se safeEncoders) Get(name string) (EncoderDecoder, error) {
	encodersMu.Lock()
	defer encodersMu.Unlock()

	enc, ok := encoders[name]
	if !ok {
		return nil, ErrNotFoundEncoder
	}

	return enc, nil
}

var (
	encodersMu sync.Mutex
	_encoders  = make(map[EncoderDecoder]struct{})
	encoders   = safeEncoders(map[string]EncoderDecoder{})

	_ = Register(JSON, JSONEncoder{})
	_ = Register(GOB, GOBEncoder{})
	_ = Register(PROTO, ProtoEncoder{})
)

func Register(name string, enc EncoderDecoder) error {
	encodersMu.Lock()
	defer encodersMu.Unlock()

	if _, ok := encoders[name]; ok {
		return errors.New("already registered name")
	}
	if _, ok := _encoders[enc]; ok {
		return errors.New("already registered encoder")
	}

	encoders[name] = enc
	_encoders[enc] = struct{}{}

	return nil
}

type JSONEncoder struct{}

var _ EncoderDecoder = JSONEncoder{}

func (j JSONEncoder) Encode(a any) ([]byte, error) {
	return json.Marshal(a)
}
func (j JSONEncoder) Decode(b []byte, a any) error {
	return json.Unmarshal(b, a)
}

type GOBEncoder struct{}

var _ EncoderDecoder = GOBEncoder{}

func (j GOBEncoder) Encode(a any) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(a); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func (j GOBEncoder) Decode(b []byte, a any) error {
	return gob.
		NewDecoder(bytes.NewReader(b)).
		Decode(a)
}

type ProtoEncoder struct{}

var _ EncoderDecoder = ProtoEncoder{}

func (j ProtoEncoder) Encode(a any) ([]byte, error) {
	if msg, ok := a.(proto.Message); ok {
		return proto.Marshal(msg)
	}

	return nil, ErrValueNotImplementsProto
}

var ErrValueNotImplementsProto = errors.New("value does not implements proto.Message")

func (j ProtoEncoder) Decode(b []byte, a any) error {
	if msg, ok := a.(proto.Message); ok {
		return proto.Unmarshal(b, msg)
	}

	return ErrValueNotImplementsProto
}
