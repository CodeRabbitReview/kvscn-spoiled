import java.lang.Float.parseFloat

val binaryDir = "bin"
val binaryName = "storage_server"
val dockerName = "miprokop/storage_server"
val port = "8080"

description = "Storage gradle"
version = "1.0.0"

tasks.register("format") {
    group = "Clean code"
    description = "Format project by go fmt action"
    doLast {
        exec {
            commandLine = listOf("go", "fmt", "./...")
        }
    }
}

tasks.register("optimizeDependencies") {
    group = "Clean code"
    description = "Remove unused and download dependencies"
    doLast {
        exec {
            commandLine = listOf("go", "mod", "tidy")
        }
    }
}

tasks.register("staticCheck") {
    group ="Clean code"
    description = "Run staticcheck util"
    doLast {
        exec {
            commandLine = listOf("staticcheck", "-f", "json", "./...")
        }
    }
}

tasks.register("cleanCode") {
    group ="Clean code"
    description = "Run format, optimizeDependencies, staticCheck tasks"
    dependsOn("format", "optimizeDependencies", "staticCheck")
}

tasks.register("build") {
    group = "build"
    description = "Build binary of project"
    doLast {
        exec {
            commandLine = listOf("go", "build", "-o", "./${binaryDir}/${binaryName}", "-a", ".")
        }
    }
}

tasks.register("test") {
    group = "tests"
    description = "Run all tests in project"
    doLast {
        exec {
            commandLine = listOf("go", "test", "--cover", "./...")
        }
    }
}

tasks.register("golint") {
    group = "Clean code"
    description = "Run go lint util"
    doLast {
        exec {
            commandLine = listOf("golangci-lint", "run", "--timeout=5m", "-c", ".golangci.yml")
        }
    }
}

tasks.register("dockerBuild") {
    group = "docker"
    description = "Build docker image by Dockerfile"
    doLast {
        exec {
            commandLine = listOf("docker", "build", "-t", dockerName, ".")
        }
    }
}

tasks.register("push") {
    group = "docker"
    description = "Push the key-value docker image to dockerhub"
    doLast {
        exec {
            commandLine = listOf("docker", "push",  dockerName)
        }
    }
}

tasks.register("testingDone") {
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

tasks.register("sanityCheck") {
    group = "sanity"
    description = "Run key-value server and execute sender.sh script"
    dependsOn("dockerBuild")
    doLast {
        exec {
            commandLine = listOf("bash", "run.sh")
        }
    }
}
