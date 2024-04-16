<h1 align="center">
  <br>
  <img src="https://github.com/Safe3/firefly/blob/main/logo.png" alt="firefly" width="70px">
</h1>
<h4 align="center">萤火虫 WireGuard 服务器</h4>

<p align="center">
<a href="https://github.com/Safe3/firefly/releases"><img src="https://img.shields.io/github/downloads/Safe3/firefly/total">
<a href="https://github.com/Safe3/firefly/graphs/contributors"><img src="https://img.shields.io/github/contributors-anon/Safe3/firefly">
<a href="https://github.com/Safe3/firefly/releases/"><img src="https://img.shields.io/github/release/Safe3/firefly">
<a href="https://github.com/Safe3/firefly/issues"><img src="https://img.shields.io/github/issues-raw/Safe3/firefly">
<a href="https://github.com/Safe3/firefly/discussions"><img src="https://img.shields.io/github/discussions/Safe3/firefly">
</p>
<p align="center">
  <a href="#dart-特色">特色</a> •
  <a href="#rocket-使用">使用</a> •
  <a href="#gift_heart-感谢">感谢</a> •
  <a href="#kissing_heart-联系">联系</a> •
  <a href="#key-授权">授权</a>
</p>




<p align="center">
  <a href="https://github.com/Safe3/firefly/blob/main/README.md">中文</a>
  <a href="https://github.com/Safe3/firefly/blob/main/README_EN.md">English</a>
</p>



---

萤火虫是一款简单、易架设的 WireGuard 服务端软件，可广泛用于异地组网、远程办公、内网穿透等场景。



## :dart: 特色

 :green_circle: 提供简单、易用的web管理后台

 :purple_circle: 支持所有 WireGuard 客户端接入

 :yellow_circle: 无需系统安装 WireGuard 组件

 :orange_circle: 单文件、无额外库依赖

 :red_circle: 自动申请免费SSL证书



<h3 align="center">
  <img src="https://github.com/Safe3/firefly/blob/main/firefly_cn.png" alt="firefly" width="700px">
  <br>
</h3>




## :rocket: 使用

萤火虫支持Linux x86、ARM等CPU架构环境，萤火虫服务端和WireGuard客户端下载地址:  https://github.com/Safe3/firefly/releases  ,其中以firefly-x-x开头的是服务端，包含wireguard名称的是客户端。



- ### 服务端安装

选择对应的服务端，如x86环境请下载firefly-linux-amd64

添加可执行权限：

```bash
chmod +x ./firefly-linux-amd64
```

前台运行：

```bash
./firefly-linux-amd64
```

后台运行：

```bash
nohup ./firefly-linux-amd64 >/dev/null 2>&1 &
```

容器中运行：下载docker-compose.yml文件然后执行

```bash
docker compose up -d
```

访问 http://ip:50121 登录管理后台，默认密码firefly

> :biohazard: ***如果服务器使用的是各种云服务，记得在云服务管理后台上开放萤火虫所需的udp端口50120、tcp端口50121***



- ### 服务端配置


首次运行firefly会在软件目录生成conf/config.json配置文件，配置说明如下：

```json
{
 "version": "1",              // 萤火虫当前版本
 "host": "7.7.7.7",           // 萤火虫web管理后台ip或域名
 "port": 50121,               // 萤火虫web管理后台端口
 "auto_ssl": false,           // 萤火虫web管理后台是否启用自动获取Let's Encrypt签发证书，若启用请将端口改为443
 "password": "firefly",       // 萤火虫web管理后台登录认证密码
 "lang": "en",                // 萤火虫web管理后台多语言支持，中文请将en改为cn
 "ui_traffic_stats": true,    // 萤火虫web管理后台是否开启流量图特效
 "ui_chart_type": 2,          // 萤火虫web管理后台流量特效图类型
 "log_level": "error",        // 萤火虫服务端日志记录等级
 "wg_private_key": "YBw5KAo1vM2mz35GLhZB01ZNYWJYWdGZNQT1MebuCHk=",  // 萤火虫服务端 WireGuard 私钥
 "wg_device": "eth0",                   // 萤火虫服务端 WireGuard 出入流量网卡名称
 "wg_port": 50120,                      // 萤火虫服务端 WireGuard UDP端口
 "wg_mtu": 1280,                        // 萤火虫服务端 WireGuard MTU值
 "wg_persistent_keepalive": 25,         // 萤火虫客户端存活包发送间隔时间
 "wg_address": "198.18.0.1/16",         // 萤火虫服务端ip和网段范围
 "wg_dns": "1.1.1.1",                   // 萤火虫客户端dns配置
 "wg_allowed_ips": "0.0.0.0/0, ::/0"    // 萤火虫客户端流量要转发到服务端的ip地址范围
}
```



- ### 客户端安装

萤火虫的客户端为WireGuard官方客户端，支持Windows、Linux、Mac、iOS、Android，这里以Windows为例，对应的客户端为wireguard-amd64-0.5.3.msi ，下载后根据提示一步步安装。



- ### 客户端配置

登录萤火虫服务端web管理后台，新建2个客户端，通过以下方式导入WireGuard客户端配置。

1.移动客户端可直接扫描萤火虫后台二维码导入配置

2.PC客户端可下载萤火虫后台配置文件到本地后导入配置

两个客户端开启之后，可以通过萤火虫服务端分配的ip 198.18.0.x 直接相互访问



## :gift_heart: 感谢

感谢所有了不起的[社区贡献者发送PR](https://github.com/safe3/cvs/graphs/contributors)并不断更新此项目。请支持我们的朋友点个 :heart: 赞。

如果你有想法或某种改进，欢迎你贡献并参与该项目，随时发送你的PR。

<p align="center">
<a href="https://github.com/Safe3/firefly/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Safe3/firefly&max=500">
</a>
</p>
捐赠请扫描如下二维码：
<img src="https://waf.uusec.com/_media/sponsor.jpg" alt="捐赠"  height="300px" />



## :kissing_heart: 联系

若想支持更多功能，如权限分组、高级路由、堡垒机、点对点传输等，请访问: https://fahi.uusec.com

- 问题提交：https://github.com/Safe3/firefly/issues

- 讨论社区：https://github.com/Safe3/firefly/discussions

- 官方 QQ 群：11500614

- 官方微信群：微信扫描以下二维码加入

  <img src="https://waf.uusec.com/_media/weixin.jpg" alt="微信群"  height="200px" />



## :key: 授权

firefly 仅用于个人免费使用，本项目前端部分来源于[wg-easy](https://github.com/wg-easy/wg-easy) ，遵循原项目CC 4.0协议。
