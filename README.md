# kparanoid - isolated kubectl environments for safer production cluster administration

kparanoid helps you avoid polluting your local environment with myriad
production kubectl contexts. Environments with multiple kubectl contexts are
dangerous, because if you forget to add `--context` even once to a kubectl
command, you can accidentally run a command on the wrong cluster. It happens.

A common "solution" to this problem is to put the active kubectl context in the
terminal prompt. For the paranoid among us, this is not good enough: kubectl is
stateful and the context may have changed since the prompt was drawn, perhaps
from another terminal window or even a script.

kparanoid takes an extreme approach to isolation, where all kubectl commands
take place in isolated, single-context docker containers. Kparanoid will error
if you give it a kubectl config with more than one context, allowing it to
display a kube context in the terminal prompt which you can have full confidence
in. If you daily drive kparanoid, you can delete production contexts from your
local `kubectl`, decreasing your risk of accidents.

There may be much simpler ways to achieve something similar, without the
usability issues of working inside a docker container. Admittedly, kparanoid is
partially optimized for coolness.

## Bonus features

- different kubectl versions for each cluster
- persistent shell history which is scoped to each cluster.

## Installation

You must have `docker` installed, and have an environment with `bash`. The
install script has been tested on Mac and Linux. It is distributed via the
public kparanoid docker image.


```shell
docker run --rm forestgagnon/kparanoid:1 install bash --config-dir="$HOME/.kparanoid" | bash
```

You can install kparanoid to any directory by changing `--config-dir`.

When the install script completes, it will print out a command which generates a
shell statement for adding `kparanoid` to your PATH.

### Building and installing from source

```shell
make build-cli && docker run --rm forestgagnon/kparanoid:1 install bash --config-dir="$HOME/.kparanoid" | bash
```

## Usage

```shell
kparanoid --help
```

### Adding a cluster

#### GKE Cluster

```
kparanoid add-cluster gke mycluster --gcloud-cluster-creds-cmd='gcloud container clusters get-credentials my-gke-cluster --region us-central1 --project some-gcp-project'
```

#### Vanilla cluster (bring-your-own kubeconfig file)

Any kubeconfig should work with vanilla mode unless it needs to shell out to
some other tool for authentication information, like how GKE kubeconfigs require
`gcloud` to be present.

```
kparanoid add-cluster vanilla some-cluster

> You must manually place the kubeconfig file at exactly '/home/forest/.kparanoid/clusters/some-cluster/container-env/.kube/config'
```

### Interacting with a cluster

```
# Start an interactive session
kparanoid open mycluster
kparanoid open my-old-cluster --kubectl-version=1.9.0-00

# Use exec to leverage your fancy local env for processing output
kparanoid exec mycluster kubectl get deployment some-deployed-thing -oyaml | yq e '.metadata.name' | cowsay

 _____________________
< some-deployed-thing >
 ---------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||

```

### Editing cluster config

You can edit the cluster config file to change the default kubectl version

```
vim $(kparanoid cluster-config get-filepath mycluster)
```

Available kubectl versions can be fetched for your flavor of cluster

```
kparanoid kubectl-versions vanilla
```

## Limitations

Currently, there are some things you can't really do well with kparanoid, or at
all:

- port-forwarding is broken
- commands which accept input from filepaths are annoying to work with
- piping input to `kparanoid exec` doesn't work

## Future work

- workarounds to allow `kubectl port-forward` to be used
- make it possible to pipe input to `kparanoid exec`
- make `kparanoid exec` usable with kubectl commands which expect filepath
  parameters
- customizable environments which don't get overwritten by the installer
