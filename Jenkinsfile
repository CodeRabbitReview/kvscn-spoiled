node {
    checkout scm
    stage('Build') {
        sh go version
    }
    stage('Test') {
        sh gradle test
    }
    stage('Deploy') {
        sh go version
    }
}