# Deploying to Dokku

> Note: This walkthrough uses the hostname `dokku.me` in commands. When deploying to your own server, you should set the `DOKKU_DOMAIN` value in the `Vagrantfile` before you initialize the Vagrant VM.

## Deploy tutorial

Once you have configured Dokku with at least one user, you can deploy applications using `git push`. To quickly see Dokku deployment in action, try using [the Heroku Ruby on Rails "Getting Started" app](https://github.com/heroku/ruby-getting-started).

```shell
# from your local machine
# SSH access to github must be enabled on this host
git clone https://github.com/heroku/ruby-getting-started
```

### Create the app

SSH into the Dokku host and create the application as follows:

```shell
# on the Dokku host
dokku apps:create ruby-getting-started
```

### Create the backing services

Dokku by default **does not** provide datastores (e.g. MySQL, PostgreSQL) on a newly created app. You can add datastore support by installing plugins, and the Dokku project [provides official plugins](/docs/community/plugins.md#official-plugins-beta) for common datastores.

The Getting Started app requires a PostgreSQL service, so install the plugin and create the related service as follows:

```shell
# on the Dokku host
# install the postgres plugin
# plugin installation requires root, hence the user change
sudo dokku plugin:install https://github.com/dokku/dokku-postgres.git

# create a postgres service with the name railsdatabase
dokku postgres:create railsdatabase
```

Each service may take a few moments to create.

### Linking backing services to applications

Once the services have been created, you then set the `DATABASE_URL` environment variable by linking the service, as follows:

```shell
# on the Dokku host
# each official datastore offers a `link` method to link a service to any application
dokku postgres:link railsdatabase ruby-getting-started
```

Dokku supports linking a single service to multiple applications as well as linking only one service per application.

### Deploy the app

> Warning: Your app should respect the `PORT` environment variable, otherwise it may not respond to web requests. You can find more information in the [port management documentation](/docs/networking/port-management.md).**

Now you can deploy the `ruby-getting-started` app to your Dokku server. All you have to do is add a remote to name the app. Applications are created on-the-fly on the Dokku server.

```shell
# from your local machine
# the remote username *must* be dokku or pushes will fail
cd ruby-getting-started
git remote add dokku dokku@dokku.me:ruby-getting-started
git push dokku master
```

> Note: Some tools may not support the short-upstream syntax referenced above, and you may need to prefix
> the upstream with the scheme `ssh://` like so: `ssh://dokku@dokku.me:ruby-getting-started`
> Please see the [Git](https://git-scm.com/docs/git-clone#_git_urls_a_id_urls_a) documentation for more details.

> Note: Your private key should be registered with `ssh-agent` in your local development environment. If you get a `permission denied` error when pushing, you can register your private key as follows: `ssh-add -k ~/<your private key>`.

After running `git push dokku master`, you should have output similar to this in your terminal:

```
Counting objects: 231, done.
Delta compression using up to 8 threads.
Compressing objects: 100% (162/162), done.
Writing objects: 100% (231/231), 36.96 KiB | 0 bytes/s, done.
Total 231 (delta 93), reused 147 (delta 53)
-----> Cleaning up...
-----> Building ruby-getting-started from herokuish...
-----> Adding BUILD_ENV to build environment...
-----> Ruby app detected
-----> Compiling Ruby/Rails
-----> Using Ruby version: ruby-2.2.1
-----> Installing dependencies using 1.9.7
       Running: bundle install --without development:test --path vendor/bundle --binstubs vendor/bundle/bin -j4 --deployment
       Fetching gem metadata from https://rubygems.org/...........
       Fetching version metadata from https://rubygems.org/...
       Fetching dependency metadata from https://rubygems.org/..
       Using rake 10.4.2

...

=====> Application deployed:
       http://ruby-getting-started.dokku.me
```

Once the deploy is complete, the application's web URL will be generated as above.

Dokku supports deploying applications via [Heroku buildpacks](https://devcenter.heroku.com/articles/buildpacks) with [Herokuish](https://github.com/gliderlabs/herokuish#buildpacks), as well as by using a project's [Dockerfile](https://docs.docker.com/reference/builder/).


### Skipping deployment

If you only want to rebuild and tag a container, you can skip the deployment phase by setting `$DOKKU_SKIP_DEPLOY` to `true` by running:

``` shell
# on the Dokku host
dokku config:set ruby-getting-started DOKKU_SKIP_DEPLOY=true
```

### Redeploying or restarting

If you need to redeploy or restart your app: 

```shell
# on the Dokku host
dokku ps:rebuild ruby-getting-started
```

See the [process scaling documentation](/docs/deployment/process-management.md) for more information.

### Deploying with private Git submodules

Dokku uses Git locally (i.e. not a Docker image) to build its own copy of your app repo, including submodules, as the `dokku` user. This means that in order to deploy private Git submodules, you need to put your deploy key in `/home/dokku/.ssh/` and potentially add `github.com` (or your VCS host key) into `/home/dokku/.ssh/known_hosts`. You can use the following test to confirm your setup is correct:

```shell
# on the Dokku host
su - dokku
ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
ssh -T git@github.com
```

> Warning: if the buildpack or Dockerfile build process require SSH key access for other reasons, the above may not always apply.

## Deploying to subdomains

If you do not enter a fully qualified domain name when pushing your app, Dokku deploys the app to `<remotename>.yourdomain.tld` as follows:

```shell
# from your local machine
# the remote username *must* be dokku or pushes will fail
git remote add dokku dokku@dokku.me:ruby-getting-started
git push dokku master
```

```
remote: -----> Application deployed:
remote:        http://ruby-getting-started.dokku.me
```

You can also specify the fully qualified name as follows:

```shell
# from your local machine
# the remote username *must* be dokku or pushes will fail
git remote add dokku dokku@dokku.me:app.dokku.me
git push dokku master
```

```
remote: -----> Application deployed:
remote:        http://app.dokku.me
```

This is useful when you want to deploy to the root domain:

```shell
# from your local machine
# the remote username *must* be dokku or pushes will fail
git remote add dokku dokku@dokku.me:dokku.me
git push dokku master
```

    ... deployment ...

    remote: -----> Application deployed:
    remote:        http://dokku.me

## Dokku/Docker container management compatibility

Dokku is, at its core, a Docker container manager. Thus, it does not necessarily play well with other out-of-band processes interacting with the Docker daemon.

Prior to every deployment, Dokku will execute a cleanup function. As of 0.5.x, the cleanup removes all containers with the `dokku` label where the status is either `dead` or `exited` (previous versions would remove _all_ `dead` or `exited` containers). The cleanup function also removes all images with `dangling` status.

## Adding deploy users

See the [user management documentation](/docs/deployment/user-management.md).

## Default vhost

See the [nginx documentation](/docs/configuration/nginx.md#default-site).

## Deploying non-master branch

See the [Git documentation](/docs/deployment/methods/git.md#changing-the-deploy-branch).

## Dockerfile deployment

See the [Dockerfile documentation](/docs/deployment/methods/dockerfiles.md).

## Image tagging

See the [image tagging documentation](/docs/deployment/methods/images.md).

## Specifying a custom buildpack

See the [buildpack documentation](/docs/deployment/methods/buildpacks.md).

## Removing a deployed app

See the [application management documentation](/docs/deployment/application-management.md#removing-a-deployed-app).

## Renaming a deployed app

See the [application management documentation](/docs/deployment/application-management.md#renaming-a-deployed-app).

## Zero downtime deploy

See the [zero-downtime deploy documentation](/docs/deployment/zero-downtime-deploys.md).
