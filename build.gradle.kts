import java.lang.Float.parseFloat

val binaryDir = "bin"
val binaryName = "storage_server"
val dockerName = "miprokop/storage_server"


task<Exec>("format") {
    group = "Clean code"
    description = "Format project by go fmt action"
    commandLine("go", "fmt", "./...")
}

task<Exec>("optimizeDependencies") {
    group = "Clean code"
    description = "Remove and download dependencies"
    doLast {
        commandLine("go", "mod", "tidy")
    }
}

task<Exec>("staticCheck") {
    group ="Clean code"
    description = "Run staticcheck util"
    doLast {
        commandLine("staticcheck", "-f", "json", "./...")
    }
}

task("cleanCode") {
    group ="Clean code"
    description = "Run format, optimizeDependencies, staticCheck tasks"
    dependsOn("format", "optimizeDependencies", "staticCheck")
}

task<Exec>("build") {
    group = "build"
    description = "Build binary of project"
    doLast {
        commandLine("go", "build", "-o", "./${binaryDir}/${binaryName}", "-a", ".")
    }
}

task<Exec>("test") {
    group = "tests"
    description = "Run all tests in project"
    doLast {
        commandLine("go", "test", "--cover", "./...")
    }
}

task<Exec>("golint") {
    group = "Clean code"
    description = "Run go lint util"
    doLast {
        commandLine("golangci-lint", "run", "--timeout=5m", "-c", ".golangci.yml")
    }
}

task<Exec>("dockerBuild") {
    group = "docker"
    description = "Build docker image by Dockerfile"
    commandLine("docker", "build", "-t", dockerName, ".")
}

task<Exec>("push") {
    group = "docker"
    description = "Push the key-value docker image to dockerhub"
    doLast {
        commandLine("docker", "push",  dockerName)
    }
}

task("testingDone") {
    group = "tests"
    description = "Run go tests with coverage report and fail if doesn’t match the low boundary of 75 percent"
    val out = java.io.ByteArrayOutputStream()
    val ps = java.io.PrintStream(out)
    val old = System.out
    val minPercentage = 75
    System.setOut(ps)
    dependsOn("test")
    doFirst {
        System.out.flush()
        System.setOut(old)
    }
    doLast {
        val resp = out.toString().split("\n")
        for (r in resp) {
            if (!r.contains("[no test files]")) {
                val t = r.split("coverage: ").last()
                    .split("%").first()
                if (t.isEmpty()) {
                    continue
                }
                val percentage = parseFloat(t)
                if (percentage < minPercentage) {
                    throw TaskExecutionException(this,
                        Exception(r))
                }
            }
        }
    }
}

task("sanityCheck") {
    group = "sanity"
    description = "Run key-value server and execute sender.sh script"
    dependsOn("runServer", "sendToServer")
    doLast {
        println("Final Task Completed!")
    }
}
