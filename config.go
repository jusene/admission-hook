package main

import (
	"crypto/tls"
	"k8s.io/klog/v2"
)

// Config contains the server (the webhook) cert and key.
type Config struct {
	CertFile string
	KeyFile  string
}

func configTLS(config Config) *tls.Config {
	sCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
	if err != nil {
		klog.Fatal(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
}