FROM golang:latest  

COPY ./ ./
ENV GOPATH=/

RUN go mod download
RUN go build -o gvapi ./cmd/main.go
CMD [ "./cmd/main" ]