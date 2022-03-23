package messagebroker

import (
	"fmt"

	"github.com/nsqio/go-nsq"
)

// ProducerImpl is
type ProducerImpl struct {
	defaultProducer *nsq.Producer
	namespace       string
}

func producerImplInit(url string) *nsq.Producer {
	nsqConfig := nsq.NewConfig()
	nsqProducer, err := nsq.NewProducer(url, nsqConfig)
	if err != nil {
		panic(err.Error())
	}
	return nsqProducer
}

func (p *ProducerImpl) Publish(topic string, body []byte) error {
	if p.namespace != "" {
		topic = fmt.Sprintf("%s.%s", p.namespace, topic)
	}
	return p.defaultProducer.Publish(topic, body)
}

// NewProducer is
func NewProducer(url, namespace string) Producer {
	return &ProducerImpl{
		defaultProducer: producerImplInit(url),
		namespace:       namespace,
	}
}
