# 初始化环境搭建
## 设置固定ip
> Ubuntu 18.04 版本 服务端版本默认开启为 systemd-networkd ,桌面版本默认开启为 NetworkManager
### systemd-networkd 设置
**启用 systemd-networkd**

```shell
sudo systemctl unmask systemd-networkd.service
sudo systemctl enable systemd-networkd.service
sudo systemctl start systemd-networkd.service
```
**编辑网卡配置**
dir:/etc/netplan 目录下的以 yaml 后缀结尾的都会加载,若没有则创建

```shell
sudo vim /etc/netplan/config.yaml
network:
 version: 2
 renderer: networkd
 ethernets:
  eth0: # 网卡名称
   addresses:
    - 192.168.1.88/24 # 地址和子网掩码
   gateway4: 192.168.1.0 # 网关
   nameservers:
    addresses: [114.114.114.114,8.8.8.8] # dns
```
**应用**

```shell
sudo netplan apply
```
### NetworkManager 设置
桌面端默认开启,直接在系统设置里面找到相应的配置即可;

**若 systemd-networkd 有进程开启，则需要进去 netplan 的配置文件将 renderer 字段的值修改为NetworkManager ，并且执行 netplan apply**

## 开启 ssh

**编辑配置文件**

```shell
vim  /etc/ssh/sshd_config
# 修改内容:
PermitEmptyPasswords yes
PermitRootLogin yes
```

**重启ssh服务**

```shell
sudo service sshd restart
```

## 在线安装docker、docker-compose
基础操作在线网站：https://yeasy.gitbook.io/docker_practice/install/ubuntu 

安装后记得配置国内的镜像地址，Docker Hub 大陆拉取较慢。

**踩坑：**
1. RK3399 板子安装 Docker 若启动失败则是内核需要更新，找购买的商家要即可(云盘有)
2. RK3399 板子安装 docker-compose 以下版本可用：
```shell
sudo curl -L "https://github.com/docker/compose/releases/download/v2.2.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo docker-compose --version
```