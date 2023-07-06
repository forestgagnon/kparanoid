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
  gnupg
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && ~/.fzf/install

RUN mkdir -p /etc/apt/keyrings && \
  curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg && \
  echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list

RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
  curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add - && \
  apt-get update && apt-get install -y google-cloud-sdk-gke-gcloud-auth-plugin

ARG KUBECTL_VERSION=1.27.3-00
RUN apt-get update && apt-get install -y kubectl=${KUBECTL_VERSION}

COPY .bash_profile /root/
ENTRYPOINT ["/bin/bash", "-l"]
