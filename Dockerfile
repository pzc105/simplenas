FROM myenv

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

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

# RUN pip install pysocks -i https://pypi.tuna.tsinghua.edu.cn/simple --trusted-host pypi.tuna.tsinghua.edu.cn
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

RUN cd /thirdparty && \
    wget https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2204/x86_64/cuda-keyring_1.1-1_all.deb && \
    dpkg -i cuda-keyring_1.1-1_all.deb && \
    apt-fast update && \
    apt-fast -y install cuda-toolkit-12-3 && \
    rm -r cuda-keyring_1.1-1_all.deb
RUN export PATH="/usr/local/cuda-12.3/bin:$PATH"

RUN apt clean && \
    rm -rf /var/lib/apt/lists

RUN cd /thirdparty && \
    git clone --recursive https://github.com/FFmpeg/FFmpeg.git && \
    cd FFmpeg && \
    ./configure --enable-gpl --enable-version3 --enable-static --disable-w32threads --disable-autodetect --enable-fontconfig --enable-iconv --enable-gnutls --enable-libxml2 --enable-gmp --enable-bzlib --enable-lzma --enable-zlib --enable-libsrt --enable-libssh --enable-libzmq --enable-avisynth --enable-sdl2 --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxvid --enable-libaom --enable-libopenjpeg --enable-libvpx --enable-libass --enable-libfreetype --enable-libfribidi --enable-libvidstab --enable-libvmaf --enable-libzimg  --enable-cuda-llvm --enable-cuvid --enable-ffnvcodec --enable-nvdec --enable-nvenc --enable-libvpl --enable-libgme --enable-libopenmpt --enable-libopencore-amrwb --enable-libmp3lame --enable-libtheora --enable-libvo-amrwbenc --enable-libgsm --enable-libopencore-amrnb --enable-libopus --enable-libspeex --enable-libvorbis --enable-librubberband --enable-gpl --enable-libnpp --enable-cuda-nvcc --enable-nonfree --extra-cflags=-I/usr/local/cuda-12.3/targets/x86_64-linux/include --extra-ldflags=-L/usr/local/cuda-12.3/targets/x86_64-linux/lib --nvccflags="-gencode arch=compute_52,code=sm_52 -O2" && \
    make -j$(nproc) && make install

RUN export LANG=C.UTF-8 && echo "export LANG=C.UTF-8" >> /etc/profile