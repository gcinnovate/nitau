# nitau

# About 

nitau is a Golang  adapter for the NITA-U Mobile Service Delivery Gateway (MSDG) 
used to handle both MT and MO SMS messaging.
# Deploying

As nitau is a go application, it compiles to a binary and that binary along with the config file is all
you need to run it on your server. You can find bundles for each platform in the
[releases directory](https://github.com/gcinnovate/nitau/releases). We recommend running nitau
behind a reverse proxy such as nginx or Elastic Load Balancer that provides HTTPs encryption.

# Configuration

nitau uses a tiered configuration system, each option takes precendence over the ones above it:
 1. The configuration file
 2. Environment variables starting with `NITAU_` 
 3. Command line parameters

We recommend running nitau with no changes to the configuration and no parameters, using only
environment variables to configure it. You can use `% nitau --help` to see a list of the
environment variables and parameters and for more details on each option.

#  Configuration

For use with the NITA-U gateway, you will want to configure these settings:

 * `NITAU_API_USER`: The NITA-U gateway userid for your account
 * `NITAU_API_PASSWORD`: The NITA-U gateway password 
 * `NITAU_API_EMAIL`: The NITA-U gateway email address for your account 
 * `NITAU_API_AUTH_TOKEN`: The NITA-U gateway authorization JSON Web Token for your account 
 * `NITAU_API_MO_URL`: The NITA-U gateway MO endpoint for posting MO messages from the gateway 
 * `NITAU_API_SMS_URL`: The NITA-U gateway SMS API endpoint for MT messaging 
 * `NITAU_API_ROOT_URI`: The root uri for the api endpoints for the gateway (default: `https://msdg.uconnect.go.ug/api/v1`)
 * `NITAU_SERVER_PORT`: The port on which to run the nitau service

 
# Development

Install nitau source in your workspace with:

```
go get github.com/gcinnovate/nitau
```

Build nitau with:

```
go install github.com/gcinnovate/nitau/..
```

This will create a new executable in $GOPATH/bin called `nitau`. 
