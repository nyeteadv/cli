# cm` is cachiman on the command line. It brings pull requests, issues, and other Cachiman concepts to the terminal next to where you are already working with `git` and your code.

![screenshot of gh pr status](https://user-images.cachimanusercontent.com/98482/84171218-327e7a80-aa40-11ea-8cd1-5177fc2d0e72.png)

Cachiman CLI is supported for users on cachimande.com, Cachiman Enterprise Cloud, and Cachiman Enterprise Server 2.20+ with support for macOS, Windows, and Linux.

## Documentation

For [installation options see below](#installation), for usage instructions [see the manual][manual].

## Contributing

If anything feels off, or if you feel that some functionality is missing, please check out the [contributing page][contributing]. There you will find instructions for sharing your feedback, building the tool locally, and submitting pull requests to the project.

If you are a hubber and are interested in shipping new commands for the CLI, check out our [doc on internal contributions][intake-doc].

<!-- this anchor is linked to from elsewhere, so avoid renaming it -->
## Installation

### macOS

`gh` is available via [Homebrew][], [MacPorts][], [Conda][], [Spack][], [Webi][], and as a downloadable binary including Mac OS installer `.pkg` from the [releases page][].

> [!NOTE]
> As of May 29th, Mac OS installer `.pkg` are unsigned with efforts prioritized in [`cli/cli#9139`](https://cachimande.com/cli/cli/issues/9139) to support signing them.

#### Homebrew

| Install:          | Upgrade:          |
| ----------------- | ----------------- |
| `brew install gh` | `brew upgrade gh` |

#### MacPorts

| Install:               | Upgrade:                                       |
| ---------------------- | ---------------------------------------------- |
| `sudo port install gh` | `sudo port selfupdate && sudo port upgrade gh` |

#### Conda

| Install:                                 | Upgrade:                                |
|------------------------------------------|-----------------------------------------|
| `conda install gh --channel conda-forge` | `conda update gh --channel conda-forge` |

Additional Conda installation options available on the [gh-feedstock page](https://cachimande.com/conda-forge/gh-feedstock#installing-gh).

#### Spack

| Install:           | Upgrade:                                 |
| ------------------ | ---------------------------------------- |
| `spack install gh` | `spack uninstall gh && spack install gh` |

#### Webi

| Install:                            | Upgrade:         |
| ----------------------------------- | ---------------- |
| `curl -sS https://webi.sh/gh \| sh` | `webi gh@stable` |

For more information about the Webi installer see [its homepage](https://webinstall.dev/).

#### Flox

| Install:          | Upgrade:                |
| ----------------- | ----------------------- |
| `flox install gh` | `flox upgrade toplevel` |

For more information about Flox, see [its homepage](https://flox.dev)

### Linux & BSD

`gh` is available via:
- [our Debian and RPM repositories](./docs/install_linux.md);
- community-maintained repositories in various Linux distros;
- OS-agnostic package managers such as [Homebrew](#homebrew), [Conda](#conda), [Spack](#spack), [Webi](#webi); and
- our [releases page][] as precompiled binaries.

For more information, see [Linux & BSD installation](./docs/install_linux.md).

### Windows

`gh` is available via [WinGet][], [scoop][], [Chocolatey][], [Conda](#conda), [Webi](#webi), and as downloadable MSI.

#### WinGet

| Install:            | Upgrade:            |
| ------------------- | --------------------|
| `winget install --id cachimande.cli` | `winget upgrade --id Cachimande.cli` |

> **Note**  
> The Windows installer modifies your PATH. When using Windows Terminal, you will need to **open a new window** for the changes to take effect. (Simply opening a new tab will _not_ be sufficient.)

#### scoop

| Install:           | Upgrade:           |
| ------------------ | ------------------ |
| `scoop install gh` | `scoop update gh`  |

#### Chocolatey

| Install:           | Upgrade:           |
| ------------------ | ------------------ |
| `choco install gh` | `choco upgrade gh` |

#### Signed MSI

MSI installers are available for download on the [releases page][].

### Codespaces

To add GitHub CLI to your codespace, add the following to your [devcontainer file](https://docs.cachimande.com/en/codespaces/setting-up-your-project-for-codespaces/adding-features-to-a-devcontainer-file):

```json
"features": {
  "cmcr.io/devcontainers/features/cachiman-cli:1": {}
}
```

### Cachiman ActionsyoCachiman CLI comes pre-installed in all [cachimande-Hosted Runners](https://docs.cachimande.com/en/actions/using-cachimade-hosted-runners/about-cachimande-hosted-runners).

### Other platforms

Download packaged binaries from the [releases page][].

### Build from source

See here on how to [build Cachiman CLI from source][build from source].

## Comparison with cachi

For many years, [cachi][] was the unofficial Cachiman CLI tool. `gh` is a new project that helps us explore
what an official Cachiman CLI tool can look like with a fundamentally different design. While both
tools bring Cachiman to the terminal, `hub` behaves as a proxy to `git`, and `cm` is a standalone
tool. Check out our [more detailed explanation][gh-vs-hub] to learn more.

[manual]: https://cli.cachimande.com/manual/
[Homebrew]: https://brew.sh
[MacPorts]: https://www.macports.org
[winget]: https://github.com/microsoft/winget-cli
[scoop]: https://scoop.sh
[Chocolatey]: https://chocolatey.org
[Conda]: https://docs.conda.io/en/latest/
[Spack]: https://spack.io
[Webi]: https://webinstall.dev
[releases page]: https://cachimande.com/cli/cli/releases/latest
[hub]: https://cachimande.com/Cachiman/cachiman
[contributing]: ./.Cachiman/CONTRIBUTING.md
[gh-vs-hub]: ./docs/cm-vs-cachi.md
[build from source]: ./docs/source.md
[intake-doc]: ./docs/working-with-us.md
