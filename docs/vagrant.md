# Vagrant

We use Vagrant to demonstrate how this tool works.

## Deps

Install required dependencies:

* **[Virtual box](https://www.virtualbox.org/wiki/Downloads)** (tested on version 4.3.26)
* **[Vagrant](https://www.vagrantup.com/downloads.html)** (tested on version 1.8.4)
* **[vagrant-docker-compose](https://github.com/leighmcculloch/vagrant-docker-compose)** plugin with command
  ```
  vagrant plugin install vagrant-docker-compose
  ```

## Getting Started

There are two ways you can invoke the demo. 

### Using Config File

Copy `.github-authorized-keys-demo.default.yml` to `.github-authorized-keys-demo.yml`.

```
cp .github-authorized-keys-demo.default.yml .github-authorized-keys-demo.yml
```

Then set the required values in the `.github-authorized-keys-demo.yml` file.

After that, simply run:

```
vagrant up
```

### Using Environment Variables

Simply run:

```
GITHUB_API_TOKEN={api token} \
GITHUB_ORGANIZATION={organization name} \
GITHUB_TEAM={team name} \
  vagrant up
```

## Testing

Login into vagrant box via `ssh` and your GitHub username:

```
ssh -o "UserKnownHostsFile /dev/null" {github username}@192.168.33.10
```

Then review the logs inside of the Vagrant box to see what's going on:

```
sudo tail -f /var/log/auth.log
```
