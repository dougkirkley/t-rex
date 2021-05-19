# T-Rex
cli tool for creating serverless api using terraform and aws lambda

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

You will need AWS env variables added to the CI variables for the project

The AWS credentials used during testing was the administrator group policy in AWS.

### Installing

To install for local development on Mac via Homebrew
```
brew install awscli
brew install go
```

### Deployment

`trex init` will build an example yaml file for deployment in the current dir

`trex deploy` will deploy the example project in your aws account.

in the root directory of the project.

## Authors

* **Douglass Kirkley** - *Initial work* - [trex](https://github.com/dougkirkley)
