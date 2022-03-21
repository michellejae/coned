#!/bin/bash
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm coned
echo "and scene! (on local machhine)"

echo "build local binary on ... the local!"
GOOS=linux GOARCH=amd64 go build

echo "connect to droplet and delete current coned binaries as well as data, css files, images and views"
ssh root@want2breakfree.com "rm -rf /var/www/go"
echo "done deleting"

echo "connect to droplet and make folders"
ssh root@want2breakfree.com "mkdir /var/www/go; mkdir /var/www/go/data; mkdir -p /var/www/go/public/css; mkdir /var/www/go/public/images; mkdir -p /var/www/go/views/layouts"


echo "send coned file to server"
scp coned root@want2breakfree.com:/var/www/go/
echo "coned file is up in space"

echo "send data files up up and away"
scp ./data/* root@want2breakfree.com:/var/www/go/data/
echo "done"

echo "send css files up up and away"
scp ./public/css/* root@want2breakfree.com:/var/www/go/public/css/
echo "done"

echo "send images files up up and away"
scp ./public/images/* root@want2breakfree.com:/var/www/go/public/images/
echo "done"

echo "send html files up up and away"
scp ./views/*.html root@want2breakfree.com:/var/www/go/views/
echo "done"

echo "send html files up up and away"
scp ./views/layouts/*.html root@want2breakfree.com:/var/www/go/views/layouts/
echo "done"

echo "  Restarting the server..."
ssh root@want2breakfree.com "sudo service coned restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@want2breakfree.com "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing coned ===="