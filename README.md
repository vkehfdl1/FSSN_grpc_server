# Deployment
First, you must build your project.
Then, you can deploy your project to fly.io

```shell
$ flyctl deploy
```

You can check grpc connection with the code below.

```shell
$ grpcurl -proto .\proto\FSSN_grpc_proto\hello_grpc.proto fssn.fly.dev:443 MyService/MyFunction
```
