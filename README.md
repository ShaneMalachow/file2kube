# file2kube: Convert files to Kubernetes objects

file2kube provides a command-line utility to take existing configuration files and convert them automatically to base64 encoded Secret files or include them as plaintext in ConfigMap files. This way whenever you need to include these files in your GitOps based repos, you can just generate the config files automatically. This is especially helpful when using operators like Mozilla sops that do not support encrypting whole files (such as RSA keys). This does the heavy lifting of converting it to a Kubernetes resource definition and then you can quickly use your tools to encrypt those definitions.

## Installation

`go get -u github.com/shanemalachow/file2kube`

## Usage

Convert files into base64 encoded Secret:
`file2kube secret [FILES...]`

Convert files into ConfigMap:
`file2kube configmap [FILES...]`
