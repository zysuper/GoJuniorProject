package samarax

import (
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"github.com/IBM/sarama"
	"github.com/prometheus/client_golang/prometheus"
)

type Handler[T any] struct {
	l  logger.LoggerV1
	fn func(msg *sarama.ConsumerMessage, event T) error
	g  prometheus.Gauge
}

func NewHandler[T any](
	l logger.LoggerV1,
	fn func(msg *sarama.ConsumerMessage, event T) error,
	g prometheus.Gauge) *Handler[T] {
	return &Handler[T]{l: l, fn: fn, g: g}
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		// 在这里调用业务处理逻辑
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			// 你也可以在这里引入重试的逻辑
			h.l.Error("反序列消息体失败",
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
		}
		err = h.fn(msg, t)
		if err != nil {
			h.l.Error("处理消息失败",
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
		}
		session.MarkMessage(msg, "")
		h.g.Dec()
		//h.l.Info("Consume Success.")
	}
	return nil
}
