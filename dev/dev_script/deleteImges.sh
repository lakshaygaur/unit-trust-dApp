docker rmi -f $(docker images | grep  dev | awk '{print $3}')
docker volume prune -f
rm -rf fabric-client-kv-*
docker rm $( docker ps -qa)