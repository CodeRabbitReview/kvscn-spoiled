import java.lang.Float.parseFloat

val binaryDir = "bin"
val binaryName = "storage_server"
val dockerName = "miprokop/storage_server:v2"
val port = "8080"
val k8sTemplatesPath = "templates/"
val storageDeployManifestFile = "storage.yaml"
val metricsServerManifestFile = "metrics_server.yaml"
val chartVersion = "0.1.0"

description = "Storage gradle"
version = "1.0.0"

tasks.register("format") {
    group = "Clean code"
    description = "Formats project by go fmt action"
    doLast {
        exec {
            commandLine = listOf("go", "fmt", "./...")
        }
    }
}

tasks.register("deploy") {
    group = "k8s"
    description = "Deploys the key-value storage app on the Kubernetes cluster"
    doFirst {
        exec {
            commandLine = listOf("kubectl", "apply", "-f", "https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml")
        }
    }
    doFirst {
        exec {
            commandLine = listOf("kubectl", "apply", "-f", "${k8sTemplatesPath}${metricsServerManifestFile}")
        }
    }
    doLast {
        exec {
            commandLine = listOf("kubectl", "apply", "-f", "${k8sTemplatesPath}${storageDeployManifestFile}")
        }
    }
}

tasks.register("undeploy") {
    group = "k8s"
    description = "Removes the key-value storage app from the Kubernetes cluster"
    doLast {
        exec {
            commandLine = listOf("kubectl", "delete", "-f", "${k8sTemplatesPath}${storageDeployManifestFile}")
        }
    }
}

tasks.register("optimizeDependencies") {
    group = "Clean code"
    description = "Removes unused and download dependencies"
    doLast {
        exec {
            commandLine = listOf("go", "mod", "tidy")
        }
    }
}

tasks.register("staticCheck") {
    group ="Clean code"
    description = "Runs staticcheck util"
    doLast {
        exec {
            commandLine = listOf("staticcheck", "-f", "json", "./...")
        }
    }
}

tasks.register("cleanCode") {
    group ="Clean code"
    description = "Runs format, optimizeDependencies, staticCheck tasks"
    dependsOn("format", "optimizeDependencies", "staticCheck")
}

tasks.register("build") {
    group = "build"
    description = "Builds binary of project"
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("go", "build", "-o", "./${binaryDir}/${binaryName}", "-a", ".")
        }
    }
}

tasks.register("test") {
    group = "tests"
    description = "Runs all tests in project"
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("go", "test", "--cover", "./...")
        }
    }
//    doLast {
//        exec {
//            workingDir("operator")
//        }
//        exec {
//            workingDir("operator")
//            commandLine = listOf("go", "test", "--cover", "./...")
//        }
//    }
}

tasks.register("golint") {
    group = "Clean code"
    description = "Runs go lint util"
    doLast {
        exec {
            commandLine = listOf("golangci-lint", "run", "--timeout=5m", "-c", ".golangci.yml")
        }
    }
}

tasks.register("dockerBuild") {
    group = "docker"
    description = "Builds docker image by Dockerfile"
    doLast {
        exec {
            commandLine = listOf("docker", "build", "-t", dockerName, ".")
        }
    }
}

tasks.register("push") {
    group = "docker"
    description = "Pushes the key-value docker image to dockerhub"
    doLast {
        exec {
            commandLine = listOf("docker", "push",  dockerName)
        }
    }
}

tasks.register("testingDone") {
    group = "tests"
    description = "Runs go tests with coverage report and fail if doesn’t match the low boundary of 70 percent"
    val out = java.io.ByteArrayOutputStream()
    val ps = java.io.PrintStream(out)
    val old = System.out
    val minPercentage = 70
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
    description = "Runs key-value server and execute sender.sh script"
    dependsOn("dockerBuild")
    doLast {
        exec {
            commandLine = listOf("bash", "run.sh")
        }
    }
}

tasks.register("createChart") {
    group = "helm"
    description = "creates a Helm chart with all needed resources for the key-value storage"
    doLast {
        exec {
            commandLine = listOf("find", ".", "-name", "kv-bundle*.tgz", "-type", "f", "-delete")
        }
    }
    doLast {
        exec {
            workingDir("kv-bundle")
            commandLine = listOf("helm", "dependency", "update")
        }
    }
    doLast {
        exec {
            commandLine = listOf("helm", "package", "kv-bundle")
        }
    }
}

//tasks.register("installChart") {
//    group = "helm"
//    description = "installs chart on Kubernetes cluster"
//    doFirst {
//        exec {
//            commandLine = listOf("kubectl", "apply", "-f", "https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml")
//        }
//    }
//    doLast {
//        Thread.sleep(5000)
//        exec {
//            commandLine = listOf("helm", "install", "kv-bundle", "kv-bundle-${chartVersion}.tgz")
//        }
//    }
//}

tasks.register("uninstallChart") {
    group = "helm"
    description = "uninstalls chart from Kubernetes cluster"
    doLast {
        exec {
            commandLine = listOf("helm", "uninstall", "kv-bundle")
        }
    }
}
