# ssl 证书配置 
## 安装说明
见github说明,地址：https://github.com/acmesh-official/acme.sh

## 更新证书

### 下载最新证书
```shell
// 两个域名
acme.sh --issue -d henrywaltz.cn --standalone 
acme.sh --issue -d www.henrywaltz.cn --standalone
```
### 证书迁移至指定目录
```shell
// 两个域名
// 需先删除目标目录的原文件
acme.sh --install-cert -d henrywaltz.cn \
--key-file       /etc/nginx/ssl/key.pem  \
--fullchain-file /etc/nginx/ssl/cert.pem \
--reloadcmd     "service nginx force-reload"

acme.sh --install-cert -d www.henrywaltz.cn \
--key-file       /etc/nginx/sslwww/key.pem  \
--fullchain-file /etc/nginx/sslwww/cert.pem \
--reloadcmd     "service nginx force-reload"
```

## nginx 配置参考
```conf
upstream php {
    server unix:/run/php/php7.4-fpm.sock;  #使用unix socket
    #server 127.0.0.1:9000;  #使用TCP端口
    # 上面两个取决于php-fpm的配置，二选一即可
}


server {
    listen       443 ssl http2;

    root /var/www/html/henrywaltzsite/;

    # Add index.php to the list if you are using PHP
    index index.php index.html;

    server_name henrywaltz.cn;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    ssl on;
    ssl_certificate /etc/nginx/ssl/cert.pem;     # 指定ssl证书路径
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_dhparam /etc/nginx/ssl/dhparam.pem;
    ssl_session_cache    shared:SSL:15m;
    ssl_session_timeout  30m;
    ssl_protocols       TLSv1 TLSv1.1 TLSv1.2;
    ssl_ciphers   ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA:ECDHE-RSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-RSA-AES256-SHA256:DHE-RSA-AES256-SHA:ECDHE-ECDSA-DES-CBC3-SHA:ECDHE-RSA-DES-CBC3-SHA:EDH-RSA-DES-CBC3-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:DES-CBC3-SHA:!DSS;
    ssl_prefer_server_ciphers  on;
    client_max_body_size 64m;
    location / {
        # First attempt to serve request as file, then
        # as directory, then fall back to displaying a 404.
        # try_files $uri $uri/ /index.php?$args?;
        try_files $uri $uri/ /index.html;
    }

    location ~ \.php$ {
        # NOTE: You should have "cgi.fix_pathinfo = 0;" in php.ini
        proxy_buffer_size          128k;
        proxy_buffers              4 256k;
        proxy_busy_buffers_size    256k;
        fastcgi_pass   php;
        fastcgi_index  index.php;
        fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
        include        fastcgi_params;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico)$ {
        expires max;
        log_not_found off;
    }

    location = /favicon.ico {
        log_not_found off;
        access_log off;
    }

    location = /robots.txt {
        allow all;
        log_not_found off;
        access_log off;
    }
}

server {
    listen       80;
    server_name henrywaltz.cn www.henrywaltz.cn;
#ACME_NGINX_START
    location ~ "^/\.well-known/acme-challenge/([-_a-zA-Z0-9]+)$" {
        default_type text/plain;
        return 200 "$1.Iy11835keRzzjRpwVARjVeC2dAoLfbz-kJptSPZY2tU";
    }
#NGINX_START
    return 301 https://henrywaltz.cn$request_uri;
}

server {
    listen       443 ssl http2;

    server_name www.henrywaltz.cn;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    ssl on;
    ssl_certificate /etc/nginx/sslwww/cert.pem;     # 指定ssl证书路径
    ssl_certificate_key /etc/nginx/sslwww/key.pem;

    ssl_dhparam /etc/nginx/ssl/dhparam.pem;
    ssl_session_cache    shared:SSL:15m;
    ssl_session_timeout  30m;
    ssl_protocols       TLSv1 TLSv1.1 TLSv1.2;
    ssl_ciphers   ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES256-SHA384:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA:ECDHE-RSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES128-SHA:DHE-RSA-AES256-SHA256:DHE-RSA-AES256-SHA:ECDHE-ECDSA-DES-CBC3-SHA:ECDHE-RSA-DES-CBC3-SHA:EDH-RSA-DES-CBC3-SHA:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:DES-CBC3-SHA:!DSS;
    ssl_prefer_server_ciphers  on;
    client_max_body_size 64m;

    location ~ "^/\.well-known/acme-challenge/([-_a-zA-Z0-9]+)$" {
        default_type text/plain;
        return 200 "$1.Iy11835keRzzjRpwVARjVeC2dAoLfbz-kJptSPZY2tU";
    }
#NGINX_START
    return 301 https://henrywaltz.cn$request_uri;
}
```