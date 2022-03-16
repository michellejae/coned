#!/bin/bash
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm coned
echo "and scene! (on local machhine)"
echo "build local binary on ... the local!"
GOOS=linux GOARCH=amd64 go build

echo "connect to droplet and delete current coned binaries as well as data, css files, images and views"
ssh root@want2breakfree.com "rm /var/www/go/coned; rm -f /var/www/go/public/css/*; rm -f /var/www/go/public/images/*; rm -f /var/www/go/data/*; rm -f /var/www/go/view/*"
echo "done deleting"

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
scp ./view/* root@want2breakfree.com:/var/www/go/view/
echo "done"

echo "  Restarting the server..."
ssh root@want2breakfree.com "sudo service coned restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@want2breakfree.com "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing coned ===="