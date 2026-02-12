package jobtype

import "github.com/wwwzy/ZeroMicroServices/app/order/model"

// 异步任务执行的数据结构

// DeferCloseHomestayOrderPayload defer close homestay order
type DeferCloseHomestayOrderPayload struct {
	Sn string
}

// PaySuccessNotifyUserPayload pay success notify user
type PaySuccessNotifyUserPayload struct {
	Order *model.HomestayOrder
}
