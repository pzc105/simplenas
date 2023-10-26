import os, sys, getopt

git_proxy = ""
all_proxy = ""
server_config = "./server.yml"
bt_config = "./bt.yml"
tls_config = "./tls"

def gen_docker_image(content, image_name):
  tmp_dir = "./tmp"
  if not os.path.exists(tmp_dir):
    os.makedirs(tmp_dir)
  f = open(tmp_dir + "/Dockerfile", "w")
  f.write(content)
  f.close()
  os.system("sudo docker image build -t {0} {1}".format(image_name, tmp_dir))


def gen_myenv():
  f = open("./myenv/Dockerfile")
  dc = f.read()
  f.close()
  if len(git_proxy) != 0:
    dc = dc.replace("#git_proxy", "RUN git config --global http.proxy {0}".format(git_proxy))
  if len(all_proxy) != 0:
    dc = dc.replace("#all_proxy", "ARG all_proxy={0}".format(all_proxy))
  gen_docker_image(dc, "myenv")

def gen_sn():
  f = open("./Dockerfile")
  dc = f.read()
  f.close()
  if len(git_proxy) != 0:
    dc = dc.replace("#git_proxy", "RUN git config --global http.proxy {0}".format(git_proxy))
  if len(all_proxy) != 0:
    dc = dc.replace("#all_proxy", "ARG all_proxy={0}".format(all_proxy))
  gen_docker_image(dc, "sn")

def main():
  global git_proxy
  global server_config
  global bt_config
  global tls_config
  global all_proxy
  try:
    opts, args = getopt.getopt(sys.argv[1:], "rn:i", ["rpc-server=", "git-proxy=", "all-proxy=", "server-config=", "bt-config=", "crt-path=", "build-myenv", "build-sn", "nvidia-init"])
  except getopt.GetoptError:
    sys.exit(2)
  build_myenv = False
  build_sn = False
  run_container = False
  container_name = ""
  init_container = False
  nvidia_init = False
  rpc_server_address = "https://rpc.pnas105.top:11236"
  for opt, arg in opts:
    if opt == '--git-proxy':
      git_proxy = arg
    elif opt == '--all-proxy':
      all_proxy = arg
    elif opt == '--server-config':
      server_config = arg
    elif opt == '--bt-config':
      bt_config = arg
    elif opt == '--crt-path':
      tls_config = arg
    elif opt == 'nvidia-init':
      nvidia_init = True
    elif opt == '--build-myenv':
      build_myenv = True
    elif opt == '--build-sn':
      build_sn = True
    elif opt == '-r':
      run_container = True
    elif opt == '-n':
      container_name = arg
    elif opt == '-i':
      init_container = True
    elif opt == '--rpc-address':
      if not arg.startswith("https://") and arg.startswith("http://"):
        print("must be https rpc-address")
        exit(-1)
      rpc_server_address = arg
      if not rpc_server_address.startswith("https://"):
        rpc_server_address = "https://" + rpc_server_address
  if build_myenv:
    gen_myenv()
  if build_sn:
    gen_sn()

  if nvidia_init:
    os.system("distribution=$(. /etc/os-release;echo $ID$VERSION_ID) && \
              curl -s -L https://nvidia.github.io/nvidia-docker/gpgkey | sudo apt-key add - && \
              curl -s -L https://nvidia.github.io/nvidia-docker/$distribution/nvidia-docker.list | sudo tee /etc/apt/sources.list.d/nvidia-docker.list && \
              sudo apt update && sudo apt install -y nvidia-container-toolkit && \
              sudo service docker restart")
  if run_container:
    os.system("sudo docker run --gpus all,capabilities=video -p 3000:3000 -p 6881:6881 -p 6881:6881/udp -p 6771:6771 -p 6771:6771/udp -p 22345:22345 -p 11236:11236 --name {0} -dti sn".format(container_name))
  
  if init_container:
    os.system("sudo docker cp {0} {1}:/app".format(server_config, container_name))
    os.system("sudo docker cp {0} {1}:/app".format(bt_config, container_name))
    os.system("sudo docker cp {0}/http.crt {1}:/app/tls".format(tls_config, container_name))
    os.system("sudo docker cp {0}/http.key {1}:/app/tls".format(tls_config, container_name))
    os.system("sudo docker cp {0}/rpc.crt {1}:/app/tls".format(tls_config, container_name))
    os.system("sudo docker cp {0}/rpc.key {1}:/app/tls".format(tls_config, container_name))
    
    os.system("sudo docker exec {0} /bin/bash -c 'service mysql start && service redis-server start'".format(container_name))
    os.system("sudo docker cp ./wait_db.sh {0}:/app".format(container_name))
    os.system("sudo docker exec {0} /bin/bash -c '/bin/bash /app/wait_db.sh'".format(container_name))
    os.system("sudo docker exec {0} /bin/bash -c 'mysql -uroot -p123 < /source/simplenas/src/server/tables.sql'".format(container_name))

    os.system("sudo docker exec {0} /bin/bash -c 'mkdir -p /app/media/poster && \
                cp /source/simplenas/default_folder.png /app/media/poster/ && \
                cp /source/simplenas/house.png /app/media/poster/'".format(container_name))

    os.system("sudo docker exec {0} /bin/bash -c \"cd /app && (nohup ./bt &) && (nohup ./pnas &)\"".format(container_name))
    os.system("sudo docker exec {0} /bin/bash -c \"echo 'REACT_APP_RPC_SERVER={1}' > /source/simplenas/src/frontend/.env.local\"".format(container_name, rpc_server_address))
    os.system("sudo docker exec {0} /bin/bash -c 'cd /source/simplenas/src/frontend && \
      npm run build && mkdir -p /app/frontend && cp -rf build/* /app/frontend'".format(container_name))
    os.system("sudo docker cp ./nginx.conf {0}:/etc/nginx/".format(container_name))
    os.system("sudo docker exec {0} /bin/bash -c 'service nginx restart'".format(container_name))

if __name__ == "__main__":
  main()