package client

import (
	"errors"
	"github.com/IBM/sarama"
	"time"
)

var ErrUnableToParseOffsetInitial = errors.New("unable to parse offset initial")

func NewOrderKafkaClient(clientDSN string,
	groupId string,
	offsetInitial string, // oldest | newest
	autocommit bool,
	autocommitInterval time.Duration,
) (sarama.ConsumerGroup, error) {

	// Инициализируем клиента кафки
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	offset, err := parseOffsetInitial(offsetInitial)
	if err != nil {
		return nil, err
	}

	config.Consumer.Offsets.Initial = offset

	if autocommit {
		config.Consumer.Offsets.AutoCommit.Enable = true
		config.Consumer.Offsets.AutoCommit.Interval = autocommitInterval
	}

	kafkaClient, err := sarama.NewConsumerGroup([]string{clientDSN}, groupId, config)

	return kafkaClient, err
}

func parseOffsetInitial(offsetInitial string) (int64, error) {
	switch offsetInitial {
	case "oldest":
		return sarama.OffsetOldest, nil
	case "newest":
		return sarama.OffsetNewest, nil
	default:
		return -1, ErrUnableToParseOffsetInitial
	}
}
