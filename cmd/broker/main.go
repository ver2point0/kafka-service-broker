package main

import (
	"net/http"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi"
	// "github.com/pivotal-cf/brokerapi/auth"

	"github.com/starkandwayne/kafka-service-broker/broker"
	"github.com/starkandwayne/kafka-service-broker/brokerconfig"
	"github.com/starkandwayne/kafka-service-broker/kafka"
)

func main() {
	// brokerConfigPath := configPath()

	brokerLogger := lager.NewLogger("kafka-service-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))

	brokerLogger.Info("Starting Kafka service broker")

	config := brokerconfig.LoadConfig()

	topicRepo := kafka.NewTopicRepository(config.KafkaConfiguration, brokerLogger)

	serviceBroker := &broker.KafkaServiceBroker{
		InstanceCreators: map[string]broker.InstanceCreator{
			"topic": topicRepo,
		},
		InstanceBinders: map[string]broker.InstanceBinder{
			"topic": topicRepo,
		},
		Config: config,
	}

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: "broker",
		Password: "broker",
	}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)

	port := "5000"
	http.Handle("/", brokerAPI)
	brokerLogger.Fatal("http-listen", http.ListenAndServe("0.0.0.0:"+port, nil))
}
