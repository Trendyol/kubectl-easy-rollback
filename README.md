![alt text](https://img.icons8.com/nolan/2x/time-machine.png) 

# EasyRollback

EasyRollback is aim to easy rollback to previous images that  deployed on k8s environment
## Installation

You should have go installation first [go](https://golang.org/dl/) to install Golang.
For OSX

```bash
brew install go
```
Then get project

```bash
go get -v github.com/Trendyol/easy-rollback
```

Or you can easily download project via brew (OSX package manager)
```bash
brew tap trendyol/trendyol-tap
brew install easy-rollback
```

Or you can install the pre-compiled binary
```bash
wget https://github.com/Trendyol/easy-rollback/releases/download/v1.0.10/easy-rollback_1.0.10_Linux_x86_64.tar.gz
tar xvfz easy-rollback_1.0.10_Linux_x86_64.tar.gz
mv easy-rollback /usr/local/bin
```

## Usage
Project look at your .kube/config file to read current-context configs hence of you should have kubernetes environment configurations inside .kube/config.

```bash
easy-rollback list --deployment <deployment> --namespace <namespace> --> Will list all of your previous deployed images.
easy-rollback rollback --to-image <image> --deployment <deployment> --namespace <namespace> --> Will rolback your deployment to given image.
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
