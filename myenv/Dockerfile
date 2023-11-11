FROM ubuntu:22.04

ARG TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
ENV LANG="C.UTF-8"

#all_proxy
ARG DEBIAN_FRONTEND=noninteractive
RUN apt update -y && \
  apt install -y
RUN apt install curl wget software-properties-common aria2 -y && \
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
RUN git clone --recursive https://github.com/grpc/grpc.git /thirdparty/grpc && \
  cd /thirdparty/grpc && \
  git switch v1.59.x && \
  cd ./third_party/protobuf && \
  git switch 24.x && \
  cd - && \
  mkdir cmake/build && \
  cd cmake/build && \
  cmake -DgRPC_SSL_PROVIDER=package -DOPENSSL_USE_STATIC_LIBS=ON -DgRPC_MSVC_STATIC_RUNTIME=ON ../.. && \
  make -j$(nproc) && \
  make install

RUN apt-fast install -y \
  nasm \
  doxygen \
  python3 \
  python3-dev \
  python3-pip \
  python3-setuptools \
  python3-wheel \
  python3-tk \
  build-essential \
  ninja-build \
  checkinstall

RUN all_proxy="" pip install pysocks -i https://pypi.tuna.tsinghua.edu.cn/simple --trusted-host pypi.tuna.tsinghua.edu.cn
RUN pip3 install --upgrade pip -i https://pypi.tuna.tsinghua.edu.cn/simple --trusted-host pypi.tuna.tsinghua.edu.cn
RUN pip3 install --no-cache-dir meson cython numpy -i https://pypi.tuna.tsinghua.edu.cn/simple --trusted-host pypi.tuna.tsinghua.edu.cn
RUN cd /thirdparty && \
  git clone https://github.com/Netflix/vmaf.git && \
  cd ./vmaf && \
  make -j$(nproc) && make install
RUN cd /thirdparty && \
  git clone https://github.com/AviSynth/AviSynthPlus && \
  cd AviSynthPlus && \
  mkdir avisynth-build && cd avisynth-build && \
  cmake ../ -G Ninja && ninja && \
  checkinstall --pkgname=avisynth --pkgversion="$(grep -r \
  Version avs_core/avisynth.pc | cut -f2 -d " ")-$(date --rfc-3339=date | \
  sed 's/-//g')-git" --backup=no --deldoc=yes --delspec=yes --deldesc=yes  --strip=yes --stripso=yes --addso=yes --fstrans=no --default ninja install
RUN cd /thirdparty && \
  git clone https://git.videolan.org/git/ffmpeg/nv-codec-headers.git && \
  cd nv-codec-headers && \
  make && make install
RUN cd /thirdparty && \
  wget https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2204/x86_64/cuda-keyring_1.1-1_all.deb && \
  dpkg -i cuda-keyring_1.1-1_all.deb && \
  apt-fast update && \
  apt-fast -y install cuda-toolkit-12-3 && \
  rm -r cuda-keyring_1.1-1_all.deb
ENV PATH="/usr/local/cuda-12.3/bin:$PATH"

RUN apt-fast install -y \
  libvpl-dev \
  libgsm1-dev \
  libgme-dev \
  libaom-dev \
  libgnutls28-dev \
  libgmp-dev \
  libvpl2 \
  libmp3lame-dev \
  libopencore-amrnb-dev \
  libopencore-amrwb-dev \
  libopenjp2-7-dev \
  libopenjp2-7-dev \
  libopus-dev \
  librubberband-dev \
  libssh-dev \
  libspeex-dev \
  libsrt-gnutls-dev \
  libtheora-dev \
  libvidstab-dev \
  libvo-amrwbenc-dev \
  libvpx-dev \
  libwebp-dev \
  libx264-dev \
  libx265-dev \
  libxvidcore-dev \
  libzimg-dev \
  libczmq-dev \
  libghc-bzlib-dev \
  liblzma-dev \
  libass-dev \
  libsdl2-dev \
  libopenmpt-dev \
  libcodec2-dev

RUN apt clean && \
  rm -rf /var/lib/apt/lists

RUN cd /thirdparty && \
  git clone --recursive https://github.com/FFmpeg/FFmpeg.git && \
  cd FFmpeg && \
  ./configure --enable-gpl --enable-version3 --enable-static --disable-w32threads --disable-autodetect --enable-fontconfig --enable-iconv --enable-gnutls --enable-libxml2 --enable-gmp --enable-bzlib --enable-lzma --enable-zlib --enable-libsrt --enable-libssh --enable-libzmq --enable-avisynth --enable-sdl2 --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxvid --enable-libaom --enable-libopenjpeg --enable-libvpx --enable-libass --enable-libfreetype --enable-libfribidi --enable-libvidstab --enable-libvmaf --enable-libzimg  --enable-cuda-llvm --enable-cuvid --enable-ffnvcodec --enable-nvdec --enable-nvenc --enable-libvpl --enable-libgme --enable-libopenmpt --enable-libopencore-amrwb --enable-libmp3lame --enable-libtheora --enable-libvo-amrwbenc --enable-libgsm --enable-libopencore-amrnb --enable-libopus --enable-libspeex --enable-libvorbis --enable-librubberband --enable-gpl --enable-libnpp --enable-cuda-nvcc --enable-nonfree --extra-cflags=-I/usr/local/cuda-12.3/targets/x86_64-linux/include --extra-ldflags=-L/usr/local/cuda-12.3/targets/x86_64-linux/lib --nvccflags="-gencode arch=compute_52,code=sm_52 -O2" && \
  make -j$(nproc) && make install