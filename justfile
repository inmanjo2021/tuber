token:
  sed -i "" "s/GCLOUD_.*/GCLOUD_TOKEN=`gcloud auth print-access-token`/" ./.env

start:
  sed -i "" "s/GCLOUD_.*/GCLOUD_TOKEN=`gcloud auth print-access-token`/" ./.env; go run main.go start

install:
  kubectl apply -f bootstrap.yaml
