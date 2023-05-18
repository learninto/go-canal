package init

import (
	"github.com/learninto/goutil/metrics"
	_ "go.uber.org/automaxprocs" // 根据容器配额设置 maxprocs
)

func init() {
	metrics.InitMetrics("canal")
}
