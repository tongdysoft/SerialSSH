ssh-keygen -m PEM -t ed25519 -C "ssh.ed25519" -N "" -f server
RENAME server server.pem
ssh-keygen -m PEM -t ed25519 -C "ssh.ed25519" -N "" -f client
RENAME client client.pem
DIR *.pem *.pub
