{
  description = "Nix shell";
  nixConfig.bash-prompt = "[nix] ";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

  };

  outputs =
    { nixpkgs, ... }:
    {
      devShell.aarch64-darwin =
        let
          pkgs = nixpkgs.legacyPackages.aarch64-darwin;
        in
        pkgs.mkShell {
          buildInputs = with pkgs; [
            delve
            go
            gopls
            graphviz
          ];

          shellHook = ''
            for p in $NIX_PROFILES; do
              GOPATH="$p/share/go:$GOPATH"
            done
          '';
        };
    };
}
