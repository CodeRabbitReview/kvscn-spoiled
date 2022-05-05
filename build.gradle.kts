import java.lang.Float.parseFloat

val binaryDir = "bin"
val binaryName = "storage_server"
val dockerName = "miprokop/storage_server"
val port = "8080"

description = "Storage gradle"
version = "1.0.0"

tasks {
    task("format") {
        group = "Clean code"
        description = "Format project by go fmt action"
        doLast {
            exec {
                commandLine = listOf("go", "fmt", "./...")
            }
        }
    }

    task("optimizeDependencies") {
        group = "Clean code"
        description = "Remove unused and download dependencies"
        doLast {
            exec {
                commandLine = listOf("go", "mod", "tidy")
            }
        }
    }

    task("staticCheck") {
        group ="Clean code"
        description = "Run staticcheck util"
        doLast {
            exec {
                commandLine = listOf("staticcheck", "-f", "json", "./...")
            }
        }
    }

    task("cleanCode") {
        group ="Clean code"
        description = "Run format, optimizeDependencies, staticCheck tasks"
        dependsOn("format", "optimizeDependencies", "staticCheck")
    }

    task("build") {
        group = "build"
        description = "Build binary of project"
        doLast {
            exec {
                commandLine("go", "build", "-o", "./${binaryDir}/${binaryName}", "-a", ".")
            }
        }
    }

    task("test") {
        group = "tests"
        description = "Run all tests in project"
        doLast {
            exec {
                commandLine("go", "test", "--cover", "./...")
            }
        }
    }

    task("golint") {
        group = "Clean code"
        description = "Run go lint util"
        doLast {
           exec {
               commandLine = listOf("golangci-lint", "run", "--timeout=5m", "-c", ".golangci.yml")
           }
        }
    }

    task("dockerBuild") {
        group = "docker"
        description = "Build docker image by Dockerfile"
        doLast {
            exec {
                commandLine = listOf("docker", "build", "-t", dockerName, ".")
            }
        }
    }

    task("push") {
        group = "docker"
        description = "Push the key-value docker image to dockerhub"
        doLast {
            exec {
                commandLine = listOf("docker", "push",  dockerName)
            }
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
        doLast {
            exec {
                commandLine = listOf("bash", "run.sh")
            }
        }
    }
}
