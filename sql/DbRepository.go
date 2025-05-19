package sql

type DbRepository struct {
	Subscription ISubscriptionRepository
}

func NewDbRepository(subscription ISubscriptionRepository) *DbRepository {
	return &DbRepository{
		Subscription: subscription,
	}
}
