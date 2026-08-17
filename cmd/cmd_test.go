package cmd

import (
	"bytes"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestVersionCommandPrintsBuildInfo(t *testing.T) {
	buf := new(bytes.Buffer)
	versionCmd.SetOut(buf)
	versionCmd.SetErr(buf)
	versionCmd.Run(versionCmd, nil)

	out := buf.String()
	assert.Contains(t, out, "Git Revision:")
	assert.Contains(t, out, gitRevision)
	assert.Contains(t, out, "UTC Build Time:")
	assert.Contains(t, out, buildStamp)
}

func TestRootRegistersCoreSubcommands(t *testing.T) {
	names := make([]string, 0, len(rootCmd.Commands()))
	for _, c := range rootCmd.Commands() {
		names = append(names, c.Name())
	}
	assert.Contains(t, names, "version")
	assert.Contains(t, names, "selenoid")
	assert.True(t, strings.HasPrefix(rootCmd.Short, "cm is a configuration manager"))
}

func TestRootWithoutArgsShowsUsage(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Usage:")
}

func TestSelenoidStartWarmPoolFlags(t *testing.T) {
	t.Cleanup(func() {
		warmPool = false
		hotPool = false
		_ = selenoidStartCmd.Flags().Set("warm-pool", "false")
		_ = selenoidStartCmd.Flags().Set("hot-pool", "false")
	})

	assert.NotNil(t, selenoidStartCmd.Flags().Lookup("warm-pool"))
	assert.NotNil(t, selenoidStartCmd.Flags().Lookup("hot-pool"))
	assert.NotNil(t, selenoidUpdateCmd.Flags().Lookup("warm-pool"))
	assert.NotNil(t, selenoidStopCmd.Flags().Lookup("hot-pool"))

	warmPool, hotPool = false, false
	assert.NoError(t, selenoidStartCmd.ParseFlags([]string{"--warm-pool"}))
	assert.True(t, warmPool)
	assert.False(t, hotPool)

	warmPool, hotPool = false, false
	assert.NoError(t, selenoidStartCmd.ParseFlags([]string{"--hot-pool"}))
	assert.True(t, hotPool)
}

func TestSelenoidStartHelpMentionsWarmPool(t *testing.T) {
	usages := selenoidStartCmd.Flags().FlagUsages()
	assert.Contains(t, usages, "--warm-pool")
	assert.Contains(t, usages, "--hot-pool")
	assert.Contains(t, usages, "-warm-pool-url")
}
