# ndv (N's dev-tools)

Personal dev tools configurations for later reuse

## Prerequisites

- yq

## Features

Configure AWS MFA session

- Create `config_data.yml`. Here is an example

``` yaml
aws:
  mfas:
    - profile: prophet # name of the profile to be authenticated, with access key configured
      device: arn:aws:iam::xxxx:mfa/user.name # MFA device ARN
      session_duration: 129600
      output_profile: prophet-eks # name of the profile to be created
```

- Run `make aws_mfa TOKEN=123456`
