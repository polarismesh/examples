# book info exmaple

本项目与 istio 的 bookinfo demo 基本相同。修改了其中的 reviews 服务和 ratings 服务，details、productpage 与 istio 的 demo 相同。

reviews 在本项目中服务基于 spring cloud tencent 构建，自注册服务。

ratings 增加了参数，控制是否返回错误，用来触发 reviews 服务实例产生熔断。

## 部署 demo

reviews 服务为 spring cloud tencent 应用，不需要注入 Envoy ，同时也不需要 Polaris Controller 托管服务注册。这里需要为 reviews 服务指定 Polaris Server 的地址，需要修改 bookinfo.yaml ，将其中 `spring.cloud.polaris.address` 的参数指定为 Polaris Server 的地址。

修改完 reviews 相关的配置后，您可以使用下面命令来创建部署本 demo：

```
kubectl create -f bookinfo.yaml
```

等待所有 pod 转为 `Running` 状态：
```
NAME                              READY   STATUS    RESTARTS   AGE
details-v1-66b6955995-jhnvw       2/2     Running   0          39m
productpage-v1-7955cdc67f-qt5zj   2/2     Running   0          39m
ratings-v1-846b7966bc-gq9wh       2/2     Running   0          39m
reviews-v1-5dd557658c-74h6x       1/1     Running   0          39m
reviews-v2-6bb7d454bc-wznj9       1/1     Running   0          39m
reviews-v3-85d446d4bf-qbf24       1/1     Running   0          39m
```

可以看到 details、productpage、ratings 除了业务容器，还运行了 Envoy Sidecar 容器。reviews 只有一个容器。reviews 服务有三个版本，每个版本返回的星星颜色不同。

## 验证 demo 是否正常工作

使用下面的命令，可以看到结果，表示 pod 在正常运行：

```
kubectl exec -it "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl productpage:9080/productpage | grep -o "<title>.*</title>"

<title>Simple Bookstore App</title>
```

同时可以修改 productpage 服务的 service 类型为 NodePort/Loadbalancer ，接着通过网页访问 productpage 服务，一切运行正常的话，您可以在浏览器中看到下面的内容：

![image](pic/productpage.png)

左侧的 Book Details 是 productpage 请求 details 返回的结果。右侧的 Book Reviews 是 productpage 请求 reviews 服务返回的结果。

每次刷新页面，可以看到右侧的星星随机在红、绿、蓝之间变化。