# notificationsloadtest

[![orca-service](https://img.shields.io/badge/orca-service-blue.svg?style=flat)](https://orca.ethos.corp.adobe.com/services)
[![moonbeam](https://img.shields.io/badge/ethos-moonbeam-yellow.svg?style=flat)](https://moonbeam.ethos.corp.adobe.com/dc/notifications_load_test)

This is a simple Go application that listens on port 8080, with the apis detailed below.

### Build Container

Building the container is a multi-step process. To learn more about this refer to the following [wiki page](https://wiki.corp.adobe.com/display/CTDxE/make+build+target).

##### Local development

Please refer to the [Local Development](https://wiki.corp.adobe.com/display/CTDxE/DxE+-+Anonymous+access+removal+in+Artifactory#DxE-AnonymousaccessremovalinArtifactory-LocalDevelopment) section of the Artifactory authentication wiki for instructions on setting the ARTIFACTORY_USER and ARTIFACTORY_API_TOKEN environment variables before running the below commands. Your generated service is already configured for Artifactory authentication and needs no changes, but the remainder of that wiki contains more details on how the authentication works in Ethos.

Note: Commit Gopkg.lock and vendor folder to git repo after building image for the first time.

##### Mac users

```
Note: set ARTIFACTORY_USER and ARTIFACTORY_API_TOKEN before running below command
make build
```

##### Windows users

<pre>
Please refer to the <a href="https://wiki.corp.adobe.com/display/CTDxE/DxE+-+Anonymous+access+removal+in+Artifactory#DxE-AnonymousaccessremovalinArtifactory-LocalDevelopment" title="Local Development">Local Development</a> section of the Artifactory authentication wiki for instructions on setting the ARTIFACTORY_USER and ARTIFACTORY_API_TOKEN environment variables before running the below commands. Your generated service is already configured for Artifactory authentication and needs no changes, but the remainder of that wiki contains more details on how the authentication works in Ethos.

docker login -u $(ARTIFACTORY_USER) -p $(ARTIFACTORY_API_TOKEN) docker-asr-release.dr.corp.adobe.com
docker build -t notificationsloadtest-img -f Dockerfile .
</pre>

### Run Container

##### Using docker compose

Docker compose is a convenient way of running docker containers. For launching the application using docker-compose, ensure the builds steps are already executed.

```
docker-compose up --build
```

##### Using docker command

```
docker run --rm -it -e ENVIRONMENT_NAME=<dev|cd|qa|sqa|stage|prod|local> \
                    -e REGION_NAME=<ap-south-1|ap-southeast-1|ap-southeast-2|ap-northeast-1|ap-northeast-2|eu-central-1|eu-west-1|sa-east-1|us-east-1|us-west-1|us-west-2|local> \
                    -p 8080:8080 notificationsloadtest-img
```

To run container locally use:

```
docker run --rm -it -e ENVIRONMENT_NAME=local -e REGION_NAME=local -p 8080:8080 notificationsloadtest-img
```

To view container logs, run following command:

```
docker logs <container id>
```

Docker clean room setup:

To ensure that we're starting fresh (useful when you're doing a training session and/or trying to debug a local set up), it's best that we start with a 'clean room' and purge any local images and volumes that could introduce any potential 'contaminants' in our setup. You can read more on the following [wiki](https://wiki.corp.adobe.com/x/khu5TQ). Here is the command for docker clean room setup:

```
make clean-room
```
### Unit testing

Testing is done with [Ginkgo](https://onsi.github.io/ginkgo/) and the [`net/http/httptest`](https://golang.org/pkg/net/http/httptest/) package.

To run unit test and code coverage simply type:

```
make test
```

### List of available APIs

API | Description
--- | ---
`GET /ping` | Returns string 'pong'. Used for basic healthcheck.
`GET /notificationsloadtest/myfirstapi` | Returns Hello World message.
`GET /notificationsloadtest/error` | Returns the sample error response.

API can be accessed via curl command: `curl http://localhost:8080/<API>`

### Manage dependencies

Despite what's talked about in [Key Takeaways](https://golang.github.io/dep/docs/daily-dep.html#key-takeaways)
there is a [gotcha in the order of operations](https://github.com/golang/dep/issues/1670)
that prevent updates to `Gopkg.toml`. Adding the import first then running `dep ensure` _won't_ correctly add the dependency to `Gopkg.toml`. Commit all the changes you see into git. This ensures when someone else clones your service and tries to run it they have the exact dependencies they need.

To add a new dependency simply type:

```
dep ensure -add <dependency>
```

To update any dependency simply type:

```
dep ensure -update <dependency>
```


### References

  * Base image: https://git.corp.adobe.com/ASR/bbc-factory/blob/master/README.md


