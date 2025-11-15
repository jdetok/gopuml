# [gopuml](https://github.com/jdetok/gopuml) plantuml generator *IN DEVELOPMENT*
#### 
- CLI tool written in Go to generate UML diagrams (via plantuml) for a codebase 
- ***[plantuml](https://plantuml.com/)*** is an open-source tool that provides a domain-specific-modeling-language for creating [UML diagrams](https://en.wikipedia.org/wiki/Unified_Modeling_Language)
    - [plantuml syntax reference](https://plantuml.com/guide)

# HOW TO USE/DEMO
- clone and cd into the repo: `git clone https://github.com/jdetok/gopuml.git`
- run gopuml on the repo
    - from demo binary: 
    `./bin/demo`
    - using go run:
    `go run ./main`
- view plantuml output:
    - by default, it will output to `./pumlout/sample_uml_class.puml`
    - the rendered UML output can be viewed using:
        - copy and paste content of the `.puml` output file into the [official online Plantuml editor](https://editor.plantuml.com/)
        - [vs-code plantuml extension](https://marketplace.visualstudio.com/items?itemName=jebbs.plantuml)
        - any other tool that renders plantuml
            

# supported languages
- ### gopuml currently only supports Go projects
- I plan to eventually add suppport for the following langauges:
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
    "languages": [
        ".go"
    ],
    "exclude_dirs": [
        ".git",
        "z_docs",
        "templates"
    ],
    "puml_out_dir": "pumlout",
    "class_diagram_file": "sample_uml_class",
    "class_diagram_title": "Sample Class Diagram",
    "activity_diagram_file": "sample_uml_activity",
    "activity_diagram_title": "Sample Activity Diagram"
}
```
#### **NOTE**: `.git` should **ALWAYS** be included in `exclude_dirs` if directory is a git repo
