{
  description = "PAM bindings for writing service modules in Golang";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    git-hooks-nix = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs:
    inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.git-hooks-nix.flakeModule
      ];
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];
      perSystem =
        {
          pkgs,
          system,
          config,
          ...
        }:
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [
              (final: prev: {
                pamtester = prev.pamtester.overrideAttrs (_: _: { meta.platform = prev.lib.platforms.unix; });
              })
            ];
            config = { };
          };
          pre-commit = {
            check.enable = true;
            settings = {
              src = ./.;
              hooks = {
                actionlint.enable = true;
                nixfmt-rfc-style.enable = true;
                gofmt.enable = true;
                govet.enable = true;
                govet.extraPackages = [ pkgs.pam ];
              };
            };
          };
          devShells = {
            default = pkgs.mkShell {
              inherit (config.pre-commit.devShell) shellHook;
              packages =
                with pkgs;
                [
                  go
                  pam
                  pamtester
                ]
                ++ config.pre-commit.settings.enabledPackages;
            };
          };
        };
    };
}
