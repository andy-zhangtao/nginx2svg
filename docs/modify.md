# 如何动态配置Nginx参数

Nginx参数众多，并且配置是非灵活，因此要达到完美的自动化配置是一件很有挑战性的事情，这个工具并不能十分完美的自动化调整参数。目前支持自动化修改的参数有:

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

### 配置规则

如果要达到自动化配置的目标，那么就需要设定一些规则。 下面是为了满足自动化而设置的规则：

+ 配置文件规则
    - 必须存在server_name。
    - 文件名以[server name].conf进行命名。 假设server_name为example.com, 则配置文件名就是example.com.conf。
    - 一个文件**有并且只有**一个server段
+ 配置内容规则
    - 同一个配置文件中location不重复(正则表达式不在限制范围内)

### 解析规则

在满足上述两个规则的前提下，我们来看如何实现Nginx参数的自动化配置。首先要明确实现nginx自动化配置的难点在哪里? 基于我的使用经验来看，难点在于以下三点:

+ nginx配置相当灵活，属于`非结构化`语义
    虽然nginx明确了配置文件的内容和格式，但在配置上可以任意组合(在执行nginx -t或者reload时才会真正验证)。因此配置文件只规定了最低门槛的`结构范式`，而并没有规定严谨的配置格式，造成了只要符合语义都可以验证成功。这一点在使用者眼里是非常灵活的优点，但从自动化角度来说则是很大的痛点，因为找不到一个统一的解析格式来理解语义。

+ 验证和回滚
    nginx是基于文本来进行配置的，每一次修改都是通过IO操作生成文本配置文件而后在加载在每个worker中。 因此当验证失败时，如何将新增/删除的内容恢复到上一个版本中，就变成了一个问题。

+ 个性化配置
    在真实业务场景中，nginx配置必然无法做到一个配置吃遍天。当某些server需要添加个性化配置参数时，如何平衡个性化配置和自动化配置，也变成了一个需要考虑的问题。

当找到上述三个问题的答案时，大体就可以满足自动化配置的要求了。

首先来看第一个问题。

如果因为nginx配置灵活而导致正面解析nginx配置文件是一个很困难的事情，那么可以尝试换个角度来理解这个问题。 **如果变化很多而不容易解析，那么就不要让它变化了**

具体怎么理解呢？ nginx是通过语义来验证的，也就是nginx自身其实对`结构`不敏感的(可以反向证明，如果nginx是依赖结构来理解配置的，那么它应该会规定严谨的配置结构)。所以我们可以事先定义好每个配置文件的配置格式，如下:

```nginx
     1
     2
     3  upstream 5d148ba37f325500011770af {
     4      server  xxxxx ;
     5  }
     6
     7
     8  server{
     9
    10    server_name web1.example.com;
    11
    12
    13
    14
    15    location /server1 {
    16      proxy_pass http://5d148ba37f325500011770af;
    17      proxy_set_header X-Real-IP $remote_addr;
    18      proxy_set_header Host $host;
    19      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    20      proxy_next_upstream error timeout http_500 http_502 http_503 http_504 non_idempotent;
    21
    22
    23
    24    }
    25
    26  }
    27
```

每个配置文件都规定好配置结构如下：

+ `upstream`都统一放置在`server`之前
+ `server_name`放置在`location`之前
+ `proxy_pass` 放置在每个`location`首行

当每个配置文件都满足上述三个条件时，自动化解析程序就可以按照设定好的规则解析并尝试理解每段语义。

只解析文件还不够，还需要能`动态修改`才可以。 再回到上面的配置内容，里面的变量有三部分，按照从上往下依次是：

1. upstream的server IP列表
2. server_name中的domain列表
3. location列表

动态修改更准确的就是如何动态修改上面三部分值，这三部分的关联关系如下：

```shell

    +-------------+
    | server_name |
    |   domain1   |
    |   domain2   |                 +-----------------+                 +-----------------+
    |   domain3   |---------------> |    location1    |-------------->  |   upstream1     |
    |   .......   |                 +-----------------+                 +-----------------+
    |   domainN   |
    +-------------+
                                    +-----------------+                 +-----------------+
                                    |    location2    |-------------->  |   upstream2     |
                                    +-----------------+                 +-----------------+


                                    +-----------------+                 +-----------------+
                                    |    locationN    |-------------->  |  upstreamN      |
                                    +-----------------+                 +-----------------+
```

同一个组的`server_name`共享所有的`location`数据，而每一个`location`则通过`proxy_pass`指向特定的`upstream`(可以是不同的，也可以是相同的upstream)。

从上图可以看出`server_name`和`location`在一个作用域中(在同一个`{}`中)而`upstream`则游离在外。

三个问题中，server_name可以通过`server_name`准确定位，`location`也可以准确定位，此时如何从`location`通过`proxy_pass`定位到`upstream`则变成了当前的难点。

在实际使用过程中，我通过添加`锚点`来解决这个问题，具体来说就是增加一组`upstream`辅助定位数据，例如下图中的数据:

```nginx
     1
     2  ### [5d148ba37f325500011770af]-[/]-[upstream]-[start]
     3  upstream 5d148ba37f325500011770af {
     4      server  xxxxx ;
     5  }
     6  ### [5d148ba37f325500011770af]-[/]-[upstream]-[end]
     7
     8  server{
     9
    10    server_name web1.example.com;
    11
    12
    13
    14
    15    location /server1 {
    16      proxy_pass http://5d148ba37f325500011770af;
    17      proxy_set_header X-Real-IP $remote_addr;
    18      proxy_set_header Host $host;
    19      proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    20      proxy_next_upstream error timeout http_500 http_502 http_503 http_504 non_idempotent;
    21
    22
    23
    24    }
    25
    26  }
    27
```

第二行和第六行就是添加的`锚点`。 锚点数据需要满足的条件是:

+ 同一个配置文件中不重复
+ 有良好的区分度

因此设计了上述的`锚点`数据，其格式如下:

```
    ### [5d148ba37f325500011770af]-[/]-[upstream]-[start]
    ----------------------------------------------------
    ### [24位随机数]-[/]-[upstream]-[开始/结束标示]
    ①       ②           ③             ④

    ① 三个#开头
    ② 满足锚点，upstream名称和proxy_pass一致，也就是第二行，第三行和第十六行使用同一个24位随机数
    ③ 固定格式,用来保证和其它注释信息不重复
    ④ start表示upstream开始， end表示upstream结束。
```

因此一个完整的自动化配置流程如下：

```shell
    // 假设配置web1.example.com的/server1 反向配置

    if web1.example.com.conf 存在
        逐行读取文件内容

        if 找到 server1的location行
            解析 proxy_pass，找到 24位随机数

            从头开始读取文件内容

            if 找到 ### [xxxx]-[/]-[upstream]-[start]
                找到锚点，此行往下两行是ip列表，开始修改
            else
                没找到锚点，配置文件出错，人工介入
        else
            // 当前没有此location配置，新建location和upstream
            新建location配置
            新建相匹配的upstream配置

    else
        // 当前没有此域名配置，新建一个
        创建 web1.example.com.conf，内容按照既定格式创建

```

### 个性化支持

从上面的解析规则来看，如果要支持个性化支持，那么在理解语义时要做到`适可而止`，也就是只需要解析到需要的数据就可以了，其它数据原样复制。例如用户在`location`中添加了个性化参数(需要满足`配置规则第三条`)，那么只要解析出`proxy_pass`就可以，后续的数据原样复制不要做变更。

