package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const nginx = `nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
nginx: configuration file /etc/nginx/nginx.conf test is successful
# configuration file /etc/nginx/nginx.conf:
user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;


events {
    worker_connections  1024;
}


http {
    vhost_traffic_status_zone;
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main 	'"$remote_addr" - "$remote_user" [$time_local] "$request" '
			            '"$status" "$body_bytes_sent" "$http_referer" '
			            '"$http_user_agent" "$http_x_forwarded_for" '
			            '"$cookie_JSESSIONID" "$host" "$upstream_addr"'
			            '"$upstream_status" "$upstream_response_time" '
			            '"$request_time" "$request_body" "$cookie__ver"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    server {
        listen  80  default_server;
        server_name  _;

        location /status {
            vhost_traffic_status_display;
            vhost_traffic_status_display_format html;
        }
    }

    include /etc/nginx/conf.d/*.conf;
}
# configuration file /etc/nginx/mime.types:

types {
    text/html                                        html htm shtml;
    text/css                                         css;
    text/xml                                         xml;
    image/gif                                        gif;
    image/jpeg                                       jpeg jpg;
    application/javascript                           js;
    application/atom+xml                             atom;
    application/rss+xml                              rss;

    text/mathml                                      mml;
    text/plain                                       txt;
    text/vnd.sun.j2me.app-descriptor                 jad;
    text/vnd.wap.wml                                 wml;
    text/x-component                                 htc;

    image/png                                        png;
    image/svg+xml                                    svg svgz;
    image/tiff                                       tif tiff;
    image/vnd.wap.wbmp                               wbmp;
    image/webp                                       webp;
    image/x-icon                                     ico;
    image/x-jng                                      jng;
    image/x-ms-bmp                                   bmp;

    application/font-woff                            woff;
    application/java-archive                         jar war ear;
    application/json                                 json;
    application/mac-binhex40                         hqx;
    application/msword                               doc;
    application/pdf                                  pdf;
    application/postscript                           ps eps ai;
    application/rtf                                  rtf;
    application/vnd.apple.mpegurl                    m3u8;
    application/vnd.google-earth.kml+xml             kml;
    application/vnd.google-earth.kmz                 kmz;
    application/vnd.ms-excel                         xls;
    application/vnd.ms-fontobject                    eot;
    application/vnd.ms-powerpoint                    ppt;
    application/vnd.oasis.opendocument.graphics      odg;
    application/vnd.oasis.opendocument.presentation  odp;
    application/vnd.oasis.opendocument.spreadsheet   ods;
    application/vnd.oasis.opendocument.text          odt;
    application/vnd.openxmlformats-officedocument.presentationml.presentation
                                                     pptx;
    application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
                                                     xlsx;
    application/vnd.openxmlformats-officedocument.wordprocessingml.document
                                                     docx;
    application/vnd.wap.wmlc                         wmlc;
    application/x-7z-compressed                      7z;
    application/x-cocoa                              cco;
    application/x-java-archive-diff                  jardiff;
    application/x-java-jnlp-file                     jnlp;
    application/x-makeself                           run;
    application/x-perl                               pl pm;
    application/x-pilot                              prc pdb;
    application/x-rar-compressed                     rar;
    application/x-redhat-package-manager             rpm;
    application/x-sea                                sea;
    application/x-shockwave-flash                    swf;
    application/x-stuffit                            sit;
    application/x-tcl                                tcl tk;
    application/x-x509-ca-cert                       der pem crt;
    application/x-xpinstall                          xpi;
    application/xhtml+xml                            xhtml;
    application/xspf+xml                             xspf;
    application/zip                                  zip;

    application/octet-stream                         bin exe dll;
    application/octet-stream                         deb;
    application/octet-stream                         dmg;
    application/octet-stream                         iso img;
    application/octet-stream                         msi msp msm;

    audio/midi                                       mid midi kar;
    audio/mpeg                                       mp3;
    audio/ogg                                        ogg;
    audio/x-m4a                                      m4a;
    audio/x-realaudio                                ra;

    video/3gpp                                       3gpp 3gp;
    video/mp2t                                       ts;
    video/mp4                                        mp4;
    video/mpeg                                       mpeg mpg;
    video/quicktime                                  mov;
    video/webm                                       webm;
    video/x-flv                                      flv;
    video/x-m4v                                      m4v;
    video/x-mng                                      mng;
    video/x-ms-asf                                   asx asf;
    video/x-ms-wmv                                   wmv;
    video/x-msvideo                                  avi;
}

# configuration file /etc/nginx/conf.d/alert.devex.yqxiu.cn.conf:

upstream alert{
    server  192.168.1.116:10000;
}

server{
	listen 443 ssl;
	server_name alert.devex.yqxiu.cn test.devex.eqshow.cn;

        ssl_certificate /etc/nginx/conf.d/ssl/1_alert.devex.yqxiu.cn_bundle.crt;
  	ssl_certificate_key /etc/nginx/conf.d/ssl/2_alert.devex.yqxiu.cn.key;

	ssl_protocols TLSv1.2;
  	ssl_prefer_server_ciphers on;
  	ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-SHA384;
  	ssl_ecdh_curve secp384r1; # Requires nginx >= 1.1.0
  	ssl_session_timeout 10m;
  	ssl_session_cache shared:SSL:10m;
  	ssl_session_tickets off; # Requires nginx >= 1.5.9

	location / {
		proxy_pass http://alert;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;



	}

}


# configuration file /etc/nginx/conf.d/max.yqxiu.cn.conf:

### [5bd963e20b5e55000af1aa73]-[/]-[upstream]-[start]
upstream 5bd963e20b5e55000af1aa73 {
    server  192.168.2.60:9000;
}
### [5bd963e20b5e55000af1aa73]-[/]-[upstream]-[end]

server{

	server_name max.yqxiu.cn;




	location / {
		proxy_pass http://5bd963e20b5e55000af1aa73;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;



	}

}


# configuration file /etc/nginx/conf.d/qa.tiny.eqshow.cn.conf:

### [5bdfe6df67609d000a5e5c4b]-[/]-[upstream]-[start]
upstream 5bdfe6df67609d000a5e5c4b {
    server  192.168.1.237:8000;
}
### [5bdfe6df67609d000a5e5c4b]-[/]-[upstream]-[end]

server{

	server_name qa.tiny.eqshow.cn;




	location / {
		proxy_pass http://5bdfe6df67609d000a5e5c4b;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;



	}

}


# configuration file /etc/nginx/conf.d/store.chinazt.com.conf:
upstream store.chinazt.com {
    server  111.231.176.40:80;
}

server{
        listen 80;
	server_name store.chinazt.com;

	location / {
		proxy_pass http://store.chinazt.com;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

	}

}

# configuration file /etc/nginx/conf.d/www.chinazt.com.conf:
upstream www.chinazt.com {
    server  212.64.45.174:80;
}

server{
        listen 80;
	server_name www.chinazt.com;

	location / {
		proxy_pass http://www.chinazt.com;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header Host $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

	}

}`

func Test_parseNginxContent(t *testing.T) {
	configure, otherinfo := parseNginxContent(nginx)
	assert.Equal(t, 1, len(otherinfo))
	assert.Equal(t, 7, len(configure))
}

func Test_extraceServer(t *testing.T) {
	configure, _ := parseNginxContent(nginx)
	server, location, upstream := extractServer(configure[2])

	assert.Equal(t, 2, len(server))
	assert.Equal(t, 2, len(location))
	assert.Equal(t, 1, len(upstream))

}

func Test_extractNginxMeta(t *testing.T) {
	meta, err := extractNginxMeta(nginx)
	assert.Nil(t, err)

	assert.Equal(t, 7, len(meta))
}
