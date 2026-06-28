package selenoid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aerokube/selenoid/config"
	"github.com/docker/docker/api/types/image"
	assert "github.com/stretchr/testify/require"
)

var (
	mockDockerServer *httptest.Server
	imageName        string
	containerName    string
	port             int
)

func init() {
	resetImageName()
	resetContainerName()
	resetPort()
	mockDockerServer = httptest.NewServer(mux())
	os.Setenv("DOCKER_HOST", "tcp://"+hostPort(mockDockerServer.URL))
}

func setImageName(name string) {
	imageName = name
}

func resetImageName() {
	setImageName("docker.io/" + selenoidWrapperImage + ":" + wrapperImageTag)
}

func setContainerName(name string) {
	containerName = name
}

func resetContainerName() {
	setContainerName(selenoidContainerName)
}

func setPort(p int) {
	port = p
}

func resetPort() {
	setPort(DefaultPort)
}

func mux() http.Handler {
	mux := http.NewServeMux()

	//Docker Registry API mock
	mux.HandleFunc("/v2/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	))

	mux.HandleFunc("/v2/qaguru/selenoid/tags/list", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			_, _ = fmt.Fprintln(w, `{"name":"selenoid", "tags": ["1.4.0", "1.4.1"]}`)
		},
	))

	mux.HandleFunc("/v2/qaguru/selenoid-ui/tags/list", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			_, _ = fmt.Fprintln(w, `{"name":"selenoid-ui", "tags": ["1.5.2"]}`)
		},
	))

	mux.HandleFunc("/v2/selenoid/firefox/tags/list", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			_, _ = fmt.Fprintln(w, `{"name":"firefox", "tags": ["46.0", "45.0", "7.0", "latest"]}`)
		},
	))

	mux.HandleFunc("/v2/selenoid/opera/tags/list", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			_, _ = fmt.Fprintln(w, `{"name":"opera", "tags": ["44.0", "latest"]}`)
		},
	))

	mux.HandleFunc("/v2/selenoid/android/tags/list", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			_, _ = fmt.Fprintln(w, `{"name":"android", "tags": ["10.0"]}`)
		},
	))

	mux.HandleFunc("/v2/browsers/edge/tags/list", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			_, _ = fmt.Fprintln(w, `{"name":"edge", "tags": ["88.0"]}`)
		},
	))

	//Docker API mock
	mux.HandleFunc("/v1.29/version", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			output := `
				{
				
					"Version": "17.04.0",
					"Os": "linux",
					"KernelVersion": "3.19.0-23-generic",
					"GoVersion": "go1.7.5",
					"GitCommit": "deadbee",
					"Arch": "amd64",
					"ApiVersion": "1.29",
					"MinAPIVersion": "1.12",
					"BuildTime": "2016-06-14T07:09:13.444803460+00:00",
					"Experimental": true
				
				}
			`
			_, _ = w.Write([]byte(output))
		},
	))
	mux.HandleFunc("/v1.29/images/create", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			output := `{"id": "a86cd3433934", "status": "Downloading layer"}`
			_, _ = w.Write([]byte(output))
		},
	))
	mux.HandleFunc("/v1.29/images/json", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			output := fmt.Sprintf(`
			[{
			
			    "Id": "sha256:e216a057b1cb1efc11f8a268f37ef62083e70b1b38323ba252e25ac88904a7e8",
			    "ParentId": "",
			    "RepoTags": [ "%s" ],
			    "RepoDigests": [],
			    "Created": 1474925151,
			    "Size": 103579269,
			    "VirtualSize": 103579269,
			    "SharedSize": 0,
			    "Labels": { },
			    "Containers": 2
			
			}]
			`, imageName)
			_, _ = w.Write([]byte(output))
		},
	))
	mux.HandleFunc("/v1.29/networks/selenoid", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`
              [{
                "Name": "selenoid",
                "Id": "39d591dabe313ed90b599e6d6515301e879c088b449a260cc02981bd25b52a6f",
                "Created": "2020-02-29T14:41:12.960257Z",
                "Scope": "local",
                "Driver": "bridge",
                "EnableIPv6": false,
                "IPAM": {
                  "Driver": "default",
                  "Options": {},
                  "Config": [
                    {
                      "Subnet": "172.18.0.0/16",
                      "Gateway": "172.18.0.1"
                    }
                  ]
                },
                "Internal": false,
                "Attachable": false,
                "Ingress": false,
                "ConfigFrom": {
                  "Network": ""
                },
                "ConfigOnly": false,
                "Containers": {},
                "Options": {},
                "Labels": {}
              }]`))
		},
	))
	mux.HandleFunc("/v1.29/networks/create", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			output := `{"id": "39d591dabe31", "warnings": []}`
			_, _ = w.Write([]byte(output))
		},
	))
	mux.HandleFunc("/v1.29/containers/create", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			output := `{"id": "e90e34656806", "warnings": []}`
			_, _ = w.Write([]byte(output))
		},
	))
	mux.HandleFunc("/v1.29/containers/e90e34656806/start", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	))
	mux.HandleFunc("/v1.29/containers/e90e34656806/stop", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	))
	mux.HandleFunc("/v1.29/containers/e90e34656806/logs", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Some logs...\n"))
		},
	))
	mux.HandleFunc("/v1.29/containers/e90e34656806", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	))
	mux.HandleFunc("/v1.29/containers/kill", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	))
	mux.HandleFunc("/v1.29/containers/json", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			output := fmt.Sprintf(`
			[{
				"Id": "e90e34656806",
				"Names": [ "%s" ],
				"Image": "%s:latest",
				"ImageID": "e216a057b1cb1efc11f8a268f37ef62083e70b1b38323ba252e25ac88904a7e8",
				"Command": "/usr/bin/some-cmd",
				"Created": 1367854154,
				"State": "Exited",
				"Status": "Exit 0",
				"Ports": [
					{
						"PrivatePort": %d,
						"PublicPort": %d,
						"Type": "tcp"
					}
				],
				"Labels": { },
				"SizeRw": 12288,
				"SizeRootFs": 0,
				"HostConfig": {},
				"NetworkSettings": {},
			    	"Mounts": [ ]
				
			}]
			`, containerName, imageName, port, port)
			_, _ = w.Write([]byte(output))
		},
	))
	return mux
}

func TestImageWithTag(t *testing.T) {
	assert.Equal(t, imageWithTag("selenoid/firefox", "tag"), "selenoid/firefox:tag")
}

func TestFetchImageTags(t *testing.T) {
	lcConfig := LifecycleConfig{
		RegistryUrl: mockDockerServer.URL,
		Quiet:       false,
	}
	c, err := NewDockerConfigurator(&lcConfig)
	assert.NoError(t, err)
	defer c.Close()
	tags := c.fetchImageTags("selenoid/firefox")
	assert.Len(t, tags, 3)
	assert.Equal(t, tags[0], "46.0")
	assert.Equal(t, tags[1], "45.0")
	assert.Equal(t, tags[2], "7.0")
}

func TestConfigureDocker(t *testing.T) {
	testConfigure(t, true)
}

func TestLimitNoPull(t *testing.T) {
	testConfigure(t, false)
}

func testConfigure(t *testing.T, download bool) {
	withTmpDir(t, "test-docker-configure", func(t *testing.T, dir string) {

		lcConfig := LifecycleConfig{
			ConfigDir:   dir,
			RegistryUrl: mockDockerServer.URL,
			Download:    download,
			Quiet:       false,
		}
		c, err := NewDockerConfigurator(&lcConfig)
		assert.NoError(t, err)
		defer c.Close()
		assert.False(t, c.IsConfigured())
		cfgPointer, err := (*c).Configure()
		assert.NoError(t, err)
		assert.NotNil(t, cfgPointer)

		cfg := *cfgPointer
		assert.Len(t, cfg, 8)
		assert.Contains(t, cfg, "chrome")
		assert.Contains(t, cfg, "firefox")
		assert.Contains(t, cfg, "msedge")
		assert.Contains(t, cfg, "playwright-chromium")
		assert.Contains(t, cfg, "playwright-firefox")
		assert.Contains(t, cfg, "playwright-webkit")
		assert.Contains(t, cfg, "playwright-chrome")
		assert.Contains(t, cfg, "playwright-msedge")
	})
}

func TestSyncWithConfig(t *testing.T) {
	withTmpDir(t, "test-sync-with-config", func(t *testing.T, dir string) {

		initialCfg := SelenoidConfig{
			"firefox": {
				Versions: map[string]*config.Browser{
					"46.0": {
						Image: "selenoid/vnc_firefox:46.0",
						Port:  "4444",
					},
				},
			},
		}

		initialCfgFile := filepath.Join(dir, "initial-browsers.json")
		data, _ := json.Marshal(initialCfg)
		_ = os.WriteFile(initialCfgFile, data, 0644)

		lcConfig := LifecycleConfig{
			ConfigDir:    dir,
			RegistryUrl:  mockDockerServer.URL,
			BrowsersJson: initialCfgFile,
			Download:     true,
			Quiet:        false,
		}
		c, err := NewDockerConfigurator(&lcConfig)
		assert.NoError(t, err)
		defer c.Close()
		assert.False(t, c.IsConfigured())
		cfgPointer, err := (*c).Configure()
		assert.NoError(t, err)
		assert.NotNil(t, cfgPointer)

		cfg := *cfgPointer
		assert.Equal(t, cfg, initialCfg)
	})

}

func TestStartStopContainer(t *testing.T) {
	withTmpDir(t, "test-start-stop", func(t *testing.T, dir string) {
		selenoidBin := filepath.Join(dir, binDirName, "selenoid")
		assert.NoError(t, os.MkdirAll(filepath.Dir(selenoidBin), 0755))
		assert.NoError(t, os.WriteFile(selenoidBin, []byte{0}, 0755))
		assert.NoError(t, os.WriteFile(getSelenoidConfigPath(dir), []byte("{}"), 0644))

		c, err := NewDockerConfigurator(&LifecycleConfig{
			ConfigDir:      dir,
			RegistryUrl:    mockDockerServer.URL,
			Port:           DefaultPort,
			SelenoidBinary: selenoidBin,
			UserNS:         "host",
		})
		assert.NoError(t, err)
		assert.True(t, c.IsRunning())
		assert.NoError(t, c.Start())
		c.Status()
		assert.NoError(t, c.Stop())
	})
}

func TestStartStopUIContainer(t *testing.T) {
	defer func() {
		resetImageName()
		resetContainerName()
		resetPort()
	}()
	withTmpDir(t, "test-start-stop-ui", func(t *testing.T, dir string) {
		uiBin := filepath.Join(dir, binDirName, "selenoid-ui")
		assert.NoError(t, os.MkdirAll(filepath.Dir(uiBin), 0755))
		assert.NoError(t, os.WriteFile(uiBin, []byte{0}, 0755))
		assert.NoError(t, os.WriteFile(getSelenoidConfigPath(dir), []byte("{}"), 0644))

		c, err := NewDockerConfigurator(&LifecycleConfig{
			ConfigDir:        dir,
			RegistryUrl:      mockDockerServer.URL,
			Port:             UIDefaultPort,
			SelenoidUIBinary: uiBin,
		})
		assert.NoError(t, err)
		setContainerName(selenoidUIContainerName)
		setImageName("docker.io/" + selenoidUIWrapperImage + ":" + wrapperImageTag)
		setPort(UIDefaultPort)
		assert.True(t, c.IsUIRunning())
		assert.NoError(t, c.StartUI())
		c.UIStatus()
		assert.NoError(t, c.StopUI())
	})
}

func TestDownload(t *testing.T) {
	withTmpDir(t, "test-download", func(t *testing.T, dir string) {
		selenoidBin := filepath.Join(dir, binDirName, "selenoid")
		assert.NoError(t, os.MkdirAll(filepath.Dir(selenoidBin), 0755))
		assert.NoError(t, os.WriteFile(selenoidBin, []byte{0}, 0755))

		c, err := NewDockerConfigurator(&LifecycleConfig{
			ConfigDir:      dir,
			RegistryUrl:    mockDockerServer.URL,
			Quiet:          true,
			SelenoidBinary: selenoidBin,
		})
		assert.NoError(t, err)
		assert.True(t, c.IsDownloaded())
		ref, err := c.Download()
		assert.NoError(t, err)
		assert.NotNil(t, ref)
		assert.NoError(t, c.PrintArgs())
	})
}

func TestDownloadUI(t *testing.T) {
	defer func() {
		resetImageName()
	}()
	withTmpDir(t, "test-download-ui", func(t *testing.T, dir string) {
		uiBin := filepath.Join(dir, binDirName, "selenoid-ui")
		assert.NoError(t, os.MkdirAll(filepath.Dir(uiBin), 0755))
		assert.NoError(t, os.WriteFile(uiBin, []byte{0}, 0755))

		c, err := NewDockerConfigurator(&LifecycleConfig{
			ConfigDir:        dir,
			RegistryUrl:      mockDockerServer.URL,
			Quiet:            true,
			SelenoidUIBinary: uiBin,
		})
		setImageName("docker.io/" + selenoidUIWrapperImage + ":" + wrapperImageTag)
		assert.NoError(t, err)
		assert.True(t, c.IsUIDownloaded())
		ref, err := c.DownloadUI()
		assert.NoError(t, err)
		assert.NotNil(t, ref)
		assert.NoError(t, c.PrintUIArgs())
	})
}

func TestGetSelenoidImage(t *testing.T) {
	defer func() {
		resetImageName()
	}()
	c, err := NewDockerConfigurator(&LifecycleConfig{
		RegistryUrl: mockDockerServer.URL,
		Quiet:       true,
		Version:     Latest,
	})
	assert.NoError(t, err)
	assert.NotNil(t, c.getSelenoidImage())
	setImageName("docker.io/" + selenoidUIWrapperImage + ":" + wrapperImageTag)
	assert.Nil(t, c.getSelenoidImage())
}

func TestFindMatchingImage(t *testing.T) {

	var (
		selenoid141 = image.Summary{
			ID:       "1",
			RepoTags: []string{"qaguru/selenoid:1.4.1"},
			Created:  100,
		}
		selenoid143 = image.Summary{
			ID:       "3",
			RepoTags: []string{"qaguru/selenoid:1.4.3"},
			Created:  300,
		}
		selenoid120CustomRegistry = image.Summary{
			ID:       "4",
			RepoTags: []string{"my-registry.com:443/qaguru/selenoid:1.2.0"},
			Created:  100,
		}
	)
	images := []image.Summary{
		selenoid141,
		{
			ID:       "2",
			RepoTags: []string{"qaguru/selenoid-ui:1.5.1"},
			Created:  200, //Intentionally using small timestamps
		},
		selenoid143,
		selenoid120CustomRegistry,
	}

	assert.Nil(t, findMatchingImage(images, "unknown-image-name", Latest))
	assert.Nil(t, findMatchingImage(images, "qaguru/selenoid", "missing-version"))

	foundSelenoid141 := findMatchingImage(images, "qaguru/selenoid", "1.4.1")
	assert.NotNil(t, foundSelenoid141)
	assert.Equal(t, *foundSelenoid141, selenoid141)

	foundSelenoidEmpty := findMatchingImage(images, "qaguru/selenoid", "")
	assert.NotNil(t, foundSelenoidEmpty)
	assert.Equal(t, *foundSelenoidEmpty, selenoid143)

	foundSelenoidLatest := findMatchingImage(images, "qaguru/selenoid", Latest)
	assert.NotNil(t, foundSelenoidLatest)
	assert.Equal(t, *foundSelenoidLatest, selenoid143)

	foundSelenoidCustomRegistry := findMatchingImage(images, "my-registry.com:443/qaguru/selenoid", "1.2.0")
	assert.NotNil(t, foundSelenoidCustomRegistry, nil)
	assert.Equal(t, *foundSelenoidCustomRegistry, selenoid120CustomRegistry)

	foundSelenoidWithoutRegistry := findMatchingImage(images, "qaguru/selenoid", "1.2.0")
	assert.NotNil(t, foundSelenoidWithoutRegistry, nil)
	assert.Equal(t, *foundSelenoidWithoutRegistry, selenoid120CustomRegistry)
}

func TestIsVideoRecordingSupported(t *testing.T) {
	logger := Logger{}
	assert.False(t, isVideoRecordingSupported(logger, "wrong-version"))
	assert.False(t, isVideoRecordingSupported(logger, "1.3.9"))
	assert.True(t, isVideoRecordingSupported(logger, "1.4.0"))
	assert.True(t, isVideoRecordingSupported(logger, "1.4.1"))
	assert.True(t, isVideoRecordingSupported(logger, "1.5.0"))
	assert.True(t, isVideoRecordingSupported(logger, "latest"))
}

func TestFilterOutLatest(t *testing.T) {
	tags := filterOutLatest([]string{"one", "latest", "latest-release", "two"})
	assert.Equal(t, tags, []string{"one", "two"})
}

func TestChooseVolumeConfigDir(t *testing.T) {
	dirWithoutVariable := chooseVolumeConfigDir("/some/dir", []string{"one", "two"})
	assert.Equal(t, dirWithoutVariable, "/some/dir")
	os.Setenv("OVERRIDE_HOME", "/test/dir")
	defer os.Unsetenv("OVERRIDE_HOME")
	dir := chooseVolumeConfigDir("/some/dir", []string{"one", "two"})
	assert.Equal(t, dir, "/test/dir/one/two")
}

func TestPostProcessPath(t *testing.T) {
	assert.Equal(t, postProcessPath("C:\\Users\\admin"), "/c/Users/admin")
	assert.Equal(t, postProcessPath("C:\\C:\\Users\\admin"), "/c/C:/Users/admin")
	assert.Equal(t, postProcessPath("1"), "1")
	assert.Empty(t, postProcessPath(""))
}

func TestValidEnviron(t *testing.T) {
	assert.Equal(t, validateEnviron([]string{"=::=::"}), []string{})
	assert.Equal(t, validateEnviron([]string{"HOMEDRIVE=C:", "DOCKER_HOST=192.168.0.1", "=::=::"}), []string{"HOMEDRIVE=C:", "DOCKER_HOST=192.168.0.1"})
}
