# docker
记录一些快速使用部署的脚本

# 注意事项
- emqx_host 命名方式不能使用 aaa 这种名称，这种命名方式集群无法通信，可以使用 ip 或者域名的方式
- swarm stack 启动方式3.0以上的版本容器顺序无法使用 depends_on 参数控制，需使用外挂 sh 脚本实现，官网推荐 [wait-for-it](https://github.com/vishnubob/wait-for-it) 项目实现。
- environment 参数的设置见[官网]( https://www.emqx.io/docs/zh/v4.4/configuration/configuration.html#cluster )文档。
- swarm 使用 nfs 共享 volume 实现 nginx 配置文件共享
- volume 和 mount 挂载逻辑不一致，volume 宿主机目录为空时，以容器目录为准，而 mount 正相反