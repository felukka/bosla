{
  buildGoModule,
}:

let
  name = "bosla";
  version = "0.0.1";
in
buildGoModule {
  pname = name;
  inherit version;

  src = ./.;
  subPackages = [
    "cmd/bosla"
  ];

  vendorHash = "sha256-OlU2sm5ol7L5BqHo8Qc2v1ayTjUfbh/UM3orew7dnd4=";

  meta = {
    mainProgram = "bosla";
    homepage = "https://bosla.felukka.org";
    description = "A serious, well-scoped versioning tool.";
  };
}
