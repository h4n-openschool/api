{
  description = "The classes microservice of OpenSchool";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      # to work with older version of flakes
      lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";

      # Generate a user-friendly version number.
      version = builtins.substring 0 8 lastModifiedDate;

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          api = pkgs.buildGoModule {
            pname = "api";
            inherit version;
            src = ./.;
            proxyVendor = true;
            vendorSha256 = "sha256-/vxP7hSd4Lq3thzBilAd9Zb/QGHP2nWmxYr2j2e2vJE=";
          };
        });

      defaultPackage = forAllSystems (system: self.packages.${system}.api);
    };
}
