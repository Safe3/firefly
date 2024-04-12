package wg

import (
	cfg "fahi/pkg/config"
	"fahi/pkg/util"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type Peer struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	PrivateKey   string    `json:"privateKey"`
	PreSharedKey string    `json:"preSharedKey"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Enabled      bool      `json:"enabled"`
}

type PeerStatus struct {
	Id                  string    `json:"id"`
	Name                string    `json:"name"`
	Address             string    `json:"address"`
	PublicKey           string    `json:"publicKey"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
	Enabled             bool      `json:"enabled"`
	PersistentKeepalive int64     `json:"persistentKeepalive"`
	LatestHandshakeAt   time.Time `json:"latestHandshakeAt"`
	TransferRx          int64     `json:"transferRx"`
	TransferTx          int64     `json:"transferTx"`
}

// configureDevice configures the wireguard device
func (w *WgIface) configureDevice(config wgtypes.Config) error {
	wgc, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer wgc.Close()

	// validate if device with name exists
	_, err = wgc.Device(w.Name)
	if err != nil {
		return err
	}
	log.Debugf("got Wireguard device %s", w.Name)

	return wgc.ConfigureDevice(w.Name, config)
}

func loadPeers() ([]Peer, error) {
	peers := []Peer{}

	path := util.RootDir + "conf/peers.json"

	data, err := os.ReadFile(path)
	if err != nil {
		data = []byte("[]")
		err = os.WriteFile(path, data, 0600)
		if err != nil {
			return peers, err
		}
	}

	err = util.Json.Unmarshal(data, &peers)
	if err != nil {
		return peers, err
	}

	return peers, nil
}

func savePeers(peers []Peer) error {
	data, err := util.Json.MarshalIndent(&peers, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(util.RootDir+"conf/peers.json", data, 0600)
}

func (w *WgIface) configure() error {
	key, err := wgtypes.ParseKey(w.privateKey)
	if err != nil {
		return err
	}
	fwmark := 0
	config := wgtypes.Config{
		PrivateKey:   &key,
		ReplacePeers: true,
		FirewallMark: &fwmark,
		ListenPort:   &w.Port,
	}

	peers, err := loadPeers()
	if err != nil {
		return err
	}

	peersConfig := []wgtypes.PeerConfig{}
	for _, peer := range peers {
		if peer.Enabled {
			_, ipNet, err := net.ParseCIDR(peer.Address + "/32")
			if err != nil {
				return err
			}

			privateKey, err := wgtypes.ParseKey(peer.PrivateKey)
			if err != nil {
				return err
			}

			peerConfig := wgtypes.PeerConfig{
				PublicKey:  privateKey.PublicKey(),
				AllowedIPs: []net.IPNet{*ipNet},
			}

			peersConfig = append(peersConfig, peerConfig)
		}
	}

	config.Peers = peersConfig

	err = w.configureDevice(config)
	if err != nil {
		return fmt.Errorf("received error \"%v\" while configuring interface %s with port %d", err, w.Name, w.Port)
	}

	return nil
}

func (w *WgIface) updatePeer(peerKey string, allowedIps string, keepAlive time.Duration) error {
	log.Debugf("updating interface %s peer %s: endpoint %s ", w.Name, peerKey)

	//parse allowed ips
	AllowedIPs := []net.IPNet{}
	ais := strings.Split(allowedIps, ",")

	for _, ai := range ais {
		_, ipNet, err := net.ParseCIDR(ai)
		if err != nil {
			return err
		}
		AllowedIPs = append(AllowedIPs, *ipNet)
	}

	peerKeyParsed, err := wgtypes.ParseKey(peerKey)
	if err != nil {
		return err
	}
	peer := wgtypes.PeerConfig{
		PublicKey:                   peerKeyParsed.PublicKey(),
		ReplaceAllowedIPs:           true,
		AllowedIPs:                  AllowedIPs,
		PersistentKeepaliveInterval: &keepAlive,
	}

	config := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peer},
	}
	err = w.configureDevice(config)
	if err != nil {
		return fmt.Errorf("received error \"%v\" while updating peer on interface %s with settings: allowed ips %s", err, w.Name, allowedIps)
	}
	return nil
}

func (w *WgIface) GetPeers() ([]PeerStatus, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	wgc, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = wgc.Close()
		if err != nil {
			log.Errorf("got error while closing wgctl: %v", err)
		}
	}()

	wgDevice, err := wgc.Device(w.Name)
	if err != nil {
		return nil, err
	}

	peers, err := loadPeers()
	if err != nil {
		return nil, err
	}

	peersStatus := []PeerStatus{}

	for _, peer := range peers {
		privateKey, err := wgtypes.ParseKey(peer.PrivateKey)
		if err != nil {
			return nil, err
		}
		publicKey := privateKey.PublicKey().String()

		peerStatus := PeerStatus{
			Id:        peer.Id,
			Name:      peer.Name,
			Address:   peer.Address,
			CreatedAt: peer.CreatedAt,
			UpdatedAt: peer.UpdatedAt,
			Enabled:   peer.Enabled,
			PublicKey: publicKey,
		}

		for _, devPeer := range wgDevice.Peers {
			if devPeer.PublicKey.String() == publicKey {
				peerStatus.PersistentKeepalive = int64(devPeer.PersistentKeepaliveInterval)
				peerStatus.LatestHandshakeAt = devPeer.LastHandshakeTime
				peerStatus.TransferRx = devPeer.ReceiveBytes
				peerStatus.TransferTx = devPeer.TransmitBytes
				break
			}
		}
		peersStatus = append(peersStatus, peerStatus)
	}

	return peersStatus, nil
}

// RemovePeer removes a Wireguard Peer from the interface iface
func (w *WgIface) removePeer(peerKey string) error {
	log.Debugf("Removing peer %s from interface %s ", peerKey, w.Name)

	peerKeyParsed, err := wgtypes.ParseKey(peerKey)
	if err != nil {
		return err
	}

	peer := wgtypes.PeerConfig{
		PublicKey: peerKeyParsed.PublicKey(),
		Remove:    true,
	}

	config := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peer},
	}
	err = w.configureDevice(config)
	if err != nil {
		return fmt.Errorf("received error \"%v\" while removing peer %s from interface %s", err, peerKey, w.Name)
	}
	return nil
}

func (w *WgIface) GetPeerConfig(id string) (string, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	peers, err := loadPeers()
	if err != nil {
		return "", err
	}

	c, err := cfg.LoadOrCreate()
	if err != nil {
		return "", err
	}

	for _, peer := range peers {
		if peer.Id == id {
			privateKey, err := wgtypes.ParseKey(c.WgPrivateKey)
			if err != nil {
				return "", err
			}
			peerConfig := "[Interface]\nPrivateKey = " + peer.PrivateKey + "\nAddress = " + peer.Address + "/16\n" + "DNS = " + c.WgDNS + "\nMTU = " + fmt.Sprintf("%d", c.WgMTU) + "\n\n"
			peerConfig = peerConfig + "[Peer]\nPublicKey = " + privateKey.PublicKey().String() + "\nAllowedIPs = " + c.WgAllowedIPs + "\nPersistentKeepalive = " + fmt.Sprintf("%d", c.WgPersistentKeepalive)
			peerConfig = peerConfig + "\nEndpoint = " + c.Host + ":" + fmt.Sprintf("%d", c.WgPort)
			return peerConfig, nil
		}
	}

	return "", nil
}

func (w *WgIface) AddPeer(name string) (*Peer, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	peers, err := loadPeers()
	if err != nil {
		return nil, err
	}

	c, err := cfg.LoadOrCreate()
	if err != nil {
		return nil, err
	}

	address := c.WgAddress[:strings.LastIndex(c.WgAddress, ".")+1]

	if len(peers) == 0 {
		address += "2"
	} else {
		for i := 2; i < 255; i++ {
			found := false
			newIp := address + fmt.Sprintf("%d", i)
			for _, peer := range peers {
				if peer.Address == newIp {
					found = true
					break
				}
			}
			if !found {
				address = newIp
				break
			}
		}
	}

	if strings.HasSuffix(address, ".") {
		return nil, fmt.Errorf("Maximum number of clients reached")
	}

	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}

	peer := Peer{
		Id:         uuid.New().String(),
		Name:       name,
		Address:    address,
		PrivateKey: privateKey.String(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Enabled:    true,
	}

	err = w.updatePeer(peer.PrivateKey, peer.Address+"/32", time.Duration(c.WgPersistentKeepalive)*time.Second)
	if err != nil {
		return nil, err
	}

	peers = append(peers, peer)
	err = savePeers(peers)
	if err != nil {
		return nil, err
	}

	return &peer, nil
}

func (w *WgIface) DelPeer(id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	peers, err := loadPeers()
	if err != nil {
		return err
	}

	i := 0
	for _, peer := range peers {
		if peer.Id != id {
			peers[i] = peer
			i++
		} else {
			err = w.removePeer(peer.PrivateKey)
			if err != nil {
				return err
			}
		}
	}

	return savePeers(peers[:i])
}

func (w *WgIface) SetPeer(id, enabled, name, address string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	peers, err := loadPeers()
	if err != nil {
		return err
	}

	for i, peer := range peers {
		if peer.Id == id {
			peers[i].UpdatedAt = time.Now()
			if name != "" {
				peers[i].Name = name
			} else if address != "" {
				peers[i].Address = address
				c, err := cfg.LoadOrCreate()
				if err != nil {
					return err
				}
				err = w.updatePeer(peer.PrivateKey, address+"/32", time.Duration(c.WgPersistentKeepalive)*time.Second)
				if err != nil {
					return err
				}
			} else {
				if enabled == "enable" {
					peers[i].Enabled = true
					c, err := cfg.LoadOrCreate()
					if err != nil {
						return err
					}
					err = w.updatePeer(peer.PrivateKey, peer.Address+"/32", time.Duration(c.WgPersistentKeepalive)*time.Second)
					if err != nil {
						return err
					}
				} else if enabled == "disable" {
					peers[i].Enabled = false
					err = w.removePeer(peer.PrivateKey)
					if err != nil {
						return err
					}
				}
			}
			break
		}
	}

	return savePeers(peers)
}
