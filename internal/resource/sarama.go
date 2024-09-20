package resource

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/abdelrahman146/digital-wallet/pkg/api"
	"github.com/abdelrahman146/digital-wallet/pkg/config"
	"github.com/abdelrahman146/digital-wallet/pkg/logger"
)

type Broker interface {
	// CreateTopic creates a new Broker topic
	CreateTopic(ctx context.Context, topicName string, topicDetail sarama.TopicDetail) error
	// DeleteTopic deletes an existing Broker topic
	DeleteTopic(ctx context.Context, topicName string) error
	// Publish sends a message to a Broker topic
	Publish(ctx context.Context, topic string, message []byte) error
	// Close closes the Broker connection
	Close()
}

func InitBroker(group string) *broker {
	s := &broker{}
	s.Init(group)
	return s
}

type broker struct {
	producer      sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup
}

// Init initializes the Broker connection
func (s *broker) Init(group string) {
	conf := sarama.NewConfig()
	kafkaBrokers := config.GetConfig().KafkaBrokers
	conf.Producer.Return.Successes = true
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	producer, err := sarama.NewSyncProducer([]string{kafkaBrokers}, conf)
	if err != nil {
		logger.GetLogger().Panic("Failed to create Broker producer", logger.Field("error", err))
	}
	consumerGroup, err := sarama.NewConsumerGroup([]string{kafkaBrokers}, group, conf)
	if err != nil {
		logger.GetLogger().Panic("Failed to create Broker consumer group", logger.Field("error", err), logger.Field("group", group))
	}
	s.producer = producer
	s.consumerGroup = consumerGroup
}

// Close closes the Broker connection
func (s *broker) Close() {
	if err := s.producer.Close(); err != nil {
		logger.GetLogger().Panic("Failed to close Broker producer", logger.Field("error", err))
	}
	if err := s.consumerGroup.Close(); err != nil {
		logger.GetLogger().Panic("Failed to close Broker consumer group", logger.Field("error", err))
	}
}

// CreateTopic creates a new Broker topic
func (s *broker) CreateTopic(ctx context.Context, topicName string, topicDetail sarama.TopicDetail) error {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, config)
	if err != nil {
		api.GetLogger(ctx).Error("Failed to create broker cluster admin", logger.Field("error", err))
		return err
	}
	defer admin.Close()
	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}
	if _, exists := topics[topicName]; exists {
		api.GetLogger(ctx).Info("Broker topic already exists", logger.Field("topic", topicName))
		return nil
	}
	if err := admin.CreateTopic(topicName, &topicDetail, false); err != nil {
		api.GetLogger(ctx).Error("Failed to create broker topic", logger.Field("error", err), logger.Field("topic", topicName))
		return err
	}
	api.GetLogger(ctx).Info("Successfully created Broker topic", logger.Field("topic", topicName))
	return nil
}

// DeleteTopic deletes an existing Broker topic
func (s *broker) DeleteTopic(ctx context.Context, topicName string) error {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin([]string{"localhost:9092"}, config)
	if err != nil {
		api.GetLogger(ctx).Error("Failed to create broker cluster admin", logger.Field("error", err))
		return err
	}
	defer admin.Close()
	topics, err := admin.ListTopics()
	if err != nil {
		return err
	}
	if _, exists := topics[topicName]; !exists {
		api.GetLogger(ctx).Info("Broker topic does not exist", logger.Field("topic", topicName))
		return nil
	}
	if err := admin.DeleteTopic(topicName); err != nil {
		api.GetLogger(ctx).Error("Failed to delete broker topic", logger.Field("error", err), logger.Field("topic", topicName))
		return err
	}
	api.GetLogger(ctx).Info("Successfully deleted Broker topic", logger.Field("topic", topicName))
	return nil
}

// Publish sends a message to a Broker topic
func (s *broker) Publish(ctx context.Context, topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		api.GetLogger(ctx).Error("Failed to send message", logger.Field("error", err))
		return err
	}
	api.GetLogger(ctx).Info("Successfully sent message", logger.Field("topic", topic), logger.Field("message", message), logger.Field("partition", partition), logger.Field("offset", offset))
	return nil
}
