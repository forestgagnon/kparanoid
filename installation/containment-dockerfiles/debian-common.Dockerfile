ARG DEBIAN_BASE_IMAGE=debian:stretch-slim
FROM ${DEBIAN_BASE_IMAGE}

RUN apt-get update && apt-get install -y \
  git \
  apt-transport-https \
  curl \
  jq \
  watch \
  ca-certificates \
  vim \
  bash-completion
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install

RUN curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg && \
  echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list

ARG KUBECTL_VERSION=1.21.1-00
RUN apt-get update && apt-get install -y kubectl=${KUBECTL_VERSION}

COPY .bash_profile /root/
ENTRYPOINT ["/bin/bash", "-l"]
