{
  description = "nichts - configuration for machines!";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    pre-commit-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      pre-commit-hooks,
      ...
    }:
    let
      forAllSystems = nixpkgs.lib.genAttrs [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];
    in
    {
      checks = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          pre-commit-check = pre-commit-hooks.lib.${system}.run {
            src = ./.;
            hooks = {
              actionlint.enable = true;
              nixfmt-rfc-style.enable = true;
              gofmt.enable = true;
              govet.enable = true;
              govet.extraPackages = [ pkgs.pam ];
              staticcheck.enable = true;
              staticcheck.extraPackages = [ pkgs.go ];
            };
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            inherit (self.checks.${system}.pre-commit-check) shellHook;
            packages =
              with pkgs;
              [
                go
                pam
                pamtester
              ]
              ++ self.checks.${system}.pre-commit-check.enabledPackages;
          };
        }
      );
    };
}
