package cmd

import (
	"fmt"
	"mallbots/plugins/jwtc"
	"mallbots/plugins/pgxc"
	"mallbots/shared/common"
	"mallbots/shared/config"
	"os"
	"time"

	sctx "github.com/phathdt/service-context"

	"github.com/spf13/cobra"
)

const (
	serviceName = "api"
)

func newServiceCtx() sctx.ServiceContext {
	return sctx.NewServiceContext(
		sctx.WithName(serviceName),
		sctx.WithComponent(pgxc.New(common.KeyPgx, "")),
		sctx.WithComponent(jwtc.New(common.KeyJwt)),
	)
}

var rootCmd = &cobra.Command{
	Use:   serviceName,
	Short: fmt.Sprintf("start %s", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		sc := newServiceCtx()

		logger := sctx.GlobalLogger().GetLogger("service")

		time.Sleep(time.Second * 1)
		cfg, err := config.LoadConfig("")
		if err != nil {
			logger.Fatalf("Failed to load configuration: %v", err)
		}

		if err := sc.Load(); err != nil {
			logger.Fatal(err)
		}

		StartRouter(sc, cfg)
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
