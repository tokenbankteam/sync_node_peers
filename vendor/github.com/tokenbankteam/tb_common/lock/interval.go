package lock

import "math"

type Interval struct {
	IntervalTime int64 //执行锁定间隔

	LockTime int64 //锁定时间

	Unit int64 //单位 换算成毫秒， 比如intervalTime=10, lockTime=1, unit=1000， 表示间隔10*1000=10s， 锁定时间1*1000=1s
}

var THREE_SECOND_INTERVAL *Interval = &Interval{IntervalTime: 3, LockTime: 3, Unit: 1000}

var FIVE_SECOND_INTERVAL *Interval = &Interval{IntervalTime: 5, LockTime: 5, Unit: 1000}

var TEN_SECOND_INTERVAL *Interval = &Interval{IntervalTime: 10, LockTime: 10, Unit: 1000}

var MIN_INTERVAL *Interval = &Interval{IntervalTime: 60, LockTime: 60, Unit: 1000}

var TWO_MIN_INTERVAL *Interval = &Interval{IntervalTime: 2 * 60, LockTime: 2 * 60, Unit: 1000}

var FIVE_MIN_INTERVAL *Interval = &Interval{IntervalTime: 5 * 60, LockTime: 5 * 60, Unit: 1000}

var TEN_MIN_INTERVAL *Interval = &Interval{IntervalTime: 10 * 60, LockTime: 120, Unit: 1000}

var TWO_HOURS_INTERVAL *Interval = &Interval{IntervalTime: 2 * 60 * 60, LockTime: 600, Unit: 1000}

var INFINITE_INTERVAL *Interval = &Interval{IntervalTime: math.MaxInt64, LockTime: math.MaxInt64, Unit: 1000}
