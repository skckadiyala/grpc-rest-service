gRPC Blogging Service using golang

Deploy Mongodb on K8s Cluster

helm upgrade --install mongodb stable/mongodb --namespace blogsvc \
--set mongodbRootPassword=secretpassword \
--set mongodbUsername=ampc \
--set mongodbPassword=changeme \
--set mongodbDatabase=systest

Create docker image

make push 

To Deploy to k8s-cluster, create a manifest file with docker image create for this service.
