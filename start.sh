
pid=$(sudo docker ps -aq)

if [ -z "$pid" ]
then
      echo "\Container Not found"
else
      echo "\$pid found.  Ending docker instance"
      sudo docker stop $pid
      sudo docker rm $pid
fi

docker pull redis
docker run --name redis-test-instance -p 6379:6379 -d redis

#docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root mysql
#docker exec -it mysql mysql -uroot -proot -e 'CREATE DATABASE hasd_covid'


