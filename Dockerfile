FROM myenv
RUN mkdir -p /app
# RUN git config --global http.proxy socks5://192.168.10.241:10808
RUN mkdir -p /source
RUN git clone --recursive https://github.com/pzc105/simplenas.git /source/simplenas && \
    cd /source/simplenas && \
    mkdir -p cmake/build && \
    cd cmake/build && \
    cmake ../.. -Dstatic_runtime=true && \
    make -j$(nproc) && \
    cp ./bin/bt /app/bt/bt
RUN cd /source/simplenas/src/server && \
    go build && \
    cp ./pnas /app/pnas/pnas
RUN npm install -g pnpm && \
    cd /source/simplenas/src/frontend && \
    npm i && npm build
