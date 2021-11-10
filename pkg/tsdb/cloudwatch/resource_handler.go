package cloudwatch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
)

func (e *cloudWatchExecutor) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/regions", resourceHandler(e.handleGetRegions))
	mux.HandleFunc("/namespaces", resourceHandler(e.handleGetNamespaces))
	mux.HandleFunc("/metrics", resourceHandler(e.handleGetMetrics))
	mux.HandleFunc("/dimensions", resourceHandler(e.handleGetDimensions))
	mux.HandleFunc("/dimension-values", resourceHandler(e.handleGetDimensionValues))
	mux.HandleFunc("/ebs-volume-ids", resourceHandler(e.handleGetEbsVolumeIds))
	mux.HandleFunc("/ec2-instance-attribute", resourceHandler(e.handleGetEc2InstanceAttribute))
	mux.HandleFunc("/resource-arns", resourceHandler(e.handleGetResourceArns))
}

type handleFn func(ctx context.Context, parameters url.Values,
	pluginCtx backend.PluginContext) ([]suggestData, error)

func resourceHandler(handleMetricFind handleFn) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		pluginContext := httpadapter.PluginConfigFromContext(ctx)
		err := req.ParseForm()
		if err != nil {
			writeResponse(rw, http.StatusBadRequest, fmt.Sprintf("unexpected error %v", err))
		}
		data, err := handleMetricFind(ctx, req.URL.Query(), pluginContext)
		if err != nil {
			writeResponse(rw, http.StatusBadRequest, fmt.Sprintf("unexpected error %v", err))
		}
		body, err := json.Marshal(data)
		if err != nil {
			writeResponse(rw, http.StatusBadRequest, fmt.Sprintf("unexpected error %v", err))
		}
		rw.WriteHeader(http.StatusOK)
		_, err = rw.Write(body)
		if err != nil {
			plog.Error("Unable to write HTTP response", "error", err)
		}
	}
}

func writeResponse(rw http.ResponseWriter, code int, msg string) {
	rw.WriteHeader(code)
	_, err := rw.Write([]byte(msg))
	if err != nil {
		plog.Error("Unable to write HTTP response", "error", err)
	}
}
