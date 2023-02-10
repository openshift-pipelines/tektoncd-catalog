# FIXME Long term, use a better image
FROM fedora:35

# FIXME: we don't need to set this HOME env variable
ENV HOME=/go/src/github.com/openshift-pipelines/tektoncd-catalog/ \
    USER_UID=1001 \
    USER_NAME=generatetasks


RUN curl -L https://go.dev/dl/go1.18.1.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH=$PATH:/usr/local/go/bin

ARG TKN_VERSION=0.29.1
RUN ARCH=$(uname -m) \
    OS=$(uname) \
    && curl -L https://github.com/tektoncd/cli/releases/download/v${TKN_VERSION}/tkn_${TKN_VERSION}_${OS}_${ARCH}.tar.gz | tar zxf - -C /usr/local/bin

WORKDIR ${HOME}
ADD . .
COPY automation/generatetasks/build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

RUN dnf install -y jq make diffutils which git gcc && rm -rf /var/cache/yum /var/cache/dnf

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
