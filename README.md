# Sparkles
Serveur d'autentification pour Oauth pour Nginx

Configuation de la zone à protéger :
```
server {
        listen 80;
        listen [::]:80;

        server_name dev.skript-mc.fr;
        root /var/www/dev.skript-mc.fr;
        index index.html index.php;
        error_log /var/log/nginx/dev.skript-mc.fr.log warn;


        location / {
       		auth_request /validate;
    		error_page 401 = @error401;
               	try_files  $uri $uri/ /index.php?$args;
        }

       	location = /validate {
        	internal;
          	proxy_set_header Host $host;
          	proxy_pass_request_body off;
          	proxy_set_header Content-Length "";

          	proxy_pass http://127.0.0.1:3000;
        }

	location @error401 {
        	return 302 https://sparkles.skript-mc.fr/login;
    	}

}
```

Configuration du reverse proxy pour l'autentification :
```
server {
    listen 80;
    listen [::]:80;

    server_name sparkles.skript-mc.fr;
    error_log /var/log/nginx/sparkles.skript-mc.fr.log warn;

    location / {
    	proxy_set_header Host $http_host;
        proxy_pass http://127.0.0.1:3000;
    }

}
```
