# openwrt 使用说明
## 安装openwrt
### pe安装
1. 写盘工具以及镜像地址放入PE，esir网盘地址：https://drive.google.com/drive/folders/1dqNUrMf9n7i3y1aSh68U5Yf44WQ3KCuh
2. 进入PE
3. 使用写盘工具写入镜像至指定磁盘（需格式化）,写入cmd命令 tools.exe -u op.img ;选择磁盘

### dd命令安装
1. 上传镜像至tmp/upload目录
2. dd if=/tmp/upload/op.img of=/dev/sda
3. reboot


## 设置wan口
1. 状态PPPoE
2. 输入宽带账号密码
3. 设置wan口指定的物理网口

## 设置lan口
1. 设置lan口指定的物理网口

## 设置ssr plus 
1. 添加订阅地址
2. 设置主页的国内ip的dns为114.114.114.114

## 扩容overlay（需要额外安装软件就需要扩容，主要用于存储软件包）
1. 新建sda3分区
2. mkfs.ext4 /dev/sda3
3. mount /dev/sda3 /mnt/sda3  // mnt/sda3 目录出现lost+found文件就说明挂载成功
4. cp -r /overlay/* /mnt/sda3
5. 去web端的挂载点把sda3挂载到overlay

## 扩容docker
1. 新建sda4分区
2. 格式化分区mkfs.ext4 /dev/sda4
3. 去web端的挂载点把sda4挂载到docker

## 开启广告大师
1. 订阅地址如下：

https://cdn.jsdelivr.net/gh/o0HalfLife0o/list@master/ad-pc.txt

https://cdn.jsdelivr.net/gh/o0HalfLife0o/list@master/ad.txt

## 设置动态dns
1. 修改ipv4
2. 主机名：jiny28.top ddns服务提供商ipv4：cloudflare.com-v4  域名：@jiny28.top  用户名：chenjinyang6@gmail.com  密码：查询账号token

## 设置sdb磁盘
1. 创建分区使用cfdisk /dev/sdb
2. 格式化sdb1，界面即可操作
3. 挂载sdb1到/data目录
4. 重启

## docker 启动 jellyfin,tmm,clouddrive2
1. clouddrive2 启动需前置条件
```
mkdir -p /etc/systemd/system/docker.service.d/

cat <<EOF > /etc/systemd/system/docker.service.d/clear_mount_propagation_flags.conf
[Service]
MountFlags=shared
EOF

mount --make-shared /data

> 同时op的web管理里面的启动项也得加入：mount --make-shared /data
```

2. docker-compose.yml 如下
```docker-compose.yml
version: "3"

services:
  # jellyfin:
  #   container_name: jellyfin
  #   image: jellyfin/jellyfin
  #   volumes:
  #     - /data/jellyfin/config:/config
  #     - /data/jellyfin/cache:/cache
  #     - /data/jellyfin/media:/media
  #     - /data/clouddrive2/media:/cloudmedia
  #   network_mode: host
  #   restart: unless-stopped
  #   devices:
  #     - /dev/dri/renderD128:/dev/dri/renderD128
  #     - /dev/dri/card0:/dev/dri/card0
  #   depends_on:
  #     - clouddrive2

  # tinymediamanager:
  #   container_name: tinymediamanager
  #   image: romancin/tinymediamanager:latest-v4
  #   volumes:
  #     - /data/tinymediamanager/config:/config
  #     - /data/jellyfin/media:/media
  #     - /data/clouddrive2/media:/cloudmedia
  #   environment:
  #     - ENABLE_CJK_FONT=1
  #     - USER_ID=0
  #     - GROUP_ID=0
  #   user: 0:0
  #   ports:
  #     - 5800:5800
  #   restart: unless-stopped
  #   extra_hosts:
  #     - "image.tmdb.org:169.150.249.169"
  #     - "api.themoviedb.org:13.226.225.52" 
  #     - "www.themoviedb.org:13.226.228.83"
  #   depends_on:
  #     - clouddrive2

  alist:
    image: 'xhofe/alist:v3.44.0'
    container_name: alist
    volumes:
        - '/data/alist:/opt/alist/data'
        - '/data/jinydata:/opt/alist/jinydata'
    ports:
        - '5244:5244'
    environment:
        - PUID=0
        - PGID=0
        - UMASK=022
        - TZ=Asia/Shanghai
    restart: unless-stopped

  # clouddrive2:
  #   container_name: clouddrive2
  #   image: cloudnas/clouddrive2-unstable
  #   restart: unless-stopped
  #   environment: 
  #     - TZ=Asia/Shanghai
  #     - CLOUDDRIVE_HOME=/Config
  #   privileged: true
  #   devices:
  #     - /dev/fuse:/dev/fuse
  #   volumes:
  #     - /data/clouddrive2/shared:/CloudNAS:shared
  #     - /data/clouddrive2/Config:/Config
  #     - /data/clouddrive2/media:/media:shared
  #   ports:
  #     - 19798:19798
  ql:
    # alpine 基础镜像版本
    image: whyour/qinglong:latest
    container_name: qinglong
    volumes:
      - /data/ql/data:/ql/data
    ports:
      - "0.0.0.0:5700:5700"
    environment:
      # 部署路径非必须，以斜杠开头和结尾，比如 /test/
      QlBaseUrl: '/'
    restart: unless-stopped
  # homeassistant:
  #   container_name: homeassistant
  #   image: homeassistant/home-assistant:stable
  #   privileged: true
  #   restart: unless-stopped
  #   environment:
  #     TZ: Asia/Shanghai
  #   volumes:
  #     - /data/homeassistant:/config
  #   network_mode: host
  
  hbbs:
    container_name: hbbs
    ports:
      - 21115:21115
      - 21116:21116 # 自定义 hbbs 映射端口
      - 21116:21116/udp # 自定义 hbbs 映射端口
    image: rustdesk/rustdesk-server
    command: hbbs 
    volumes:
      - /data/rustdesk:/root # 自定义挂载目录
    depends_on:
      - hbbr
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 64M

  hbbr:
    container_name: hbbr
    ports:
      - 21117:21117 # 自定义 hbbr 映射端口
    image: rustdesk/rustdesk-server
    command: hbbr
    volumes:
      - /data/rustdesk:/root # 自定义挂载目录
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 64M
```

## jellyfin 常见问题解决
1. 解决字幕乱码
```
apt update && apt install fonts-noto-cjk-extra
```
2. 开启硬件加速

播放设置中开启硬件加速，选择VAAPI。

播放时，在 htop 中的进程明细中看到 -hwaccel vaapi ,就代表硬件加速已启动。

## tinymediamanager 常见问题解决
1. 解决界面中文乱码问题
```
wget https://mirrors.aliyun.com/alpine/edge/community/x86_64/font-wqy-zenhei-0.9.45-r3.apk
apk add --allow-untrusted font-wqy-zenhei-0.9.45-r3.apk
>restart container.
>apk not found , OneDrive 已存.
```
2. 设置访问密码
```
x11vnc -storepasswd
>查看输出的位置，一般为：/root/.vnc/passwd
cd /run/s6/etc/services.d/x11vnc
vi run
>输出的位置一致的情况，直接替换如下区域的代码：
# Handle the VNC password.
if [ -f /root/.vnc/passwd ] && [ -n "$( cat /root/.vnc/passwd )" ]; then          
    VNC_SECURITY="-rfbauth /root/.vnc/passwd"                                     
else                                                                              
    VNC_SECURITY="-nopw"                                                          
fi

> restart container .
```

## 开启 smb
1. web设置界面网络共享
2. 注销配置界面中的 # invalid users = root
3. 打开共享家目录
4. 如果不生效可以 service samba4 reload ; service samba4 restart


## 开启 Transmission
1. web设置界面Transmission
2. 用户组需要改为root,其他的按需求设置即可.

## 开启端口转发
1. web设置界面的防火墙设置
2. 设置端口转发，入口网口+端口，出口网口+端口即可.

## ipv6设置
1. lan 口高级设置中勾选使用内置的IPv6管理和强制链路
2. lan 口的dhcp 的IPv6 设置中设置路由通告服务和DHCPv6服务为服务器模式，总是通告默认路由为勾选
3. lan 口的IPv6 的长度设置为60
4. wan 口的使用内置的IPv6管理勾上，强制链路取消，获取IPv6地址为自动
5. wan6 口请求IPv6地址设置为try 
6. DHCP/DNS中高级设置关闭禁止解析 IPv6 DNS 记录.
