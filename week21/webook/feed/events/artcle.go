package events

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/feed/domain"
	"gitee.com/geekbang/basic-go/webook/feed/service"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

const topicArticleEvent = "article_feed_event"

// ArticleFeedEvent 由业务方定义，本服务做适配
type ArticleFeedEvent struct {
	Uid int64
	Aid int64
}

type ArticleEventConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	svc    service.FeedService
}

func NewArticleEventConsumer(
	client sarama.Client,
	l logger.LoggerV1,
	svc service.FeedService) *ArticleEventConsumer {
	ac := &ArticleEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}
	return ac
}

// Start 这边就是自己启动 goroutine 了
func (r *ArticleEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("articleFeed",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicArticleEvent},
			saramax.NewHandler[ArticleFeedEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func (r *ArticleEventConsumer) Consume(msg *sarama.ConsumerMessage,
	evt ArticleFeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: service.FollowEventName,
		Ext: map[string]string{
			"uid": strconv.FormatInt(evt.Uid, 10),
			"aid": strconv.FormatInt(evt.Uid, 10),
		},
	})

}
