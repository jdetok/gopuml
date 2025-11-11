# [gopuml](https://github.com/jdetok/gopuml) plantuml generator *IN DEVELOPMENT*
#### 
- CLI tool written in Go to generate UML diagrams (via plantuml) for a codebase 
- ***[plantuml](https://plantuml.com/)*** is an open-source tool that provides a domain-specific-modeling-language for creating [UML diagrams](https://en.wikipedia.org/wiki/Unified_Modeling_Language)
    - [plantuml syntax reference](https://plantuml.com/guide)

# supported programming languages
- ### gopuml currently only supports Go projects
- #### I plan to eventually add suppport for the following langauges:
    - Python
    - Javascript

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
