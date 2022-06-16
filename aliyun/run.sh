#!/bin/sh
chmod +x /bin/*
islb -c /root/islb.toml
stupid -conf /root/stupid.sfu.toml -filter "drawbox=x=0:y=0:w=50:h=50:c=blue"
isglb -c /root/islb.toml
sxu -c /root/beijing.sfu.toml -id beijing -filter "drawtext=text='beijing %{localtime\:%Y-%m-%d %H.%M.%S}':fontsize=60:x=(w-text_w)/2:y=0"
