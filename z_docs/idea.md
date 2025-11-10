# gopuml - plantuml generator for Go projects

## idea outline
- command line tool
- json config file in root dir of project it's generating for
    - set type of diagram, where to output puml source, png, etc
- need to go into a directory and read each go file
    - parse through the files, find which should be packages, classes, etc
    