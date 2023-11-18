package domain

type Msg struct {
	Id      string   `json:"id"`
	TplId   string   `json:"tplId"`
	Args    []string `json:"args"`
	Numbers []string `json:"numbers"`
}
