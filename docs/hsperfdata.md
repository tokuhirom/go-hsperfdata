# how hsperfdata works

java creates `/tmp/hsperfdata_${USER}/${PID}` file.
It's a binary file. The file contains useful data for performance analysis.

The library placed on "github.com/tokuhirom/go-hsperfdata/hsperfdata" parses the file.

