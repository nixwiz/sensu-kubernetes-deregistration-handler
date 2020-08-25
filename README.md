# Sensu Go Kubernetes Deregistration Handler
[![Bonsai Asset Badge](https://img.shields.io/badge/Sensu%20Kubernetes%20Deregistration%20Handler-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/sensu/sensu-kubernetes-deregistration-handler)

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Handler definition](#handler-definition)
  - [Environment variables](#environment-variables)
  - [RBAC](#rbac)

## Overview

sensu-kubernetes-deregistration-handler is a [Sensu Handler][2] for
deregistering entities created by the [Kubernetes Events Check][1].  This
handler is not expected to be assigned to that, or any otther, check, however.
That check creates events using the agent API and when it detects a Kubernetes
event that necessitates the deletion of an associated Sensu entity, it
creates an event with a check named `kubernetes-delete-entity` and with this
handler specified.

## Usage Examples

Help:
```
Sensu Handler to deregister entities created by the Sensu Kubernetes Events check

Usage:
  sensu-kubernetes-deregistration-handler [flags]
  sensu-kubernetes-deregistration-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -k, --apiKey string            Sensu Go Backend API Key
  -u, --apiURL string            Sensu Go Backend URL (default "http://127.0.0.1:8080")
  -t, --trusted-ca-file string   TLS CA certificate bundle in PEM format
  -i, --insecure-skip-verify     skip TLS certificate verification (not recommended!)
  -h, --help                     help for sensu-kubernetes-deregistration-handler

Use "sensu-kubernetes-deregistration-handler [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][3] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```bash
sensuctl asset add sensu/sensu-kubernetes-deregistration-handler
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][4]

### Handler definition

```yml
---
type: Handler
api_version: core/v2
metadata:
  name: kubernetes-deregistration
  namespace: default
spec:
  command: sensu-kubernetes-deregistration-handler
  type: pipe
  runtime_assets:
    - sensu/sensu-kubernetes-deregistration-handler
  secrets:
    - name: SENSU_API_KEY
      secret: sensu_api_key
```

### Environment variables

Many arguments for this handler are available to be set via environment
variables.  However, any arguments specified directly on the command line
override the corresponding environment variable.

|Argument      |Environment Variable  |
|--------------|----------------------|
|--apiURL      |SENSU_API_URL         |
|--apiKey      |SENSU_API_KEY         |

Given the sensitivity of the API key, it is advised to use
[secrets management][5] to surface it.  The [handler definition](#handler-definition)
above references it as a secret.  Below is an example secret definition that
makes use of the built-in [env secrets provider][6].

```yml
---
type: Secret
api_version: secrets/v1
metadata:
  name: sensu_api_key
spec:
  provider: env
  id: SENSU_API_KEY
```

### RBAC

It is advised to use [RBAC][8] to create a user scoped specifically for
purposes such as this handler and to not re-use the admin account.  For
this handler, in particular, the account would need access to list and
delete entities.  The example below shows how to create a limited-scope user
and the necessary role and role-binding resources to give it the required
access.

```
sensuctl user create kubernetes --password='4yva#ko!Yq'
Created

sensuctl role create delete-entities --verb list,delete --resource entity
Created

sensuctl role-binding create kubernetes-delete-entities --role=delete-entities --user=kubernetes
Created
```

This handler only supports the use of an [API key][9] for accessing the API.
You can create the API key with sensuctl:

```
sensuctl api-key grant kubernetes
Created: /api/core/v2/apikeys/03f66dbf-6fe0-40d4-8174-95b5eab95649
```

[1]: https://github.com/sensu/sensu-kubernetes-events
[2]: https://docs.sensu.io/sensu-go/latest/reference/handlers/
[3]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[4]: https://bonsai.sensu.io/assets/sensu/sensu-kubernetes-deregistration-handler
[5]: https://docs.sensu.io/sensu-go/latest/guides/secrets-management/
[6]: https://docs.sensu.io/sensu-go/latest/guides/secrets-management/#use-env-for-secrets-management
[7]: https://docs.sensu.io/sensu-go/latest/reference/events/#metrics
[8]: https://docs.sensu.io/sensu-go/latest/reference/rbac/
[9]: https://docs.sensu.io/sensu-go/latest/reference/apikeys/

