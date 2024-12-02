# :warning: Repository not maintained :warning:

Please note that this repository is currently archived, and is no longer being maintained.

- It may contain code, or reference dependencies, with known vulnerabilities
- It may contain out-dated advice, how-to's or other forms of documentation

The contents might still serve as a source of inspiration, but please review any contents before reusing elsewhere.

# manage-aadgroup-members

This is normally done from a pipeline, but during development it is easier to run it locally.

## Development practices

Create a local file called .env with the following content:

```env
AZURE_SUBSCRIPTION_ID=<REDACTED>
AZURE_TENANT_ID=<REDACTED>
AZURE_CLIENT_ID=<REDACTED>
AZURE_CLIENT_SECRET=<REDACTED>
AZURE_GROUP_OBJECT_ID=<REDACTED>
```

Where <REDACTED> is replaced with actual values. The .env file is already specified in .gitignore.

### Build

```bash
make
```

### Run

```bash
make run
```

### Run with Docker Compose

```bash
make start
```
