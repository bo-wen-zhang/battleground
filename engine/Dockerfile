#go test -v ./worker/worker_test.go
#sudo docker build -t battleground-engine .
#sudo docker run battleground-engine
FROM golang:1.21-bullseye AS base

WORKDIR /engine

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o /bin/battleground-engine .

FROM golang:1.21-bullseye

COPY --from=0 /bin/battleground-engine battleground-engine
#COPY --from=0 /engine/engine_logs.log engine_logs.log

ENV USERNAME=worker_user
ENV USER_UID=1001
ENV USER_GID=$USER_UID

RUN groupadd -r -g $USER_GID $USERNAME 
RUN useradd -r -u $USER_UID -g $USER_GID $USERNAME
RUN chmod -R 755 . 