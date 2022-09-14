docker stop urlshortnergo
docker rm   urlshortnergo
docker rmi  urlshortnergo

docker build -t urlshortnergo .
docker run -d -it -p 18080:18080 --name=urlshortnergo urlshortnergo