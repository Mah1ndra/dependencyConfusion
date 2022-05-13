DependencyConfusion is tool used for finding any potential library project vulnerable to dependency confusion attack for the given project. 
	
	Project with following lagnuages supported:
	- Golang
	- python (still in progress)
	- c/c++ (still in progress)

	Flags:
		-u, --url  provide github go.mod raw rul
		-v, --verbose  Print verbose logs to stderr.

sample usage:

go run main.go -u URL_HERE
