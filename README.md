# byodns

Let's build a dns ad blocker like pihole, in Golang.  
The concept is same as PiHole: you can block specific domain, and you can edit blocked domain list from web interface.    

WARNING: Since this project is meant to be used in a server, **UNIX** environment is assumed   
If you are on windows, why not install WSL2?  

## Frameworks / Libraries used
- redis-server: to cach dns query logs
- github.com/miekg/dns: provides dns interface to golang
- github.com/go-redis/redis: to interact with redis
- github.com/mattn/go-sqlite3: to store blocked domain list
- github.com/gorilla/mux: to provide web interface


## Installation

### For raspi user

Just copy&paste 3 command:)

1. `sudo su`
2. `apt update && apt install docker.io -y`
3. `docker run -d --name byodns -p 53:53/udp -p 80:80/tcp 0xsuk/byodns`
   After that `redis-server & byodns` will immediately start.  
   check application log using `docker exec byodns tail -f /etc/byodns/byodns.log`  
   Or check err log only using `docker exec byodns cat /etc/byodns/err.log`

### For those who don't want to use docker
Visit [Wiki](https://github.com/0xsuk/byodns/wiki)
1. Clone repo
   `git clone https://github.com/0xsuk/byodns.git` or  
   `git clone git@github.com:0xsuk/byodns.git` for ssh.
2. Install
   `bash install.sh`  
   it's going to `go build` and place default db/config file in /etc/byodns directory  
   When byodns prompt Username/Groupname, type user/group name you want to be owner of /etc/byodns directory
3. run program as UNIX daemon (OPTIONAL)
   `bash daemon.sh`  
   This will start install byodns previously built and
   AND, start/daemonize byodns.service (so that server will immediately start even after machine rebooted)  
   NOTE: default/byodns.service is written for user "pi". If you are not user "pi", you have to modify default/byodns.service BEFORE `bash daemon.sh`, because what daemon.sh does is to copy default/byodns.service to /lib/systemd/system/byodns.service which is going to be referenced by systemd.

Check /var/log/byodns.log to confirm server running  
3. Or run manually if you don't want it to be daemonized
type `sudo ./byodns` in directory you cloned.

## Usage

- Start DNS: check [Installation](#Installation)
- Add domain to blacklist: under development
- Add domain to whitelist: udner development

## Tested on

Installation for raspi user

- Linux pi 5.10.60-v7+ #1449 SMP Wed Aug 25 15:00:01 BST 2021 armv7l GNU/Linux (Raspbian/3b+)
  Installation for non-raspi user
- Kali linux 2020


## Configuration

There's a few things you can configure yourself (you don't have to though)

NOTE:everything should be "written in double quote".

- Upperstream DNS
  - IP: the ip of upperstream dns for when local dns (byodns) could not resolve domain
  - Port: the port of upperstream. Basically 53
- Local DNS
  - Port: the dns port to be listened to by byodns
- Web Server 
  - Port: the http port tobe listened to by byodns





## Todo

- [x] Implement really basic dns feature
- [x] Employ blacklist based blocking
- [x] Caching
- [x] Query Log
- [x] Web Interface
