.\ffmpeg.exe -f dshow -i video="@device_pnp_\\?\usb#vid_2bdf&pid_028a&mi_00#6&1d424522&0&0000#{65e8773d-8f56-11d0-a3b9-00a0c9223196}\global" -s 1280x720 -vcodec libvpx -g 24 -b:v 3M -f ivf output.ivf
.\ffplay.exe -f ivf -i output.ivf