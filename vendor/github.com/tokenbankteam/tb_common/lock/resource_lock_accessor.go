package lock

//资源锁访问接口
type ResourceLockAccessor interface {
	//获取单个锁   没有获得所, 则返回nil
	GetResourceLock(intervalSecond int64, lockTime int64) (*ResourceLock, error)
	GetResourceLockByInterval(interval *Interval) (*ResourceLock, error)

	//一次获取多个锁
	GetResourceLockList(intervalSecond int64, lockTime int64, max int) ([]*ResourceLock, error)
	GetResourceLockListByInterval(interval *Interval, max int) ([]*ResourceLock, error)

	//添加一个锁
	AddResourceLock(resourceLock *ResourceLock) (int64, error)

	//一次添加多个锁
	AddResourceLockList(resourceLockList []*ResourceLock) ([]int64, error)

	//删除一个锁
	DeleteResourceLock(resourceLock *ResourceLock) (int64, error)
}
