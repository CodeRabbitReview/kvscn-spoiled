import java.lang.Float.parseFloat

val binaryDir = "bin"
val binaryName = "storage_server"
val kvServerDockerName = "miprokop/storage_server:v2"
val port = "8080"
val k8sTemplatesPath = "templates/"
val chartVersion = "0.1.0"
val operatorDockerName = "miprokop/crd:v1"
val localBIN = ""
val controllerToolVersion = "v0.8.0"

description = "Storage gradle"
version = "1.0.0"

tasks.register("operatorDeploy") {
    group = "deploy"
    description = ""
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf("kubectl", "apply", "-f", k8sTemplatesPath)
        }
    }
}

tasks.register("operatorUndeploy") {
    group = "deploy"
    description = ""
    doLast {
        exec {
            exec {
                workingDir("operator")
                commandLine = listOf("kubectl", "delete", "-f", k8sTemplatesPath)
            }
        }
    }
}

tasks.register("controllerGen") {
    group = "controller"
    description = "create controller gen binary file into bin dir"
    doLast {
        exec {
            workingDir("operator")
            environment("GOBIN", "${projectDir.absolutePath}/operator/bin")
            commandLine = listOf("go", "install", "sigs.k8s.io/controller-tools/cmd/controller-gen@${controllerToolVersion}")
        }
    }
}

tasks.register("manifests") {
    group = "controller"
    description = "create controller gen binary file into bin dir"
    val controllerGen = "${projectDir.absolutePath}/operator/bin/controller-gen"
    dependsOn("controllerGen")
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf(controllerGen, "rbac:roleName=manager-role", "crd", "webhook", "paths=./...", "output:crd:artifacts:config=config/crd/bases")
        }
    }
}

tasks.register("format") {
    group = "Clean code"
    description = "Formats project by go fmt action"
    doFirst {
        exec {
            workingDir("server")
            commandLine = listOf("go", "fmt", "./...")
        }
    }
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf("go", "fmt", "./...")
        }
    }
}

tasks.register("serverDeploy") {
    group = "k8s"
    description = "Deploys the key-value storage app on the Kubernetes cluster"
    doFirst {
        exec {
            commandLine = listOf("kubectl", "apply", "-f", "https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml")
        }
    }
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("kubectl", "apply", "-f", k8sTemplatesPath)
        }
    }
}

tasks.register("serverUndeploy") {
    group = "k8s"
    description = "Removes the key-value storage app from the Kubernetes cluster"
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("kubectl", "delete", "-f", k8sTemplatesPath)
        }
    }
}

tasks.register("serverOptimizeDependencies") {
    group = "Clean code"
    description = "Removes unused and download dependencies"
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf("go", "mod", "tidy")
        }
    }
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("go", "mod", "tidy")
        }
    }
}

tasks.register("staticCheck") {
    group ="Clean code"
    description = "Runs staticcheck util"
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf("staticcheck", "-f", "json", "./...")
        }
    }
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("staticcheck", "-f", "json", "./...")
        }
    }
}

tasks.register("cleanCode") {
    group ="Clean code"
    description = "Runs format, optimizeDependencies, staticCheck tasks"
    dependsOn("format", "optimizeDependencies", "staticCheck")
}

tasks.register("serverBuild") {
    group = "build"
    description = "Builds binary of project"
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("go", "build", "-o", "./${binaryDir}/${binaryName}", "-a", ".")
        }
    }
}

tasks.register("serverTest") {
    group = "tests"
    description = "Runs all tests in server dir"
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("go", "test", "--cover", "./...")
        }
    }
}

tasks.register("operatorTest") {
    group = "tests"
    description = "Runs all tests in operator dir"
    val controllerGen = "${projectDir.absolutePath}/operator/bin/controller-gen"
    dependsOn("manifests")
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf(controllerGen, "object:headerFile=hack/boilerplate.go.txt", "paths=./...")
        }
    }
    doLast {
        exec {
            workingDir("operator")
            environment("GOBIN", "${projectDir.absolutePath}/operator/bin")
            commandLine = listOf("go", "install", "sigs.k8s.io/controller-runtime/tools/setup-envtest@latest")
        }
    }
    doLast {
        exec {
            workingDir("operator")
            environment("KUBEBUILDER_ASSETS", "${projectDir.absolutePath}/operator/1.23.5-darwin-amd64")
            commandLine = listOf("go", "test", "--cover", "./...")
        }
    }
}

tasks.register("golint") {
    group = "Clean code"
    description = "Runs go lint util"
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf("golangci-lint", "run", "--timeout=5m", "-c", ".golangci.yml")
        }
    }
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("golangci-lint", "run", "--timeout=5m", "-c", ".golangci.yml")
        }
    }
}

tasks.register("operatorDockerBuild") {
    group = "docker"
    description = "Builds docker operator image by Dockerfile"
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf("docker", "build", "-t", operatorDockerName, ".")
        }
    }
}

tasks.register("serverDockerBuild") {
    group = "docker"
    description = "Builds docker server image by Dockerfile"
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("docker", "build", "-t", kvServerDockerName, ".")
        }
    }
}

tasks.register("serverDockerPush") {
    group = "docker"
    description = "Pushes the key-value docker image to dockerhub"
    doLast {
        exec {
            workingDir("server")
            commandLine = listOf("docker", "push",  kvServerDockerName)
        }
    }
}

tasks.register("operatorDockerPush") {
    group = "docker"
    description = "Pushes the key-value docker image to dockerhub"
    doLast {
        exec {
            workingDir("operator")
            commandLine = listOf("docker", "push",  operatorDockerName)
        }
    }
}

tasks.register("testingDone") {
    group = "tests"
    description = "Runs go tests with coverage report and fail if it does not match the low boundary of 70 percent"
    val out = java.io.ByteArrayOutputStream()
    val ps = java.io.PrintStream(out)
    val old = System.out
    val minPercentage = 70
    System.setOut(ps)
    dependsOn("serverTest")
    doFirst {
        System.out.flush()
        System.setOut(old)
        logger.info(out.toString())
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

tasks.register("installChart") {
    group = "helm"
    description = "installs chart on Kubernetes cluster"
    doFirst {
        exec {
            commandLine = listOf("kubectl", "apply", "-f", "https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml")
        }
    }
    doLast {
        Thread.sleep(5000)
        exec {
            commandLine = listOf("helm", "install", "kv-bundle", "kv-bundle-${chartVersion}.tgz")
        }
    }
}

tasks.register("uninstallChart") {
    group = "helm"
    description = "uninstalls chart from Kubernetes cluster"
    doLast {
        exec {
            commandLine = listOf("helm", "uninstall", "kv-bundle")
        }
    }
}
