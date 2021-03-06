package main

import (
	"html"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
}

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	assert.Error(checkArgs(event))
	plugin.APIURL = "NotURL"
	assert.Error(checkArgs(event))
	plugin.APIURL = "http://127.0.0.1:8080"
	assert.Error(checkArgs(event))
	plugin.APIKey = "01a23bc4-56d7-890e-fa12-3456789bcd01"
	assert.NoError(checkArgs(event))
}

func TestExecuteHandler(t *testing.T) {
	testcases := []struct {
		checkName   string
		expectError bool
		httpStatus int
	}{
		{"check1", true, 200},
		{"kubernetes-delete-entity", false, 200},
		{"kubernetes-delete-entity", false, 404},
	}
	for _, tc := range testcases {
		assert := assert.New(t)
		event := corev2.FixtureEvent("entity1", tc.checkName)
		plugin.APIKey = "blah-blah-blah"

		var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(html.EscapeString(r.URL.Path), "/api/core/v2/namespaces/default/entities/entity1")
			_, err := ioutil.ReadAll(r.Body)
			assert.NoError(err)
			w.WriteHeader(tc.httpStatus)
		}))
		_, err := url.ParseRequestURI(test.URL)
		require.NoError(t, err)
		plugin.APIURL = test.URL
		if tc.expectError {
			assert.Error(executeHandler(event))
		} else {
			assert.NoError(executeHandler(event))
		}
	}
}

func TestReadCAFile(t *testing.T) {
	// A self-signed CA PEM for testing
	caPEM := `-----BEGIN CERTIFICATE-----
MIIDpDCCAoygAwIBAgIUP8/zjYLblTJxw69FC1HXj3m46qQwDQYJKoZIhvcNAQEL
BQAwajELMAkGA1UEBhMCVVMxDzANBgNVBAgTBk9yZWdvbjERMA8GA1UEBxMIUG9y
dGxhbmQxEDAOBgNVBAoTB1Rlc3RpbmcxEzARBgNVBAsTClRlc3RpbmcgQ0ExEDAO
BgNVBAMTB1Rlc3QgQ0EwHhcNMjAwODI2MTYzNjAwWhcNMjUwODI1MTYzNjAwWjBq
MQswCQYDVQQGEwJVUzEPMA0GA1UECBMGT3JlZ29uMREwDwYDVQQHEwhQb3J0bGFu
ZDEQMA4GA1UEChMHVGVzdGluZzETMBEGA1UECxMKVGVzdGluZyBDQTEQMA4GA1UE
AxMHVGVzdCBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANAw+OXY
FImakh73dNCMn7FN9q9TMQBXenmcr21FOgnxHG60/fMkqgap/gxxVB0V7gUE0ZMX
+KHufx9zMnx6FhGaAbjNRRV71b8C/x+PnR71Od0YmTa5HmWRI81MS2AptZQRUwHx
c+AXdPf1f04QhznTwfVcAd8Iu1z0h0D3eQdX1fBrcru4LqpUAniNrD1AmcyEGVhD
xJGyYR25gFQWgRzH3gxzu3DaZ+mz4NsdmVOZLwIzZLo0mlgutFYTn62F+dv13nzk
X94vJ+5dGJtUo4MUIKSw6EkqQBNdKBVmh9lDIREY3eoE03vskXsJ/Ta1NSN0mDNl
GEUyh3YxcelRFQMCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQF
MAMBAf8wHQYDVR0OBBYEFGax+Emv0DC/WhKK2970Aoy/kuoyMA0GCSqGSIb3DQEB
CwUAA4IBAQAHAMK7ObpufJ+JIrKv261lrY3MaWiyjlcdilxkw5o1YTkEreZ+N3xd
Dm68lcvV7CwhG1pqmOLiztS0K/qtQ91c9JO5g2hwWE+Kc1kT/TpfD29KfpPqRrXB
x8p0X4r0gtIQCp6HqgO58HyfILUcsefRdkipf2MB51rNFKEF7FJ8t2UOj4NjVHCw
igAFYCoJig215prbCNSmGJml9eIRZcZ1hXYVgkmNU1LBaXk/JK4r2rOC+uZ/X8uu
vxfC5nxN0thcjRCpydrQMf/aLRynWxL05iV5+ZEqR8gcF2M+552SBA3QtW3xtXXF
K1FESSsDNHsGdZioIdZIKY8d0GTM4tEj
-----END CERTIFICATE-----`


	// A certificate, but not a CA
	notcaPEM := `-----BEGIN CERTIFICATE-----
MIIDwjCCAqqgAwIBAgIUIfeqXZpw70ZkdBNHSdgnYD7OOpgwDQYJKoZIhvcNAQEL
BQAwajELMAkGA1UEBhMCVVMxDzANBgNVBAgTBk9yZWdvbjERMA8GA1UEBxMIUG9y
dGxhbmQxEDAOBgNVBAoTB1Rlc3RpbmcxEzARBgNVBAsTClRlc3RpbmcgQ0ExEDAO
BgNVBAMTB1Rlc3QgQ0EwHhcNMjAwODI2MTcwODAwWhcNMjUwODI1MTcwODAwWjAe
MRwwGgYDVQQDExNiYWNrZW5kLmV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEA9nLGR4zo6FXWuyHJuVsQLdMxXbqeIo7X1rl0WxHhwYuS
nihXxaaJRTdqtLezUFcCJR8IYBZj+jvFI9xn6j04BvfYUlbk9ZqxC7S11qqQCFxG
pC0VNgXlnT7Ty/r3tR9kj2Z8xK+lf+ZxUR5X1NV0Oj6bbLXPgM+UlqMPubFebUUa
o+Benh6tZ1ubRTT3AI5O+1HJV/6WyslbLg9g0ju5nAwEW2tBI+XmqA7EyiEKJwwY
PlEDoiUAe1+JQXwQBe+ibQO6rP+yRzYcjlep87mllbinJRjapcDeQCA1VpbodOv3
wJwxzcIJeEKTcLI5VFjeW3MBBqFmfSs2zRDqpTlB5wIDAQABo4GrMIGoMA4GA1Ud
DwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwDAYDVR0T
AQH/BAIwADAdBgNVHQ4EFgQUvoITPNqM55EfnxAMRI9X7tdTEmIwHwYDVR0jBBgw
FoAUZrH4Sa/QML9aEorb3vQCjL+S6jIwKQYDVR0RBCIwIIIJbG9jYWxob3N0ggdi
YWNrZW5khwR/AAABhwQKAAABMA0GCSqGSIb3DQEBCwUAA4IBAQBgyKhgEmstkxvJ
5M5cfBBEB+YrZvTMjc5wch3PfI0puRUCOHeqRey8vb3x64zW18Xuo+GjLc29I0Fh
13BooBfI0chMOTMqyf7KWf0tBn2peCPh2BQikt2vinE2z7mDf6tDd7ZmC1X6HCBe
lrIY0/fRxYZlmv3Czllt846n106iVsyLDXlLlcicsYFotTQ3ZY6aNzekj8yBl3S6
3FNrZwd+s0tmiIioHT1kWo73F7IFflVPEnNMof+QUxC9TOLe6hkV8zTe2MXfH1T1
Pm7W5dUgDUMXI7BN2+GW2DJrLIANQm5EDz76AhqKiiqnz7JB4A0vKqJ94BSlIVwG
Aecfiikf
-----END CERTIFICATE-----`

	// Obviously...
	badcaPEM := "This is definitely not a certificate"

	assert := assert.New(t)
	r := strings.NewReader(caPEM)
	cert, err := readCAFile(r)
	assert.NoError(err)
	assert.Equal(true, cert.IsCA)
	assert.Equal("Test CA", cert.Subject.CommonName)

	r = strings.NewReader(notcaPEM)
	_, err = readCAFile(r)
	assert.Error(err)

	r = strings.NewReader(badcaPEM)
	_, err = readCAFile(r)
	assert.Error(err)
}
