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
    dc.replace("#git_proxy", "RUN git config --global http.proxy={0}".format(git_proxy))
  gen_docker_image(dc, "myenv")

def gen_sn():
  f = open("./Dockerfile")
  dc = f.read()
  f.close()
  if len(git_proxy) != 0:
    dc.replace("#git_proxy", "RUN git config --global http.proxy={0}".format(git_proxy))
  dc += "\nCOPY {0} /app".format(server_config)
  dc += "\nCOPY {0} /app".format(bt_config)
  dc += "\nCOPY {0}/http.crt /app/tls/".format(tls_config)
  dc += "\nCOPY {0}/http.key /app/tls/".format(tls_config)
  dc += "\nCOPY {0}/rpc.crt /app/tls/".format(tls_config)
  dc += "\nCOPY {0}/rpc.key /app/tls/".format(tls_config)
  gen_docker_image(dc, "sn")

def main():
  global git_proxy
  global server_config
  global bt_config
  global tls_config
  try:
    opts, args = getopt.getopt(sys.argv[1:], "", ["git_proxy=", "server_config=", "bt_config=", "tls_config="])
  except getopt.GetoptError:
    sys.exit(2)
  for opt, arg in opts:
    if opt == '--git_proxy':
      git_proxy = arg
    elif opt == '--server_config':
      server_config = arg
    elif opt == '--bt_config':
      bt_config = arg
    elif opt == '--tls_config':
      tls_config = arg
  
  gen_myenv()
  gen_sn()

if __name__ == "__main__":
  main()