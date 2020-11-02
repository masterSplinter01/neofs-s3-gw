package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nspcc-dev/neofs-s3-gate/api"
	"github.com/nspcc-dev/neofs-s3-gate/api/layer"
	"go.uber.org/zap"
)

func (h *handler) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err error
		inf *layer.ObjectInfo

		req = mux.Vars(r)
		bkt = req["bucket"]
		obj = req["object"]
		rid = api.GetRequestID(r.Context())
	)

	if inf, err = h.obj.GetObjectInfo(r.Context(), bkt, obj); err != nil {
		h.log.Error("could not find object",
			zap.String("request_id", rid),
			zap.String("bucket_name", bkt),
			zap.String("object_name", obj),
			zap.Error(err))

		api.WriteErrorResponse(r.Context(), w, api.Error{
			Code:           api.GetAPIError(api.ErrInternalError).Code,
			Description:    err.Error(),
			HTTPStatusCode: http.StatusInternalServerError,
		}, r.URL)

		return
	}

	params := &layer.GetObjectParams{
		Bucket: inf.Bucket,
		Object: inf.Name,
		Writer: w,
	}

	// params.Length = inf.Size

	if err = h.obj.GetObject(r.Context(), params); err != nil {
		h.log.Error("could not get object",
			zap.String("request_id", rid),
			zap.String("bucket_name", bkt),
			zap.String("object_name", obj),
			zap.Error(err))

		api.WriteErrorResponse(r.Context(), w, api.Error{
			Code:           api.GetAPIError(api.ErrInternalError).Code,
			Description:    err.Error(),
			HTTPStatusCode: http.StatusInternalServerError,
		}, r.URL)

		return
	}

	w.Header().Set("Content-Type", inf.ContentType)
	w.Header().Set("Last-Modified", inf.Created.Format(http.TimeFormat))
	w.Header().Set("Content-Length", strconv.FormatInt(inf.Size, 10))

	w.WriteHeader(http.StatusOK)
}