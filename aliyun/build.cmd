# 在Win上编译
docker run --rm -e HTTP_PROXY="http://192.168.1.2:10801" -e HTTPS_PROXY="http://192.168.1.2:10801" -v "$(pwd):/dion-frout" golang:1.18-buster sh -c 'apt update && apt install -y make && cd /dion-frout/aliyun && make all'

# 从Win上传到Linux里
scp aliyun/bin/* root@8.131.50.230:/bin
scp aliyun/conf/* root@8.131.50.230:/root
scp aliyun/bin/* root@47.105.63.37:/bin
scp aliyun/conf/* root@47.105.63.37:/root
scp aliyun/bin/* root@47.122.1.50:/bin
scp aliyun/conf/* root@47.122.1.50:/root

# 在北京运行
docker run --name nats --network host -d nats
chmod +x /bin/*
islb -c /root/islb.toml
stupid -conf /root/stupid.sfu.toml -filter "drawbox=x=0:y=0:w=50:h=50:c=blue"
isglb -c /root/islb.toml
sxu -c /root/beijing.sfu.toml -id beijing -filter "drawtext=text='beijing %{localtime\:%Y-%m-%d %H.%M.%S}':fontsize=60:x=(w-text_w)/2:y=0"

# 在青岛运行
chmod +x /bin/*
sxu -c /root/qingdao.sfu.toml -id beijing -filter "drawtext=text='beijing %{localtime\:%Y-%m-%d %H.%M.%S}':fontsize=60:x=(w-text_w)/2:y=0"

# 在南京运行
chmod +x /bin/*
sxu -c /root/nanjing.sfu.toml -id beijing -filter "drawtext=text='beijing %{localtime\:%Y-%m-%d %H.%M.%S}':fontsize=60:x=(w-text_w)/2:y=0"
