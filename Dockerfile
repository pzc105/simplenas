FROM myenv
RUN mkdir -p /app && mkdir -p /app/tls
#git_proxy
RUN mkdir -p /source
RUN git clone --recursive https://github.com/pzc105/simplenas.git /source/simplenas && \
    cd /source/simplenas && \
    mkdir -p cmake/build && \
    cd cmake/build && \
    cmake ../.. -Dstatic_runtime=true && \
    make -j$(nproc) && \
    cp ./bin/bt /app/bt
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    cd /source/simplenas/src/server && \
    go build && \
    cp ./pnas /app/pnas
RUN npm install -g pnpm && \
    cd /source/simplenas/src/frontend && \
    npm i

