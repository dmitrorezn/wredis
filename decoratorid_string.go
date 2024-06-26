// Code generated by "stringer -type=decoratorID"; DO NOT EDIT.

package wredis

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LRUDecorator-1]
	_ = x[LFUDecorator-2]
	_ = x[SingleFlightDecorator-3]
	_ = x[EncodingDecorator-4]
}

const _decoratorID_name = "LRUDecoratorLFUDecoratorSingleFlightDecoratorEncodingDecorator"

var _decoratorID_index = [...]uint8{0, 12, 24, 45, 62}

func (i decoratorID) String() string {
	i -= 1
	if i < 0 || i >= decoratorID(len(_decoratorID_index)-1) {
		return "decoratorID(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _decoratorID_name[_decoratorID_index[i]:_decoratorID_index[i+1]]
}
