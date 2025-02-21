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
  <a href="#dart-features">Features</a> •
  <a href="#rocket-usage">Usage</a> •
  <a href="#gift_heart-credits">Credits</a> •
  <a href="#kissing_heart-contact">Contact</a> •
  <a href="#key-license">License</a>
</p>






<p align="center">
  <a href="https://github.com/Safe3/firefly/blob/main/README_CN.md">中文</a>
  <br/><br/>
  ⭐Please help us with a star to support our continuous improvement, thank you!
</p>




---

Firefly is a simple and easy to install WireGuard server software, which can be widely used in scenarios such as remote networking, remote work, and  expose a local server behind a NAT or firewall to the internet.

<h3 align="center">
  <img src="https://github.com/Safe3/firefly/blob/main/firefly.png" alt="firefly" width="700px">
  <br>
</h3>

## :dart: Features
:green_circle: Provide a beautiful, simple, and easy-to-use web management UI

 :purple_circle: Support all native WireGuard client access

 :yellow_circle: Compact and lightweight, less than 13M in size, does not rely on WireGuard

 :orange_circle: Developed in Go language, single file, high-performance, supports multi CPU architecture

 :red_circle: Support automatic application and renewal of free SSL certificates

 :large_blue_circle: Support TCP protocol relay and prevent UDP QoS flow limitation (advanced version)

## :rocket: Usage

Firefly supports CPU architecture environments such as Linux x86 and ARM. The download address for the Firefly server and WireGuard clients is: https://github.com/Safe3/firefly/releases .


- ### Server Installation

Select the corresponding binary, such as x86 environment, please download [firefly-linux-amd64](https://github.com/Safe3/firefly/releases/download/v4.4/firefly-linux-amd64)

Prepare to install：

```bash
mkdir -p /opt/firefly && mv firefly-linux-amd64 /opt/firefly/firefly 
```

Install as a service：

```bash
chmod +x /opt/firefly/firefly && sudo /opt/firefly/firefly -s install
```

Start firefly service：

```bash
sudo /opt/firefly/firefly -s start
```

If you want to run it in a container, just download docker-compose.yml and execute the following:

```bash
docker compose up -d
```

Visit http://ip:50121 ,login to the management with the default password firefly

> :biohazard: ***If the server is using cloud services, remember to open the UDP port 50120 and TCP port 50121-50122 required for Firefly***




- ### Server Configuration

The first time running firefly will generate a conf/config.json configuration file in the software root directory, as follows:

```json
{
 "version": 4.3,              // Firefly current version
 "host": "7.7.7.7",           // Firefly web management IP or domain name
 "port": 50121,               // Firefly web management port
 "auto_ssl": false,           // Is the firefly web enabled to automatically obtain Let's Encrypt certificate issuance? If enabled, please change the port to 443
 "password": "firefly",       // Firefly web management login authentication password
 "lang": "en",                // Firefly web management UI language
 "ui_traffic_stats": true,    // Firefly web management traffic chart switch
 "ui_chart_type": 2,          // Firefly web management traffic chart type
 "log_level": "error",        // Firefly server logging level
 "wg_private_key": "YBw5KAo1vM2mz35GLhZB01ZNYWJYWdGZNQT1MebuCHk=",  // WireGuard server private key
 "wg_device": "eth0",                   // WireGuard server in/out traffic network card name
 "wg_port": 50120,                      // WireGuard server UDP port
 "wg_mtu": 1280,                        // WireGuard server MTU value
 "wg_persistent_keepalive": 25,         // WireGuard client keepalive packet sending interval time
 "wg_address": "198.18.0.1/15",         // WireGuard client virtual IP network range
 "wg_dns": "8.8.8.8",                   // WireGuard client DNS configuration
 "wg_allowed_ips": "0.0.0.0/0, ::/0",   // WireGuard client allowed ips
 "wg_proxy_address": ":50122"           // TCP relay listening address,which can prevent UDP QoS flow limitation
}
```



- ### Client configuration

After creating multiple clients in the web management UI on the server side, import the WireGuard client configuration in the following way.

1.WireGuard mobile client can directly scan the firefly web QR code to import configuration

2.WireGuard PC client can download the firefly web configuration file to the local device and import the configuration





## :gift_heart: Credits

Thanks to all the amazing [community contributors for sending PRs](https://github.com/Safe3/firefly/graphs/contributors) and keeping this project updated. ❤️

If you have an idea or some kind of improvement, you are welcome to contribute and participate in the Project, feel free to send your PR.

<p align="center">
<a href="https://github.com/Safe3/firefly/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Safe3/firefly&max=500">
</a>
</p>

## :kissing_heart: Contact

If you want to support more features such as access controling, advanced routing, bastion machines, peer-to-peer transmission, etc, please contact us.



## :key: License

Firefly is only for personal free use. The front-end of this project is modified from wg-easy , follow the original project CC 4.0 license, thanks for the original auther Emile Nijssen！

