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
    opts, args = getopt.getopt(sys.argv[1:], "r:i", ["git_proxy=", "server_config=", "bt_config=", "tls_config=", "build_myenv="])
  except getopt.GetoptError:
    sys.exit(2)
  build_myenv = True
  run_container = ""
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
    elif opt == '-r':
      run_container = arg
    elif opt == '-i' and (arg.lower() == "false" or arg.lower() == "0"):
      init_container = True
  if build_myenv:
    gen_myenv()
  gen_sn()

  if len(run_container) > 0:
    os.system("sudo docker run -p 3000:3000 -p 6881:6881 -p 22345:22345 -p 11236:11236 --name sn -dti sn")
  
  if init_container:
    os.system("sudo docker cp {0} sn:/app".format(server_config))
    os.system("sudo docker cp {0} sn:/app".format(bt_config))
    os.system("sudo docker cp {0}/http.crt sn:/app/tls".format(tls_config))
    os.system("sudo docker cp {0}/http.key sn:/app/tls".format(tls_config))
    os.system("sudo docker cp {0}/rpc.crt sn:/app/tls".format(tls_config))
    os.system("sudo docker cp {0}/rpc.key sn:/app/tls".format(tls_config))
    os.system("sudo docker exec -it sn sh -c 'service mysql start && service redis-server start'")
    os.system("sudo docker exec -it sn sh -c 'cd /app && ./bt &'")
    os.system("sudo docker exec -it sn sh -c 'cd /app && ./pnas &'")
    os.system("sudo docker exec -it sn sh -c 'cd /source/simplenas/src/frontend && chmod +x start_unix.sh && ./start_unix.sh -c /app/tls/http.crt -k /app/tls/http.key &'")


if __name__ == "__main__":
  main()