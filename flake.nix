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
        let
          buildGoPAMModule =
            { pname, ... }@args:
            pkgs.buildGoModule (
              args
              // {
                env.CGO_CFLAGS = "-I${pkgs.pam}/include";

                buildPhase = ''
                  runHook preBuild

                  if [ -z "$enableParallelBuilding" ]; then
                      export NIX_BUILD_CORES=1
                  fi

                  go build '-ldflags=-buildid= -extldflags="-L${pkgs.pam}/lib"' -buildmode=c-shared -o ${pname}.so -p "$NIX_BUILD_CORES" .
                  go build '-ldflags=-buildid= -extldflags="-L${pkgs.pam}/lib"' -o ${pname}-helper -p "$NIX_BUILD_CORES" .
                  chmod +x ${pname}.so
                  runHook postBuild
                '';

                installPhase = ''
                  runHook preInstall

                  mkdir -p $out/lib/security
                  mkdir -p $out/bin
                  cp ${pname}.so $out/lib/security
                  cp ${pname}-helper $out/bin

                  runHook postInstall
                '';
              }
            );
        in
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

          packages = {
            pam_spicedb = buildGoPAMModule {
              pname = "pam_spicedb";
              version = "0.0.1";
              src = ./examples/pam_spicedb;
              vendorHash = null;

              meta = {
                description = "pam_spicedb.so is a PAM module that validates a user's account by checking a permission in SpiceDB.";
                mainProgram = "lib/security/pam_spicedb.so";
                homepage = "https://github.com/squat/pam";
              };
            };

            pam_print = buildGoPAMModule {
              pname = "pam_print";
              version = "0.0.1";
              src = ./examples/pam_print;
              vendorHash = null;

              meta = {
                description = "pam_print.so is a PAM module that prints the user, flags, arguments, and environment variables provided to the module.";
                mainProgram = "lib/security/pam_print.so";
                homepage = "https://github.com/squat/pam";
              };
            };
          };

          pre-commit = {
            check.enable = true;
            settings = {
              src = ./.;
              hooks = {
                actionlint.enable = true;
                nixfmt-rfc-style.enable = true;
                gofmt.enable = true;
                gofmt.excludes = [ "examples/pam_.*/vendor" ];
                govet.enable = true;
                govet.excludes = [ "examples/pam_.*/vendor" ];
                govet.extraPackages = [ pkgs.pam ];
                readme = {
                  enable = true;
                  name = "README.md";
                  entry =
                    let
                      readmeCheck = pkgs.writeShellApplication {
                        name = "readme-check";
                        text = ''
                          for f in "$@"; do
                              if ! grep -q embedmd "$f"; then continue; fi
                              (cd "$(dirname "$f")" && go run ./... --help 2>help.txt)
                              go tool embedmd -d "$f"
                          done
                        '';
                      };
                    in
                    pkgs.lib.getExe readmeCheck;
                  files = "README\\.md$";
                  extraPackages = [ pkgs.go ];
                  excludes = [ "examples/pam_.*/vendor" ];
                };
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
