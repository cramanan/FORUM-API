docker stop server
docker rm server
docker image build -f Dockerfile -t forum .
docker container run -td -p 8080:8080/tcp -p 8081:8081/tcp --name server forum
docker exec -it server /bin/bash