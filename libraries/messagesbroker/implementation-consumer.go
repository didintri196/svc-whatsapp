package messagebroker

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

type consumerImpl struct {
	url              string
	namespace        string
	consumerHandlers []ConsumerHandler
}

func (c *consumerImpl) StartListening() {

	nsqConfig := nsq.NewConfig()

	for _, ch := range c.consumerHandlers {

		consumer, err := nsq.NewConsumer(ch.Topic, ch.Channel, nsqConfig)
		if err != nil {
			panic(err.Error())
		}

		consumer.AddHandler(&internalConsumer{handler: ch.FunctionHandler})

		if err := consumer.ConnectToNSQD(c.url); err != nil {
			panic(err.Error())
		}

	}

}

func (c *consumerImpl) Handle(topic, channel string, funcHandler FunctionHandler) {

	if funcHandler == nil {
		panic("FunctionHandler must not nil")
	}

	if c.namespace != "" {
		topic = fmt.Sprintf("%s.%s", c.namespace, topic)
	}

	c.consumerHandlers = append(c.consumerHandlers, ConsumerHandler{
		Topic:           topic,
		Channel:         channel,
		FunctionHandler: funcHandler,
	})
}

type internalConsumer struct {
	handler FunctionHandler
}

// HandleMessage is
func (mb *internalConsumer) HandleMessage(m *nsq.Message) error {
	mb.handler(m.Body)
	return nil
}

func NewConsumer(url, namespace string) Consumer {
	return &consumerImpl{
		url:       url,
		namespace: namespace,
	}
}
