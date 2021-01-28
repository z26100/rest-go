package rest

import (
	"crypto/tls"
	log "github.com/z26100/log-go"
	flag "github.com/z26100/service-config-go"
	"net/http"
	"time"
)

const (
	defaultTimeout    = 120 * time.Second
	defaultListen     = "0.0.0.0:8080"
)

var (
	restListen       = flag.String("listen", defaultListen, "the web server address:port")
	restReadTimeout  = flag.Duration("readTimeout", defaultTimeout, "the http read timeout")
	restWriteTimeout = flag.Duration("writeTimeout", defaultTimeout, "the http write timeout")
	debug            = flag.Bool("debug", false, "debug mode")
	keyFile          = flag.String("keyFile", "", "the tls key file")
	crtFile          = flag.String("certFile", "", "the tls cert file")
)


func DefaultRestConfig() ServerConfig {
	return ServerConfig{
		Listen:         *restListen,
		ProductionMode: !*debug,
		CertFile:       *crtFile,
		KeyFile:        *keyFile,
		ReadTimeout:    *restReadTimeout,
		WriteTimeout:   *restWriteTimeout,
		TlsConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			}},
	}
}


func check(condition func() bool, w http.ResponseWriter) bool {
	if condition() {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return true
	}
	return false
}

func CheckError(err error, w http.ResponseWriter) bool {
	return check(func() bool {
		if err != nil {
			log.Errorln(err)
		}
		return err != nil
	}, w)
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
}
