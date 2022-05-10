#!/bin/bash
echo "Setting up byodns as daemon service. If you don't want to daemonize service, just run \"sudo \
  ~/go/bin/byodns\" after you ran install.sh. 
This is setup script for raspberry pi or any other linux pc.
If you are not on raspberry pi, and your username is not pi, YOU SHOULD MODIFY \
  default/byodns.service and default/byodns.conf (just change "pi" to your username)"
read -p "Continue? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1
echo "Setting up.."
echo "Installing built code"
go install

sudo cp default/byodns.service /lib/systemd/system/byodns.service

sudo systemctl daemon-reload
sudo systemctl enable redis-server
sudo systemctl restart byodns.service
sudo systemctl enable byodns.service

echo "Setup complete"
echo "log can be found in /etc/byodns/byodns.log && /etc/byodns/err.log"
echo "other files are in /etc/byodns/"
