{
  inputs = {
    # Prefer nixpkgs- rather than nixos- for darwin
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    selfup = {
      url = "github:kachick/selfup/v1.1.8";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      selfup,
    }:
    let
      inherit (nixpkgs) lib;
      forAllSystems = lib.genAttrs lib.systems.flakeExposed;
    in
    {
      formatter = forAllSystems (system: nixpkgs.legacyPackages.${system}.nixfmt-rfc-style);
      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShellNoCC {
            buildInputs =
              (with pkgs; [
                bashInteractive
                findutils # xargs
                coreutils # mktemp
                nixfmt-rfc-style
                nil
                go-task

                dprint
                typos

                nixpkgs-review
                gh
                git
                tree
                fd
              ])
              ++ [
                selfup.packages.${system}.default
              ];
          };
        }
      );
    };
}
