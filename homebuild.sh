CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ovs-manager .
ssh root@192.168.1.32 "rm -rf /ovs-manager"
scp ovs-manager root@192.168.1.32:/