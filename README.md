# helm schema gen plugin
This is a fork from knechtionscoding/helm-schema-gen, with added options to make use of Helm schema validation capability.

![](https://github.com/KnechtionsCoding/helm-schema-gen/workflows/goreleaser/badge.svg)

So that you don't have to write values.schema.json by hand from scratch for your Helm 3 charts

[Helm](https://helm.sh) plugin to generate [JSON Schema for values yaml](https://helm.sh/docs/topics/charts/#schema-files)

## Code stuff

Nothing fancy about the code, all the heavy lifting is done by:

- [go-jsonschema-generator](https://github.com/KnechtionsCoding/go-jsonschema-generator) - for generating JSON schema. It's a fork of [this](https://github.com/mcuadros/go-jsonschema-generator). Thanks to [@mcuadros](https://github.com/mcuadros)
- [go-yaml](https://github.com/go-yaml/yaml/) - for YAML parsing
- [cobra](https://github.com/spf13/cobra) - for CLI stuff
- [The Go stdlib](https://golang.org/pkg/) - for everything else

## Install

The plugin works with both Helm v2 and v3 versions as it's agnostic to the Helm
binary version

```
$ helm plugin install https://github.com/KnechtionsCoding/helm-schema-gen.git
KnechtionsCoding/helm-schema-gen info checking GitHub for tag '0.0.4'
KnechtionsCoding/helm-schema-gen info found version: 0.0.4 for 0.0.4/Darwin/x86_64
KnechtionsCoding/helm-schema-gen info installed ./bin/helm-schema-gen
Installed plugin: schema-gen
```

But note that the schema feature is present only in Helm v3 charts, so Helm
chart still has to be v3, meaning - based on the Helm chart v3 spec. And the
schema validation is only done in Helm v3. Read more in the
[Schema Files](https://helm.sh/docs/topics/charts/#schema-files) section of the
Helm official docs.

## Usage

The plugin works with both Helm v2 and v3 versions

Let's take a sample `values.yaml` like the below

```
replicaCount: 1

image:
  repository: nginx
  pullPolicy: IfNotPresent

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # The name of the service account to use.
  # If not set and create is true, a name is generated using the fullname template
  name:

podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: chart-example.local
      paths: []
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

nodeSelector: {}

tolerations: []

affinity: {}
```

Now if you use the plugin and pass the `values.yaml` to it, you will
get the JSON Schema for the `values.yaml`

```
$ helm schema-gen values.yaml
{
    "$schema": "http://json-schema.org/schema#",
    "type": "object",
    "properties": {
        "affinity": {
            "type": "object"
        },
        "fullnameOverride": {
            "type": "string"
        },
        "image": {
            "type": "object",
            "properties": {
                "pullPolicy": {
                    "type": "string"
                },
                "repository": {
                    "type": "string"
                }
            }
        },
        "imagePullSecrets": {
            "type": "array"
        },
        "ingress": {
            "type": "object",
            "properties": {
                "annotations": {
                    "type": "object"
                },
                "enabled": {
                    "type": "boolean"
                },
                "hosts": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "host": {
                                "type": "string"
                            },
                            "paths": {
                                "type": "array"
                            }
                        }
                    }
                },
                "tls": {
                    "type": "array"
                }
            }
        },
        "nameOverride": {
            "type": "string"
        },
        "nodeSelector": {
            "type": "object"
        },
        "podSecurityContext": {
            "type": "object"
        },
        "replicaCount": {
            "type": "integer"
        },
        "resources": {
            "type": "object"
        },
        "securityContext": {
            "type": "object"
        },
        "service": {
            "type": "object",
            "properties": {
                "port": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "serviceAccount": {
            "type": "object",
            "properties": {
                "create": {
                    "type": "boolean"
                },
                "name": {
                    "type": "null"
                }
            }
        },
        "tolerations": {
            "type": "array"
        }
    }
}
```

You can save it to a file like this

```
helm schema-gen values.yaml > values.schema.json
```
## Extra Helm validation options
You can add comments into your values.yaml file to automatically add constraint to your value.schema.json file. Which can be usefull if using this in a CI.

 Available options are :
|option|description|
|-----|-----|
|schemaRegex|Adds a `pattern` condition to your key|
|schemaEnum|Allows to limit the number of available string (must be comma-separated, no spaces)|
|schemaMinimum|Minimum value for an integer|
|schemaMaximum|Maximum value for an integer|
|schemaForceType|Override detected type (useful when you do not want default values set, but a correct schema). Accepted values are: null, string, boolean, integer, number, object, array|
**Example:**
```yaml
## @schemaEnum image.pullPolicy Always,Never,IfNotPresent
## @schemaForceType image.tag string
image:
  repository: nginx
  pullPolicy: IfNotPresent
  tag:
## @schemaRegex host.fqdn ^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$
host:
  fqdn: this.is.ok
## @schemaMinimum service.port 1024
## @schemaMaximum service.port 30000
service:
  port: 8080
```
## Issues? Feature Requests? Proposals? Feedback?

Put them all in [GitHub issues](https://github.com/KnechtionsCoding/helm-schema-gen/issues)
