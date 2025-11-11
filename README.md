# [gopuml](https://github.com/jdetok/gopuml) | plantuml generator
- ### Go CLI app to scan a code base and create [plantuml](https://plantuml.com/) source code 
- ### written by [Justin DeKock](https://github.com/jdetok)
 
# configuration (JSON)
- gopuml uses a JSON configuration file to setup project name, root, output directory, directories to exclude
- template can be downloaded [here](https://github.com/jdetok/gopuml/blob/main/templates/.gopuml.json)
#### example `.gopuml.json`:
```json
{
    "project_name": "gopuml",
    "project_root": ".",
    "exclude_dirs": [
        ".git",
        "z_docs",
        "templates"
    ],
    "puml_out": "puml"
}
```
#### **NOTE**: `.git` should **ALWAYS** be included in `exclude_dirs` if directory is a git repo
