![alt text](https://img.icons8.com/nolan/2x/time-machine.png) 

# EasyRollback

EasyRollback is aim to easy rollback to previous images that  deployed.
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

## Usage
Project look at your .kube/config file to read current-context configs hence of you should have kubernetes environment configurations inside .kube/config.

```bash
easy-rollback list --deployment <deployment> --namespace <namespace>
easy-rollback rollback --to-image <image> --deployment <deployment> --namespace <namespace>
```

## License
[MIT](https://choosealicense.com/licenses/mit/)
