# Kubernetes Operator

## Introduction

The operator observes a namespace for services which are annotated with a special xfsc tag. This tag contains information about the plugin name, its id etc. The annotation must be in the service deployment like this: 

```
xfsc.kubernetes.io/configuration: | 
              {"version":"v1","route":"/plugin-template", "name":"Plugin template", "serviceguid":"19bafa36-8415-44e9-a2de-32852fefa6ef","routeguid":"30959362-e35b-4eb4-afe5-e8fdd6a3fecd"}

```

When the operator finds such a annotation, he will start to add it into the kong services list, which is later on discovered by the plugin discovery. 


## Dependencies

[Kong](https://konghq.com/) - The API Gateway used to manage the plugins

[Kubernetes](https://kubernetes.io/) - The container orchestration system used to deploy the microservices. Used for plugin detection and synchronization

## Plugin Synchronization Annotation

If a micro service is deployed by using 

```
xfsc.kubernetes.io/component: xfsc.pcm.plugin
```

as label, the operator will pick it up and insert into the configured kong api and service/route to the api by using the metadata konfigured in the deployment annotations which are defined as:

```
xfsc.kubernetes.io/configuration: | 
              {"version":"v1","route":"/myFancyRoute", "name":"New Plugin","serviceguid":"19bafa36-8415-44e9-a2de-32852fefa6ef","routeguid":"30959362-e35b-4eb4-afe5-e8fdd6a3fecd"}
```

## Configuration

### Environment variables
1. PLUGIN_PORT_NAME_POSTFIX - The port of the service is defined by postfix of service port name. See [Plugin Deployment chart](https://gitlab.eclipse.org/eclipse/xfsc/personal-credential-manager-cloud/plugins/deployment/-/blob/main/deployment/helm/templates/service.yaml?ref_type=heads)
2. PLUGIN_TAGS - The tags Kong associates with plugins
3. PLUGIN_HTTP_METHODS - The http methods forwarded to the plugin by Kong
4. KONG_ADMIN_API - The url of the kong admin api
5. NAMESPACE - The namespace of the kubernetes cluster where to search for new plugins
6. KUBE_FILE - The path to the kubeconfig file
7. KUBE_CLUSTER_URL - url of remote kubernetes cluster
8. OVERRIDE_INCLUSTER - If set to true, the operator will enforce the configuration from provided KUBE_FILE and KUBE_CLUSTER_URL