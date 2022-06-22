#!/usr/bin/env groovy
timestamps {
    node {
        def root = tool type: 'go', name: 'Go 1.17.10'

        withEnv(["GOROOT=${root}",
         "PATH+GO=${root}/bin",
         "GO111MODULE=on",
         "CGO_ENABLED=0"]) {
            stage('SetUp') {
                checkout scm
            }
            stage('Build') {
                withGradle {
                    sh './gradlew build'
                }
            }
            stage('Test') {
                try {
                    sh './gradlew test'
                    sh './gradlew testingDone'
                } catch(Exception e) {
                    echo "Error in testing storage project: ${e.toString()}"
                    sh 'cat test.out'
                    throw new Exception('Error in running tests(any test did not finish correctly)')
                } finally {
                    sh 'rm -f test.out'
                }
            }
        }
    }
}
