package client

import (
	"github.com/IBM/sarama"
	"time"
)

func NewOrderKafkaClient(clientDSN string) (sarama.ConsumerGroup, error) {

	// Инициализируем клиента кафки
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	groupID := "courier-group-id"

	kafkaClient, err := sarama.NewConsumerGroup([]string{clientDSN}, groupID, config)

	return kafkaClient, err

}
