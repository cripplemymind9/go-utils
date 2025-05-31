package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

// ErrorHandler преобразует gRPC ошибки в HTTP ответы с JSON
func ErrorHandler() runtime.ErrorHandlerFunc {
	return func(
		ctx context.Context,
		_ *runtime.ServeMux,
		_ runtime.Marshaler,
		w http.ResponseWriter,
		_ *http.Request,
		err error,
	) {
		st := status.Convert(err)
		code := runtime.HTTPStatusFromCode(st.Code())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)

		response := map[string]interface{}{
			"code":    st.Code().String(),
			"message": st.Message(),
		}

		if len(st.Details()) > 0 {
			response["details"] = st.Details()
		}

		json.NewEncoder(w).Encode(response)
	}
}
