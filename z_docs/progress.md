## todo on flight back home 11/12
- get as many regex patterns ready as possible
- setup architecture to map matches
    - brainstorm map structure: 
        - dir (package) outer
            - file 
                - line number ?? 
                    - struct holding different fields depending on map, func, etc

## end of 11/12
- regex patterns more robust
- implemented ordered checking to patterns - more exclusive patterns tested first

## 11/14 starting point
- reading config file is done
- reading all files in dir passed as -root is done, mapped into a workable structure
- regex structure for parsing files is setup and usable
- matches mapping into types is setup but may need more thought
- need to get together a mvp outputting a file
**GOAL:**
- create structure for outputting the .puml file