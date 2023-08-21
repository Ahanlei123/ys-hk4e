package app

import (
	"context"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"hk4e/common/config"
	"hk4e/common/rpc"
	"hk4e/dispatch/controller"
	"hk4e/dispatch/dao"
	"hk4e/pkg/logger"
)

func Run(ctx context.Context, configFile string) error {
	config.InitConfig(configFile)

	logger.InitLogger("dispatch")
	logger.Warn("dispatch start")
	defer func() {
		logger.CloseLogger()
	}()

	db := dao.NewDao()
	defer db.CloseDao()

	// natsrpc client
	discoveryClient, err := rpc.NewDiscoveryClient()
	if err != nil {
		return err
	}

	_ = controller.NewController(db, discoveryClient)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-ctx.Done():
			return nil
		case s := <-c:
			logger.Warn("get a signal %s", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				logger.Warn("dispatch exit")
				return nil
			case syscall.SIGHUP:
			default:
				return nil
			}
		}
	}
}
