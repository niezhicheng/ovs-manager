CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ovs-manager .
ssh root@10.10.10.7 "rm -rf /ovs-manager"
scp ovs-manager root@10.10.10.7:/