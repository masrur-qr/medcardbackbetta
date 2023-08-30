FROM golang:1.21rc3-bullseye
COPY . /new
WORKDIR /perent
COPY go.mod /perent/
COPY go.sum /perent/
RUN go mod download 
COPY . /perent/
EXPOSE 5500
RUN go build -o /main
CMD [ "/main" ]
