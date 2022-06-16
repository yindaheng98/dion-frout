docker run --rm -e HTTP_PROXY="http://192.168.1.2:10801" -e HTTPS_PROXY="http://192.168.1.2:10801" -v "$(pwd):/dion-frout" golang:1.18-buster sh -c 'apt update && apt install -y make && cd /dion-frout/aliyun && make all'
scp aliyun/bin/* root@8.131.50.230:/bin
scp aliyun/conf/* root@8.131.50.230:/root
