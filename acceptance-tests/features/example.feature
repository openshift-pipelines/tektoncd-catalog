Feature: Assert Tekton Installation

    Assert the installation of Tekton and its dependencies in the Kubernetes cluster

    @automated
    Scenario: Asserting Tekton Installation : TKNECO-01-TC01
        Given the Kubernetes cluster is available
        Then "tekton-pipelines-controller" is deployed and "READY" on "tekton-pipelines" namespace
        And "tekton-pipelines-webhook" is deployed and "READY" on "tekton-pipelines" namespace
 