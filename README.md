### Run

docker run -ti -p 69:69/udp -v $PWD:/data tftpd /bin/bash -c "cd /data && /go/bin/app" 
