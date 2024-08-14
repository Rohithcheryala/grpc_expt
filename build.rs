use std::env;
use std::path::PathBuf;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let proto_file = "./src/proto_src/store.proto";
    let proto_dir = "./src/proto_src";
    let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());

    tonic_build::configure()
        .protoc_arg("--experimental_allow_proto3_optional") // for older systems
        .build_client(true)
        .build_server(true)
        .file_descriptor_set_path(out_dir.join("store_descriptor.bin"))
        .out_dir("./src/proto_out/")
        .compile(&[proto_file], &[proto_dir])?;

    tonic_build::configure()
        .protoc_arg("--experimental_allow_proto3_optional") // for older systems
        .build_client(true)
        .build_server(false)
        .file_descriptor_set_path(out_dir.join("store_descriptor.bin"))
        .out_dir("./src/proto_out/client")
        .compile(&["./src/proto_src/store.proto"], &["./src/proto_src"])?;

    tonic_build::configure()
        .protoc_arg("--experimental_allow_proto3_optional") // for older systems
        .build_client(false)
        .build_server(true)
        .file_descriptor_set_path(out_dir.join("store_descriptor.bin"))
        .out_dir("./src/proto_out/server")
        .compile(&[proto_file], &[proto_dir])?;

    Ok(())
}
