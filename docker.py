import os, sys, getopt

git_proxy = ""
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
  gen_docker_image(dc, "myenv")

def gen_sn():
  f = open("./Dockerfile")
  dc = f.read()
  f.close()
  if len(git_proxy) != 0:
    dc.replace("#git_proxy", "RUN git config --global http.proxy {0}".format(git_proxy))
  gen_docker_image(dc, "sn")

def main():
  global git_proxy
  global server_config
  global bt_config
  global tls_config
  try:
    opts, args = getopt.getopt(sys.argv[1:], "rn:i", ["git_proxy=", "server_config=", "bt_config=", "tls_config=", "build_myenv="])
  except getopt.GetoptError:
    sys.exit(2)
  build_myenv = True
  run_container = False
  container_name = ""
  init_container = False
  for opt, arg in opts:
    if opt == '--git_proxy':
      git_proxy = arg
    elif opt == '--server_config':
      server_config = arg
    elif opt == '--bt_config':
      bt_config = arg
    elif opt == '--tls_config':
      tls_config = arg
    elif opt == '--build_myenv' and (arg.lower() == "false" or arg.lower() == "0"):
      build_myenv = False
    elif opt == '-r' and (arg.lower() == "true" or arg.lower() != "0"):
      run_container = True
    elif opt == '-n':
      container_name = arg
    elif opt == '-i' and (arg.lower() == "true" or arg.lower() != "0"):
      init_container = True
  if build_myenv:
    gen_myenv()
  gen_sn()

  if run_container:
    os.system("sudo docker run -p 3000:3000 -p 6881:6881 -p 22345:22345 -p 11236:11236 --name {0} -dti sn".format(container_name))
  
  if init_container:
    os.system("sudo docker cp {0} {1}:/app".format(server_config, container_name))
    os.system("sudo docker cp {0} {1}:/app".format(bt_config, container_name))
    os.system("sudo docker cp {0}/http.crt {1}:/app/tls".format(tls_config, container_name))
    os.system("sudo docker cp {0}/http.key {1}:/app/tls".format(tls_config, container_name))
    os.system("sudo docker cp {0}/rpc.crt {1}:/app/tls".format(tls_config, container_name))
    os.system("sudo docker cp {0}/rpc.key {1}:/app/tls".format(tls_config, container_name))
    
    os.system("sudo docker exec {0} /bin/bash -c 'service mysql start && service redis-server start'".format(container_name))
    os.system("sleep 3")
    os.system("sudo docker exec {0} /bin/bash -c \"cd /app && (./bt && ./pnas &)\"".format(container_name))
    os.system("sudo docker exec {0} /bin/bash -c \"echo 'REACT_APP_RPC_SERVER=https://rpc.pnas105.top:11236' > /source/simplenas/src/frontend/.env.local\"".format(container_name))
    os.system("sudo docker exec {0} /bin/bash -c 'cd /source/simplenas/src/frontend && chmod +x start_unix.sh && ./start_unix.sh -c /app/tls/http.crt -k /app/tls/http.key &'".format(container_name))


if __name__ == "__main__":
  main()