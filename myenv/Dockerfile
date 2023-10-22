FROM ubuntu:22.04
RUN apt update -y && \
    apt install -y
RUN apt install curl software-properties-common aria2 -y && \
    /bin/bash -c "$(curl -sL https://git.io/vokNn)" && \
    add-apt-repository ppa:apt-fast/stable -y && \
    apt update && \
    apt install apt-fast

RUN apt-fast update && apt-fast install -y \
    autoconf \
    automake \
    build-essential \
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
    iproute2 \
    netcat \
    net-tools \
    vim \
    nginx \
    apt check && \
    apt clean

RUN add-apt-repository -y ppa:longsleep/golang-backports && \
    apt-fast update -y && \
    apt-fast install golang-go -y

RUN apt-fast install -y ca-certificates curl gnupg && \
    mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_18.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && \
    apt-get install nodejs -y

RUN mkdir -p /tmp/src/cmake && cd /tmp/src/cmake && \
    wget https://cmake.org/files/v3.23/cmake-3.23.0.tar.gz && tar -zxvf cmake-3.23.0.tar.gz && \
    cd cmake-3.23.0 && ./configure && make -j$(nproc) && make install && \
    cd .. && rm -rf cmake-3.23.0

RUN apt-fast install mysql-server redis-server -y && \
    sed -i "s/^bind-address/#bind-address/" /etc/mysql/mysql.conf.d/mysqld.cnf && \
    sed -i "/\[mysqld\]/a\bind-address            = 0.0.0.0" /etc/mysql/mysql.conf.d/mysqld.cnf && \
    service mysql start && \
    echo "alter user 'root'@'localhost' identified with mysql_native_password by '123'; \
    use mysql; \
    update user set host='%' where user ='root'; \
    FLUSH PRIVILEGES; \
    GRANT ALL PRIVILEGES ON *.* TO 'root'@'%'WITH GRANT OPTION; \
    FLUSH PRIVILEGES;" > alter.sql && \
    mysql < alter.sql && \
    rm alter.sql && \
    service mysql start && \
    service redis-server start

#git_proxy
# RUN mkdir /thirdparty
# RUN git clone --recursive https://github.com/protocolbuffers/protobuf.git /thirdparty/protobuf && \
#     cd /thirdparty/protobuf && \
#     git switch 25.x && \
#     mkdir cmake/build && \
#     cd cmake/build && \
#     cmake ../.. && \
#     make -j$(nproc) && \
#     make install
RUN git clone --recursive https://github.com/grpc/grpc.git /thirdparty/grpc && \
    cd /thirdparty/grpc && \
    git switch v1.59.x && \
    mkdir cmake/build && \
    cd cmake/build && \
    cmake -DgRPC_SSL_PROVIDER=package -DOPENSSL_USE_STATIC_LIBS=ON -DgRPC_MSVC_STATIC_RUNTIME=ON ../.. && \
    make -j$(nproc) && \
    make install