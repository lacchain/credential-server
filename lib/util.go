package lib

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var clientAuthTypes = map[string]tls.ClientAuthType{
	"noclientcert":               tls.NoClientCert,
	"requestclientcert":          tls.RequestClientCert,
	"requireanyclientcert":       tls.RequireAnyClientCert,
	"verifyclientcertifgiven":    tls.VerifyClientCertIfGiven,
	"requireandverifyclientcert": tls.RequireAndVerifyClientCert,
}

// BytesToX509Cert converts bytes (PEM or DER) to an X509 certificate
func BytesToX509Cert(bytes []byte) (*x509.Certificate, error) {
	dcert, _ := pem.Decode(bytes)
	if dcert != nil {
		bytes = dcert.Bytes
	}
	cert, err := x509.ParseCertificate(bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Buffer was neither PEM nor DER encoding")
	}
	return cert, err
}

// UnmarshalConfig unmarshals a configuration file
func UnmarshalConfig(config interface{}, vp *viper.Viper, configFile string,
	server bool) error {

	vp.SetConfigFile(configFile)
	err := vp.ReadInConfig()
	if err != nil {
		return errors.Wrapf(err, "Failed to read config file '%s'", configFile)
	}

	err = vp.Unmarshal(config)
	if err != nil {
		return errors.Wrapf(err, "Incorrect format in file '%s'", configFile)
	}

	return nil
}
