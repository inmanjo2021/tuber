# tuber

![logo](logo.png)

Tuber is a CLI and server app for continuous delivery and cluster interactability.

So what does that mean? Server side:

- New services deployed to GCP do so through tuber
- Change config *inside your repo* for *how* it is deployed
- Canary releases that monitor sentry
- Staging and prod versions of your service just slightly different, without drift
- Autodeploy off certain branches on staging and prod
- Automatic rollbacks of all your resources if just one of them fails during in a release
- Monitor different things per resource while deploying
- Automatic cleanup of resources you no longer need
- Automigrate rails apps
- Easily adjust the levers of your resources
- Slack channels *per-app* for release notifications
- Review apps
- Dashboard for doing all of the above

&nbsp;

Here's the CLI side (most useful commands at least):
- `switch` - `tuber switch staging` points you at a cluster
- `context` - tell what cluster you're currently working with locally
- `auth` - log in and get credentials (most commands use the running tuber's graphql api)
- `exec` - interactive command runner (rails c, etc)
- `rollback` - we cache whatever was applied in the last successful release, this reapplies that.
- `pause` - pause releases, `resume` resumes
- `apps info` for everything tuber knows about your app
- `apps set` - subcommands for every field
- `deploy -t` - deploy a specific image rather than pulling latest
- `env` - get, set, unset, and list env vars
- `apps install` - deploy a new service in seconds.

&nbsp;

**all commands and subcommands** support `-h`/`--help`, and we **highly** encourage using that as you learn your way around.

&nbsp;

We do all of this with NO esoteric custom resources, and ALL of tuber's interaction with the cluster occurs through `kubectl`.

You can run our entire release process yourself just by copying the kubectl commands in the debug logs.

You need only the Kubernetes (and ok fine Istio) documentation to understand a tuber app's footprint.

So, welcome to "kubectl as a service" :tada:

&nbsp;

&nbsp;

# Project Status
We are pushing for Tuber to be more generally applicable, but it currently makes too many assumptions about how things _outside_ your cluster are configured.

Until then, Tuber remains closer to "open code" than an "open source solution".

&nbsp;

# What's a Tuber App?
```
- TuberApp
  - Name
  - ImageTag
  - Vars
  - State
    - Current
    - Previous
  - ExcludedResources
  - SlackChannel
  - Paused?
  - ReviewApp?
  - ReviewAppsConfig
    - ReviewAppsEnabled?
    - Vars
    - ExcludedResources
```

So that's our data model. This is all stored for each app, in an in-memory database loaded in as a local file to the running tuber server.

Below you'll find explanations for some of the less intuitive aspects of it all.

These should also give some context as to "why tuber is the way that it is" - these are crucial to how we can handle any app in any env:

&nbsp;

### Vars
**Vars** are custom interpolation variables. We use these all over the monolith's resources.

These are internally referred to as "App-Specific Interpolatables", or "ASIs" to the chagrin of everyone involved. It's accurate though.

Tuber offers the following Vars to every app's `.tuber/` resources automatically:
```
{{ .tuberImage }}             <- the full image tag for a release
{{ .tuberAppName }}           <- the tuber app name (vital for review apps)
{{ .clusterDefaultHost }}     <- default hostname to offer to virtualservices
{{ .clusterDefaultGateway }}  <- default gateway powering the default host
{{ .clusterAdminHost }}       <- secondary hostname to offer to virtualservices, ideal for an IAP domain
{{ .clusterAdminGateway }}    <- secondary gateway powering the secondary host
```

Those are the Go format for string interpolation (like ruby's `"#{hi}"`), and are hard-coded like that in the resources in `.tuber/`.

Tuber interpolates those in every release based on its own context.

So if these "interpolatables" help us do review apps, and differentiate Staging vs Production, what's a Var?

It's anything *specific to your app* that needs interpolation to support reuse for review apps or different environments.

A clear example is something like `{{ .sidekiqWorkers }}` different on staging and prod.

You can make a Var for anything, and interpolate anything, including booleans and integers.

&nbsp;

&nbsp;

### ExcludedResources
Sometimes Vars don't cut it. Sometimes you just need to cut entire resources out of a specific environment, or from a review app.

**ExcludedResources** is a hash of kubernetes Kind (**Deployment**, **CronJob**, etc), mapped to the Name of a resource.

When your app is deployed, any resources contained in `.tuber/` matching these Kind/Name pairs will be skipped.

We also interpolate Exclusion names prior to comparison, so a resource coded as `name: "{{.tuberAppName}}-canary"` can be excluded with the same name, to support excluding that resource on review apps.

It's a lot of setup, but that's the point. 

If this makes for a bunch of work, you are doing that work to *actively* make staging different from prod. 

We WANT that to be a bunch of work. If staging is different from production, it should be EXPLICITLY different.

&nbsp;

&nbsp;

### ReviewAppsConfig
This specifies how an app's review apps should be created.

The field has its own **Vars** and **ExcludedResources**, which are *merged with* the top-level **Vars** and **ExcludedResources**.

Let's say you want `myapp-demo` to exclude `deployment` `myapp-canary`, and review apps created off of it to ALSO exclude `deployment` `myapp-sidekiq`?

This is how. And the final merged lists of **Vars** and **ExcludedResources** are set on the review app, so if you then have a review app that needs to change something about `deployment` `myapp-sidekiq`, you can DE-exclude it from that review app specifically, after it's created.

The Vars are the most useful here - for example, letting you specify lower CPU limits and memory limits for an app's review apps.

&nbsp;

&nbsp;


# Installation

### Prerequisites
- `gcloud`
- `kubectl`

### Homebrew
```
brew tap freshly/taps
brew install tuber
```

### Scoop
```
scoop bucket add freshly https://github.com/freshly/scoops
scoop install tuber
```

### Download binary
Download the binary file from the latest release: https://github.com/Freshly/tuber/releases/
