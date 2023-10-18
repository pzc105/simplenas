FROM ubuntu:22.04
RUN apt update && apt install -y
RUN apt-get update && apt-get install -y \
  autoconf \
  automake \
  build-essential \
  curl \
  g++ \
  git \
  libtool \
  make \
  pkg-config \
  unzip \
  zlib1g-dev \
  libssl-dev \
  wget \
  libboost-all-dev \
  golang \
  && apt-get clean

RUN apt-get update && \
    apt-get install -y ca-certificates curl gnupg && \
    mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_18.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && \
    apt-get install nodejs -y

RUN mkdir -p /tmp/src/cmake && cd /tmp/src/cmake && \
    wget https://cmake.org/files/v3.23/cmake-3.23.0.tar.gz && tar -zxvf cmake-3.23.0.tar.gz && \
    cd cmake-3.23.0 && ./configure && make -j$(nproc) && make install

RUN git config --global http.proxy socks5://192.168.10.241:10808
RUN mkdir /thirdparty
RUN git clone --recursive https://github.com/protocolbuffers/protobuf.git /thirdparty/protobuf && \
    cd /thirdparty/protobuf && \
    git switch 25.x && \
    mkdir cmake/build && \
    cd cmake/build && \
    cmake ../.. && \
    make -j$(nproc) && \
    make install
RUN git clone --recursive https://github.com/grpc/grpc.git /thirdparty/grpc && \
    cd /thirdparty/grpc && \
    git switch v1.59.x && \
    mkdir cmake/build && \
    cd cmake/build && \
    cmake -DgRPC_SSL_PROVIDER=package -DOPENSSL_USE_STATIC_LIBS=ON -DgRPC_MSVC_STATIC_RUNTIME=ON ../.. && \
    make -j$(nproc) && \
    make install