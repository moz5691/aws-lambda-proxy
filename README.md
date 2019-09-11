## Serverless Start Guide for Golang.

### Simple three steps you can start instantly. Assume that Serverless Framework is already installed.

### 1. Create a project

```s
$ serverless create -t aws-go-dep -p <project_name>
```

### 2. Build

```s
$ make
```

### 3. Deploy

```s
$ serverless deploy
```

### If you need to tear down, Remove the whole deployment.

```s
$ serverless remove
```

### You can still invoke a function manually.

```s
$ serverless invoke -f <function>
```

If you have credentials stored in ~/.aws , you can choose one of profiles in it.

```s
$ ls ~/.aws
config          credentials

$ cat ~/.aws/config
[default]
region=us-east-1
output=json

[profile sls-admin]
region=us-east-1
output=json

$ cat ~/.aws/credentials
[default]
aws_access_key_id = ***********
aws_secret_access_key = **************

[sls-admin]
aws_access_key_id = ************
aws_secret_access_key = ***************

```

More tips for profile.

```s
$ aws configure --profile account1
$ aws configure --profile account2

$ aws dynamodb list-tables --profile account1
$ aws s3 ls --profile account2
```

If you want to use one of profiles as default, you could set it and test it.

```s
$ export AWS_DEFAULT_PROFILE=account1
$ aws dynamodb list-tables
```

### Add a profile, if not provided, default one is used.

```yaml
provider:
name: aws
runtime: go1.x
profile: sls-admin
```

### Seeding Data to Dynamo

```s
$ go run cmd/seed-dynamo/main.go --table-name="$YOUR_TABLE_NAME"
```

You can set the configured (set by default) profile as SharedConfigState in your code.
Again, if your chosen profile is not default, run the following to change to the one you need to use.

```s
export AWS_DEFAULT_PROFILE=account1
```

```go
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(session)
```

### Testing

Use Curl

```s
$ curl -X POST https://<apigw-url> -d <request data>
```

Run client

```s
$ make dev-crun
```
