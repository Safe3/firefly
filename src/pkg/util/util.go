package util

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var (
	Json    = jsoniter.ConfigCompatibleWithStandardLibrary
	RootDir = ""
)

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	RootDir = filepath.Dir(ex) + "/"
}

func GenerateKey() string {
	key, err := wgtypes.GenerateKey()
	if err == nil {
		return key.String()
	}
	return ""
}

func GeneratePrivateKey() string {
	key, err := wgtypes.GeneratePrivateKey()
	if err == nil {
		return key.String()
	}
	return ""
}

func GetExternalIP(timeout time.Duration) (net.IP, error) {
	// Define the GET method with the correct url,
	// setting the User-Agent to our library
	req, err := http.NewRequest("GET", "https://checkip.amazonaws.com/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "firefly")

	// transport to avoid goroutine leak
	tr := &http.Transport{
		MaxIdleConns:      1,
		IdleConnTimeout:   3 * time.Second,
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
			Control: func(network, address string, c syscall.RawConn) error {
				return nil
			},
		}).DialContext,
	}

	client := &http.Client{Timeout: timeout, Transport: tr}

	// Do the request and read the body for non-error results.
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// validate the IP
	raw := string(bytes)
	externalIP := net.ParseIP(strings.TrimSpace(raw))
	if externalIP == nil {
		return nil, fmt.Errorf("[ERROR] returned an invalid IP: %s\n", raw)
	}

	// returned the parsed IP
	return externalIP, nil
}
