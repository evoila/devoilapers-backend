# Build
FROM golang:latest as build

WORKDIR /usr/src/app
COPY . .

WORKDIR /usr/src/app/cmd/service
ENV GOPATH=/usr/src/app/cmd/service
RUN go build -o service .

# COPY CAs
# configure path in .env.prod
COPY ${KUBERNETES_CERTIFICATE_AUTHORITY} /usr/src/app/configs/kubernetes_ca.crt

EXPOSE 8080
WORKDIR /usr/src/app
CMD /usr/src/app/cmd/service/service start -c "configs/appconfig_docker.json"

# Commands
# export KUBERNETES_CERTIFICATE_AUTHORITY='/home/bene/Dokumente/evoila/ca.crt'
# docker build -t devoilapers-backend .
# docker run --name devoilapers-backend -p 8080:8080 devoilapers-backend