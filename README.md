### Cloud Foundry Scheduler Service Broker

A Cloud Foundry Service Broker that acts as a scheduler. You can feed it schedules, jobs (cf tasks) and calls, and it will execute them.  
See also [Open Service Broker API](https://github.com/openservicebrokerapi/servicebroker/blob/v2.16/spec.md)


### Intro

The configuration for the template broker consists of the following environment variables:
* **SSB_DEBUG** - Debugging on or off, default is false.
* **SSB_HTTP_TIMEOUT** - Timeout in seconds for connecting to UAA endpoint, default is 10.
* **SSB_CLIENT_ID** - The uaa client to use for logging in, should have ??? scope.
* **SSB_CFAPI_URL** - The URL where to reach the cf api.
* **SSB_SCHEDULER_ENDPOINT** - The URL where to reach the scheduler rest endpoint (i.e. https://scheduler.sys.<cf domain>). This is used by the scheduler-plugin.
* **SSB_BROKER_USER** - The userid for the broker (should be specified issuing the _cf create-service-broker_ cmd).
* **SSB_CATALOG_DIR** - The directory where to find the cf catalog for the broker, the directory should contain a file called catalog.json.
* **SSB_LISTEN_PORT** - The port that the broker should listen on, default is 8080.
* **SSB_DB_NAME** - The name of the database to be used by the broker (usually this database is on the cf deployment's database VMs), default is ``schedulerdb``
* **SSB_DB_USER** - The user to use while connecting to the database.

The following are properties to be set in credhub, do this by creating a credhub service instance, and binding the scheduler-service-broker app to it:
* ``cf create-service --wait credhub default scheduler-service-broker-credentials -c '{ "SSB_CLIENT_SECRET": "secret1", "SSB_BROKER_PASSWORD": "secret2" , "SSB_DB_PASSWORD": "secret3" }'``
* ``cf bind-service scheduler-service-broker scheduler-service-broker-credentials``
* **SSB_CLIENT_SECRET** - The password for SSB_CLIENT_ID.
* **SSB_BROKER_PASSWORD** - The password for the broker (should be specified issuing the _cf create-service-broker_ cmd).
* **SSB_DB_PASSWORD** - The password for SSB_DB_NAME

Besides the broker protocol it also provides REST endpoints for handling schedules, jobs, calls and histories (basic CRUD operations), these endpoints are called by the scheduler-plugin.  
In the background it also runs a routine that checks the existing schedules to see if it needs to run cf tasks, or needs to call URLs.

### Deploying/installing the broker

## Prepare the database

Do a bosh login into one of the cf database servers (_bosh -d cf ssh database/0_), and switch to root.
Do a _mysql -u scheduler -p -h 127.0.0.1 --database=scheduler_, you will be prompted for the password, get this from credhub entry **/bosh/cf/scheduler_database_password**

# Update cf-deployment
We recommend using the cf mysql database.  
Add the following operator to cf-deployment to seed the db, create the dbuser and generate a database password:
```yaml
- type: replace
  path: /instance_groups/name=database/jobs/name=pxc-mysql/properties/seeded_databases/-
  value:
    name: scheduler
    password: ((scheduler_database_password))
    username: scheduler
- type: replace
  path: /variables/-
  value:
    name: scheduler_database_password
    type: password
```

# Deploy the broker 
Make sure the broker itself runs (most probably as a cf app), and the URL is available to the Cloud Controller.
Then install the broker:
```
#  the user and password should match with the user/pass you use when starting the broker app
cf create-service-broker scheduler-service-broker scheduler-service-broker-user scheduler-service-broker-password https://scheduler-service-broker.apps.cfd04.aws.rabo.cloud
```
Give access to the service(s) (all plans to all orgs):
```
cf enable-service-access scheduler
```
