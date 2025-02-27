// Code generated by "stringer -type=ErrorCode --linecomment"; DO NOT EDIT.

package code

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unexpected-1000]
	_ = x[OutOfRange-1001]
	_ = x[SignFailed-1002]
	_ = x[AssertFailed-1003]
	_ = x[ParseFailed-1004]
	_ = x[RequestFailed-1005]
	_ = x[InvalidResponse-1006]
	_ = x[GetAddressFailed-1007]
	_ = x[SelectAddressFailed-1008]
	_ = x[NoValidAddress-1009]
	_ = x[NoValidProduct-1010]
	_ = x[NoReserveTime-1011]
	_ = x[NoReserveTimeAndRetry-1012]
	_ = x[ReserveTimeIsDisabled-1013]
}

const _ErrorCode_name = "意外错误数组索引越界签名失败断言失败解析失败请求失败无效的响应获取收货地址失败选择收货地址错误当前没有可用的收货地址当前购物车中没有可购商品当前没有可用的配送时段当前没有可用的配送时段, 请稍后再试您选择的送达时间已经失效, 请重新选择"

var _ErrorCode_index = [...]uint16{0, 12, 30, 42, 54, 66, 78, 93, 117, 141, 174, 210, 243, 293, 346}

func (i ErrorCode) String() string {
	i -= 1000
	if i < 0 || i >= ErrorCode(len(_ErrorCode_index)-1) {
		return "ErrorCode(" + strconv.FormatInt(int64(i+1000), 10) + ")"
	}
	return _ErrorCode_name[_ErrorCode_index[i]:_ErrorCode_index[i+1]]
}
