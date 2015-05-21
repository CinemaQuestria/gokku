gokku
=====

Gogs commit hook -> dokku deploy

Usage
-----

Build a docker image off of this repo. Next make a new repo for your 
configuration.

```Dockerfile
FROM cinemaquestria/gokku

ADD gokku.key /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
ADD known_hosts /root/.ssh/known_hosts

RUN mkdir -p /gokku/repo
```

Next you need to set your environment variables. As we are mainly on PonyChat, 
gokku is currently hardcoded to connect to `irc.ponychat.net`.

```Dockerfile
ENV BOT_CHANNEL commits      # An IRC channel name without the leading hash
ENV BOT_PASS foobang barbaz  # NickServ credentials

# What to git pull from
ENV GOKKU_REPO git@github.com:cinemaquestria/site
# Where to git push to
ENV GOKKU_DOKKU_REMOTE dokku@cq.internal:cinemaquestria
```
