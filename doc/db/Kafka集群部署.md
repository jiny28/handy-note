# Kafka 教程
## Kafka 集群部署
> 该教程仅适用于 ``kafka 3.1`` 以上的版本,且是使用的``kraft``模式。
1. 下载安装包
   - [教程版本](https://archive.apache.org/dist/kafka/3.1.0/kafka_2.13-3.1.0.tgz)
   - [官网下载](https://kafka.apache.org/downloads)
2. 安装包上传至集群中的每台服务器
3. 操作每台机器，并且解压安装包

```shell
tar -zxvf kafka_2.13-3.1.0.tgz
```
4. 跳转至配置目录

```shell
cd kafka_2.13-3.1.0/config/kraft/
```

5. 编写配置

```shell
vim server.properties
```

6. 以0.7、0.13、0.14服务器为例,分别的配置文件如下：
> 特别注意需要更改：log.dirs参数值，该值代表的是所有主题日志

<font color="red">0.7 服务器:</font>
```properties
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# This configuration file is intended for use in KRaft mode, where
# Apache ZooKeeper is not present.  See config/kraft/README.md for details.
#

############################# Server Basics #############################

# The role of this server. Setting this puts us in KRaft mode
process.roles=broker,controller

# The node id associated with this instance's roles
node.id=1
message.max.bytes=104857600
replica.fetch.max.bytes=104857600
# The connect string for the controller quorum
controller.quorum.voters=1@10.88.0.7:9093,2@10.88.0.14:9093,3@10.88.0.13:9093

############################# Socket Server Settings #############################

# The address the socket server listens on. It will get the value returned from
# java.net.InetAddress.getCanonicalHostName() if not configured.
#   FORMAT:
#     listeners = listener_name://host_name:port
#   EXAMPLE:
#     listeners = PLAINTEXT://your.host.name:9092
listeners=PLAINTEXT://10.88.0.7:9092,CONTROLLER://10.88.0.7:9093
inter.broker.listener.name=PLAINTEXT

# Hostname and port the broker will advertise to producers and consumers. If not set,
# it uses the value for "listeners" if configured.  Otherwise, it will use the value
# returned from java.net.InetAddress.getCanonicalHostName().
advertised.listeners=PLAINTEXT://10.88.0.7:9092

# Listener, host name, and port for the controller to advertise to the brokers. If
# this server is a controller, this listener must be configured.
controller.listener.names=CONTROLLER

# Maps listener names to security protocols, the default is for them to be the same. See the config documentation for more details
listener.security.protocol.map=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,SSL:SSL,SASL_PLAINTEXT:SASL_PLAINTEXT,SASL_SSL:SASL_SSL

# The number of threads that the server uses for receiving requests from the network and sending responses to the network
num.network.threads=3

# The number of threads that the server uses for processing requests, which may include disk I/O
num.io.threads=8

# The send buffer (SO_SNDBUF) used by the socket server
socket.send.buffer.bytes=102400

# The receive buffer (SO_RCVBUF) used by the socket server
socket.receive.buffer.bytes=102400

# The maximum size of a request that the socket server will accept (protection against OOM)
socket.request.max.bytes=104857600


############################# Log Basics #############################

# A comma separated list of directories under which to store log files
log.dirs=/home/kafka/kafka_2.13-3.1.0/log/kraft-combined-logs

# The default number of log partitions per topic. More partitions allow greater
# parallelism for consumption, but this will also result in more files across
# the brokers.
num.partitions=1

# The number of threads per data directory to be used for log recovery at startup and flushing at shutdown.
# This value is recommended to be increased for installations with data dirs located in RAID array.
num.recovery.threads.per.data.dir=1

############################# Internal Topic Settings  #############################
# The replication factor for the group metadata internal topics "__consumer_offsets" and "__transaction_state"
# For anything other than development testing, a value greater than 1 is recommended to ensure availability such as 3.
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1

############################# Log Flush Policy #############################

# Messages are immediately written to the filesystem but by default we only fsync() to sync
# the OS cache lazily. The following configurations control the flush of data to disk.
# There are a few important trade-offs here:
#    1. Durability: Unflushed data may be lost if you are not using replication.
#    2. Latency: Very large flush intervals may lead to latency spikes when the flush does occur as there will be a lot of data to flush.
#    3. Throughput: The flush is generally the most expensive operation, and a small flush interval may lead to excessive seeks.
# The settings below allow one to configure the flush policy to flush data after a period of time or
# every N messages (or both). This can be done globally and overridden on a per-topic basis.

# The number of messages to accept before forcing a flush of data to disk
#log.flush.interval.messages=10000

# The maximum amount of time a message can sit in a log before we force a flush
#log.flush.interval.ms=1000

############################# Log Retention Policy #############################

# The following configurations control the disposal of log segments. The policy can
# be set to delete segments after a period of time, or after a given size has accumulated.
# A segment will be deleted whenever *either* of these criteria are met. Deletion always happens
# from the end of the log.

# The minimum age of a log file to be eligible for deletion due to age
log.retention.hours=168

# A size-based retention policy for logs. Segments are pruned from the log unless the remaining
# segments drop below log.retention.bytes. Functions independently of log.retention.hours.
#log.retention.bytes=1073741824

# The maximum size of a log segment file. When this size is reached a new log segment will be created.
log.segment.bytes=1073741824

# The interval at which log segments are checked to see if they can be deleted according
# to the retention policies
log.retention.check.interval.ms=300000

```
<font color="red">0.13 服务器:</font>
```properties
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# This configuration file is intended for use in KRaft mode, where
# Apache ZooKeeper is not present.  See config/kraft/README.md for details.
#

############################# Server Basics #############################

# The role of this server. Setting this puts us in KRaft mode
process.roles=broker,controller

# The node id associated with this instance's roles
node.id=3
message.max.bytes=104857600
replica.fetch.max.bytes=104857600
# The connect string for the controller quorum
controller.quorum.voters=1@10.88.0.7:9093,2@10.88.0.14:9093,3@10.88.0.13:9093

############################# Socket Server Settings #############################

# The address the socket server listens on. It will get the value returned from
# java.net.InetAddress.getCanonicalHostName() if not configured.
#   FORMAT:
#     listeners = listener_name://host_name:port
#   EXAMPLE:
#     listeners = PLAINTEXT://your.host.name:9092
listeners=PLAINTEXT://10.88.0.13:9092,CONTROLLER://10.88.0.13:9093
inter.broker.listener.name=PLAINTEXT

# Hostname and port the broker will advertise to producers and consumers. If not set,
# it uses the value for "listeners" if configured.  Otherwise, it will use the value
# returned from java.net.InetAddress.getCanonicalHostName().
advertised.listeners=PLAINTEXT://10.88.0.13:9092

# Listener, host name, and port for the controller to advertise to the brokers. If
# this server is a controller, this listener must be configured.
controller.listener.names=CONTROLLER

# Maps listener names to security protocols, the default is for them to be the same. See the config documentation for more details
listener.security.protocol.map=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,SSL:SSL,SASL_PLAINTEXT:SASL_PLAINTEXT,SASL_SSL:SASL_SSL

# The number of threads that the server uses for receiving requests from the network and sending responses to the network
num.network.threads=3

# The number of threads that the server uses for processing requests, which may include disk I/O
num.io.threads=8

# The send buffer (SO_SNDBUF) used by the socket server
socket.send.buffer.bytes=102400

# The receive buffer (SO_RCVBUF) used by the socket server
socket.receive.buffer.bytes=102400

# The maximum size of a request that the socket server will accept (protection against OOM)
socket.request.max.bytes=104857600


############################# Log Basics #############################

# A comma separated list of directories under which to store log files
log.dirs=/home/kafka/kafka_2.13-3.1.0/log/kraft-combined-logs

# The default number of log partitions per topic. More partitions allow greater
# parallelism for consumption, but this will also result in more files across
# the brokers.
num.partitions=1

# The number of threads per data directory to be used for log recovery at startup and flushing at shutdown.
# This value is recommended to be increased for installations with data dirs located in RAID array.
num.recovery.threads.per.data.dir=1

############################# Internal Topic Settings  #############################
# The replication factor for the group metadata internal topics "__consumer_offsets" and "__transaction_state"
# For anything other than development testing, a value greater than 1 is recommended to ensure availability such as 3.
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1

############################# Log Flush Policy #############################

# Messages are immediately written to the filesystem but by default we only fsync() to sync
# the OS cache lazily. The following configurations control the flush of data to disk.
# There are a few important trade-offs here:
#    1. Durability: Unflushed data may be lost if you are not using replication.
#    2. Latency: Very large flush intervals may lead to latency spikes when the flush does occur as there will be a lot of data to flush.
#    3. Throughput: The flush is generally the most expensive operation, and a small flush interval may lead to excessive seeks.
# The settings below allow one to configure the flush policy to flush data after a period of time or
# every N messages (or both). This can be done globally and overridden on a per-topic basis.

# The number of messages to accept before forcing a flush of data to disk
#log.flush.interval.messages=10000

# The maximum amount of time a message can sit in a log before we force a flush
#log.flush.interval.ms=1000

############################# Log Retention Policy #############################

# The following configurations control the disposal of log segments. The policy can
# be set to delete segments after a period of time, or after a given size has accumulated.
# A segment will be deleted whenever *either* of these criteria are met. Deletion always happens
# from the end of the log.

# The minimum age of a log file to be eligible for deletion due to age
log.retention.hours=168

# A size-based retention policy for logs. Segments are pruned from the log unless the remaining
# segments drop below log.retention.bytes. Functions independently of log.retention.hours.
#log.retention.bytes=1073741824

# The maximum size of a log segment file. When this size is reached a new log segment will be created.
log.segment.bytes=1073741824

# The interval at which log segments are checked to see if they can be deleted according
# to the retention policies
log.retention.check.interval.ms=300000

```
<font color="red">0.14 服务器:</font>
```properties
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# This configuration file is intended for use in KRaft mode, where
# Apache ZooKeeper is not present.  See config/kraft/README.md for details.
#

############################# Server Basics #############################

# The role of this server. Setting this puts us in KRaft mode
process.roles=broker,controller

# The node id associated with this instance's roles
node.id=2
message.max.bytes=104857600
replica.fetch.max.bytes=104857600
# The connect string for the controller quorum
controller.quorum.voters=1@10.88.0.7:9093,2@10.88.0.14:9093,3@10.88.0.13:9093

############################# Socket Server Settings #############################

# The address the socket server listens on. It will get the value returned from
# java.net.InetAddress.getCanonicalHostName() if not configured.
#   FORMAT:
#     listeners = listener_name://host_name:port
#   EXAMPLE:
#     listeners = PLAINTEXT://your.host.name:9092
listeners=PLAINTEXT://10.88.0.14:9092,CONTROLLER://10.88.0.14:9093
inter.broker.listener.name=PLAINTEXT

# Hostname and port the broker will advertise to producers and consumers. If not set,
# it uses the value for "listeners" if configured.  Otherwise, it will use the value
# returned from java.net.InetAddress.getCanonicalHostName().
advertised.listeners=PLAINTEXT://10.88.0.14:9092

# Listener, host name, and port for the controller to advertise to the brokers. If
# this server is a controller, this listener must be configured.
controller.listener.names=CONTROLLER

# Maps listener names to security protocols, the default is for them to be the same. See the config documentation for more details
listener.security.protocol.map=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,SSL:SSL,SASL_PLAINTEXT:SASL_PLAINTEXT,SASL_SSL:SASL_SSL

# The number of threads that the server uses for receiving requests from the network and sending responses to the network
num.network.threads=3

# The number of threads that the server uses for processing requests, which may include disk I/O
num.io.threads=8

# The send buffer (SO_SNDBUF) used by the socket server
socket.send.buffer.bytes=102400

# The receive buffer (SO_RCVBUF) used by the socket server
socket.receive.buffer.bytes=102400

# The maximum size of a request that the socket server will accept (protection against OOM)
socket.request.max.bytes=104857600


############################# Log Basics #############################

# A comma separated list of directories under which to store log files
log.dirs=/home/kafka/kafka_2.13-3.1.0/log/kraft-combined-logs

# The default number of log partitions per topic. More partitions allow greater
# parallelism for consumption, but this will also result in more files across
# the brokers.
num.partitions=1

# The number of threads per data directory to be used for log recovery at startup and flushing at shutdown.
# This value is recommended to be increased for installations with data dirs located in RAID array.
num.recovery.threads.per.data.dir=1

############################# Internal Topic Settings  #############################
# The replication factor for the group metadata internal topics "__consumer_offsets" and "__transaction_state"
# For anything other than development testing, a value greater than 1 is recommended to ensure availability such as 3.
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1

############################# Log Flush Policy #############################

# Messages are immediately written to the filesystem but by default we only fsync() to sync
# the OS cache lazily. The following configurations control the flush of data to disk.
# There are a few important trade-offs here:
#    1. Durability: Unflushed data may be lost if you are not using replication.
#    2. Latency: Very large flush intervals may lead to latency spikes when the flush does occur as there will be a lot of data to flush.
#    3. Throughput: The flush is generally the most expensive operation, and a small flush interval may lead to excessive seeks.
# The settings below allow one to configure the flush policy to flush data after a period of time or
# every N messages (or both). This can be done globally and overridden on a per-topic basis.

# The number of messages to accept before forcing a flush of data to disk
#log.flush.interval.messages=10000

# The maximum amount of time a message can sit in a log before we force a flush
#log.flush.interval.ms=1000

############################# Log Retention Policy #############################

# The following configurations control the disposal of log segments. The policy can
# be set to delete segments after a period of time, or after a given size has accumulated.
# A segment will be deleted whenever *either* of these criteria are met. Deletion always happens
# from the end of the log.

# The minimum age of a log file to be eligible for deletion due to age
log.retention.hours=168

# A size-based retention policy for logs. Segments are pruned from the log unless the remaining
# segments drop below log.retention.bytes. Functions independently of log.retention.hours.
#log.retention.bytes=1073741824

# The maximum size of a log segment file. When this size is reached a new log segment will be created.
log.segment.bytes=1073741824

# The interval at which log segments are checked to see if they can be deleted according
# to the retention policies
log.retention.check.interval.ms=300000

```
7. 用``kafka-storage.sh`` 生成一个唯一的集群ID,只需要在集群中的某一台执行即可，但该台的``role``配置必须是``broker,controller``
```shell
cd bin/
./kafka-storage.sh random-uuid
# 复制生成的uuid
```
8. 用``kafka-storage.sh`` 格式化存储数据的目录，<font color = "red">每一个节点</font>都需要执行，执行时用上一步复制的``uuid``替换下面命令中的``uuid``
```shell
./kafka-storage.sh format -t uuid -c ../config/kraft/server.properties
```
9. 启动每个节点的``kafka``
```shell
./kafka-server-start.sh -daemon ../config/kraft/server.properties
```

## Kafka 服务端常用命令
1. 停止``kafka``
```shell
./kafka-server-stop.sh
```
2. 创建``topic``
```shell
./kafka-topics.sh --create --bootstrap-server 10.88.0.7:9092 --replication-factor 1 --partitions 1 --topic test
```
3. 列出所有topic
```shell
./kafka-topics.sh --list --bootstrap-server 10.88.0.7:9092 
```
4. 查看指定主题的详细信息
```shell
./kafka-topics.sh --describe --bootstrap-server 10.88.0.7:9092  --topic test 
```
5. 删除主题
```shell
./kafka-topics.sh --bootstrap-server 10.88.0.7:9092 --topic test --delete 
```
6. 订阅主题
```shell
./kafka-console-consumer.sh --bootstrap-server 10.88.0.7:9092 --topic test --offset latest --partition 0
```
7. 给主题发送消息
```shell
./kafka-console-producer.sh --broker-list 10.88.0.7:9092 --topic test 
```
## Java 客户端实现
1. 引入依赖
```xml
<dependency>
    <groupId>org.springframework.cloud</groupId>
    <artifactId>spring-cloud-stream-binder-kafka</artifactId>
    <version>3.0.6.RELEASE</version>
</dependency>
```
2. 配置``yml``文件
```yaml
spring:
  cloud:
    stream:
      kafka:
        binder:
          # 节点地址
          brokers: 10.88.0.7:9092,10.88.0.13:9092,10.88.0.14:9092
          producer-properties:
             acks: 1  # 收到server端确认次数
             retries: 1 # 发送失败重试次数，默认值0
             batch.size: 262144
             linger.ms: 0 # 0ms的延迟
             buffer.memory: 67108864
             max.request.size: 104857600 # 100M
             key.serializer: org.apache.kafka.common.serialization.StringSerializer
             value.serializer: org.apache.kafka.common.serialization.ByteArraySerializer
          consumer-properties:
             max.partition.fetch.bytes: 104857600 # 100M
             key.deserializer: org.apache.kafka.common.serialization.StringDeserializer
             value.deserializer: org.apache.kafka.common.serialization.ByteArrayDeserializer
             session.timeout.ms: 60000
        bindings:
           data-input:
              consumer:
                 auto-commit-offset: true # 是否自动提交，默认为true，
                 enable-dlq: true # 是否开启死信队列，默认为 false 关闭
                 dlq-name: # 死信队列名，默认为 `errors.{topicName}.{consumerGroup}`
      bindings:
        data-input: # 通道名称
          destination: datainput # 主题名称
          group: data-group # 分组
          consumer:
             max-attempts: 3 # 重试次数，默认为 3 次。
             back-off-initial-interval: 3000 # 重试间隔的初始值，单位毫秒，默认为 1000
             back-off-multiplier: 2.0 # 重试间隔的递乘系数，默认为 2.0
             back-off-max-interval: 10000 # 重试间隔的最大值，单位毫秒，默认为 10000
        data-output:
           destination: datainput # 主题名称
           content-type: application/json # 内容格式
           producer:
              partition-count: 2 # 几个partition
              partition-key-expression: payload.tableName # partition 根据什么key分组
```
3. 创建通道接口
```java
// input 接口
public interface DataInputInterface {

    String DATA = "data-input";//需要与配置文件中的通道名称一一对应
    @Input(DATA)
    SubscribableChannel input();
}
// output 接口
public interface DataOutputInterface {
   String DATA = "data-output";//需要与配置文件中的通道名称一一对应
   @Output(DataOutputInterface.DATA)
   MessageChannel output();
}
// 注入接口，启动类上添加如下注解
@EnableBinding({DataInputInterface.class,DataOutputInterface.class})
```
4. 业务中使用
发送消息:
```java
public class Send {
   @Resource
   private DataOutputInterface dataOutputInterface;

   public Boolean send(String msg) {
      return dataOutputInterface.output().send(MessageBuilder.withPayload(msg).build());
   }
}
```
接受消息:
```java
@Service
public class DataInput {

   @StreamListener(DataInputInterface.DATA)
   public void receive(Object o) {
      System.out.println("接收到消息:" + o.toString());
   }
}
```