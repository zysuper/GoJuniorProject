package events

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/follow/repository"
	"gitee.com/geekbang/basic-go/webook/follow/repository/dao"
	"gitee.com/geekbang/basic-go/webook/pkg/canalx"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"gitee.com/geekbang/basic-go/webook/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.LoggerV1
	repo   repository.FollowRepository
}

func NewFollowConsumer(client sarama.Client, l logger.LoggerV1, repo repository.FollowRepository) Consumer {
	return &MySQLBinlogConsumer{client: client, l: l, repo: repo}
}

func (r *MySQLBinlogConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("pub_follow_cache",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"webook_binlog"},
			saramax.NewHandler[canalx.Message[dao.FollowRelation]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[dao.FollowRelation]) error {
	if val.Table != "follow_relations" {
		// 我不关心的
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, data := range val.Data {
		var err error
		switch data.Status {
		case dao.FollowRelationStatusActive:
			err = r.repo.Cache().Follow(ctx, data.Follower, data.Followee)
		case dao.FollowRelationStatusInactive:
			err = r.repo.Cache().CancelFollow(ctx, data.Follower, data.Followee)
		}
		if err != nil {
			// 你可以继续，也可以中断
			return err
		}
	}
	return nil
}
