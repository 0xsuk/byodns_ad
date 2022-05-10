echo "byodns"
echo "Don't type until you get prompted"
echo "Building code"
go build

echo "Checking redis-server installed"
redis=/etc/init.d/redis-server
if [[ ! -f  "$redis" ]]; then
  echo "Installing redis-server."
  echo "Password might be required"
  sudo apt install redis-server
  echo "installed redis server"
else
  echo "redis-server is already installed."
fi
sudo systemctl restart redis-server
echo "redis-server installed/started (not daemonized)."


echo "Creating required file in /etc/byodns..."
read -p "Enter Username: " user
read -p "Enter Group name: " group
DIR=/etc/byodns
sudo mkdir -p $DIR
sudo chown $user:$group $DIR
cp default/domains.db $DIR
cp default/config.ini $DIR
chmod 644 $DIR/*


echo "Created required file in /etc/byodns!"
echo "run \"bash daemon.sh\" if you want to daemonize process!"
echo "run \"sudo ./byodns\" if you want to start process!"
