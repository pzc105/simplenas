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
    opts, args = getopt.getopt(sys.argv[1:], " ", ["git_proxy=", "server_config=", "bt_config=", "tls_config=", "build_myenv="])
  except getopt.GetoptError:
    sys.exit(2)
  build_myenv = True
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
  if build_myenv:
    gen_myenv()
  gen_sn()

  os.system("sudo docker run --name sn -dti sn")
  os.system("sudo docker cp {0} sn:/app".format(server_config))
  os.system("sudo docker cp {0} sn:/app".format(bt_config))
  os.system("sudo docker cp {0}/http.crt sn:/app/tls".format(tls_config))
  os.system("sudo docker cp {0}/http.key sn:/app/tls".format(tls_config))
  os.system("sudo docker cp {0}/rpc.crt sn:/app/tls".format(tls_config))
  os.system("sudo docker cp {0}/rpc.key sn:/app/tls".format(tls_config))


if __name__ == "__main__":
  main()