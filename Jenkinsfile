#!/usr/bin/env groovy

node() {
    def root = tool name: 'Go 1.18', type: 'go'
    stage('Preparation') {
        checkout scm
    }
    stage('Build') {
        sh "{$root}/bin/go build"
    }
   stage ('Test'){
       withEnv(["GOPATH=${WORKSPACE}", "PATH+GO=${root}/bin:${WORKSPACE}/bin", "GOBIN=${WORKSPACE}/bin"]){
         sh "go get github.com/golang/lint/golint"

         try{
           sh "golint ."
           sh "${tool} build ./..."
         } catch (err){
           sh "echo static analyis failed.  See report"
         }

         warnings canComputeNew: true, canResolveRelativePaths: true, categoriesPattern: '', consoleParsers: [[parserName: 'Go Vet'], [parserName: 'Go Lint']], defaultEncoding: '', excludePattern: '', healthy: '', includePattern: '', messagesPattern: '', unHealthy: ''
       }
     }
}