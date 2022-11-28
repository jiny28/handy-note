FROM openjdk:8-jre

MAINTAINER "jiny"

ADD *.jar /opt/jar/app.jar
EXPOSE 8761
WORKDIR /opt/jar
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
ENTRYPOINT ["java","-jar","app.jar"]