ARG GO_VERSION=1.18
# use ssh to get private git. set 'https' if you want.
ARG con=ssh

FROM golang:${GO_VERSION} as golang


FROM golang as build-ssh
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
COPY . /app
WORKDIR /app
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
RUN --mount=type=ssh go mod tidy -compat=1.18
RUN go build -o main ./cmd


FROM golang as build-https
COPY . /app
WORKDIR /app
COPY .netrc /root/.netrc
RUN go mod tidy
RUN go build -o main ./cmd


FROM build-${con} AS build


FROM golang:${GO_VERSION}
RUN mkdir /app
WORKDIR /app

COPY --from=build app/main .

CMD ["./main"]
