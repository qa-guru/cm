package cmd

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/aerokube/cm/selenoid"
	"github.com/spf13/cobra"
)

var (
	operatingSystem string
	arch            string
	version         string
	browsers        string
	useDrivers      bool
	browsersJson    string
	driversInfoUrl  string
	configDir       string
	uiConfigDir     string
	skipDownload    bool
	force           bool
	graceful        bool
	gracefulTimeout time.Duration
	args            string
	env             string
	browserEnv      string
	port            uint16
	uiPort          uint16
	userNS          string
	disableLogs     bool
	selenoidBinary  string
	selenoidUIBinary string
)

func init() {
	initFlags()

	selenoidCmd.AddCommand(selenoidDownloadCmd)
	selenoidCmd.AddCommand(selenoidArgsCmd)
	selenoidCmd.AddCommand(selenoidConfigureCmd)
	selenoidCmd.AddCommand(selenoidStartCmd)
	selenoidCmd.AddCommand(selenoidStopCmd)
	selenoidCmd.AddCommand(selenoidUpdateCmd)
	selenoidCmd.AddCommand(selenoidCleanupCmd)
	selenoidCmd.AddCommand(selenoidStatusCmd)

	selenoidUICmd.AddCommand(selenoidDownloadUICmd)
	selenoidUICmd.AddCommand(selenoidUIArgsCmd)
	selenoidUICmd.AddCommand(selenoidStartUICmd)
	selenoidUICmd.AddCommand(selenoidStopUICmd)
	selenoidUICmd.AddCommand(selenoidUpdateUICmd)
	selenoidUICmd.AddCommand(selenoidCleanupUICmd)
	selenoidUICmd.AddCommand(selenoidUIStatusCmd)
}

func initFlags() {
	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidArgsCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidStopCmd,
		selenoidUpdateCmd,
		selenoidCleanupCmd,
		selenoidStatusCmd,
		selenoidDownloadUICmd,
		selenoidUIArgsCmd,
		selenoidStartUICmd,
		selenoidStopUICmd,
		selenoidUpdateUICmd,
		selenoidCleanupUICmd,
		selenoidUIStatusCmd,
	} {
		c.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress output")
		c.Flags().BoolVarP(&useDrivers, "use-drivers", "d", false, "use drivers mode instead of Docker")
	}
	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidArgsCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidStopCmd,
		selenoidUpdateCmd,
		selenoidCleanupCmd,
		selenoidStatusCmd,
	} {
		c.Flags().StringVarP(&configDir, "config-dir", "c", selenoid.GetSelenoidConfigDir(), "directory to save files")
		c.Flags().Uint16VarP(&port, "port", "p", selenoid.DefaultPort, "override listen port")
	}
	for _, c := range []*cobra.Command{
		selenoidDownloadUICmd,
		selenoidUIArgsCmd,
		selenoidStartUICmd,
		selenoidStopUICmd,
		selenoidUpdateUICmd,
		selenoidCleanupUICmd,
		selenoidUIStatusCmd,
	} {
		c.Flags().StringVarP(&uiConfigDir, "config-dir", "c", selenoid.GetSelenoidConfigDir(), "directory to save files")
		c.Flags().Uint16VarP(&uiPort, "port", "p", selenoid.UIDefaultPort, "override listen port")
	}

	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidArgsCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidUpdateCmd,
		selenoidDownloadUICmd,
		selenoidUIArgsCmd,
		selenoidStartUICmd,
		selenoidUpdateUICmd,
	} {
		c.Flags().StringVarP(&operatingSystem, "operating-system", "o", runtime.GOOS, "target operating system (drivers only)")
		c.Flags().StringVarP(&arch, "architecture", "a", runtime.GOARCH, "target architecture (drivers only)")
	}
	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidArgsCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidUpdateCmd,
		selenoidDownloadUICmd,
		selenoidUIArgsCmd,
		selenoidStartUICmd,
		selenoidUpdateUICmd,
	} {
		c.Flags().StringVarP(&version, "version", "v", selenoid.Latest, "desired version; default is latest release")
		c.Flags().StringVarP(&registry, "registry", "r", selenoid.DefaultRegistryUrl, "Docker registry to use")
	}
	for _, c := range []*cobra.Command{
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidUpdateCmd,
	} {
		c.Flags().StringVarP(&browsersJson, "browsers-json", "j", "", "browsers JSON file to sync with")
		c.Flags().BoolVarP(&skipDownload, "no-download", "n", false, "only output config file without downloading images or drivers")
		c.Flags().StringVar(&selenoidBinary, "selenoid-binary", "", "path to Selenoid binary (default: download from qa-guru/selenoid release)")
		c.Flags().StringVar(&selenoidUIBinary, "selenoid-ui-binary", "", "path to Selenoid UI binary (default: download from qa-guru/selenoid-ui release)")
	}
	for _, c := range []*cobra.Command{
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidUpdateCmd,
	} {
		c.Flags().StringVarP(&browsers, "browsers", "b", "", "semicolon separated list of browser names (drivers mode only)")
		c.Flags().StringVarP(&browserEnv, "browser-env", "w", "", "override driver environment variables (drivers mode only, e.g. \"KEY1=value1 KEY2=value2\")")
		c.Flags().StringVarP(&driversInfoUrl, "drivers-info", "", "", "drivers catalog JSON URL (required with --use-drivers)")
	}
	for _, c := range []*cobra.Command{
		selenoidDownloadCmd,
		selenoidArgsCmd,
		selenoidConfigureCmd,
		selenoidStartCmd,
		selenoidDownloadUICmd,
		selenoidUIArgsCmd,
		selenoidStartUICmd,
	} {
		c.Flags().BoolVarP(&force, "force", "f", false, "force action")
	}
	for _, c := range []*cobra.Command{
		selenoidStopCmd,
		selenoidStopUICmd,
	} {
		c.Flags().BoolVarP(&graceful, "graceful", "", false, "do action gracefully (e.g. gracefully stop Selenoid)")
		c.Flags().DurationVarP(&gracefulTimeout, "graceful-timeout", "", 30*time.Second, "graceful timeout value (how much time to wait for graceful action execution)")
	}
	for _, c := range []*cobra.Command{
		selenoidStartCmd,
		selenoidUpdateCmd,
		selenoidStartUICmd,
		selenoidUpdateUICmd,
	} {
		c.Flags().StringVarP(&args, "args", "g", "", "additional service arguments (e.g. \"-limit 5\")")
		c.Flags().StringVarP(&env, "env", "e", "", "override service environment variables (e.g. \"KEY1=value1 KEY2=value2\")")
		c.Flags().StringVarP(&userNS, "userns", "", "", "override user namespace, similarly to \"docker run --userns host ...\" (Docker only)")
		c.Flags().BoolVarP(&disableLogs, "disable-logs", "", false, "start with log saving feature disabled")
	}
}

func createLifecycle(configDir string, port uint16) (*selenoid.Lifecycle, error) {
	config := selenoid.LifecycleConfig{
		Quiet:           quiet,
		Force:           force,
		Graceful:        graceful,
		GracefulTimeout: gracefulTimeout,
		ConfigDir:       configDir,
		UseDrivers:      useDrivers,
		Browsers:        browsers,
		BrowserEnv:      browserEnv,
		Download:        !skipDownload,
		Args:            args,
		Env:             env,
		Port:            int(port),
		DisableLogs:     disableLogs,

		RegistryUrl:  registry,
		BrowsersJson: browsersJson,
		UserNS:       userNS,

		SelenoidBinary:   selenoidBinary,
		SelenoidUIBinary: selenoidUIBinary,

		DriversInfoUrl: driversInfoUrl,
		OS:             operatingSystem,
		Arch:           arch,
		Version:        version,
	}
	return selenoid.NewLifecycle(&config)
}

var selenoidCmd = &cobra.Command{
	Use:   "selenoid",
	Short: "Download, configure and run Selenoid",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Usage()
	},
}

func stderr(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
}
