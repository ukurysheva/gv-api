FROM golang:latest  

COPY ./ ./
ENV GOPATH=/

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh
RUN migrate -path ./schema -database 'postgres://postgres:qwerty@globalavia-api.ru:5435/postgres?sslmode=disable' up
RUN go mod download
RUN go build -o gvapi ./cmd/main.go
CMD [ "./cmd/main" ]