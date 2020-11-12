#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installSTE() {
    curl -Lo- https://github.com/sunshineplan/cipher/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/cipher* /var/www/cipher
    cd /var/www/cipher/webapp
    go build -ldflags "-s -w" -o ste
}

configSTE() {
    read -p 'Please enter unix socket(default: /run/ste.sock): ' unix
    [ -z $unix ] && unix=/run/ste.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/ste.log): ' log
    [ -z $log ] && log=/var/log/app/ste.log
    mkdir -p $(dirname $log)
    sed "s,\$unix,$unix," /var/www/cipher/webapp/config.ini.default > /var/www/cipher/webapp/config.ini
    sed -i "s,\$log,$log," /var/www/cipher/webapp/config.ini
    sed -i "s/\$host/$host/" /var/www/cipher/webapp/config.ini
    sed -i "s/\$port/$port/" /var/www/cipher/webapp/config.ini
}

setupsystemd() {
    cp -s /var/www/cipher/webapp/scripts/ste.service /etc/systemd/system
    systemctl enable ste
    service ste start
}

writeLogrotateScrip() {
    if [ ! -f '/etc/logrotate.d/app' ]; then
	cat >/etc/logrotate.d/app <<-EOF
		/var/log/app/*.log {
		    copytruncate
		    rotate 12
		    compress
		    delaycompress
		    missingok
		    notifempty
		}
		EOF
    fi
}

setupNGINX() {
    cp -s /var/www/cipher/webapp/scripts/ste.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/cipher/webapp/scripts/ste.conf
    sed -i "s,\$unix,$unix," /var/www/cipher/webapp/scripts/ste.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installSTE
    configSTE
    setupsystemd
    writeLogrotateScrip
    setupNGINX
}

main
