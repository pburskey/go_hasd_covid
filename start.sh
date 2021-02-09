docker pull redis
docker run --name redis-test-instance -p 6379:6379 -d redis

#docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root mysql
#docker exec -it mysql mysql -uroot -proot -e 'CREATE DATABASE hasd_covid'

cd go_src
go run service/main.go /home/patrickburskey/IdeaProjects/go/go_hasd_covid/data
