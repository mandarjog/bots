# Policy Bot

The policy bot enforces a number of Istio-wide policies around how we manage
issues, pull requests, test flakes, and more.

- [Overview](#overivew)
- [Plugins](#plugins)
- [Startup options](#startup-options)
- [Configuration file](#configuration-file)
- [Deployment](#deployment)
- [Credentials and secrets](#credentials-and-secrets)

## Overview

The bot consists of a single binary running in a container. Use `make container` to build both the binary and the container.

At runtime, The bot has the following external dependencies:

- Google Cloud Spanner. The bot uses spanner for primary storage.

- GitHub. The bot acts as a GitHub webhook to receive notifications of GitHub activity. It
also calls the GitHub API.

- ZenHub. The bot calls into the ZenHub API to get information about GitHub issues.

- SendGrid. The bot sends email using SendGrid.

## Plugins

The bot consists of a simple framework and distinct plugins that have specific isolated responsibilities. There are three
flavors of plugins:
. API plugins expose a REST API, whereas Webhook Plugins listen to GitHub notifications. At the moment,
the plugins are:

Handler plugins expose a REST API to trigger some actions:

- syncer. Synchronizes GitHub data to Google Cloud Spanner, where the data can then be used
for analysis. The syncer needs to be invoked on a periodic basis to refresh the data. This is handled by using
Google Cloud Schedule to invoke the syncer's REST API.

- github. Handles GitHub web hook events by dispatching to the various Webhook plugins deacribed below

Webhook plugins respond to incoming GitHub events:

- cfgmonitor. Monitors GitHub for changes to the bot's configuration file. When it sees such a change, it triggers a
partial shutdown and restart of the bot, which will reread the config and start back up fully.

- labeler. Attached labels to issues and pull requests if specific conditions are detected. This is primarily used
to perform initial triage on incoming issues by assigning an area-specific label to issues based on patterns
found in newly-opened issues.

- nagger. Injects nagging comments in pull requests if specific conditions are detected. This is primarily used to
remind developers to include tests whenever they fix bugs, but the engine is general-purpose and could be used
creatively for other nagging comments.

- refresher. Updates the local Google Cloud Spanner copy of GitHub data based on events
reported by the GitHub webhook.

Finally topic plugins are responsible for the dashboard UI exposed by the bot:

- maintainers. Lets users view the set of project maintainers and perform operations on
them.

## Startup options

The bot supports a number of startup options. These can be specified as environment variables or
via command-line options. Command-line options take precedence over environment variables. The
available startup options are:

- GITHUB_WEBHOOK_SECRET / --github_webhook_secret. Indicates the GitHub secret necessary to authenticate with
the GitHub webhook.

- GITHUB_TOKEN / --github_token. The access token necessary to let the bot invoke the GitHub
API.

- GITHUB_OAUTH_CLIENT_SECRET / --github_oauth_client_secret. The client secret to use in the GitHub OAuth flow,
as obtained in the GitHub admin UI for the target organization.

- GITHUB_OAUTH_CLIENT_ID / --github_oauth_client_id. The client ID to use in the GitHub OAuth flow,
as obtained in the GitHub admin UI for the target organization.

- ZENHUB_TOKEN / --zenhub_token. The access token necessary to let the bot invoke the ZenHub
API.

- GCP_CREDS / --gcp_creds. Base64-encoded JSON credentials for GCP, enabling the bot to invoke
Google Cloud Spanner.

- SENDGRID_APIKEY / --sendgrid_apikey. An API Key for the SendGrid service, enabling the bot to
send emails.

- CONFIG_REPO / --config_repo. The bot can read its configuration directly from a GitHub repository. As
changes are made to the repository, the bot automatically refreshes its configuration. This option lets
you indicate the GitHub organization, repository, and branch where the configuration file can be found.
This is specified as a single string in the form of org/repo/branch.

- CONFIG_FILE / --config_file. Indicates the path to the bot's YAML configuration file. If the config
repo is specified as a startup option, then this file path is relative to the repo. Otherwise, it is
treated as a local file path within the bot's container.

- PORT / --port. The TCP port to listen to for incoming traffic.

## REST API

The bot exposes a REST API at https://eng.istio.io:

- /sync - triggers the bot to synchronize GitHub issues into Google Cloud Spanner. This is called periodically  by 
a job scheduled in Google Cloud scheduler. You can filter what gets synced using a filter query string with a 
command-separated list of things to sync [members, maintainers, issues, prs, labels, zenhub]

- /githubwebhook - used to report events in GitHub. This is called by GitHub whenever anything interesting happens in
the Istio repos.

- /maintainersapi - used to query information about project maintainers

## Configuration file

The bot's behavior is controlled entirely through its configuration file. The
format of the configuration file is described by the `Args` struct in
`pkg/config/args.go` 

## Deployment

Here's how the bot is currently deployed:

- The bot executes using Google Cloud Run. You can do a full build and deploy
a new revision of the bot to Google Cloud Run using `make deploy`. In order for this
to work, you'll need to have setup `gcloud` previously to authenticate
with GCP. Once a new revision is deployed to Cloud Run, it immediately starts receiving traffic.

- The bot depends on having a configured Google Cloud Spanner database. The schema of the database
is described by `spanner.ddl`

- The various credentials used by the bot are set via environment variables specified within the Google Cloud Run
UI.

- TBD: The bot'ss configuration is maintained in the file `policybot/policybot.yaml` in the <https://github.com/istio/test-infra> repo.
Changes pushed to this file are automatically picked up by the bot.

## Credentials and secrets

The bot needs a bunch of credentials to operate. As explained above, these credentials are supplied
to the bot via environment variables or command-line flags. If you're a Googler, you can get access
to the credentials at `go/valentine`.
