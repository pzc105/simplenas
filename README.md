## 特点
  * 一键部署，注册登录即可使用
  * 支持多端同时登录
  * 支持爬虫爬取magnet uri
  * BT下载
  * 高性能文件目录管理
  * 在线观看，支持弹幕
  * 在线聊天
  * 目录分享、视频分享
  * 用户目录读写权限隔离
  * 一键整理视频名字，按集数排序
## Docker
  * 创建一个目录A, 接着在这个目录下创建2个证书私钥（http.crt http.key和rpc.http rpc.key）
  * sudo python3 docker.py --crt-path \<path to A\> --rpc-server \<rpc server address\> --build-myenv --build-sn --nvidia-init -r -n sn -i