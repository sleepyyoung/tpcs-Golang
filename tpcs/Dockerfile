FROM golang:alpine

ENV TZ Asia/Shanghai

WORKDIR $GOPATH/src/tpcs
COPY . $GOPATH/src/tpcs

RUN apk -U upgrade && apk add tzdata && cp /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone

RUN mkdir -p /tpcs/logs
RUN mkdir -p /tpcs/uploads
RUN mkdir -p /tpcs/files

EXPOSE 8000
ENTRYPOINT ["./tpcs"]

# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o tpcs .

# docker build -t tpcs .

# docker run --name mysql -p 3307:3306 -e MYSQL_ROOT_PASSWORD=19834044876 -e TZ="Asia/Shanghai" -d philpy/tpcs-mysql:v1 --lower_case_table_names=1

# docker run --name redis -p 6380:6379 --privileged=true -v D:/Code/GitProject/tpcs-Golang/docker/redis/conf/redis.conf:/etc/redis/redis.conf -v D:/Code/GitProject/tpcs-Golang/docker/redis/data:/data -d redis redis-server /etc/redis/redis.conf --appendonly yes

# docker run --name tpcs -p 8080:8080 --link mysql:mysql --link redis:redis -v D:/Code/GitProject/tpcs-Golang:/tpcs tpcs

