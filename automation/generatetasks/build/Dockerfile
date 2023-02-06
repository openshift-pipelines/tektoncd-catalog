# FIXME Long term, use a better image
FROM fedora:35

# FIXME: we don't need to set this HOME env variable
ENV HOME=/go/src/github.com/openshift-pipelines/tektoncd-catalog/ \
    USER_UID=1001 \
    USER_NAME=generatetasks


RUN curl -L https://go.dev/dl/go1.18.1.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH=$PATH:/usr/local/go/bin

WORKDIR ${HOME}
ADD . .
COPY automation/generatetasks/build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

RUN dnf install -y jq make diffutils which git gcc && rm -rf /var/cache/yum /var/cache/dnf

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
