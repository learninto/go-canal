package job

import (
	"context"
	"fmt"
	"github.com/learninto/go-canal/pkg/oslib"
	"github.com/learninto/goutil"
	"github.com/learninto/goutil/conf"
	"github.com/learninto/goutil/log"
	"github.com/spf13/cobra"
	httpD "net/http"
)

var port int

// Cmd run job once or periodically
var Cmd = &cobra.Command{
	Use:   "job",
	Short: "Run job",
	Long:  `Go Run`,
	Run: func(cmd *cobra.Command, args []string) {

		server := &httpD.Server{Addr: fmt.Sprintf(":%d", port)} // 不指定 handler 则会使用默认 handler
		go func() {
			goutil.PrometheusHandleFunc("/metrics")
			goutil.Ping("/monitor/ping")

			if err := server.ListenAndServe(); err != nil {
				panic(err)
			}
		}()

		go func() {
			conf.OnConfigChange(func() {
				log.Reset()
				oslib.RestartApp()
			})
			conf.WatchConfig()
		}()

		ctx := context.Background()
		// 开始运行
		_ = run(ctx)
	},
}

func init() {
	Cmd.Flags().IntVar(&port, "port", 8080, "listen port")
	//Cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
}
