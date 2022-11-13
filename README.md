[![Go Badge](https://img.shields.io/badge/-Go-blue?style=flat&logo=Go&logoColor=white)](https://go.dev/)
[![Twitch Badge](https://img.shields.io/badge/-Internship%20Project-blueviolet?style=flat&logo=Twitch&logoColor=white)](https://www.twitch.tv/)

# My PagerDuty Wrapper
My PagerDuty wrapper (mPagerDuty) is an independent Go package that any other projects can leverage to create PagerDuty API requests using a specific authentication token

## Initialize a Client

You can retrieve a pagerduty-go client struct that is authenticated using PagerDuty token from the mPagerDuty package like this:

```go
mPD, err := mPagerDuty.GetMPagerDutyClient()
```

Once retrieved, you can access any of the receiver functions tied to the client, which is a [IMPagerDuty interface](./pkg/mPagerDuty.go#L13), like this:

```go
userID, err := mPD.GetUserIDbyName("Timur Kalandarov")
```

### Test Stub

The mPagerDuty package also implements an API stub that does not send live traffic data and instead returns static responses to function calls. The functions are implemented the same way as the live functions, so parameters and return objects will be exactly the same, but responses behave predictably given certain arguments and return known responses. See the [mPagerDuty_fake.go](./pkg/mPagerDuty_fake.go) file for the stub function implementations.

You can retrieve a faked client like this:

```go
fmPD := mPagerDuty.FakePDClient{}
```

**Note:** A faked client is automatically returned from the `mPagerDuty.newMPagerDutyClient()` if either of the following environment variables are set in the environment where you're running Mercy:

```Go
RUNNING_IN_JENKINS=true
LOCAL_DEV_TESTING=true
```

With this automatic fake client based on environment variables, a single test case can be written for a function instead of separate tests for Unit and Integration purposes. Separate tests may still be best practice, and combined-function tests may not necessarily behave predictably depending on how (non) robust the data returned by the fake client is.

## Adding Functions

You can add your own functions to this package to expand what it's capable of. If you want to hit a certain Jira endpoint that isn't covered by any existing function, implement it here instead of in your own package so everyone can take advantage of it. In the simplest sense, here's what you'll need to do to add a new function:

- Add the full function declaration, including parameters and return types, into the [IMPagerDuty interface](./pkg/mPagerDuty.go#L13)
- Implement the function so it's publicly exported, **taking care to properly make it a receiver function**, somewhere in [mPagerDuty.go](./pkg/mPagerDuty.go)

```go
// In the following function declaration, the `(c *client)` part makes it a receiver function
func (c *client) GetScheduleIDbyName(name string) (string, string, error)
```
- Inside of your function implementation, make calls to the PagerDuty API through the pagerduty-go SDK like this:
```go
schedules, err := c.pdClient.ListSchedulesWithContext(/*parameters*/)
```
- You can also makes calls to other functions inside of the mjira package itself like this:
```go
response, err := c.GetOnCalls(scheduleIDs)
```
- **Make sure to implement a stubbed version of your function**. If you do not, the package will break. Stubbed versions of your functions go in [mPagerDuty_fake.go](./pkg/mPagerDuty_fake.go) and have the following declration style (note the **(fakeClient \*FakePDClient)** specifically):
```go
func (fakeClient *FakePDClient) GetScheduleIDbyName(name string) (string, string, error)
```
- Ensure that the stubbed version does not call the actual PagerDuty API, but instead returns a static response that you create. You can see the existing functions for examples
