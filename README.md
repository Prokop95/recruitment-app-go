# Getting Started

### Reference Documentation
For further reference, please consider the following sections:

* CREATE DOCKER NETWORK `docker network create project_network`
* RUN CASSANDRA `docker run --name cassandra -d -p 9042:9042 --net=project_network cassandra`
* BUILD PROJECT `docker build -t recruitment .`
* RUN PROJECT `docker run --net=project_network -d -p 8080:8080 --name=recruitment_go recruitment`
