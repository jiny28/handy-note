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
