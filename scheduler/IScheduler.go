package scheduler

type IScheduler interface {
	Start()
	GetCache() ISubscriptionCache
}
