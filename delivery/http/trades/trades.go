package trades

import (
	constant "VincentLimarus/stock-analyzer-performance/constant/tradeConst"
	service "VincentLimarus/stock-analyzer-performance/service"
	"net/http"

	responseModel "VincentLimarus/stock-analyzer-performance/model/dto/http/response"
	dto "VincentLimarus/stock-analyzer-performance/model/dto/http/trades"

	"github.com/gin-gonic/gin"
)

type Trade struct {
	serviceRegistry service.IRegistry
}

func NewTrade(
	serviceRegistry service.IRegistry,
) *Trade {
	return &Trade{
		serviceRegistry: serviceRegistry,
	}
}

type ITrade interface {
	UpsertTradeData(c *gin.Context)
}

func (r *Trade) UpsertTradeData(c *gin.Context) {
	var payload []dto.TradeData
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusUnprocessableEntity, responseModel.Response{
			Message: http.StatusText(http.StatusUnprocessableEntity),
			Code:    constant.CodeInvalidRequest,
		})
		return
	}

	err := r.serviceRegistry.GetTrade().UpsertTradeData(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responseModel.Response{
			Message: http.StatusText(http.StatusInternalServerError),
			Code:    constant.CodeInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, responseModel.Response{
		Message: http.StatusText(http.StatusOK),
		Code:    constant.CodeSuccess,
	})
}
