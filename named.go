package wredis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type NamedClient struct {
	Name string
	UniversalClient
}

func NewNamedClient(client UniversalClient, name string) *NamedClient {
	return &NamedClient{
		Name:            name,
		UniversalClient: client,
	}
}

type CMDABLE interface {
	redis.BitMapCmdable
	redis.GenericCmdable
	redis.GeoCmdable
	redis.HashCmdable
	redis.HyperLogLogCmdable
	redis.ListCmdable
	redis.ProbabilisticCmdable
	redis.SetCmdable
	redis.SortedSetCmdable
	redis.StringCmdable
	//	redis.StreamCmdable todo
	redis.TimeseriesCmdable
	redis.JSONCmdable
}

var _ CMDABLE = (*NamedClient)(nil)

func (nc *NamedClient) BFCard(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFCard(ctx, arg1)
}
func (nc *NamedClient) BFInfo(ctx context.Context, arg1 string) *redis.BFInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFInfo(ctx, arg1)
}
func (nc *NamedClient) BFInfoCapacity(ctx context.Context, arg1 string) *redis.BFInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFInfoCapacity(ctx, arg1)
}
func (nc *NamedClient) BFInfoExpansion(ctx context.Context, arg1 string) *redis.BFInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFInfoExpansion(ctx, arg1)
}
func (nc *NamedClient) BFInfoFilters(ctx context.Context, arg1 string) *redis.BFInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFInfoFilters(ctx, arg1)
}
func (nc *NamedClient) BFInfoItems(ctx context.Context, arg1 string) *redis.BFInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFInfoItems(ctx, arg1)
}
func (nc *NamedClient) BFInfoSize(ctx context.Context, arg1 string) *redis.BFInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFInfoSize(ctx, arg1)
}
func (nc *NamedClient) BFInsert(ctx context.Context, arg1 string, arg2 *redis.BFInsertOptions, arg3 ...interface{}) *redis.BoolSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFInsert(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) BFLoadChunk(ctx context.Context, arg1 string, arg2 int64, arg3 interface{}) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFLoadChunk(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) BFReserve(ctx context.Context, arg1 string, arg2 float64, arg3 int64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFReserve(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) BFReserveExpansion(ctx context.Context, arg1 string, arg2 float64, arg3 int64, arg4 int64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFReserveExpansion(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) BFReserveNonScaling(ctx context.Context, arg1 string, arg2 float64, arg3 int64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BFReserveNonScaling(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) BLMPop(ctx context.Context, arg1 time.Duration, arg2 string, arg3 int64, keys ...string) *redis.KeyValuesCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.BLMPop(ctx, arg1, arg2, arg3, keys...)
}
func (nc *NamedClient) BLMove(ctx context.Context, arg1 string, arg2 string, arg3 string, arg4 string, arg5 time.Duration) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BLMove(ctx, arg1, arg2, arg3, arg4, arg5)
}
func (nc *NamedClient) BRPopLPush(ctx context.Context, arg1 string, arg2 string, arg3 time.Duration) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BRPopLPush(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) BZMPop(ctx context.Context, arg1 time.Duration, arg2 string, arg3 int64, keys ...string) *redis.ZSliceWithKeyCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.BZMPop(ctx, arg1, arg2, arg3, keys...)
}
func (nc *NamedClient) BitPos(ctx context.Context, arg1 string, arg2 int64, arg3 ...int64) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BitPos(ctx, arg1, arg2, arg3...)
}
func (nc *NamedClient) BitPosSpan(ctx context.Context, arg1 string, arg2 int8, arg3 int64, arg4 int64, arg5 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.BitPosSpan(ctx, arg1, arg2, arg3, arg4, arg5)
}
func (nc *NamedClient) CFInfo(ctx context.Context, arg1 string) *redis.CFInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CFInfo(ctx, arg1)
}
func (nc *NamedClient) CFInsert(ctx context.Context, arg1 string, arg2 *redis.CFInsertOptions, arg3 ...interface{}) *redis.BoolSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CFInsert(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) CFInsertNX(ctx context.Context, arg1 string, arg2 *redis.CFInsertOptions, arg3 ...interface{}) *redis.IntSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CFInsertNX(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) CFLoadChunk(ctx context.Context, arg1 string, arg2 int64, arg3 interface{}) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CFLoadChunk(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) CFReserveBucketSize(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CFReserveBucketSize(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) CFReserveExpansion(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CFReserveExpansion(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) CFReserveMaxIterations(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CFReserveMaxIterations(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) CMSInfo(ctx context.Context, arg1 string) *redis.CMSInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CMSInfo(ctx, arg1)
}
func (nc *NamedClient) CMSInitByDim(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CMSInitByDim(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) CMSInitByProb(ctx context.Context, arg1 string, arg2 float64, arg3 float64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.CMSInitByProb(ctx, arg1, arg2, arg3)
}

func (nc *NamedClient) Copy(ctx context.Context, arg1 string, arg2 string, arg3 int, arg4 bool) *redis.IntCmd {
	arg1 = nc.Name + arg1
	arg2 = nc.Name + arg2

	return nc.UniversalClient.Copy(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) DebugObject(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.DebugObject(ctx, arg1)
}
func (nc *NamedClient) Decr(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Decr(ctx, arg1)
}
func (nc *NamedClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.Del(ctx, keys...)
}

func (nc *NamedClient) Dump(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Dump(ctx, arg1)
}
func (nc *NamedClient) Eval(ctx context.Context, arg1 string, arg2 []string, arg3 ...interface{}) *redis.Cmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Eval(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) EvalRO(ctx context.Context, arg1 string, arg2 []string, arg3 ...interface{}) *redis.Cmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.EvalRO(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) EvalSha(ctx context.Context, arg1 string, arg2 []string, arg3 ...interface{}) *redis.Cmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.EvalSha(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) EvalShaRO(ctx context.Context, arg1 string, arg2 []string, arg3 ...interface{}) *redis.Cmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.EvalShaRO(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}

	return nc.UniversalClient.Exists(ctx, keys...)
}
func (nc *NamedClient) ExpireTime(ctx context.Context, arg1 string) *redis.DurationCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ExpireTime(ctx, arg1)
}

func (nc *NamedClient) GeoDist(ctx context.Context, arg1 string, arg2 string, arg3 string, arg4 string) *redis.FloatCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.GeoDist(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) GeoRadius(ctx context.Context, arg1 string, arg2 float64, arg3 float64, arg4 *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.GeoRadius(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) GeoRadiusByMember(ctx context.Context, arg1 string, arg2 string, arg3 *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.GeoRadiusByMember(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) GeoRadiusByMemberStore(ctx context.Context, arg1 string, arg2 string, arg3 *redis.GeoRadiusQuery) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.GeoRadiusByMemberStore(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) GeoRadiusStore(ctx context.Context, arg1 string, arg2 float64, arg3 float64, arg4 *redis.GeoRadiusQuery) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.GeoRadiusStore(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) GeoSearchStore(ctx context.Context, arg1 string, arg2 string, arg3 *redis.GeoSearchStoreQuery) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.GeoSearchStore(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) Get(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Get(ctx, arg1)
}
func (nc *NamedClient) GetDel(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.GetDel(ctx, arg1)
}
func (nc *NamedClient) GetRange(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.GetRange(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) HGetAll(ctx context.Context, arg1 string) *redis.MapStringStringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.HGetAll(ctx, arg1)
}
func (nc *NamedClient) HIncrBy(ctx context.Context, arg1 string, arg2 string, arg3 int64) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.HIncrBy(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) HIncrByFloat(ctx context.Context, arg1 string, arg2 string, arg3 float64) *redis.FloatCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.HIncrByFloat(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) HKeys(ctx context.Context, arg1 string) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.HKeys(ctx, arg1)
}
func (nc *NamedClient) HLen(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.HLen(ctx, arg1)
}
func (nc *NamedClient) HScan(ctx context.Context, arg1 string, arg2 uint64, arg3 string, arg4 int64) *redis.ScanCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.HScan(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) HSetNX(ctx context.Context, arg1 string, arg2 string, arg3 interface{}) *redis.BoolCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.HSetNX(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) HVals(ctx context.Context, arg1 string) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.HVals(ctx, arg1)
}
func (nc *NamedClient) Incr(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Incr(ctx, arg1)
}
func (nc *NamedClient) JSONArrAppend(ctx context.Context, arg1 string, arg2 string, arg3 ...interface{}) *redis.IntSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONArrAppend(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) JSONArrIndex(ctx context.Context, arg1 string, arg2 string, arg3 ...interface{}) *redis.IntSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONArrIndex(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) JSONArrIndexWithArgs(ctx context.Context, arg1 string, arg2 string, arg3 *redis.JSONArrIndexArgs, arg4 ...interface{}) *redis.IntSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONArrIndexWithArgs(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) JSONArrInsert(ctx context.Context, arg1 string, arg2 string, arg3 int64, arg4 ...interface{}) *redis.IntSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONArrInsert(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) JSONArrPop(ctx context.Context, arg1 string, arg2 string, arg3 int) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONArrPop(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) JSONArrTrimWithArgs(ctx context.Context, arg1 string, arg2 string, arg3 *redis.JSONArrTrimArgs) *redis.IntSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONArrTrimWithArgs(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) JSONGetWithArgs(ctx context.Context, arg1 string, arg2 *redis.JSONGetArgs, arg3 ...string) *redis.JSONCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONGetWithArgs(ctx, arg1, arg2, arg3...)
}
func (nc *NamedClient) JSONMSet(ctx context.Context, kv ...interface{}) *redis.StatusCmd {
	for i := range kv {
		if i%2 == 0 {
			key := kv[i].(string)
			kv[i] = nc.Name + key
		}
	}
	return nc.UniversalClient.JSONMSet(ctx, kv...)
}

//	func (nc *NamedClient)JSONMSetArgs(ctx context.Context,arg1 []redis.JSONSetArgs) *redis.StatusCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.JSONMSetArgs(ctx, arg1)
//	}
func (nc *NamedClient) JSONMerge(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONMerge(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) JSONNumIncrBy(ctx context.Context, arg1 string, arg2 string, arg3 float64) *redis.JSONCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONNumIncrBy(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) JSONSet(ctx context.Context, arg1 string, arg2 string, arg3 interface{}) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONSet(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) JSONSetMode(ctx context.Context, arg1 string, arg2 string, arg3 interface{}, arg4 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONSetMode(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) JSONStrAppend(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.IntPointerSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.JSONStrAppend(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) Keys(ctx context.Context, arg1 string) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Keys(ctx, arg1)
}

//	func (nc *NamedClient)LCS(ctx context.Context,arg1 *redis.LCSQuery) *redis.LCSCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.LCS(ctx, arg1)
//	}
func (nc *NamedClient) LInsert(ctx context.Context, arg1 string, arg2 string, arg3 interface{}, arg4 interface{}) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LInsert(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) LInsertAfter(ctx context.Context, arg1 string, arg2 interface{}, arg3 interface{}) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LInsertAfter(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) LInsertBefore(ctx context.Context, arg1 string, arg2 interface{}, arg3 interface{}) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LInsertBefore(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) LLen(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LLen(ctx, arg1)
}
func (nc *NamedClient) LMPop(ctx context.Context, arg1 string, arg2 int64, keys ...string) *redis.KeyValuesCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}

	return nc.UniversalClient.LMPop(ctx, arg1, arg2, keys...)
}
func (nc *NamedClient) LMove(ctx context.Context, arg1 string, arg2 string, arg3 string, arg4 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LMove(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) LPop(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LPop(ctx, arg1)
}
func (nc *NamedClient) LPos(ctx context.Context, arg1 string, arg2 string, arg3 redis.LPosArgs) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LPos(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) LPosCount(ctx context.Context, arg1 string, arg2 string, arg3 int64, arg4 redis.LPosArgs) *redis.IntSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LPosCount(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) LRange(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LRange(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) LRem(ctx context.Context, arg1 string, arg2 int64, arg3 interface{}) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LRem(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) LSet(ctx context.Context, arg1 string, arg2 int64, arg3 interface{}) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LSet(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) LTrim(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.LTrim(ctx, arg1, arg2, arg3)
}

func (nc *NamedClient) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.MGet(ctx, keys...)
}

func (nc *NamedClient) MSet(ctx context.Context, kv ...interface{}) *redis.StatusCmd {
	for i := range kv {
		if i%2 == 0 {
			key := kv[i].(string)
			kv[i] = nc.Name + key
		}
	}

	return nc.UniversalClient.MSet(ctx, kv...)
}
func (nc *NamedClient) MSetNX(ctx context.Context, kv ...interface{}) *redis.BoolCmd {
	for i := range kv {
		if i%2 == 0 {
			key := kv[i].(string)
			kv[i] = nc.Name + key
		}
	}

	return nc.UniversalClient.MSetNX(ctx, kv...)
}
func (nc *NamedClient) Migrate(ctx context.Context, arg1 string, arg2 string, key string, arg4 int, arg5 time.Duration) *redis.StatusCmd {
	key = nc.Name + key

	return nc.UniversalClient.Migrate(ctx, arg1, arg2, key, arg4, arg5)
}

//	func (nc *NamedClient)ModuleLoadex(ctx context.Context,arg1 *redis.ModuleLoadexConfig) *redis.StringCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.ModuleLoadex(ctx, arg1)
//	}
func (nc *NamedClient) ObjectEncoding(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ObjectEncoding(ctx, arg1)
}
func (nc *NamedClient) ObjectFreq(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ObjectFreq(ctx, arg1)
}
func (nc *NamedClient) ObjectIdleTime(ctx context.Context, arg1 string) *redis.DurationCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ObjectIdleTime(ctx, arg1)
}
func (nc *NamedClient) ObjectRefCount(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ObjectRefCount(ctx, arg1)
}
func (nc *NamedClient) PExpireTime(ctx context.Context, arg1 string) *redis.DurationCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.PExpireTime(ctx, arg1)
}
func (nc *NamedClient) PFCount(ctx context.Context, keys ...string) *redis.IntCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.PFCount(ctx, keys...)
}

func (nc *NamedClient) PTTL(ctx context.Context, arg1 string) *redis.DurationCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.PTTL(ctx, arg1)
}
func (nc *NamedClient) Persist(ctx context.Context, arg1 string) *redis.BoolCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Persist(ctx, arg1)
}

func (nc *NamedClient) RPop(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.RPop(ctx, arg1)
}

func (nc *NamedClient) Restore(ctx context.Context, arg1 string, arg2 time.Duration, arg3 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Restore(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) RestoreReplace(ctx context.Context, arg1 string, arg2 time.Duration, arg3 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.RestoreReplace(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) SCard(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SCard(ctx, arg1)
}
func (nc *NamedClient) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.SDiff(ctx, keys...)
}
func (nc *NamedClient) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.SInter(ctx, keys...)
}
func (nc *NamedClient) SMembers(ctx context.Context, arg1 string) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SMembers(ctx, arg1)
}
func (nc *NamedClient) SMembersMap(ctx context.Context, arg1 string) *redis.StringStructMapCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SMembersMap(ctx, arg1)
}
func (nc *NamedClient) SMove(ctx context.Context, arg1 string, arg2 string, arg3 interface{}) *redis.BoolCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SMove(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) SPop(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SPop(ctx, arg1)
}
func (nc *NamedClient) SRandMember(ctx context.Context, arg1 string) *redis.StringCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SRandMember(ctx, arg1)
}
func (nc *NamedClient) SScan(ctx context.Context, arg1 string, arg2 uint64, arg3 string, arg4 int64) *redis.ScanCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SScan(ctx, arg1, arg2, arg3, arg4)
}

func (nc *NamedClient) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.SUnion(ctx, keys...)
}
func (nc *NamedClient) Scan(ctx context.Context, arg1 uint64, arg2 string, arg3 int64) *redis.ScanCmd {
	arg2 = nc.Name + arg2

	return nc.UniversalClient.Scan(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ScanType(ctx context.Context, arg1 uint64, arg2 string, arg3 int64, arg4 string) *redis.ScanCmd {
	arg2 = nc.Name + arg2

	return nc.UniversalClient.ScanType(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) Set(ctx context.Context, arg1 string, arg2 interface{}, arg3 time.Duration) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Set(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) SetArgs(ctx context.Context, arg1 string, arg2 interface{}, arg3 redis.SetArgs) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SetArgs(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) SetBit(ctx context.Context, arg1 string, arg2 int64, arg3 int) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SetBit(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) SetEx(ctx context.Context, arg1 string, arg2 interface{}, arg3 time.Duration) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SetEx(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) SetNX(ctx context.Context, arg1 string, arg2 interface{}, arg3 time.Duration) *redis.BoolCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SetNX(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) SetRange(ctx context.Context, arg1 string, arg2 int64, arg3 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SetRange(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) SetXX(ctx context.Context, arg1 string, arg2 interface{}, arg3 time.Duration) *redis.BoolCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SetXX(ctx, arg1, arg2, arg3)
}

func (nc *NamedClient) SortStore(ctx context.Context, arg1 string, arg2 string, arg3 *redis.Sort) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.SortStore(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) StrLen(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.StrLen(ctx, arg1)
}

func (nc *NamedClient) TDigestCreate(ctx context.Context, arg1 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TDigestCreate(ctx, arg1)
}
func (nc *NamedClient) TDigestInfo(ctx context.Context, arg1 string) *redis.TDigestInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TDigestInfo(ctx, arg1)
}
func (nc *NamedClient) TDigestMax(ctx context.Context, arg1 string) *redis.FloatCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TDigestMax(ctx, arg1)
}
func (nc *NamedClient) TDigestMerge(ctx context.Context, arg1 string, arg2 *redis.TDigestMergeOptions, keys ...string) *redis.StatusCmd {
	arg1 = nc.Name + arg1
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}

	return nc.UniversalClient.TDigestMerge(ctx, arg1, arg2, keys...)
}
func (nc *NamedClient) TDigestMin(ctx context.Context, arg1 string) *redis.FloatCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TDigestMin(ctx, arg1)
}
func (nc *NamedClient) TDigestReset(ctx context.Context, arg1 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TDigestReset(ctx, arg1)
}
func (nc *NamedClient) TDigestTrimmedMean(ctx context.Context, arg1 string, arg2 float64, arg3 float64) *redis.FloatCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TDigestTrimmedMean(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TFCall(ctx context.Context, arg1 string, arg2 string, arg3 int) *redis.Cmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TFCall(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TFCallASYNC(ctx context.Context, arg1 string, arg2 string, arg3 int) *redis.Cmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TFCallASYNC(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TFCallASYNCArgs(ctx context.Context, arg1 string, arg2 string, arg3 int, arg4 *redis.TFCallOptions) *redis.Cmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TFCallASYNCArgs(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) TFCallArgs(ctx context.Context, arg1 string, arg2 string, arg3 int, arg4 *redis.TFCallOptions) *redis.Cmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TFCallArgs(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) TFunctionDelete(ctx context.Context, arg1 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TFunctionDelete(ctx, arg1)
}

//	func (nc *NamedClient)TFunctionListArgs(ctx context.Context,arg1 *redis.TFunctionListOptions) *redis.MapStringInterfaceSliceCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.TFunctionListArgs(ctx, arg1)
//	}
func (nc *NamedClient) TFunctionLoad(ctx context.Context, arg1 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TFunctionLoad(ctx, arg1)
}
func (nc *NamedClient) TSAdd(ctx context.Context, arg1 string, arg2 interface{}, arg3 float64) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSAdd(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TSAddWithArgs(ctx context.Context, arg1 string, arg2 interface{}, arg3 float64, arg4 *redis.TSOptions) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSAddWithArgs(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) TSCreate(ctx context.Context, arg1 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSCreate(ctx, arg1)
}
func (nc *NamedClient) TSCreateRule(ctx context.Context, arg1 string, arg2 string, arg3 redis.Aggregator, arg4 int) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSCreateRule(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) TSCreateRuleWithArgs(ctx context.Context, arg1 string, arg2 string, arg3 redis.Aggregator, arg4 int, arg5 *redis.TSCreateRuleOptions) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSCreateRuleWithArgs(ctx, arg1, arg2, arg3, arg4, arg5)
}
func (nc *NamedClient) TSDecrByWithArgs(ctx context.Context, arg1 string, arg2 float64, arg3 *redis.TSIncrDecrOptions) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSDecrByWithArgs(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TSDel(ctx context.Context, arg1 string, arg2 int, arg3 int) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSDel(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TSGet(ctx context.Context, arg1 string) *redis.TSTimestampValueCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSGet(ctx, arg1)
}
func (nc *NamedClient) TSIncrByWithArgs(ctx context.Context, arg1 string, arg2 float64, arg3 *redis.TSIncrDecrOptions) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSIncrByWithArgs(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TSInfo(ctx context.Context, arg1 string) *redis.MapStringInterfaceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSInfo(ctx, arg1)
}
func (nc *NamedClient) TSMGet(ctx context.Context, keys []string) *redis.MapStringSliceInterfaceCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.TSMGet(ctx, keys)
}

func (nc *NamedClient) TSRange(ctx context.Context, arg1 string, arg2 int, arg3 int) *redis.TSTimestampValueSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSRange(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TSRangeWithArgs(ctx context.Context, arg1 string, arg2 int, arg3 int, arg4 *redis.TSRangeOptions) *redis.TSTimestampValueSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSRangeWithArgs(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) TSRevRange(ctx context.Context, arg1 string, arg2 int, arg3 int) *redis.TSTimestampValueSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSRevRange(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) TSRevRangeWithArgs(ctx context.Context, arg1 string, arg2 int, arg3 int, arg4 *redis.TSRevRangeOptions) *redis.TSTimestampValueSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TSRevRangeWithArgs(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) TTL(ctx context.Context, arg1 string) *redis.DurationCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TTL(ctx, arg1)
}
func (nc *NamedClient) TopKInfo(ctx context.Context, arg1 string) *redis.TopKInfoCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TopKInfo(ctx, arg1)
}
func (nc *NamedClient) TopKList(ctx context.Context, arg1 string) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TopKList(ctx, arg1)
}
func (nc *NamedClient) TopKListWithCount(ctx context.Context, arg1 string) *redis.MapStringIntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TopKListWithCount(ctx, arg1)
}
func (nc *NamedClient) TopKReserveWithOptions(ctx context.Context, arg1 string, arg2 int64, arg3 int64, arg4 int64, arg5 float64) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.TopKReserveWithOptions(ctx, arg1, arg2, arg3, arg4, arg5)
}
func (nc *NamedClient) Touch(ctx context.Context, keys ...string) *redis.IntCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}

	return nc.UniversalClient.Touch(ctx, keys...)
}
func (nc *NamedClient) Type(ctx context.Context, arg1 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.Type(ctx, arg1)
}
func (nc *NamedClient) Unlink(ctx context.Context, keys ...string) *redis.IntCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}

	return nc.UniversalClient.Unlink(ctx, keys...)
}
func (nc *NamedClient) XAck(ctx context.Context, arg1 string, arg2 string, arg3 ...string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XAck(ctx, arg1, arg2, arg3...)
}

//	func (nc *NamedClient)XAdd(ctx context.Context,arg1 *redis.XAddArgs) *redis.StringCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.XAdd(ctx, arg1)
//	}
//
//	func (nc *NamedClient)XAutoClaim(ctx context.Context,arg1 *redis.XAutoClaimArgs) *redis.XAutoClaimCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.XAutoClaim(ctx, arg1)
//	}
//
//	func (nc *NamedClient)XAutoClaimJustID(ctx context.Context,arg1 *redis.XAutoClaimArgs) *redis.XAutoClaimJustIDCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.XAutoClaimJustID(ctx, arg1)
//	}
//
//	func (nc *NamedClient)XClaim(ctx context.Context,arg1 *redis.XClaimArgs) *redis.XMessageSliceCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.XClaim(ctx, arg1)
//	}
//
//	func (nc *NamedClient)XClaimJustID(ctx context.Context,arg1 *redis.XClaimArgs) *redis.StringSliceCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.XClaimJustID(ctx, arg1)
//	}
func (nc *NamedClient) XGroupCreate(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XGroupCreate(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) XGroupCreateConsumer(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XGroupCreateConsumer(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) XGroupCreateMkStream(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XGroupCreateMkStream(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) XGroupDelConsumer(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XGroupDelConsumer(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) XGroupSetID(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.StatusCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XGroupSetID(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) XInfoGroups(ctx context.Context, arg1 string) *redis.XInfoGroupsCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XInfoGroups(ctx, arg1)
}
func (nc *NamedClient) XInfoStream(ctx context.Context, arg1 string) *redis.XInfoStreamCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XInfoStream(ctx, arg1)
}
func (nc *NamedClient) XLen(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XLen(ctx, arg1)
}

//	func (nc *NamedClient)XPendingExt(ctx context.Context,arg1 *redis.XPendingExtArgs) *redis.XPendingExtCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.XPendingExt(ctx, arg1)
//	}
func (nc *NamedClient) XRange(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.XMessageSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XRange(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) XRangeN(ctx context.Context, arg1 string, arg2 string, arg3 string, arg4 int64) *redis.XMessageSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XRangeN(ctx, arg1, arg2, arg3, arg4)
}

//	func (nc *NamedClient)XRead(ctx context.Context,arg1 *redis.XReadArgs) *redis.XStreamSliceCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.XRead(ctx, arg1)
//	}
//
//	func (nc *NamedClient)XReadGroup(ctx context.Context,arg1 *redis.XReadGroupArgs) *redis.XStreamSliceCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.XReadGroup(ctx, arg1)
//	}
func (nc *NamedClient) XReadStreams(ctx context.Context, keys ...string) *redis.XStreamSliceCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}

	return nc.UniversalClient.XReadStreams(ctx, keys...)
}
func (nc *NamedClient) XRevRange(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.XMessageSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XRevRange(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) XRevRangeN(ctx context.Context, arg1 string, arg2 string, arg3 string, arg4 int64) *redis.XMessageSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XRevRangeN(ctx, arg1, arg2, arg3, arg4)
}
func (nc *NamedClient) XTrimMaxLenApprox(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XTrimMaxLenApprox(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) XTrimMinIDApprox(ctx context.Context, arg1 string, arg2 string, arg3 int64) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.XTrimMinIDApprox(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZCard(ctx context.Context, arg1 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZCard(ctx, arg1)
}
func (nc *NamedClient) ZCount(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZCount(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.ZDiff(ctx, keys...)
}
func (nc *NamedClient) ZDiffWithScores(ctx context.Context, keys ...string) *redis.ZSliceCmd {
	for i := range keys {
		keys[i] = nc.Name + keys[i]
	}
	return nc.UniversalClient.ZDiffWithScores(ctx, keys...)
}
func (nc *NamedClient) ZIncrBy(ctx context.Context, arg1 string, arg2 float64, arg3 string) *redis.FloatCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZIncrBy(ctx, arg1, arg2, arg3)
}

// func (nc *NamedClient)ZInter(ctx context.Context,arg1 *redis.ZStore) *redis.StringSliceCmd {
// //	arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.ZInter(ctx, arg1)
//	}
//
//	func (nc *NamedClient)ZInterWithScores(ctx context.Context,arg1 *redis.ZStore) *redis.ZSliceCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.ZInterWithScores(ctx, arg1)
//	}
func (nc *NamedClient) ZLexCount(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZLexCount(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZMPop(ctx context.Context, arg1 string, arg2 int64, arg3 ...string) *redis.ZSliceWithKeyCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZMPop(ctx, arg1, arg2, arg3...)
}
func (nc *NamedClient) ZRange(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZRange(ctx, arg1, arg2, arg3)
}

//	func (nc *NamedClient)ZRangeArgs(ctx context.Context,arg1 redis.ZRangeArgs) *redis.StringSliceCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.ZRangeArgs(ctx, arg1)
//	}
//
//	func (nc *NamedClient)ZRangeArgsWithScores(ctx context.Context,arg1 redis.ZRangeArgs) *redis.ZSliceCmd {
//		arg1 = nc.Name + arg1
//
//		return nc.UniversalClient.ZRangeArgsWithScores(ctx, arg1)
//	}
func (nc *NamedClient) ZRangeWithScores(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.ZSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZRangeWithScores(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZRemRangeByLex(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZRemRangeByLex(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZRemRangeByRank(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZRemRangeByRank(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZRemRangeByScore(ctx context.Context, arg1 string, arg2 string, arg3 string) *redis.IntCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZRemRangeByScore(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZRevRange(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.StringSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZRevRange(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZRevRangeWithScores(ctx context.Context, arg1 string, arg2 int64, arg3 int64) *redis.ZSliceCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZRevRangeWithScores(ctx, arg1, arg2, arg3)
}
func (nc *NamedClient) ZScan(ctx context.Context, arg1 string, arg2 uint64, arg3 string, arg4 int64) *redis.ScanCmd {
	arg1 = nc.Name + arg1

	return nc.UniversalClient.ZScan(ctx, arg1, arg2, arg3, arg4)
}

//func (nc *NamedClient)ZUnion(ctx context.Context,arg1 redis.ZStore) *redis.StringSliceCmd {
//	arg1 = nc.Name + arg1
//
//	return nc.UniversalClient.ZUnion(ctx, arg1)
//}
//func (nc *NamedClient)ZUnionWithScores(ctx context.Context,arg1 redis.ZStore) *redis.ZSliceCmd {
//	arg1 = nc.Name + arg1
//
//	return nc.UniversalClient.ZUnionWithScores(ctx, arg1)
//}
