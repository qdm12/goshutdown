package mock

import (
	"github.com/golang/mock/gomock"
	"github.com/qdm12/goshutdown/goroutine/mock_goroutine"
	"github.com/qdm12/goshutdown/group/mock_group"
	"github.com/qdm12/goshutdown/order/mock_order"
)

func NewGoRoutineMockHandler(ctrl *gomock.Controller) *mock_goroutine.MockHandler {
	return mock_goroutine.NewMockHandler(ctrl)
}

func NewGroupMockHandler(ctrl *gomock.Controller) *mock_group.MockHandler {
	return mock_group.NewMockHandler(ctrl)
}

func NewOrderMockHandler(ctrl *gomock.Controller) *mock_order.MockHandler {
	return mock_order.NewMockHandler(ctrl)
}
