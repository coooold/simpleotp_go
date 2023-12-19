# simpleotp_go
port of [newhouseb/simpleotp](https://github.com/newhouseb/simpleotp)


# Super Basic TOTP `auth_request` Server for nginx

## What is this for?

Have you ever wanted to add more security to a web application without modifying the web application itself? Take for example Jupter Notebook/Lab, which allows you to run arbitrary code from a web browser. It supports a built-in password / token-based authentication. Hopefully you're using a unique password, but if you're following proper security practices it's generally a good idea to protect stuff with "something you know and something you have." Chances are that if you've gotten this far you don't need me to convince you of the merits of two factor authentication.

## How does it work?

I use nginx in front of a variety of web services to handle SSL termination (using letsencrypt, which is amazing and you should also use). Nginx has a handy module called auth_request that you can use to specify an endpoint to check if a user is authenticated. If the endpoint returns 200, the parent request is allowed to succeed, otherwise a 401 error is returned. You can set up nginx to then redirect the user to a login page where they can do whatever they need to assert proof of identity.

In this case, the auth endpoint is reverse proxied to the simple script in this repo, which does things like token checking and presenting a login form.

## Example Configuration

In something like `/etc/nginx/sites-enabled/default`

```
server {
        server_name jupyter.example.com;

        location /totp/login {
                proxy_pass http://127.0.0.1:8000; # This is the TOTP Server
                proxy_set_header X-Original-URI $request_uri;
        }

        # This ensures that if the TOTP server returns 401 we redirect to login
        error_page 401 = @error401;
        location @error401 {
            return 302 /totp/login;
        }

        location / {
                auth_request /totp/check;
                proxy_pass http://127.0.0.1:8888; # This is Jupyter

                # This is needed for Jupyter to proxy websockets correctly,
                # it's unrelated to auth but handy to have written down here
                # for reference anyhow...
                proxy_http_version 1.1;
                proxy_set_header Upgrade $http_upgrade;
                proxy_set_header Connection $connection_upgrade;
        }

# The rest of the server definition, including SSL and whatnot
```

## Run TOTP backend
simpleotp_go --help

  -cookie string
        cookie name
  -port int
        listen port default 8000
  -secret string
        TOTP secret key


Configurations are also loaded from env.

1. SECRET_KEY
2. COOKIE_NAME
3. PORT

```

## FAQ

**Wait, this checks the TOTP secret before you enter a password?**

Yep, it feels kinda backwards, but I only have one login anyhow and I've rate-limited TOTP checks, so you can't hammer auth to figure out the TOTP secret.

