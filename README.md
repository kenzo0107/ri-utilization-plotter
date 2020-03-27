![logo](logo.png)

[![](https://github.com/kenzo0107/ri-utilization-plotter/workflows/static%20check/badge.svg)](https://github.com/kenzo0107/ri-utilization-plotter/actions?query=workflow%3A%22static+check%22)
[![](https://github.com/kenzo0107/ri-utilization-plotter/workflows/Test/badge.svg)](https://github.com/kenzo0107/ri-utilization-plotter/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/kenzo0107/ri-utilization-plotter)](https://goreportcard.com/report/github.com/kenzo0107/ri-utilization-plotter)
[![codecov](https://codecov.io/gh/kenzo0107/ri-utilization-plotter/branch/master/graph/badge.svg)](https://codecov.io/gh/kenzo0107/ri-utilization-plotter)
[![GoDoc](https://godoc.org/github.com/kenzo0107/ri-utilization-plotter?status.svg)](https://pkg.go.dev/github.com/kenzo0107/ri-utilization-plotter)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/kenzo0107/ri-utilization-plotter/master/LICENSE)

This project provides an AWS Lambda application that created and deployed by Serverless Framework for the following purpose:

* Plot below metrics of your AWS Account to Datadog
  - AWS Reserved Instance Utilization
  - AWS Reserved Instance Coverage

### Set SSM Parameter store with description in your AWS account

* datadog_api_key
* datadog_app_key

## Invoke Lambda Function in Local

```sh
make local-invoke
```

## Deploy Lambda Function

At first, execute `sam deploy --guided` and generate samconfig.toml.

From then, run the following command:

```sh
make deploy
```

## Invoke Lambda Function

```sh
aws lambda invoke --function-name ri-utilization-plotter --log-type Tail out.log
```

## LICENSE

[MIT License](https://github.com/kenzo0107/ri-utilization-plotter/blob/master/LICENSE)

## Note

Icon made by bqlqn from [www.flaticon.com](https://www.flaticon.com)
