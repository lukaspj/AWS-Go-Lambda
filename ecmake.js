//----------------------
// # ECMake build script
// This is a simple build system using JavaScript, it's home-brewed so use at your own risk.
// To install, run `go install github.com/lukasj/ecmake` or download the latest binary from:
// https://github.com/lukaspj/ecmake/releases
//
// Build with `ecmake Build`
// Deploy with `ecmake Deploy`
//
// **Only ECMAScript 5 is supported.**

var sh = require('sh');
var io = require('io');

function BuildAllProjects() {
    io.Walk("cmd/lambdas", function (path, fileinfo, error) {
        if (fileinfo.Name() !== "lambdas" && fileinfo.IsDir()) {
            sh.RunWithV({
                    "GO111MODULE": "on",
                    "GOOS": "linux"
                },
                "go",
                "build",
                "-ldflags",
                "-s -w",
                "-o",
                "out/",
                "BryrupTeater.Backend/cmd/lambdas/" + fileinfo.Name())
        }
    })
}

SetTargets({
    "Build": BuildAllProjects,
    "Deploy": () => {
        BuildAllProjects();
        sh.RunV(
            "terraform",
            "-chdir=terraform",
            "apply"
        )
    }
});