ARG VARIANT=1.18
ARG OPENCV=4.5.5
# use ssh to get private git. set 'https' if you want.
ARG con=ssh
# path from workdir
ARG MAIN_PATH=./cmd


FROM ghcr.io/machinemapplatform/devcontainer-go${VARIANT}-addlibkdslam-opencv-${OPENCV} as golang-opencv


FROM golang-opencv as build-ssh
# ビルドに必要なため、tmpからlibKdSlam.soを移動(サイズが大きい為)
RUN mv /tmp/libKdSlam.so /usr/local/lib
RUN mv /tmp/ORBvoc.kdbow /usr/local/lib
ARG MAIN_PATH
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
COPY . /app
WORKDIR /app
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
RUN --mount=type=ssh go mod tidy
RUN go build -o main ${MAIN_PATH}


FROM golang-opencv as build-https
# ビルドに必要なため、tmpからlibKdSlam.soを移動(サイズが大きい為)
RUN mv /tmp/libKdSlam.so /usr/local/lib
RUN mv /tmp/ORBvoc.kdbow /usr/local/lib
ARG MAIN_PATH
COPY . /app
WORKDIR /app
COPY .netrc /root/.netrc
RUN go mod tidy
RUN go build -o main ${MAIN_PATH}


FROM build-${con} AS build


FROM ghcr.io/machinemapplatform/devcontainer-go${VARIANT}-addlibkdslam-opencv-${OPENCV}:latest
RUN mkdir /app
WORKDIR /app

ENV LD_LIBRARY_PATH=/usr/local/lib:/app/lib
ENV CGO_LDFLAGS="-L/usr/local/lib -lKdSlam"
COPY --from=build /usr/local/lib /usr/local/lib
COPY --from=build /app/main .
# 軽量化
RUN rm -f /tmp/libKdSlam.so /tmp/*.kdmp /tmp/*.kdbow

CMD ["./main"]
