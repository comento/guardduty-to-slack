[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/comento/guardduty-to-slack.svg)](https://github.com/comento/guardduty-to-slack)
[![Go Report Card](https://goreportcard.com/badge/github.com/comento/guardduty-to-slack)](https://goreportcard.com/report/github.com/comento/guardduty-to-slack)

# GuardDuty To Slack

AWS GuardDuty to Slack using Serverless

## Getting Started

### Prerequisites

* Create `.env` file and set WEBHOOK_URL
* Set stage, region, profile, deploymentBucket according to personal settings in `serverless.yml`
* 

### Deploy

```
make build
sls deploy
```

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.

## Contact

tech@comento.kr


