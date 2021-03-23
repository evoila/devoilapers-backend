# Build
FROM golang:latest as build

WORKDIR /usr/src/app
COPY . .

WORKDIR /usr/src/app/cmd/service
ENV GOPATH=/usr/src/app/cmd/service
RUN go build -o service .


VOLUME /usr/src/app/configs

EXPOSE 8080
WORKDIR /usr/src/app
CMD /usr/src/app/cmd/service/service start -c "configs/appconfig_docker.json"

# Commands
# docker build -t devoilapers-backend .
# docker run --name devoilapers-backend -v /home/bene/Dokumente/evoila/:/usr/src/app/configs  -p 8080:8080 devoilapers-backend