# Acceptance tests for the tektoncd-catalog

## Requirements

This test suite assumes that Red Hat OpenShift Pipelines Operator is installed on the cluster, in the
`openshift-operators` namespace.

The system executing the tests must have following tools installed:

* `kuttl` kubectl plugin (v0.11.1 at time of writing)
* `oc` and `kubectl` client
* `jq` for parsing JSON data
* `curl` for interacting with some things

There should be a `kubeconfig` pointing to your cluster, user should have full
admin privileges (i.e. `kubeadm`).

## Environments

### `Regular OCP clusters`

The suites should work fine on regular OCP clusters (PSI/AWS/GCP)

### `Disconnected clusters`

The suites don't work flawlessly on disconnected clusters, for reasons unknown
yet.

To get started, you need to at least whitelist `.docker.io` and `.docker.com`
domains in the disconnected proxy's configuration.

### `How to get proxy host`

```bash
oc edit proxy cluster
```
Or

After successful disconnected cluster spin up, you can get proxy host details from the flexy job artifacts.

### `How to ssh into disconnected cluster proxy host`
1. Click on the link - https://github.com/openshift/shared-secrets/blob/master/aws/openshift-qe.pem
2. Click Raw option of openshift-qe.pem file
3. Open terminal and execute
4. wget <copy the URL from Raw view of openshift-qe.pem file> -O ~/.ssh/openshift-qe.pem
5. chmod 400 ~/.ssh/openshift-qe.pem
6. ssh -i ~/.ssh/openshift-qe.pem ec2-user@<HOST>

### `Use squid.config file /srv/squid/etc/squid.conf for whitelisting urls.`
```bash
sudo vi /srv/squid/etc/squid.conf
```
```bash
acl whitelist dstdomain tagging.us-east-1.amazonaws.com route53.amazonaws.com ec2.us-east-2.amazonaws.com iam.amazonaws.com .s3.us-east-2.amazonaws.com elasticloadbalancing.us-east-2.amazonaws.com .apps.airgap45-amit.qe.devcluster.openshift.com ec2-18-222-179-45.us-east-2.compute.amazonaws.com .github.com .rubygems.org
```
append the domain you want to whitelist for example [.gitlab.com](https://about.gitlab.com/)

Save and restart squid service
```bash
sudo systemctl restart squid-proxy
```

However, disconnected clusters behave strangely, and some of the test cases
may fail right now.

## Test suites

The tests in the test suite can be executed in any order, but must be run 
sequentially - no two tests should run in parallel, because they may potentially 
affect each other. Also, they must not run parallel to tests from the other 
test suite.

The test suite can only be run against an Operator installed into
a cluster, because it will manipulate resources such as `Subscription` and
assert for certain behavior.

## Running the tests

In any case, you should have set up your `kubeconfig` in such a way that your
default context points to the cluster you want to test. You can use the
`oc login ...` command to set this up for you or `export KUBECONFIG=<path/to/kubeconfig>`

To run the test suite:

```
make acceptance-test
```

### `Running manual with kuttl`

To run the test suite:

```
kubectl kuttl test --artifacts-dir ./_output --config ./tests/kuttl-test.yaml
```
The name of the test is the name of the directory containing its steps and
assertions.

If you are troubleshooting, you may want to prevent `kuttl` from deleting the
test's namespace afterwards. In order to do so, just pass the additional flag
`--skip-delete` to above command.

## Writing new tests

### `Name of the test`

Each test comes in its own directory, containing all its test steps. The name
of the test is defined by the name of this directory.

The name of the test should be short, but expressive. The format for naming a
test is currently `<test ID>_<short description>`.

The `<test ID>` is the serial number of the test as defined in the Test Plan
document. The `<short description>` is exactly that, a short description of
what happens in the test.

### `Name of the test steps`

Each test step is a unique YAML file within the test's directory. The name of
the step is defined by its file name.

The test steps must be named `XX-<name>.yaml`. This is a `kuttl` convention
and cannot be overriden. `XX` is a number (prefixed with `0`, so step `1` must
be `01`), and `<name>` is a free form value for the test step.

There are two reserved words you cannot use for `<name>`:

* `assert` contains positive assertions (i.e. resources that must exist) and
* `errors` contains negative assertions (i.e. resources that must not exist)

Refer to the
[kuttl documentation](https://kuttl.dev/docs)
for more information.

### `Documentation`

Documentation is important, even for tests. You can should provide inline
documentation in your YAML files (using comments) and a `README.md` in your
test case's directory. The `README.md` should provide some context for the
test case, e.g. what it tries to assert for under which circumstances. This
will help others in troubleshooting failing tests.

### `Recipes`

`kuttl` unfortunately neither encourages or supports re-use of your test steps
and assertions yet.

Generally, you should try to use `assert` and `errors` declaration whenever
possible and viable. For some cases, you may need to use custom scripts to
get the results you are looking for.

#### `Scripts general`

Scripts can be executed in a `kuttl.dev/TestStep` resources from a usual test
step declaration.

Your script probably will retrieve some information, and asserts it state. If
the assertion fails, the script should exit with a code > 0, and also print
some information why it failed, e.g.

```yaml
apiVersion: kuttl.dev/v1beta1
kind: TestStep
commands:
- script: |
    # Get some piece of information...
    if test "$result" != "expected"; then
      echo "Expectation failed, should 'expected', is '$result'"
      exit 1
    fi
```

Also, you may want to use `set -e` and `set -o pipefail` at the top of your
script to catch unexpected errors as test case failures, e.g.

```yaml
apiVersion: kuttl.dev/v1beta1
kind: TestStep
commands:
- script: |
    set -e
    set -o pipefail
    # rest of your script
```
#### `Getting values of a resource's environment variables`

YAML declarations used in `assert` or `errors` files unfortunately don't handle
arrays very well yet. You will always have to specify the complete expectation,
i.e. the complete array.

If you are just interested in a certain variable, and don't care about the rest,
you can use a script similar to the following using `jq`. E.g. to get the value
of a variable named `FOO` for the `argocd-server` deployment in the test's
namespace:

```yaml
apiVersion: kuttl.dev/v1beta1
kind: TestStep
commands:
- script: |
    val=$(oc get -n $NAMESPACE deployments argocd-server -o json \
      | jq -r '.spec.templates.spec.containers[0].env[]|select(.name=="FOO").value')
    if test "$val" != "bar"; then
      echo "Expectation failed for for env FOO in argocd-server: should 'bar', is '$val'"
      exit 1
    fi
```