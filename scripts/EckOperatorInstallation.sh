minikube delete --all
minikube start
kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v0.17.0/crds.yaml
kubectl apply -f https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v0.17.0/olm.yaml
kubectl create -f https://operatorhub.io/install/elastic-cloud-eck.yaml
minikube dashboard
# Search for kubernetes-dashboard-token (Copy token into config/appconfig.json )