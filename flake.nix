{
  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs { inherit system; }; in rec {
        defaultPackage = pkgs.buildGoModule {
          pname = "unified-search";
          version = "0";
          src = ./.;
          vendorHash = "sha256-U/AuOIC1eYVi+lwjD5jS8lDInqnwOwdlq0X+9bC3R7c=";
        };

        devShell = pkgs.mkShell {
          inputsFrom = [ defaultPackage ];
          packages = with pkgs; [
            brotab
            gopls
            goreleaser
          ];
        };
      });
}
