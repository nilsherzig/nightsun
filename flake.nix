{
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; }; in rec {
        defaultPackage = pkgs.buildGoModule {
          pname = "nightsun";
          version = "1";
          src = ./.;
          vendorHash = "sha256-XOchm6hHRvyL9t9GMFjVlGlsFoAo8qPnH5VZnhy28jM=";
        };

        devShell = pkgs.mkShell {
          inputsFrom = [ defaultPackage ];
          packages = with pkgs; [
            brotab
            gopls
            goreleaser
            wmctrl
          ];
        };
      });
}
