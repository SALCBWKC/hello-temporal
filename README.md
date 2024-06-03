# hello-world-temporal
## 1. 启动Minikube集群
```
# brew install minikube
minikube start
```

## 2. 部署Zookeeper
```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install zookeeper bitnami/zookeeper --set persistence.volumeSize="5Gi"
```

## 3. 部署Kafka
```
helm install kafka bitnami/kafka --set persistence.volumeSize="20Gi"
```

## 4. 部署Temporal
```
git clone https://github.com/temporalio/helm-charts.git 
cd helm-charts/charts
helm dependencies update
helm install \
    --set server.replicaCount=1 \
    --set cassandra.config.cluster_size=1 \
    --set prometheus.enabled=false \
    --set grafana.enabled=false \
    --set elasticsearch.enabled=false \
    my-temporal . --timeout 15m
```

## 5.构建镜像
```
cd hello-world-temporal
docker build -t hello-world-temporal .
```
## 6.镜像上传到GitHub Container Registry
a) 创建token https://github.com/settings/tokens,  确保选择了“写入packages”权限

b) docker login ghcr.io -u username -p token  

c)docker tag your-local-image:tag ghcr.io/<YOUR_GITHUB_USERNAME>/<IMAGE_NAME>:<TAG>

d)docker push ghcr.io/<YOUR_GITHUB_USERNAME>/<IMAGE_NAME>:<TAG>

e)仓库或组织下的“Packages”页面，查看是否成功上传了Docker镜像。

## 7.部署应用
```
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```