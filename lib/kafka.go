package lib

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	defaultPartition = 0
	defaultAddress   = "localhost:9092"
	network          = "tcp"

	writeDeadline = 10 * time.Second
	readDeadline  = 10 * time.Second

	// fetch 10KB min, 1MB max
	readMinBytes = 10e3
	readMaxBytes = 1e6

	msgBytesLimit = 10e3 // 10KB max per message
)

// Produce to produce messages
func Produce(ctx context.Context, topic string, messages []string) error {
	conn, err := kafka.DialLeader(ctx, network, defaultAddress, topic, defaultPartition)
	if err != nil {
		log.Println("failed to dial leader:", err)
		return err
	}

	err = conn.SetWriteDeadline(time.Now().Add(writeDeadline))
	if err != nil {
		log.Println("failed to set write deadline:", err)
		return err
	}

	kafkaMsg := make([]kafka.Message, 0, len(messages))
	for _, msg := range messages {
		kafkaMsg = append(kafkaMsg, kafka.Message{Value: []byte(msg)})
	}
	_, err = conn.WriteMessages(kafkaMsg...)
	if err != nil {
		log.Println("failed to write messages:", err)
		return err
	}

	if err = conn.Close(); err != nil {
		log.Println("failed to close writer:", err)
		return err
	}

	return nil
}

// Consume to consume messages
func Consume(ctx context.Context, topic string, nums int) ([]string, error) {
	conn, err := kafka.DialLeader(ctx, network, defaultAddress, topic, defaultPartition)
	if err != nil {
		log.Println("failed to dial leader:", err)
		return nil, err
	}

	err = conn.SetReadDeadline(time.Now().Add(readDeadline))
	if err != nil {
		log.Println("failed to set read deadline:", err)
		return nil, err
	}

	batch := conn.ReadBatch(readMinBytes, readMaxBytes)
	b := make([]byte, msgBytesLimit)
	res := make([]string, 0, nums)
	for i := 0; i < nums; i++ {
		n, err := batch.Read(b)
		if err != nil {
			log.Println("read message failed:", err)
			break
		}
		res = append(res, string(b[:n]))
	}

	if err = batch.Close(); err != nil {
		log.Println("failed to close batch:", err)
		return nil, err
	}

	if err = conn.Close(); err != nil {
		log.Println("failed to close connection:", err)
		return nil, err
	}

	return res, nil
}

func GroupConsume(ctx context.Context, topic string, nums int) ([]string, error) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		GroupID:  "consumer-group-id-lucky-12580",
		Topic:    topic,
		MaxBytes: readMaxBytes,
	})

	res := make([]string, 0, nums)
	for i := 0; i < nums; i++ {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("read message err", err)
			break
		}
		log.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition,
			m.Offset, string(m.Key), string(m.Value))
		res = append(res, string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Println("failed to close reader:", err)
		return []string{"reader close err"}, err
	}
	return res, nil
}

// CreateTopic to create topics when auto.create.topics.enable='false'
func CreateTopic(ctx context.Context, topics []string) error {
	conn, err := kafka.Dial(network, defaultAddress)
	if err != nil {
		log.Println("failed to dial:", err)
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Println("failed to get controller:", err)
		return err
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial(network, net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Println("failed to dial controller:", err)
		return err
	}
	defer controllerConn.Close()

	topicConfigs := make([]kafka.TopicConfig, 0, len(topics))
	for _, topic := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Println("failed to create topic:", err)
		return err
	}

	return nil
}
