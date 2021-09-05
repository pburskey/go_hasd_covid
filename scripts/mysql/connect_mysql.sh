#use letmein
docker exec -it go_hasd_covid_mysql_1 mysql -p

docker exec -it go_hasd_covid_mysql_1 mysql -uroot -pletmein

docker exec -it go_hasd_covid_mysql_1 mysql -uroot -pletmein -e 'connect hasd_covid; select * from school;'

docker exec -it go_hasd_covid_mysql_1 mysql -uroot -pletmein hasd_covid < path-to-file.sql