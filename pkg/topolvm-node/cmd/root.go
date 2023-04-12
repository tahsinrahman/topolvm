package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/topolvm/topolvm"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var config struct {
	csiSocket   string
	lvmdSocket  string
	metricsAddr string
	zapOpts     zap.Options
}

var rootCmd = &cobra.Command{
	Use:     "topolvm-node",
	Version: topolvm.Version,
	Short:   "TopoLVM CSI node",
	Long: `topolvm-node provides CSI node service.
It also works as a custom Kubernetes controller.

The node name where this program runs must be given by either
NODE_NAME environment variable or --nodename flag.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return subMain()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	fs := rootCmd.Flags()
	fs.StringVar(&config.csiSocket, "csi-socket", topolvm.DefaultCSISocket, "UNIX domain socket filename for CSI")
	fs.StringVar(&config.lvmdSocket, "lvmd-socket", topolvm.DefaultLVMdSocket, "UNIX domain socket of lvmd service")
	fs.StringVar(&config.metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	fs.String("nodename", "", "The resource name of the running node")
	fs.String("affinity-key", topolvm.GetTopologyNodeKey(), "The key of the node label to be used for topology-aware scheduling")

	viper.BindEnv("nodename", "NODE_NAME")
	viper.BindPFlag("nodename", fs.Lookup("nodename"))
	viper.BindPFlag("affinity-key", fs.Lookup("affinity-key"))

	goflags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(goflags)
	config.zapOpts.BindFlags(goflags)

	fs.AddGoFlagSet(goflags)
}
