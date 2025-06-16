{
  description = "Go application with Swagger";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {inherit system;};

        # Custom swag package for the specific version
        swag = pkgs.buildGoModule rec {
          pname = "swag";
          version = "v2.0.0-rc4";

          src = pkgs.fetchFromGitHub {
            owner = "swaggo";
            repo = "swag";
            rev = version;
            sha256 = "sha256-3dX01PAVEkn8ciVNbIn1IwiSkwogPJLYNo6xTg9jhDA=";
          };

          vendorHash = "sha256-bUMW9wjPIT3JLMw9F/NvWqZv1M62o/Y4gIpp6XyHbek=";
          subPackages = ["cmd/swag"];
        };

        app = pkgs.buildGoModule {
          pname = "go-app";
          version = "0.1.0";

          src = ./.;

          vendorHash = "sha256-jx70HhXLbBt63Vt0iZK8aUwoQnpe57MIjUMzXsnaRgA=";

          nativeBuildInputs = [swag];

          preBuild = ''
            # Generate Swagger docs
            swag init -v3.1 -o docs -g main.go --parseDependency --parseInternal
          '';

          meta = with pkgs.lib; {
            description = "Go application with Swagger documentation";
            license = licenses.mit;
          };
        };
      in {
        packages = {
          default = app;
          app = app;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_24
            swag
            air
          ];

          shellHook = ''
            echo "Go development environment ready"
            echo "Go version: $(go version)"
            echo "Swag version: $(swag --version)"
          '';
        };
      }
    );
}
