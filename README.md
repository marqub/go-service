# Simple Go microservice

Sample code just to demo resiliency testing.
Just one endpoint 'status' that will invoke the 'status' endpoint of a dependent microservice.

This simple "dependencies" graph is built by passing an environment variable named DEPENDENCY_NAME.
The dependency name is the internal name of the K8s service `<service name>.<namespace>`

The name of the microservice is set using the environment variable SVC_NAME or using the helm option `--name`

Plenty of opportunities to make it crashes without having a cyclic dependency graph: dependency not reachable, call timing out, invalid response ...

Deploy it using the helm chart. 
For example, to deploy 2 microservices and have one depends on the other:

``` helm install chart/ --name service2 --namespace resiliency-testing ```

``` helm install chart/ --name service1 --namespace resiliency-testing --set DEPENDENCY_NAME=service2.resiliency-testing ```

Ingresses, services and pods are created in the same namespace
``` 
$ kubectl get pods,services,ing -n resiliency-testing
NAME                                       READY     STATUS    RESTARTS   AGE
pod/service1-go-service-79bd4b8669-fqgq2   1/1       Running   0          53s
pod/service2-go-service-77f64948b8-p2mtz   1/1       Running   0          40s

NAME                          TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/service1-go-service   ClusterIP   xxx.xxx.xxx.xxx   <none>        8080/TCP   53s
service/service2-go-service   ClusterIP   xxx.xxx.xxx.xxx    <none>        8080/TCP   40s

NAME                                     HOSTS                             ADDRESS   PORTS     AGE
ingress.extensions/service1-go-service   service1.resiliency-testing.com             80        53s
ingress.extensions/service2-go-service   service2.resiliency-testing.com             80        40s
```