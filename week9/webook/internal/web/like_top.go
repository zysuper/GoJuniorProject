package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LikeTopHandler struct {
	service service.TopNService
}

func NewLikeTopHandler(service service.TopNService) *LikeTopHandler {
	return &LikeTopHandler{service: service}
}

func (h *LikeTopHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("/likeTop100", h.GetTopLikeN)
}

func (h *LikeTopHandler) GetTopLikeN(cxt *gin.Context) {
	l, err := h.service.GetTopLikeN(cxt, "article", 100)
	if err != nil {
		cxt.JSON(http.StatusOK, Result{
			Code: 5, Msg: "系统错误",
		})
	}
	cxt.JSON(http.StatusOK, Result{Data: l, Msg: "ok", Code: 0})
}
