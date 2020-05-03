# rps_financial [![CircleCI Build Status](https://circleci.com/gh/agundry/rps_financial.svg?style=shield)](https://circleci.com/gh/agundry/rps_financial)

Before going any further, let's set some expectations: This is a dumb application. That said, it is a cloud distributed, highly available dumb application written in golang that is scalable with Kubernetes.

So what does it do?

My girlfriend and I distribute the financial burdens of life via games of rock paper scissors, henceforth referred to as RPS.

Grocery shopping? RPS. Date night? RPS.

Somehow, I almost always lose.

This is an api meant to enable tracking these games and the costs associated, hopefully one day answering the question: which of rock, paper, and scissors is the most fiscally responsible choice?

# Dev setup

Make sure docker and docker-compose are installed, then in the root directory run `docker-compose up -d` to initialize the database.

### Single local instance

Assuming the local db is already running

`go mod download`

`make`

`./bin/rps_financial -config config/local-config.yml`

### Local Kubernetes Cluster via Docker Image

Start up minikube with docker driver

`minikube start --driver=docker`

Next, check what the host machine ip is in minikube

`minikube ssh "host host.docker.internal"`

You'll need to set the `addresses` ip for the mysql-service endpoint to this value in `local-k8s-deployment.yml`

Create the kubernetes cluster

`kubectl apply -f deploy/local-k8s-deployment.yml`

The services should have started automatically

Start a tunnel for the kubernetes service url

`minikube service rps-financial-service --url`

Test curl

`curl <url>/expenses`

### Local Metrics

If you want to see live metrics, you can start a local prometheus server and use the grafana docker container.

Edit the [prometheus.yml](config/prometheus.yml) file to point to the url minikube rps-financial-service url.

`prometheus --config.file="config/prometheus.yml"`

Go to grafana at [localhost:3000](localhost:3000) 
1. Add prometheus as a data source (source url `http://docker.for.mac.host.internal:9090`)
2. Create a new dashboard, you can import from any of these options
    * [Go Processes](https://grafana.com/grafana/dashboards/6671)
    * [Go Runtime](https://grafana.com/grafana/dashboards/1144)
    * [Go Metrics](https://grafana.com/grafana/dashboards/10826)