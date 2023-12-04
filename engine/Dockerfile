#sudo docker build -t battleground-engine .
#sudo docker run battleground-engine
FROM golang:1.19.2-bullseye

ENV USERNAME=worker_user
ENV USER_UID=1001
ENV USER_GID=$USER_UID

WORKDIR /engine

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o /battleground-engine

RUN groupadd -r -g $USER_GID $USERNAME 
RUN useradd -r -u $USER_UID -g $USER_GID $USERNAME
RUN chmod -R 755 . 

#EXPOSE 8089 <- exposing programmatically