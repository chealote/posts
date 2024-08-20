#!/bin/bash
container="postgres"
image="postgres:latest"
password="mysecretpassword"
port=5432

output=$(sudo docker ps -a | grep $container | wc -l)
echo "Running instances: $output"
if [ $output -gt 0 ]; then
	sudo docker start $container
	exit 0
else
	sudo docker run --name $container -e POSTGRES_PASSWORD=$password -p $port:$port -d $image
fi


# ti: docker run -ti --rm --network postgres postgres psql -h postgres -U postgres
