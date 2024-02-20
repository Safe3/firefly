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
  <a href="#特色">特色</a> •
  <a href="#使用">使用</a> •
  <a href="#感谢">感谢</a> •
  <a href="#联系">联系</a> •
  <a href="#授权">授权</a>
</p>





<p align="center">
  <a href="https://github.com/Safe3/firefly/blob/main/README.md">English</a>
  <a href="https://github.com/Safe3/firefly/blob/main/README_CN.md">中文</a>
</p>


---

萤火虫是一款简单、易架设的 WireGuard 服务端软件，可广泛用于异地组网、远程办公、内网穿透等场景。



## 特色

<h3 align="center">
  <img src="https://github.com/Safe3/firefly/blob/main/firefly.png" alt="firefly" width="700px">
  <br>
</h3>


 - 提供简单、易用的web管理后台
 - 支持所有 WireGuard 客户端接入
 - 无需系统安装 WireGuard 组件
 - 单文件、无额外库依赖





## 使用

萤火虫支持Linux x86、ARM等CPU架构环境，萤火虫服务端下载地址:  https://github.com/Safe3/firefly/releases  ,WireGuard客户端下载地址:  https://www.wireguard.com/install/ 。




### 服务端配置

首次运行firefly会在软件根目录生成config.json配置文件，如下：

```json
{
 "version": "1",              // 萤火虫当前版本
 "host": "7.7.7.7",           // 萤火虫web管理后台ip或域名
 "port": 50121,               // 萤火虫web管理后台端口
 "auto_ssl": true,            // 萤火虫web管理后台是否启用自动获取Let's Encrypt签发证书，若启用请将端口改为443
 "password": "firefly",       // 萤火虫web管理后台登录认证密码
 "lang": "en",                // 萤火虫web管理后台多语言支持，中文请将en改为cn
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


### 客户端配置

在服务端web管理后台创建完多个客户端后，通过以下方式导入WireGuard客户端配置。

1.WireGuard移动端可直接扫描萤火虫后台二维码导入配置

2.WireGuard PC端可下载萤火虫后台配置文件到本地后导入配置





## 感谢

感谢所有了不起的[社区贡献者发送PR](https://github.com/safe3/cvs/graphs/contributors)并不断更新此项目。请支持我们的朋友点个:heart:赞。

如果你有想法或某种改进，欢迎你贡献并参与该项目，随时发送你的PR。

<p align="center">
<a href="https://github.com/Safe3/firefly/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Safe3/firefly&max=500">
</a>
</p>

## 联系

若想支持更多功能，如权限分组、高级路由、堡垒机、点对点传输等，请访问: https://fahi.uusec.com



## 授权

firefly 仅用于个人免费使用，如要进行商业用途请联系我们获取商业授权。
