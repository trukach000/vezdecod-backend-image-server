docker build -t imloader .  

docker run -dp 30001:80 -e GLOBAL_PREFIX=/imloader -e DB_HOST=host.docker.internal --add-host host.docker.internal:host-gateway  imloader 
