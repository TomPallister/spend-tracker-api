start-stop-daemon -SbmCv -x /usr/bin/nohup -d /home/mijitt0m/godutch-api api


start-stop-daemon -SbmCv -x /usr/bin/nohup -p /home/mijitt0m/godutch-api/godutch-api.pid -d /home/mijitt0m/godutch-api ./api


#Build the api
CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' .

#copy files to remote
scp ./dist/*.* mijitt0m@godutch-app-01.northeurope.cloudapp.azure.com:/home/mijitt0m/godutch-app
scp ./dist/assets/*.* mijitt0m@godutch-app-01.northeurope.cloudapp.azure.com:/home/mijitt0m/godutch-app/dist

scp ./api mijitt0m@godutch-app-01.northeurope.cloudapp.azure.com:/home/mijitt0m/godutch-api
scp ./.env mijitt0m@godutch-app-01.northeurope.cloudapp.azure.com:/home/mijitt0m/godutch-api
scp ./godutch.sql mijitt0m@godutch-app-01.northeurope.cloudapp.azure.com:/home/mijitt0m/godutch-api

#ember app
scp ~/.ssh/id_rsa.pub mijitt0m@godutch-pgsql-01.northeurope.cloudapp.azure.com:~/.ssh/id_rsa.pub
scp ~/.ssh/id_rsa mijitt0m@godutch-pgsql-01.northeurope.cloudapp.azure.com:~/.ssh/id_rsa

pg_dump -C godutch | ssh -C mijitt0m@godutch-app-01.northeurope.cloudapp.azure.com "psql godutch"


#build the container
docker build -t godutch/api .

#run containter interactive and headless

#prod 
docker run -itd -p 3001:3001 godutch/api

#prod with conn env passed in
docker run -itd -e PGSQL_CONNECTIONSTRING='postgres://godutch:password@godutch-pgsql-01.northeurope.cloudapp.azure.com/godutch?sslmode=require' -p 3001:3001 godutch/api

#dev
docker run -itd -e PGSQL_CONNECTIONSTRING='postgres://godutch:password@172.16.42.1:5432/godutch?sslmode=false' -p 3001:3001 godutch/api

#connect to container
docker exec -i -t a8c89bf7a91d  /bin/bash

#Delete all containers

docker rm $(docker ps -a -q)
#Delete all images

docker rmi $(docker images -q)

#Remote db connect
psql -h godutch-pgsql-01.northeurope.cloudapp.azure.com -p 5432 -U godutch -W password godutch

#Running the api
In order to run the example you need to have go and goget installed.

You also need to set the ClientSecret, ClientID and Domain for 
your Auth0 app as environment variables with the following names respectively: 
AUTH0_CLIENT_SECRET, AUTH0_CLIENT_ID and AUTH0_DOMAIN.

For that, if you just create a file named .env in the directory and set 
the values like the following, the app will just work:

````
# .env file
AUTH0_CLIENT_SECRET=myCoolSecret
AUTH0_CLIENT_ID=myCoolClientId
AUTH0_DOMAIN=myCoolDomain
````

Once you've set those 3 environment variables, you need to install all Go dependencies. For that, just run `go get .`.

Finally, run `go run main.go` to start the app and try calling http://localhost:3001/ping