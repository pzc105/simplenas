import os
import glob

def GenBtRpc():
    oldpath = os.getcwd()
    script_path = os.path.dirname(os.path.abspath(__file__))
    os.chdir(script_path)
    paths = ["../proto ./prpc"]
    google_api_protos = ["../proto/google/api/annotations.proto", "../proto/google/api/http.proto"]
    commont_path = "../proto"

    for in_path in paths:

        rpc_path, out_path = in_path.split(" ")

        proto_path = "--proto_path={0} --proto_path={1}".format(rpc_path, commont_path)
        cmd_pattern = "protoc \
                        --cpp_out={out_path} \
                        --grpc_out={out_path} --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` \
                        {proto_path} {file_name}"

        rpc_files = glob.glob("{0}/*.proto".format(rpc_path))
        rpc_files += google_api_protos
        for fn in rpc_files:
            cmd = cmd_pattern.format(
                proto_path=proto_path, out_path=out_path, file_name=fn)
            fd = os.popen(cmd)
            ot = fd.read()
            if (len(ot) > 0):
                print("e:", ot)
    os.chdir(oldpath)
    
if __name__ == "__main__":
    GenBtRpc()