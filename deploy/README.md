# Putzing with the Docker image

> putz | To engage in aimless recreation or frivolous time-wasting; to fool around

After making your changes, you'll need to make sure the docker image builds, then tag and upload the image

Run commands from the project root directory:

Build the image

`docker build -f deploy/Dockerfile .`

Optional: debug a failed image/container

`docker container run -it --name=debug <container> /bin/sh`

Tag the image

`docker build -t rps_financial -f deploy/Dockerfile .`

Associate the tag

`docker tag rps_financial agundry/rps_financial:<version>`

Push the image

`docker push agundry/rps_financial:<version>`

Running in prod

`microk8s.start`

`kubectl apply -f prod-k8s-deployment.yml`

`sudo prometheus --config.file=prometheus.yml`

`sudo service grafana-server start`
