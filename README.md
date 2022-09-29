# mataroa-cli

A CLI tool for mataroa.blog

## Usage

Run `mata init` to get started. Read the [man pages](docs/manpages) to learn
about all commands.

## Community

We have a mailing list at
[~sirodoht/mataroa-community@lists.sr.ht](mailto:~sirodoht/mataroa-community@lists.sr.ht)
for the mataroa community to introduce themselves, their blogs, and discuss
anything that’s on their mind!

Archives at
[lists.sr.ht/~sirodoht/mataroa-community](https://lists.sr.ht/~sirodoht/mataroa-community)

## Contributing

Feel free to open a PR on [GitHub](https://github.com/mataroa-blog/mata) or send
an email patch to
[~sirodoht/public-inbox@lists.sr.ht](mailto:~sirodoht/public-inbox@lists.sr.ht).

On how to contribute using email patches see

[git-send-email.io](https://git-send-email.io/).

Also checkout our docs on:

* [Git Commit Message Guidelines](docs/commit-messages.md)
* [File Structure Walkthrough](docs/file-structure-walkthrough.md)
* [Dependencies](docs/dependencies.md)

## Development

This is a [Go](https://go.dev) codebase. Check out the [Go
docs](https://go.dev/doc/) for general technical documentation.

### Structure

The CLI-related code lives under the `cmd/mata` directory and the API client
lives under the `mataroa` directory.

### Dependencies with Nix

This project is configured with [Nix](https://nixos.org). You can simply run
`nix develop` and have all the needed right away.

You can also run the program directly with nix: `nix run github:mataroa-blog/mata -- <commands>`

## Documentation

Also available as man pages:

- [mata(1)](docs/manpages/mata.1.md)
- [mata-config(5)](docs/manpages/mata.5.md)

## License

This software is licensed under the MIT license. For more information, read the
[LICENSE](LICENSE) file.
