FROM google/debian:wheezy
MAINTAINER Tom Pallister <tom@threemammals.co.uk>
ADD api api
ADD .env .
ENV PGSQL_CONNECTIONSTRING postgres://godutch:password@godutch-pgsql-01.northeurope.cloudapp.azure.com/godutch?sslmode=require
EXPOSE 3001
ENTRYPOINT ["/api"]
