{
  rustPlatform,
  makeWrapper,
  lib,
  gh,
  coreutils,
}:
rustPlatform.buildRustPackage {
  pname = "resume-rust";
  version = "0.1.0";
  src = ./.;

  cargoHash = "sha256-3C3XEi9TUKqSR0Tr7ioefUjh/ud0/eTrFOrBpYfeJ+c=";

  nativeBuildInputs = [
    makeWrapper
  ];

  postInstall = ''
    wrapProgram "$out/bin/resume-rust" --prefix PATH : ${
      lib.makeBinPath [
        gh
        coreutils
      ]
    }
  '';

  meta = {
    description = "Resume to track a running job in console (Rust version)";
    mainProgram = "resume-rust";
  };
}
