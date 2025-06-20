# RepoFork

> Helper to maintain forks of upstream repositories with a rewrite of the history to remove unwanted files and directories.

**Required Tools:**

- git
- git-filter-repo

## Features

**Empty Repository**:

- Clones the upstream repository
- Rewrites the history, removing filtered files and directories and adding the upstream commit id into the commit message
- Pushes the rewritten history to a new repository

**Existing Repository**:

- Clones the mirror repository
- Searches for the latest upstream commit id in the commit messages
- Individually cherry-picks commits from the upstream repository that are not present in the mirror repository
- Rewrites the history of cherry-picked commits, removing filtered files and directories and adding the upstream commit id into the commit message
- Pushes the rewritten history to the mirror repository

## Usage

```bash
repofork update --origin <mirror-repo> --upstream <upstream-repo>
repofork update --origin <mirror-repo> --origin-branch main --upstream <upstream-repo> --upstream-branch master --full-rewrite=true --push=true
```

## License

Released under the [MIT license](./LICENSE).
