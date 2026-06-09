package mercure

import (
	"flowforge-api/infrastructure/config"
	infrastructureMercure "flowforge-api/infrastructure/mercure"
)

func Publish(topic string, data any, env *config.Config) error {
	publisher := infrastructureMercure.NewPublisher(env.MercureURL, env.MercurePublisherJWTKey)
	return publisher.Publish(topic, data)
}
