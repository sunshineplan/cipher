#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
    apt -qq -y -t $(lsb_release -sc)-backports install golang-go
}

installSTE() {
    curl -Lo- https://github.com/sunshineplan/ste-go/archive/v1.0.tar.gz | tar zxC /var/www
    mv /var/www/ste-go* /var/www/ste-go
    cd /var/www/ste-go
    go build
}

configSTE() {
    read -p 'Please enter unix socket(default: /run/ste-go.sock): ' unix
    [ -z $unix ] && unix=/run/ste-go.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/ste-go.log): ' log
    [ -z $log ] && log=/var/log/app/ste-go.log
    mkdir -p $(dirname $log)
    sed "s,\$unix,$unix," /var/www/ste-go/config.ini.default > /var/www/ste-go/config.ini
    sed -i "s,\$log,$log," /var/www/ste-go/config.ini
    sed -i "s/\$host/$host/" /var/www/ste-go/config.ini
    sed -i "s/\$port/$port/" /var/www/ste-go/config.ini
}

setupsystemd() {
    cp -s /var/www/ste-go/scripts/ste-go.service /etc/systemd/system
    systemctl enable ste-go
    service ste-go start
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
    cp -s /var/www/ste-go/scripts/ste-go.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /var/www/ste-go/scripts/ste-go.conf
    sed -i "s,\$unix,$unix," /var/www/ste-go/scripts/ste-go.conf
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
