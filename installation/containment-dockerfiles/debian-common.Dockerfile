ARG DEBIAN_BASE_IMAGE=debian:bullseye-slim
FROM ${DEBIAN_BASE_IMAGE}

RUN apt-get update && apt-get install -y \
  git \
  apt-transport-https \
  curl \
  jq \
  watch \
  ca-certificates \
  vim \
  bash-completion \
  gnupg \
  google-cloud-sdk-gke-gcloud-auth-plugin

RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install

ARG KUBECTL_VERSION=1.27.10
RUN curl -LO https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl \
  && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

COPY .bash_profile /root/
ENTRYPOINT ["/bin/bash", "-l"]
