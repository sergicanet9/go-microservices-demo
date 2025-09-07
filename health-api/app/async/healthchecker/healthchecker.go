package healthchecker

import (
	"context"
	"net/http"
	"time"

	"github.com/sergicanet9/scv-go-tools/v4/observability"
)

const contentType = "application/json"

func RunHTTP(ctx context.Context, cancel context.CancelFunc, url string, interval time.Duration) {
	defer cancel()
	defer func() {
		if rec := recover(); rec != nil {
			observability.Logger().Printf("FATAL - recovered panic in HTTP healthchecker process: %v", rec)
		}
	}()

	for ctx.Err() == nil {
		<-time.After(interval)

		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			observability.Logger().Printf("HTTP healthchecker process - error: %s", err)
			continue
		}
		req.Header.Set("Content-Type", contentType)

		start := time.Now()
		resp, err := http.DefaultClient.Do(req)
		elapsed := time.Since(start)

		if err != nil {
			observability.Logger().Printf("HTTP healthchecker process - error: %s", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			observability.Logger().Printf("HTTP healthchecker process - error: %s", err)
			continue
		}

		observability.Logger().Printf("HTTP healthchecker process - health Check complete, time elapsed: %s", elapsed)
	}
}
