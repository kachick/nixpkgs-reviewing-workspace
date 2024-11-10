# How to develop

## Setup

1. Install [Nix](https://nixos.org/) package manager and enable [Flakes](https://nixos.wiki/wiki/Flakes)\
   Or use Nix installed containers. For example, look at [this repo](https://github.com/kachick/containers)
2. Run dev shell as one of the following
   - with [direnv](https://github.com/direnv/direnv): `direnv allow`
   - nix only: `nix develop`
3. You can use development tools

```console
> nix develop
(prepared bash)
> dprint --version
...
```

## Note

- nixpkgs-review supports shallow clone since <https://github.com/Mic92/nixpkgs-review/issues/267> => <https://github.com/Mic92/nixpkgs-review/commit/055465e55d131ffb1e1617f46d3bade0b87bbe69>
