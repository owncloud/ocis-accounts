Enhancement: Add moq to generate grpc server mock

We are now using [`github.com/matryer/moq`](http://github.com/matryer/moq) to generate a mock implementation of the service so that other services can use it in unit tests. See ocis-proxy for an example.

https://github.com/owncloud/ocis-accounts/103
