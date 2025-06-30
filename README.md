# RepoFork

RepoFork makes it simple to keep a mirror repository in sync with an upstream project while stripping out unwanted files or directories in your fork (e.g. CI configs, internal tooling, etc.).
It supports both initializing a new fork and incrementally updating an existing fork with only the new commits since your last sync.

## How it works

- RepoFork clones both the mirror (origin) and upstream repositories into a temporary working directory.
- For first-time forks, it performs a full history rewrite using git-filter-repo.
- For incremental updates, it:
  - Finds the last upstream commit already mirrored
  - Cherry-picks all newer upstream commits one by one (supporting merge commits)
  - Performs a partial history rewrite of the new commits to remove unwanted files or directories
- Rewrites commit messages to embed the original upstream commit ID for traceability.
- Optionally pushes the updated branch back to the mirror repository.

## Getting Started

> RepoFork uses your local git configuration for authentication.

### Requirements

- git
- git-filter-repo

### Installation

You can install `repofork` using the following command:

TODO: add installation instructions

## Usage

**Minimal**

```bash
repofork update --origin <mirror-repo> --upstream <upstream-repo> --push=true
```

**Full**

```bash
repofork update \
  --origin <mirror-repo> \
  --origin-branch main \
  --upstream <upstream-repo> \
  --upstream-branch master \
  --push=true
```

**Recreate**

```bash
repofork update \
  --origin <mirror-repo> \
  --origin-branch main \
  --upstream <upstream-repo> \
  --upstream-branch master \
  --full-rewrite=true \
  --push=true
```

## License

Released under the [MIT license](./LICENSE).
