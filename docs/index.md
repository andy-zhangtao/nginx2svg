# 如何动态配置Nginx参数

Nginx参数众多，并且配置是非灵活，因此要达到完美的自动化配置是一件很有挑战性的事情，这个工具并不能十分完美的自动化调整参数。目前能自动化修改的参数仅有:

- server
- upstream
- proxy_pass
- root

下面将介绍`Nginx2Svg`是如何实现自动化修改参数的。

### 预备知识

为了更好的理解`Nginx2Svg`，需要一些很简单的预备知识。 首先需要了解Nginx的配置文件格式，一个典型的Nginx配置文件(假设此处Nginx作为7层反向负载使用)看起来应该是下面的样子:

```nginx
# 抄自nginx官网 http://nginx.org/en/docs/example.html
     1	user  www www;
     2
     3	worker_processes  2;
     4
     5	pid /var/run/nginx.pid;
     6
     7	#                          [ debug | info | notice | warn | error | crit ]
     8
     9	error_log  /var/log/nginx.error_log  info;
    10
    11	events {
    12	    worker_connections   2000;
    13
    14	    # use [ kqueue | epoll | /dev/poll | select | poll ];
    15	    use kqueue;
    16	}
    17
    18	http {
    19
    20	    include       conf/mime.types;
    21	    default_type  application/octet-stream;
    22
    23
    24	    log_format main      '$remote_addr - $remote_user [$time_local] '
    25	                         '"$request" $status $bytes_sent '
    26	                         '"$http_referer" "$http_user_agent" '
    27	                         '"$gzip_ratio"';
    28
    29	    log_format download  '$remote_addr - $remote_user [$time_local] '
    30	                         '"$request" $status $bytes_sent '
    31	                         '"$http_referer" "$http_user_agent" '
    32	                         '"$http_range" "$sent_http_content_range"';
    33
    34	    client_header_timeout  3m;
    35	    client_body_timeout    3m;
    36	    send_timeout           3m;
    37
    38	    client_header_buffer_size    1k;
    39	    large_client_header_buffers  4 4k;
    40
    41	    gzip on;
    42	    gzip_min_length  1100;
    43	    gzip_buffers     4 8k;
    44	    gzip_types       text/plain;
    45
    46	    output_buffers   1 32k;
    47	    postpone_output  1460;
    48
    49	    sendfile         on;
    50	    tcp_nopush       on;
    51	    tcp_nodelay      on;
    52	    send_lowat       12000;
    53
    54	    keepalive_timeout  75 20;
    55
    56	    #lingering_time     30;
    57	    #lingering_timeout  10;
    58	    #reset_timedout_connection  on;
    59
    60
    61	    server {
    62	        listen        one.example.com;
    63	        server_name   one.example.com  www.one.example.com;
    64
    65	        access_log   /var/log/nginx.access_log  main;
    66
    67	        location / {
    68	            proxy_pass         http://127.0.0.1/;
    69	            proxy_redirect     off;
    70
    71	            proxy_set_header   Host             $host;
    72	            proxy_set_header   X-Real-IP        $remote_addr;
    73	            #proxy_set_header  X-Forwarded-For  $proxy_add_x_forwarded_for;
    74
    75	            client_max_body_size       10m;
    76	            client_body_buffer_size    128k;
    77
    78	            client_body_temp_path      /var/nginx/client_body_temp;
    79
    80	            proxy_connect_timeout      70;
    81	            proxy_send_timeout         90;
    82	            proxy_read_timeout         90;
    83	            proxy_send_lowat           12000;
    84
    85	            proxy_buffer_size          4k;
    86	            proxy_buffers              4 32k;
    87	            proxy_busy_buffers_size    64k;
    88	            proxy_temp_file_write_size 64k;
    89
    90	            proxy_temp_path            /var/nginx/proxy_temp;
    91
    92	            charset  koi8-r;
    93	        }
    94
    95	        error_page  404  /404.html;
    96
    97	        location = /404.html {
    98	            root  /spool/www;
    99	        }
   100
   101	        location /old_stuff/ {
   102	            rewrite   ^/old_stuff/(.*)$  /new_stuff/$1  permanent;
   103	        }
   104
   105	        location /download/ {
   106
   107	            valid_referers  none  blocked  server_names  *.example.com;
   108
   109	            if ($invalid_referer) {
   110	                #rewrite   ^/   http://www.example.com/;
   111	                return   403;
   112	            }
   113
   114	            #rewrite_log  on;
   115
   116	            # rewrite /download/*/mp3/*.any_ext to /download/*/mp3/*.mp3
   117	            rewrite ^/(download/.*)/mp3/(.*)\..*$
   118	                    /$1/mp3/$2.mp3                   break;
   119
   120	            root         /spool/www;
   121	            #autoindex    on;
   122	            access_log   /var/log/nginx-download.access_log  download;
   123	        }
   124
   125	        location ~* \.(jpg|jpeg|gif)$ {
   126	            root         /spool/www;
   127	            access_log   off;
   128	            expires      30d;
   129	        }
   130	    }
   131	}
```

从18行到131行属于`http`配置内容，在这部分参数中，第61行到130行属于`server`配置内容，(一个server对应一个虚拟主机)，`server`的参数属于`http`参数的子集，当相同参数出现时，`server`优先级会高于`http`。按照作用域来做类比，`http`就是全局变量，`server`就是局部变量。

所以18行到60行属于全局变量,而61行到130则属于局部变量。 为了简化后面的操作，我们可以简化`http`和`server`之间的包含关系，如下:

```nginx
     1	user  nginx;
     2	worker_processes  1;
     3
     4	error_log  /var/log/nginx/error.log warn;
     5	pid        /var/run/nginx.pid;
     6
     7
     8	events {
     9	    worker_connections  1024;
    10	}
    11
    12
    13	http {
    15	    include       /etc/nginx/mime.types;
    16	    default_type  application/octet-stream;
    17
    18	    log_format main      '$remote_addr - $remote_user [$time_local] '
    19	                         '"$request" $status $bytes_sent '
    20	                         '"$http_referer" "$http_user_agent" '
    21	                         '"$gzip_ratio"';
    22
    23	    log_format download  '$remote_addr - $remote_user [$time_local] '
    24	                         '"$request" $status $bytes_sent '
    25	                         '"$http_referer" "$http_user_agent" '
    26	                         '"$http_range" "$sent_http_content_range"';
    27
    28	    access_log  /var/log/nginx/access.log  main;
    29
    30	    sendfile        on;
    31
    32	    keepalive_timeout  65;
    33
    34
    35	    server {
    36	        listen  80  default_server;
    37	        server_name  _;
    38
    39	        location /status {
    40	            vhost_traffic_status_display;
    41	            vhost_traffic_status_display_format html;
    42	        }
    43	    }
    44
    45	    include /etc/nginx/conf.d/*.conf;
    46	}
```

通过`include`引入其它server配置文件，而上面的内容可以作为`nginx.conf`全局默认配置文件，基本就不再修改了。而以后我们所要动态修改的配置文件就是`/etc/nginx/conf.d/*.conf`这部分。
