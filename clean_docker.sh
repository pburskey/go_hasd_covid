dockerImageName=$1
echo "Searching for docker image name: ${dockerImageName}"
pids=$(docker ps | grep $dockerImageName | sed -e 's/\s.*$//')

if [ -z "$pids" ]
then
      echo "Container Not found"
else
    for pid in $pids
    do
      echo "pid: ${pid} found.  Ending docker instance"
      docker stop $pid
      docker rm $pid
    done
fi