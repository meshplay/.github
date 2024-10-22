package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/khulnasoft/meshplay/meshplayctl/internal/cli/root/config"
	"github.com/khulnasoft/meshplay/meshplayctl/pkg/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TestHelper struct {
	Version string
	BaseURL string
}

type MockURL struct {
	// method such as GET or POST
	Method string
	// url to mock the request
	URL string
	// response for the request
	Response string
	// response code
	ResponseCode int
}

func NewTestHelper(_ *testing.T) *TestHelper {
	return &TestHelper{
		Version: "v0.5.10",
		BaseURL: MeshplayEndpoint,
	}
}

type CmdTestInput struct {
	Name                 string
	Args                 []string
	ExpectedResponse     string
	ExpectedResponseYaml string
	ExpectError          bool
	ErrorStringContains  []string
}

type GoldenFile struct {
	t    *testing.T
	name string
	dir  string
}

func NewGoldenFile(t *testing.T, name string, directory string) *GoldenFile {
	return &GoldenFile{t: t, name: name, dir: directory}
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// Path to the current file
func GetBasePath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}

	return filepath.Dir(filename)
}

// Load a Golden file
func (tf *GoldenFile) Load() string {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	content, err := os.ReadFile(path)
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}
	// ensuring that the newline characters in the content are consistent and match the expected newline representation
	normalizedContent := strings.ReplaceAll(string(content), "\r\n", "\n")
	return normalizedContent
}

// Load a Golden file
func (tf *GoldenFile) LoadByte() []byte {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	content, err := os.ReadFile(path)
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return content
}

// write a Golden file
func (tf *GoldenFile) Write(content string) {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.WriteFile(path, []byte(content), 0755)
			if err != nil {
				fmt.Printf("Unable to write file: %v", err)
			}
			return
		}
		tf.t.Fatal(err)
	}

	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

// write a Golden file
func (tf *GoldenFile) WriteInByte(content []byte) {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

// use default context /pkg/utils/TestConfig.yaml
func SetupContextEnv(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshplay directory")
	}
	viper.Reset()
	viper.SetConfigFile(path + "/../../../../pkg/utils/TestConfig.yaml")
	DefaultConfigPath = path + "/../../../../pkg/utils/TestConfig.yaml"
	//fmt.Println(viper.ConfigFileUsed())
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	_, err = config.GetMeshplayCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}

// setup logrus formatter and return the buffer in which commands output is to be set.
func SetupLogrusGrabTesting(_ *testing.T, _ bool) *bytes.Buffer {
	b := bytes.NewBufferString("")
	logrus.SetOutput(b)
	SetupLogrusFormatter()
	return b
}

// setup meshkit logger for testing and return the buffer in which commands output is to be set.
func SetupMeshkitLoggerTesting(_ *testing.T, verbose bool) *bytes.Buffer {
	b := bytes.NewBufferString("")
	Log = SetupMeshkitLogger("meshplayctl", verbose, b)
	return b
}

// setup custom context with SetupCustomContextEnv
func SetupCustomContextEnv(t *testing.T, pathToContext string) {
	viper.Reset()
	ViperCompose = viper.New()
	ViperMeshconfig = viper.New()

	viper.SetConfigFile(pathToContext)
	DefaultConfigPath = pathToContext
	//fmt.Println(viper.ConfigFileUsed())
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	_, err = config.GetMeshplayCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}

// Start mock HTTP client to mock requests
func StartMockery(t *testing.T) {
	// activate http mocking
	httpmock.Activate()

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiResponse := NewGoldenFile(t, "validate.version.github.golden", fixturesDir).Load()

	// For validate version requests
	url1 := "https://github.com/" + constants.GetMeshplayGitHubOrg() + "/" + constants.GetMeshplayGitHubRepo() + "/releases/tag/" + "v0.5.54"
	httpmock.RegisterResponder("GET", url1,
		httpmock.NewStringResponder(200, apiResponse))
}

// stop HTTP mock client
func StopMockery(_ *testing.T) {
	httpmock.DeactivateAndReset()
}

// Set file location for testing stuff
func SetFileLocationTesting(dir string) {
	MeshplayFolder = filepath.Join(dir, "fixtures", MeshplayFolder)
	DockerComposeFile = filepath.Join(MeshplayFolder, DockerComposeFile)
	AuthConfigFile = filepath.Join(MeshplayFolder, AuthConfigFile)
}

func Populate(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func StartMockMeshplayServer(t *testing.T) error {
	serverAddr := strings.TrimPrefix(MeshplayEndpoint, "http://")
	l, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				// Log error and break the loop if it's a permanent error
				if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
					t.Logf("Failed to accept connection: %v", err)
					break
				}
				continue
			}
			// Close the connection to verify IsServerRunning() in auth.go
			conn.Close()
		}
	}()

	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)
	return nil
}
