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
