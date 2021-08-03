# Local Cluster

fully functional cluster all in native local environment _inspired by k8s_.

# Architecture

components and their counterparts in k8s:

| Local Cluster      | Kubernetes              |
| ------------------ | ----------------------- |
| application        | image(in docker)        |
| instance           | pod                     |
| instance interface | expose(in docker)       |
| instance group     | deployment/statefulset  |
| gateway            | service                 |
| entrance           | ingress                 |
| storage            | persistent volume       |
| storage binding    | persistent volume claim |

**note: the instance only holds one application whereas pod in k8s can have many containers**

## Application

Application defines everything about how to run the app. An app can be of 3 types: **local**, **static** and **remote**.

### Local

A local application runs every part of it locally. It can interact with the outside world through a set of well-defined interfaces(a.k.a ports). A local application must have defined how to run itself on the local system(defined by architecture and platform). If the local system is not supported, the app will not run.

### Static

A static application acts as a static HTTP server. It must download all the files required then starts serving them. Hence the url of the archive must be provided. The site will be served through a single interface.

### Remote

A remote application points to an instance running on a remote system. The interfaces should be defined in advance.

## Instance

An Instance is an attempt to run an application. The status is retrieved by liveness probes. An instance must have all interfaces defined by its application. If any interface failed to setup, the instance cannot be run.

Every instance should be defined with the criteria of **ready** and **failed** by liveness probe and max retries respectively.

An instance can be created either directly, or by an instance group. An instance can be deleted either by instance group or manually.

## Instance Group

An instance is not reliable in terms of availability. Instance group is here to help. Basically it maintains multiple instances running. When any instance goes out of service, it tries to spin up another instance to take over.

## Gateway

The interfaces of instance and instance groups are generated at runtime, hence not reliable. Gateway can be defined with a pre-defined port, a service name(pointing to an instance group or an instance) and an application interface name(pointing to the desired interface to expose).

When multiple interfaces are matched, round robin algorithm is applied to select a target for each request. A gateway can be created only by user and deleted only by user.

## Entrance

One of the original purpose of the local cluster is to bundle all web servers together and access them from the same host. An entrance defines the path to access a certain gateway. e.g. if `http://localhost/` points to the cluster server, `http://localhost/entrances/abc` is the root of the entrance named `abc`.

An entrance can only created by user and deleted by user.

## Storage

For stateful applications there is data to keep among different instances, e.g. database files. A storage is basically a directory in local system. When it is bound to any(many) instance(s), the folder is symlinked to the runtime of the instance(s) as defined in the storage binding.

The storage can only be created manually and deleted manually. A storage binding can be created by the instance it binds to and deleted with the deletion of the binding instance.
