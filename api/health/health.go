package health

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/Kaibling/chino/models"
	"github.com/Kaibling/chino/pkg/utils"
)

func health(w http.ResponseWriter, r *http.Request) {
	envelope := utils.GetContext("envelope", r).(*utils.Envelope)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// https://github.com/nelkinda/health-go
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	hc := models.MemoryCheck{
		Alloc:      formatB(m.Alloc),
		TotalAlloc: formatB(m.TotalAlloc),
		Sys:        formatB(m.Sys),
		NumGC:      fmt.Sprintf("%v", m.NumGC),
	}

	utils.Render(w, r, envelope.SetResponse(hc))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func formatB(b uint64) string {
	if b < 1024 {
		return fmt.Sprintf("%d B", b)
	}
	if b/1024 < 1024 {
		return fmt.Sprintf("%d KiB", b/1024)
	}
	if b/1024/1024 < 1024 {
		return fmt.Sprintf("%d MiB", b/1024/1024)
	}

	return fmt.Sprintf("%d GiB", b/1024/1024/1024)
}
