# nginx2svg

Parse nginx configrue and generate a svg.

# Usage

Nginx2Svg will read nginx configure via `nginx -T`. So this topo graph is configure that nginx used. In some scenes, user will use other nginx configrue dir, so U can specify nginx command via `NGINX_COMMAND`. e.g.

```shell
export NGINX_COMMAND="/usr/local/bin/nginx -p /home/nginx"
```

Then nginx2svg will load configures for this nginx.

# Nginx

Now Nginx2Svg support load the below directives:

- server
- upsteam
- proxy_pass
- root

Nginx2Svg parse nginx configure via [znginx](https://github.com/andy-zhangtao/gogather/tree/master/tools/znginx). znginx is a nginx parse tool package.

# SVG

Nginx2Svg show topo graph via `d3js`.
