{
  inputs = {
    # Prefer nixpkgs- rather than nixos- for darwin
    nixpkgs.url = "https://channels.nixos.org/nixpkgs-unstable/nixexprs.tar.zst";
  };

  outputs =
    {
      nixpkgs,
      self,
      ...
    }:
    let
      inherit (nixpkgs) lib;
      forAllSystems = lib.genAttrs lib.systems.flakeExposed;
    in
    rec {
      formatter = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        pkgs.writeShellScriptBin "dprint-fmt" ''
          exec "${lib.getExe pkgs.dprint}" fmt "$@"
        ''
      );

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

            inputsFrom = [
              self.packages.${system}.resume
            ];

            buildInputs = (
              with pkgs;
              [
                bashInteractive
                coreutils # mktemp
                nixd
                go-task

                dprint
                typos

                nixpkgs-reviewFull
                bubblewrap # Require to run nixpkgs-review with sandbox mode. See https://github.com/Mic92/nixpkgs-review/pull/441
                gh
                git
                tree
                fd
                fzf

                go
                gopls

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
