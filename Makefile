migration:
	goose -dir=migrations create create_subscriptions sql

gen-swagger:
	swagger generate server -t ./gen/http -f  ./swagger/swagger.yaml --exclude-main -A SubscriptionsService