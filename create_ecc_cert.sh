ssh-keygen -m PEM -t ed25519 -C "ssh.ed25519" -N "" -f server
mv server server.pem
ssh-keygen -m PEM -t ed25519 -C "ssh.ed25519" -N "" -f client
mv client client.pem
ls -ahl *.pem *.pub
