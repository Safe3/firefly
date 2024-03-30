<h1 align="center">
  <br>
  <img src="https://github.com/Safe3/firefly/blob/main/logo.png" alt="firefly" width="70px">
</h1>
<h4 align="center">Firefly WireGuard Server</h4>

<p align="center">
<a href="https://github.com/Safe3/firefly/releases"><img src="https://img.shields.io/github/downloads/Safe3/firefly/total">
<a href="https://github.com/Safe3/firefly/graphs/contributors"><img src="https://img.shields.io/github/contributors-anon/Safe3/firefly">
<a href="https://github.com/Safe3/firefly/releases/"><img src="https://img.shields.io/github/release/Safe3/firefly">
<a href="https://github.com/Safe3/firefly/issues"><img src="https://img.shields.io/github/issues-raw/Safe3/firefly">
<a href="https://github.com/Safe3/firefly/discussions"><img src="https://img.shields.io/github/discussions/Safe3/firefly">
</p>
<p align="center">
  <a href="#features">Features</a> •
  <a href="#usage">Usage</a> •
  <a href="#credits">Credits</a> •
  <a href="#contact">Contact</a> •
  <a href="#license">License</a>
</p>






<p align="center">
  <a href="https://github.com/Safe3/firefly/blob/main/README.md">English</a>
  <a href="https://github.com/Safe3/firefly/blob/main/README_CN.md">中文</a>
</p>


---

Firefly is a simple and easy to install WireGuard server software, which can be widely used in scenarios such as remote networking, remote work, and  expose a local server behind a NAT or firewall to the internet.



## Features

<h3 align="center">
  <img src="https://github.com/Safe3/firefly/blob/main/firefly.png" alt="firefly" width="700px">
  <br>
</h3>


 - Provide a simple and easy-to-use web management UI
 - Supports access to all WireGuard clients
 - No need for system installation of WireGuard components
 - Single file, no additional library dependencies





## Usage

Firefly supports CPU architecture environments such as Linux x86 and ARM. The download address for the Firefly server is: https://github.com/Safe3/firefly/releases , WireGuard client download address: https://www.wireguard.com/install/ .







### Server Configuration

The first time running firefly will generate a config.json configuration file in the software root directory, as follows:

```json
{
 "version": "1",              // Firefly current version
 "host": "7.7.7.7",           // Firefly web management IP or domain name
 "port": 50121,               // Firefly web management port
 "auto_ssl": false,           // Is the firefly web enabled to automatically obtain Let's Encrypt certificate issuance? If enabled, please change the port to 443
 "password": "firefly",       // Firefly web management login authentication password
 "lang": "en",                // Firefly web management UI language
 "log_level": "error",        // Firefly server logging level
 "wg_private_key": "YBw5KAo1vM2mz35GLhZB01ZNYWJYWdGZNQT1MebuCHk=",  // WireGuard server private key
 "wg_device": "eth0",                   // WireGuard server in/out traffic network card name
 "wg_port": 50120,                      // WireGuard server UDP port
 "wg_mtu": 1280,                        // WireGuard server MTU value
 "wg_persistent_keepalive": 25,         // WireGuard client keepalive packet sending interval time
 "wg_address": "198.18.0.1/16",         // WireGuard server IP and network range
 "wg_dns": "1.1.1.1",                   // WireGuard client DNS configuration
 "wg_allowed_ips": "0.0.0.0/0, ::/0"    // WireGuard client allowed ips
}
```



### Client configuration

After creating multiple clients in the web management UI on the server side, import the WireGuard client configuration in the following way.

1.WireGuard mobile client can directly scan the firefly web QR code to import configuration

2.WireGuard PC client can download the firefly web configuration file to the local device and import the configuration





## Credits

Thanks to all the amazing [community contributors for sending PRs](https://github.com/Safe3/firefly/graphs/contributors) and keeping this project updated. ❤️

If you have an idea or some kind of improvement, you are welcome to contribute and participate in the Project, feel free to send your PR.

<p align="center">
<a href="https://github.com/Safe3/firefly/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Safe3/firefly&max=500">
</a>
</p>

## Contact

If you want to support more features such as access controling, advanced routing, bastion machines, peer-to-peer transmission, etc., please contact us.



## License

Firefly is only for personal free use. If you want to use it for commercial purposes, please contact us for commercial authorization.
