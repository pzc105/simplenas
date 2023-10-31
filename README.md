## 特点
  * 一键部署，注册登录即可使用
  * 后台BT下载，下载越久下载越快
  * 高性能文件目录管理，随时随地创建删除
  * 自动识别BT下载的视频文件，可任意添加到你创建的目录中
  * 支持HLS传输协议，随心拖动时间条，播放视频无卡顿
  * 还有支持多音轨视频、在线聊天室、分享你创建的目录文件、记录播放时间等功能哦
## Docker
  * 创建一个目录A, 接着在这个目录下创建2个证书私钥（http.crt http.key和rpc.http rpc.key）
  * sudo python3 docker.py --crt-path \<path to A\> --rpc-server \<rpc server address\> --build-myenv --build-sn --nvidia-init -r -n sn -i