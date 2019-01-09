# PUP
Maybe it's the wrong name but it's working pretty well. It's a aws-based *serverless* architecture, that let's you publish static content secured via passwords and usernames.

## Compnents
### Client
In order to make it work you'll need a config file `pup.json`. The client can run in a CI or on your local machine. You need to setup AWS Keys with the right permissions. You'll need read permissions for dynamodb and read/write permission to the s3 bucket.

**pup.json**
```
{
	"version": 1,
	"customer-code": "bio",
	"project": "test-project",
	"root": "./dist",
	"users": [
		{
			"username": "interal-user",
			"password": "supersecret",
			"groups": [
				"internal-group"
			]
		},
		{
			"username": "external-user",
			"password": "supersecret",
			"groups": [
				"external-group"
			]
		}
	],
	"rules": [
		{
			"/*/resources/**/*": {
				"group-permissions": [
					"public"
				]
			}
		},
		{
			"/master/**/*": {
				"group-permissions": [
					"external-group"
				]
			}
		},
		{
			"/**/*": {
				"group-permissions": [
					"dev",
					"internal-group"
				]
			}
		}
	],
	"aws-profile": "personal",
	"aws-bucket": "test-zip-bucket",
	"aws-region": "eu-central-1",
	"aws-dynamodb-logs-table": "test-logs-table"
}
```


### Worker
The worker is a docker container that runs in AWS Fargate. The worker is responsible to unpack files, put them in the right place, update users and whatever else is needed.

## Setup
0. Install terraform
1. Create a route53 Zone for your desired *root-domain*
2. Issue a wildcard certficate for that domain in us-east-1
3. Update `infrastructure/main.tf` to match your settings
4. Run `terraform init` to setup everything
5. Run `terraform apply` to setup aws
6. Add the following to your *logs-table*
```
{
  "baseUri": "/login",
  "r": [
    {
      "/**/*": {
        "groups": [
          "public"
        ]
      }
    }
  ]
}
```
7. Rerun `terraform apply` to initialize lambda edge the right way.


## Development
### Build Worker
```
cd worker
./build-linux.sh
docker build -t jurekbarth/worker:v0.0.1 .
docker push jurekbarth/worker:v0.0.1
```


## Todo
1. Better verification of configs

