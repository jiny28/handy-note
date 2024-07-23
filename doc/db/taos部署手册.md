# taos使用手册
## java 服务客户端部署
> 核心逻辑：打包 java 镜像时把客户端所需的驱动安装至镜像里面;如果遇到没网络情况请先下载好文件放置根目录并修改``dockerfile``文件
1. ``Dockerfile``文件如下:
```dockerfile
FROM openjdk:8-jre

MAINTAINER "jiny"
EXPOSE 8888
ENV TDENGINE_VERSION=2.6.0.6
RUN wget -c https://www.taosdata.com/assets-download/TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz \
   && tar xvf TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz \
   && cd TDengine-client-${TDENGINE_VERSION} \
   && ./install_client.sh \
   && cd ../ \
   && rm -rf TDengine-client-${TDENGINE_VERSION}-Linux-x64.tar.gz TDengine-client-${TDENGINE_VERSION}
ADD *.jar /opt/jar/app.jar
WORKDIR /opt/jar
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' >/etc/timezone
ENTRYPOINT ["java","-jar","app.jar"]
```
2. build 镜像
3. 执行该镜像的run命令，注意：执行镜像时需要添加参数如下
```
# 该参数修改容器hosts映射，客户端需要配置FQDN，host设置映射到容器id或者是容器hostname
# 不使用FQDN方式访问，有一些API会直接异常
--add-host=taos-server:10.88.0.14
```

## taos 服务端部署
### Docker 单机部署
1. 执行``Docker``命令：
```
docker run -d --name taos --hostname="taos-server" -e TZ=Asia/Shanghai  -p 6030-6049:6030-6049 -p 6030-6049:6030-6049/udp tdengine/tdengine:2.6.0.6
```
### Docker 集群部署

1. 规划物理服务器 `ip` 和 `FQDN`

```dockerfile
- "c1.taosdata.com:10.88.0.36"
- "c2.taosdata.com:10.88.0.38"
- "c3.taosdata.com:10.88.0.47"
- "c4.taosdata.com:10.88.0.39"
```

2. 准备`docker-compose.yml` 文件

**文件更改说明**

每台物理服务器都需要启动该容器，并且镜像版本需一致，每台需差异化更改内容如下：

a. `TAOS_FQDN` 每台修改为规划好的`FQDN`

b. `TAOS_FIRST_EP` 每台一样，需统一填写集群中第一台启动的节点（管理节点）

c. `TAOS_SECOND_EP` 每台一样，在集群中任选一台即可，相当于备用管理节点，形成 end point 的高可用

d. `extra_hosts` 每台一样，修改为规划好的`FQDN`的内容即可，注意：部署时需要把本机`ip`的host给注释掉

```dockerfile
version: "3.5"

services:
  taos:
    image: tdengine/tdengine:3.0.3.1
    # hostname: taos-server
    container_name: taos
    ports:
      - 6030-6049:6030-6049
      - 6030-6049:6030-6049/udp
    environment:
      TZ: Asia/Shanghai
      TAOS_FQDN: "c1.taosdata.com"
      TAOS_FIRST_EP: "c1.taosdata.com:6030"
      TAOS_SECOND_EP: "c2.taosdata.com:6030"
    extra_hosts:
      - "c1.taosdata.com:10.88.0.36"
      - "c2.taosdata.com:10.88.0.38"
      - "c3.taosdata.com:10.88.0.47"
      - "c4.taosdata.com:10.88.0.39"
    restart: always
    volumes:
      - ./taos/data:/var/lib/taos/
      - ./taos/log:/var/log/taos/
      
```

3. 启动`firstep`配置的服务器内的容器

4. 进去该`taos`容器的客户端执行，看到如下信息说明启动成功

```
taos> show dnodes;
     id|    endpoint     | vnodes | support_vnodes |   status   |  create_time   | reboot_time   |     note      |
========================================================================================================================
     1 | c1.taosdata.com:6030  |      0 |   8 | ready | 2023-03-23 18:03:04.013 | 2023-03-23 18:02:52.399 |     |
Query OK, 1 row(s) in set (0.005001s)
```

5. 依次启动剩余的节点后并在`firstep`节点`taos`客户端执行添加节点操作

   > 如果按first、second顺序启动，启动完后可以show dnodes;看一眼，大多数时候会自动加入。
   >
   > 集群创建完后可以show mnodes;看一眼有多少管理节点，如果只有一个，并且又有高可用的需求，需要多添加几个管理节点，添加语法为:create mnode on dnode [dnodeid]

```
taos> CREATE DNODE "c2.taosdata.com";
taos> CREATE DNODE "c3.taosdata.com";
taos> CREATE DNODE "c4.taosdata.com";
taos> show dnodes;
id |   endpoint       | vnodes | support_vnodes |   status   |       create_time       |  reboot_time  | note    |
==========================================================================================================================
1 | c1.taosdata.com:6030 |      0 |       8 | ready      | 2023-03-23 18:03:04.013 | 2023-03-23 18:02:52.399 | |
2 | c2.taosdata.com:6030 |      0 |       8 | ready      | 2023-03-23 18:55:41.444 | 2023-03-23 19:01:24.058 | |
3 | c3.taosdata.com:6030 |      0 |       8 | ready      | 2023-03-23 19:06:14.952 | 2023-03-23 19:13:05.931 | |
4 | c4.taosdata.com:6030 |      0 |       8 | ready      | 2023-03-23 19:10:25.189 | 2023-03-23 19:13:52.209 | |
Query OK, 4 row(s) in set (0.003706s)
```

6. `java`客户端访问方式

a. 原生方式访问：访问如果需要高可用（也可以按URL方式访问），则可以把`JDBC`的`IP`和`PORT`填写为空后，把客户端配置文件内容中的`firstEp`和`secondEp`做相应的更改

```
############### 1. Cluster End point ############################

# The end point of the first dnode in the cluster to be connected to when this dnode or a CLI `taos` is started
firstEp                   c1.taosdata.com:6030

# The end point of the second dnode to be connected to if the firstEp is not available
secondEp                  c2.taosdata.com:6030
```

b. `REST`连接方式：直接访问`nginx` ，形成集群高可用访问，`nginx`监听`6041`与`6044 udp`端口，该端口都是`taos`默认端口，注意`nginx`所在服务器需要配置`taos`访问的`host`。

```
############### nginx point conf ################
http {
	upstream cli {
        server c1.taosdata.com:6041;
        server c2.taosdata.com:6041;
        server c3.taosdata.com:6041;
    }
    server{
        listen 6041;
        location /{
            proxy_pass http://cli;
        }
    }
}
stream {
    upstream scli {
        server c1.taosdata.com:6044;
        server c2.taosdata.com:6044;
        server c3.taosdata.com:6044;
    }
    server {
        listen 6044 udp;
        proxy_pass scli;
    }
}
```


