package veneur

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	exampleConfig, err := os.Open("example.yaml")
	assert.NoError(t, err)
	defer exampleConfig.Close()

	c, err := readConfig(exampleConfig)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "https://app.datadoghq.com", c.APIHostname)
	assert.Equal(t, 96, c.NumWorkers)

	interval, err := c.ParseInterval()
	assert.NoError(t, err)
	assert.Equal(t, interval, 10*time.Second)

	assert.Equal(t, c.TraceAddress, "127.0.0.1:8128")

}

func TestReadBadConfig(t *testing.T) {
	const exampleConfig = `--- api_hostname: :bad`
	r := strings.NewReader(exampleConfig)
	c, err := readConfig(r)

	assert.NotNil(t, err, "Should have encountered parsing error when reading invalid config file")
	assert.Equal(t, c, Config{}, "Parsing invalid config file should return zero struct")
}

func TestHostname(t *testing.T) {
	const hostnameConfig = "hostname: foo"
	r := strings.NewReader(hostnameConfig)
	c, err := readConfig(r)
	assert.Nil(t, err, "Should parsed valid config file: %s", hostnameConfig)
	assert.Equal(t, c, Config{Hostname: "foo", ReadBufferSizeBytes: defaultBufferSizeBytes},
		"Should have parsed hostname into Config")

	const noHostname = "hostname: ''"
	r = strings.NewReader(noHostname)
	c, err = readConfig(r)
	assert.Nil(t, err, "Should parsed valid config file: %s", noHostname)
	currentHost, err := os.Hostname()
	assert.Nil(t, err, "Could not get current hostname")
	assert.Equal(t, c, Config{Hostname: currentHost, ReadBufferSizeBytes: defaultBufferSizeBytes},
		"Should have used current hostname in Config")

	const omitHostname = "omit_empty_hostname: true"
	r = strings.NewReader(omitHostname)
	c, err = readConfig(r)
	assert.Nil(t, err, "Should parsed valid config file: %s", omitHostname)
	assert.Equal(t, c, Config{
		Hostname:            "",
		ReadBufferSizeBytes: defaultBufferSizeBytes,
		OmitEmptyHostname:   true})
}
