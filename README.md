DependencyConfusion is tool used for finding any library used by the project that might be vulnerable to dependency confusion attack. 
	
	Project with following lagnuages supported:
	- Golang
	- python (still in progress)
	- c/c++ (still in progress)

	Flags:
		-u, --url  provide github go.mod raw rul
		-v, --verbose  Print verbose logs to stderr.

sample usage:

go run main.go -u URL_HERE
