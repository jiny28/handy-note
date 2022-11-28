# 直接从docker提供的远程镜像仓库中拉取最新的tomcat组件
FROM tomcat:8.5.40

# 本镜像的制作者的信息
MAINTAINER "jiny" 


#定义工作目录
ENV WORK_PATH /usr/local/tomcat
#定义要替换的文件名
ENV USER_CONF_FILE_NAME tomcat-users.xml
#定义要替换的server.xml文件名
ENV SERVER_CONF_FILE_NAME server.xml
#定义需要替换的context.xml文件名
ENV WEBAPP_CONTEXT_FILE_NAME context.xml
#定义context.xml的路径
ENV MANAGER_CONTEXT_PATH /usr/local/tomcat/webapps/manager/META-INF
ENV HOST_MANAGER_CONTEXT_PATH /usr/local/tomcat/webapps/host-manager/META-INF



#删除原文件tomcat-users.xml
RUN rm $WORK_PATH/conf/$USER_CONF_FILE_NAME && \
    rm $WORK_PATH/conf/$SERVER_CONF_FILE_NAME && \
    rm $MANAGER_CONTEXT_PATH/$WEBAPP_CONTEXT_FILE_NAME && \
    rm $HOST_MANAGER_CONTEXT_PATH/$WEBAPP_CONTEXT_FILE_NAME

#替换为新的文件
COPY ./conf/$SERVER_CONF_FILE_NAME $WORK_PATH/conf/
COPY ./conf/$USER_CONF_FILE_NAME $WORK_PATH/conf/
COPY ./conf/$WEBAPP_CONTEXT_FILE_NAME $MANAGER_CONTEXT_PATH/
COPY ./conf/$WEBAPP_CONTEXT_FILE_NAME $HOST_MANAGER_CONTEXT_PATH/
