package input

import (
	"context"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/wswz/go_commons/log"
)

type Kafka struct {
	client  sarama.ConsumerGroup
	stopped chan struct{}
	msgs    chan ([]byte)

	Name          string
	Version       string
	Earliest      bool
	Brokers       string
	ConsumerGroup string
	Topic         string
// TOK_ID_KRB_AP_REQ   = 256
// GSS_API_GENERIC_TAG = 0x60
// KRB5_USER_AUTH      = 1
// KRB5_KEYTAB_AUTH    = 2
// GSS_API_INITIAL     = 1
// GSS_API_VERIFY      = 2
// GSS_API_FINISH      = 3
	Sasl struct {
		Username string
		Password string
		Mechanism string
		GSSAPI struct {
			AuthType int
			Username string
			Password string
			ServiceName string
			Realm string
			KerberosConfigPath string
			KeyTabPath string
		}
	}
	consumer *Consumer
	context  context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

func NewKafka() *Kafka {
	return &Kafka{}
}

func (k *Kafka) Init() error {
	k.msgs = make(chan []byte, 300000)
	k.stopped = make(chan struct{})
	k.consumer = &Consumer{
		msgs:  k.msgs,
		ready: make(chan bool),
	}
	k.context, k.cancel = context.WithCancel(context.Background())
	return nil
}

func (k *Kafka) Msgs() chan []byte {
	return k.msgs
}

func (k *Kafka) Start() error {
	config := sarama.NewConfig()

	if k.Version != "" {
		version, err := sarama.ParseKafkaVersion(k.Version)
		if err != nil {
			return err
		}
		config.Version = version
	}
	// sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	if k.Sasl.Username != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = k.Sasl.Username
		config.Net.SASL.Password = k.Sasl.Password
	}

	if k.Sasl.Mechanism != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLMechanism(k.Sasl.Mechanism)
		config.Net.SASL.GSSAPI.AuthType = k.Sasl.GSSAPI.AuthType
		config.Net.SASL.GSSAPI.Username = k.Sasl.GSSAPI.Username
		config.Net.SASL.GSSAPI.ServiceName = k.Sasl.GSSAPI.ServiceName
		config.Net.SASL.GSSAPI.Realm = k.Sasl.GSSAPI.Realm
		config.Net.SASL.GSSAPI.KerberosConfigPath = k.Sasl.GSSAPI.KerberosConfigPath
		config.Net.SASL.GSSAPI.KeyTabPath = k.Sasl.GSSAPI.KeyTabPath
	}

	if k.Earliest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	log.Info("start to dial kafka ", k.Brokers)
	client, err := sarama.NewConsumerGroup(strings.Split(k.Brokers, ","), k.ConsumerGroup, config)
	if err != nil {
		return err
	}

	k.client = client

	go func() {
		k.wg.Add(1)
		defer k.wg.Done()
		for {
			if err := k.client.Consume(k.context, strings.Split(k.Topic, ","), k.consumer); err != nil {
				log.Error("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if k.context.Err() != nil {
				return
			}
			k.consumer.ready = make(chan bool, 0)
		}
	}()

	<-k.consumer.ready
	return nil
}

func (k *Kafka) Stop() error {
	k.cancel()
	k.wg.Wait()

	k.client.Close()
	close(k.msgs)
	return nil
}

func (k *Kafka) Description() string {
	return "kafka consumer:" + k.Topic
}

func (k *Kafka) GetName() string {
	return k.Name
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	msgs  chan []byte
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		consumer.msgs <- message.Value
		session.MarkMessage(message, "")
	}

	return nil
}
