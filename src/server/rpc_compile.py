import os
import glob

def GenServerRpc():
    oldpath = os.getcwd()
    script_path = os.path.dirname(os.path.abspath(__file__))
    os.chdir(script_path)
    paths = ["../proto ./"]
    common_rpc_path = "../proto"
    google_api_protos = ["../proto/google/api/annotations.proto", "../proto/google/api/http.proto"]
    for in_path in paths:

        rpc_path, out_path = in_path.split(" ")

        proto_path = "--proto_path=. --proto_path={0} --proto_path={1}".format(rpc_path, common_rpc_path)
        cmd_pattern = "protoc \
                        --grpc-gateway_out={out_path} \
                        --grpc-gateway_opt logtostderr=true \
                        --grpc-gateway_opt generate_unbound_methods=true \
                        --go_out={out_path} \
                        --go-grpc_out={out_path} {proto_path} {file_name}"

        rpc_files = glob.glob("{0}/*.proto".format(rpc_path))
        for fn in rpc_files:
            cmd = cmd_pattern.format(
                proto_path=proto_path, out_path=out_path, file_name=fn)
            fd = os.popen(cmd)
            ot = fd.read()
            if (len(ot) > 0):
                print("e:", ot)
    os.chdir(oldpath)

if __name__ == "__main__":
    GenServerRpc()