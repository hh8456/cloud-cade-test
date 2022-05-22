docker build -t test:4567/app .

imageid=$(docker images |grep test:4567|awk '{print $3}')
echo $imageid
docker run -d -p 4567:4567 $imageid
docker ps -a
