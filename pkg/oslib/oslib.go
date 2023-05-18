package oslib

import (
	"context"
	"github.com/siddontang/go-log/log"
	"os"
	"syscall"
)

// RestartApp 重启命令
func RestartApp() {
	executablePath, err := os.Executable()
	if err != nil {
		log.Error(context.Background(), "os.Executable 重启失败：", err)
		return
	}
	if err := syscall.Exec(executablePath, os.Args, os.Environ()); err != nil {
		log.Error(context.Background(), "syscall.Exec 重启失败：", err)
		return
	}
}

// GetConfPath 获取配置路径
func GetConfPath(ctx context.Context) (confPath string, err error) {
	// 获取conf路径
	if confPath = os.Getenv("CONF_PATH"); confPath != "" {
		return
	}

	log.Info(ctx, "env CONF_PATH is empty")
	if confPath, err = os.Getwd(); err != nil {
		log.Error(ctx, "os.Getwd Error：", err)
	}

	return
}
