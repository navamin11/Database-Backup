#
# STAGE 1: prepare
#
# Start from golang base image
FROM golang:1.22.3-alpine as build

# Updates the repository and installs git
RUN apk update && apk upgrade && apk --no-cache add git && apk --no-cache add tzdata && rm -rf /var/cache/apk/*

WORKDIR /go/src

COPY go.mod .
COPY go.sum .
COPY config.json .
COPY houseKeeping.sh .

#
# STAGE 2: build
#
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./dbBackup-v1.0.0 ./cmd

RUN ls -la .

#########################################################
#
# STAGE 3: run
#
# The project has been successfully built and we will use a
# lightweight alpine image to run the server
FROM alpine:latest as Dev

# Adds Package to the image
RUN apk update && \
    apk upgrade && \
    apk add --no-cache ca-certificates && \
    apk add --no-cache tzdata && \
    # apk add --no-cache doas && \
    apk add --no-cache bash && \
    apk add --no-cache sudo && \
    apk add --no-cache openssh && \
    rm -rf /var/cache/apk/*
    
# RUN mkdir -p /var/run/sshd

RUN adduser -D support; \
    echo 'support:password' | chpasswd

RUN ssh-keygen -A && \
    mkdir /root/.ssh && \
    chmod 0700 /root/.ssh && \
    echo 'root:password' | chpasswd && \
    ln -s /etc/ssh/ssh_host_ed25519_key.pub /root/.ssh/authorized_keys
    
RUN echo 'PermitRootLogin no' >> /etc/ssh/sshd_config && \
    echo 'PasswordAuthentication yes' >> /etc/ssh/sshd_config

WORKDIR /home/support

RUN mkdir -p bin/logs && \
    chown -R support:support bin/logs


# Copies the binary file from the BUILD container to /app folder
COPY --from=build --chown=support /go/src/dbBackup-v1.0.0 ./bin/dbBackup-v1.0.0
COPY --from=build --chown=support /go/src/config.json ./bin/config.json
COPY --from=build --chown=support /go/src/houseKeeping.sh ./bin/houseKeeping.sh

EXPOSE 22

WORKDIR /home/support/bin

RUN chmod 0755 dbBackup-v1.0.0

# Runs the binary once the container starts
# CMD ./dbBackup-v1.0.0 && \
#     /usr/sbin/sshd -D

CMD /usr/sbin/sshd -D