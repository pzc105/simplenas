## 项目前端为React工程，后端为C++和Golang。
  * 通过session管理登录，支持多端同时登录。
  * 支持BT下载。
  * 支持自定义文件夹。
  * 支持将BT视频添加到自定义文件夹，视频转成HLS传输格式，H264编码格式。
  * 支持分享自定义文件夹。
  * 支持多字幕和多音轨，记录播放进度。
  * 支持在线聊天室，通过分享的URL，其他用户可以加入到自定义文件夹对应的聊天室中。
  * 支持字幕上传，智能匹配
  * 支持自动适配系统主题
## Docker
  * 创建一个目录A, 接着在这个目录下创建2个证书私钥（http.crt http.key和rpc.http rpc.key）
  * sudo python3 docker.py --crt-path \<path to A\> --build-myenv --nvidia-init -r -n sn -i