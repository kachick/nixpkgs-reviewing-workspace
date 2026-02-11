{
  buildGoModule,
  makeWrapper,
  lib,
  gh,
  coreutils,
}:
buildGoModule {
  pname = "resume";
  version = "0.1.0";
  src = ./.;

  vendorHash = "sha256-uPqabZgQGQulf+F3BvMLhv4O0h5jOq12F7K60u5xjtA=";

  nativeBuildInputs = [
    makeWrapper
  ];

  postInstall = ''
    wrapProgram "$out/bin/resume" \
      --prefix PATH : ${
        lib.makeBinPath [
          gh
          coreutils
        ]
      }
  '';

  meta = {
    description = "Resume to track a running job in console";
    mainProgram = "resume";
  };
}
