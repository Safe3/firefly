package config

import (
	"fahi/pkg/util"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Version               string `json:"version"`
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	AutoSSL               bool   `json:"auto_ssl"`
	Password              string `json:"password"`
	Lang                  string `json:"lang"`
	LogLevel              string `json:"log_level"`
	WgPrivateKey          string `json:"wg_private_key"`
	WgDevice              string `json:"wg_device"`
	WgPort                int    `json:"wg_port"`
	WgMTU                 int    `json:"wg_mtu"`
	WgPersistentKeepalive int    `json:"wg_persistent_keepalive"`
	WgAddress             string `json:"wg_address"`
	WgDNS                 string `json:"wg_dns"`
	WgAllowedIPs          string `json:"wg_allowed_ips"`
}

func LoadOrCreate() (*Config, error) {
	var cfg Config

	cfgPath := util.RootDir + "conf/config.json"
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		wgDevice := "eth0"
		r, err := util.NewRouter()
		if err == nil {
			iface, _, _, err := r.Route(net.IPv4(0, 0, 0, 0))
			if err == nil {
				wgDevice = iface.Name
			}
		}

		host := ""
		ip, err := util.GetExternalIP(7 * time.Second)
		if err == nil {
			host = ip.String()
		}

		wgDeviceEnv := os.Getenv("FIREFLY_DEVICE")
		if wgDeviceEnv != "" {
			wgDevice = wgDeviceEnv
		}

		password := os.Getenv("FIREFLY_PASSWORD")
		if password == "" {
			password = "firefly"
		}

		port, err := strconv.Atoi(os.Getenv("FIREFLY_PORT"))
		if err != nil {
			port = 50121
		}

		autoSSL := false
		autoSslEnv := strings.ToLower(os.Getenv("FIREFLY_AUTO_SSL"))
		if autoSslEnv == "true" {
			autoSSL = true
		}

		cfg = Config{
			Version:               "2",
			Host:                  host,
			Port:                  port,
			AutoSSL:               autoSSL,
			Lang:                  "cn",
			LogLevel:              "error",
			Password:              password,
			WgPrivateKey:          util.GeneratePrivateKey(),
			WgDevice:              wgDevice,
			WgPort:                50120,
			WgMTU:                 1280,
			WgPersistentKeepalive: 25,
			WgAddress:             "198.18.0.1/16",
			WgDNS:                 "1.1.1.1",
			WgAllowedIPs:          "0.0.0.0/0, ::/0",
		}

		err = Save(&cfg)
		if err != nil {
			return nil, err
		}

		return &cfg, nil
	}

	err = util.Json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	data, err := util.Json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return err
	}

	path := util.RootDir + "conf"
	if _, err = os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0751)
	}

	return os.WriteFile(path+"/config.json", data, 0600)
}
