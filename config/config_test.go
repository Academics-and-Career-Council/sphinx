package config

import (
	"bytes"
	"strings"
	"testing"
)

// test example config file is loaded correctly
func TestConfigurationFileLoading(t *testing.T) {

	config, err := New("../example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Error("could not load example configuration")
	}

	if config.Proxy.Handler != "http" {
		t.Error("expected http for Proxy.Handler")
	}

	if len(config.Limits) != 4 {
		t.Error("expected 4 bucket definitions")
	}

	for name, limit := range config.Limits {
		if limit.Interval < 1 {
			t.Errorf("limit interval should be greator than 1 for limit: %s", name)
		}
		if limit.Max < 1 {
			t.Errorf("limit max should be greator than 1 for limit: %s", name)
		}
	}
}

func TestInvalidConfigurationPath(t *testing.T) {
	if _, err := New("./does-not-exist.yaml"); err == nil {
		t.Fatalf("Expected error for invalid config path")
	} else if !strings.Contains(err.Error(), "no such file or directory") {
		t.Fatalf("Expected no file error got %s", err.Error())
	}
}

func TestInvalidYaml(t *testing.T) {
	invalidYaml := []byte(`
forward
  host$$: proxy.example.com
`)

	if _, err := LoadAndValidateYaml(invalidYaml); !strings.Contains(err.Error(), "YAML error:") {
		t.Errorf("expected yaml error, got %s", err.Error())
	}
}

// Incorrect configuration file should return errors
func TestInvalidProxyConfig(t *testing.T) {

	invalidConfig := []byte(`
proxy:
  host: http://proxy.example.com
`)
	_, err := LoadAndValidateYaml(invalidConfig)
	if err == nil || !strings.Contains(err.Error(), "handler") {
		t.Errorf("Expected proxy handler error. Got: %s", err.Error())
	}

	invalidConfig = []byte(`
proxy:
  handler: http
  host: proxy.example.com
`)
	_, err = LoadAndValidateYaml(invalidConfig)
	if err == nil || !strings.Contains(err.Error(), "host:port") {
		t.Errorf("Expected proxy host error. Got: %s", err.Error())
	}

	invalidConfig = []byte(`
proxy:
  handler: http
  host: proxy.example.com
  listen: :8000
`)
	_, err = LoadAndValidateYaml(invalidConfig)
	if err == nil || !strings.Contains(err.Error(), "proxy") {
		t.Errorf("Expected proxy host error. Got: %s", err.Error())
	}
}

func TestInvalidLimitConfig(t *testing.T) {

	baseBuf := bytes.NewBufferString(`
proxy:
  handler: http
  host: http://proxy.example.com
  listen: "0.0.0.0:8080"
storage:
  type: memory
`)

	configBuf := baseBuf
	configBuf.WriteString(`
limits:
  bearer/events:
    keys:
      headers: 
        - 'Authentication'
`)
	_, err := LoadAndValidateYaml(configBuf.Bytes())
	if err == nil || !strings.Contains(err.Error(), "interval") {
		t.Errorf("Expected Limit Interval error. Got: %s", err.Error())
	}

	configBuf = baseBuf
	configBuf.WriteString(`
limits:
  bearer/events:
    interval: 10
    keys:
      headers: 
        - 'Authentication'
`)
	_, err = LoadAndValidateYaml(configBuf.Bytes())
	if err == nil || !strings.Contains(err.Error(), "max") {
		t.Errorf("Expected Limit Interval error. Got: %s", err.Error())
	}
}

func TestInvalidStorageConfig(t *testing.T) {
	baseBuf := bytes.NewBufferString(`
proxy:
  handler: http
  host: http://proxy.example.com
  listen: localhost:8080
limits:
  test:
    keys:
      headers:
        - Authorization
    interval: 15  # in seconds
    max: 200
`)

	_, err := LoadAndValidateYaml(baseBuf.Bytes())
	if err == nil || !strings.Contains(err.Error(), "storage type must be set") {
		t.Errorf("Expected Storage error. Got: %s", err.Error())
	}

	configBuf := baseBuf
	configBuf.WriteString(`
storage:
  type: redis
`)
	_, err = LoadAndValidateYaml(configBuf.Bytes())
	if err == nil || !strings.Contains(err.Error(), "host") {
		t.Errorf("Expected redis Storage host error. Got: %s", err.Error())
	}

	configBuf = baseBuf
	configBuf.WriteString(`
storage:
  type: redis
  host: localhost
`)
	_, err = LoadAndValidateYaml(configBuf.Bytes())
	if err == nil || !strings.Contains(err.Error(), "port") {
		t.Errorf("Expected redis Storage host error. Got: %s", err.Error())
	}
}