package result

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

// ---------------------------
// @brief http请求正常返回
// ---------------------------
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		// 成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		// 错误返回
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*xerr.CodeError); ok {
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok {
				//若是grpc的错误
				grpcCode := uint32(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) {
					//若是grpc自定义错误
					errcode = grpcCode
					errmsg = gstatus.Message()
				}
			}
		}

		logx.WithContext(r.Context()).Errorf("【API-ERR】 : %+v ", err)
		httpx.WriteJson(w, http.StatusBadRequest, Error(errcode, errmsg))
	}
}

// ---------------------------
// @brief http鉴权错误返回(用于鉴权中间件)
// ---------------------------
func AuthHttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil {
		// 成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		// 错误返回
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*xerr.CodeError); ok {
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok {
				//若是grpc的错误
				grpcCode := uint32(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) {
					//若是grpc自定义错误
					errcode = grpcCode
					errmsg = gstatus.Message()
				}
			}
		}

		logx.WithContext(r.Context()).Errorf("【GATEWAY-ERR】 : %+v ", err)
		httpx.WriteJson(w, http.StatusUnauthorized, Error(errcode, errmsg))
	}
}

// ---------------------------
// @brief http参数错误返回
// ---------------------------
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error) {
	errMsg := fmt.Sprintf("%s ,%s", xerr.MapErrMsg(xerr.REUQEST_PARAM_ERROR), err.Error())
	httpx.WriteJson(w, http.StatusBadRequest, Error(xerr.REUQEST_PARAM_ERROR, errMsg))
}
