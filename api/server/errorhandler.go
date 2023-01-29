package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rumpl/bof/api/server/httpstatus"
	"github.com/rumpl/bof/api/server/httputils"
	"github.com/rumpl/bof/api/types"
	"github.com/rumpl/bof/api/types/versions"
	"google.golang.org/grpc/status"
)

// makeErrorHandler makes an HTTP handler that decodes a Docker error and
// returns it in the response.
func makeErrorHandler(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusCode := httpstatus.FromError(err)
		vars := mux.Vars(r)
		if apiVersionSupportsJSONErrors(vars["version"]) {
			response := &types.ErrorResponse{
				Message: err.Error(),
			}
			_ = httputils.WriteJSON(w, statusCode, response)
		} else {
			http.Error(w, status.Convert(err).Message(), statusCode)
		}
	}
}

func apiVersionSupportsJSONErrors(version string) bool {
	const firstAPIVersionWithJSONErrors = "1.23"
	return version == "" || versions.GreaterThan(version, firstAPIVersionWithJSONErrors)
}
