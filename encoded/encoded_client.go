package wredis

import (
	"context"
	"encoding"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

type EncodedUniversalClient interface {
	StringCmdable
	HashCmdable
}

type StringCmdable interface {
	Append(ctx context.Context, key string, value string) *IntCmd
	Decr(ctx context.Context, key string) *IntCmd
	DecrBy(ctx context.Context, key string, decrement int64) *IntCmd
	Get(ctx context.Context, key string) *StringCmd
	GetRange(ctx context.Context, key string, start int64, end int64) *StringCmd
	GetSet(ctx context.Context, key string, value interface{}) *StringCmd
	GetEx(ctx context.Context, key string, expiration time.Duration) *StringCmd
	GetDel(ctx context.Context, key string) *StringCmd
	Incr(ctx context.Context, key string) *IntCmd
	IncrBy(ctx context.Context, key string, value int64) *IntCmd
	IncrByFloat(ctx context.Context, key string, value float64) *FloatCmd
	LCS(ctx context.Context, q *LCSQuery) *LCSCmd
	MGet(ctx context.Context, keys ...string) *SliceCmd
	MSet(ctx context.Context, values ...interface{}) *StatusCmd
	MSetNX(ctx context.Context, values ...interface{}) *BoolCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	SetArgs(ctx context.Context, key string, value interface{}, a SetArgs) *StatusCmd
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd
	SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd
	SetRange(ctx context.Context, key string, offset int64, value string) *IntCmd
	StrLen(ctx context.Context, key string) *IntCmd
}
type HashCmdable interface {
	HDel(ctx context.Context, key string, fields ...string) *IntCmd
	HExists(ctx context.Context, key, field string) *BoolCmd
	HGet(ctx context.Context, key, field string) *StringCmd
	HGetAll(ctx context.Context, key string) *MapStringStringCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd
	HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd
	HKeys(ctx context.Context, key string) *StringSliceCmd
	HLen(ctx context.Context, key string) *IntCmd
	HMGet(ctx context.Context, key string, fields ...string) *SliceCmd
	HSet(ctx context.Context, key string, values ...interface{}) *IntCmd
	HMSet(ctx context.Context, key string, values ...interface{}) *BoolCmd
	HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmd
	HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd
	HVals(ctx context.Context, key string) *StringSliceCmd
	HRandField(ctx context.Context, key string, count int) *StringSliceCmd
	HRandFieldWithValues(ctx context.Context, key string, count int) *KeyValueSliceCmd
}

var _ EncodedUniversalClient = (*EncodedClient)(nil)

type EncodedClient struct {
	redis.UniversalClient
	encoder encoder
	decoder decoder
}

type (
	BoolCmd            redis.BoolCmd
	StatusCmd          redis.StatusCmd
	MapStringStringCmd redis.MapStringStringCmd
	IntCmd             redis.IntCmd
	FloatCmd           redis.FloatCmd
	StringSliceCmd     redis.StringSliceCmd
	SliceCmd           redis.SliceCmd
	ScanCmd            redis.ScanCmd
	KeyValueSliceCmd   redis.KeyValueSliceCmd
	LCSQuery           redis.LCSQuery
	SetArgs            redis.SetArgs
	LCSCmd             redis.LCSCmd
)

func boolCmdWithErr(ctx context.Context, err error) *BoolCmd {
	var cmd = (*BoolCmd)(redis.NewBoolCmd(ctx))
	cmd.SetErr(err)

	return cmd
}

func intCmdWithErr(ctx context.Context, err error) *IntCmd {
	var cmd = (*IntCmd)(redis.NewIntCmd(ctx))
	cmd.SetErr(err)

	return cmd
}

func statusCmdWithErr(ctx context.Context, err error) *StatusCmd {
	var cmd = (*StatusCmd)(redis.NewStatusCmd(ctx))
	cmd.SetErr(err)

	return cmd
}

func NewEncodedClient(encoderName string, client redis.UniversalClient) (*EncodedClient, error) {
	enc, err := encoders.Get(encoderName)
	if err != nil {
		return nil, err
	}
	ec := &EncodedClient{
		encoder:         enc,
		decoder:         enc,
		UniversalClient: client,
	}

	return ec, nil
}

type StringCmd struct {
	decoder
	*redis.StringCmd
}

func NewStringCmd(d decoder, cmd *redis.StringCmd) *StringCmd {
	return &StringCmd{
		decoder:   d,
		StringCmd: cmd,
	}
}

var _ encoding.BinaryMarshaler = BinaryEncoder{}
var _ encoding.BinaryUnmarshaler = BinaryDecoder{}

type BinaryEncoder struct {
	encoder
	Val any
}

func (b BinaryEncoder) MarshalBinary() (data []byte, err error) {
	return b.encoder.Encode(b.Val)
}

type BinaryDecoder struct {
	decoder
	Val any
}

func (b BinaryDecoder) UnmarshalBinary(data []byte) error {
	return b.decoder.Decode(data, b.Val)
}

func isSupportedByRedisLicType(v interface{}) bool {
	switch v.(type) {
	case *string, *[]byte, *int, *int8, *int16, *int32, *int64, *uint, *uint16, *uint32, *uint64, *float32, *float64, *time.Time, *time.Duration, *net.IP:
		return true
	}

	return false
}

func (sc *StringCmd) Scan(v interface{}) error {
	if v != nil && !isSupportedByRedisLicType(v) {
		v = BinaryDecoder{
			decoder: sc.decoder,
			Val:     v,
		}
	}

	return sc.StringCmd.Scan(v)
}

func (c *EncodedClient) Get(ctx context.Context, key string) *StringCmd {
	return &StringCmd{decoder: c.decoder, StringCmd: c.UniversalClient.Get(ctx, key)}
}

func (c *EncodedClient) Decode(b []byte, value interface{}) error {
	return c.decoder.Decode(b, value)
}

func (c *EncodedClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) *StatusCmd {
	var err error
	if value, err = c.encodeValue(value); err != nil {
		return statusCmdWithErr(ctx, err)
	}

	return (*StatusCmd)(c.UniversalClient.Set(ctx, key, value, ttl))
}

func (c *EncodedClient) GetDel(ctx context.Context, key string) *StringCmd {
	return &StringCmd{
		decoder:   c.decoder,
		StringCmd: c.UniversalClient.GetDel(ctx, key),
	}
}

func (c *EncodedClient) HGet(ctx context.Context, key, field string) *StringCmd {
	return &StringCmd{
		decoder:   c.decoder,
		StringCmd: c.UniversalClient.HGet(ctx, key, field),
	}
}

func (c *EncodedClient) encodeValues(values ...interface{}) error {
	var err error
	for i := range values {
		if (i+1)%2 == 0 {
			values[i], err = c.encodeValue(values[i])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (c *EncodedClient) encodeValue(value interface{}) (interface{}, error) {
	if !isSupportedByRedisLicType(value) {
		var err error
		value, err = c.encoder.Encode(value)
		if err != nil {
			return value, err
		}
	}
	return value, nil
}

func (c *EncodedClient) HSet(ctx context.Context, key string, values ...interface{}) *IntCmd {
	if err := c.encodeValues(values); err != nil {
		var cmd = (*IntCmd)(redis.NewIntCmd(ctx))
		cmd.SetErr(err)

		return cmd
	}

	return (*IntCmd)(c.UniversalClient.HSet(ctx, key, values...))
}

func (c *EncodedClient) HMSet(ctx context.Context, key string, values ...interface{}) *BoolCmd {
	if err := c.encodeValues(values); err != nil {
		return boolCmdWithErr(ctx, err)
	}

	return (*BoolCmd)(c.UniversalClient.HMSet(ctx, key, values...))
}

func (c *EncodedClient) HSetNX(ctx context.Context, key, field string, value interface{}) *BoolCmd {
	var err error
	if value, err = c.encodeValue(value); err != nil {
		return boolCmdWithErr(ctx, err)
	}

	return (*BoolCmd)(c.UniversalClient.HSetNX(ctx, key, field, value))
}

func (c *EncodedClient) Append(ctx context.Context, key string, value string) *IntCmd {
	return (*IntCmd)(c.UniversalClient.Append(ctx, key, value))
}

func (c *EncodedClient) Decr(ctx context.Context, key string) *IntCmd {
	return (*IntCmd)(c.UniversalClient.Decr(ctx, key))
}

func (c *EncodedClient) DecrBy(ctx context.Context, key string, decrement int64) *IntCmd {
	return (*IntCmd)(c.UniversalClient.DecrBy(ctx, key, decrement))
}

func (c *EncodedClient) GetRange(ctx context.Context, key string, start int64, end int64) *StringCmd {
	return &StringCmd{
		decoder:   c.decoder,
		StringCmd: c.UniversalClient.GetRange(ctx, key, start, end),
	}
}

func (c *EncodedClient) GetSet(ctx context.Context, key string, value interface{}) *StringCmd {
	var err error
	if value, err = c.encodeValue(value); err != nil {
		cmd := NewStringCmd(c.decoder, redis.NewStringCmd(ctx))
		cmd.SetErr(err)

		return cmd
	}

	return &StringCmd{
		decoder:   c.decoder,
		StringCmd: c.UniversalClient.GetSet(ctx, key, value),
	}
}

func (c *EncodedClient) MGet(ctx context.Context, keys ...string) *SliceCmd {
	return (*SliceCmd)(c.UniversalClient.MGet(ctx, keys...))
}

func (c *EncodedClient) MSet(ctx context.Context, values ...interface{}) *StatusCmd {
	if err := c.encodeValues(values); err != nil {
		return statusCmdWithErr(ctx, err)
	}

	return (*StatusCmd)(c.UniversalClient.MSet(ctx, values))
}

func (c *EncodedClient) MSetNX(ctx context.Context, values ...interface{}) *BoolCmd {
	if err := c.encodeValues(values); err != nil {
		return boolCmdWithErr(ctx, err)
	}

	return (*BoolCmd)(c.UniversalClient.MSetNX(ctx, values))
}

func (c *EncodedClient) SetArgs(ctx context.Context, key string, value interface{}, a SetArgs) *StatusCmd {
	var err error
	if value, err = c.encodeValue(value); err != nil {
		return statusCmdWithErr(ctx, err)
	}

	return (*StatusCmd)(c.UniversalClient.SetArgs(ctx, key, value, redis.SetArgs(a)))
}

func (c *EncodedClient) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd {
	var err error
	if value, err = c.encodeValue(value); err != nil {
		return statusCmdWithErr(ctx, err)
	}

	return (*StatusCmd)(c.UniversalClient.SetEx(ctx, key, value, expiration))
}

func (c *EncodedClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd {
	var err error
	if value, err = c.encodeValue(value); err != nil {
		return boolCmdWithErr(ctx, err)
	}

	return (*BoolCmd)(c.UniversalClient.SetNX(ctx, key, value, expiration))
}

func (c *EncodedClient) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *BoolCmd {
	var err error
	if value, err = c.encodeValue(value); err != nil {
		return boolCmdWithErr(ctx, err)
	}

	return (*BoolCmd)(c.UniversalClient.SetXX(ctx, key, value, expiration))
}

func (c *EncodedClient) SetRange(ctx context.Context, key string, offset int64, value string) *IntCmd {
	return (*IntCmd)(c.UniversalClient.SetRange(ctx, key, offset, value))
}

func (c *EncodedClient) GetEx(ctx context.Context, key string, expiration time.Duration) *StringCmd {
	return &StringCmd{
		decoder:   c.decoder,
		StringCmd: c.UniversalClient.GetEx(ctx, key, expiration),
	}
}

func (c *EncodedClient) Incr(ctx context.Context, key string) *IntCmd {
	return (*IntCmd)(c.UniversalClient.Incr(ctx, key))
}

func (c *EncodedClient) IncrBy(ctx context.Context, key string, value int64) *IntCmd {
	return (*IntCmd)(c.UniversalClient.IncrBy(ctx, key, value))
}

func (c *EncodedClient) IncrByFloat(ctx context.Context, key string, value float64) *FloatCmd {
	return (*FloatCmd)(c.UniversalClient.IncrByFloat(ctx, key, value))
}

func (c *EncodedClient) LCS(ctx context.Context, q *LCSQuery) *LCSCmd {
	return (*LCSCmd)(c.UniversalClient.LCS(ctx, (*redis.LCSQuery)(q)))
}

func (c *EncodedClient) StrLen(ctx context.Context, key string) *IntCmd {
	return (*IntCmd)(c.UniversalClient.StrLen(ctx, key))
}

func (c *EncodedClient) HDel(ctx context.Context, key string, fields ...string) *IntCmd {
	return (*IntCmd)(c.UniversalClient.HDel(ctx, key, fields...))
}

func (c *EncodedClient) HExists(ctx context.Context, key, field string) *BoolCmd {
	return (*BoolCmd)(c.UniversalClient.HExists(ctx, key, field))
}

func (c *EncodedClient) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *ScanCmd {
	return (*ScanCmd)(c.UniversalClient.HScan(ctx, key, cursor, match, count))
}

func (c *EncodedClient) HVals(ctx context.Context, key string) *StringSliceCmd {
	return (*StringSliceCmd)(c.UniversalClient.HVals(ctx, key))
}

func (c *EncodedClient) HRandField(ctx context.Context, key string, count int) *StringSliceCmd {
	return (*StringSliceCmd)(c.UniversalClient.HRandField(ctx, key, count))
}

func (c *EncodedClient) HRandFieldWithValues(ctx context.Context, key string, count int) *KeyValueSliceCmd {
	return (*KeyValueSliceCmd)(c.UniversalClient.HRandFieldWithValues(ctx, key, count))
}

func (c *EncodedClient) HGetAll(ctx context.Context, key string) *MapStringStringCmd {
	return (*MapStringStringCmd)(c.UniversalClient.HGetAll(ctx, key))
}

func (c *EncodedClient) HIncrBy(ctx context.Context, key, field string, incr int64) *IntCmd {
	return (*IntCmd)(c.UniversalClient.HIncrBy(ctx, key, field, incr))
}

func (c *EncodedClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) *FloatCmd {
	return (*FloatCmd)(c.UniversalClient.HIncrByFloat(ctx, key, field, incr))
}

func (c *EncodedClient) HKeys(ctx context.Context, key string) *StringSliceCmd {
	return (*StringSliceCmd)(c.UniversalClient.HKeys(ctx, key))
}

func (c *EncodedClient) HLen(ctx context.Context, key string) *IntCmd {
	return (*IntCmd)(c.UniversalClient.HLen(ctx, key))
}

func (c *EncodedClient) HMGet(ctx context.Context, key string, fields ...string) *SliceCmd {
	return (*SliceCmd)(c.UniversalClient.HMGet(ctx, key, fields...))
}
