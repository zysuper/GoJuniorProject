package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

const UserIndexName = "user_index"

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

type UserElasticDAO struct {
	client *elastic.Client
}

func (h *UserElasticDAO) Search(ctx context.Context, keywords []string) ([]User, error) {
	queryString := strings.Join(keywords, " ")
	query := elastic.NewMatchQuery("nickname", queryString)
	resp, err := h.client.Search(UserIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]User, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var u User
		err = json.Unmarshal(hit.Source, &u)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return nil, err
}

func (h *UserElasticDAO) InputUser(ctx context.Context, user User) error {
	_, err := h.client.Index().Index(UserIndexName).
		Id(strconv.FormatInt(user.Id, 10)).
		BodyJson(user).Do(ctx)
	return err
}

func NewUserElasticDAO(client *elastic.Client) UserDAO {
	return &UserElasticDAO{
		client: client,
	}
}
