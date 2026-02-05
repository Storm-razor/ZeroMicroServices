package rpcserver

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"
)

// ---------------------------
// @brief 自定义的一元拦截器,统一处理 RPC 服务端的错误日志,将业务自定义错误转换为 gRPC 标准错误
// ---------------------------
func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		caseErr := errors.Cause(err)
		if e, ok := caseErr.(*xerr.CodeError); ok {
			// 若是自定义的错误类型则日志打印
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)

			// 否则转为grpc错误
			err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
		} else {
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)
		}
	}

	return resp, err
}
