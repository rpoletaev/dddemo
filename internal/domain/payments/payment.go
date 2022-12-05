package payments

import "time"

type category string
type status string

const (
	categoryBinding category = "binding"
	categoryPayment category = "payment"

	statusCreated   status = "created"   // a payment created on local storage but doesn't send to the payments service
	statusWaiting   status = "waiting"   // a payment is sent to the payments service but is waiting for processing result
	statusSuccessed status = "successed" // successed payment's result is received
	statusFailed    status = "failed"    // failed payment's result is received
)

type Payment struct {
	uuid            string
	externalID      string
	subscriptionID  uint64
	amount          uint
	category        category
	status          status
	failReason      string
	createdAt       time.Time
	statusUpdatedAt time.Time
}

func New(
	uuid string,
	externalID string,
	subscriptionID uint64,
	amount uint,
	category category,
	status status,
	failReason string,
	createdAt time.Time,
	statusUpdatedAt time.Time,
) Payment {
	return Payment{
		uuid:            uuid,
		externalID:      externalID,
		subscriptionID:  subscriptionID,
		amount:          amount,
		category:        category,
		status:          status,
		failReason:      failReason,
		createdAt:       createdAt,
		statusUpdatedAt: statusUpdatedAt,
	}
}

func (p Payment) Uuid() string {
	return p.uuid
}
func (p Payment) ExternalId() string {
	return p.externalID
}
func (p Payment) SubscriptionId() uint64 {
	return p.subscriptionID
}
func (p Payment) Category() category {
	return p.category
}
func (p Payment) Status() status {
	return p.status
}
func (p Payment) CreatedAt() time.Time {
	return p.createdAt
}
func (p Payment) StatusUpdatedAt() time.Time {
	return p.statusUpdatedAt
}
func (p Payment) FailReason() string {
	return p.failReason
}

func InitNewBinding(uuid string, subscriptionID uint64, amount uint, dateTime time.Time) Payment {
	return New(
		uuid,
		"",
		subscriptionID,
		amount,
		categoryBinding,
		statusCreated,
		"",
		dateTime,
		dateTime,
	)
}

func InitNewPayment(uuid string, subscriptionID uint64, amount uint, dateTime time.Time) Payment {
	return New(
		uuid,
		"",
		subscriptionID,
		amount,
		categoryPayment,
		statusCreated,
		"",
		dateTime,
		dateTime,
	)
}

func (p Payment) Register(externalID string, dateTime time.Time) {
	p.externalID = externalID
	p.status = statusWaiting
	p.statusUpdatedAt = dateTime
}

func (p Payment) Succeessed(dateTime time.Time) {
	p.status = statusSuccessed
	p.statusUpdatedAt = dateTime
}

func (p Payment) Failed(reason string, dateTime time.Time) {
	p.status = statusFailed
	p.failReason = reason
	p.statusUpdatedAt = dateTime
}
