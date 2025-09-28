package initializer

import (
	"github.com/nuanxinqing123/QLToolsV2/internal/middleware"
)

// StartRateLimitCleanup 启动限速器清理任务
func StartRateLimitCleanup() {
	middleware.StartCleanupTask()
}
