use protobuf_codegen::Codegen;
use regex::Regex;
use std::fs;
use std::io;
use std::io::Write;
use std::path::Path;

fn generate_v1_protos() {
    Codegen::new()
        .out_dir("src/v1")
        .include("../protos")
        .inputs([
            "../protos/in_toto_attestation/v1/resource_descriptor.proto",
            "../protos/in_toto_attestation/v1/statement.proto",
        ])
        .capture_stderr()
        .run()
        .expect("Protobuf codegen failed for v1 protos");
}

// inspired by https://github.com/stepancheg/rust-protobuf/blob/7131fb244fb1246d2835f5ad7426e607ee7c4a1f/protobuf-codegen/src/gen/mod_rs.rs
fn gen_interm_mod_rs(path: &Path, mods: Vec<String>) -> io::Result<()> {
    // skip if we have no mods
    if mods.is_empty() {
        return Ok(());
    }

    let mod_path = path.join("mod.rs");
    let mut f = fs::File::create(mod_path)?;
    f.write_all(b"// @generated\n")?;

    let mut sorted: Vec<String> = mods.into_iter().collect();
    sorted.sort();
    for m in sorted {
        f.write_fmt(format_args!("pub mod {};\n", m))?;
    }

    Ok(())
}

fn replace_resource_desc_imports(path: &Path) -> io::Result<()> {
    let re = Regex::new(r"(?<super>super::)(?<rd>resource_descriptor)").unwrap();

    let gen_code = fs::read_to_string(path)?;

    // only replace if the module uses resource descriptors
    if re.is_match(&gen_code) {
        let mut f = fs::File::create(path)?;

        for line in gen_code.lines() {
            let replaced = re.replace(line, r"${rd}").into_owned();
            f.write_fmt(format_args!("{}\n", replaced))?;

            if line.starts_with("const _PROTOBUF_VERSION_CHECK") {
                // want to add a single correct import for the resource descriptor module
                f.write_all(b"\nuse crate::v1::resource_descriptor;\n")?;
            }
        }
    }

    Ok(())
}

// this function recurses through the predicates directory
fn generate_predicate_protos(dir: &Path) -> io::Result<()> {
    let prefix = Path::new("../protos/in_toto_attestation/");
    let input_path = prefix.join(dir);
    if input_path.is_dir() {
        // this is to auto-generate mod.rs files at each layer
        let mut mods = Vec::<String>::new();

        // need to convert vX.Y directories into Rust package syntax
        let re = Regex::new(r"(?<pred>\D+)/v(?<major>\d+)\.(?<minor>\d+)").unwrap();
        let proto_dir_str = dir.to_str().unwrap();
        // replaced will either be a borrowed reference to proto_dir_str if it already
        // matched, or a new string if the regex was replaced
        let replaced = re
            .replace(proto_dir_str, r"${pred}/v${major}_${minor}")
            .into_owned();
        let rust_path = Path::new("src/").join(replaced);
        for entry in fs::read_dir(&input_path)? {
            let entry = entry?;
            let path = entry.path();
            if path.is_dir() {
                // save the sub mod name so we can include it in the mod.rs
                let submod_dir_str = path.to_str().unwrap();
                let r = re
                    .replace(submod_dir_str, r"${pred}/v${major}_${minor}")
                    .into_owned();
                let mod_name = Path::new(&r).file_name().unwrap();

                if let Some(m) = mod_name.to_str() {
                    mods.push(m.to_string());
                }

                let child_dir = path.strip_prefix(prefix).unwrap();
                generate_predicate_protos(child_dir)?;
            } else {
                fs::create_dir_all(&rust_path)?;
                Codegen::new()
                    .out_dir(&rust_path)
                    .include("../protos")
                    .input(&path)
                    .capture_stderr()
                    .run()
                    .expect("Protobuf codegen failed for {path}");

                let gen_mod = path.file_stem().unwrap();
                let gen_file = gen_mod.to_str().unwrap().to_owned() + ".rs";
                replace_resource_desc_imports(&Path::new(&rust_path).join(gen_file))?;
            }
        }
        gen_interm_mod_rs(&rust_path, mods)?;
    }
    Ok(())
}

fn main() {
    generate_v1_protos();
    generate_predicate_protos(Path::new("predicates")).unwrap();
}
