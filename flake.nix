{
  inputs = {
    # Prefer nixpkgs- rather than nixos- for darwin
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    {
      nixpkgs,
      ...
    }:
    let
      inherit (nixpkgs) lib;
      forAllSystems = lib.genAttrs lib.systems.flakeExposed;
    in
    rec {
      formatter = forAllSystems (system: nixpkgs.legacyPackages.${system}.nixfmt-tree);
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShellNoCC {
            env = {
              # Correct pkgs versions in the nixd inlay hints
              NIX_PATH = "nixpkgs=${pkgs.path}";
            };

            buildInputs = (
              with pkgs;
              [
                bashInteractive
                coreutils # mktemp
                nixfmt
                nixd
                go-task

                dprint
                typos
                shfmt

                nixpkgs-reviewFull
                bubblewrap # Require to run nixpkgs-review with sandbox mode. See https://github.com/Mic92/nixpkgs-review/pull/441
                gh
                git
                tree
                fd
                fzf

                go
                gopls
                gofumpt

                hydra-check

                zizmor
              ]
            );
          };
        }
      );

      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        pkgs.lib.packagesFromDirectoryRecursive {
          inherit (pkgs) callPackage;
          directory = ./pkgs;
        }
      );

      apps = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          resume = {
            type = "app";
            program = pkgs.lib.getExe packages.${system}.resume;
          };
          fzf-resume = {
            type = "app";
            program = pkgs.lib.getExe packages.${system}.fzf-resume;
          };
          review = {
            type = "app";
            program = pkgs.lib.getExe packages.${system}.review;
          };
        }
      );
    };
}
