# Demo App: outerspace-go

Outerspace is a demo app that uses the SpaceX API to fetch data about SpaceX launches, rockets, and capsules. It also talks to a numbers API that consistently generates random results.

![outerspace-go](/img/outerspace-go.png)

# Working Locally

## How to run the go app locally

This will run a server on port `:8080` with the API endpoints.
```
go run main.go
```

Once the app is running you can make API requests.

The simplest way to make a bunch of requests is to run the build in client
script.
```
cd cmd/client
go run .
```

But you can also make requests directly.

```
curl localhost:8080/api/latest-launch
```

Request with no path to see the full list of API endpoints:
```
curl localhost:8080/ | jq
{
  "/": "Shows this list of available endpoints",
  "/api/latest-launch": "Get the latest SpaceX launch",
  "/api/numbers": "Get a random math fact",
  "/api/rocket": "Get a specific rocket by ID (use ?id=[rocket_id])",
  "/api/rockets": "Get a list of all SpaceX rockets"
}

```

## How to run the tests locally

There is a test file under `lib/api_testify_test.go` that loops through a set of API calls that were recorded using `proxymock` and replays each one. It then compares the response of each one to what was previously recorded. You can run it like so:
```
go test -v ./...
```

# proxymock

For apps that make lots of API calls, [proxymock](https://proxymock.io/) can be used to record, mock, and replay those downstream calls.

## Record with proxymock

All you have to do is set the following environment variables before you run the go code:
```
export http_proxy=http://localhost:4140
export https_proxy=http://localhost:4140
export grpc_proxy=http://$(hostname):4140
go run main.go
```

Then in a new terminal window run `proxymock run` to start recording:
```
proxymock run
```

Now, in yet another new terminal window, run the script to make requests to the server:
```
cd cmd/client
go run .
```

Note that this will only record requests made from the API server to the external APIs.  To record the requests from the client you will need to export the same proxy environment variables in that terminal.

For examples of the kind of data collected, check the `.speedscale` directory of the repository.

## Running tests with proxymock

You may find that because the Numbers API returns random values that the test usually fails. The solution for this is to run `proxymock` first and then run the tests. First import the existing data into `proxymock`, and take note of the snapshot id that is created.
```
proxymock import --file .speedscale/raw.jsonl
```

If you want to see what the data looks like, you can use the inspect command:
```
proxymock inspect snapshot "$SNAPSHOT_ID"
```

### API List

You can see a list of every API call that was recorded, both inbound to your application and outbound to downstream systems.

![proxymock](/img/inspect-list.png)

### Drill-down

Using your arrow keys you can navigate to a particular call and see all the details of what request was sent and how the application responded. Here we can see the exact value returned from the Numbers API.

![proxymock](/img/inspect-drill-down.png)

Once you are ready to use that data you can then run with that snapshot:
```
proxymock run --snapshot-id "$SNAPSHOT_ID"
```

Now in your window where you are going to run the tests, make sure to export the environment variables and run the test again:
```
export http_proxy=http://localhost:4140
export https_proxy=http://localhost:4140
go test -v ./...
```

You should see that all the tests are now passing because `proxymock` is mocking out the Numbers API with a consistent and repeatable result.
