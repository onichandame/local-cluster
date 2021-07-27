# Local Cluster

fully functional cluster all in native local environment _inspired by k8s_.

# Architecture

components and their counterparts in k8s:

| Local Cluster     | Kubernetes             |
| ----------------- | ---------------------- |
| application       | image(in docker)       |
| instance          | pod                    |
| instance group    | deployment/statefulset |
| service           | app                    |
| service interface | expose(in docker)      |
| gateway           | service                |
| entrance          | ingress                |

**note: the instance only holds one application whereas pod in k8s can have many containers**
