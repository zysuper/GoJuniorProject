package article

import (
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

const TopicReadEvent = "article_read"

type Producer interface {
	ProduceReadEvent(evt ReadEvent) error
}

type ReadEvent struct {
	Aid int64
	Uid int64
}

type BatchReadEvent struct {
	Aids []int64
	Uids []int64
}

type SaramaSyncProducer struct {
	producer sarama.SyncProducer
	g        prometheus.Gauge
	l        logger.LoggerV1
}

func NewSaramaSyncProducer(producer sarama.SyncProducer, l logger.LoggerV1) Producer {
	return &SaramaSyncProducer{producer: producer, g: gauge, l: l}
}

func (s *SaramaSyncProducer) ProduceReadEvent(evt ReadEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: TopicReadEvent,
		Value: sarama.StringEncoder(val),
	})

	if err == nil {
		//s.l.Info("Producer Success.")
		s.g.Inc()
	}
	return err
}
