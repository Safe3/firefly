package wg

import (
	"fahi/pkg/config"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/coreos/go-iptables/iptables"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/device"
)

type WgIface struct {
	Name         string
	Device       string
	AdminPort    int
	Port         int
	MTU          int
	Address      *WgAddress
	Interface    NetInterface
	privateKey   string
	mu           sync.Mutex
	tunDevice    *device.Device
	uapiListener net.Listener
	iptables     *iptables.IPTables
}

// WGAddress Wireguard parsed address
type WgAddress struct {
	IP      net.IP
	Network *net.IPNet
}

func (addr *WgAddress) String() string {
	maskSize, _ := addr.Network.Mask.Size()
	return fmt.Sprintf("%s/%d", addr.IP.String(), maskSize)
}

func parseAddress(address string) (*WgAddress, error) {
	ip, network, err := net.ParseCIDR(address)
	if err != nil {
		return nil, err
	}

	return &WgAddress{
		IP:      ip,
		Network: network,
	}, nil
}

// NetInterface represents a generic network tunnel interface
type NetInterface interface {
	Close() error
}

func New(cfg *config.Config) (*WgIface, error) {
	wgIface := &WgIface{
		Name:       "firefly",
		Device:     cfg.WgDevice,
		AdminPort:  cfg.Port,
		Port:       cfg.WgPort,
		MTU:        cfg.WgMTU,
		mu:         sync.Mutex{},
		privateKey: cfg.WgPrivateKey,
	}

	wgAddress, err := parseAddress(cfg.WgAddress)
	if err != nil {
		return wgIface, err
	}
	wgIface.Address = wgAddress

	return wgIface, nil
}

func (w *WgIface) Create() (err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if WireguardModuleIsLoaded() {
		log.Info("using kernel WireGuard")
		err = w.createWithKernel()
	} else {
		if !tunModuleIsLoaded() {
			return fmt.Errorf("couldn't check or load tun module")
		}
		log.Info("using userspace WireGuard")
		err = w.createWithUserspace()
	}
	if err != nil {
		return
	}

	err = setIPForwarding(true)
	if err != nil {
		return
	}

	w.iptables, err = iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		return
	}
	err = w.iptables.Insert("filter", "FORWARD", 1, "-i", w.Name, "-j", "ACCEPT")
	if err != nil {
		return
	}
	err = w.iptables.Insert("filter", "FORWARD", 1, "-o", w.Name, "-j", "ACCEPT")
	if err != nil {
		return
	}
	err = w.iptables.Insert("nat", "POSTROUTING", 1, "-s", w.Address.Network.String(), "-o", w.Device, "-j", "MASQUERADE")
	if err != nil {
		return
	}
	err = w.iptables.Insert("filter", "INPUT", 1, "-p", "udp", "--dport", fmt.Sprintf("%d", w.Port), "-j", "ACCEPT")
	if err != nil {
		return
	}
	err = w.iptables.Insert("filter", "INPUT", 1, "-p", "tcp", "--dport", fmt.Sprintf("%d", w.AdminPort), "-j", "ACCEPT")
	return
}

func (w *WgIface) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var err error

	if w.tunDevice != nil {
		w.tunDevice.Close()
		w.tunDevice = nil
	} else if w.Interface != nil {
		err = w.Interface.Close()
		w.Interface = nil
		if err != nil {
			log.Debugf("failed to close interface: %s", err)
		}
	}

	sockPath := "/var/run/wireguard/" + w.Name + ".sock"
	if _, statErr := os.Stat(sockPath); statErr == nil {
		_ = os.Remove(sockPath)
	}

	if w.uapiListener != nil {
		err = w.uapiListener.Close()
		w.uapiListener = nil
		if err != nil {
			log.Errorf("failed to close uapi listener: %v", err)
		}
	}

	setIPForwarding(false)

	if w.iptables != nil {
		w.iptables.Delete("filter", "FORWARD", "-i", w.Name, "-j", "ACCEPT")
		w.iptables.Delete("filter", "FORWARD", "-o", w.Name, "-j", "ACCEPT")
		w.iptables.Delete("nat", "POSTROUTING", "-s", w.Address.Network.String(), "-o", w.Device, "-j", "MASQUERADE")
		w.iptables.Delete("filter", "INPUT", "-p", "udp", "--dport", fmt.Sprintf("%d", w.Port), "-j", "ACCEPT")
		w.iptables.Delete("filter", "INPUT", "-p", "tcp", "--dport", fmt.Sprintf("%d", w.AdminPort), "-j", "ACCEPT")
		w.iptables = nil
	}

	return err
}

func setIPForwarding(enabled bool) error {
	ipv4ForwardingPath := "/proc/sys/net/ipv4/ip_forward"
	bytes, err := os.ReadFile(ipv4ForwardingPath)
	if err != nil {
		return err
	}

	if len(bytes) > 0 {
		if enabled && bytes[0] == 49 {
			return nil
		} else if !enabled && bytes[0] == 48 {
			return nil
		}
	}

	if enabled {
		return os.WriteFile(ipv4ForwardingPath, []byte("1"), 0644)
	}

	return os.WriteFile(ipv4ForwardingPath, []byte("0"), 0644)
}
