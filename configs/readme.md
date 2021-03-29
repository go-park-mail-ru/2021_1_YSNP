
docker build --tag=mytnt . 
docker run --name mytnt-inst -p 3301:3301 -d mytnt 
docker exec -t -i mytnt-inst console
sudo docker exec -it 2021_1_ysnp_postgres_1 psql -U postgres
